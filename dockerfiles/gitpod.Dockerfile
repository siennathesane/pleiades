FROM gitpod/workspace-full

RUN brew install \
    mage \
    fzf \
    kubectl \
    helm \
    capnp && \
    go install capnproto.org/go/capnp/v3/capnpc-go@latest
