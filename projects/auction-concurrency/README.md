# Auction Concurrency - Leilões com Fechamento Automático

Sistema de leilões que **fecha automaticamente** em 20 segundos usando **goroutines**.

## Como Executar

```bash
docker-compose up --build
```

A API fica disponível em `http://localhost:8085`

## Teste Simples

Execute os comandos abaixo para ver o fechamento automático funcionando:

```bash
# 1. Criar um leilão
curl -X POST http://localhost:8085/auction \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "iPhone 15",
    "category": "Electronics", 
    "description": "iPhone 15 usado em bom estado",
    "condition": 2
  }'

# 2. Ver leilão ativo (status = 0)
curl -s "http://localhost:8085/auction?status=0" | jq '.'

# 3. Aguardar 25 segundos
sleep 25

# 4. Ver que foi fechado automaticamente (status = 1)  
curl -s "http://localhost:8085/auction?status=1" | jq '.'
```

**Resultado:** O leilão muda de `status: 0` para `status: 1` automaticamente após 20 segundos.

---
