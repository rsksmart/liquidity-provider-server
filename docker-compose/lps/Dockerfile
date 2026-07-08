FROM --platform=linux/amd64 node:24-bookworm-slim AS ui-builder

ENV COREPACK_ENABLE_DOWNLOAD_PROMPT=0
ENV PLAYWRIGHT_SKIP_BROWSER_DOWNLOAD=1
WORKDIR /code/ui

COPY ui/package.json ui/pnpm-lock.yaml ./
RUN corepack enable && corepack prepare pnpm@9.15.4 --activate \
	&& pnpm install --frozen-lockfile

COPY ui/ ./
RUN pnpm run build:app

FROM --platform=linux/amd64 golang:1.25.9@sha256:5ab234a9519e05043f4a97a505a59f21dc40eee172d6b17d411863d6bba599bb AS builder

ARG COMMIT_HASH
ARG COMMIT_TAG
ARG LPS_STAGE=regtest
ENV COMMIT_HASH_VALUE=${COMMIT_HASH}
ENV COMMIT_TAG_VALUE=${COMMIT_TAG}

SHELL ["/bin/bash", "-o", "pipefail", "-c"]

WORKDIR /code

COPY Makefile go.mod go.sum ./
COPY cmd ./cmd
COPY pkg ./pkg
COPY internal ./internal
COPY managementnextui/fs.go ./managementnextui/fs.go
COPY --from=ui-builder /code/ui/dist ./managementnextui/dist

RUN if [[ "${LPS_STAGE}" != "regtest" ]]; then \
	test -n "${COMMIT_TAG_VALUE}" && test -n "${COMMIT_HASH_VALUE}"; \
	fi \
	&& make build-server

FROM --platform=linux/amd64 alpine:3.19.1@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b

COPY --from=builder /code/build/liquidity-provider-server /usr/local/bin/liquidity-provider-server

ARG HOME="/home/lps"
RUN adduser -u 1000 --home="$HOME" lps -D lps

RUN mkdir -p "$HOME/db"; chown 1000 "$HOME/db"
RUN mkdir -p "$HOME/logs"; chown 1000 "$HOME/logs"

WORKDIR "$HOME"

USER lps

EXPOSE 8080
