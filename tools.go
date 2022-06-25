// +build tools

package tools

import (
    _ "capnproto.org/go/capnp/v3/capnpc-go"
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
    _ "storj.io/drpc/cmd/protoc-gen-go-drpc"
)
