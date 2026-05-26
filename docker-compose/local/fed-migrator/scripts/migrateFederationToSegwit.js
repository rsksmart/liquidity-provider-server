'use strict';

/**
 * Migrates the regtest RSK federation from legacy P2SH to segwit-compatible P2SH-P2WSH.
 *
 * Background:
 *   RSKIP-305 is active from block 0 in regtest (reed800 = 0, reed810 = 0 in regtest.conf). Any
 *   federation committed after genesis is automatically created as P2SH-P2WSH by the bridge.
 *   The genesis federation is legacy P2SH, so a federation change vote is required to
 *   produce the first segwit-compatible federation.
 *
 * Vote mechanism:
 *   Federation changes are voted on by 5 authorized addresses hardcoded in the regtest
 *   bridge configuration. A majority (3/5) is required for each step.
 *   The private keys are derived as keccak256("auth-a"), keccak256("auth-b"),
 *   keccak256("auth-c").
 *
 * Steps:
 *   1. Wait for the RSK node to be reachable.
 *   2. Fund auth-a/b/c from the pre-unlocked coinbase so they can pay for gas.
 *   3. Vote createFederation() — 3 times to reach majority.
 *   4. Vote addFederatorPublicKeyMultikey() for each of the 3 federation members — 3 times each.
 *   5. Call getPendingFederationHash() to retrieve the hash for the commit step.
 *   6. Vote commitFederation(hash) — 3 times to reach majority.
 *   7. Mine enough blocks for the new federation to become active (federationActivationAge=150).
 *   8. Verify the active federation address changed (confirms activation succeeded).
 *   9. Write USE_SEGWIT_FEDERATION=true to the env file so LPS starts with the right flag.
 */

const fs = require('fs');
const { ethers } = require('ethers');

// ── Environment ───────────────────────────────────────────────────────────────

const RSK_URL = process.env.RSK_ENDPOINT;
if (!RSK_URL) throw new Error('RSK_ENDPOINT environment variable is required');

const ENV_FILE = process.env.ENV_FILE;
if (!ENV_FILE) throw new Error('ENV_FILE environment variable is required');
const ENV_FILE_PATH = '/' + ENV_FILE;

// ── Constants ─────────────────────────────────────────────────────────────────

// RSK Bridge precompile address
const BRIDGE_ADDRESS = '0x0000000000000000000000000000000001000006';

// Pre-funded coinbase
const RSK_FUNDER = '0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826';

// 0.1 RBTC per auth address — sufficient for all bridge voting transactions.
const AUTH_FUND_AMOUNT = ethers.parseEther('0.1');

// Federation change authorized private keys derived as keccak256("<generator>")
const AUTH_PRIVATE_KEYS = [
    '0xc14991a187e185ca1442e75eb8f60a6a5efd4ca57ce31e50d6e841d9381e996b', // auth-a
    '0x488cdd0c11d602598225fe96c4b85c2afbec3f1d938cd88f4655831cb6ff454b', // auth-b
    '0x72255947e1aff21d3fc9c077c6a70912aede3674913d5c76b4128c1ec5692499', // auth-c
];

// Regtest powpeg member compressed public keys.
// In regtest, BTC = RSK keys are all derived from the same key file.
// These are the same 3 members as the genesis federation so the composition
// matches and fund migration can proceed cleanly.
//   reg1 → powpeg-pegin  (reg1.key: 45c5b07f...)
//   reg2 → powpeg-pegout (reg2.key: 505334c7...)
//   02cd53fc → third genesis member (same as in the regtest genesis federation)
const FEDERATION_MEMBER_PUBKEYS = [
    '0362634ab57dae9cb373a5d536e66a8c4f67468bbcfb063809bab643072d78a124', // reg1
    '03c5946b3fbae03a654237da863c9ed534e0878657175b132b8ca630f245df04db', // reg2
    '02cd53fc53a07f211641a677d250f6de99caf620e8e77071e811a28b3bcddf0be1', // third genesis member
];

// Blocks required after commitFederation for the new federation to become active:
//   validationPeriodDurationInBlocks (125) + federationActivationAge (150) + safety buffer (5)
const BLOCKS_AFTER_COMMIT = 280;

// ── Bridge ABI ────────────────────────────────────────────────────────────────

const BRIDGE_ABI = [
    'function createFederation() returns (int256)',
    'function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns (int256)',
    'function getPendingFederationHash() view returns (bytes)',
    'function commitFederation(bytes hash) returns (int256)',
    'function getFederationAddress() view returns (string)',
];

// ── Helpers ───────────────────────────────────────────────────────────────────

async function waitForNode(provider, maxRetries = 30, delayMs = 3000) {
    for (let i = 0; i < maxRetries; i++) {
        try {
            const network = await provider.getNetwork();
            console.log(`RSK node is ready (chainId: ${network.chainId})`);
            return;
        } catch {
            console.log(`Waiting for RSK node... (attempt ${i + 1}/${maxRetries})`);
            await new Promise(r => setTimeout(r, delayMs));
        }
    }
    throw new Error('RSK node did not become available in time');
}

async function rpcCall(method, params = []) {
    const res = await fetch(RSK_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ jsonrpc: '2.0', method, params, id: 1 }),
    });
    if (!res.ok) throw new Error(`RPC ${method} failed with HTTP ${res.status}`);
    const json = await res.json();
    if (json.error) throw new Error(`RPC ${method} returned error: ${JSON.stringify(json.error)}`);
    return json.result;
}

async function fundAuthAddresses(addresses) {
    for (const addr of addresses) {
        await rpcCall('eth_sendTransaction', [{
            from: RSK_FUNDER,
            to: addr,
            value: '0x' + AUTH_FUND_AMOUNT.toString(16),
        }]);
        console.log(`  Funded ${addr} with 0.1 RBTC`);
    }
}

async function sendBridgeVote(wallet, data, label) {
    const tx = await wallet.sendTransaction({
        to: BRIDGE_ADDRESS,
        data,
        gasLimit: 200000,
    });
    const receipt = await tx.wait();
    console.log(`  [${label}] from ${wallet.address} — block ${receipt.blockNumber}, tx ${receipt.hash}`);
}

async function voteAll(wallets, bridgeInterface, fnName, args, label) {
    console.log(`Voting: ${label}`);
    const data = bridgeInterface.encodeFunctionData(fnName, args ?? []);
    for (const wallet of wallets) {
        await sendBridgeVote(wallet, data, label);
    }
}

// Uses eth_sendTransaction from the pre-unlocked coinbase to advance the chain.
// autoMine=true mines one block per transaction. A short delay between calls prevents
// the RSK node from failing with "transaction wasn't mined" under rapid load.
async function mineBlocks(count) {
    console.log(`Mining ${count} blocks via coinbase self-transfers...`);
    const delay = ms => new Promise(r => setTimeout(r, ms));
    for (let i = 0; i < count; i++) {
        await rpcCall('eth_sendTransaction', [{
            from: RSK_FUNDER,
            to: RSK_FUNDER,
            value: '0x0',
        }]);
        await delay(50);
        if ((i + 1) % 50 === 0) {
            console.log(`  ${i + 1}/${count} blocks mined`);
        }
    }
    console.log(`  Done — ${count} blocks mined`);
}

function updateEnvFile(filePath) {
    let content = fs.readFileSync(filePath, 'utf8');
    if (!content.includes('USE_SEGWIT_FEDERATION=')) {
        content += '\nUSE_SEGWIT_FEDERATION=true\n';
    } else {
        content = content.replace(/^USE_SEGWIT_FEDERATION=.*$/m, 'USE_SEGWIT_FEDERATION=true');
    }
    fs.writeFileSync(filePath, content, 'utf8');
    console.log(`  Updated ${filePath}: USE_SEGWIT_FEDERATION=true`);
}

// ── Main ──────────────────────────────────────────────────────────────────────

async function main() {
    console.log('╔══════════════════════════════════════════╗');
    console.log('║   Federation Segwit Migration (RSKIP-305) ║');
    console.log('╚══════════════════════════════════════════╝');
    console.log(`RSK endpoint : ${RSK_URL}`);
    console.log(`Env file     : ${ENV_FILE_PATH}`);

    const provider = new ethers.JsonRpcProvider(RSK_URL);
    const bridgeInterface = new ethers.Interface(BRIDGE_ABI);
    const authWallets = AUTH_PRIVATE_KEYS.map(pk => new ethers.Wallet(pk, provider));

    // ── Step 1: Wait for RSK node ─────────────────────────────────────────────
    console.log('\n[1/9] Waiting for RSK node...');
    await waitForNode(provider);

    // ── Step 2: Fund authorized addresses ────────────────────────────────────
    console.log('\n[2/9] Funding authorized addresses from coinbase...');
    await fundAuthAddresses(authWallets.map(w => w.address));

    // ── Step 3: createFederation ──────────────────────────────────────────────
    // Three votes reach the majority threshold (3/5) required by the regtest bridge.
    // The third vote executes the function and creates a pending federation.
    console.log('\n[3/9] createFederation (3 votes)...');
    await voteAll(authWallets, bridgeInterface, 'createFederation', [], 'createFederation');

    // ── Step 4: addFederatorPublicKeyMultikey ─────────────────────────────────
    // Add each powpeg member to the pending federation.
    // In regtest BTC/RSK/MST keys are all the same compressed public key.
    console.log('\n[4/9] addFederatorPublicKeyMultikey for each member (3 votes each)...');
    for (const pubkeyHex of FEDERATION_MEMBER_PUBKEYS) {
        const keyBytes = '0x' + pubkeyHex;
        await voteAll(
            authWallets,
            bridgeInterface,
            'addFederatorPublicKeyMultikey',
            [keyBytes, keyBytes, keyBytes],
            `addMember(${pubkeyHex.slice(0, 14)}...)`,
        );
    }

    // ── Step 5: getPendingFederationHash ──────────────────────────────────────
    console.log('\n[5/9] Fetching pending federation hash...');
    const hashCallResult = await provider.call({
        to: BRIDGE_ADDRESS,
        data: bridgeInterface.encodeFunctionData('getPendingFederationHash'),
    });
    const [pendingHash] = bridgeInterface.decodeFunctionResult('getPendingFederationHash', hashCallResult);
    console.log(`  Pending federation hash: ${pendingHash}`);

    // ── Step 6: commitFederation ──────────────────────────────────────────────
    // After commit: federation activates after federationActivationAge (150 blocks).
    // Because RSKIP-305 is active at block 0, the bridge automatically produces
    // a segwit-compatible P2SH-P2WSH address for any new federation.
    const oldFederationAddress = await provider.call({
        to: BRIDGE_ADDRESS,
        data: bridgeInterface.encodeFunctionData('getFederationAddress'),
    }).then(r => bridgeInterface.decodeFunctionResult('getFederationAddress', r)[0]);
    console.log('\n[6/9] commitFederation (3 votes)...');
    await voteAll(authWallets, bridgeInterface, 'commitFederation', [pendingHash], 'commitFederation');

    // ── Step 7: Mine blocks ───────────────────────────────────────────────────
    // rskj has autoMine=true: one block is mined per incoming transaction.
    // We send self-transfers from RSK_FUNDER to advance the block height.
    console.log(`\n[7/9] Mining ${BLOCKS_AFTER_COMMIT} blocks for federation activation...`);
    await mineBlocks(BLOCKS_AFTER_COMMIT);

    // ── Step 8: Verify federation activated ──────────────────────────────────
    console.log('\n[8/9] Verifying federation activation...');
    const newFederationAddress = await provider.call({
        to: BRIDGE_ADDRESS,
        data: bridgeInterface.encodeFunctionData('getFederationAddress'),
    }).then(r => bridgeInterface.decodeFunctionResult('getFederationAddress', r)[0]);
    if (newFederationAddress === oldFederationAddress) {
        throw new Error(
            `Federation address did not change after ${BLOCKS_AFTER_COMMIT} blocks.\n` +
            `  Old address: ${oldFederationAddress}\n` +
            `  New address: ${newFederationAddress}\n` +
            'The new federation may not have activated yet — increase BLOCKS_AFTER_COMMIT.'
        );
    }
    console.log('  Federation activation confirmed.');

    // ── Step 9: Update env file ───────────────────────────────────────────────
    // lps-local.sh re-sources the env file after this container exits, so LPS
    // will start with USE_SEGWIT_FEDERATION=true.
    console.log('\n[9/9] Updating env file...');
    updateEnvFile(ENV_FILE_PATH);

    console.log('\n Migration complete.');
    console.log('  The new segwit-compatible (P2SH-P2WSH) federation is now active.\n');
}

main().catch(err => {
    console.error('\nFederation migration failed:', err.message ?? err);
    process.exit(1);
});
