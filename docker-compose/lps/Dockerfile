FROM --platform=linux/amd64 golang:1.23.8@sha256:e54daaadd35ebb90fc1404ecdc6eb7338ae13555f71a71856ad96976ae084e44 AS builder

ARG COMMIT_HASH
ARG COMMIT_TAG
ENV COMMIT_HASH_VALUE=${COMMIT_HASH}
ENV COMMIT_TAG_VALUE=${COMMIT_TAG}

WORKDIR /code

COPY Makefile go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN if [[ "$LPS_STAGE" != "regtest" ]]; \
    then test -n "$COMMIT_TAG_VALUE" && \
    test -n "$COMMIT_HASH_VALUE" &&     \
    make build;                         \
    else make build;                    \
    fi

FROM --platform=linux/amd64 alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

COPY --from=builder /code/build/liquidity-provider-server /usr/local/bin/liquidity-provider-server

ARG HOME="/home/lps"
RUN adduser -u 1000 --home="$HOME" lps -D lps

RUN mkdir -p "$HOME/db"; chown 1000 "$HOME/db"
RUN mkdir -p "$HOME/logs"; chown 1000 "$HOME/logs"

WORKDIR "$HOME"

USER lps

EXPOSE 8080
