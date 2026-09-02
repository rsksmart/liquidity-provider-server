FROM --platform=linux/amd64 eclipse-temurin:17-jdk@sha256:08295ab0f5007a37cbcc6679a8447a7278d9403f9f82acd80ed08cd10921e026 AS builder

RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends curl git gnupg && \
    apt-get clean

WORKDIR /code/powpeg

ARG POWPEG_RELEASE="VETIVER"
ARG POWPEG_VERSION="9.0.3.0"
# Must match the rskj node major.minor (regtest genesis federation type).
ARG RSKJ_CORE_VERSION="9.0.4-VETIVER"
ARG RSKJ_CORE_JAR_SHA256="6190d64ee7cfcf66345b4bd8b70e7d64a97f24995b89952038a148d56998aa4a"
ARG RSKJ_CORE_MODULE_SHA256="7d7fb1870ce397e8a375ce375b31f9448d7af1444cb8899e8d2e6e16a510128c"

RUN gitrev="${POWPEG_RELEASE}-${POWPEG_VERSION}" && \
    git config --global http.version HTTP/1.1 && \
    git init && \
    git remote add origin https://github.com/rsksmart/powpeg-node.git && \
    git fetch --depth 1 origin tag "$gitrev" && \
    git checkout "$gitrev"

RUN gpg --keyserver https://secchannel.rsk.co/SUPPORT.asc --recv-keys A6DBEAC640C5A14B && \
    gpg --finger A6DBEAC640C5A14B && \
    gpg --verify --output SHA256SUMS SHA256SUMS.asc && \
    sha256sum --check SHA256SUMS && \
    sed -i "s/rskjcoreVersion = '9.0.3-VETIVER'/rskjcoreVersion = '${RSKJ_CORE_VERSION}'/" build.gradle && \
    sed -i \
      -e "s/version=\"9.0.3-VETIVER\"/version=\"${RSKJ_CORE_VERSION}\"/" \
      -e "s/rskj-core-9.0.3-VETIVER.jar/rskj-core-${RSKJ_CORE_VERSION}.jar/" \
      -e "s/rskj-core-9.0.3-VETIVER.module/rskj-core-${RSKJ_CORE_VERSION}.module/" \
      -e "s/sha256 value=\"ad5c791752f328a3c6d272e924d67ba78939269e0e68b56f9987997fe8c80fb9\"/sha256 value=\"${RSKJ_CORE_JAR_SHA256}\"/" \
      -e "s/sha256 value=\"576c3964105a1e2e542596e6ddd8c66c06adca81420866cd89f9a455a074c5ac\"/sha256 value=\"${RSKJ_CORE_MODULE_SHA256}\"/" \
      gradle/verification-metadata.xml && \
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
