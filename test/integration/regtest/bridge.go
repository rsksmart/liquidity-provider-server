package regtest

import (
	"encoding/hex"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	bridgeBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/bridge"
	"github.com/stretchr/testify/require"
)

// Bridge reads what the Rootstock bridge knows about Bitcoin. It only learns BTC headers from
// powpeg, so its view lags bitcoind until powpeg relays them.
type Bridge struct {
	rsk     *Rootstock
	address common.Address
	binding *bridgeBindings.RskBridge
}

func NewBridge(stack Stack, rsk *Rootstock) *Bridge {
	return &Bridge{rsk: rsk, address: stack.Contracts.Bridge, binding: bridgeBindings.NewRskBridge()}
}

func (bridge *Bridge) BtcHeight(t *testing.T) int64 {
	t.Helper()
	height, err := bridge.binding.UnpackGetBtcBlockchainBestChainHeight(bridge.call(t, bridge.binding.PackGetBtcBlockchainBestChainHeight()))
	require.NoError(t, err)
	return height.Int64()
}

// TxConfirmations is how many confirmations the bridge sees for a mined transaction. Negative
// values are the bridge's error codes, for example when it does not have the block yet.
func (bridge *Bridge) TxConfirmations(t *testing.T, txID string, proof Proof) int64 {
	t.Helper()
	// The bridge takes both hashes in display order, as bitcoind prints them.
	txHash, err := hex.DecodeString(txID)
	require.NoError(t, err)
	require.Len(t, txHash, 32)
	data := bridge.binding.PackGetBtcTransactionConfirmations([32]byte(txHash), proof.BlockHash, proof.Branch.Path, proof.Branch.Hashes)
	confirmations, err := bridge.binding.UnpackGetBtcTransactionConfirmations(bridge.call(t, data))
	require.NoError(t, err)
	return confirmations.Int64()
}

func (bridge *Bridge) call(t *testing.T, data []byte) []byte {
	t.Helper()
	output, err := bridge.rsk.Client().CallContract(t.Context(), ethereum.CallMsg{To: &bridge.address, Data: data}, nil)
	require.NoError(t, err, "the bridge call failed")
	return output
}
