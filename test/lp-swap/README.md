# Liquidity Provider Swap environment
This is a simple environment for simulating a liquidity provider swap.
This is the process where one provider is added to the Liquidity Bridge Contract to coexist with the existing provider for some time.
After that time, the old provider is removed and the new provider is the only one providing liquidity.
This environment can be used to perform all the required validations over this process.

## How to run
1. Complete the missing variables in [sample-config.env](../../sample-config.env) since this file will be used as a template for the LPs env files.
Most likely the changes that you will need to do are setting the captcha keys (you can use the test ones) and setting the GITHUB_TOKEN in case the
environment is using any internal version.
2. Create a `config.json` with the same structure as `config.example.json` in the lp-swap folder (the same one where this file is located).
3. Run the following command:
```bash
./setup.sh
```
This script will set up the whole environment for you and will stop when both LPs are running. The user will be able to use both providers until he decides
to sunset one of them, at that point he just need to specify the ID in the console (the script will be waiting for that ID). After that, the script will
proceed to resign the LP and withdraw its collateral from the contract.

### Miners
Most of the operations that the LP needs to perform require a miner to be running so the transactions can achieve certain number of confirmations.
There is one script for each network; [btc-miner.sh](btc-miner.sh) and [rsk-miner.sh](rsk-miner.sh). These scripts will start a miner for the specified local network.
They only require one argument which is the number of seconds between blocks. E.g.:
```bash
./btc-miner.sh 60
```

## Clarifications
- This environment doesn't include the powpeg nodes so the LPs won't be refunded for the peins  or pegouts done in it.
