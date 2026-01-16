#!/bin/bash

set -e

COMMIT_HASH=$(git rev-parse HEAD)
COMMIT_TAG=$(git describe --exact-match --tags 2>/dev/null || echo "")
export COMMIT_HASH
export COMMIT_TAG

# Detect OS
OS_TYPE="$(uname)"
ARCH_TYPE="$(uname -m)"

case "$ARCH_TYPE" in
  arm64|aarch64) LPS_DOCKER_ARCH_DEFAULT="arm64" ;;
  x86_64|amd64) LPS_DOCKER_ARCH_DEFAULT="amd64" ;;
  *)
    echo "Unsupported architecture: $ARCH_TYPE"
    exit 1
    ;;
esac

if [[ "$OS_TYPE" == "Darwin" ]]; then
    # macOS
    echo "Running on macOS"
    SED_INPLACE=("sed" "-i" "")
    LPS_DOCKERFILE_DEFAULT="docker-compose/lps/Dockerfile.prebuilt"
elif [[ "$OS_TYPE" == "Linux" ]]; then
    # Assume Ubuntu or other Linux
    echo "Running on Linux"
    SED_INPLACE=("sed" "-i")
    LPS_DOCKERFILE_DEFAULT="docker-compose/lps/Dockerfile"
else
    echo "Unsupported OS: $OS_TYPE"
    exit 1
fi

if [ -z "${LPS_DOCKERFILE}" ]; then
  export LPS_DOCKERFILE="$LPS_DOCKERFILE_DEFAULT"
fi

if [ -z "${LPS_DOCKER_ARCH}" ]; then
  export LPS_DOCKER_ARCH="$LPS_DOCKER_ARCH_DEFAULT"
fi

if [ -z "${LPS_STAGE}" ]; then
  echo "LPS_STAGE is not set. Exit 1"
  exit 1
elif [ "$LPS_STAGE" = "regtest" ]; then
  ENV_FILE=".env.regtest"
  if [ ! -f "$ENV_FILE" ]; then
    echo "Creating $ENV_FILE from sample-config.env..."
    cp ../../sample-config.env "$ENV_FILE"
  else
    echo "Using existing $ENV_FILE"
  fi
elif [ "$LPS_STAGE" = "testnet" ]; then
  ENV_FILE=".env.testnet"
else
  echo "Invalid LPS_STAGE: $LPS_STAGE"
  exit 1
fi

if [ -z "${LPS_UID}" ]; then
  LPS_UID=$(id -u)
  export LPS_UID
  if [ "$LPS_UID" = "0" ]; then
    echo "Please set LPS_UID env var or run as a non-root user"
    exit 1
  fi
fi

echo "LPS_STAGE: $LPS_STAGE; ENV_FILE: $ENV_FILE; LPS_UID: $LPS_UID"

if [ -f "$ENV_FILE" ]; then
  # Force Management API to be enabled
  "${SED_INPLACE[@]}" 's/^ENABLE_MANAGEMENT_API=.*/ENABLE_MANAGEMENT_API=true/' "$ENV_FILE"
  # Don't use extra sources on local env
  "${SED_INPLACE[@]}" 's/^RSK_EXTRA_SOURCES=.*/RSK_EXTRA_SOURCES=/' "$ENV_FILE"
  "${SED_INPLACE[@]}" 's/^BTC_EXTRA_SOURCES=.*/BTC_EXTRA_SOURCES=/' "$ENV_FILE"
fi

SCRIPT_CMD=$1
if [ -z "${SCRIPT_CMD}" ]; then
  echo "Command is not provided"
  exit 1
elif [ "$SCRIPT_CMD" = "up" ]; then
  echo "Starting LPS env up..."
elif [ "$SCRIPT_CMD" = "down" ]; then
  echo "Shutting LPS env down..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml down
  exit 0
elif [ "$SCRIPT_CMD" = "build" ]; then
  echo "Building LPS env..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml build
  exit 0
elif [ "$SCRIPT_CMD" = "stop" ]; then
  echo "Stopping LPS env..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml stop
  exit 0
elif [ "$SCRIPT_CMD" = "ps" ]; then
  echo "List of running services:"
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml ps
  exit 0
elif [ "$SCRIPT_CMD" = "deploy" ]; then
  echo "Stopping LPS..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml stop lps
  echo "Building LPS..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml build lps
  echo "Starting LPS..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml start lps
  exit 0
elif [ "$SCRIPT_CMD" = "import-rsk-db" ]; then
  echo "Importing rsk db..."
  docker compose --env-file "$ENV_FILE" run -d rskj java -Xmx6g -Drpc.providers.web.http.bind_address=0.0.0.0 -Drpc.providers.web.http.hosts.0=localhost -Drpc.providers.web.http.hosts.1=rskj -cp rskj-core.jar co.rsk.Start --"${LPS_STAGE}" --import
  exit 0
elif [ "$SCRIPT_CMD" = "start-bitcoind" ]; then
  echo "Starting bitcoind..."
  docker compose --env-file "$ENV_FILE" -f docker-compose.yml up -d bitcoind
  exit 0
else
  echo "Invalid command: $SCRIPT_CMD"
  exit 1
fi

BTCD_HOME="${BTCD_HOME:-./volumes/bitcoind}"
POWPEG_PEGIN_HOME="${POWPEG_PEGIN_HOME:-./volumes/powpeg/pegin}"
POWPEG_PEGOUT_HOME="${POWPEG_PEGOUT_HOME:-./volumes/powpeg/pegout}"
RSKJ_HOME="${RSKJ_HOME:-./volumes/rskj}"
LPS_HOME="${LPS_HOME:-./volumes/lps}"
MONGO_HOME="${MONGO_HOME:-./volumes/mongo}"
LOCALSTACK_HOME="${LOCALSTACK_HOME:-./volumes/localstack}"

# Set LOG_FILE environment variable for LPS to write logs to file
export LOG_FILE="/home/lps/logs/lps.log"

# Fixed directory creation with proper operator precedence
[ -d "$BTCD_HOME" ] || (mkdir -p "$BTCD_HOME" && chown "$LPS_UID" "$BTCD_HOME")
[ -d "$RSKJ_HOME" ] || (mkdir -p "$RSKJ_HOME/db" && mkdir -p "$RSKJ_HOME/logs" && chown -R "$LPS_UID" "$RSKJ_HOME")
[ -d "$POWPEG_PEGIN_HOME" ] || (mkdir -p "$POWPEG_PEGIN_HOME/db" && mkdir -p "$POWPEG_PEGIN_HOME/logs" && chown -R "$LPS_UID" "$POWPEG_PEGIN_HOME" && chmod -R 777 "$POWPEG_PEGIN_HOME")
[ -d "$POWPEG_PEGOUT_HOME" ] || (mkdir -p "$POWPEG_PEGOUT_HOME/db" && mkdir -p "$POWPEG_PEGOUT_HOME/logs" && chown -R "$LPS_UID" "$POWPEG_PEGOUT_HOME" && chmod -R 777 "$POWPEG_PEGOUT_HOME")
[ -d "$LPS_HOME" ] || (mkdir -p "$LPS_HOME/logs" && chmod -R 777 "$LPS_HOME")
[ -d "$MONGO_HOME" ] || (mkdir -p "$MONGO_HOME/db" && chown -R "$LPS_UID" "$MONGO_HOME")
[ -d "$LOCALSTACK_HOME" ] || (mkdir -p "$LOCALSTACK_HOME/db" && mkdir -p "$LOCALSTACK_HOME/logs" && chown -R "$LPS_UID" "$LOCALSTACK_HOME")
[ -d "./volumes/loki" ] || (mkdir -p "./volumes/loki" && chmod 777 "./volumes/loki")

echo "LPS_UID: $LPS_UID; BTCD_HOME: '$BTCD_HOME'; RSKJ_HOME: '$RSKJ_HOME'; LPS_HOME: '$LPS_HOME'; MONGO_HOME: '$MONGO_HOME'; POWPEG_PEGIN_HOME: '$POWPEG_PEGIN_HOME'; POWPEG_PEGOUT_HOME: '$POWPEG_PEGOUT_HOME'; LOCALSTACK_HOME: '$LOCALSTACK_HOME'; LOKI_HOME: './volumes/loki'"

# start bitcoind and RSKJ dependant services
docker compose --env-file "$ENV_FILE" up -d bitcoind rskj mongodb localstack

# shellcheck disable=SC1090
. ./"$ENV_FILE"

echo "Waiting for RskJ to be up and running..."
while true
do
  sleep 3
  curl -s "http://127.0.0.1:4444" -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"eth_chainId","params": [],"id":1}' \
    && echo "RskJ is up and running" \
    && break
done

echo "Waiting for Bitcoind to be up and running..."
while true
do
  sleep 3
  curl -s "http://127.0.0.1:5555" -X POST --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getnetworkinfo", "params": [], "id":"1"}' | grep "\"result\":{" \
    && echo "Bitcoind is up and running" \
    && break
done

curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getwalletinfo", "params": [], "id":"getwallet"}' | grep "{\"result\":null,\"error\":{\"code\":-18" \
  && echo "No default wallet" \
  && echo "Creating wallet" \
  && curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "createwallet", "params": ["main", false, false, "test-password", true, false, true], "id":"createwallet"}' \
  && curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getnewaddress", "params": ["main"], "id":"getnewaddress"}' \
  | jq .result | xargs -I ADDRESS curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "generatetoaddress", "params": [500, "ADDRESS"], "id":"generatetoaddress"}' \
  && echo "Wallet created and generated 500 blocks" \
  && curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "createwallet", "params": ["rsk-wallet", true, true, "", false, false, true], "id":"createwallet"}' \
  && curl -s "http://127.0.0.1:5555/wallet/rsk-wallet" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "importpubkey", "params": ["0232858a5faa413101831afe7a880da9a8ac4de6bd5e25b4358d762ba450b03c22", "", false], "id":"importpubkey"}' \
  && curl -s "http://127.0.0.1:5555/wallet/rsk-wallet" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "importaddress", "params": ["n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6", "", false], "id":"importaddress"}' \
  && curl -s "http://127.0.0.1:5555/wallet/main" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "walletpassphrase", "params": ["test-password", 30000], "id":"walletpassphrase"}' \
  && curl -s "http://127.0.0.1:5555/wallet/main" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "sendtoaddress", "params": { "amount": 500, "fee_rate": 25, "address": "n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6" }, "id":"sendtoaddress"}' \
  && curl -s "http://127.0.0.1:5555/wallet/main" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "getnewaddress", "params": ["main"], "id":"getnewaddress"}' \
    | jq .result | xargs -I ADDRESS curl -s "http://127.0.0.1:5555" --user "$BTC_USERNAME:$BTC_PASSWORD" -H "Content-Type: application/json" -d '{"jsonrpc": "1.0", "method": "generatetoaddress", "params": [1, "ADDRESS"], "id":"generatetoaddress"}'

if [ "$LPS_STAGE" = "regtest" ]; then
  PROVIDER_TX_COUNT=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getTransactionCount\",\"params\": [\"$LIQUIDITY_PROVIDER_RSK_ADDR\",\"latest\"],\"id\":1}" | jq -r ".result")
  if [ "$PROVIDER_TX_COUNT" = "0x0" ]; then
    echo "Transferring funds to $LIQUIDITY_PROVIDER_RSK_ADDR..."

    TX_HASH=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendTransaction\",\"params\": [{\"from\": \"0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826\", \"to\": \"$LIQUIDITY_PROVIDER_RSK_ADDR\", \"value\": \"0x3635C9ADC5DEA00000\"}],\"id\":1}" | jq -r ".result")
    echo "Result: $TX_HASH"
    sleep 10
  else
    echo "No need to fund the '$LIQUIDITY_PROVIDER_RSK_ADDR' provider. Nonce: $PROVIDER_TX_COUNT"
  fi

  # Path to liquidity-bridge-contract repo for Foundry deployment
  LBC_REPO_PATH="${LBC_REPO_PATH:-$HOME/liquidity-bridge-contract}"
  # Private key for local deployment (default Hardhat/Foundry test key)
  DEPLOYER_PRIVATE_KEY="${DEPLOYER_PRIVATE_KEY:-0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80}"
  DEPLOYER_ADDRESS="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

  # Check if contracts are already deployed by checking for code at the addresses
  echo "Checking if Flyover contracts are deployed..."
  CONTRACTS_MISSING=false

  for CONTRACT_VAR in PEGIN_CONTRACT_ADDRESS PEGOUT_CONTRACT_ADDRESS COLLATERAL_MANAGEMENT_ADDRESS DISCOVERY_ADDRESS; do
    CONTRACT_ADDR=$(eval echo "\$$CONTRACT_VAR")
    if [ -z "$CONTRACT_ADDR" ] || [ "$CONTRACT_ADDR" = "0x0000000000000000000000000000000000000000" ]; then
      echo "  $CONTRACT_VAR is not set or is zero address"
      CONTRACTS_MISSING=true
      continue
    fi

    CODE=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" \
      -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getCode\",\"params\": [\"$CONTRACT_ADDR\",\"latest\"],\"id\":1}" | jq -r ".result")

    if [ "$CODE" = "0x" ] || [ -z "$CODE" ]; then
      echo "  $CONTRACT_VAR ($CONTRACT_ADDR) has no code deployed"
      CONTRACTS_MISSING=true
    else
      echo "  ✓ $CONTRACT_VAR ($CONTRACT_ADDR) is deployed"
    fi
  done

  if [ "$CONTRACTS_MISSING" = true ]; then
    echo ""
    echo "Some contracts are missing. Deploying Flyover contracts using Foundry..."

    if [ ! -d "$LBC_REPO_PATH" ]; then
      echo "ERROR: liquidity-bridge-contract repo not found at $LBC_REPO_PATH"
      echo "Please clone the repo or set LBC_REPO_PATH environment variable"
      exit 1
    fi

    # Fund the deployer account from RSK cow account
    echo "Funding deployer account $DEPLOYER_ADDRESS..."
    DEPLOYER_BALANCE=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" \
      -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getBalance\",\"params\": [\"$DEPLOYER_ADDRESS\",\"latest\"],\"id\":1}" | jq -r ".result")

    if [ "$DEPLOYER_BALANCE" = "0x0" ]; then
      echo "Transferring funds to deployer..."
      curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" \
        -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendTransaction\",\"params\": [{\"from\": \"0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826\", \"to\": \"$DEPLOYER_ADDRESS\", \"value\": \"0x3635C9ADC5DEA00000\"}],\"id\":1}"
      sleep 10
    fi

    # Deploy contracts using Foundry
    pushd "$LBC_REPO_PATH" > /dev/null

    # Save current branch and checkout the deployment script branch
    ORIGINAL_BRANCH=$(git rev-parse --abbrev-ref HEAD)
    echo "Switching to QA-Test branch for deployment..."
    git fetch origin QA-Test 2>/dev/null || true
    git checkout QA-Test 2>/dev/null || git checkout -b QA-Test origin/QA-Test

    echo "Running Foundry deployment script..."
    export DEV_SIGNER_PRIVATE_KEY="$DEPLOYER_PRIVATE_KEY"
    DEPLOY_OUTPUT=$(forge script script/deployment/DeployFlyover.s.sol:DeployFlyover \
      --rpc-url http://127.0.0.1:4444 \
      --private-key "$DEPLOYER_PRIVATE_KEY" \
      --broadcast \
      --legacy \
      --slow 2>&1) || {
        echo "Foundry deployment failed:"
        echo "$DEPLOY_OUTPUT"
        git checkout "$ORIGINAL_BRANCH" 2>/dev/null || true
        popd > /dev/null
        exit 1
      }

    echo "$DEPLOY_OUTPUT"

    # Switch back to original branch
    git checkout "$ORIGINAL_BRANCH" 2>/dev/null || true
    popd > /dev/null

    # Parse the deployment output to get addresses (portable grep)
    COLLATERAL_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'CollateralManagement proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
    DISCOVERY_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'FlyoverDiscovery proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
    PEGIN_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'PegInContract proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)
    PEGOUT_PROXY=$(echo "$DEPLOY_OUTPUT" | grep -o 'PegOutContract proxy: 0x[a-fA-F0-9]*' | sed 's/.*: //' | head -1)

    if [ -z "$COLLATERAL_PROXY" ] || [ -z "$DISCOVERY_PROXY" ] || [ -z "$PEGIN_PROXY" ] || [ -z "$PEGOUT_PROXY" ]; then
      echo "ERROR: Failed to parse contract addresses from deployment output"
      echo "Please check the deployment output above and update .env.regtest manually"
      exit 1
    fi

    echo ""
    echo "Deployed contract addresses:"
    echo "  CollateralManagement: $COLLATERAL_PROXY"
    echo "  FlyoverDiscovery: $DISCOVERY_PROXY"
    echo "  PegIn: $PEGIN_PROXY"
    echo "  PegOut: $PEGOUT_PROXY"

    # Update .env.regtest with new addresses
    echo ""
    echo "Updating $ENV_FILE with deployed addresses..."
    "${SED_INPLACE[@]}" "s/^PEGIN_CONTRACT_ADDRESS=.*/PEGIN_CONTRACT_ADDRESS=$PEGIN_PROXY/" "$ENV_FILE"
    "${SED_INPLACE[@]}" "s/^PEGOUT_CONTRACT_ADDRESS=.*/PEGOUT_CONTRACT_ADDRESS=$PEGOUT_PROXY/" "$ENV_FILE"
    "${SED_INPLACE[@]}" "s/^COLLATERAL_MANAGEMENT_ADDRESS=.*/COLLATERAL_MANAGEMENT_ADDRESS=$COLLATERAL_PROXY/" "$ENV_FILE"
    "${SED_INPLACE[@]}" "s/^DISCOVERY_ADDRESS=.*/DISCOVERY_ADDRESS=$DISCOVERY_PROXY/" "$ENV_FILE"

    # Re-source the env file to get updated addresses
    # shellcheck disable=SC1090
    . ./"$ENV_FILE"

    echo "Contract addresses updated in $ENV_FILE"

    # Verify deployment
    echo ""
    echo "Verifying deployed contracts..."
    for CONTRACT_VAR in PEGIN_CONTRACT_ADDRESS PEGOUT_CONTRACT_ADDRESS COLLATERAL_MANAGEMENT_ADDRESS DISCOVERY_ADDRESS; do
      CONTRACT_ADDR=$(eval echo "\$$CONTRACT_VAR")
      CODE=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" \
        -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_getCode\",\"params\": [\"$CONTRACT_ADDR\",\"latest\"],\"id\":1}" | jq -r ".result")

      if [ "$CODE" = "0x" ] || [ -z "$CODE" ]; then
        echo "  ✗ $CONTRACT_VAR ($CONTRACT_ADDR) - NO CODE (deployment may have failed)"
      else
        echo "  ✓ $CONTRACT_VAR ($CONTRACT_ADDR) - verified"
      fi
    done
  else
    echo "All Flyover contracts are already deployed!"
  fi
fi

echo ""
echo "Flyover contracts configuration:"
echo "  PEGIN_CONTRACT_ADDRESS: $PEGIN_CONTRACT_ADDRESS"
echo "  PEGOUT_CONTRACT_ADDRESS: $PEGOUT_CONTRACT_ADDRESS"
echo "  COLLATERAL_MANAGEMENT_ADDRESS: $COLLATERAL_MANAGEMENT_ADDRESS"
echo "  DISCOVERY_ADDRESS: $DISCOVERY_ADDRESS"

docker compose --env-file "$ENV_FILE" up -d powpeg-pegin powpeg-pegout

if [ "$LPS_DOCKERFILE" = "docker-compose/lps/Dockerfile.prebuilt" ]; then
  # Build LPS binary locally (cross-compile) to avoid Go segfault in Docker on Mac
  echo "Building LPS binary locally (cross-compile for linux/${LPS_DOCKER_ARCH})..."
  pushd ../../ > /dev/null
  COMMIT_HASH_VALUE=$(git rev-parse HEAD)
  COMMIT_TAG_VALUE=$(git describe --exact-match --tags 2>/dev/null || echo "")
  mkdir -p build
  GOOS=linux GOARCH="${LPS_DOCKER_ARCH}" CGO_ENABLED=0 go build -mod=mod -a -trimpath \
    -ldflags="-s -w -X 'main.BuildVersion=${COMMIT_HASH_VALUE}' -X 'main.BuildTime=$(date)' -X 'github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider.BuildVersion=${COMMIT_TAG_VALUE}' -X 'github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider.BuildRevision=${COMMIT_HASH_VALUE}'" \
    -o ./build/liquidity-provider-server ./cmd/application/main.go

  if [ ! -f "./build/liquidity-provider-server" ]; then
    echo "ERROR: Binary build failed!"
    exit 1
  fi

  # Verify it's a Linux binary for the requested architecture
  if [ "$LPS_DOCKER_ARCH" = "arm64" ]; then
    if ! file ./build/liquidity-provider-server | grep -q "ELF.*aarch64"; then
      echo "WARNING: Binary might not be correct Linux/arm64 format"
      file ./build/liquidity-provider-server
    fi
  elif [ "$LPS_DOCKER_ARCH" = "amd64" ]; then
    if ! file ./build/liquidity-provider-server | grep -q "ELF.*x86-64"; then
      echo "WARNING: Binary might not be correct Linux/amd64 format"
      file ./build/liquidity-provider-server
    fi
  fi

  popd > /dev/null
fi

# start LPS
docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml build lps
docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.lps.yml up -d lps

FAIL=true
for _ in $(seq 1 10);
do
  sleep 5
  curl -s "http://localhost:8080/health" \
    && echo "LPS is up and running" \
    && FAIL=false \
    || echo "LPS is not up yet"
  if [ "$FAIL" = false ]; then
    break
  fi
done

if [ "$FAIL" = true ]; then
  echo "LPS failed to start"
  exit 1
fi

rm -f cookie_jar.txt

PASSWORD_FILE_PATH="/tmp/management_password.txt"

echo "Checking for management_password.txt..."
if ! docker exec lps01 test -f "$PASSWORD_FILE_PATH"; then
  echo "management_password.txt not found. Skipping configuration steps"
  exit 0
fi

echo "management_password.txt found. Proceeding with configuration."

MANAGEMENT_PWD=$(docker exec lps01 cat "$PASSWORD_FILE_PATH")

CSRF_TOKEN=$(curl -s -c cookie_jar.txt \
                      -H 'Accept: */*' \
                      -H 'Connection: keep-alive' \
                      -H 'Content-Type: application/json' \
                      -H 'Origin: http://localhost:8080' \
                      -H 'Sec-Fetch-Dest: empty' \
                      -H 'Sec-Fetch-Mode: cors' \
                      -H 'Sec-Fetch-Site: same-origin' \
  "http://localhost:8080/management" | sed -n 's/.*name="csrf"[^>]*value="\([^"]*\)".*/\1/p')

# shellcheck disable=SC2001
CSRF_TOKEN=$(echo "$CSRF_TOKEN" | sed 's/&#43;/+/g')
curl -s -b cookie_jar.txt -c cookie_jar.txt "http://localhost:8080/management/login" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Origin: http://localhost:8080' \
  -H 'Referer: http://localhost:8080/management' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  --data "{
     \"username\": \"admin\",
     \"password\": \"$MANAGEMENT_PWD\"
  }" || { echo "Error: login to Management UI failed"; exit 1; }

echo "Setting up general regtest configuration"
curl -sfS -b cookie_jar.txt 'http://localhost:8080/configuration' \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Origin: http://localhost:8080' \
  -H 'Referer: http://localhost:8080/management' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  --data '{
      "configuration": {
          "rskConfirmations": {
              "100000000000000000": 4,
              "2000000000000000000": 20,
              "400000000000000000": 12,
              "4000000000000000000": 40,
              "8000000000000000000": 80
          },
          "btcConfirmations": {
              "100000000000000000": 2,
              "2000000000000000000": 10,
              "400000000000000000": 6,
              "4000000000000000000": 20,
              "8000000000000000000": 40
          },
          "publicLiquidityCheck": true
      }
  }' || { echo "Error in configuring general regtest configuration"; exit 1; }

echo "Setting up pegin regtest configuration"
CURL_OUTPUT=$(curl -s -w '\n%{http_code}' -b cookie_jar.txt 'http://localhost:8080/pegin/configuration' \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Origin: http://localhost:8080' \
  -H 'Referer: http://localhost:8080/management' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  --data '{
      "configuration": {
          "timeForDeposit": 3600,
          "callTime": 7200,
          "penaltyFee": "1000000000000000",
          "maxValue": "10000000000000000000",
          "minValue": "600000000000000000",
          "feePercentage": 0.33,
          "fixedFee": "200000000000000"
      }
  }')

HTTP_STATUS=$(echo "$CURL_OUTPUT" | tail -n1)
RESPONSE_BODY=$(echo "$CURL_OUTPUT" | sed '$d')

if [ "$HTTP_STATUS" -lt 200 ] || [ "$HTTP_STATUS" -ge 300 ]; then
  echo "Error in configuring pegin regtest configuration"
  echo "HTTP Status: $HTTP_STATUS"
  echo "Response Body:"
  echo "$RESPONSE_BODY"
  exit 1
fi

echo "Setting up pegout regtest configuration"
CURL_OUTPUT=$(curl -s -w '\n%{http_code}' -b cookie_jar.txt 'http://localhost:8080/pegout/configuration' \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Origin: http://localhost:8080' \
  -H 'Referer: http://localhost:8080/management' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  --data '{
      "configuration": {
          "timeForDeposit": 3600,
          "expireTime": 10800,
          "penaltyFee": "1000000000000000",
          "maxValue": "10000000000000000000",
          "minValue": "600000000000000000",
          "expireBlocks": 500,
          "bridgeTransactionMin": "1500000000000000000",
          "feePercentage": 0.33,
          "fixedFee": "200000000000000"
      }
  }')

HTTP_STATUS=$(echo "$CURL_OUTPUT" | tail -n1)
RESPONSE_BODY=$(echo "$CURL_OUTPUT" | sed '$d')

if [ "$HTTP_STATUS" -lt 200 ] || [ "$HTTP_STATUS" -ge 300 ]; then
  echo "Error in configuring pegout regtest configuration"
  echo "HTTP Status: $HTTP_STATUS"
  echo "Response Body:"
  echo "$RESPONSE_BODY"
  exit 1
fi

echo "Creating trusted account for regtest"
CURL_OUTPUT=$(curl -s -w '\n%{http_code}' -b cookie_jar.txt 'http://localhost:8080/management/trusted-accounts' \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Origin: http://localhost:8080' \
  -H 'Referer: http://localhost:8080/management' \
  -H 'Sec-Fetch-Dest: empty' \
  -H 'Sec-Fetch-Mode: cors' \
  -H 'Sec-Fetch-Site: same-origin' \
  --data "{
      \"address\": \"$TRUSTED_ACCOUNT_ADDRESS\",
      \"name\": \"Boletaz\",
      \"btcLockingCap\": 3000000000000000000,
      \"rbtcLockingCap\": 3000000000000000000
  }")

HTTP_STATUS=$(echo "$CURL_OUTPUT" | tail -n1)
RESPONSE_BODY=$(echo "$CURL_OUTPUT" | sed '$d')

if [ "$HTTP_STATUS" -ge 200 ] && [ "$HTTP_STATUS" -lt 300 ]; then
  echo "✓ Trusted account created successfully!"
  echo "  Address: $TRUSTED_ACCOUNT_ADDRESS"
elif echo "$RESPONSE_BODY" | grep -q "already exists"; then
  echo "✓ Trusted account already exists (OK)"
  echo "  Address: $TRUSTED_ACCOUNT_ADDRESS"
else
  echo "Error creating trusted account"
  echo "HTTP Status: $HTTP_STATUS"
  echo "Response Body:"
  echo "$RESPONSE_BODY"
  exit 1
fi

docker compose --env-file "$ENV_FILE" -f docker-compose.yml -f docker-compose.metrics.yml up -d prometheus loki alloy grafana mailhog

echo ""
echo "============================================"
echo "✓ LPS environment is ready!"
echo "  LPS API:    http://localhost:8080"
echo "  Health:     http://localhost:8080/health"
echo "  Management: http://localhost:8080/management"
echo "  Grafana:    http://localhost:3000"
echo "============================================"
