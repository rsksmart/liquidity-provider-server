ARG GO_VERSION=1.21.6
ARG ALPINE_VERSION=3.19.1

FROM --platform=linux/amd64 golang:${GO_VERSION} AS builder

WORKDIR /code

COPY Makefile go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN make build

FROM --platform=linux/amd64 alpine:${ALPINE_VERSION}

COPY --from=builder /code/build/liquidity-provider-server /usr/local/bin/liquidity-provider-server

ARG HOME="/home/lps"
RUN adduser -u 1000 --home="$HOME" lps -D lps

RUN mkdir -p "$HOME/db"; chown 1000 "$HOME/db"
RUN mkdir -p "$HOME/logs"; chown 1000 "$HOME/logs"

WORKDIR "$HOME"

USER lps

EXPOSE 8080
