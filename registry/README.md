# Registry

## Etcd

```
$ docker run \
    -p 2379:2379 \
    -p 2380:2380 \
    --name etcd-gcr-v3.4.0 \
    --rm -d \
    quay.io/coreos/etcd:v3.4.0

```