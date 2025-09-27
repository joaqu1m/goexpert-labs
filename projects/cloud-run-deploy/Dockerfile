FROM golang:1.24 AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cloudrun

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/cloudrun .
EXPOSE 8080
ENTRYPOINT ["./cloudrun"]
