#!/usr/bin/env bash

set -e
shopt -s extglob

if ! command -v protoc  &> /dev/null; then
  echo "missing protoc"
fi

# don't ask...
# # generate the pb bindings while on macos
if [[ "$(uname)" == "Darwin" ]]; then
  pb_files=()
  while IFS= read -r line; do
      pb_files+=( "$line" )
  done < <(ls -1 ./pkg/protos/!(*server*|*pb*))

  protoc --proto_path=./pkg/protos \
    --go_out=./pkg/types \
    --go_opt=paths=source_relative \
    "${pb_files[@]}"
fi

# generate grpc bindings
protoc --proto_path=./pkg/protos \
  --go-grpc_out=./pkg/servers \
  --go-grpc_opt=paths=source_relative \
  ./pkg/protos/*server.proto
