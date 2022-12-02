#!/bin/bash

echo "${BUF_API_TOKEN}" | buf registry login --username "${BUF_USER}" --token-stdin

for f in databaseapi errorapi raftapi; do
  cd $f || exit
  buf push
  cd .. || exit
done
