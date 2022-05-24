package pb

//go:generate protoc --go_out=paths=source_relative:. --plugin protoc-gen-go-drpc="${GOBIN}/protoc-gen-go-drpc" --go-drpc_out=paths=source_relative:. --go-vtproto_out=paths=source_relative:. --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" --go-vtproto_opt=features=all config.proto
