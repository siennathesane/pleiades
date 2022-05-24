#!/usr/bin/env bash

set -e
shopt -s extglob

if ! command -v protoc  &> /dev/null; then
  echo "missing protoc"
fi

pb_files=()
while IFS= read -r line; do
    pb_files+=( "$line" )
done < <(ls -1 ./pkg/protos/!(*server*|*pb*))

# generate the pb bindings
protoc --proto_path=./pkg/protos \
    --go_out=./pkg/types \
    --go_opt=paths=source_relative \
    --go-vtproto_out=paths=source_relative:./pkg/types \
    --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
    --go-vtproto_opt=features=all \
    "${pb_files[@]}"

# generate grpc bindings
protoc --proto_path=./pkg/protos \
  --go-drpc_out=./pkg/servers \
  --go-drpc_opt=paths=source_relative \
  ./pkg/protos/*server.proto
