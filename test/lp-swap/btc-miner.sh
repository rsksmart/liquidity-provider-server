#!/bin/bash

if ! [[ "$1" =~ ^[0-9]+$ ]]; then
  echo "Argument must be a number"
  exit 1
fi

RPC_IP=127.0.0.1
RPC_USER="test"
RPC_PASSWORD="test"
RPC_PORT=5555

WALLETS=$(bitcoin-cli -rpcuser=$RPC_USER -rpcpassword=$RPC_PASSWORD -rpcport=$RPC_PORT -rpcconnect=$RPC_IP listwallets)
if ! [[ $WALLETS == *"main"* ]]; then
  bitcoin-cli -rpcuser=$RPC_USER  -rpcpassword=$RPC_PASSWORD -rpcport=$RPC_PORT -rpcconnect=$RPC_IP \
    createwallet "main" false false "test-password" true false true
  ADDRESS=$(bitcoin-cli -rpcuser=$RPC_USER  -rpcpassword=$RPC_PASSWORD -rpcport=$RPC_PORT -rpcconnect=$RPC_IP \
      -rpcwallet=main getnewaddress)
  bitcoin-cli -rpcuser=$RPC_USER  -rpcpassword=$RPC_PASSWORD -rpcport=$RPC_PORT -rpcconnect=$RPC_IP \
      generatetoaddress 100 "$ADDRESS"
fi

while true; do
  bitcoin-cli -rpcuser=$RPC_USER  -rpcpassword=$RPC_PASSWORD -rpcport=$RPC_PORT -rpcconnect=$RPC_IP \
      -rpcwallet=main walletpassphrase "test-password" "$1"
  ADDRESS=$(bitcoin-cli -rpcuser=$RPC_USER  -rpcpassword=$RPC_PASSWORD -rpcport=$RPC_PORT -rpcconnect=$RPC_IP \
    -rpcwallet=main getnewaddress)
  bitcoin-cli -rpcuser=$RPC_USER  -rpcpassword=$RPC_PASSWORD -rpcport=$RPC_PORT -rpcconnect=$RPC_IP \
    -rpcwallet=main -named sendtoaddress address="$ADDRESS" fee_rate=25 amount=0.00001
  bitcoin-cli -rpcuser=$RPC_USER  -rpcpassword=$RPC_PASSWORD -rpcport=$RPC_PORT -rpcconnect=$RPC_IP \
    -rpcwallet=main -generate 1
  sleep "$1"
done
