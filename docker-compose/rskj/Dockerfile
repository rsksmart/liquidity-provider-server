FROM openjdk:8-jdk-slim-buster@sha256:221a061dbd1a81efe42e6ce9bc099f7ce318b1f7f3ba186e387264f14e61dec2 AS builder

RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends curl git gnupg && \
    apt-get clean

WORKDIR /code/rskj

ARG RSKJ_RELEASE="ARROWHEAD"
ARG RSKJ_VERSION="6.3.1"

RUN gitrev="${RSKJ_RELEASE}-${RSKJ_VERSION}" && \
    git init && \
    git remote add origin https://github.com/rsksmart/rskj.git && \
    git fetch --depth 1 origin tag "$gitrev" && \
    git checkout "$gitrev"

RUN gpg --keyserver https://secchannel.rsk.co/SUPPORT.asc --recv-keys 1DC9157991323D23FD37BAA7A6DBEAC640C5A14B && \
    gpg --verify --output SHA256SUMS SHA256SUMS.asc && \
    sha256sum --check SHA256SUMS && \
    ./configure.sh && \
    ./gradlew --no-daemon clean build -x test && \
    file=rskj-core/src/main/resources/version.properties && \
    version_number=$(sed -n 's/^versionNumber=//p' "$file" | tr -d "\"'") && \
    modifier=$(sed -n 's/^modifier=//p' "$file" | tr -d "\"'") && \
    cp "rskj-core/build/libs/rskj-core-$version_number-$modifier-all.jar" rsk.jar

    FROM --platform=linux/amd64 openjdk:8-jdk@sha256:86e863cc57215cfb181bd319736d0baf625fe8f150577f9eb58bd937f5452cb8 AS runner

ARG UID=1000
ARG HOME="/home/rsk"
RUN useradd -m -u "$UID" --home="$HOME" rsk

COPY --from=builder --chown=root:root --chmod=555 "/code/rskj/rskj-core/build/libs/rskj-core-*-all.jar" "$HOME/rskj-core.jar"
COPY --chown=root:root --chmod=444 rsk.conf $HOME

RUN mkdir -p "$HOME/.rsk"; chown "$UID" "$HOME/.rsk"
RUN mkdir -p "$HOME/logs"; chown -R "$UID" "$HOME/logs"

WORKDIR "$HOME"

USER rsk

HEALTHCHECK --start-period=15s --retries=5 --interval=20s \
    CMD curl -s "http://127.0.0.1:4444" -X POST -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_chainId","params": [],"id":1}' || exit 1
