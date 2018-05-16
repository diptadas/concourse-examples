```console
go build main.go
```

```console
docker build -t diptadas/concourse-k8s .
docker push diptadas/concourse-k8
```

```console
docker run -it diptadas/concourse-k8s ./opt/resource/check

{
  "source": {
    "kubeconfig": "/kube/config",
    "namespace": "default",
    "name": "my-config"
  },
  "version": {"resourceVersion": ""}
}
```

```console
docker run -it diptadas/concourse-k8s ./opt/resource/in /tmp/concourse

{
  "source": {
    "kubeconfig": "/kube/config",
    "namespace": "default",
    "name": "my-config"
  },
  "version": {"resourceVersion": "1036"},
  "params": {"filename": "my-config.yaml"}
}
```

```console
docker run -it diptadas/concourse-k8s ./opt/resource/out /tmp/concourse

{
  "source": {
    "kubeconfig": "/kube/config",
    "namespace": "default",
    "name": "my-config"
  },
  "params": {"filename": "my-config.yaml"}
}
```

```console
docker-compose up
```

```console
fly -t tutorial sp -p k8s-test -c pipeline.yaml
fly -t tutorial up -p k8s-test
```