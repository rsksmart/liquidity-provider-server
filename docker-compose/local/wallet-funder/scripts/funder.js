#!/bin/node

const env = process.env
const BTC_URL = `http://${env.BTC_ENDPOINT}`;
const RSK_URL = env.RSK_ENDPOINT;
const BTC_AUTH = 'Basic ' + Buffer.from(`${env.BTC_USERNAME}:${env.BTC_PASSWORD}`).toString('base64');

const fundedRskWallets = JSON.parse(env.FUNDED_RSK_WALLETS);
const fundedBtcWallets = JSON.parse(env.FUNDED_BTC_WALLETS);
const RSK_FUNDER = "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826";
const btcFunderPassword = "test-password";
const btcFunderName = "main";
const btcFundingBlocks = 500;

async function createBitcoinFunderWallet({name, password, fundingBlocks}) {
    const created = await fetch(BTC_URL+"/wallet/"+name, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': BTC_AUTH,
        },
        body: JSON.stringify({
            jsonrpc: '1.0',
            method: 'getwalletinfo',
            params: [],
            id: 1
        })
    }).then(response => response.ok)

    if (created) {
        console.log('Bitcoin funder wallet already exists, skipping creation');
        return;
    }

    console.log('Creating Bitcoin funder wallet...');
    return fetch(BTC_URL, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': BTC_AUTH,
        },
        body: JSON.stringify({
            jsonrpc: '1.0',
            method: 'createwallet',
            params: [name, false, false, password, true, false, true],
            id: 1
        })
    }).then(response => {
        if (!response.ok) {
            throw new Error(`Error creating Bitcoin funder wallet! status: ${response.status}`);
        }
        return fetch(BTC_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': BTC_AUTH,
            },
            body: JSON.stringify({
                jsonrpc: '1.0',
                method: 'getnewaddress',
                params: [name],
                id: 1
            })
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Error getting new address for Bitcoin funder wallet! status: ${response.status}`);
        }
        return response.json();
    }).then(data => {
        const address = data.result;
        return fetch(BTC_URL, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': BTC_AUTH,
            },
            body: JSON.stringify({
                jsonrpc: '1.0',
                method: 'generatetoaddress',
                params: [fundingBlocks, address],
                id: 1
            })
        })
    }).then(response => {
        if (!response.ok) {
            throw new Error(`Error generating Bitcoin blocks for funder wallet! status: ${response.status}`);
        }

    })
}

function fundRskWallet(address, amount) {
    return fetch(RSK_URL, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            jsonrpc: '2.0',
            method: 'eth_sendTransaction',
            params: [{
                from: RSK_FUNDER,
                to: address,
                value: '0x' + BigInt(amount).toString(16),
                id: 1
            }],
        })
    }).then(response => {
        if (!response.ok) {
            throw new Error(`Error funding wallets! status: ${response.status}`);
        }
    })
}

function fundBtcAddress({funderPassword, funderName, address, amount}) {
    return fetch(BTC_URL+"/wallet/"+funderName, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': BTC_AUTH,
        },
        body: JSON.stringify({
            jsonrpc: '1.0',
            method: 'walletpassphrase',
            params: [funderPassword, 60],
            id: 1
        })
    }).then(response => {
        if (!response.ok) {
            throw new Error(`Error unlocking funder wallet! status: ${response.status}`);
        }
        return fetch(BTC_URL+"/wallet/"+funderName, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': BTC_AUTH,
            },
            body: JSON.stringify({
                jsonrpc: '1.0',
                method: 'sendtoaddress',
                params: {
                    amount: amount,
                    fee_rate: 25,
                    address: address
                },
                id: 1
            })
        })
    }).then(async response => {
        if (!response.ok) {
            console.log(await response.json())
            throw new Error(`Error funding Bitcoin address ${address}! status: ${response.status}`);
        }
        return fetch(BTC_URL+"/wallet/"+funderName, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': BTC_AUTH,
            },
            body: JSON.stringify({
                jsonrpc: '1.0',
                method: 'generatetoaddress',
                params: [1, address],
                id: 1
            })
        })
    }).then(response => {
        if (!response.ok) {
            throw new Error(`Error confirming funding tx to address ${address}! status: ${response.status}`);
        }
    })
}

try {
    await createBitcoinFunderWallet({name: btcFunderName, password: btcFunderPassword, fundingBlocks: btcFundingBlocks});
    for (const [address, amount] of Object.entries(fundedRskWallets)) {
        await fundRskWallet(address, amount);
        console.log(`Funded wallet ${address} with ${amount} wei`);
    }
    for (const [address, amount] of Object.entries(fundedBtcWallets)) {
        await fundBtcAddress({
            funderPassword: btcFunderPassword,
            funderName: btcFunderName,
            address,
            amount
        });
        console.log(`Funded Bitcoin address ${address} with ${amount} btc`);
    }
} catch (error) {
    console.error(error);
    process.exit(1);
}
