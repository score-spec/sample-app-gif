# sample-app-gif

The sample app used for the gif and video embedded in the docs. The purpose is for it to be complex enough to show the core value propositions of Score:

- Convert to more than one deployment format
- Provision ingress and database resources with dynamic credentials
- Launch locally or remotely
- Show evidence that the request was routed correctly and hit the target database

All with as few lines of code as possible. We're going to for simple and short rather than _correct_.

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
