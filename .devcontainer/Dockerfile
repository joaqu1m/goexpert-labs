FROM mcr.microsoft.com/devcontainers/base:bullseye

RUN apt update -y && \
    apt upgrade -y && \
    apt install -y \
    protobuf-compiler \
    sqlite3

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
RUN go install github.com/google/wire/cmd/wire@latest
