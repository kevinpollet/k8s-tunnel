# k8s-tunnel

This repository contains an example of an mTLS tunnel between two Kubernetes clusters. Each cluster contains a `gateway`
allowing to setup a secured mTLS tunnel.

The `East` cluster contains a whoami service. The `West` cluster contains a shadow service allowing to access the whoami
service in the `East` cluster through the secured tunnel.

### Requirements

- k3d [v5.0.0](https://github.com/rancher/k3d/releases/tag/v5.0.0)

### 1. Start the East cluster

```bash
$ cd east
$ ./start.sh
```

### 2. Start the West cluster

```bash
$ cd west
$ ./start.sh
```

### 3. Send a request to the Whoami service

```bash
$ kubectl exec deployment/client -- curl http://whoami.default.svc
```

### License

[MIT](./LICENSE.md)
