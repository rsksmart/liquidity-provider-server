FROM --platform=linux/amd64 golang:1.23.8@sha256:e54daaadd35ebb90fc1404ecdc6eb7338ae13555f71a71856ad96976ae084e44 AS builder

WORKDIR /code

COPY Makefile go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN make utils

FROM scratch
COPY --from=builder /code/utils /
