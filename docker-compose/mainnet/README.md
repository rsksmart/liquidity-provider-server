# Usage

* `docker-compose-block.yml` must be used on a separate VM where bitcoind is
  installed via docker and the rskj node is installed locally.
* `docker-compose.yml` must be used in a separate VM in order to run LPS +
  MongoDB.
* In both installations you must have a `.env.mainnet` file with appropriate
  variables.
* Requires Docker Compose v2 (`docker compose`); legacy `docker-compose` v1 is not supported.

# Combined upgrade: `mongo:4` standalone → `mongo:8.0` single-node replica set (rs0)

`docker-compose.yml` now pins `mongodb` to `mongo:8.0` **and** runs it as a
single-node replica set named `rs0`, with internal authentication via a keyfile
generated into a docker-managed volume (`mongo_keyfile`).

Production deployments currently running `mongo:4` standalone must perform both
changes in a single maintenance window:

- **Major version 4 → 8.** MongoDB does not allow skipping major versions on an
  existing WiredTiger data directory. The new `mongo:8.0` binary refuses to
  open a `mongo:4` dbpath; an in-place FCV chain (4.4 → 5.0 → 6.0 → 7.0 → 8.0)
  is technically possible but requires five sequential image swaps with
  `setFeatureCompatibilityVersion` between each. We use **dump and restore**
  instead — one image swap, no FCV gymnastics.
- **Standalone → replica set.** Required so the Go driver can open a
  replica-set topology and multi-document transactions used by
  `pegoutMongoRepository.UpdateRetainedQuotes` succeed instead of erroring with
  `Transaction numbers are only allowed on a replica set member or mongos`.

The two transformations land together on a **fresh empty dbpath** seeded by
`mongorestore` from a pre-upgrade `mongodump` archive. The old `mongo:4` volume
is moved aside (not deleted) and serves as the byte-identical rollback.

Two compose services handle the new topology:

| Service               | Image       | Purpose                                                                 |
|-----------------------|-------------|-------------------------------------------------------------------------|
| `mongo-keyfile-init`  | `alpine:3.20` | One-shot. Generates `/etc/mongo/mongo-keyfile` (0400, owned by mongodb) inside the `mongo_keyfile` named volume. Idempotent. |
| `mongo-rs-init`       | `mongo:8.0` | One-shot. Runs `rs.initiate({_id:"rs0", members:[{_id:0, host:"mongo01:27017"}]})` and waits for PRIMARY via `mongosh`. Idempotent. |

`mongodb` is healthy only after `rs.isMaster().ismaster == true` (probed with
`mongosh` in the healthcheck), so `docker compose up --wait` blocks until the
set is PRIMARY before starting `lps`.

## Pre-flight

Before the maintenance window:

- Confirm you can `mongodump` from the production `mongo:4` instance with the
  credentials in `.env.mainnet` (`MONGODB_USER` / `MONGODB_PASSWORD`).
- Make sure there is enough free space on the LPS VM to hold both the moved-aside
  `mongo:4` dbpath and a `mongodump --gzip` archive of the same database.
- Have the previous (pre-this-PR) compose file on hand in case rollback is needed.

## Procedure

Run these commands on the LPS VM, in this order. Replace `<TS>` with a UTC
timestamp captured at the start of the window (e.g. `$(date -u +%Y%m%dT%H%M%SZ)`)
and reuse it throughout — every artifact gets the same suffix so rollback is
unambiguous.

1. **Quiesce writers.** Mongo data must be quiescent before `mongodump`
   starts. `lps` is the only writer.

   ```bash
   docker compose --env-file .env.mainnet stop lps
   docker ps --filter name=lps01 --format '{{.Names}} {{.Status}}'   # expect empty
   ```

2. **Take a `mongodump` archive from the running `mongo:4` instance.** This is
   the portable, version-independent rollback artifact. Keep it on the LPS VM
   filesystem, **not** inside the mongo data volume.

   ```bash
   TS=$(date -u +%Y%m%dT%H%M%SZ); echo "$TS"
   mkdir -p ./backups

   docker compose --env-file .env.mainnet exec mongodb \
     mongodump \
     -u "$MONGODB_USER" -p "$MONGODB_PASSWORD" \
     --authenticationDatabase admin \
     --db flyover \
     --archive=/data/db/flyover-${TS}.archive --gzip

   docker cp mongo01:/data/db/flyover-${TS}.archive ./backups/flyover-${TS}.archive
   sha256sum ./backups/flyover-${TS}.archive          # record this hash
   ```

3. **Stop the old `mongo:4` container and archive the dbpath.** Move
   (do not delete) `${MONGO_HOME:-/mnt/mongo}/logs` aside; that directory is
   the byte-identical fallback rollback if the `mongodump` archive itself turns
   out to be unusable.

   ```bash
   docker compose --env-file .env.mainnet stop mongodb
   docker compose --env-file .env.mainnet rm -f mongodb

   sudo mv "${MONGO_HOME:-/mnt/mongo}/logs" "${MONGO_HOME:-/mnt/mongo}/logs.bak-${TS}"
   ls -ld "${MONGO_HOME:-/mnt/mongo}/logs.bak-"*
   ```

   Additionally, take a tarball of the moved-aside dbpath if your operations
   policy requires off-host backup:

   ```bash
   sudo tar -czf "./backups/mongo-volume-${TS}.tgz" \
     -C "${MONGO_HOME:-/mnt/mongo}" "logs.bak-${TS}"
   ```

4. **Pull the new compose changes and bring up the new mongo stack on an
   empty dbpath.** The new compose file (`mongo:8.0` + `--replSet rs0
   --keyFile`) creates an empty `${MONGO_HOME:-/mnt/mongo}/logs` and lets
   `mongo-keyfile-init` and `mongo-rs-init` finish the topology setup.

   ```bash
   docker compose --env-file .env.mainnet up -d --wait --wait-timeout 300 \
     mongo-keyfile-init mongodb mongo-rs-init
   ```

   What happens:
   - `mongo-keyfile-init` generates `/etc/mongo/mongo-keyfile` in the
     `mongo_keyfile` named volume and exits 0.
   - `mongodb` starts on a fresh empty dbpath, runs the standard `mongo:8.0`
     entrypoint that creates the root user from `MONGO_INITDB_ROOT_*`, then
     re-execs `mongod --replSet rs0 --bind_ip_all --keyFile
     /etc/mongo/mongo-keyfile`.
   - `mongo-rs-init` calls `rs.initiate(...)` via `mongosh` and waits for
     PRIMARY. On reruns (e.g. `docker compose up` after a restart) it detects
     the existing set and only waits for PRIMARY, which makes the step
     idempotent.
   - The `mongodb` healthcheck flips to `healthy` once `ismaster == true`.

5. **Verify the new server.** Confirm both the version jump and the topology.

   ```bash
   docker compose --env-file .env.mainnet exec mongodb \
     mongosh --quiet \
     -u "$MONGODB_USER" -p "$MONGODB_PASSWORD" \
     --authenticationDatabase admin \
     --eval 'JSON.stringify({version: db.version(), fcv: db.adminCommand({getParameter:1, featureCompatibilityVersion:1}).featureCompatibilityVersion.version})'

   docker compose --env-file .env.mainnet exec mongodb \
     mongosh --quiet \
     -u "$MONGODB_USER" -p "$MONGODB_PASSWORD" \
     --authenticationDatabase admin \
     --eval 'JSON.stringify(rs.status().members.map(m => ({name: m.name, stateStr: m.stateStr})))'
   ```

   Expected:

   ```
   {"version":"8.0.x","fcv":{"version":"8.0"}}
   [{"name":"mongo01:27017","stateStr":"PRIMARY"}]
   ```

   The `flyover` database should be empty at this point — restore happens
   next:

   ```bash
   docker compose --env-file .env.mainnet exec mongodb \
     mongosh --quiet \
     -u "$MONGODB_USER" -p "$MONGODB_PASSWORD" \
     --authenticationDatabase admin \
     --eval 'db.getSiblingDB("flyover").getCollectionNames()'   # expect []
   ```

6. **Restore the dump into MongoDB 8.0.** A cross-version warning (4.4 → 8.0)
   is expected and harmless.

   ```bash
   docker cp ./backups/flyover-${TS}.archive mongo01:/data/db/restore.archive

   docker compose --env-file .env.mainnet exec mongodb \
     mongorestore \
     -u "$MONGODB_USER" -p "$MONGODB_PASSWORD" \
     --authenticationDatabase admin \
     --archive=/data/db/restore.archive --gzip --drop

   docker compose --env-file .env.mainnet exec mongodb \
     rm -f /data/db/restore.archive
   ```

   Re-run the count query and compare every collection against your
   pre-upgrade audit:

   ```bash
   docker compose --env-file .env.mainnet exec mongodb \
     mongosh --quiet \
     -u "$MONGODB_USER" -p "$MONGODB_PASSWORD" \
     --authenticationDatabase admin flyover \
     --eval 'db.getCollectionNames().forEach(c => print(c + " " + db.getCollection(c).countDocuments({})));'
   ```

   Every collection must have the same document count as before step 1. If
   any count is missing or wrong, **stop and roll back** — do not start
   `lps`.

7. **Resume LPS.** No env or URI change is required — the Go driver
   auto-discovers the replica set from `isMaster` and opens replica-set
   topology with the existing `MONGODB_HOST` / `MONGODB_PORT` connection
   string.

   ```bash
   docker compose --env-file .env.mainnet up -d --wait lps
   curl -fsS http://localhost:8080/health
   ```

   Expected:

   ```
   {"status":"ok","services":{"db":"ok","rsk":"ok","btc":"ok"}}
   ```

   In the LPS logs, expect:
   - `Connecting to MongoDB`
   - `Database already at version 1, no migrations needed` (the restored
     `schema_migrations` row carries forward)
   - the standard index-creation lines on `depositEvents.tx_hash`,
     `trustedAccounts`, `batchPegOutEvents.transaction_hash`
   - `Connected to MongoDB`

   Drive a small canary pegin and pegout quote through the API before
   declaring the window complete.

## Rollback

If anything between steps 4 and 7 fails irrecoverably, roll back to the
pre-upgrade `mongo:4` standalone state.

```bash
# 1. stop the new 8.0 stack
docker compose --env-file .env.mainnet stop lps mongo-rs-init mongodb
docker compose --env-file .env.mainnet rm -f mongo-rs-init mongodb

# 2. drop the new (now-tainted) data dir and restore the moved-aside 4.4 volume
sudo rm -rf "${MONGO_HOME:-/mnt/mongo}/logs"
sudo mv "${MONGO_HOME:-/mnt/mongo}/logs.bak-${TS}" "${MONGO_HOME:-/mnt/mongo}/logs"

# 3. revert docker-compose.yml to the previous version (mongo:4, no --replSet,
#    no keyfile, no mongo-keyfile-init / mongo-rs-init services). Then:
docker compose --env-file .env.mainnet up -d --wait mongodb
docker compose --env-file .env.mainnet up -d --wait lps
```

Second fallback — if the moved-aside `logs.bak-<TS>` directory is also
unusable, restore the `mongodump` archive into a fresh `mongo:4` instance:

```bash
docker run --rm -d --name mongo-rollback \
  -e MONGO_INITDB_ROOT_USERNAME=$MONGODB_USER \
  -e MONGO_INITDB_ROOT_PASSWORD=$MONGODB_PASSWORD \
  -v "${MONGO_HOME:-/mnt/mongo}/logs:/data/db" \
  mongo:4
docker cp ./backups/flyover-${TS}.archive mongo-rollback:/tmp/restore.archive
docker exec mongo-rollback mongorestore \
  -u "$MONGODB_USER" -p "$MONGODB_PASSWORD" --authenticationDatabase admin \
  --archive=/tmp/restore.archive --gzip
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
- **Shell.** `mongo:8.0` ships `mongosh` only; the legacy `mongo` shell is
  gone. Every `docker exec`/healthcheck example above uses `mongosh`.
- **Why no in-place FCV chain.** Stepping the existing volume through 4.4 → 5.0
  → 6.0 → 7.0 → 8.0 would also work, but each step requires its own image swap
  and `setFeatureCompatibilityVersion` cycle. Dump and restore is one swap and
  produces a clean WiredTiger data directory written by 8.0 from scratch,
  which is also what the rehearsal in
  `.vscode/my-doc/tasks/FLY-2290/prod-upgrade-regtest.md` validates.
- **Future cluster expansion.** Adding more members later is a non-breaking
  `rs.add("host:27017")` from PRIMARY. Members joining will need the same
  keyfile contents; mount the `mongo_keyfile` volume on each, or distribute
  the same generated keyfile to each node's dbpath.
