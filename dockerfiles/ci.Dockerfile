FROM golang:latest

ENV MAGE_VERSION "1.13.0"
ADD ci/scripts/install-mage.sh .
RUN ./install-mage.sh && \
    rm LICENSE install-mage.sh mage_"${MAGE_VERSION}"_Linux-64bit.tar.gz
