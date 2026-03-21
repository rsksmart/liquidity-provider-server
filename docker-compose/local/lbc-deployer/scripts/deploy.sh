#!/bin/bash

DEPLOYER_PRIVATE_KEY=$(cast wallet derive-private-key "$DEPLOYER_MNEMONIC")
export DEV_SIGNER_PRIVATE_KEY=$DEPLOYER_PRIVATE_KEY

DEPLOY_OUTPUT=$(forge script script/deployment/DeployFlyover.s.sol:DeployFlyover \
    --rpc-url  "$RSK_ENDPOINT" \
    --private-key "$DEPLOYER_PRIVATE_KEY" \
    --broadcast \
    --legacy \
    --slow 2>&1) || {
      echo "Foundry deployment failed:"
      echo "$DEPLOY_OUTPUT"
      exit 1
    }

echo "$DEPLOY_OUTPUT"

COLLATERAL_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'CollateralManagement proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
DISCOVERY_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'FlyoverDiscovery proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
PEGIN_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'PegInContract proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
PEGOUT_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'PegOutContract proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)

if [ -z "$COLLATERAL_PROXY" ] || [ -z "$DISCOVERY_PROXY" ] || [ -z "$PEGIN_PROXY" ] || [ -z "$PEGOUT_PROXY" ]; then
    echo "ERROR: Failed to parse contract addresses from deployment output"
    exit 1
fi


echo ""
echo "Verifying deployed contracts..."
for CONTRACT_VAR in PEGIN_PROXY PEGOUT_PROXY COLLATERAL_PROXY DISCOVERY_PROXY; do
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
grep -vE "^(PEGIN_CONTRACT_ADDRESS|PEGOUT_CONTRACT_ADDRESS|COLLATERAL_MANAGEMENT_ADDRESS|DISCOVERY_ADDRESS)=" /"$ENV_FILE" > "$temp_env_file"
{
  echo "PEGIN_CONTRACT_ADDRESS=$PEGIN_PROXY"
  echo "PEGOUT_CONTRACT_ADDRESS=$PEGOUT_PROXY"
  echo "COLLATERAL_MANAGEMENT_ADDRESS=$COLLATERAL_PROXY"
  echo "DISCOVERY_ADDRESS=$DISCOVERY_PROXY"
} >> "$temp_env_file"
cat "$temp_env_file" > /"$ENV_FILE"
rm "$temp_env_file"
echo "Deployment complete. Updated $ENV_FILE with new contract addresses."
