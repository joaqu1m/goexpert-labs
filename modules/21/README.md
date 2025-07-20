```bash
docker build -t 21:latest -f Dockerfile.prod .
```

```bash
kind create cluster --name=goexpert
```

```bash
kubectl get nodes
```

```bash
kind load docker-image 21:latest --name goexpert
```

```bash
kubectl apply -f k8s/deployment.yaml # rerun manually on every change
```

```bash
kubectl get pods
```

```bash
kubectl apply -f k8s/service.yaml
```

```bash
kubectl get services
```

```bash
kubectl port-forward svc/server 80:80
```

let's test the connection!

```bash
curl localhost
```
