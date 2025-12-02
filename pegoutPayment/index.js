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
 *   1. Paste your quote data from getQuote response in the QUOTE section below
 *   2. Paste your signature from acceptQuote response in the SIGNATURE section below
 *   3. Run: npm run pay
 */

const { ethers } = require('ethers');
const bs58check = require('bs58check');

// ============================================================================
// PASTE YOUR DATA HERE
// ============================================================================

// Paste the quote object from the getQuote response (the "quote" field inside the array)
// Example: If your response was [{ "quote": {...}, "quoteHash": "..." }]
// Copy the {...} part and paste it below:

const QUOTE = {
    "lbcAddress": "0x03f23ae1917722d5a27a2ea0bcc98725a2a2a49a",
    "liquidityProviderRskAddress": "0x9d93929a9099be4355fc2389fbf253982f9df47c",
    "btcRefundAddress": "n2NxBRs6aDvuL3qZ3e6vDxWst6wGCyLpHK",
    "rskRefundAddress": "0x45400C53eBd0853Cd26b21C3d479f0eedc46bc44",
    "lpBtcAddr": "n1jGDaxCW6jemLZyd9wmDHddseZwEMV9C6",
    "callFee": "2180000000000000",
    "penaltyFee": "1000000000000000",
    "nonce": "3504134078398023607",  // IMPORTANT: Must be string to avoid precision loss
    "depositAddr": "n2NxBRs6aDvuL3qZ3e6vDxWst6wGCyLpHK",
    "value": "600000000000000000",
    "agreementTimestamp": 1764715974,
    "depositDateLimit": 1764719574,
    "depositConfirmations": 20,
    "transferConfirmations": 10,
    "transferTime": 3600,
    "expireDate": 1764726774,
    "expireBlocks": 1962,
    "gasFee": "67250000000000",
    "productFeeAmount": "0"
};

// Paste the signature from the acceptQuote response (without 0x prefix)
const SIGNATURE = "30a3913896dac33d47ec29d45fc6592f0da7a3b665a5c758d8bcdddff48b4c4016cbb01309a2861aefc90ec192678897d0eec00fa1462068b4503e52b72a04071c";

// ============================================================================
// CONFIGURATION (usually no need to change for local environment)
// ============================================================================

const RSK_RPC_URL = "http://127.0.0.1:4444";
// This is the pre-funded account in the local RSKj regtest
const PAYMENT_ADDRESS = "0xcd2a3d9f938e13cd947ec05abc7fe734df8dd826";

// ============================================================================
// LBC Contract ABI (only the depositPegout function)
// ============================================================================

const LBC_ABI = [
    {
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
    }
];

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

function stringToBytes(str) {
    return ethers.toUtf8Bytes(str);
}

// Decode Bitcoin Base58Check address to bytes (with version byte)
function decodeBtcAddress(address) {
    const decoded = bs58check.decode(address);
    return decoded;
}

function toChecksumAddress(address) {
    // Convert to lowercase first, then get proper checksum
    return ethers.getAddress(address.toLowerCase());
}

function buildQuoteStruct(quote) {
    return {
        lbcAddress: toChecksumAddress(quote.lbcAddress),
        lpRskAddress: toChecksumAddress(quote.liquidityProviderRskAddress),
        btcRefundAddress: decodeBtcAddress(quote.btcRefundAddress),
        rskRefundAddress: toChecksumAddress(quote.rskRefundAddress),
        lpBtcAddress: decodeBtcAddress(quote.lpBtcAddr),
        callFee: BigInt(quote.callFee),
        penaltyFee: BigInt(quote.penaltyFee),
        nonce: BigInt(quote.nonce),
        deposityAddress: decodeBtcAddress(quote.depositAddr),
        value: BigInt(quote.value),
        agreementTimestamp: quote.agreementTimestamp,
        depositDateLimit: quote.depositDateLimit,
        depositConfirmations: quote.depositConfirmations,
        transferConfirmations: quote.transferConfirmations,
        transferTime: quote.transferTime,
        expireDate: quote.expireDate,
        expireBlock: quote.expireBlocks,
        productFeeAmount: BigInt(quote.productFeeAmount),
        gasFee: BigInt(quote.gasFee)
    };
}

function calculateTotalValue(quote) {
    const value = BigInt(quote.value);
    const callFee = BigInt(quote.callFee);
    const gasFee = BigInt(quote.gasFee);
    const productFee = BigInt(quote.productFeeAmount);
    return value + callFee + gasFee + productFee;
}

// ============================================================================
// MAIN
// ============================================================================

async function main() {
    console.log("===================================");
    console.log("Pegout Payment Script");
    console.log("===================================\n");

    // Validate input
    if (QUOTE.lbcAddress === "PASTE_LBC_ADDRESS_HERE" || SIGNATURE === "PASTE_SIGNATURE_HERE") {
        console.error("❌ ERROR: Please paste your quote and signature data in the script before running!");
        console.error("   Open index.js and replace the placeholder values in the QUOTE and SIGNATURE sections.");
        process.exit(1);
    }

    // Connect to RSKj
    const provider = new ethers.JsonRpcProvider(RSK_RPC_URL);
    
    console.log("Connected to RSKj at:", RSK_RPC_URL);
    
    // Get network info
    const network = await provider.getNetwork();
    console.log("Chain ID:", network.chainId.toString());
    
    // Build the quote struct for the contract
    const quoteStruct = buildQuoteStruct(QUOTE);
    console.log("\nQuote struct built successfully");
    
    // Calculate total value to send
    const totalValue = calculateTotalValue(QUOTE);
    console.log("\nTotal value to send:", ethers.formatEther(totalValue), "RBTC");
    console.log("  - Value:", ethers.formatEther(BigInt(QUOTE.value)), "RBTC");
    console.log("  - Call Fee:", ethers.formatEther(BigInt(QUOTE.callFee)), "RBTC");
    console.log("  - Gas Fee:", ethers.formatEther(BigInt(QUOTE.gasFee)), "RBTC");
    console.log("  - Product Fee:", ethers.formatEther(BigInt(QUOTE.productFeeAmount)), "RBTC");
    
    // Prepare signature
    const signatureBytes = "0x" + SIGNATURE;
    console.log("\nSignature prepared");
    
    // Create contract interface for encoding
    const iface = new ethers.Interface(LBC_ABI);
    const data = iface.encodeFunctionData("depositPegout", [quoteStruct, signatureBytes]);
    
    console.log("\nTransaction data encoded");
    console.log("LBC Address:", QUOTE.lbcAddress);
    
    // Send raw transaction using eth_sendTransaction (for regtest with unlocked accounts)
    console.log("\nSending transaction...");
    
    const txParams = {
        from: PAYMENT_ADDRESS,
        to: QUOTE.lbcAddress,
        value: "0x" + totalValue.toString(16),
        data: data,
        gas: "0x" + (500000).toString(16)
    };
    
    try {
        const txHash = await provider.send("eth_sendTransaction", [txParams]);
        console.log("\n✅ Transaction sent successfully!");
        console.log("Transaction Hash:", txHash);
        console.log("\nMonitor the pegout state in MongoDB Compass (flyover/retainedPegoutQuote collection)");
    } catch (error) {
        console.error("\n❌ Transaction failed!");
        console.error("Error:", error.message);
        if (error.data) {
            console.error("Error data:", error.data);
        }
    }
}

main().catch(console.error);

