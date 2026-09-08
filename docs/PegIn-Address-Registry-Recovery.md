# PegIn address registry recovery

Use this runbook when the PegIn address registry watcher reports a root mismatch that automatic recovery cannot correct. Also use it after an incorrect registry address, deployment block, or network configuration.

## Automatic recovery

The watcher folds the current MongoDB `rsk_address` values in chain order. It compares the result with the registry root at the same Rootstock block.

If the roots differ, the watcher finds a block boundary where the stored prefix is valid. It deletes the invalid suffix from the `peginWatch` collection and replays canonical `AddressRegistered` events from that boundary. It then reads the stored rows again and verifies the new root before it writes a checkpoint.

A Rootstock reorg starts the same recovery operation early. It does not use a separate recovery path.

An RPC error stops the operation before suffix deletion. The watcher does not delete data based on an RPC timeout or an unavailable historical state.

Automatic recovery supports changes to the ordered registry address sequence:

- a missing `RskAddress`;
- an extra `RskAddress`;
- a changed, syntactically valid `RskAddress`;
- reordered addresses.

The registry root does not cover these fields:

- `BtcAddress`;
- `Encoding`;
- `State`;
- `LastSeenAt`;
- `LastError`;
- timestamps.

Automatic recovery does not correct these fields. It also cannot detect a block-number-only change that preserves the address order because the registry root does not include block numbers. A malformed `RskAddress` stops recovery because the watcher cannot fold it.

## Manual full replay

Use this procedure only after automatic recovery fails or an operator confirms that the Rootstock node cannot provide the required historical state.

1. Stop the LPS process. Do not modify watcher data while LPS runs.
2. Back up MongoDB.
3. Verify the active network and configuration:
   - `LPS_STAGE`, `CHAIN_ID`, and `RSK_ENDPOINT` select the intended Rootstock network.
   - `PEGIN_ADDRESS_REGISTRY_ADDRESS` is the registry for that network.
   - `PEGIN_ADDRESS_REGISTRY_WATCHER_START_BLOCK` is the registry deployment block and is not above the current head.
   - `LBC_ADDR` and the registry configuration refer to the same active-network LBC.
4. Delete the registry watch collection contents:

   ```javascript
   db.peginWatch.deleteMany({})
   ```

   This command deletes the registry watch rows and their checkpoint. It does not delete quotes or deposits.
5. Restart LPS with the verified configuration. The watcher finds the first mismatching block at or after the deployment block and replays registry events from that boundary.
6. Read the checkpoint with `db.peginWatch.findOne({_id: "checkpoint"})`. Compare its `local_root` with the result of `getRegistrationRoot()` at its `last_processed_block`.
7. If the roots differ, stop LPS again. Record the block and both roots. Verify the Rootstock endpoint and registry configuration, then investigate or escalate.

Bitcoin Core cannot remove addresses that LPS imported with `importaddress`. Recovery can therefore leave old imported addresses in the Bitcoin wallet. These addresses add rescan cost, but they do not restore deleted MongoDB watch rows or change the registry root check.
