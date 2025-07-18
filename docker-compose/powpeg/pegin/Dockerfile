FROM --platform=linux/amd64 eclipse-temurin:17-jdk@sha256:08295ab0f5007a37cbcc6679a8447a7278d9403f9f82acd80ed08cd10921e026 AS builder

RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends curl git gnupg && \
    apt-get clean

WORKDIR /code/powpeg

ARG POWPEG_RELEASE="LOVELL"
ARG POWPEG_VERSION="7.1.0.0"

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

FROM --platform=linux/amd64 eclipse-temurin:17-jre@sha256:f1515395c0695910a3ca665e973cc11013d1f50d265e61cb8c9156e999d914b4 AS runner

ARG HOME="/home/powpeg"
RUN useradd -m -u 1001 --home="$HOME" powpeg

COPY --from=builder --chown=root:root --chmod=555 "/code/powpeg/build/libs/federate-node-*-all.jar" "$HOME/powpeg.jar"

COPY --chown=root:root --chmod=444 reg1.key "$HOME/reg1.key"
COPY --chown=root:root --chmod=444 regtest-fed.conf "$HOME/regtest-fed.conf"
COPY --chown=root:root --chmod=444 regtest-fed.conf "/etc/rsk/node.conf"

RUN chmod 400 "$HOME/reg1.key" && chown 1001 "$HOME/reg1.key"

RUN mkdir -p "$HOME/.rsk"; chown 1001 "$HOME/.rsk"
RUN mkdir -p "$HOME/logs"; chown 1001 "$HOME/logs"

WORKDIR "$HOME"

USER powpeg
