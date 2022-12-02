#!/bin/bash
set -eux

buf lint

# kvstore
cd databaseapi || exit
buf breaking --against buf.build/anthropos-labs/kvstore
cd .. || exit

# errors
cd errorapi || exit
buf breaking --against buf.build/anthropos-labs/errors
cd .. || exit

# raft
cd raftapi || exit
buf breaking --against buf.build/anthropos-labs/raft
cd .. || exit
