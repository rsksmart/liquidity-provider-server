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
 * SVP (RSKIP-419) flow (VETIVER / 9.0.3+):
 *   After commitFederation the bridge requires spend-validation before the new federation
 *   activates. The federators must sign + broadcast a fund tx and a spend tx to BTC, both
 *   of which must confirm on-chain. Only then does the bridge accept the proposed federation
 *   and count down federationActivationAge blocks.
 *
 *   Pre-conditions for SVP to succeed:
 *     1. The bridge must have synced BTC headers (getBtcBlockchainBestChainHeight > 0) before
 *        commit, or it cannot build the svpFundTx.
 *     2. The active (genesis) federation must hold BTC, or the fund tx fails with
 *        INSUFFICIENT_MONEY.
 *
 *   This script satisfies both conditions before committing:
 *     a. primeBtcRelay — interleaves RSK + BTC block production until the bridge BTC height
 *        catches up to the bitcoind tip.
 *     b. fundFederation — sends BTC to the genesis federation address so SVP has funds.
 *
 * Steps:
 *   1.  Wait for the RSK node to be reachable.
 *   2.  Fund auth-a/b/c from the pre-unlocked coinbase so they can pay for gas.
 *   3.  Vote createFederation() — 3 times to reach majority.
 *   4.  Vote addFederatorPublicKeyMultikey() for each of the 3 federation members — 3 times each.
 *   5.  Call getPendingFederationHash() to retrieve the hash for the commit step.
 *   6.  Prime BTC header relay, then fund the active (genesis) federation.
 *   7.  Vote commitFederation(hash) — 3 times, then drive SVP + activation via interleaved
 *       BTC/RSK block production until the federation address changes.
 *   8.  Verify the active federation address changed (confirms activation succeeded).
 *   9.  Write USE_SEGWIT_FEDERATION=true to the env file so LPS starts with the right flag.
 */

const fs = require('fs');
const { ethers } = require('ethers');

// ── Environment ───────────────────────────────────────────────────────────────

const RSK_URL = process.env.RSK_ENDPOINT;
if (!RSK_URL) throw new Error('RSK_ENDPOINT environment variable is required');

const ENV_FILE = process.env.ENV_FILE;
if (!ENV_FILE) throw new Error('ENV_FILE environment variable is required');
const ENV_FILE_PATH = '/' + ENV_FILE;

const BTC_RPC_URL = process.env.BTC_RPC_URL;
if (!BTC_RPC_URL) throw new Error('BTC_RPC_URL environment variable is required');

const BTC_RPC_USER = process.env.BTC_RPC_USER || 'test';
const BTC_RPC_PASSWORD = process.env.BTC_RPC_PASSWORD || 'test';
const BTC_RPC_WALLET = process.env.BTC_RPC_WALLET || 'main';
const BTC_WALLET_PASSPHRASE = process.env.BTC_WALLET_PASSPHRASE || '';

// ── Constants ─────────────────────────────────────────────────────────────────

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

// BTC to send to the active federation so the bridge can build the SVP fund tx.
// Must cover the 500000 sat SVP amount plus fees.
const FEDERATION_FUNDING_BTC = 5;
// BTC blocks to mine after the federation funding tx so federators register the UTXO.
const FEDERATION_FUNDING_CONFIRMATIONS = 20;

// Relay-priming: rounds of (3 RSK blocks + 2 BTC blocks + 2 s pause) before checking bridge height.
const RELAY_PRIME_MAX_ROUNDS = 30;
// Bridge BTC height must be within this many blocks of the bitcoind tip to consider relay primed.
const RELAY_PRIME_TOLERANCE = 5;

// SVP activation loop: interleave BTC + RSK blocks with wall-clock pauses so the federators
// can sign/broadcast the fund and spend txs and the bridge can validate them.
const SVP_MAX_ROUNDS = 220;
const BTC_BLOCKS_PER_ROUND = 4;
const RSK_BLOCKS_PER_ROUND = 2;
const ROUND_DELAY_MS = 2000;

// ── Bridge ABI ────────────────────────────────────────────────────────────────

const BRIDGE_ABI = [
    'function createFederation() returns (int256)',
    'function addFederatorPublicKeyMultikey(bytes btcKey, bytes rskKey, bytes mstKey) returns (int256)',
    'function getPendingFederationHash() view returns (bytes)',
    'function commitFederation(bytes hash) returns (int256)',
    'function getFederationAddress() view returns (string)',
    'function getBtcBlockchainBestChainHeight() view returns (uint256)',
];

// ── Helpers ───────────────────────────────────────────────────────────────────

const delay = ms => new Promise(r => setTimeout(r, ms));

async function waitForNode(provider, maxRetries = 30, delayMs = 3000) {
    for (let i = 0; i < maxRetries; i++) {
        try {
            const network = await provider.getNetwork();
            console.log(`RSK node is ready (chainId: ${network.chainId})`);
            return;
        } catch {
            console.log(`Waiting for RSK node... (attempt ${i + 1}/${maxRetries})`);
            await delay(delayMs);
        }
    }
    throw new Error('RSK node did not become available in time');
}

async function rskRpcCall(method, params = []) {
    const res = await fetch(RSK_URL, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ jsonrpc: '2.0', method, params, id: 1 }),
    });
    if (!res.ok) throw new Error(`RSK RPC ${method} failed with HTTP ${res.status}`);
    const json = await res.json();
    if (json.error) throw new Error(`RSK RPC ${method} returned error: ${JSON.stringify(json.error)}`);
    return json.result;
}

async function btcRpc(method, params = [], useWallet = false) {
    const url = useWallet
        ? `${BTC_RPC_URL}/wallet/${BTC_RPC_WALLET}`
        : BTC_RPC_URL;
    const credentials = Buffer.from(`${BTC_RPC_USER}:${BTC_RPC_PASSWORD}`).toString('base64');
    const res = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Basic ${credentials}`,
        },
        body: JSON.stringify({ jsonrpc: '1.0', method, params, id: 1 }),
    });
    if (!res.ok) throw new Error(`BTC RPC ${method} failed with HTTP ${res.status}: ${await res.text()}`);
    const json = await res.json();
    if (json.error) throw new Error(`BTC RPC ${method} returned error: ${JSON.stringify(json.error)}`);
    return json.result;
}

async function btcHeight() {
    return btcRpc('getblockcount');
}

async function mineBtc(toAddr, n) {
    return btcRpc('generatetoaddress', [n, toAddr]);
}

async function getBtcMiningAddress() {
    return btcRpc('getnewaddress', [], true);
}

async function getBridgeBtcHeight(provider, bridgeInterface) {
    const result = await provider.call({
        to: BRIDGE_ADDRESS,
        data: bridgeInterface.encodeFunctionData('getBtcBlockchainBestChainHeight'),
    });
    return Number(bridgeInterface.decodeFunctionResult('getBtcBlockchainBestChainHeight', result)[0]);
}

async function getFederationAddress(provider, bridgeInterface) {
    const result = await provider.call({
        to: BRIDGE_ADDRESS,
        data: bridgeInterface.encodeFunctionData('getFederationAddress'),
    });
    return bridgeInterface.decodeFunctionResult('getFederationAddress', result)[0];
}

async function fundAuthAddresses(addresses) {
    for (const addr of addresses) {
        await rskRpcCall('eth_sendTransaction', [{
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

// Mines RSK blocks by sending coinbase self-transfers (autoMine=true mints one block per tx).
async function mineRsk(count) {
    for (let i = 0; i < count; i++) {
        await rskRpcCall('eth_sendTransaction', [{
            from: RSK_FUNDER,
            to: RSK_FUNDER,
            value: '0x0',
        }]);
        await delay(50);
    }
}

// Interleaves RSK + BTC block production until the bridge BTC height is within
// RELAY_PRIME_TOLERANCE blocks of the bitcoind tip.
async function primeBtcRelay(provider, bridgeInterface, btcAddr) {
    const target = await btcHeight();
    console.log(`  BTC tip: ${target}, driving relay to bridge...`);
    for (let i = 0; i < RELAY_PRIME_MAX_ROUNDS; i++) {
        await mineRsk(3);
        await mineBtc(btcAddr, 2);
        await delay(2000);
        const bridgeHeight = await getBridgeBtcHeight(provider, bridgeInterface);
        console.log(`    relay round ${i}: bridge BTC height = ${bridgeHeight} / ${target}`);
        if (bridgeHeight >= target - RELAY_PRIME_TOLERANCE) {
            console.log(`  BTC header relay primed: bridge BTC height = ${bridgeHeight}`);
            return;
        }
    }
    const bridgeHeight = await getBridgeBtcHeight(provider, bridgeInterface);
    console.warn(`  WARNING: relay priming exhausted ${RELAY_PRIME_MAX_ROUNDS} rounds; bridge height = ${bridgeHeight}`);
}

// Sends BTC to the active (genesis) federation so the bridge has funds for the SVP fund tx.
async function fundFederation(provider, bridgeInterface, btcAddr) {
    const fedAddr = await getFederationAddress(provider, bridgeInterface);
    console.log(`  Funding active federation ${fedAddr} with ${FEDERATION_FUNDING_BTC} BTC (peg-in)...`);

    if (BTC_WALLET_PASSPHRASE) {
        await btcRpc('walletpassphrase', [BTC_WALLET_PASSPHRASE, 120], true);
    }
    await btcRpc('settxfee', [0.0001], true);
    const txid = await btcRpc('sendtoaddress', [fedAddr, FEDERATION_FUNDING_BTC], true);
    console.log(`  Sent ${FEDERATION_FUNDING_BTC} BTC to federation, txid: ${txid}`);

    await mineBtc(btcAddr, FEDERATION_FUNDING_CONFIRMATIONS);
    console.log(`  Mined ${FEDERATION_FUNDING_CONFIRMATIONS} BTC blocks to confirm funding`);

    // Drive RSK + BTC interleaved so federators can register the UTXO into the bridge wallet.
    for (let i = 0; i < 12; i++) {
        await mineRsk(2);
        await mineBtc(btcAddr, 2);
        await delay(2000);
    }
    console.log('  Federation funded and UTXO registration driven.');
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
    console.log(`BTC RPC URL  : ${BTC_RPC_URL}`);
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
    console.log('\n[3/9] createFederation (3 votes)...');
    await voteAll(authWallets, bridgeInterface, 'createFederation', [], 'createFederation');

    // ── Step 4: addFederatorPublicKeyMultikey ─────────────────────────────────
    // Add each powpeg member to the pending federation.
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

    const oldFederationAddress = await getFederationAddress(provider, bridgeInterface);
    console.log(`  Current (old) federation address: ${oldFederationAddress}`);

    // Get a BTC address to mine rewards to throughout the rest of the migration.
    const btcMiningAddr = await getBtcMiningAddress();
    console.log(`  BTC mining address: ${btcMiningAddr}`);

    // ── Step 6: Prime relay + fund active federation (before commit) ──────────
    // These must happen BEFORE commitFederation so the SVP validation window is
    // not consumed by pre-conditions work.
    console.log('\n[6/9] Priming BTC header relay before commit...');
    await primeBtcRelay(provider, bridgeInterface, btcMiningAddr);

    console.log('      Funding the active federation for SVP...');
    await fundFederation(provider, bridgeInterface, btcMiningAddr);

    // ── Step 7: commitFederation + SVP driving ────────────────────────────────
    // After commit the bridge starts the SVP validation period. Interleave BTC and
    // RSK block production so the federators can complete the fund→confirm→spend→confirm
    // cycle within the validation window, then count down federationActivationAge.
    console.log('\n[7/9] commitFederation (3 votes)...');
    await voteAll(authWallets, bridgeInterface, 'commitFederation', [pendingHash], 'commitFederation');

    console.log('      Driving SVP validation + activation (interleaved RSK+BTC mining)...');
    let activated = false;
    for (let round = 0; round < SVP_MAX_ROUNDS; round++) {
        await mineBtc(btcMiningAddr, BTC_BLOCKS_PER_ROUND);
        await mineRsk(RSK_BLOCKS_PER_ROUND);
        await delay(ROUND_DELAY_MS);

        if (round % 5 === 0) {
            const bridgeBtcHeight = await getBridgeBtcHeight(provider, bridgeInterface);
            const btcTip = await btcHeight();
            const curFedAddr = await getFederationAddress(provider, bridgeInterface);
            console.log(
                `      round ${String(round).padStart(3)}: fedAddr=${curFedAddr}` +
                ` bridgeBtcHeight=${bridgeBtcHeight} btcTip=${btcTip}`,
            );
            if (curFedAddr !== oldFederationAddress) {
                activated = true;
                break;
            }
        }
    }

    if (!activated) {
        throw new Error(
            `SVP activation did not complete within ${SVP_MAX_ROUNDS} rounds.\n` +
            'Increase SVP_MAX_ROUNDS or check the powpeg federator logs for errors.',
        );
    }

    // ── Step 8: Verify federation activated ──────────────────────────────────
    console.log('\n[8/9] Verifying federation activation...');
    const newFederationAddress = await getFederationAddress(provider, bridgeInterface);
    if (newFederationAddress === oldFederationAddress) {
        throw new Error(
            `Federation address did not change.\n` +
            `  Old address: ${oldFederationAddress}\n` +
            `  New address: ${newFederationAddress}`,
        );
    }
    console.log(`  Federation activation confirmed. New address: ${newFederationAddress}`);

    // ── Step 9: Update env file ───────────────────────────────────────────────
    console.log('\n[9/9] Updating env file...');
    updateEnvFile(ENV_FILE_PATH);

    console.log('\n Migration complete.');
    console.log(`  The new segwit-compatible (P2SH-P2WSH) federation is now active.\n`);
}

main().catch(err => {
    console.error('\nFederation migration failed:', err.message ?? err);
    process.exit(1);
});
