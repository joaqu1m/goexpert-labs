### Como Funciona:
- **Por IP**: Limite de 10 req/s (configurável via `RATE_LIMIT_IP_REQUESTS_PER_SECOND`)
- **Por Token**: Limite de 100 req/s (configurável via `RATE_LIMIT_TOKEN_REQUESTS_PER_SECOND`) 
- **Precedência**: Token sempre sobrepõe IP
- **Bloqueio**: 5 minutos quando limite excedido (configurável via `*_BLOCK_TIME_MINUTES`)
- **Strategy (chamado de Storage)**: Redis com sliding window por segundo

- **Testes automatizados demonstrando a eficácia e a robustez do rate limiter.**

Execute `./demo.sh` para ver demonstração completa funcionando.

### Como rodar o projeto?

```bash
docker-compose up --build
```

API disponível em `http://localhost:8080`

## Teste Rápido

```bash
# 1. Testar rate limiting por IP (limite: 10 req/s)
for i in {1..12}; do
  echo "Request $i: $(curl -s -w "%{http_code}" http://localhost:8080/ | tail -c 3)"
done

# 2. Testar com token (limite: 100 req/s)
curl -H "API_KEY: test-token" http://localhost:8080/

# 3. Ver headers de rate limiting
curl -I http://localhost:8080/ | grep X-Rate
```

**Resultado esperado**: Primeiras 10 requisições retornam `200`, próximas retornam `429` com mensagem "you have reached the maximum number of requests or actions allowed within a certain time frame"
