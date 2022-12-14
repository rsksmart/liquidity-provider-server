# docker-compose for Liquidity Provider Server

This docker-compose files can be used to quickly spin up an environment with the server and its
dependant services (`bitcoind` and `rskj`) for either `regtest` or `testnet`

## Deploy locally
* Use scripts located on `local` directory

## Deploy on dev server with testnet config

```
docker-compose --env-file .env.testnet down && 
docker-compose --env-file .env.testnet build --no-cache && docker-compose --env-file .env.testnet up -d
```
