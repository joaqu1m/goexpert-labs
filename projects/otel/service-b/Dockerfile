FROM golang:1.24-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o service-b ./

FROM alpine:3.18
COPY --from=build /app/service-b /service-b
EXPOSE 8081
CMD ["/service-b"]
