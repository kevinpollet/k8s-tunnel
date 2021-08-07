package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var (
	ca   = flag.String("ca", "", "give me a ca")
	cert = flag.String("cert", "", "give me a cert")
	key  = flag.String("key", "", "give me a key")
)

func main() {
	flag.Parse()

	ca, err := os.ReadFile(*ca)
	if err != nil {
		log.Fatal(err)
	}

	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(ca)

	cert, err := tls.LoadX509KeyPair(*cert, *key)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}

	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn, tlsConfig)
	}
}

func handleConnection(conn net.Conn, config *tls.Config) {
	var hello *tls.ClientHelloInfo
	tlsConn := tls.Server(conn, &tls.Config{
		GetConfigForClient: func(argHello *tls.ClientHelloInfo) (*tls.Config, error) {
			hello = argHello
			return config, nil
		},
	})
	defer func() { _ = tlsConn.Close() }()

	err := tlsConn.Handshake()
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("Received SNI: %s\n", hello.ServerName)

	dstConn, err := net.Dial("tcp", "whoami.default.svc:80")
	if err != nil {
		log.Println(err)
		return
	}
	defer func() { _ = dstConn.Close() }()

	go func() { _, _ = io.Copy(dstConn, tlsConn) }() // Goroutine exits when connection is closed
	_, _ = io.Copy(tlsConn, dstConn)
}
