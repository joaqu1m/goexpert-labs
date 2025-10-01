FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o stress-test .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/stress-test .

ENTRYPOINT ["./stress-test"]