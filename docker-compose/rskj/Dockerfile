FROM --platform=linux/amd64 ubuntu:noble@sha256:d22e4fb389065efa4a61bb36416768698ef6d955fe8a7e0cdb3cd6de80fa7eec

# rskj is installed from the Ubuntu package published for ppa:rsksmart/rskj, so
# the image keeps the packaged layout and its maintainer scripts run: the
# configs in /etc/rsk, the jar in /usr/share/rsk, the licence in
# /usr/share/doc/rskj, and "dpkg -l rskj" reporting the installed version.
#
# The package is pinned by checksum rather than by apt version. A PPA index
# holds only the newest build of a package, so "apt-get install rskj=<version>"
# stops resolving as soon as the next release is published, and the build then
# fails with "Version 'x.y.z~noble' for 'rskj' was not found". Launchpad keeps
# serving the file of a superseded build, so the bytes below stay fetchable.
#
# The regtest stack needs this exact rskj version. The powpeg federators embed
# their own rskj-core, and the node and the federators must agree on consensus.
# If they do not, the federators reject the block that commits a federation
# change and stop relaying BTC headers, which stalls federation migration.
# Do not bump this without bumping POWPEG_VERSION in
# local/powpeg/{pegin,pegout}/Dockerfile to a powpeg-node release built on the
# matching rskj-core. The URL and the checksum always move together.
ARG RSKJ_DEB_URL="https://launchpad.net/~rsksmart/+archive/ubuntu/rskj/+files/rskj_9.0.4~noble_all.deb"
ARG RSKJ_DEB_SHA256="342a107fa9c3ebc46ae3e220f38ba8875a3436251d8bd00182d5ec3feb79b511"

# The package creates the rsk user with uid 888. Here the uid comes from the
# compose file so the bind-mounted volumes stay writable by the host user.
ARG UID=1001
ARG HOME="/var/lib/rsk"

ENV DEBIAN_FRONTEND=noninteractive

# The two debconf answers are the ones the package asks for: debian/rskj.preinst
# gates on the licence, and debian/rskj.postinst asks which network to select.
# postinst also runs "systemctl enable rsk.service". Nothing in this image boots
# systemd, so the unit is inert and that call only has to succeed. A stub keeps
# the systemd packages out of the image and is removed in the same layer.
RUN (getent passwd "$UID" && userdel -r "$(getent passwd "$UID" | cut -d: -f1)" || true) && \
    useradd -m -u "$UID" --home="$HOME" -s /sbin/nologin rsk && \
    apt-get update && \
    apt-get install --no-install-recommends -y \
      ca-certificates \
      curl \
      openjdk-17-jre && \
    rm -rf /var/lib/apt/lists/* && \
    curl -fsSL "$RSKJ_DEB_URL" -o /tmp/rskj.deb && \
    echo "$RSKJ_DEB_SHA256  /tmp/rskj.deb" | sha256sum -c - && \
    echo 'rskj shared/accepted-rsk-license-v1-1 boolean true' | debconf-set-selections && \
    echo 'rskj shared/config select regtest' | debconf-set-selections && \
    printf '#!/bin/sh\nexit 0\n' > /usr/local/sbin/systemctl && \
    chmod +x /usr/local/sbin/systemctl && \
    dpkg -i /tmp/rskj.deb && \
    rm -f /tmp/rskj.deb /usr/local/sbin/systemctl

# The package selects a network by symlinking /etc/rsk/node.conf to one of the
# packaged configs. This stack needs regtest settings the packaged regtest.conf
# does not carry: the powpeg peer enodes, autoMine, and the hardfork activation
# heights the federators expect. Repointing the symlink leaves every packaged
# config untouched. Do not COPY onto node.conf, because COPY writes through the
# symlink and would overwrite the packaged regtest.conf.
COPY rsk.conf /etc/rsk/lps.conf
RUN ln -sfn /etc/rsk/lps.conf /etc/rsk/node.conf

WORKDIR "$HOME"

USER rsk

# Mirrors ExecStart in the packaged /lib/systemd/system/rsk.service.
ENTRYPOINT ["/usr/bin/java", "-Dlogback.configurationFile=/etc/rsk/logback.xml", "-cp", "/usr/share/rsk/rsk.jar", "co.rsk.Start"]

HEALTHCHECK --start-period=15s --retries=5 --interval=20s \
    CMD curl -s "http://127.0.0.1:4444" -X POST -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","method":"eth_chainId","params": [],"id":1}' || exit 1
