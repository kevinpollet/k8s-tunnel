#!/usr/bin/env bash

k3d cluster delete west
k3d cluster create west --agents 1 --k3s-server-arg '--no-deploy=traefik'

docker build -f gateway/Dockerfile gateway/ -t gateway:west
k3d image import -c west gateway:west

kubectl create secret generic gateway-certs --from-file=tls.ca=../ca/minica.pem --from-file=tls.crt=certs/cert.pem --from-file=tls.key=certs/key.pem
kubectl apply -f gateway.yml -f shadow-svc.yml -f client.yml