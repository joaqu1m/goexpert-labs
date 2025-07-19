### first time
```bash
mysql -h 127.0.0.1 -P 3306 -u root -proot -e "CREATE TABLE orders.orders (id VARCHAR(36) PRIMARY KEY, price DOUBLE, tax DOUBLE, final_price DOUBLE);"
```

### on every update
```bash
protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto
```
```bash
go run github.com/99designs/gqlgen generate --config=internal/infra/graph/gqlgen.yml
```
```bash
cd wired && wire
```

### testing resources used:

#### http:
modules/20/api/create_order.http

#### grpc:
```bash
evans -r repl
$ package pb
$ service OrderService
$ call CreateOrder
```

#### graphql:
```graphql
mutation createOrder {
    createOrder(input: {
        id: "ccc",
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
