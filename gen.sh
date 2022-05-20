#!/usr/bin/env bash

shopt -s extglob

if ! command -v protoc  &> /dev/null; then
  echo "missing protoc"
fi

# don't ask...
# # generate the pb bindings for macos
if [[ "$(uname)" == "Darwin" ]]; then
  pb_files=()
  IFS=" " read -r -a pb_files <<< $(ls ./pkg/protos/!(*server*|*pb*))
  protoc --proto_path=./pkg/protos \
    --go_out=./pkg/types \
    --go_opt=paths=source_relative \
    "${pb_files[*]}"
fi

# generate grpc bindings
protoc --proto_path=./pkg/protos \
  --go-grpc_out=./pkg/servers \
  --go-grpc_opt=paths=source_relative \
  ./pkg/protos/*server.proto
