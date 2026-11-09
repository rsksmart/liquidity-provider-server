# PegIn address registry recovery

Use this runbook when the LPS cannot repair its peg-in watch list by itself. Also use it after a wrong registry address, start block or network configuration.

## How the LPS keeps the watch list

The LPS stores two things in MongoDB:

- `peginWatch`: one row per registered RSK address, with its BTC deposit address and import state.
- `peginWatchCheckpoint`: one document (`_id: "verified"`) with the last verified registry root and the block it was verified at.

Every minute the LPS reads the `AddressRegistered` events after the stored block, page by page (`PEGIN_ADDRESS_REGISTRY_WATCHER_PAGE_SIZE`), up to the finalized head: the current head minus `RSK_MAX_REORG_DEPTH` (default `2`). It folds each address into the stored root and checks the result against:

- the root in each event;
- `getRegistrationRoot()` at the last block of each page.

When every check matches, it saves the new rows and the new point in one MongoDB transaction. It keeps reading pages until it reaches the finalized head.

## Automatic repair

| Situation | What the LPS does |
| --- | --- |
| No stored point: a first start, or an upgrade from a version without one | Rebuilds the list from `PEGIN_ADDRESS_REGISTRY_WATCHER_START_BLOCK`. |
| The new blocks do not match, for example the node dropped an event | Reads the registry root at the stored block. If it still matches, reads the new blocks again from the stored point, in the same tick. |
| The chain no longer matches at or below the stored block, for example after a reorg deeper than `RSK_MAX_REORG_DEPTH` | Finds the first block where the stored rows and the registry differ, with a binary search of about log2(stored block − start block) reads, then rebuilds from that block. |
| The stored point is below the start block | Rebuilds from the start block. |
| After a mismatch, the stored rows do not match the stored point | Rebuilds from the start block, without a binary search. |
| The stored point is above the finalized head | Changes nothing and waits for the node to catch up. |

A rebuild reads every page from the rebuild block to the finalized head, with the same checks as above, then writes everything in one MongoDB transaction:

- rows at or above the rebuild block that the chain no longer has are removed;
- rows that are still on the chain keep their BTC address and import state;
- new addresses are added as `discovered` and imported on the same tick;
- the new point is saved.

If any step fails, nothing changes and the LPS tries again on the next tick.

### What to look for in the logs

| Log line | Meaning |
| --- | --- |
| `PegIn address registry catch-up failed, rebuilding the watch list: ...` | A mismatch was found. A rebuild starts. |
| `Rebuilt the PegIn watch list from block X to block Y with N registrations` | The rebuild succeeded. No action needed. |
| `Stored PegIn watch rows do not match the stored point at block X, rebuilding from the start block Y` | The rows and the point disagree. A full rebuild starts. |
| `Stored PegIn address registry point at block X is above the finalized RSK head Y, waiting for the node to catch up` | The node is behind or on another network. |
| `cannot read the PegIn address registry root at old block X, check that the RSK node can answer calls at old blocks` | The node cannot serve the binary search. See manual repair. |
| `PegIn address registry watcher scan failed: ...` | The tick failed. The rest of the line says why. |

### Requirements and limits

- The binary search calls `getRegistrationRoot()` at old blocks, so the RSK node must keep the state of old blocks. A full rebuild from the start block reads the root only at the finalized head. A catch-up reads it at the end of each page, so after a long stop it also reads old blocks.
- The registry root covers only the ordered list of RSK addresses. It does not cover the BTC address, encoding, state, errors or timestamps of a row.
- Rows below the start block are never changed.
- Bitcoin Core cannot remove an address imported with `importaddress`. An address removed from the list stays in the wallet. It only adds rescan cost.

## Manual repair

Use this when the same error repeats on every tick, the waiting warning stays while the node is in sync, or the configuration was wrong.

1. Stop the LPS. Do not change watcher data while it runs.
2. Back up MongoDB.
3. Check the configuration:
   - `LPS_STAGE`, `CHAIN_ID` and `RSK_ENDPOINT` point to the intended Rootstock network.
   - `PEGIN_ADDRESS_REGISTRY_ADDRESS` is the registry of that network.
   - `PEGIN_ADDRESS_REGISTRY_WATCHER_START_BLOCK` is at or below the block where the registry was deployed. The LPS does not check this.
4. If the error says the node cannot answer calls at old blocks, point `RSK_ENDPOINT` to a node that keeps old state, or go on to step 5, which needs no binary search.
5. Delete the stored point:

   ```javascript
   db.peginWatchCheckpoint.deleteMany({})
   ```

   On the next start the LPS rebuilds the list from the start block. Rows still on the chain keep their import state, and rows the chain does not have are removed.
6. Only if rows were edited or damaged by hand, also delete the rows:

   ```javascript
   db.peginWatch.deleteMany({})
   ```

   Every address is then imported into the Bitcoin wallet again, with a rescan.
7. Start the LPS and wait for `Rebuilt the PegIn watch list ...` in the logs.
8. Compare the stored point with the chain. `db.peginWatchCheckpoint.findOne()` gives `root` and `block`. `getRegistrationRoot()` called at that block must return the same root.
9. If the roots still differ, stop the LPS, record the block and both roots from the error, check the node and the configuration again, then escalate.
