# Stress Test

Sistema CLI em Go para realizar testes de carga em serviços web.

## Como rodar o projeto?

```bash
docker build -t stress-test .
docker run stress-test --url=http://google.com --requests=1000 --concurrency=10
```

## Teste Rápido

```bash
# Teste básico
docker run stress-test --url=http://google.com --requests=100 --concurrency=5

# Teste mais intenso
docker run stress-test --url=https://httpbin.org/status/200 --requests=1000 --concurrency=50
```
