# PegIn Process
This process consist in the conversion of BTC to RBTC. Meaning that tokens will be transferred from Bitcoin network to RSK network. To achieve that, Liquidity Provider Server makes following steps:

* Precondition: the user has executed /pegin/getQuote endpoint and selected one of the returned quotes
1. The user executes accept quote with the selected quote's hash and LPS responds with an Bitcoin address to deposit
2. LPS creates a BTC watcher for that particular address to be able to monitor when the deposit is made and whe the required confirmation blocks have been mined
3. The user makes the deposit to the Bitcoin derivation address
4. Once the deposit to the Bitcoin address is done and the required confirmations have passed, the LPS executes the callForUser function of LBC to send the RBTC to its destination.
5. The LPS waits until required bridge confirmations have passed
6. LPS executes the registerPegin function of LBC, refunding LP and paying his fee