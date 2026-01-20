---
title: LP Migration Utilities
---

# LP Migration Utilities

This document describes the utility scripts that help Liquidity Providers migrate liquidity to the new contract set.

## Prerequisites

- Access to the Rootstock RPC endpoint for the target network.
- Wallet secrets configured via the same mechanisms used by the LPS (`--secret-src env|aws`).

## Build the utilities

```bash
make utils
```

The binaries will be placed in the `./utils` directory.

## Common options

Both scripts reuse the base input flags from `cmd/utils/scripts`:

- `--network` (required): `regtest`, `testnet`, or `mainnet`
- `--rsk-endpoint` (required): RPC URL (e.g. `http://localhost:4444`)
- `--secret-src` (required): `env` or `aws`
- `--keystore-file`: required when `--secret-src=env`
- `--keystore-secret` and `--password-secret`: required when `--secret-src=aws`
- `--custom-pegin-address`, `--custom-collateral-address`, `--custom-pegout-address`, `--custom-discovery-address`: optional overrides

### Default regtest env file

If you want to run the utilities with minimal flags (for example `./utils/resign_utils --resign`), you can provide defaults via an env file:

- Set `LPS_UTILS_ENV_FILE` to a file containing `LPS_STAGE`, `RSK_ENDPOINT`, `SECRET_SRC`, `AWS_LOCAL_ENDPOINT`, and wallet keys such as `WALLET_FILE`, `WALLET_SECRET`, and `PASSWORD_SECRET`.
- If `LPS_UTILS_ENV_FILE` is not set, the utilities will try `docker-compose/local/.env.regtest`, then `regtest.env` at the repo root.
- A minimal `regtest.env` is included for local regtest defaults; adjust it to match your environment.

## Resign utility

### Resign

```bash
./utils/resign_utils \
  --network testnet \
  --rsk-endpoint http://localhost:4444 \
  --secret-src env \
  --keystore-file /path/to/keystore.json \
  --resign
```

### Withdraw collateral

```bash
./utils/resign_utils \
  --network testnet \
  --rsk-endpoint http://localhost:4444 \
  --secret-src env \
  --keystore-file /path/to/keystore.json \
  --withdraw-collateral
```

Note: `--withdraw-collateral` requires the resignation delay to have elapsed after calling `--resign`.

## Withdraw utility

### Withdraw full balance

```bash
./utils/withdraw \
  --network testnet \
  --rsk-endpoint http://localhost:4444 \
  --secret-src env \
  --keystore-file /path/to/keystore.json \
  --all
```

### Withdraw a specific amount

```bash
./utils/withdraw \
  --network testnet \
  --rsk-endpoint http://localhost:4444 \
  --secret-src env \
  --keystore-file /path/to/keystore.json \
  --amount 1000000000000000000
```

The `--amount` value is expressed in wei.
