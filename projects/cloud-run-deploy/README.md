# Entrega:

- O código-fonte completo da implementação. ✅

Código inteiro disponível em `main.go`.

- Testes automatizados demonstrando o funcionamento. ✅

Testes disponíveis em `main_test.go`, rode com `go test .`

- Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação. ✅

Dockerfile disponível, rode com `docker build -t cloud-run-deploy:latest .`

- Deploy realizado no Google Cloud Run (free tier) e endereço ativo para ser acessado. ✅

Deploy realizado em: `https://cloudrun-goexpert-joaqu1m-344054204477.us-central1.run.app/`

Com o comando:

```bash
gcloud run deploy cloudrun-goexpert-joaqu1m \
  --image gcr.io/fullcycle-classes/cloudrun-goexpert-joaqu1m \
  --set-env-vars WEATHERAPI_KEY='MY_SECRET_KEY' \
  --platform managed \
  --project fullcycle-classes \
  --region us-central1 \
  --allow-unauthenticated
```

Para testar, rode:

```bash
curl "https://cloudrun-goexpert-joaqu1m-344054204477.us-central1.run.app/weather?cep=01001000"
```
