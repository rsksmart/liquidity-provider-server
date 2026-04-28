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
docker cp temp_lbc_deployer:/home/node/artifacts bindings/artifacts
FILE=$(ls -t bindings/artifacts/build-info/*.json | head -n1)


jq '{
  contracts:
    ( .output.contracts
    | to_entries
    | map(select(.key|test("^contracts/(interfaces|libraries)/"))) # only files under contracts/interfaces
    | map(
        .key as $path
        | (.value | to_entries
          | map({
              key: ($path + ":" + .key),
              value: {
                abi: .value.abi,
                bin: (.value.evm.bytecode.object // "")
              }
            })
        )
      )
    | add
    | from_entries )
}' "$FILE" > bindings/abigen_combined.json

abigen --combined-json bindings/abigen_combined.json --pkg bindings --out bindings/flyover_contracts.go
mv bindings/flyover_contracts.go internal/adapters/dataproviders/rootstock/bindings

docker rm temp_lbc_deployer && rm -r bindings