# k8s-tunnel

This repository contains an example of an mTLS tunnel between two Kubernetes clusters. 
Each cluster contains a `gateway` allowing to setup a secured mTLS tunnel. 
The `East` cluster contains a whoami service.
The `West` cluster contains a shadow service allowing to access the whoami service in the `East` cluster through the secured tunnel.

## Start the East cluster

```bash
$ cd east
$ ./start.sh
```

## Start the West cluster

```bash
$ cd west
$ ./start.sh
```

## Send a request to the whoami service

```bash
$ kubectl exec deployment/client -- curl http://whoami.default.svc
```