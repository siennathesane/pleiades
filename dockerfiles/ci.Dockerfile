FROM homebrew/brew

RUN brew install \
    go \
    helm \
    jq \
    kubectl \
    minio-mc \
    vim
