#!/bin/bash

echo "${BUF_API_TOKEN}" | buf registry login --username "${BUF_USER}" --token-stdin

for f in databaseapi errorapi raftapi; do
  pushd $f || exit
  buf push
  popd || exit
done
