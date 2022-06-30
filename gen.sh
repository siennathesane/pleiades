#!/bin/sh

set -e

echo "generating config protocols"
capnp compile \
  -I$GOPATH/src/capnproto.org/go/capnp/std \
  -ogo:pkg \
  protocols/v1/config/*.capnp

echo "generating database protocols"
capnp compile \
  -I$GOPATH/src/capnproto.org/go/capnp/std \
  -ogo:pkg \
  protocols/v1/database/*.capnp
