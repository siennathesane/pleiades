FROM gitpod/workspace-full

RUN brew install \
    mage \
    fzf \
    kubectl \
    helm \
    capnp && \
    go install capnproto.org/go/capnp/v3/capnpc-go@latest

ENV GO111MODULE=off
RUN go get -u capnproto.org/go/capnp/v3/
ENV GO111MODULE=on
