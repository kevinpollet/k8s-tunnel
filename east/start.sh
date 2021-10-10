#!/usr/bin/env bash

k3d cluster delete east
k3d cluster create east --agents 1 --k3s-arg '--no-deploy=traefik@servers:*' -p "3000:3000@loadbalancer"

docker build -f gateway/Dockerfile gateway/ -t gateway:east
k3d image import -c east gateway:east

kubectl create secret generic gateway-certs --from-file=tls.ca=../ca/minica.pem --from-file=tls.crt=certs/cert.pem --from-file=tls.key=certs/key.pem
kubectl apply -f gateway.yml -f whoami.yml