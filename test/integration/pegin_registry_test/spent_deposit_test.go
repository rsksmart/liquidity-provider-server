package pegin_registry_test

import (
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// spendFeeBtc leaves the node a fee well above the regtest relay minimum without depending on fee
// estimation, which the stack's node cannot answer.
const spendFeeBtc = 0.0001

// A reconciler reading the address's unspent list would demote a valid entry the moment its deposit
// output is spent, and erase the txid with it. Wallet history keeps the transaction. Both halves are
// properties of the node rather than of this repository, so only a live node can settle them.
func TestSpentDepositStaysVisibleInWalletHistory(t *testing.T) {
	env := btcEnvironment(t)
	nodeClient := newRpcClient(t, env, "")
	fundingClient := openFundingWallet(t, env, nodeClient)
	// Registered before the wallet below so that LIFO cleanup unloads it while the node client is
	// still connected.
	t.Cleanup(nodeClient.Shutdown)
	t.Cleanup(fundingClient.Shutdown)

	depositAddress, depositTxID := fundConfirmedDeposit(t, nodeClient, fundingClient)
	monitoringWallet := newWatchOnlyWallet(t, env, nodeClient, "pegin-registry-spent")
	require.NoError(t, monitoringWallet.ImportAddress(depositAddress))
	tip, err := nodeClient.GetBlockCount()
	require.NoError(t, err)
	fromHeight := tip - 100
	if fromHeight < 0 {
		fromHeight = 0
	}
	_, err = monitoringWallet.RescanBlockchain(fromHeight)
	require.NoError(t, err)

	deposits, err := monitoringWallet.GetTransactions(depositAddress)
	require.NoError(t, err)
	require.Len(t, deposits, 1, "the imported address must hold the confirmed deposit before it is spent")
	require.Equal(t, depositTxID, deposits[0].Hash)

	spendDepositOutput(t, nodeClient, fundingClient, depositAddress, depositTxID)

	deposits, err = monitoringWallet.GetTransactions(depositAddress)
	require.NoError(t, err)
	assert.Empty(t, deposits, "spending the deposit empties the unspent view the reconciler must not rely on")

	observed, err := monitoringWallet.GetTransaction(depositTxID)
	require.NoError(t, err)
	assert.Equal(t, depositTxID, observed.Hash)
	assert.Positive(t, observed.Confirmations, "a spent deposit keeps its confirmations in wallet history")
}

// Spending exactly the deposit output, rather than letting the node choose inputs, keeps the
// assertions that follow about that output rather than about the wallet's coin selection.
func spendDepositOutput(
	t *testing.T,
	nodeClient, fundingClient *rpcclient.Client,
	depositAddress, depositTxID string,
) {
	t.Helper()
	hash, err := chainhash.NewHashFromStr(depositTxID)
	require.NoError(t, err)
	deposit, err := fundingClient.GetTransaction(hash)
	require.NoError(t, err)

	var received *btcjson.GetTransactionDetailsResult
	for index, detail := range deposit.Details {
		if detail.Address == depositAddress && detail.Category == "receive" {
			received = &deposit.Details[index]
			break
		}
	}
	require.NotNil(t, received, "the funding wallet must own the deposit output to be able to spend it")

	destination, err := fundingClient.GetNewAddress("")
	require.NoError(t, err)
	amount, err := btcutil.NewAmount(received.Amount - spendFeeBtc)
	require.NoError(t, err)
	spend, err := fundingClient.CreateRawTransaction(
		[]btcjson.TransactionInput{{Txid: depositTxID, Vout: received.Vout}},
		map[btcutil.Address]btcutil.Amount{destination: amount},
		nil,
	)
	require.NoError(t, err)
	signed, complete, err := fundingClient.SignRawTransactionWithWallet(spend)
	require.NoError(t, err)
	require.True(t, complete, "the funding wallet must be able to sign its own deposit output")
	_, err = fundingClient.SendRawTransaction(signed, false)
	require.NoError(t, err)
	_, err = nodeClient.GenerateToAddress(1, destination, nil)
	require.NoError(t, err)
}
