FROM golang:1.16-alpine
RUN apk add git
RUN apk add gcc
RUN apk add musl-dev
WORKDIR /app

COPY config.json ./
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN git clone https://github.com/rsksmart/liquidity-provider-server.git
RUN cd liquidity-provider-server  && go get github.com/rsksmart/liquidity-provider-server/connectors
RUN cd liquidity-provider-server  && go get github.com/rsksmart/liquidity-provider-server/http
RUN cd liquidity-provider-server  && go get github.com/rsksmart/liquidity-provider-server/storage

RUN cd liquidity-provider-server && go build -o /liquidity-provider-server

EXPOSE 8080

CMD [ "/liquidity-provider-server" ]