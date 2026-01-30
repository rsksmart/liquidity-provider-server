FROM --platform=linux/amd64 golang:1.24.11@sha256:28252a22a79f4995895749d69853d4c8f750167bd94113ae4da4a570ec5ce7b2 AS builder

WORKDIR /code

COPY Makefile go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN make utils

FROM scratch
COPY --from=builder /code/utils /
