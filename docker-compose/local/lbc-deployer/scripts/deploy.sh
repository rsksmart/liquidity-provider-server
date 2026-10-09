#!/bin/bash

DEPLOYER_PRIVATE_KEY=$(cast wallet derive-private-key "$DEPLOYER_MNEMONIC")
export DEV_SIGNER_PRIVATE_KEY=$DEPLOYER_PRIVATE_KEY

run_forge_script() {
  forge script "$@" \
    --rpc-url  "$RSK_ENDPOINT" \
    --private-key "$DEPLOYER_PRIVATE_KEY" \
    --broadcast \
    --legacy \
    --slow 2>&1
}

parse_address() {
  echo "$1" | grep -o "$2: 0x[a-fA-F0-9]*" | sed 's/.*: //' | head -1
}

# DeployFlyover deploys through OpenZeppelin Upgrades, which needs the libraries already deployed and linked
LIBRARIES_OUTPUT=$(run_forge_script script/deployment/DeployLibraries.s.sol:DeployLibraries) || {
      echo "Foundry library deployment failed:"
      echo "$LIBRARIES_OUTPUT"
      exit 1
    }

echo "$LIBRARIES_OUTPUT"

QUOTES_LIB=$(parse_address "$LIBRARIES_OUTPUT" 'Quotes')
SIGNATURE_VALIDATOR_LIB=$(parse_address "$LIBRARIES_OUTPUT" 'SignatureValidator')
BTC_UTILS_LIB=$(parse_address "$LIBRARIES_OUTPUT" 'BtcUtils')

if [ -z "$QUOTES_LIB" ] || [ -z "$SIGNATURE_VALIDATOR_LIB" ] || [ -z "$BTC_UTILS_LIB" ]; then
    echo "ERROR: Failed to parse library addresses from deployment output"
    exit 1
fi

DEPLOY_OUTPUT=$(run_forge_script script/deployment/DeployFlyover.s.sol:DeployFlyover \
    --libraries "src/libraries/Quotes.sol:Quotes:$QUOTES_LIB" \
    --libraries "src/libraries/SignatureValidator.sol:SignatureValidator:$SIGNATURE_VALIDATOR_LIB" \
    --libraries "node_modules/@rsksmart/btc-transaction-solidity-helper/contracts/BtcUtils.sol:BtcUtils:$BTC_UTILS_LIB") || {
      echo "Foundry deployment failed:"
      echo "$DEPLOY_OUTPUT"
      exit 1
    }

echo "$DEPLOY_OUTPUT"

COLLATERAL_PROXY=$(parse_address "$DEPLOY_OUTPUT" 'CollateralManagement proxy')
DISCOVERY_PROXY=$(parse_address "$DEPLOY_OUTPUT" 'FlyoverDiscovery proxy')
PEGIN_PROXY=$(parse_address "$DEPLOY_OUTPUT" 'PegInContract proxy')
PEGOUT_PROXY=$(parse_address "$DEPLOY_OUTPUT" 'PegOutContract proxy')
PEGIN_ADDRESS_REGISTRY_PROXY=$(parse_address "$DEPLOY_OUTPUT" 'PegInAddressRegistry proxy')
FLYOVER_CONFIGURATIONS_PROXY=$(parse_address "$DEPLOY_OUTPUT" 'FlyoverConfigurations proxy')
PAUSE_REGISTRY_PROXY=$(parse_address "$DEPLOY_OUTPUT" 'PauseRegistry proxy')

CONTRACT_VARS="PEGIN_PROXY PEGOUT_PROXY COLLATERAL_PROXY DISCOVERY_PROXY PEGIN_ADDRESS_REGISTRY_PROXY FLYOVER_CONFIGURATIONS_PROXY PAUSE_REGISTRY_PROXY"

for CONTRACT_VAR in $CONTRACT_VARS; do
  if [ -z "$(eval echo "\$$CONTRACT_VAR")" ]; then
    echo "ERROR: Failed to parse $CONTRACT_VAR from deployment output"
    exit 1
  fi
done


echo ""
echo "Verifying deployed contracts..."
for CONTRACT_VAR in $CONTRACT_VARS; do
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
