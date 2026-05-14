# Usage

* `docker-compose-block.yml` must be used on a separate VM where bitcoind is
  installed via docker and the rskj node is installed locally.
* `docker-compose.yml` must be used in a separate VM in order to run LPS +
  MongoDB.
* In both installations you must have a `.env.mainnet` file with appropriate
  variables.
* Requires Docker Compose v2 (`docker compose`); legacy `docker-compose` v1 is not supported.

# In-place upgrade: standalone MongoDB → single-node replica set (rs0)

`docker-compose.yml` now runs MongoDB as a single-node replica set named `rs0`,
backed by the existing `${MONGO_HOME:-/mnt/mongo}/logs` dbpath. Internal
authentication uses a keyfile generated into a docker-managed volume
(`mongo_keyfile`) by the `mongo-keyfile-init` service; the host filesystem is
not touched.

This is required so the application driver can open a replica-set topology and
multi-document transactions (used by `pegoutMongoRepository.UpdateRetainedQuotes`)
succeed instead of erroring with `Transaction numbers are only allowed on a
replica set member or mongos`.

Two compose services handle setup:

| Service               | Purpose                                                                 |
|-----------------------|-------------------------------------------------------------------------|
| `mongo-keyfile-init`  | One-shot. Generates `/etc/mongo/mongo-keyfile` (0400, owned by mongodb) inside the `mongo_keyfile` named volume. Idempotent. |
| `mongo-rs-init`       | One-shot. Runs `rs.initiate({_id:"rs0", members:[{_id:0, host:"mongo01:27017"}]})` and waits for PRIMARY. Idempotent. |

`mongodb` is healthy only after `rs.isMaster().ismaster == true`, so
`docker compose up --wait` blocks until the set is PRIMARY before starting `lps`.

## Procedure

Run these commands on the LPS VM, in this order.

1. **Stop application writers.** Mongo data must be quiescent before the
   restart. The `lps` service is the only writer.

   ```bash
   docker compose --env-file .env.mainnet stop lps
   ```

2. **Back up the existing dbpath.** Take a cold backup of
   `${MONGO_HOME:-/mnt/mongo}/logs` before flipping any flags. This is the only
   recovery path if anything goes wrong.

   ```bash
   docker compose --env-file .env.mainnet stop mongodb
   tar -czf "mongo-pre-replset-$(date +%F-%H%M).tgz" \
     -C "${MONGO_HOME:-/mnt/mongo}" logs
   ```

3. **Pull the new compose changes and bring up the mongo stack.** The
   `mongo-keyfile-init` and `mongo-rs-init` services are picked up
   automatically; the existing dbpath is reused as-is.

   ```bash
   docker compose --env-file .env.mainnet up -d --wait --wait-timeout 300 \
     mongo-keyfile-init mongodb mongo-rs-init
   ```

   What happens:
   - `mongo-keyfile-init` generates the keyfile (or reuses an existing one) and
     exits 0.
   - `mongodb` restarts with `--replSet rs0 --keyFile /etc/mongo/mongo-keyfile`
     against the same dbpath. Existing collections and the existing root user
     in the `admin` db are preserved.
   - `mongo-rs-init` runs `rs.initiate(...)` (or detects an already-initialized
     set on subsequent runs) and waits for PRIMARY.
   - The `mongodb` healthcheck flips to `healthy` once `ismaster == true`.

4. **Verify the replica set.**

   ```bash
   docker compose --env-file .env.mainnet exec mongodb \
     mongo --quiet \
     -u "$MONGODB_USER" -p "$MONGODB_PASSWORD" \
     --authenticationDatabase admin \
     --eval 'JSON.stringify(rs.status().members.map(m => ({name: m.name, stateStr: m.stateStr})))'
   ```

   Expected:

   ```
   [{"name":"mongo01:27017","stateStr":"PRIMARY"}]
   ```

5. **Resume LPS.** No env or URI change is required — the Go driver
   auto-discovers the replica set from `isMaster` and opens replica-set
   topology with the existing connection string.

   ```bash
   docker compose --env-file .env.mainnet up -d --wait lps
   curl -fsS http://localhost:8080/health
   ```

   Expected:

   ```
   {"status":"ok","services":{"db":"ok","rsk":"ok","btc":"ok"}}
   ```

## Notes

- **Member host name.** `mongo01` is the replica-set member host advertised by
  `isMaster`. The application container resolves it through the
  `net_lps` network. Connecting from outside the docker network with a
  replica-set-aware driver requires either a `127.0.0.1 mongo01` hosts entry or
  `directConnection=true` in the URI; both are workarounds for ad-hoc shell
  access only.
- **The keyfile.** Lives in the `mongo_keyfile` named docker volume, not on
  the host. Generated once on first `mongo-keyfile-init` run and reused
  thereafter. `docker compose down` (without `-v`) keeps it; `docker volume rm`
  or `docker compose down -v` deletes it — after which `mongodb` will fail to
  start until `mongo-keyfile-init` re-runs.
- **Going back.** Removing `--replSet rs0` from `mongodb.command` and restarting
  reverts to standalone mode against the same dbpath. The replica-set
  metadata in `local.system.replset` is ignored when `--replSet` is absent,
  but is left in place; if you also want to wipe it, drop the `local` db
  while running without `--replSet`.
- **Future cluster expansion.** Adding more members later is a non-breaking
  `rs.add("host:27017")` from PRIMARY. Members joining will need the same
  keyfile contents; mount the `mongo_keyfile` volume on each, or distribute
  the same generated keyfile to each node's dbpath.

# MongoDB 8.0 upgrade (recommended procedure)

The `mongodb` service in this compose file is pinned to `mongo:8.0`. MongoDB does
not support skipping major versions on an existing data directory, so any
deployment that still has a `mongo:4` volume on disk must migrate the data
before pulling the new image. The exact mechanics depend on how the production
database is operated (volume mounts, managed service, snapshot tooling), so the
steps below are a reference outline for the operator to adapt:

1. **Quiesce writers.** Stop the `lps` service / container so no new documents
   are written during the migration window.
2. **Dump the database.** Run `mongodump --gzip --archive=<path>.archive`
   against the running `mongo:4` instance and copy the archive to a safe
   location outside the data volume.
3. **Retire the old data directory.** Stop the `mongodb` service and either
   delete the existing data volume or move it aside so the new image starts on
   a clean directory. Keep the moved-aside copy until the upgrade has been
   validated, in addition to the archive from step 2.
4. **Start MongoDB 8.0 on an empty volume.** Bring up the `mongodb` service
   like in this compose file. Confirm `db.version()` reports `8.0.x` and that
   `featureCompatibilityVersion` is `8.0`.
5. **Restore the dump.** Run `mongorestore --gzip --archive=<path>.archive`
   against the new instance. A cross-version warning from `4.4` to `8.0` is
   expected; verify the restored collection counts match the pre-upgrade
   audit.
6. **Resume writers.** Start the `lps` service again and check the application
   health endpoint and logs for a clean MongoDB connection and a successful
   schema migration check.

Roll back by reverting the compose image tag to the previous version, removing
the 8.0 data directory, and restoring either the moved-aside volume (byte
identical) or the `mongodump` archive (portable) into a `mongo:4` instance.
