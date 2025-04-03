FROM --platform=linux/amd64 golang:1.23.7@sha256:cb45cf739cf6bc9eaeacf75d3cd7c157e7d39b757216d813d8115d026ee32e75 AS builder

WORKDIR /code

COPY Makefile go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN make utils

FROM scratch
COPY --from=builder /code/utils /
