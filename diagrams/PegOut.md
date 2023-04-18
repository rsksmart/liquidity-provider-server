# PegOut Process
This process consist in the conversion of RBTC to BTC. Meaning that tokens will be transferred from RSK network to Bitcoin network. To achieve that, Liquidity Provider Server makes following steps:

* Precondition: the user has executed /pegout/getQuotes endpoint and selected one of the returned quotes
1. The user executes accept quote with the selected quote's hash and LPS responds with an RSK address to deposit 
2. LPS creates a RSK watcher for that particular address to be able to monitor when the deposit is made and whe the required confirmation blocks have been mined
3. Once the deposit to the RSK address is done and the required confirmations have passed, the LPS checks if the LP has balance available and locks it, then registers the pegout on the LBC and unlocks the balance to send to its destination
4. After sending balance to destination, LPS creates an BTC watcher to check when the required confirmations of the deposit made in previous step have passed
5. When deposit confirmations have passed, LPS calls refundPegout method of LBC to mark quote as finished and verify and punish LP if necessary 
6. After calling refundPegout and do the proper verifications, LPS sends RBTC to the bridge to convert it to BTC, refunding LP and giving him his fee 