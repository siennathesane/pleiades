FROM homebrew/brew

RUN brew install \
    go \
    helm \
    jq \
    kubectl \
    mage \
    minio-mc \
    vim
