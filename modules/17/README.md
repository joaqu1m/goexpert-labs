```bash
migrate create -ext=sql -dir=sql/migrations -seq init
```

```bash
migrate -path=sql/migrations -database "mysql://root:root@tcp(localhost:3306)/courses" -verbose up
```

```bash
sqlc generate
```
