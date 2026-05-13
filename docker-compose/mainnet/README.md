# Usage
* docker-compose-block must be used on a separate vm where bitcoind is intalled via docker and 
rskj node is installed locally
* docker-compose must be used in a separate VM in order to run LPS + mongodb 
* In both intallations you must have a .env.mainnet file with appropiate variables. 

## MongoDB 8.0 upgrade (recommended procedure)

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

