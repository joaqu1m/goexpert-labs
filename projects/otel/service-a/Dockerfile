FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o service-a ./

FROM alpine:3.18
COPY --from=build /app/service-a /service-a
EXPOSE 8080
CMD ["/service-a"]
