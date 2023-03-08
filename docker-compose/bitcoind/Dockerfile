FROM --platform=linux/amd64 ubuntu:18.04

ARG BITCOIN_VERSION=23.1
ARG BITCOIN_DOWNLOAD_FILENAME=bitcoin-${BITCOIN_VERSION}-x86_64-linux-gnu.tar.gz

ENV BITCOIN_VERSION=$BITCOIN_VERSION
ENV BITCOIN_DOWNLOAD_FILENAME=$BITCOIN_DOWNLOAD_FILENAME

RUN apt-get update && apt-get install -y wget gnupg2

RUN wget https://bitcoincore.org/bin/bitcoin-core-$BITCOIN_VERSION/SHA256SUMS.asc \
    && wget https://bitcoincore.org/bin/bitcoin-core-$BITCOIN_VERSION/$BITCOIN_DOWNLOAD_FILENAME \
        && tar -xzvf $BITCOIN_DOWNLOAD_FILENAME \
        && cd bitcoin-$BITCOIN_VERSION/bin \
        && mv bitcoin-cli bitcoin-tx bitcoind /usr/local/bin \
        && cd / \
        && rm -Rfv bitcoin-* *.asc

ARG UID=1000
ARG HOME="/home/bitcoin"
RUN useradd -m -u "$UID" --home="$HOME" bitcoin

RUN mkdir -p "$HOME/.bitcoin"; chown "$UID" "$HOME/.bitcoin"

WORKDIR "$HOME"

USER bitcoin
