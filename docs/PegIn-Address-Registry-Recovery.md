# PegIn address registry recovery

Use this runbook when the PegIn address registry watcher reports a root mismatch that automatic recovery cannot correct. Also use it after an incorrect registry address, deployment block, or network configuration.

## Automatic recovery

The watcher folds the current MongoDB `rsk_address` values in chain order. It compares the result with the registry root at the same Rootstock block.

If the roots differ, the watcher finds a block boundary where the stored prefix is valid. It fetches and validates the canonical `AddressRegistered` event suffix in memory. It deletes the invalid suffix from the `peginWatch` collection only after that validation succeeds. It then stores the replayed events, reads the stored rows again, and verifies the new root. The watcher keeps its verified checkpoint in memory only.

A Rootstock reorg starts the same recovery operation early. It does not use a separate recovery path.

An RPC error stops the operation before suffix deletion. An event-fetch error leaves the stored suffix unchanged. The watcher does not delete data based on an RPC timeout or an unavailable historical state.

### Scan sequence

Each tick, and once at startup through `Prepare`, the watcher runs this sequence:

1. Replay verifies registry integrity and persists `peginWatch` rows. It does not import Bitcoin addresses.
2. The watcher lists rows in state `discovered`. These rows are the pending Bitcoin imports.
3. The watcher imports each pending row, one call per row.
4. The watcher runs one Bitcoin rescan for the rows that need it.
5. The watcher finalizes the import results.

`Prepare` runs Replay with no checkpoint. Each later tick passes the checkpoint from the last successful Replay.

### Recovery reason values

Replay returns a recovery reason. The watcher publishes the events. Replay does not publish.

| Reason | Replay returns it when | EventBus `Reason` |
| --- | --- | --- |
| empty | The head is below the configured start block, or the local root already equals the chain root. | none published |
| `catch_up` | A trusted checkpoint let Replay advance through a suffix. | `catch_up` |
| `root_mismatch` | Binary search found the first bad block and Replay rebuilt the suffix. | `root_mismatch` |

The watcher publishes a `PegInAddressRegistryRootMismatch` event **only after the rebuild succeeds**. It carries the roots captured before the rebuild. A failed plan, fetch, or persist-verify writes an error log and publishes nothing. An operator therefore never sees a mismatch event for a state that LPS did not repair.

### Import failure

If a Bitcoin import fails, the tick stops at that row. The checkpoint from the successful Replay stays in memory. The remaining `discovered` rows retry on the next tick. A wallet import failure does not roll back a verified registry state.

### In-memory checkpoint

The checkpoint holds the verified local root and the block it was verified through. It lives in watcher memory. There is no checkpoint document in MongoDB, so there is nothing to query or edit.

The `peginWatch` rows stay the durable source of truth. After a restart the watcher starts with no checkpoint. It pays one binary-search root walk, about `log2(head - start)` historical `getRegistrationRoot` reads, and then replays the suffix as usual.

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
   - `PEGIN_ADDRESS_REGISTRY_WATCHER_START_BLOCK` is at or below the first registry event and is not above the current head. LPS uses this value as the lower scan bound. It does not verify that the value is the deployment block.
   - `LBC_ADDR` and the registry configuration refer to the same active-network LBC.
4. Delete the registry watch collection contents:

   ```javascript
   db.peginWatch.deleteMany({})
   ```

   This command deletes the registry watch rows. It does not delete quotes or deposits.
5. Restart LPS with the verified configuration. The watcher finds the first mismatching block at or after the deployment block and replays registry events from that boundary.
6. Recompute the root from the ordered `peginWatch` rows. Compare it with the result of `getRegistrationRoot()` at the current Rootstock head.
7. If the roots differ, stop LPS again. Record the block and both roots. Verify the Rootstock endpoint and registry configuration, then investigate or escalate.

Bitcoin Core cannot remove addresses that LPS imported with `importaddress`. Recovery can therefore leave old imported addresses in the Bitcoin wallet. These addresses add rescan cost, but they do not restore deleted MongoDB watch rows or change the registry root check.
