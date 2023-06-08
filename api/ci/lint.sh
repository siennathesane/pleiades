#!/bin/bash
set -eux

buf lint

# kvstore
cd databaseapi || exit
buf breaking --against buf.build/mxplusb/pleiades/kvstore
cd .. || exit

# errors
cd errorapi || exit
buf breaking --against buf.build/mxplusb/pleiades/errors
cd .. || exit

# raft
cd raftapi || exit
buf breaking --against buf.build/mxplusb/pleiades/raft
cd .. || exit
