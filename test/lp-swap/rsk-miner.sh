#!/bin/bash

if ! [[ "$1" =~ ^[0-9]+$ ]]; then
  echo "Argument must be a number"
  exit 1
fi

while true; do
  curl --location 'http://localhost:4444' \
  --header 'Content-Type: application/json' \
  --data '{
      "method": "evm_mine",
      "params": [],
      "id": 1,
      "jsonrpc": "2.0"
  }'
  sleep "$1"
done
