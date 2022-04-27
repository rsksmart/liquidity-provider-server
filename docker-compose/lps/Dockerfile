FROM golang:1.16-alpine

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

RUN go build -o /liquidity-provider-server

WORKDIR /lps

COPY docker-compose/lps/config.json ./
COPY docker-compose/lps/start.sh /

EXPOSE 8080
