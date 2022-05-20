protoc --proto_path=./pkg/protos  --go_out=./pkg/protos --go_opt=paths=source_relative ./pkg/protos/raft.proto

protoc --proto_path=./pkg/protos --go-grpc_out=./pkg/servers --go-grpc_opt=paths=source_relative ./pkg/protos/raft-server.proto