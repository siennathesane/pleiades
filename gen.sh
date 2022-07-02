#!/bin/sh

#
# Copyright (c) 2022 Sienna Lloyd
#
# Licensed under the PolyForm Strict License 1.0.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License here:
#  https://github.com/mxplusb/pleiades/blob/mainline/LICENSE
#

set -e

echo "generating config protocols"
capnp compile \
  -I$GOPATH/src/capnproto.org/go/capnp/std \
  -ogo:pkg \
  protocols/v1/host/*.capnp

echo "generating database protocols"
capnp compile \
  -I$GOPATH/src/capnproto.org/go/capnp/std \
  -ogo:pkg \
  protocols/v1/database/*.capnp
