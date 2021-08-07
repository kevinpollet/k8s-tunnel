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
		RootCAs:      caPool,
		Certificates: []tls.Certificate{cert},
		ServerName:   "sni.east.io",
	}

	l, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn, tlsConfig)
	}
}

func handleConnection(conn net.Conn, tlsConfig *tls.Config) {
	defer func() { _ = conn.Close() }()

	dstConn, err := tls.Dial("tcp", "host.k3d.internal:3000", tlsConfig)
	if err != nil {
		log.Println(err)
		return
	}
	defer func() { _ = dstConn.Close() }()

	go func() { _, _ = io.Copy(dstConn, conn) }() // Goroutine exits when connection is closed
	_, _ = io.Copy(conn, dstConn)
}
