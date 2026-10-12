package regtest

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/btc_bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/require"
)

const (
	// FundingWalletID is the wallet the tests pay from. It is never one of the stack's wallets, which
	// may be encrypted or belong to the LPS.
	FundingWalletID = "regtest-integration-funding"
	// The node runs without -fallbackfee, so regtest fee estimation has nothing to answer with.
	fundingFeeRate   = btcutil.Amount(10000)
	coinbaseMaturity = 101
	// Covers the fee of a payment on top of its amount when checking the balance.
	paymentFeeMargin = btcutil.Amount(btcutil.SatoshiPerBitcoin / 100)
)

// Proof is the SPV proof of a transaction, built the way the LPS builds it for requestPegIn.
type Proof struct {
	RawTx     []byte
	BlockHash [32]byte
	Branch    blockchain.MerkleBranch
}

// Bitcoin acts as the user and the miner on the regtest bitcoind.
type Bitcoin struct {
	env     Stack
	node    *rpcclient.Client
	wallet  *rpcclient.Client
	network blockchain.BitcoinNetwork
}

// NewBitcoin connects to bitcoind and opens the funding wallet. The clients close when the test ends.
func NewBitcoin(t *testing.T, stack Stack) *Bitcoin {
	t.Helper()
	node := newRpcClient(t, stack, "")
	t.Cleanup(node.Shutdown)
	wallet := openFundingWallet(t, stack, node)
	t.Cleanup(wallet.Shutdown)

	connection, err := btc_bootstrap.Bitcoin(stack.Btc)
	require.NoError(t, err)
	t.Cleanup(func() { connection.Shutdown(make(chan bool, 1)) })

	return &Bitcoin{env: stack, node: node, wallet: wallet, network: bitcoin.NewBitcoindRpc(connection)}
}

// Node is the walletless client, for tests that need a node call this helper does not cover.
func (btc *Bitcoin) Node() *rpcclient.Client {
	return btc.node
}

// Wallet is the funding wallet client, for tests that need a wallet call this helper does not cover.
func (btc *Bitcoin) Wallet() *rpcclient.Client {
	return btc.wallet
}

func (btc *Bitcoin) Params(t *testing.T) *chaincfg.Params {
	t.Helper()
	params, err := btc.env.Btc.GetNetworkParams()
	require.NoError(t, err)
	return params
}

func (btc *Bitcoin) Height(t *testing.T) int64 {
	t.Helper()
	height, err := btc.node.GetBlockCount()
	require.NoError(t, err)
	return height
}

// Mine adds blocks paying the coinbase to the funding wallet.
func (btc *Bitcoin) Mine(t *testing.T, blocks int64) {
	t.Helper()
	_, err := btc.node.GenerateToAddress(blocks, btc.NewAddress(t), nil)
	require.NoError(t, err)
}

func (btc *Bitcoin) NewAddress(t *testing.T) btcutil.Address {
	t.Helper()
	address, err := btc.wallet.GetNewAddress("")
	require.NoError(t, err)
	return address
}

// EnsureFunds mines until the funding wallet can spend amount.
func (btc *Bitcoin) EnsureFunds(t *testing.T, amount btcutil.Amount) {
	t.Helper()
	for range 3 {
		balance, err := btc.wallet.GetBalance("*")
		require.NoError(t, err)
		if balance >= amount {
			return
		}
		btc.Mine(t, coinbaseMaturity)
	}
	require.FailNowf(t, "the funding wallet cannot be funded", "it still holds less than %s", amount)
}

// Pay sends amount to address from the funding wallet and returns the transaction id. It does not mine.
func (btc *Bitcoin) Pay(t *testing.T, address string, amount btcutil.Amount) string {
	t.Helper()
	btc.EnsureFunds(t, amount+paymentFeeMargin)
	decoded, err := btcutil.DecodeAddress(address, btc.Params(t))
	require.NoErrorf(t, err, "%s is not a %s address", address, btc.env.Btc.Network)
	txID, err := btc.wallet.SendToAddress(decoded, amount)
	require.NoErrorf(t, err, "cannot pay %s to %s", amount, address)
	t.Logf("paid %s to %s in %s", amount, address, txID)
	return txID.String()
}

// Confirmations is 0 while the transaction is in the mempool.
func (btc *Bitcoin) Confirmations(t *testing.T, txID string) uint64 {
	t.Helper()
	hash, err := chainhash.NewHashFromStr(txID)
	require.NoError(t, err)
	tx, err := btc.node.GetRawTransactionVerbose(hash)
	require.NoErrorf(t, err, "bitcoind does not know transaction %s", txID)
	return tx.Confirmations
}

// Proof builds the SPV proof of a mined transaction with the LPS's own Bitcoin adapter.
func (btc *Bitcoin) Proof(t *testing.T, txID string) Proof {
	t.Helper()
	rawTx, err := btc.network.GetRawTransaction(txID)
	require.NoError(t, err)
	require.NoError(t, blockchain.RejectWitnessSerializedTx(rawTx))
	block, err := btc.network.GetTransactionBlockInfo(txID)
	require.NoErrorf(t, err, "transaction %s is not mined yet", txID)
	branch, err := btc.network.BuildMerkleBranch(txID)
	require.NoError(t, err)
	return Proof{RawTx: rawTx, BlockHash: block.Hash, Branch: branch}
}

func newRpcClient(t *testing.T, stack Stack, walletID string) *rpcclient.Client {
	t.Helper()
	host := stack.Btc.Endpoint
	if walletID != "" {
		host = fmt.Sprintf("%s/wallet/%s", host, walletID)
	}
	params, err := stack.Btc.GetNetworkParams()
	require.NoError(t, err)
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         host,
		User:         stack.Btc.Username,
		Pass:         stack.Btc.Password,
		Params:       params.Name,
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)
	require.NoError(t, err)
	return client
}

func openFundingWallet(t *testing.T, stack Stack, node *rpcclient.Client) *rpcclient.Client {
	t.Helper()
	// Loading fails when the wallet is already loaded and creating fails when it already exists, so
	// neither error decides anything on its own; the wallet-scoped probe below does.
	var createError error
	_, loadError := node.LoadWallet(FundingWalletID)
	if loadError != nil {
		_, createError = node.CreateWallet(FundingWalletID)
	}
	wallet := newRpcClient(t, stack, FundingWalletID)
	_, err := wallet.GetWalletInfo()
	require.NoErrorf(t, err, "bitcoind is not reachable at %s, or the funding wallet was neither loaded (%v) nor created (%v)",
		stack.Btc.Endpoint, loadError, createError)
	require.NoError(t, wallet.SetTxFee(fundingFeeRate))
	return wallet
}
