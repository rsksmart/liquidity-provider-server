FROM --platform=linux/amd64 golang:1.16-alpine

RUN apk add git
RUN apk add gcc
RUN apk add musl-dev
RUN apk add curl
RUN apk add jq

WORKDIR /code

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY connectors ./connectors
COPY http ./http
COPY storage ./storage
COPY pegout ./pegout
COPY mongo ./mongo

RUN go build -o /usr/local/bin/liquidity-provider-server

ARG HOME="/home/lps"
RUN adduser -u 1000 --home="$HOME" lps -D lps

COPY --chown=lps "docker-compose/lps/config.json" "$HOME/"
COPY --chown=lps "docker-compose/lps/start.sh" "$HOME/"

RUN mkdir -p "$HOME/db"; chown 1000 "$HOME/db"
RUN mkdir -p "$HOME/logs"; chown 1000 "$HOME/logs"

WORKDIR "$HOME"

USER lps

EXPOSE 8080
