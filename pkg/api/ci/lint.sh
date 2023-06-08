#!/bin/bash
#
# Copyright (c) 2023 Sienna Lloyd
#
# Licensed under the PolyForm Internal Use License 1.0.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License here:
#  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
#

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
