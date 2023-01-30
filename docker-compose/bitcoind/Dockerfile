FROM --platform=linux/amd64 ubuntu:18.04

ARG BITCOIN_VERSION=0.17.1
ARG BITCOIN_RELEASES_KEY=01EA5486DE18A882D4C2684590C8019E36C2E964
ARG BITCOIN_DOWNLOAD_FILENAME=bitcoin-${BITCOIN_VERSION}-x86_64-linux-gnu.tar.gz

ENV BITCOIN_VERSION=$BITCOIN_VERSION
ENV BITCOIN_RELEASES_KEY=$BITCOIN_RELEASES_KEY
ENV BITCOIN_DOWNLOAD_FILENAME=$BITCOIN_DOWNLOAD_FILENAME

RUN apt-get update && apt-get install -y wget gnupg2

RUN wget https://bitcoin.org/bin/bitcoin-core-$BITCOIN_VERSION/SHA256SUMS.asc \
    && wget https://bitcoin.org/bin/bitcoin-core-$BITCOIN_VERSION/$BITCOIN_DOWNLOAD_FILENAME \
        && wget https://bitcoin.org/laanwj-releases.asc \
        && gpg2 --import laanwj-releases.asc \
        && gpg2 --fingerprint $BITCOIN_RELEASES_KEY \
        && gpg2 --verify SHA256SUMS.asc \
        && grep -o "$(sha256sum $BITCOIN_DOWNLOAD_FILENAME)" SHA256SUMS.asc \
        && tar -xzvf $BITCOIN_DOWNLOAD_FILENAME \
        && cd bitcoin-$BITCOIN_VERSION/bin \
        && mv bitcoin-cli bitcoin-tx bitcoind /usr/local/bin \
        && cd / \
        && rm -Rfv bitcoin-* *.asc

ARG HOME="/home/bitcoin"
RUN useradd -m -u 1000 --home="$HOME" bitcoin

RUN mkdir -p "$HOME/.bitcoin"; chown 1000 "$HOME/.bitcoin"

WORKDIR "$HOME"

USER bitcoin
