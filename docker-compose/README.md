## Deploy on Regtest Environment
Go to the [`local` directory](./local) for a detailed explanation of the local set up.

## Deploy on Development Server with Testnet Config

For testnet or mainnet environments, use the docker-compose files directly:

```bash
docker-compose --env-file .env.testnet down &&
docker-compose --env-file .env.testnet build --no-cache &&
docker-compose --env-file .env.testnet up -d
```

:::danger[Troubleshooting]
Encountering difficulties with the Docker setup or Flyover issues? Join the [Rootstock Discord community](http://discord.gg/rootstock) for expert support and assistance. Our dedicated team is ready to help you resolve any problems you may encounter.
:::
