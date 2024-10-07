#!/bin/bash

function cleanup() {
     docker compose --project-name lp-swap down -v
     exit 1
}

if [ "$1" == "-clean" ]; then
  cleanup
fi

trap cleanup INT

LP1="0x"$(jq -r ".lp1.key.address" config.json)
LP2="0x"$(jq -r ".lp2.key.address" config.json)

jq -r .ghToken config.json > gh_token.txt

docker compose build btc01 btc02 rskj
docker compose -f docker-compose.lps.yml build

docker compose up -d btc01 btc02 rskj

echo "Waiting for RskJ to be up and running..."
while true
do
  sleep 3
  curl -s "http://127.0.0.1:4444" -X POST -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"eth_chainId","params": [],"id":1}' && echo "RskJ is up and running" && break
done

echo "Transferring 10 RBTC to $LP1..."
TX_HASH=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendTransaction\",\"params\": [{\"from\": \"0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826\", \"to\": \"$LP1\", \"value\": \"0x8AC7230489E80000\"}],\"id\":1}" | jq -r ".result")
echo "Result: $TX_HASH"

echo "Transferring 10 RBTC to $LP2..."
TX_HASH=$(curl -s -X POST "http://127.0.0.1:4444" -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"eth_sendTransaction\",\"params\": [{\"from\": \"0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826\", \"to\": \"$LP2\", \"value\": \"0x8AC7230489E80000\"}],\"id\":1}" | jq -r ".result")
echo "Result: $TX_HASH"

LBC_ADDR_LINE=$(docker compose run --rm lbc-deployer bash deploy-lbc.sh)
echo "LBC_ADDR_LINE: $LBC_ADDR_LINE"
docker exec btc01 bitcoin-cli -rpcuser=test -rpcpassword=test -rpcport=5555 -rpcconnect=127.0.0.1 addnode btc02:18444 add

jq -r ".lp1.key" config.json > lp1-key.json
jq -r ".lp2.key" config.json > lp2-key.json

LBC_ADDRESS=$(echo "$LBC_ADDR_LINE" | grep "LBC_ADDR=")

sed -e 's/SECRET_SRC=aws/SECRET_SRC=env/g' \
  -e 's/PROVIDER_NAME=\"Default Provider\"/PROVIDER_NAME=\"Provider 1\"/g' \
  -e 's/LOG_FILE=\/home\/lps\/logs\/lps.log/LOG_FILE=/g' \
  -e 's/MONGODB_HOST=mongodb/MONGODB_HOST=lp-swap-mongodb-1/g' \
  -e 's/MONGODB_USER=root/MONGODB_USER=flyover-user/g' \
  -e 's/MONGODB_PASSWORD=root/MONGODB_PASSWORD=flyover-password/g' \
  -e "s/LBC_ADDR=/$LBC_ADDRESS/g" \
  -e "s/BTC_ENDPOINT=bitcoind:5555/BTC_ENDPOINT=btc01:5555/g" \
  -e "s/ENABLE_MANAGEMENT_API=false/ENABLE_MANAGEMENT_API=true/g" \
  -e 's/KEYSTORE_FILE=geth_keystore\/UTC--2024-01-29T16-36-09.688642000Z--9d93929a9099be4355fc2389fbf253982f9df47c/KEYSTORE_FILE=\/mnt\/lp1-key.json/g' \
  ../../sample-config.env > lp1.env

sed -e 's/SECRET_SRC=aws/SECRET_SRC=env/g' \
  -e 's/PROVIDER_NAME=\"Default Provider\"/PROVIDER_NAME=\"Provider 2\"/g' \
  -e 's/LOG_FILE=\/home\/lps\/logs\/lps.log/LOG_FILE=/g' \
  -e 's/MONGODB_HOST=mongodb/MONGODB_HOST=lp-swap-mongodb-2/g' \
  -e 's/MONGODB_USER=root/MONGODB_USER=flyover-user/g' \
  -e 's/MONGODB_PASSWORD=root/MONGODB_PASSWORD=flyover-password/g' \
  -e "s/LBC_ADDR=/$LBC_ADDRESS/g" \
  -e "s/BASE_URL=\"http:\/\/localhost:8080\"/BASE_URL=\"http:\/\/localhost:8081\"/g" \
  -e "s/BTC_ENDPOINT=bitcoind:5555/BTC_ENDPOINT=btc02:5555/g" \
  -e "s/ENABLE_MANAGEMENT_API=false/ENABLE_MANAGEMENT_API=true/g" \
  -e 's/KEYSTORE_FILE=geth_keystore\/UTC--2024-01-29T16-36-09.688642000Z--9d93929a9099be4355fc2389fbf253982f9df47c/KEYSTORE_FILE=\/mnt\/lp2-key.json/g' \
  ../../sample-config.env > lp2.env

docker compose --project-name lp-swap -f docker-compose.lps.yml --env-file=lp1.env up -d
docker compose --project-name lp-swap -f docker-compose.lps.yml --env-file=lp2.env up -d --scale lps=2 --scale mongodb=2 --no-recreate

while [[ $(docker container inspect -f '{{.State.Status}}' lp-swap-lps-1) != "exited" ]]; do
  sleep 1
done
docker start lp-swap-lps-1

while [[ $(docker container inspect -f '{{.State.Status}}' lp-swap-lps-2) != "exited" ]]; do
  sleep 1
done
docker start lp-swap-lps-2

docker compose build --build-arg LBC_ADDRESS="${LBC_ADDRESS//LBC_ADDR=/}" --build-arg GH_TOKEN="$(jq -r .ghToken config.json)" ui
docker compose --project-name lp-swap up -d ui

echo "Both providers are up and running."

while true; do
    echo "Insert the ID of LP to sunset"
    read -r -p "-> " choice

    if [ "$choice" -eq 1 ] || [ "$choice" -eq 2 ]; then
        SUNSET_ID=$choice
        break
    else
        echo "Invalid choice. The only LPs are 1 and 2"
    fi
done

echo "Sunsetting LP $SUNSET_ID..."
rm -f cookie_jar.txt

echo "Insert LP $SUNSET_ID management user"
read -r -p "-> " MANAGEMENT_USER

echo "Insert LP $SUNSET_ID management password"
read -r -p "-> " MANAGEMENT_PWD

if [ "$choice" -eq 1 ]; then
  SUNSET_URL="http://localhost:8080"
else
  SUNSET_URL="http://localhost:8081"
fi

CSRF_TOKEN=$(curl -s -c cookie_jar.txt -H 'Content-Type: application/json' \
  "$SUNSET_URL/management" | sed -n 's/.*name="csrf"[^>]*value="\([^"]*\)".*/\1/p')

CSRF_TOKEN=${CSRF_TOKEN//&#43;/+}
echo "CSRF_TOKEN -> $CSRF_TOKEN"
curl -s -b cookie_jar.txt -c cookie_jar.txt "$SUNSET_URL/management/login" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' \
  --data "{
     \"username\": \"$MANAGEMENT_USER\",
     \"password\": \"$MANAGEMENT_PWD\"
  }"

curl -s -b cookie_jar.txt "$SUNSET_URL/providers/resignation" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' -X POST && \
  echo "LP $SUNSET_ID resigned successfully" || echo "Error resigning LP $SUNSET_ID"

read -r -p "Press enter to withdraw funds from LP $SUNSET_ID. Remember before withdrawing you need to wait the resign blocks!" choice

curl -s -b cookie_jar.txt "$SUNSET_URL/providers/withdrawCollateral" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -H 'Content-Type: application/json' -X POST && \
  echo "LP $SUNSET_ID withdrew the collateral successfully" || echo "Error withdrawing LP $SUNSET_ID collateral"
