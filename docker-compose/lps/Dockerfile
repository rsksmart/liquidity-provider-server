FROM --platform=linux/amd64 golang:1.21.6-alpine

RUN apk add git gcc musl-dev curl jq

WORKDIR /code

COPY Makefile go.mod go.sum cmd/ internal/ pkg/ ./
COPY go.sum ./
RUN make build

RUN mv ./build/liquidity-provider-server /usr/local/bin/liquidity-provider-server

ARG HOME="/home/lps"
RUN adduser -u 1000 --home="$HOME" lps -D lps

COPY --chown=lps "docker-compose/lps/start.sh" "$HOME/"

RUN mkdir -p "$HOME/db"; chown 1000 "$HOME/db"
RUN mkdir -p "$HOME/logs"; chown 1000 "$HOME/logs"

WORKDIR "$HOME"

USER lps

EXPOSE 8080
