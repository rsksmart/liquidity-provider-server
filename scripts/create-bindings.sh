#!/bin/sh

IMAGE=$1

if [ -z "$IMAGE" ]; then
  echo "Usage: $0 <docker-image>"
  echo "Example: $0 ghcr.io/liquidity-bridge-contracts/liquidity-bridge-contracts:latest"
  exit 1
fi

echo "Generating bindings using the image $IMAGE from the Liquidity Bridge Contract docker registry"
docker create --name temp_lbc_deployer "$IMAGE"
mkdir -p bindings
mkdir -p bindings/abi
docker cp temp_lbc_deployer:/home/node/out bindings/out

CONTRACTS=(
  "IPegIn                  PeginContract                  pegin"
  "IPegOut                 PegoutContract                 pegout"
  "ICollateralManagement   CollateralManagementContract   collateral_management"
  "IBridge                 RskBridge                      bridge"
  "IFlyoverDiscovery       FlyoverDiscovery               discovery"
  "Flyover                 Flyover                        flyover"
)

for entry in "${CONTRACTS[@]}"; do
  read -r CONTRACT TYPE NAME <<< "$entry"

  ABI_PATH="bindings/abi/${CONTRACT}.abi.json"
  ARTIFACT="bindings/out/${CONTRACT}.sol/${CONTRACT}.json"
  OUT_GO="bindings/${NAME}_contract.go"
  DEST="internal/adapters/dataproviders/rootstock/bindings/${NAME}/${NAME}_contract.go"

  echo "Generating bindings for $CONTRACT ($TYPE → $NAME)"

  mkdir -p bindings/abi
  mkdir -p "$(dirname "$DEST")"

  jq '.abi' "$ARTIFACT" > "$ABI_PATH"

  abigen --v2 \
    --abi "$ABI_PATH" \
    --type "$TYPE" \
    --pkg bindings \
    --out "$OUT_GO"

  mv "$OUT_GO" "$DEST"
done

docker rm temp_lbc_deployer && rm -r bindings