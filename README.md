# sample-app-gif

The sample app used for the gif and video embedded in the docs. The purpose is for it to be complex enough to show the core value propositions of Score:

- Convert to more than one deployment format
- Provision ingress and database resources with dynamic credentials
- Launch locally or remotely
- Show evidence that the request was routed correctly and hit the target database

All with as few lines of code as possible. We're going to for simple and short rather than _correct_.

This starts a simple server connected to a postgres database, and then on each request returns something like:

```
HTTP/1.1 200 OK
Server: nginx/1.25.4
Date: Thu, 23 May 2024 16:39:34 GMT
Content-Type: text/plain; charset=utf-8
Content-Length: 127
Connection: keep-alive

SQL VERSION: PostgreSQL 16.1 on aarch64-unknown-linux-musl, compiled by gcc (Alpine 13.2.1_git20231014) 13.2.1 20231014, 64-bit%
```

## How to record the sample gif

Use <https://docs.asciinema.org/>.

Preparation:

```
$ rm -rfv .score-compose .score-k8s compose.yaml manifests.yaml
$ docker pull ghcr.io/score-spec/sample-app-gif:main

$ kind create cluster
$ kubectl use-context kind-kind
$ kubectl --context kind-kind apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.0.0/standard-install.yam
$ helm --kube-context kind-kind install ngf oci://ghcr.io/nginxinc/charts/nginx-gateway-fabric --create-namespace -n nginx-gateway --set service.type=ClusterIP
```

Instructions to record:

```
$ score-compose init
$ score-compose generate score.yaml
$ docker compose up -d 
$ curl http://$(score-compose resources get-outputs 'dns.default#sample.dns' --format '{{ .host }}'):8080/ -i
$ docker logs sample-app-gif-sample-main-1
$ docker compose down -v

$ score-k8s init
$ score-k8s generate score.yaml
$ kubectl apply -f manifests.yaml
$ kubectl wait deployments/sample --for=condition=Ready
$ curl http://$(score-compose resources get-outputs 'dns.default#sample.dns' --format '{{ .host }}')
```
