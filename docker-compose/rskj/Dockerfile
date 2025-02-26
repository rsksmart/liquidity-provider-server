FROM eclipse-temurin:17-jdk@sha256:0613a19436dc8f745914b25235d43f3b0eddb8d432d19edce30ffaf2d2f95403 AS builder

RUN apt-get update -y && \
    apt-get install -y -qq --no-install-recommends curl git gnupg && \
    apt-get clean

WORKDIR /code/rskj

ARG RSKJ_RELEASE="LOVELL"
ARG RSKJ_VERSION="7.0.0"

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

FROM --platform=linux/amd64 eclipse-temurin:17-jre@sha256:38e0afc86a10bf4cadbf1586fb617b3a9a4d09c9a0be882e29ada4ed0895fc84

ARG UID=1001
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
