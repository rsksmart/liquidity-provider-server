package pegin_registry_test

import (
	"errors"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfirmedDepositIsDiscoverableOnlyAfterARescanImport(t *testing.T) {
	stack := regtest.LoadStack(t)
	btc := regtest.NewBitcoin(t, stack)

	depositAddress, depositTxID := fundConfirmedDeposit(t, btc)
	t.Logf("deposit %s confirmed at %s before any import", depositTxID, depositAddress)

	withoutRescan := newWatchOnlyWallet(t, stack, btc, "pegin-registry-norescan")
	require.NoError(t, withoutRescan.ImportAddress(depositAddress))

	deposits, err := withoutRescan.GetTransactions(depositAddress)
	require.NoError(t, err)
	assert.Empty(t, deposits, "a rescan=false import must not find a deposit that confirmed before it")
	_, err = withoutRescan.GetTransaction(depositTxID)
	require.Error(t, err)
	assert.True(
		t,
		isWalletTransactionNotFound(err),
		"the reconciler demotes on this exact node error, so it must be the one an unaware wallet returns: %v",
		err,
	)

	withRescan := newWatchOnlyWallet(t, stack, btc, "pegin-registry-rescan")
	require.NoError(t, withRescan.ImportAddress(depositAddress))
	fromHeight := max(btc.Height(t)-100, 0)
	_, err = withRescan.RescanBlockchain(fromHeight)
	require.NoError(t, err)

	deposits, err = withRescan.GetTransactions(depositAddress)
	require.NoError(t, err)
	require.Len(t, deposits, 1, "a bounded rescan must find the already-confirmed deposit")
	assert.Equal(t, depositTxID, deposits[0].Hash)

	observed, err := withRescan.GetTransaction(depositTxID)
	require.NoError(t, err)
	assert.Equal(t, depositTxID, observed.Hash)
	assert.GreaterOrEqual(t, observed.Confirmations, uint64(depositBlocks))
}

func isWalletTransactionNotFound(err error) bool {
	var rpcError *btcjson.RPCError
	return errors.As(err, &rpcError) && rpcError.Code == btcjson.ErrRPCInvalidAddressOrKey
}
