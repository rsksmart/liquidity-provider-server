# Liquidity Provider Server local environment
The provided docker-compose files can be used to quickly spin up an environment with the Liquidity Provider Server and its dependent services (`bitcoind` and `rskj`) for `regtest`. The regtest environment requires a regtest federation and localstack, making the setup process different from testnet or mainnet.

## Deploy Locally (Regtest Environment)

* Use scripts located in the `local` directory while being inside that directory. Using them from other directories might cause issues with the relative paths defined in the different compose files.
* Create an env file and export it as an environment variable `export ENV_FILE=regtest`. If you don't want to create it, the script will use by default the `sample-config.env` file located at the root directory.
* Run the following command to create the environment:
```bash
    ./lps-local.sh
```

## Configurations
The `sample-config.env` file contains a set of flags at the end of the file that can be used to enable or disable certain features of the Liquidity Provider Server. These flags are:
* `FUND_WALLETS`: if enabled, the wallet-funder script will be executed every time the local script is executed. This script is responsible for funding the wallets with the necessary funds (as configured in the env file) to operate in the regtest environment. This script is not idempotent.
* `DEPLOY_CONTRACTS`: if enabled, the deploy-contracts script will be executed every time the local script is executed. This script is responsible for deploying the necessary contracts to operate in the regtest environment. This script is not idempotent.
* `CREATE_POWPEG`: if enabled, the regtest federation nodes will be created and the powpeg will be set up. This script is idempotent.
* `CREATE_MONITORING`: if enabled, the monitoring stack will be created. This script is idempotent.

## Extending the environment
The provided docker-compose files can be extended to include additional services or configurations as needed. The general guidelines are:
* If you want a new service or initialization process, create it in a directory inside the `local` directory and add it to the `lps-local.sh` script.
* Use the `lps-local.sh` as entrypoint only. Avoid adding logic there. If you need to add logic, do it inside the container you're creating and call it from the `lps-local.sh` script.
* Use the `sample-config.env` file as a reference for the environment variables that you can use in your new service or initialization process. You can also add new environment variables to the `sample-config.env` file if needed.
* Make sure to update the documentation with the new service or initialization process and how to use it.
* If you need to add a new service that depends on the existing services, make sure to add the necessary dependencies in the docker-compose files and update the `lps-local.sh` script accordingly.
* Avoid installing new dependencies in the host machine. Use docker images to run the necessary scripts and services. If you need to install new dependencies, create a new docker image for it and use it in the `lps-local.sh` script.
