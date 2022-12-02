from ubuntu:jammy

run apt update -y && \
    apt upgrade -y && \
    apt install -y \
        build-essential \
        jq \
        curl \
        golang \
        git \
        protobuf-compiler

## install node and es
run curl -sL https://deb.nodesource.com/setup_16.x | bash - && \
    apt install -y nodejs && \
    npm install -g @bufbuild/protoc-gen-connect-web @bufbuild/protoc-gen-es

# install go tools
run go env -w GOBIN=/usr/local/bin && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest && \
    go install github.com/bufbuild/connect-go/cmd/protoc-gen-connect-go@latest

# install buf
run curl -sSL \
    "https://github.com/bufbuild/buf/releases/download/v1.8.0/buf-$(uname -s)-$(uname -m)" \
        -o "/usr/local/bin/buf" && \
    chmod +x "/usr/local/bin/buf"
