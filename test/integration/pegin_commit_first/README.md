# Commit-first peg-in integration tests

These tests drive a running local LPS. They do not start LPS themselves. They pay Bitcoin, call `registerAddress` as the deployer (a gas sponsor; Watchtower is a different server), then watch Mongo, the `PegInRequested` event, the registration root, and the user Rootstock balance.

Watchtower, the SDK, `resolvePegIn`, and the quote-protocol suite in `pegin_test` are not part of this suite.

## What the tests need

1. An LBC image whose `DeployFlyover` script deploys `FlyoverConfigurations`, calls `setPegInDependencies`, and logs `FlyoverConfigurations proxy:` and `PauseRegistry proxy:`.
2. The local LPS Docker stack, including the Rootstock Bridge (`CREATE_POWPEG=true`).
3. `test/integration/integration-test.config.json` filled with the addresses written to `.env.regtest` after deploy: `peginContract`, `peginAddressRegistry`, `flyoverConfigurations`, `pauseRegistry` (`PEGIN_CONTRACT_ADDRESS`, `PEGIN_ADDRESS_REGISTRY_ADDRESS`, `FLYOVER_CONFIGURATIONS_ADDRESS`, `PAUSE_REGISTRY_ADDRESS`).

## Start the stack

From the liquidity-bridge-contract repository that matches this LPS checkout:

```bash
docker build -t lbc:local .
```

From this repository:

```bash
cd docker-compose/local
```

If the Flyover ABI or LBC image changed, wipe `volumes` and write a new `.env.regtest` before you start. Do not reuse the old chain. Do not run `./lps-local.sh --reset` alone: it recopies `sample-config.env` and pins a GHCR digest.

If `.env.regtest` does not exist, copy `sample-config.env` and set:

```
LBC_IMAGE=lbc:local
LBC_PULL_POLICY=never
CREATE_POWPEG=true
```

```bash
./lps-local.sh
```

Copy `test/integration/integration-test.config.example.json` to `test/integration/integration-test.config.json`. Fill `peginContract`, `peginAddressRegistry`, `flyoverConfigurations`, and `pauseRegistry` from the new `.env.regtest`.

## Run the tests

From the LPS repository root:

```bash
go test -count=1 -timeout 15m ./test/integration/pegin_commit_first
```

The tests run in one suite, one after another. They share the LPS wallet and Mongo. `TestFirstDepositRegisterImportClaim` must pass before `TestSecondDepositSameWatch`.

The suite also proves: a deposit below Flyover `minAmount` registers and does not claim; hard pause blocks claim and claim resumes after unpause; a pay with no `registerAddress` never creates a `peginWatch`.

The tests wait until the Rootstock Bridge Bitcoin height includes the deposit block before `registerAddress` or a claim. That is the same wait `fed-migrator` uses (`primeBtcRelay`): mine a few Rootstock blocks, then read `getBtcBlockchainBestChainHeight`.

Local compose sets the discovery and claim watcher intervals to 2 seconds. Each Mongo/balance wait polls every 500 ms and fails after 20 seconds. The Bridge wait fails after 3 minutes.

## Host endpoints

The test process runs on the host, not inside Docker:

- Bitcoin RPC: `127.0.0.1:5555` (user/password `test`/`test`)
- Rootstock HTTP: `http://localhost:4444`
- Mongo: `mongodb://root:root@127.0.0.1:27017` with `directConnection=true`

Do not use Docker DNS names (`bitcoind`, `mongodb`, `mongo01`) from the host.
