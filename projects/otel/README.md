# Como rodar (dev)

1. Crie um arquivo `.env` na raiz (`projects/otel`) com sua chave WeatherAPI:

```env
WEATHERAPI_KEY=YOUR_WEATHERAPI_KEY
```

2. Suba os serviços com:

```bash
docker compose up --build
```

3. Envie a sua primeira requisição:

```bash
curl -X POST http://localhost:8085/cep -H "Content-Type: application/json" -d '{"cep": "01001000"}'
```

4. Acesse o Zipkin em `http://localhost:9411/zipkin/` e veja os traces.
