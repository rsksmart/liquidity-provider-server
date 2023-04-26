# Liquidity Provider Server

This is a server that interacts with a [Liquidity Bridge Contract (LBC)](https://github.com/rsksmart/liquidity-bridge-contract) to provide liquidity for users 
as part of the Flyover protocol. The server runs a local [Liquidity Provider (LP)](https://github.com/rsksmart/liquidity-provider).

The server's functionality is provided through a JSON HTTP interface. In addition, the server needs access to a Bitcoin and an RSK node.

## Configuration

#### Configuration File

    - logfile (string): the path where the logs are saved to. If empty, it prints logs to the console.
    - debug (bool): the value that indicates whether the server is run in debug mode.
    - irisActivationHeight: the block height at where Iris was activated, so the federation goes into ERP.
    - maxQuoteValue: The maximum allowed value to a quote globaly.
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
        - maxQuoteValue: The maximum allowed value to a quote to the liquidity provider, will be overrided by global value if that value where greater than 0.
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

### Pegout
You can see the detail of the pegout process [here](./diagrams/PegOut.md)

### /pegout/getQuotes

Computes and returns a quote for the pegout service.

#### Parameters
    
    to (string) - BTC destination address
    valueToTransfer (int) - value to transfer to BTC address
    rskRefundAddress (string) - Hex-encoded user RSK refund address.
    bitcoinRefundAddress (string) - Base58-encoded user Bitcoin refund address.

#### Returns

    quotes - a list of pegout quotes for the service, where each quote consists of:
        quote - a pegout quote
            lbcAddress;                                 // the address of the LBC
            liquidityProviderRskAddress;                // the RSK address of the LP
            btcRefundAddress;                           // a user BTC refund address
            rskRefundAddress;                           // a user RSK refund address 
            lpBtcAddr;                                  // the BTC address of the LP
            callFee;                                    // the fee charged by the LP
            penaltyFee;                                 // the penalty that the LP pays if it fails to deliver the service
            nonce;                                      // a nonce that uniquely identifies this quote
            depositAddr;                                // the destination address of the peg-out
            gasLimit;                                   // the gas limit -> Calculated based on the estimated gas in the network
            value;                                      // the value to transfer in the call
            agreementTimestamp;                         // the timestamp of the agreement
            depositDateLimit;                           // time in seconds to do the deposit
            depositConfirmations;                       // number of confirmations to do the pegout
            transferConfirmations;                      // number of pegout confirmations to do the refund
            transferTime;                               // time in seconds to do the pegout without punishing the LP
            expireDate;                                 // the timestamp of the expiration
            expireBlocks;                               // amount of blocks for the quote to be expired

        quoteHash - the corresponding quote hash

### /pegout/acceptQuote

Accepts one of the LPs pegout quotes.

#### Parameters

    quoteHash (string) - Hex-encoded quote hash as computed by LBC.hashQuote

#### Returns

    signature - Signature of the quote
    lbcAddress - Address of the contract to execute the depositPegout function

### Pegin
You can see the detail of the pegin process [here](./diagrams/PegIn.md)

### /pegin/getQuote

Computes and returns a quote for the pegin service.

#### Parameters

    callEoaOrContractAddress (string) - Hex-encoded contract address.
    data (string) - Hex-encoded contract data.
    valueToTransfer (int) - Value to send in the call.
    rskRefundAddress (string) - Hex-encoded user RSK refund address.
    bitcoinRefundAddress (string) - Base58-encoded user Bitcoin refund address.

#### Returns

    quotes - a list of pegin quotes for the service, where each quote consists of:
        quote - a pegin quote
            fedBtcAddress;                    // the BTC address
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
        quoteHash - the corresponding quote hash
    
### /pegin/acceptQuote

Accepts one of the LPs pegin quotes.

#### Parameters

    quoteHash (string) - Hex-encoded quote hash as computed by LBC.hashQuote

#### Returns

    signature - Signature of the quote
    bitcoinDepositAddressHash - Hash of the deposit BTC address
    
    
### getProviders

Gets the registered providers List.

#### Parameters


#### Returns

    Array of registered providers with the fields.

    Id - The Id of the provider
    Provider - The address of the Liquidity Provider in RSK Network


## Run Integration Tests

1. You should run LPS env into [docker-compose](./docker-compose/README.md)
2. Then you need to change directory to it
3. Then into it folder you only need to run `go test -integration`


#### Note: 
It is required to run LPS env to run integration tests because in `it` folder there is a config.json file which will be updated when contracts will be deployed.