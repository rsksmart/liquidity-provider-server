FROM openjdk:8-jdk-slim-buster AS builder

RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends git curl gnupg && \
    apt-get clean

WORKDIR /code/powpeg

ARG POWPEG_RELEASE="ARROWHEAD"
ARG POWPEG_VERSION="6.1.0.0"

RUN gitrev="${POWPEG_RELEASE}-${POWPEG_VERSION}" && \
    git init && \
    git remote add origin https://github.com/rsksmart/powpeg-node.git && \
    git fetch --depth 1 origin tag "$gitrev" && \
    git checkout "$gitrev"

RUN gpg --keyserver https://secchannel.rsk.co/SUPPORT.asc --recv-keys A6DBEAC640C5A14B && \
    gpg --finger A6DBEAC640C5A14B && \
    gpg --verify --output SHA256SUMS SHA256SUMS.asc && \
    sha256sum --check SHA256SUMS && \
    ./configure.sh && \
    ./gradlew --no-daemon clean build -x test

FROM --platform=linux/amd64 openjdk:8-jdk AS runner

ARG HOME="/home/powpeg"
RUN useradd -m -u 1000 --home="$HOME" powpeg

COPY --from=builder --chown=root:root --chmod=555 "/code/powpeg/build/libs/federate-node-*-all.jar" "$HOME/powpeg.jar"

COPY --chown=root:root --chmod=444 reg2.key "$HOME/reg2.key"
COPY --chown=root:root --chmod=444 regtest-fed.conf "$HOME/regtest-fed.conf"
COPY --chown=root:root --chmod=444 regtest-fed.conf "/etc/rsk/node.conf"

RUN chmod 400 "$HOME/reg2.key" && chown 1000 "$HOME/reg2.key"

RUN mkdir -p "$HOME/.rsk"; chown 1000 "$HOME/.rsk"
RUN mkdir -p "$HOME/logs"; chown 1000 "$HOME/logs"

WORKDIR "$HOME"

USER powpeg