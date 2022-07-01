FROM ubuntu:jammy

# core stuff
RUN apt update -y && \
    apt upgrade -y && \
    apt install -y \
        apt-transport-https \
        ca-certificates \
        curl \
        gnupg    \
        jq \
        vim

# helm
RUN curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | tee /usr/share/keyrings/helm.gpg > /dev/null && \
    apt install apt-transport-https -y && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | tee /etc/apt/sources.list.d/helm-stable-debian.list && \
    apt update && \
    apt install helm

# kubectl
RUN curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg && \
    echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list && \
    apt update && \
    apt install kubectl


# minio client
RUN curl -fsSLo /usr/bin/mc https://dl.min.io/client/mc/release/linux-amd64/mc && \
    chmod +x /usr/bin/mc
