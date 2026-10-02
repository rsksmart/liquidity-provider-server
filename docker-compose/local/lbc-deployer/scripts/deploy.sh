#!/bin/bash

DEPLOYER_PRIVATE_KEY=$(cast wallet derive-private-key "$DEPLOYER_MNEMONIC")
export DEV_SIGNER_PRIVATE_KEY=$DEPLOYER_PRIVATE_KEY

# LBC v3 PegInContract calls external functions of the Quotes, SignatureValidator, and BtcUtils
# libraries, so its compiled bytecode keeps library address placeholders. The OpenZeppelin
# Upgrades helper in DeployFlyover reads that bytecode with vm.getCode, which fails on unlinked
# bytecode ("vm.getCode: no bytecode for contract; is it abstract or unlinked?"). Deploy the
# libraries first and link them with --libraries.
LIB_OUTPUT=$(forge script script/deployment/DeployLibraries.s.sol:DeployLibraries \
    --rpc-url  "$RSK_ENDPOINT" \
    --private-key "$DEPLOYER_PRIVATE_KEY" \
    --broadcast \
    --legacy \
    --slow 2>&1) || {
      echo "Foundry library deployment failed:"
      echo "$LIB_OUTPUT"
      exit 1
    }

echo "$LIB_OUTPUT"

QUOTES=$(echo "$LIB_OUTPUT" | grep -o 'Quotes: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
SIGNATURE_VALIDATOR=$(echo "$LIB_OUTPUT" | grep -o 'SignatureValidator: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
BTC_UTILS=$(echo "$LIB_OUTPUT" | grep -o 'BtcUtils: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)

if [ -z "$QUOTES" ] || [ -z "$SIGNATURE_VALIDATOR" ] || [ -z "$BTC_UTILS" ]; then
    echo "ERROR: Failed to parse library addresses from deployment output"
    exit 1
fi

echo "Linking libraries Quotes=$QUOTES SignatureValidator=$SIGNATURE_VALIDATOR BtcUtils=$BTC_UTILS"

DEPLOY_OUTPUT=$(forge script script/deployment/DeployFlyover.s.sol:DeployFlyover \
    --rpc-url  "$RSK_ENDPOINT" \
    --private-key "$DEPLOYER_PRIVATE_KEY" \
    --broadcast \
    --legacy \
    --slow \
    --libraries "src/libraries/Quotes.sol:Quotes:${QUOTES}" \
    --libraries "src/libraries/SignatureValidator.sol:SignatureValidator:${SIGNATURE_VALIDATOR}" \
    --libraries "node_modules/@rsksmart/btc-transaction-solidity-helper/contracts/BtcUtils.sol:BtcUtils:${BTC_UTILS}" \
    2>&1) || {
      echo "Foundry deployment failed:"
      echo "$DEPLOY_OUTPUT"
      exit 1
    }

echo "$DEPLOY_OUTPUT"

COLLATERAL_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'CollateralManagement proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
DISCOVERY_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'FlyoverDiscovery proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
PEGIN_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'PegInContract proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
PEGOUT_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'PegOutContract proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
PEGIN_ADDRESS_REGISTRY_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'PegInAddressRegistry proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
FLYOVER_CONFIGURATIONS_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'FlyoverConfigurations proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
PAUSE_REGISTRY_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'PauseRegistry proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)

if [ -z "$COLLATERAL_PROXY" ] || [ -z "$DISCOVERY_PROXY" ] || [ -z "$PEGIN_PROXY" ] || [ -z "$PEGOUT_PROXY" ] || [ -z "$PEGIN_ADDRESS_REGISTRY_PROXY" ] || [ -z "$FLYOVER_CONFIGURATIONS_PROXY" ] || [ -z "$PAUSE_REGISTRY_PROXY" ]; then
    echo "ERROR: Failed to parse contract addresses from deployment output"
    exit 1
fi


echo ""
echo "Verifying deployed contracts..."
for CONTRACT_VAR in PEGIN_PROXY PEGOUT_PROXY COLLATERAL_PROXY DISCOVERY_PROXY PEGIN_ADDRESS_REGISTRY_PROXY FLYOVER_CONFIGURATIONS_PROXY PAUSE_REGISTRY_PROXY; do
  CONTRACT_ADDR=$(eval echo "\$$CONTRACT_VAR")
  CODE=$(curl -s -X POST "$RSK_ENDPOINT" -H "Content-Type: application/json" \
    -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getCode\",\"params\": [\"$CONTRACT_ADDR\",\"latest\"],\"id\":1}" | jq -r ".result")

  if [ "$CODE" = "0x" ] || [ -z "$CODE" ]; then
    echo "  ✗ $CONTRACT_VAR ($CONTRACT_ADDR) - NO CODE (deployment may have failed)"
    exit 1
  else
    echo "  ✓ $CONTRACT_VAR ($CONTRACT_ADDR) - verified"
  fi
done

# Update .env.regtest with new addresses
echo ""
echo "Updating $ENV_FILE with deployed addresses..."
temp_env_file=$(mktemp)
grep -vE "^(PEGIN_CONTRACT_ADDRESS|PEGOUT_CONTRACT_ADDRESS|COLLATERAL_MANAGEMENT_ADDRESS|DISCOVERY_ADDRESS|PEGIN_ADDRESS_REGISTRY_ADDRESS|FLYOVER_CONFIGURATIONS_ADDRESS|PAUSE_REGISTRY_ADDRESS)=" /"$ENV_FILE" > "$temp_env_file"
{
  echo "PEGIN_CONTRACT_ADDRESS=$PEGIN_PROXY"
  echo "PEGOUT_CONTRACT_ADDRESS=$PEGOUT_PROXY"
  echo "COLLATERAL_MANAGEMENT_ADDRESS=$COLLATERAL_PROXY"
  echo "DISCOVERY_ADDRESS=$DISCOVERY_PROXY"
  echo "PEGIN_ADDRESS_REGISTRY_ADDRESS=$PEGIN_ADDRESS_REGISTRY_PROXY"
  echo "FLYOVER_CONFIGURATIONS_ADDRESS=$FLYOVER_CONFIGURATIONS_PROXY"
  echo "PAUSE_REGISTRY_ADDRESS=$PAUSE_REGISTRY_PROXY"
} >> "$temp_env_file"
cat "$temp_env_file" > /"$ENV_FILE"
rm "$temp_env_file"
echo "Deployment complete. Updated $ENV_FILE with new contract addresses."
