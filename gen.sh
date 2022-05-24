#!/usr/bin/env bash

set -e
shopt -s extglob

if ! command -v protoc  &> /dev/null; then
  echo "missing protoc"
fi

pb_files=()
while IFS= read -r line; do
    pb_files+=( "$line" )
done < <(ls -1 ./protobufs/*.proto)

protoc \
  -I ./protobufs \
  --go_out=./pkg/pb \
  --go_opt=paths=source_relative \
  --plugin protoc-gen-go-drpc="${GOBIN}/protoc-gen-go-drpc" \
  --plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
  --go-drpc_out=./pkg/pb \
  --go-drpc_opt=paths=source_relative \
  --go-vtproto_out=./pkg/pb \
  --go-vtproto_opt=features=all \
   --go-vtproto_opt=paths=source_relative \
  "${pb_files[@]}"
