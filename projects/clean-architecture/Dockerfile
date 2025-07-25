FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o server cmd/orders/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/.env .

CMD ["./server"]
