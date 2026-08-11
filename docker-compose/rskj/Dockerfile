# rsk.jar comes from the official rskj image rather than from ppa:rsksmart/rskj,
# because a PPA publishes only the newest build of a package. The pinned version
# left the index as soon as the next one was released and the image build then
# failed with "Version 'x.y.z~noble' for 'rskj' was not found". A registry digest
# never disappears, so the version below stays buildable.
#
# The regtest stack needs this exact rskj version. The powpeg federators embed
# their own rskj-core, and the node and the federators must agree on consensus.
# If they do not, the federators reject the block that commits a federation
# change and stop relaying BTC headers, which stalls federation migration.
# Do not bump this without bumping POWPEG_VERSION in
# local/powpeg/{pegin,pegout}/Dockerfile to a powpeg-node release built on the
# matching rskj-core.
FROM --platform=linux/amd64 rsksmart/rskj:VETIVER-9.0.3@sha256:eb1a15ed808a5f204d5c2307e8b3fd171b8709f3154a0e1582798bd3f62a0b59 AS rskj

FROM --platform=linux/amd64 ubuntu:noble@sha256:d22e4fb389065efa4a61bb36416768698ef6d955fe8a7e0cdb3cd6de80fa7eec

ARG UID=1001
ARG HOME="/home/rsk"

ENV DEBIAN_FRONTEND=noninteractive

RUN (getent passwd "$UID" && userdel -r "$(getent passwd "$UID" | cut -d: -f1)" || true) && \
    useradd -m -u "$UID" --home="$HOME" rsk && \
    apt-get update && \
    apt-get install --no-install-recommends -y \
      ca-certificates \
      curl \
      openjdk-17-jre-headless && \
    rm -rf /var/lib/apt/lists/* && \
    mkdir -p /usr/share/rsk /var/log/rsk && \
    chown "$UID" /var/log/rsk

COPY --from=rskj /var/lib/rsk/rsk.jar /usr/share/rsk/rsk.jar

# Both files ship in the rskj .deb, which this image no longer installs.
# logback.xml is the packaged one: it writes /var/log/rsk/rsk.log, the path the
# compose file mounts a volume on.
COPY logback.xml /etc/rsk/logback.xml
COPY rsk.conf /etc/rsk/node.conf

WORKDIR "$HOME"

USER rsk

ENTRYPOINT ["/usr/bin/java", "-Dlogback.configurationFile=/etc/rsk/logback.xml", "-cp", "/usr/share/rsk/rsk.jar", "co.rsk.Start"]

HEALTHCHECK --start-period=15s --retries=5 --interval=20s \
    CMD curl -s "http://127.0.0.1:4444" -X POST -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_chainId","params": [],"id":1}' || exit 1
