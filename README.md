# Liquidity Provider Server
[![OpenSSF Scorecard](https://api.scorecard.dev/projects/github.com/rsksmart/liquidity-provider-server/badge)](https://scorecard.dev/viewer/?uri=github.com/rsksmart/liquidity-provider-server)

This is a server that interacts with a [Liquidity Bridge Contract (LBC)](https://github.com/rsksmart/liquidity-bridge-contract) to provide liquidity for users
as part of the Flyover protocol. This server performs all the necessary operations to play the role of the Liquidity Provider, involving transactions in both
Rootstock and Bitcoin networks.

## How to run
To run the project locally you can follow these steps:

1. `git clone git@github.com:rsksmart/liquidity-provider-server.git`
2. `cd docker-compose/local`
3. `export LPS_STAGE=regtest`
4. `./lps-env.sh up`

This will set up a local environment, please keep in mind that a productive set-up could vary in multiple aspects.

### How to run the tests
For the unit tests you can run `make test` in the root of the repository and for the integration tests please [check this file](test/integration/Readme.md)

### Installing the project
If you want to play with the code and make modifications to it then run the following commands (remember that you need to have Go installed with the version
specified in the `go.mod` file):
1. `git clone git@github.com:rsksmart/liquidity-provider-server.git`
2. `make tools`

## Configuration

### Environment variables
To see the required environment variables to run an instance of this server and its description check the [Environment.md](docs/Environment.md) file.

### API
The HTTP API provided in this server is divided in two; the public API and the Management (private) API. To understand the purpose of each one of those
API check the [LP Management file](docs/LP-Management.md#context).

To see the details of the interface itself and the structures involved check the [OpenAPI.yml](OpenApi.yml) file that is in the root of the repository.

### Dependencies
The server has the following dependencies:
- Rootstock node
- Bitcoin node
- MongoDB instance

**IMPORTANT**: liquidity provider server performs sensitive operations and uses non publicly enabled functionality of both Rootstock and Bitcoin nodes.
This means that the nodes used to run this server must be private and well protected, the usage of public nodes or nodes that are not properly secured
might lead to a loss of funds.

P.S.: if you run the server locally you'll see that the docker compose includes more services than the previously mentioned, that is because the ones
mentioned before are the minimal dependencies, but in order to run a fully functional environment more dependencies might be required.

## Main operations
- **PegIn**: process of converting BTC into RBTC. [Here](docs/diagrams/PegIn.mmd) is a diagram with a detailed view of the process.
- **PegOut**: process of converting RBTC into BTC. [Here](docs/diagrams/PegOut.mmd) is a diagram with a detailed view of the process.

## LPS Utilities
The [cmd/utils](cmd/utils) directory contains scripts with different utilities for the liquidity providers. You can either run them directly
with `go run` or build them with `make utils`. You can run the scripts with the `--help` flag to see the available options. The current utilities are:
- **update_provider_url**: updates the URL of a liquidity provider provided when the discovery function of the Liquidity Bridge Contract is executed.
- **register_pegin**: register a PegIn transaction within the Liquidity Bridge Contract. Most times, this script is only required to execute refunds
on special cases. This script requires an input file whose structure can be found [the input-example.json](cmd/utils/register_pegin/input-example.json) file.
- **refund_user_pegout**: executes a refund for a user's peg-out operation through the Liquidity Bridge Contract. This is used when a peg-out operation needs to be refunded back to the user's RSK address. The script requires the quote hash of the operation to refund.
- **key_conversion**: shows the corresponding BTC and RSK address for a given private key and encrypts it into a keystore, accepts the key either in WIF or hex format. The key can be provided through the terminal, a file or an existing keystore.

### Monitoring Service
The project includes a Bitcoin balance monitoring service that tracks specified BTC addresses and exposes metrics at `http://<host>:8080/metrics` using Prometheus `https://prometheus.io/`.

To run the monitoring service:
```bash
make monitoring
```
The service is configured in `docker-compose/monitoring/src/config.ts` and supports both testnet and mainnet monitoring:

- MONITORED_ADDRESSES: The set of addresses to be monitored. Each address should have an alias that will be used in the metrics.
- MONITOR_CONFIG: The configuration for the monitoring service.
  - pollingIntervalSeconds: How often the service will check the bitcoin balance of the monitored addresses in seconds.
  - monitorName: The name of the monitoring service.
  - network: The network to monitor (mainnet or testnet).

The service can be configured to monitor other addresses by modifying the `MONITORED_ADDRESSES` array in `docker-compose/monitoring/src/config.ts`.


### More information
If you're looking forward to integrate with Flyover Protocol then you can check the [Flyover SDK repository](https://github.com/rsksmart/unified-bridges-sdk/tree/main/packages/flyover-sdk).

If you're interested in becoming a liquidity provider then you can read the [Liquidity Provider Management](docs/LP-Management.md) file.
