FROM mcr.microsoft.com/devcontainers/base:bullseye

RUN apt update -y && \
    apt upgrade -y && \
    apt install -y \
    protobuf-compiler \
    sqlite3 \
    default-mysql-client

RUN [ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.29.0/kind-linux-amd64 && \
    chmod +x ./kind && \
    mv ./kind /usr/local/bin/kind

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
