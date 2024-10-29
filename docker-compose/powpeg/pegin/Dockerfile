FROM openjdk:8-jdk-slim-buster@sha256:221a061dbd1a81efe42e6ce9bc099f7ce318b1f7f3ba186e387264f14e61dec2 AS builder

RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends curl git gnupg && \
    apt-get clean

WORKDIR /code/powpeg

ARG POWPEG_RELEASE="ARROWHEAD"
ARG POWPEG_VERSION="6.3.1.0"

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

    FROM --platform=linux/amd64 openjdk:8-jdk@sha256:86e863cc57215cfb181bd319736d0baf625fe8f150577f9eb58bd937f5452cb8 AS runner

ARG HOME="/home/powpeg"
RUN useradd -m -u 1000 --home="$HOME" powpeg

COPY --from=builder --chown=root:root --chmod=555 "/code/powpeg/build/libs/federate-node-*-all.jar" "$HOME/powpeg.jar"

COPY --chown=root:root --chmod=444 reg1.key "$HOME/reg1.key"
COPY --chown=root:root --chmod=444 regtest-fed.conf "$HOME/regtest-fed.conf"
COPY --chown=root:root --chmod=444 regtest-fed.conf "/etc/rsk/node.conf"

RUN chmod 400 "$HOME/reg1.key" && chown 1000 "$HOME/reg1.key"

RUN mkdir -p "$HOME/.powpeg"; chown 1000 "$HOME/.powpeg"
RUN mkdir -p "$HOME/logs"; chown 1000 "$HOME/logs"

WORKDIR "$HOME"

USER powpeg