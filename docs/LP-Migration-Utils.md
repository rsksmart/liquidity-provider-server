---
title: LP Migration Utilities
---

# LP Migration Utilities

This document describes the utility scripts that help Liquidity Providers migrate from the **legacy single Liquidity Bridge Contract** to the new 4-contract setup (FlyoverDiscovery, CollateralManagement, Pegin, Pegout).

## Overview

These utilities help Liquidity Providers migrate from the **legacy single LiquidityBridgeContract** by allowing them to:
- Resign from the legacy contract
- Withdraw funds used for pegins

These utilities are specifically designed for the legacy contract and point all contract address fields to the same legacy contract address (since all functionality was in one contract).

### Legacy Contract Addresses

These utilities are designed ONLY for migrating from the legacy monolithic LiquidityBridgeContract. The legacy contract addresses are:

- **Regtest**: Deployed locally via `./lps-env.sh up` (address set automatically)
- **Testnet**: `0xc2a630c053d12d63d32b025082f6ba268db18300`
- **Mainnet**: `0xaa9caf1e3967600578727f975f283446a3da6612`

In the legacy setup, all functionality (collateral management, pegin, pegout, discovery) existed in a single contract. The utilities point all contract address fields to the same legacy contract address.

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
--lbc-address`: optional override for custom LBC contract address
### Default regtest env file

If you want to run the utilities with minimal flags (for example `./utils/resign_utils --resign`), you can provide defaults via an env file:

- Set `LPS_UTILS_ENV_FILE` to a file containing `LPS_STAGE`, `RSK_ENDPOINT`, `SECRET_SRC`, `AWS_LOCAL_ENDPOINT`, and wallet keys such as `WALLET_FILE`, `WALLET_SECRET`, and `PASSWORD_SECRET`.
- If `LPS_UTILS_ENV_FILE` is not set, the utilities will try `docker-compose/local/.env.regtest`, then `regtest.env` at the repo root.
- A minimal `regtest.env` is included for local regtest defaults; adjust it to match your environment.

## Resign utility

Resign from the legacy Liquidity Bridge Contract:

```bash
./utils/resign_utils \
  --network testnet \
  --rsk-endpoint http://localhost:4444 \
  --secret-src env \
  --keystore-file docker-compose/localstack/keystore.json \
  --resign
```

## Withdraw utility

The withdraw utility allows you to withdraw funds used for pegins from the Liquidity Bridge Contract. This utility supports both full and partial withdrawals.

### Withdraw all funds

To withdraw all available funds:

```bash
./utils/withdraw \
  --network testnet \
  --rsk-endpoint http://localhost:4444 \
  --secret-src env \
  --keystore-file /path/to/keystore.json \
  --all
```

### Withdraw specific amount

To withdraw a specific amount (in wei):

```bash
./utils/withdraw \
  --network testnet \
  --rsk-endpoint http://localhost:4444 \
  --secret-src env \
  --keystore-file /path/to/keystore.json \
  --amount 1000000000000000000
```

**Note**: You must provide either `--all` or `--amount` flag. The `--amount` value should be specified in wei (1 RBTC = 10^18 wei).

## Testing Guide

Testing can be done on regtest (local) for development, or on testnet/mainnet for actual migration.

### Testing on Testnet or Mainnet

**Legacy Contract Addresses:**
- **Testnet**: `0xc2a630c053d12d63d32b025082f6ba268db18300`
- **Mainnet**: `0xaa9caf1e3967600578727f975f283446a3da6612`

1. **Prerequisites**:
   - RPC access (e.g., `https://public-node.testnet.rsk.co` for testnet)
   - Wallet with RBTC and collateral in the legacy contract

2. **Build utilities**:
   ```bash
   make utils
   ```

3. **Test resignation**:
   ```bash
   ./utils/resign_utils --resign \
     --network testnet \
     --rsk-endpoint https://public-node.testnet.rsk.co \
     --secret-src env \
     --keystore-file /path/to/keystore.json
   # Password: <your-password>
   ```

4. **Withdraw liquidity** (all funds):
   ```bash
   ./utils/withdraw \
     --network testnet \
     --rsk-endpoint https://public-node.testnet.rsk.co \
     --secret-src env \
     --keystore-file /path/to/keystore.json \
     --all
   ```
   
   Or to withdraw a specific amount:
   ```bash
   ./utils/withdraw \
     --network testnet \
     --rsk-endpoint https://public-node.testnet.rsk.co \
     --secret-src env \
     --keystore-file /path/to/keystore.json \
     --amount 1000000000000000000
   ```

### Testing on Regtest

To test legacy contract migration on regtest, start the environment which will deploy the legacy contract for testing:

1. **Deploy the legacy contract**:
   ```bash
   cd docker-compose/local
   export LPS_STAGE=regtest
   ./lps-env.sh up
   ```
   
   This will:
   - Start the local RSK regtest node and Bitcoin node
   - Deploy the LEGACY LiquidityBridgeContract for migration testing
   - **Automatically update** `.env.regtest`with the deployed contract address

2. **The migration utilities are ready to use** - no additional setup needed!
   ```bash
   cd ../..  # Back to repo root
   # .env.regtest is already configured with correct addresses
   ```

3. **Register as a Liquidity Provider** (required before adding collateral):
   ```bash
   # Note: Legacy contract requires minimum 0.5 RBTC for single provider type
   # or 1 RBTC for "both" (pegin+pegout) provider type
   # The --legacy flag is required for RSK regtest compatibility
   cast send $LBC_ADDR \
     "register(string,string,bool,string)" \
     "TestLP" "http://localhost:8080" true "both" \
     --value 1ether \
     --rpc-url http://localhost:4444 \
     --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
     --legacy
   ```

4. **Add collateral to the legacy contract**:
   ```bash
   cast send $LBC_ADDR \
     "addCollateral()" \
     --value 1ether \
     --rpc-url http://localhost:4444 \
     --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
     --legacy
   ```

5. **Test resignation with the migration utilities**:
   ```bash
   # Use the regtest.env file with the deployed contract address
   export LPS_UTILS_ENV_FILE=regtest.env
   ```

6. **Run the resignation utility**:
   ```bash
   ./utils/resign_utils \
     --network regtest \
     --rsk-endpoint http://localhost:4444 \
     --secret-src env \
     --keystore-file docker-compose/localstack/local-key.json \
     --resign
   # Password: test
   ```

7. **Withdraw funds used for pegins**:
   
   Withdraw all funds:
   ```bash
   ./utils/withdraw \
     --network regtest \
     --rsk-endpoint http://localhost:4444 \
     --secret-src env \
     --keystore-file docker-compose/localstack/local-key.json \
     --all
   # Password: test
   ```
   
   Or withdraw a specific amount (in wei):
   ```bash
   ./utils/withdraw \
     --network regtest \
     --rsk-endpoint http://localhost:4444 \
     --secret-src env \
     --keystore-file docker-compose/localstack/local-key.json \
     --amount 500000000000000000
   # Password: test
   ```

**Note**: The legacy contract is a single contract that contains all functionality (collateral management, pegin, pegout, discovery). These migration utilities are designed specifically for migrating FROM the legacy contract.

