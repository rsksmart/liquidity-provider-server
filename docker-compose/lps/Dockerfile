ARG GO_VERSION=1.21.12
ARG ALPINE_VERSION=3.19.1

FROM --platform=linux/amd64 golang:${GO_VERSION} AS builder

ARG COMMIT_HASH
ARG COMMIT_TAG
ENV COMMIT_HASH_VALUE=${COMMIT_HASH}
ENV COMMIT_TAG_VALUE=${COMMIT_TAG}

WORKDIR /code

COPY Makefile go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN test -n "$COMMIT_TAG_VALUE" && test -n "$COMMIT_HASH_VALUE" && make build

FROM --platform=linux/amd64 alpine:${ALPINE_VERSION}

COPY --from=builder /code/build/liquidity-provider-server /usr/local/bin/liquidity-provider-server

ARG HOME="/home/lps"
RUN adduser -u 1000 --home="$HOME" lps -D lps

RUN mkdir -p "$HOME/db"; chown 1000 "$HOME/db"
RUN mkdir -p "$HOME/logs"; chown 1000 "$HOME/logs"

WORKDIR "$HOME"

USER lps

EXPOSE 8080
