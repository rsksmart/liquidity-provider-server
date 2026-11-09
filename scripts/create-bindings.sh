#!/usr/bin/env bash

IMAGE=$1

if [ -z "$IMAGE" ]; then
  echo "Usage: $0 <docker-image>"
  echo "Example: $0 ghcr.io/rsksmart/liquidity-bridge-contract@sha256:<digest>"
  exit 1
fi

echo "Generating bindings using the image $IMAGE from the Liquidity Bridge Contract docker registry"
docker create --name temp_lbc_deployer "$IMAGE"
mkdir -p bindings
mkdir -p bindings/abi
docker cp temp_lbc_deployer:/home/node/out bindings/out

CONTRACTS=(
  "IPegIn                  PeginContract                  pegin"
  "IPegInCommitFirst       PeginCommitFirstContract       pegin_commit_first"
  "IPegOut                 PegoutContract                 pegout"
  "ICollateralManagement   CollateralManagementContract   collateral_management"
  "IBridge                 RskBridge                      bridge"
  "IFlyoverDiscovery       FlyoverDiscovery               discovery"
  "Flyover                 Flyover                        flyover"
  "IPegInAddressRegistry   PegInAddressRegistryContract   pegin_address_registry"
  "IPauseRegistry          PauseRegistryContract          pause_registry"
  "EmergencyPause          EmergencyPauseContract         emergency_pause"
  "IFlyoverConfigurations  FlyoverConfigurationsContract  flyover_configurations"
  "IPegOutEscrow           PegOutEscrowContract           pegout_escrow"
)

for entry in "${CONTRACTS[@]}"; do
  read -r CONTRACT TYPE NAME <<< "$entry"

  ABI_PATH="bindings/abi/${CONTRACT}.abi.json"
  ARTIFACT="bindings/out/${CONTRACT}.sol/${CONTRACT}.json"
  ABI_EXTENSION="scripts/abi-extensions/${CONTRACT}.abi.json"
  OUT_GO="bindings/${NAME}_contract.go"
  DEST="internal/adapters/dataproviders/rootstock/bindings/${NAME}/${NAME}_contract.go"

  echo "Generating bindings for $CONTRACT ($TYPE → $NAME)"

  mkdir -p bindings/abi
  mkdir -p "$(dirname "$DEST")"

  if [ -f "$ABI_EXTENSION" ]; then
    jq -s '.[0].abi + .[1]' "$ARTIFACT" "$ABI_EXTENSION" > "$ABI_PATH"
  else
    jq '.abi' "$ARTIFACT" > "$ABI_PATH"
  fi

  abigen --v2 \
    --abi "$ABI_PATH" \
    --type "$TYPE" \
    --pkg bindings \
    --out "$OUT_GO"

  mv "$OUT_GO" "$DEST"
done

docker rm temp_lbc_deployer && rm -r bindings
