---
title: LP Migration Utilities
---

# LP Migration Utilities

This document describes the utility scripts that help Liquidity Providers migrate from the **legacy single Liquidity Bridge Contract** to the new 4-contract setup (FlyoverDiscovery, CollateralManagement, Pegin, Pegout).

## Overview

These utilities help Liquidity Providers migrate from the **legacy single LiquidityBridgeContract** by allowing them to:
- Resign from the legacy contract
- Withdraw collateral after resignation
- Withdraw locked liquidity balances

These utilities are specifically designed for the legacy contract and point all contract address fields to the same legacy contract address (since all functionality was in one contract).

### Legacy Contract Addresses

These utilities are designed ONLY for migrating from the legacy monolithic LiquidityBridgeContract. The legacy contract addresses are:

- **Regtest**: Deployed locally via `./lps-env.sh up` (address set automatically)
- **Testnet**: `0xc2a630c053d12d63d32b025082f6ba268db18300`
- **Mainnet**: `0xaa9caf1e3967600578727f975f283446a3da6612`

In the legacy setup, all functionality (collateral management, pegin, pegout, discovery) existed in a single contract. The utilities point all contract address fields to the same legacy contract address.

## Prerequisites

- Access to the Rootstock RPC endpoint for the target network
- Wallet secrets configured via the same mechanisms used by the LPS (`--secret-src env|aws`)
- [Foundry](https://book.getfoundry.sh/getting-started/installation) installed (for `cast` commands used in testing)

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

4. **Wait for resignation delay** (check contract for delay in blocks)

5. **Withdraw collateral**:
   ```bash
   ./utils/resign_utils --withdraw-collateral \
     --network testnet \
     --rsk-endpoint https://public-node.testnet.rsk.co \
     --secret-src env \
     --keystore-file /path/to/keystore.json
   ```

6. **Withdraw liquidity**:
   ```bash
   ./utils/withdraw --all \
     --network testnet \
     --rsk-endpoint https://public-node.testnet.rsk.co \
     --secret-src env \
     --keystore-file /path/to/keystore.json
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
   - **Automatically update** `.env.regtest` and `regtest-legacy.env` with the deployed contract address

2. **The migration utilities are ready to use** - no additional setup needed!
   ```bash
   cd ../..  # Back to repo root
   # regtest-legacy.env is already configured with correct addresses
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
   # The regtest-legacy.env file was automatically updated during deployment
   export LPS_UTILS_ENV_FILE=regtest-legacy.env
   ```

6. **Run the resignation utility**:
   ```bash
   ./utils/resign_utils --resign
   # Password: test
   ```

7. **Wait for resignation delay** - Mine blocks to skip the delay on regtest:
   ```bash
   # Check the resignation delay
   DELAY=$(cast to-dec $(cast call $LBC_ADDR "getResignDelayBlocks()" --rpc-url http://localhost:4444))
   echo "Resignation delay: $DELAY blocks"
   
   # Mine the required number of blocks
   for i in $(seq 1 $DELAY); do
     cast rpc evm_mine --rpc-url http://localhost:4444 > /dev/null
   done
   
   echo "Mined $DELAY blocks"
   ```

8. **Withdraw collateral**:
   ```bash
   ./utils/resign_utils --withdraw-collateral
   # Password: test
   ```

9. **Withdraw liquidity** (if you added any PegIn balance):
   ```bash
   ./utils/withdraw --all
   ```

**Note**: The legacy contract is a single contract that contains all functionality (collateral management, pegin, pegout, discovery). These migration utilities are designed specifically for migrating FROM the legacy contract.

