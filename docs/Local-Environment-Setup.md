# Flyover Local Environment Setup Tutorial

## Overview

This tutorial guides you through setting up and running the Liquidity Provider Server (LPS) in a local development environment on macOS.

## Prerequisites

Before starting, ensure you have the following installed:

- Docker Desktop
- Python 3.8 or higher
- Go 1.23 or higher (check `go.mod` for exact version)
- Git
- Make
- MongoDB Compass
- MetaMask browser extension configured for local development

**System requirements:**
- macOS
- At least 8GB RAM (16GB recommended)
- At least 10GB free disk space
- Multi-core processor

---

## Set Up Local Environment

### 1. Clone the repository

```bash
git clone https://github.com/rsksmart/liquidity-provider-server.git
cd liquidity-provider-server
```

### 2. Switch to tutorial branch

```bash
git checkout tutorial/local_env
```

### 3. Create Python virtual environment

Strongly recommended to avoid dependency conflicts. While some systems may work without it, many users experience issues with dependency installation and pre-commit tools if they skip this step.

```bash
python3 -m venv .venv
source .venv/bin/activate
```

You should see `(.venv)` in your terminal prompt, indicating the virtual environment is active.

If you encounter issues during `make tools`, ensure you created and activated the virtual environment.

### 4. Install development tools

Install all required development tools:

```bash
make tools
```

This installs:
- go-swagger3: API documentation generator
- govulncheck: Go vulnerability scanner
- pre-commit: Git hooks framework for code quality
- golangci-lint: Comprehensive Go linter
- mockery: Mock generation tool for testing

### 5. Clean Docker environment (Recommended)

**Important for troubleshooting:** This step is highly recommended if you've run the LPS before or are experiencing issues. If this is your **first time** setting up the environment, you can skip this step, but make sure the required ports are available.

**Required ports that must be free:**
- `8080` - LPS HTTP server
- `4444` - RSKj RPC endpoint
- `5555` - Bitcoind RPC endpoint
- `27017` - MongoDB
- `4566` - Localstack
- `4450`, `4451` - Powpeg nodes

**When to use Docker cleanup:**
- You've run the LPS environment before and want a fresh start
- You're experiencing port conflicts or container startup issues
- Services aren't starting correctly or behaving unexpectedly
- You want to ensure a completely clean state

**Why this helps:**
- Prevents port conflicts from previous containers
- Ensures fresh blockchain state
- Allows proper contract deployment
- Eliminates stale data that can cause synchronization issues

**If you need to clean your Docker environment, run these commands:**

**Step 1: Stop and remove all containers**

```bash
# Stop all running containers
docker stop $(docker ps -aq) 2>/dev/null || true

# Remove all containers
docker rm $(docker ps -aq) 2>/dev/null || true
```

**Step 2: Remove Docker resources**

```bash
# Remove all images, networks, and build cache
docker system prune -a -f --volumes
```

**Step 3: Remove bind mount data**

This project uses bind mounts (not Docker volumes) to store blockchain data in `docker-compose/local/volumes/`. Docker prune won't remove these directories, so you need to delete them manually:

```bash
# Remove local volumes directory entirely
rm -rf docker-compose/local/volumes
```

This removes all blockchain state and database files stored locally.

**Note:** The `2>/dev/null || true` handles cases where there are no containers to stop or remove. If you see "Error: No such container", that's expected and means the cleanup is working.

### 6. Start the local environment

Navigate to the docker-compose directory and start the environment:

```bash
cd docker-compose/local
export LPS_STAGE=regtest
./lps-env.sh up
```

**What happens during startup:**

1. **Environment setup** (~5 seconds)
   - Creates configuration from `sample-config.env`
   - Enables Management API
   - Sets up directory structure

2. **Base services** (~30 seconds)
   - Starts MongoDB, Localstack, Bitcoind, and RSKj
   - Waits for services to become responsive

3. **Bitcoin wallet setup** (~1-2 minutes)
   - Creates Bitcoin wallet with initial blocks
   - Funds the liquidity provider's address

4. **Liquidity Bridge contract deployment** (~30 seconds)
   - Deploys the LBC smart contract
   - Look for this key message:
     ```
     LBC deployed at 0xefb80db9e2d943a492...
     ```

5. **Powpeg nodes** (~1 minute)
   - Starts federation nodes for peg operations

6. **LPS build and start** (~2-3 minutes)
   - Builds the LPS Docker image
   - Starts the LPS server
   - Performs internal bootstrapping

7. **Automated configuration** (~15 seconds)
   - Configures peg-in and peg-out parameters
   - Creates trusted account for testing

**Expected duration:** 5-15 minutes total, depending on hardware.
- Fast machines (SSD, 16GB+ RAM): ~5-7 minutes
- Average machines: ~10-15 minutes

**Success indicators (watch for these messages):**

```bash
RskJ is up and running
Bitcoind is up and running
LBC deployed at 0xefb80db9e2d943a492...
LPS is up and running
management_password.txt found. Proceeding with configuration.
Trusted account created successfully!
```

**At this point, your local environment is fully up and running.**

You should see all containers running in Docker Desktop:
- `lps01`
- `rskj01`
- `bitcoind01`
- `mongo01`
- `localstack`
- `powpeg-pegin`
- `powpeg-pegout`

### Verify everything is working

Check the health of your LPS:

```bash
curl http://localhost:8080/health
```

Expected output:

```json
{
  "status": "ok",
  "services": {
    "db": "ok",
    "rsk": "ok",
    "btc": "ok"
  }
}
```

---

## ✅ Milestone 1: Local environment running

Your local LPS environment is now fully operational!

---

## Management Console Access

Now that your environment is running, let's access the Management Console.

**What is the Management Console?**

The Management Console is a web-based interface that allows you to manage your Liquidity Provider Server. Through this console, you can:
- Configure peg-in and peg-out parameters (fees, limits, confirmations)
- Add or remove trusted accounts
- View your provider registration details

It's the main control panel for operating your LP and requires authentication to prevent unauthorized access.

### 1. Retrieve the Initial Password

The LPS generates a one-time password on startup. Retrieve it with this command:

```bash
docker exec lps01 cat /tmp/management_password.txt
```

You'll see output similar to:

```
FQW3V7NYDL4X2AIIWBCFZTRLXEMOKYNOM5D3DHXO5YD55BCYG7RGOZ4ULWB7UBMQZENAYAXYLDRSXJ64RPFA3ZQ7GSVIKR6STYBUI4Q
```

Copy this password - you'll need it in the next step.

### 2. Open the Management Console

Open your browser and navigate to:

```
http://localhost:8080/management
```

### 3. Complete First-Time Login

You'll see a login form with fields for both initial and new credentials on the same screen.

Fill in the form as follows:

**Current Credentials (top section):**
- **Username**: `admin`
- **Password**: *[paste the password you retrieved in step 1]*

**New Credentials (bottom section):**
- **New Username**: `admin`
- **New Password**: `Password0*`

**Note:** You can use any password that meets the requirements, Password0* is just an example.

Click **Login** or **Submit**.

### 4. Login with New Credentials

After completing the first-time setup, you'll be redirected to a regular login screen with just two fields:

- **Username**: `admin` (the username you just created)
- **Password**: `Password0*` (the password you just created)

Enter your newly created credentials and click **Login**.

### 5. You're In!

After logging in, you'll be in the Management Console dashboard.

From here you can:
- View your liquidity provider information
- Configure fees and limits
- Manage peg-in and peg-out operations
- Manage trusted accounts

### 6. Configure Pegin Limits

Before creating pegin quotes, configure the accepted value range:

1. In the Management Console, navigate to **Configuration** → **Pegin**
2. Set the following values:
   - **Min Value**: `0.5`
   - **Max Value**: `1`
3. Save the changes

These values define the minimum and maximum BTC amounts that will be accepted for pegin quotes.

---

## ✅ Milestone 2: Management Console Access Complete

Congratulations! You've successfully accessed the Management Console and set up your credentials.

---

## Configure Mining and Liquidity

For the LPS to process transactions, we need to continuously mine blocks on both Bitcoin and Rootstock networks, and ensure both the LPS and your wallet have sufficient liquidity. We'll set up automated miners and fund the necessary wallets.

### 1. Bitcoin Block Miner

Open a new terminal window and run the following commands:

**Set up the Bitcoin CLI alias:**

```bash
alias bi='bitcoin-cli -rpcport=5555 -rpcuser=test -rpcpassword=test -rpcconnect=127.0.0.1'
```

**Start the automatic miner** (mines a block every 2 seconds):

```bash
while true; do bi -rpcwallet=main generatetoaddress 1 mni1YpzHTXrrTtP2AVzwzdkTY6ni5uhJ3U && bi -rpcwallet=main -named sendtoaddress fee_rate=25 address=mni1YpzHTXrrTtP2AVzwzdkTY6ni5uhJ3U amount=0.00001 && sleep 2; done
```

**Note for slower machines:** Mining every 2 seconds can impact performance on slower machines. If you experience slowdowns, you can mine blocks manually instead.

**To mine blocks manually** (example: mine 20 blocks):

```bash
bi -rpcwallet=main -generate 20
```

### 2. Rootstock Block Miner

Open another new terminal window and run:

**Start the automatic miner** (mines a block every 2 seconds):

```bash
while true; do curl --location 'http://localhost:4444' \
--header 'Content-Type: application/json' \
--data '{
    "method": "evm_mine",
    "params": [],
    "id": 1,
    "jsonrpc": "2.0"
}' && sleep 2; done
```

**To mine a block manually:**

```bash
curl --location 'http://localhost:4444' \
--header 'Content-Type: application/json' \
--data '{
    "method": "evm_mine",
    "params": [],
    "id": 1,
    "jsonrpc": "2.0"
}'
```

**For this tutorial:** Keep both mining processes running continuously in their respective terminals. This ensures the blockchains progress and transactions are confirmed.

### 3. Fund Your MetaMask Wallet

You'll need a MetaMask wallet with RBTC to interact with the Flyover protocol. If you don't have MetaMask configured for the local Rootstock network, create and configure it now.

Once your MetaMask wallet is ready, fund it with RBTC using the following command (replace the address with your MetaMask wallet address):

```bash
curl -s -X POST "http://127.0.0.1:4444" \
-H "Content-Type: application/json" \
-d '{"jsonrpc":"2.0","method":"eth_sendTransaction","params": [{"from": "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826", "to": "0x1538283abbD198DcD966f43230363A68108c6373", "value": "0x21e19e0c9bab2400000"}],"id":1}' | jq
```

**Important:**
- Replace `0x1538283abbD198DcD966f43230363A68108c6373` with your MetaMask wallet address
- The `value` parameter (`0x21e19e0c9bab2400000`) represents 10,000 RBTC in wei (hexadecimal)

**Verify the balance:**

```bash
curl -s -X POST "http://127.0.0.1:4444" \
-H "Content-Type: application/json" \
-d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0x1538283abbD198DcD966f43230363A68108c6373","latest"],"id":1}' | jq
```

Replace the address with your MetaMask wallet address. The result will be in wei (hexadecimal).

---

## ✅ Milestone 3: Mining and Liquidity Setup Complete

Both Bitcoin and Rootstock miners are running, and you've deposited liquidity to the LPS wallet and funded your MetaMask wallet!

---

## Create a Pegin Quote

Now let's create a pegin quote (BTC → RBTC) by requesting, accepting, and paying for it.

### 1. Connect MongoDB Compass

Before creating quotes, connect MongoDB Compass to visualize the database changes:

1. Open MongoDB Compass
2. Create a new connection with the following URI:
   ```
   mongodb://root:root@localhost:27017/admin
   ```
3. Click **Connect**
4. Navigate to the `flyover` database

Keep MongoDB Compass open to observe the quote records as they're created.

### 2. Get a Pegin Quote

Request a quote from the LPS to convert BTC to RBTC.

**Important:** Replace the addresses in the command below with your MetaMask wallet address. This ensures the RBTC will be sent to your wallet after the pegin completes.

In your **main terminal**, run:

```bash
curl -X POST 'http://localhost:8080/pegin/getQuote' \
-H 'Content-Type: application/json' \
-d '{
    "callEoaOrContractAddress": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
    "rskRefundAddress": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
    "valueToTransfer": 612345678900000000,
    "callContractArguments": "0x"
}' | jq
```

**Note:** Replace both address fields (`0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826`) with your MetaMask address.

**Response:** You'll receive a quote with details like this:

```json
[
    {
        "quote": {
            "fedBTCAddr": "2N5muMepJizJE1gR7FbHJU6CD18V3BpNF9p",
            "lbcAddr": "0xefb80db9e2d943a492bd988f4c619495ca815643",
            "lpRSKAddr": "0x9d93929a9099be4355fc2389fbf253982f9df47c",
            "btcRefundAddr": "mfWxJ45yp2SFn7UciZyNpvDKrzbhyfKrY8",
            "rskRefundAddr": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
            "lpBTCAddr": "n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6",
            "callFee": 2220740740370000,
            "penaltyFee": 1000000000000000,
            "contractAddr": "0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826",
            "data": "",
            "gasLimit": 21000,
            "nonce": 3828365102111464907,
            "value": 612345678900000000,
            "agreementTimestamp": 1764197556,
            "timeForDeposit": 3600,
            "lpCallTime": 7200,
            "confirmations": 10,
            "callOnRegister": false,
            "gasFee": 0,
            "productFeeAmount": 0
        },
        "quoteHash": "1f4d284f9e4fb2f93b21176a49ac2140a0f2517b933596591f5e37dcb0c7b595"
    }
]
```

**Important:** Copy the `quoteHash` value - you'll need it for the next step.

**Check MongoDB Compass:** After running this command, refresh the `flyover/peginQuote` collection in MongoDB Compass. You should see a new document with your quote details.

### 3. Accept the Pegin Quote

Now accept the quote to proceed with the pegin operation.

In your **main terminal**, run:

```bash
curl -X POST 'http://localhost:8080/pegin/acceptQuote' \
-H 'Content-Type: application/json' \
-d '{
    "QuoteHash": "1f4d284f9e4fb2f93b21176a49ac2140a0f2517b933596591f5e37dcb0c7b595"
}'
```

**Note:** Replace the `QuoteHash` value with the one you received from the previous step.

**Response:** You'll receive acceptance confirmation:

```json
{
    "signature": "fc4738590489c91888d64c3aa9d96baf917fba88d56d53801eb889fb11b4bef84660d569cdd8614d0be4ea43dc9c34b7c0b0621b833eb49d32c97beea1445a691c",
    "bitcoinDepositAddressHash": "2N5D6gUCrgUD5aFPmTAJsdSiEW1bAewGtU5"
}
```

**Important:** Copy the `bitcoinDepositAddressHash` - you'll need it for the payment.

**Check MongoDB Compass:** After running this command, refresh the `flyover/retainedPeginQuote` collection in MongoDB Compass. You should see a new document representing your accepted quote with its initial state.

### 4. Pay the Pegin Quote

Now send Bitcoin to the deposit address to fulfill the quote.

First, set up the Bitcoin CLI alias in your **main terminal** if you haven't already:

```bash
alias bi='bitcoin-cli -rpcport=5555 -rpcuser=test -rpcpassword=test -rpcconnect=127.0.0.1'
```

Then send the payment:

```bash
bi -rpcwallet=main -named sendtoaddress fee_rate=25 address=2N5D6gUCrgUD5aFPmTAJsdSiEW1bAewGtU5 amount=0.7
```

**Note:** Replace the `address` value with the `bitcoinDepositAddressHash` you received in the previous step.

**About the amount:** The `amount` parameter (0.7 BTC in this example) should cover the value you requested in the quote plus fees. This ensures the transaction can be processed successfully.

### 5. Monitor the Pegin Process

After making the payment, the pegin will go through several state changes. As Bitcoin blocks are mined, the LPS will detect the payment and process the pegin.

**Monitor in MongoDB Compass:**
1. Keep the `retainedPeginQuote` collection open and refresh it periodically
2. Find your quote document (search by the quote hash)
3. Watch the `state` field change as the process progresses

**Expected state transitions (happy path):**
1. `WaitingForDeposit` - Initial state after accepting the quote, waiting for you to send BTC
2. `WaitingForDepositConfirmations` - BTC payment detected, waiting for sufficient Bitcoin confirmations
3. `CallForUserSucceeded` - LPS has successfully sent RBTC to your MetaMask address
4. `RegisterPegInSucceeded` - Pegin has been registered with the bridge contract for LP refund

**Wait for completion:** With your Bitcoin and Rootstock miners running, this process will take several minutes as blocks are mined and confirmations accumulate. Be patient and keep refreshing MongoDB Compass to see the `state` field change through these stages.

**Speed up the process:** If the automatic miners are running slowly (every 2 seconds), you can manually mine additional blocks to accelerate the process:

For Bitcoin blocks (in your **main terminal**):

```bash
bi -rpcwallet=main -generate 10
```

For Rootstock blocks (in your **main terminal**):

```bash
curl --location 'http://localhost:4444' \
--header 'Content-Type: application/json' \
--data '{"method": "evm_mine", "params": [], "id": 1, "jsonrpc": "2.0"}'
```

Run these commands multiple times to generate more blocks and meet the confirmation requirements faster.

**Verify in MetaMask:** Once the state reaches `CallForUserSucceeded`, check your MetaMask wallet - you should see the RBTC balance increase!

**Alternative - Check balance via command line:**

You can also verify your RBTC balance using a curl command:

```bash
curl -X POST "http://127.0.0.1:4444" \
-H "Content-Type: application/json" \
-d '{"jsonrpc":"2.0","method":"eth_getBalance","params":["0xYourMetaMaskAddress","latest"],"id":1}' | jq
```

**Note:** Replace `0xYourMetaMaskAddress` with your actual MetaMask wallet address.

The response shows the balance in wei (hexadecimal format):

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": "0x21e19e0c9bab2400000"
}
```

The `result` field is in wei. To convert to RBTC, use an online hex-to-decimal converter and divide by 10^18.

---

## ✅ Milestone 4: Pegin Quote Completed

You've successfully created, accepted, paid for, and completed a pegin! You should see the state changes in MongoDB Compass and receive RBTC in your MetaMask wallet.

---

## Create a Pegout Quote

Now let's create a pegout quote (RBTC → BTC) to convert your RBTC back to Bitcoin.

### 1. Get a Pegout Quote

Request a pegout quote from the LPS.

In your **main terminal**, run:

```bash
curl -X POST 'http://localhost:8080/pegout/getQuotes' \
-H 'Content-Type: application/json' \
-d '{
    "to": "n2NxBRs6aDvuL3qZ3e6vDxWst6wGCyLpHK",
    "rskRefundAddress": "0x45400C53eBd0853Cd26b21C3d479f0eedc46bc44",
    "valueToTransfer": 600000000000000000
}' | jq
```

**Note:**
- Replace `to` with your Bitcoin address where you want to receive BTC
- Replace `rskRefundAddress` with your MetaMask address for refunds
- The `valueToTransfer` is in wei (60000000000000000 wei = 0.06 RBTC)

**Response:** You'll receive a quote similar to this:

```json
[
    {
        "quote": {
            "lbcAddress": "0xefb80db9e2d943a492bd988f4c619495ca815643",
            "liquidityProviderRskAddress": "0x9d93929a9099be4355fc2389fbf253982f9df47c",
            "btcRefundAddress": "n2NxBRs6aDvuL3qZ3e6vDxWst6wGCyLpHK",
            "rskRefundAddress": "0x45400C53eBd0853Cd26b21C3d479f0eedc46bc44",
            "lpBtcAddr": "n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6",
            "callFee": 2180000000000000,
            "penaltyFee": 1000000000000000,
            "nonce": 1458658024694514325,
            "depositAddr": "n2NxBRs6aDvuL3qZ3e6vDxWst6wGCyLpHK",
            "value": 600000000000000000,
            "agreementTimestamp": 1764199603,
            "depositDateLimit": 1764203203,
            "depositConfirmations": 20,
            "transferConfirmations": 10,
            "transferTime": 3600,
            "expireDate": 1764210403,
            "expireBlocks": 524,
            "gasFee": 67250000000000,
            "productFeeAmount": 0
        },
        "quoteHash": "3ddaaec152406a6094927769866a8682e4ffda54eee5e17f86cc14e37bf7ff11"
    }
]
```

**Important:** Save the following values - you'll need them for the next steps:
- `quoteHash` - For accepting the quote
- `quote.lbcAddress` - For payment
- `quote.value` - For payment

**Check MongoDB Compass:** Refresh the `pegoutQuote` collection to see your new quote.

### 2. Accept the Pegout Quote

Accept the quote to proceed with the pegout.

In your **main terminal**, run:

```bash
curl -X POST 'http://localhost:8080/pegout/acceptQuote' \
-H 'Content-Type: application/json' \
-d '{
    "QuoteHash": "3ddaaec152406a6094927769866a8682e4ffda54eee5e17f86cc14e37bf7ff11"
}'
```

**Note:** Replace the `QuoteHash` with the one from your previous response.

**Response:** You'll receive the signature needed for payment:

```json
{
    "signature": "d64a6294782dd3e03d47e13c5f69a34b58400946fa3e2dad083fb9c7f5c813552acd01eac935eaebf53cfb02db7ff693db052892882e06d3289fb7a5ad2a4cd11b",
    "lbcAddress": "0xefb80db9e2d943a492bd988f4c619495ca815643"
}
```

**Important:** Save the `signature` value for the payment step.

**Check MongoDB Compass:** Refresh the `retainedPegoutQuote` collection to see your accepted quote.

### 3. Pay the Pegout Quote

To make payment easier, we've created a Node.js helper script in the `pegoutPayment/` folder.

**First time setup:**

```bash
cd pegoutPayment
npm install
```

**Configure the payment:**

Open `pegoutPayment/index.js` and paste your quote data from the `getQuotes` response into the `QUOTE` object:

```javascript
const QUOTE = {
    "lbcAddress": "0x03f23ae1917722d5a27a2ea0bcc98725a2a2a49a",
    "liquidityProviderRskAddress": "0x9d93929a9099be4355fc2389fbf253982f9df47c",
    "btcRefundAddress": "n2NxBRs6aDvuL3qZ3e6vDxWst6wGCyLpHK",
    "rskRefundAddress": "0x45400C53eBd0853Cd26b21C3d479f0eedc46bc44",
    "lpBtcAddr": "n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6",
    "callFee": "2180000000000000",
    "penaltyFee": "1000000000000000",
    "nonce": "3504134078398023607",  // IMPORTANT: Must be string
    "depositAddr": "n2NxBRs6aDvuL3qZ3e6vDxWst6wGCyLpHK",
    "value": "600000000000000000",
    // ... other fields from your quote
};

const SIGNATURE = "your_signature_from_acceptQuote";
```

**Important:** 
- Large numbers (nonce, value, fees) must be quoted as strings to avoid JavaScript precision loss.
- Update `PAYMENT_ADDRESS` in the script to your MetaMask address (the one you funded earlier with RBTC).

**Run the payment:**

```bash
npm run pay
```

The script will:
1. Properly encode the quote struct with Base58-decoded BTC addresses
2. Send the `depositPegout` transaction to the LBC contract
3. Display the transaction hash

### 4. Monitor the Pegout Process

After payment, the pegout will progress through several states as Rootstock blocks are mined.

**Monitor in MongoDB Compass:**
1. Refresh the `retainedPegoutQuote` collection
2. Find your quote (search by quote hash)
3. Watch the `state` field change

**Expected state transitions (happy path):**
1. `WaitingForDeposit` - Initial state, waiting for RBTC payment
2. `WaitingForDepositConfirmations` - Payment detected, waiting for confirmations
3. `SendPegoutSucceeded` - LPS has sent BTC to your Bitcoin address
4. `RefundPegOutSucceeded` - LP has been refunded from the bridge

**Note:** Similar to pegin, there are additional states after the LP completes the bridge rebalancing. For this tutorial, we'll focus on the flow up to `SendPegoutSucceeded`, which means you've received your BTC!

**Speed up the process:** Just like with pegin, you can manually mine blocks to accelerate:

For Rootstock blocks:

```bash
curl --location 'http://localhost:4444' \
--header 'Content-Type: application/json' \
--data '{"method": "evm_mine", "params": [], "id": 1, "jsonrpc": "2.0"}'
```

For Bitcoin blocks:

```bash
bi -rpcwallet=main -generate 10
```

**Verify BTC received:** Once the state reaches `SendPegoutSucceeded`, check your Bitcoin address balance:

```bash
bi -rpcwallet=main getreceivedbyaddress n2NxBRs6aDvuL3qZ3e6vDxWst6wGCyLpHK
```

Replace the address with your Bitcoin address from the quote.

---

## ✅ Milestone 5: Pegout Quote Completed

You've successfully created, accepted, paid for, and completed a pegout! You should see the BTC in your Bitcoin address.

---

**Tutorial Version**: 1.0  
**Last Updated**: November 26, 2025  
**Applicable Branch**: QA-Test
