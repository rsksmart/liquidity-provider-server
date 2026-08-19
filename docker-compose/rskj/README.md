# Running rskj in this stack

This image installs rskj from the Ubuntu package Rootstock publishes, pinned to one
exact file.

## Why a checksum, not a version number

A PPA is a distribution channel, not an archive. It carries only the newest build of a
package, so a version pin stops resolving as soon as the next release lands:

```
E: Version '9.0.3~noble' for 'rskj' was not found
```

A version number names whatever is current. A checksum names one file, permanently. The
Dockerfile fetches the package from its stable Launchpad URL and refuses to build if the
bytes differ:

```dockerfile
ARG RSKJ_DEB_URL="https://launchpad.net/~rsksmart/+archive/ubuntu/rskj/+files/rskj_9.0.3~noble_all.deb"
ARG RSKJ_DEB_SHA256="9373866f46063917d38a75b2d077cd007d4dd1b1d937794d51701a4fad5f83e3"
```

You run the same node on every machine, on every build, at any point in the future. If
the file ever changes, the build stops instead of shipping something you did not test.

## Why the package, not a bare jar

Installing the `.deb` keeps what the maintainer ships: the configs in `/etc/rsk`, the
licence in `/usr/share/doc/rskj`, and the declared `openjdk-17-jre`, which dpkg enforces.
It also registers the node, so `dpkg -l rskj` reports the version you are running and any
inventory or vulnerability scan that reads the package database can see it. A jar copied
into an image appears in no such listing.

## rskj and powpeg move together

The powpeg federators embed their own rskj-core, and the node and the federators must
agree on consensus. When they disagree, the federators reject the block that commits a
federation change and stop relaying BTC headers, which stalls federation migration.

Never bump one without the other.

## Bumping the version

1. Choose the powpeg-node release, and read which rskj-core it embeds.
2. Download that rskj `.deb` from `https://launchpad.net/~rsksmart/+archive/ubuntu/rskj/+files/`.
3. Take its checksum with `sha256sum`.
4. Update `RSKJ_DEB_URL` and `RSKJ_DEB_SHA256` together. They are one pin, not two.
5. Update `POWPEG_VERSION` in `../local/powpeg/{pegin,pegout}/Dockerfile`.
6. Rebuild and run the regtest E2E suite before merging.

## Checking what you are running

```sh
docker exec rskj01 dpkg -l rskj    # installed version
docker exec rskj01 dpkg -V rskj    # files that differ from the package
```

`dpkg -V` reports two entries by design: `regtest.conf` and `rsk.service`, which the
package's own install script edits. Anything else is drift worth investigating.
