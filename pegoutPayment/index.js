/**
 * Pay Pegout Script
 *
 * This script pays for a pegout quote by calling depositPegout on the LBC contract.
 *
 * SETUP:
 *   cd pegoutPayment
 *   npm install
 *
 * USAGE:
 *   1. Paste your quote from getQuotes response into quote.json
 *   2. Paste your signature below (from acceptQuote response)
 *   3. Run: npm run pay
 */

const { ethers } = require('ethers');
const bs58check = require('bs58check');
const fs = require('fs');
const path = require('path');

// ============================================================================
// PASTE YOUR SIGNATURE HERE (from acceptQuote response)
// ============================================================================

const SIGNATURE = "f89638d98b828a4428ea84f8dbdf11af0396769f1d9a1e0cfa8335888384fd13129bd4e255d75d3065d546a46dc52167dade9586ba3d847faa19e4b809a5cf7c1b";

// ============================================================================
// CONFIGURATION
// ============================================================================

const RSK_RPC_URL = "http://127.0.0.1:4444";
const PAYMENT_ADDRESS = "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826";

// ============================================================================
// LOAD QUOTE FROM FILE
// ============================================================================

function loadQuote() {
    const quotePath = path.join(__dirname, 'quote.json');

    if (!fs.existsSync(quotePath)) {
        console.error("❌ ERROR: quote.json not found!");
        console.error("   Paste the response from getQuotes into quote.json");
        process.exit(1);
    }

    // Read as raw text and convert numbers to strings to preserve precision
    const rawJson = fs.readFileSync(quotePath, 'utf8');
    const safeJson = rawJson.replace(/:\s*(-?\d+\.?\d*)\s*([,}\]])/g, ': "$1"$2');

    try {
        const data = JSON.parse(safeJson);
        return { quote: data.quote, quoteHash: data.quoteHash };
    } catch (e) {
        console.error("❌ ERROR: Failed to parse quote.json:", e.message);
        process.exit(1);
    }
}

const { quote: QUOTE, quoteHash: QUOTE_HASH } = loadQuote();

// ============================================================================
// LBC Contract ABI
// ============================================================================

const LBC_ABI = [{
    "inputs": [
        {
            "components": [
                { "internalType": "address", "name": "lbcAddress", "type": "address" },
                { "internalType": "address", "name": "lpRskAddress", "type": "address" },
                { "internalType": "bytes", "name": "btcRefundAddress", "type": "bytes" },
                { "internalType": "address", "name": "rskRefundAddress", "type": "address" },
                { "internalType": "bytes", "name": "lpBtcAddress", "type": "bytes" },
                { "internalType": "uint256", "name": "callFee", "type": "uint256" },
                { "internalType": "uint256", "name": "penaltyFee", "type": "uint256" },
                { "internalType": "int64", "name": "nonce", "type": "int64" },
                { "internalType": "bytes", "name": "deposityAddress", "type": "bytes" },
                { "internalType": "uint256", "name": "value", "type": "uint256" },
                { "internalType": "uint32", "name": "agreementTimestamp", "type": "uint32" },
                { "internalType": "uint32", "name": "depositDateLimit", "type": "uint32" },
                { "internalType": "uint16", "name": "depositConfirmations", "type": "uint16" },
                { "internalType": "uint16", "name": "transferConfirmations", "type": "uint16" },
                { "internalType": "uint32", "name": "transferTime", "type": "uint32" },
                { "internalType": "uint32", "name": "expireDate", "type": "uint32" },
                { "internalType": "uint32", "name": "expireBlock", "type": "uint32" },
                { "internalType": "uint256", "name": "productFeeAmount", "type": "uint256" },
                { "internalType": "uint256", "name": "gasFee", "type": "uint256" }
            ],
            "internalType": "struct Quotes.PegOutQuote",
            "name": "quote",
            "type": "tuple"
        },
        { "internalType": "bytes", "name": "signature", "type": "bytes" }
    ],
    "name": "depositPegout",
    "outputs": [],
    "stateMutability": "payable",
    "type": "function"
}];

// ============================================================================
// HELPERS
// ============================================================================

function decodeBtcAddress(address) {
    return bs58check.decode(address);
}

function buildQuoteStruct(q) {
    return {
        lbcAddress: ethers.getAddress(q.lbcAddress.toLowerCase()),
        lpRskAddress: ethers.getAddress(q.liquidityProviderRskAddress.toLowerCase()),
        btcRefundAddress: decodeBtcAddress(q.btcRefundAddress),
        rskRefundAddress: ethers.getAddress(q.rskRefundAddress.toLowerCase()),
        lpBtcAddress: decodeBtcAddress(q.lpBtcAddr),
        callFee: BigInt(q.callFee),
        penaltyFee: BigInt(q.penaltyFee),
        nonce: BigInt(q.nonce),
        deposityAddress: decodeBtcAddress(q.depositAddr),
        value: BigInt(q.value),
        agreementTimestamp: Number(q.agreementTimestamp),
        depositDateLimit: Number(q.depositDateLimit),
        depositConfirmations: Number(q.depositConfirmations),
        transferConfirmations: Number(q.transferConfirmations),
        transferTime: Number(q.transferTime),
        expireDate: Number(q.expireDate),
        expireBlock: Number(q.expireBlocks),
        productFeeAmount: BigInt(q.productFeeAmount),
        gasFee: BigInt(q.gasFee)
    };
}

function calculateTotal(q) {
    return BigInt(q.value) + BigInt(q.callFee) + BigInt(q.gasFee) + BigInt(q.productFeeAmount);
}

// ============================================================================
// MAIN
// ============================================================================

async function main() {
    console.log("=== Pegout Payment ===\n");

    // Validate
    if (!QUOTE?.lbcAddress) {
        console.error("❌ Invalid quote.json");
        process.exit(1);
    }

    if (SIGNATURE === "PASTE_YOUR_SIGNATURE_HERE") {
        console.error("❌ Missing signature!");
        console.error("   1. Call acceptQuote with quoteHash:", QUOTE_HASH);
        console.error("   2. Edit index.js and paste the signature");
        process.exit(1);
    }

    console.log("Quote Hash:", QUOTE_HASH);

    const provider = new ethers.JsonRpcProvider(RSK_RPC_URL);
    const network = await provider.getNetwork();
    console.log("Chain ID:", network.chainId.toString());

    const quoteStruct = buildQuoteStruct(QUOTE);
    const totalValue = calculateTotal(QUOTE);

    console.log("\nTotal:", ethers.formatEther(totalValue), "RBTC");
    console.log("  Value:", ethers.formatEther(BigInt(QUOTE.value)));
    console.log("  Call Fee:", ethers.formatEther(BigInt(QUOTE.callFee)));
    console.log("  Gas Fee:", ethers.formatEther(BigInt(QUOTE.gasFee)));

    const iface = new ethers.Interface(LBC_ABI);
    const data = iface.encodeFunctionData("depositPegout", [quoteStruct, "0x" + SIGNATURE]);

    console.log("\nSending to LBC:", QUOTE.lbcAddress);

    try {
        const txHash = await provider.send("eth_sendTransaction", [{
            from: PAYMENT_ADDRESS,
            to: QUOTE.lbcAddress,
            value: "0x" + totalValue.toString(16),
            data: data,
            gas: "0x" + (500000).toString(16)
        }]);
        console.log("\n✅ Success! TX:", txHash);
    } catch (error) {
        console.error("\n❌ Failed:", error.message);
    }
}

main().catch(console.error);
