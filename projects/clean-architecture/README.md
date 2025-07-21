## Como realizar o setup do desafio?

Comece executando, comando por comando:

```bash
docker-compose up -d
```

```bash
go mod tidy
```

```bash
go run cmd/orders/main.go
```

## Agora, realize as operações de listagem e criação de pedidos pelos seguintes entrypoints:

### HTTP

Acesse a partir da porta `8000`, consulte os endpoints disponíveis no caminho `projects/clean-architecture/api`

### GRPC

Acesse a partir da porta `50051`, no package `pb` e service `OrderService`

### GraphQL

Acesse a partir da porta `8080`, com uma interface disponível em `http://localhost:8080/`. Algumas queries de exemplo são:

```graphql
mutation createOrder {
    createOrder(input: {
        id: "12345",
        Price: 400.0,
        Tax: 150.0
    }) {
        id,
        Price,
        Tax,
        FinalPrice
    }
}
```

```graphql
query listOrders {
    listOrders {
        id
        Price
        Tax
        FinalPrice
    }
}
```

## Mensageria

### Rabbit MQ

Na URL `http://localhost:15672/` duas filas podem ser acessadas, nesses respectivos binds:

---

Exchange: `amq.direct`

Routing key: `order.created`

---

Exchange: `amq.direct`

Routing key: `order.listed`
