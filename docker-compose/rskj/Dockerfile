FROM --platform=linux/amd64 ubuntu:noble@sha256:d22e4fb389065efa4a61bb36416768698ef6d955fe8a7e0cdb3cd6de80fa7eec

ARG RSKJ_VERSION="8.1.0~noble"
ARG UID=1001
ARG HOME="/home/rsk"

ENV DEBIAN_FRONTEND=noninteractive

RUN useradd -m -u "$UID" --home="$HOME" rsk && \
    apt-get update && \
    apt-get install --no-install-recommends -y \
      dbus debconf-utils \
      init-system-helpers \
      software-properties-common \
      systemd systemd-sysv && \
    apt-add-repository -y ppa:rsksmart/rskj && \
    apt-get update && \
    echo 'rskj    shared/accepted-rsk-license-v1-1        boolean true' | debconf-set-selections && \
    echo 'rskj    shared/config   select  regtest' | debconf-set-selections && \
    apt-get install -y rskj=$RSKJ_VERSION && \
    rm -rf /var/lib/apt/lists/*

COPY rsk.conf /etc/rsk/node.conf

WORKDIR "$HOME"

USER rsk

ENTRYPOINT ["/usr/bin/java", "-Dlogback.configurationFile=/etc/rsk/logback.xml", "-cp", "/usr/share/rsk/rsk.jar", "co.rsk.Start"]

HEALTHCHECK --start-period=15s --retries=5 --interval=20s \
    CMD curl -s "http://127.0.0.1:4444" -X POST -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_chainId","params": [],"id":1}' || exit 1
