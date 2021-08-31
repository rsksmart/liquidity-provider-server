# Liquidity Provider Server

This is a server that interacts with a [Liquidity Bridge Contract (LBC)](https://github.com/rsksmart/liquidity-bridge-contract) to provide liquidity for users 
as part of the Flyover protocol. The server runs a local [Liquidity Provider (LP)](https://github.com/rsksmart/liquidity-provider), and also allows connections
from remote LPs.

The server's functionality is provided through a JSON HTTP interface. In addition, the server needs access to a Bitcoin and an RSK node.

## API

### getQuote

Computes and returns a quote for the service.

#### Parameters

    contractAddr (string) - Hex-encoded contract address.
    data (string) - Hex-encoded contract data.
    value (int) - Value to send in the call.
    gasLimit (int) - Gas limit to use in the call.
    rskRefundAddr (string) - Hex-encoded user RSK refund address.
    btcRefundAddr (string) - Base58-encoded user Bitcoin refund address.

#### Returns

    quotes - a list of quotes for the service, where each quote consists of:

        fedBtcAddress;                    // the BTC address of the Powpeg
        lbcAddress;                       // the address of the LBC
        lpRSKAddr;                        // the RSK address of the LP
        btcRefundAddress;                 // a user BTC refund address
        rskRefundAddress;                 // a user RSK refund address 
        lpBTCAddr;                        // the BTC address of the LP
        callFee;                          // the fee charged by the LP
        penaltyFee;                       // the penalty that the LP pays if it fails to deliver the service
        contractAddr;                     // the destination address of the peg-in
        data;                             // the arguments to send in the call
        gasLimit;                         // the gas limit
        nonce;                            // a nonce that uniquely identifies this quote
        value;                            // the value to transfer in the call
        agreementTimestamp;               // the timestamp of the agreement
        timeForDeposit;                   // the time (in seconds) that the user has to achieve one confirmation on the BTC deposit
        callTime;                         // the time (in seconds) that the LP has to perform the call on behalf of the user after the deposit achieves the number of confirmations
        confirmations;                   // the number of confirmations that the LP requires before making the call
    
### acceptQuote

Accepts one of the LPs quotes.

#### Parameters

    quoteHash (string) - Hex-encoded quote hash as computed by LBC.hashQuote

#### Returns

    signature - Signature of the quote
    bitcoinDepositAddressHash - Hash of the deposit BTC address
