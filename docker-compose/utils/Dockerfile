FROM --platform=linux/amd64 golang:1.25.9@sha256:5ab234a9519e05043f4a97a505a59f21dc40eee172d6b17d411863d6bba599bb AS builder

WORKDIR /code

COPY Makefile go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal

RUN make utils

FROM scratch
COPY --from=builder /code/utils /
