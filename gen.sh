#!/bin/sh

set -e

echo "generating config protocols"
capnp compile \
  -I$GOPATH/src/capnproto.org/go/capnp/std \
  -ogo:pkg \
  protocols/config/v1/*.capnp
