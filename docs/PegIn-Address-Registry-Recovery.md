# PegIn address registry recovery

Use this runbook when the LPS cannot repair its peg-in watch list by itself. Also use it after a wrong registry address, start block or network configuration.

## How the LPS keeps the watch list

The LPS stores two things in MongoDB:

- `peginWatch`: one document per registered RSK address, with its BTC deposit address and import state.
- `peginWatchCheckpoint`: the checkpoint, one document (`_id: "verified"`) with the last verified registry root and the block it was verified at.

Every minute the LPS reads the `AddressRegistered` events after the checkpoint block, page by page (`PEGIN_ADDRESS_REGISTRY_WATCHER_PAGE_SIZE`), up to the finalized head: the current head minus `RSK_MAX_REORG_DEPTH` (default `2`). It folds each address into the checkpoint root and checks the result against:

- the root in each event;
- `getRegistrationRoot()` at the last block of each page.

When every check matches, it saves the new documents and the new checkpoint in one MongoDB transaction. It keeps reading pages until it reaches the finalized head.

## Automatic repair

| Situation | What the LPS does |
| --- | --- |
| No checkpoint: a first start, or an upgrade from a version without one | Rebuilds the list from `PEGIN_ADDRESS_REGISTRY_WATCHER_START_BLOCK`. |
| The new blocks do not match, for example the node dropped an event | Reads the registry root at the checkpoint block. If it still matches, reads the new blocks again from the checkpoint, in the same tick. |
| The chain no longer matches at or below the checkpoint block, for example after a reorg deeper than `RSK_MAX_REORG_DEPTH` | Finds the first block where the stored documents and the registry differ, with a binary search of about log2(checkpoint block − start block) reads, then rebuilds from that block. |
| The checkpoint is below the start block | Rebuilds from the start block. |
| The registry root at the checkpoint block differs from the checkpoint, and the stored documents do not match the checkpoint either | Rebuilds from the start block, without a binary search. |
| The checkpoint is above the finalized head | Changes nothing and waits for the node to catch up. |

A rebuild reads every page from the rebuild block to the finalized head and checks the root in each event. The rebuild read itself calls `getRegistrationRoot()` only once, at the finalized head. Then it writes everything in one MongoDB transaction:

- documents at or above the rebuild block that the chain no longer has are removed;
- documents that are still on the chain keep their BTC address and import state;
- new addresses are added as `discovered`;
- the new checkpoint is saved.

If any of these steps fails, nothing changes and the LPS tries again on the next tick.

On the same tick, after the transaction, the LPS imports the `discovered` addresses into the Bitcoin wallet and rescans. This is not part of the transaction: if it fails, the rebuild stays saved, the document keeps the error in `last_error` and stays `discovered`, and the LPS tries the import again on the next tick.

### What to look for in the logs

| Log line | Meaning |
| --- | --- |
| `PegIn address registry catch-up failed, rebuilding the watch list: ...` | A mismatch was found. A rebuild starts. |
| `Rebuilt the PegIn watch list from block X to block Y with N registrations` | The rebuild succeeded. No action needed. |
| `Stored PegIn watch documents do not match the checkpoint at block X, rebuilding from the start block Y` | The documents and the checkpoint disagree. A full rebuild starts. |
| `Stored PegIn address registry checkpoint at block X is above the finalized RSK head Y, waiting for the node to catch up` | The node is behind or on another network. |
| `PegIn address registry watcher scan failed: ... cannot read the PegIn address registry root at old block X, check that the RSK node can answer calls at old blocks: ...` | A repair after a mismatch needs the registry root at an old block, and the node cannot answer. See manual repair. |
| `PegIn address registry watcher scan failed: ... get PegIn address registry root at block X: ...` | The node could not answer a root read at the end of a page or at the finalized head. After a long stop that block can be old. See manual repair if it repeats. |
| `PegIn address registry watcher scan failed: ...` | The tick failed. The rest of the line says why. |

### Requirements and limits

- Every repair after a mismatch calls `getRegistrationRoot()` at the checkpoint block, and the binary search calls it at older blocks, so the RSK node must keep the state of old blocks. A full rebuild from the start block reads the root only at the finalized head. A catch-up reads it at the end of each page, so after a long stop it also reads old blocks.
- The registry root covers only the ordered list of RSK addresses. It does not cover the BTC address, encoding, state, errors or timestamps of a document.
- A rebuild rewrites only the documents at or above the rebuild block. If a reorg moved a registration below that block to another block, its document keeps the old tx hash and block number. The address list is still right.
- Documents below the start block are never changed.
- Bitcoin Core cannot remove an address imported with `importaddress`. An address removed from the list stays in the wallet. It only adds rescan cost.
- The rescan after an import starts 100 BTC blocks (about 16 hours) below the BTC tip. The wallet does not see older deposits to a newly imported address. This matters after step 6 of the manual repair, and for addresses registered while the LPS was stopped for longer than that.

## Manual repair

Use this when the same error repeats on every tick, the waiting warning stays while the node is in sync, or the configuration was wrong.

1. Stop the LPS. Do not change watcher data while it runs.
2. Back up MongoDB.
3. Check the configuration:
   - `LPS_STAGE`, `CHAIN_ID` and `RSK_ENDPOINT` point to the intended Rootstock network.
   - `PEGIN_ADDRESS_REGISTRY_ADDRESS` is the registry of that network.
   - `PEGIN_ADDRESS_REGISTRY_WATCHER_START_BLOCK` is at or below the block where the registry was deployed. The LPS does not check this.
4. If the error says the node cannot answer calls at old blocks, or `get PegIn address registry root at block X` repeats after a long stop, point `RSK_ENDPOINT` to a node that keeps old state, or go on to step 5, which reads the root only at the finalized head.
5. Delete the checkpoint:

   ```javascript
   db.peginWatchCheckpoint.deleteMany({})
   ```

   On the next start the LPS rebuilds the list from the start block. Documents still on the chain keep their import state, and documents the chain does not have are removed.
6. Also delete the documents if the registry address or the network was wrong, or if documents were edited or damaged by hand. A rebuild keeps the BTC address and import state of the documents it keeps, and never changes documents below the start block:

   ```javascript
   db.peginWatch.deleteMany({})
   ```

   Every address is then imported into the Bitcoin wallet again. The rescan covers only the last 100 BTC blocks.
7. Start the LPS and wait for `Rebuilt the PegIn watch list ...` in the logs.
8. If you deleted the documents in step 6, or the LPS was stopped for more than 100 BTC blocks, and an older deposit may still matter, rescan the wallet from a BTC height before that deposit. The rescan can take a long time and blocks the wallet while it runs:

   ```bash
   bitcoin-cli -rpcwallet=pegin-watchonly-wallet rescanblockchain <height>
   ```

9. Compare the checkpoint with the chain. `db.peginWatchCheckpoint.findOne()` gives `root` and `block`, and `db.peginWatchCheckpoint.findOne().root.toString('hex')` prints the root in hex. `getRegistrationRoot()` called at that block must return the same root.
10. If the roots still differ, stop the LPS, record the block and both roots from the error, check the node and the configuration again, then escalate.
