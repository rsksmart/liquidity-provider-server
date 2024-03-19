FROM --platform=linux/amd64 ubuntu:18.04

ARG BITCOIN_VERSION=23.1
ARG BITCOIN_DOWNLOAD_FILENAME=bitcoin-${BITCOIN_VERSION}-x86_64-linux-gnu.tar.gz

ENV BITCOIN_VERSION=$BITCOIN_VERSION
ENV BITCOIN_DOWNLOAD_FILENAME=$BITCOIN_DOWNLOAD_FILENAME

RUN apt-get update && apt-get install --no-install-recommends -y ca-certificates wget gnupg2 curl && apt-get clean

RUN wget "https://bitcoincore.org/bin/bitcoin-core-$BITCOIN_VERSION/SHA256SUMS.asc" \
    && wget "https://bitcoincore.org/bin/bitcoin-core-$BITCOIN_VERSION/$BITCOIN_DOWNLOAD_FILENAME" \
        && tar -xzvf "$BITCOIN_DOWNLOAD_FILENAME" \
        && cd "bitcoin-$BITCOIN_VERSION/bin" \
        && mv bitcoin-cli bitcoin-tx bitcoind /usr/local/bin \
        && cd / \
        && rm -Rfv ./bitcoin-* ./*.asc

ARG HOME="/home/bitcoin"
RUN useradd -m -u 1000 --home="$HOME" bitcoin

RUN mkdir -p "$HOME/.bitcoin"; chown 1000 "$HOME/.bitcoin"

WORKDIR "$HOME"

USER bitcoin

HEALTHCHECK --start-period=15s --retries=5 --interval=20s \
    CMD curl -s "http://127.0.0.1:5555" -X POST --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" \
        -d '{"jsonrpc": "1.0", "method": "getnetworkinfo", "params": [], "id":"1"}' | grep "\"result\":{" || exit 1
