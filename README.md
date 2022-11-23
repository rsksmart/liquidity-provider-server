# Liquidity Provider Server

This is a server that interacts with a [Liquidity Bridge Contract (LBC)](https://github.com/rsksmart/liquidity-bridge-contract) to provide liquidity for users 
as part of the Flyover protocol. The server runs a local [Liquidity Provider (LP)](https://github.com/rsksmart/liquidity-provider), and also allows connections
from remote LPs.

The server's functionality is provided through a JSON HTTP interface. In addition, the server needs access to a Bitcoin and an RSK node.

## Configuration

#### Configuration File

    - logfile (string): the path where the logs are saved to. If empty, it prints logs to the console.
    - debug (bool): the value that indicates whether the server is run in debug mode.
    - irisActivationHeight: the block height at where Iris was activated, so the federation goes into ERP.
    - erpKeys (array[string]): the public keys of the erp pegnatories to be used in p2sh scripts.
    - server (object): object that holds settings for the http server.
        - port (int): port where the api is served.
    - db (object): object that holds settings for the database.
        - path (string): path to the sqlite db file.
    - rsk (object): object that holds settings for the rsk connector.
        - endpoint (string): endpoint to the json-rpc api where the RSK node is listening.
        - lbcAddr (string): address of the Liquidity Bridge Contract.
        - bridgeAddr (string): address of the Bridge Contract.
        - requiredBridgeConfirmations (int): amount of confirmations required by the Bridge Contract.
    - btc (object): object that holds settings for the bitcoin connector.
        - endpoint (string): Url where the Bitcoin node is hosted (in the format IP:PORT).
        - username (string): username to be used in the connection to the bitcoin node.
        - password (string): password to be used in the connection to the bitcoin node.
        - network (string): network to be used in the connection to the bitcoin node.
    - provider (object): object that holds settings for the local liquidity provider.
        - keydir (string): directory where the keystore is located (by default "keystore").
        - pwdFile (string): The path to the file that contains the password that matches the keystore specified above. 
        - btcAddr (string): bitcoin address of the local Liquidity Provider registered with the server.
        - accountNum (int): RSK Account number for the local Liquidity Provider.
        - chainId (int): id of the RSK network to be used.
        - maxConf: Maximum amount of confirmations required for a callForUser (liquidity advancement).
        - confirmations (object): object that holds pairs of value-confirmations so the Liquidity provider
                can specify different amount of required confirmations depending on the value to transfer.
            (string): (int): The objects within this container specify value-confirmations in the following
                format: 
                    {
                        "50000000000": 5,
                        "1000000000": 3 
                        ...
                    }
    - timeForDeposit (int): the default time threshold for deposit to be set in quotes, in seconds.
    - callTime (int): the default time the Liquidity Provider has to advance the funds.
    - callFee (int): the default fee to be applied to a quote.
    - penaltyFee (int): the penalty fee to be applied in case of missbehaviour.



## API

### getQuote

Computes and returns a quote for the service.

#### Parameters

    contractAddr (string) - Hex-encoded contract address.
    data (string) - Hex-encoded contract data.
    value (int) - Value to send in the call.
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
        gasLimit;                         // the gas limit -> Calculated based on the estimated gas in the network
        nonce;                            // a nonce that uniquely identifies this quote
        value;                            // the value to transfer in the call
        agreementTimestamp;               // the timestamp of the agreement
        timeForDeposit;                   // the time (in seconds) that the user has to achieve one confirmation on the BTC deposit
        callTime;                         // the time (in seconds) that the LP has to perform the call on behalf of the user after the deposit achieves the number of confirmations
        confirmations;                    // the number of confirmations that the LP requires before making the call
        callOnRegister:                   // a boolean value indicating whether the callForUser can be called on registerPegIn.
    
### acceptQuote

Accepts one of the LPs quotes.

#### Parameters

    quoteHash (string) - Hex-encoded quote hash as computed by LBC.hashQuote

#### Returns

    signature - Signature of the quote
    bitcoinDepositAddressHash - Hash of the deposit BTC address
