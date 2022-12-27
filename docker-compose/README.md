# docker-compose for Liquidity Provider Server

This docker-compose files can be used to quickly spin up an environment with the server and its
dependant services (`bitcoind` and `rskj`) for either `regtest` or `testnet`

## Deploy locally
* Use scripts located on `local` directory
* If there is any changes on the Liquidity Bridge Contratcs you need to deploy localy in your environemnt and then grab the  LiquidityBridgeContractProxy address.
- - export LBC_ADDR="NEW ADDRESS"
- - export LPS_STAGE=regtest
- chmod +x lps-env.sh
- - ./lps-env.sh up

## Deploy on dev server with testnet config

```
docker-compose --env-file .env.testnet down && 
docker-compose --env-file .env.testnet build --no-cache && docker-compose --env-file .env.testnet up -d
```
