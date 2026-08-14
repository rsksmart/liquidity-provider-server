// Package pegin_registry_test holds the registry discovery claims that only a live node can settle:
// how a watch-only wallet answers about a deposit that confirmed before the address was imported,
// and how the watch set behaves under a replayed event.
//
// It is outside the `make test` scope because it needs running Bitcoin and Mongo services. Point it
// at them through BTC_NETWORK, BTC_USERNAME, BTC_PASSWORD, BTC_ENDPOINT, MONGODB_HOST, MONGODB_PORT,
// MONGODB_USER and MONGODB_PASSWORD; the defaults match the repository's local regtest stack.
package pegin_registry_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/btc_bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	defaultBtcEndpoint = "127.0.0.1:5555"
	defaultBtcNetwork  = "regtest"
	defaultBtcUser     = "test"
	defaultBtcPassword = "test"

	fundingWalletID    = "pegin-registry-funding"
	depositBtc         = 0.005
	depositBlocks      = 6
	coinbaseMaturity   = 101
	minimumFundingCoin = btcutil.Amount(btcutil.SatoshiPerBitcoin)
	fixedFeeRate       = btcutil.Amount(10000)
)

func TestConfirmedDepositIsDiscoverableOnlyAfterARescanImport(t *testing.T) {
	env := btcEnvironment(t)
	nodeClient := newRpcClient(t, env, "")
	fundingClient := openFundingWallet(t, env, nodeClient)
	// Registered before the wallets below so that LIFO cleanup unloads them while the node client
	// is still connected.
	t.Cleanup(nodeClient.Shutdown)
	t.Cleanup(fundingClient.Shutdown)

	depositAddress, depositTxID := fundConfirmedDeposit(t, nodeClient, fundingClient)
	t.Logf("deposit %s confirmed at %s before any import", depositTxID, depositAddress)

	// Importing without a rescan is what leaves a real, confirmed payment unreported.
	withoutRescan := newWatchOnlyWallet(t, env, nodeClient, "pegin-registry-norescan")
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

	withRescan := newWatchOnlyWallet(t, env, nodeClient, "pegin-registry-rescan")
	require.NoError(t, withRescan.ImportAddress(depositAddress))
	tip, err := nodeClient.GetBlockCount()
	require.NoError(t, err)
	fromHeight := tip - 100
	if fromHeight < 0 {
		fromHeight = 0
	}
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

func btcEnvironment(t *testing.T) environment.BtcEnv {
	t.Helper()
	var env environment.Environment
	require.NoError(t, environment.Load(&env))
	btc := env.Btc
	if btc.Network == "" {
		btc.Network = defaultBtcNetwork
	}
	if btc.Username == "" {
		btc.Username = defaultBtcUser
	}
	if btc.Password == "" {
		btc.Password = defaultBtcPassword
	}
	if btc.Endpoint == "" {
		btc.Endpoint = defaultBtcEndpoint
	}
	return btc
}

// newRpcClient talks to the node directly, so the test can act as the paying user and the miner
// without borrowing capabilities the watch-only wallet deliberately does not have.
func newRpcClient(t *testing.T, env environment.BtcEnv, walletID string) *rpcclient.Client {
	t.Helper()
	host := env.Endpoint
	if walletID != "" {
		host = fmt.Sprintf("%s/wallet/%s", host, walletID)
	}
	params, err := env.GetNetworkParams()
	require.NoError(t, err)
	client, err := rpcclient.New(&rpcclient.ConnConfig{
		Host:         host,
		User:         env.Username,
		Pass:         env.Password,
		Params:       params.Name,
		DisableTLS:   true,
		HTTPPostMode: true,
	}, nil)
	require.NoError(t, err)
	return client
}

func newWatchOnlyWallet(
	t *testing.T,
	env environment.BtcEnv,
	nodeClient *rpcclient.Client,
	prefix string,
) blockchain.BitcoinWallet {
	t.Helper()
	walletID := fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	connection, err := btc_bootstrap.BitcoinWallet(env, walletID)
	require.NoError(t, err)
	wallet, err := bitcoin.NewWatchOnlyWallet(connection)
	require.NoError(t, err)
	t.Cleanup(func() {
		assert.NoError(t, nodeClient.UnloadWallet(&walletID))
	})
	return wallet
}

// openFundingWallet gives the test its own spending wallet rather than borrowing whatever the stack
// already holds, which may be encrypted or in use by another exercise.
func openFundingWallet(t *testing.T, env environment.BtcEnv, nodeClient *rpcclient.Client) *rpcclient.Client {
	t.Helper()
	// Loading fails when the wallet is already loaded and creating fails when it already exists, so
	// neither error decides anything on its own; the wallet-scoped probe below does.
	var createError error
	_, loadError := nodeClient.LoadWallet(fundingWalletID)
	if loadError != nil {
		_, createError = nodeClient.CreateWallet(fundingWalletID)
	}
	fundingClient := newRpcClient(t, env, fundingWalletID)
	_, err := fundingClient.GetWalletInfo()
	require.NoErrorf(t, err, "the funding wallet was neither loaded (%v) nor created (%v)", loadError, createError)
	// The stack's node runs without -fallbackfee, so regtest fee estimation has nothing to answer
	// with. Fixing the wallet's rate keeps the deviation inside the test instead of the node config.
	require.NoError(t, fundingClient.SetTxFee(fixedFeeRate))
	return fundingClient
}

// fundConfirmedDeposit pays a fresh address and buries the payment, reproducing the order the
// registry enforces: the deposit is already confirmed by the time a registration can be observed.
func fundConfirmedDeposit(t *testing.T, nodeClient, fundingClient *rpcclient.Client) (string, string) {
	t.Helper()
	miningAddress, err := fundingClient.GetNewAddress("")
	require.NoError(t, err)
	balance, err := fundingClient.GetBalance("*")
	require.NoError(t, err)
	if balance < minimumFundingCoin {
		_, err = nodeClient.GenerateToAddress(coinbaseMaturity, miningAddress, nil)
		require.NoError(t, err)
	}

	depositAddress, err := fundingClient.GetNewAddress("")
	require.NoError(t, err)
	amount, err := btcutil.NewAmount(depositBtc)
	require.NoError(t, err)
	depositTxID, err := fundingClient.SendToAddress(depositAddress, amount)
	require.NoError(t, err)
	_, err = nodeClient.GenerateToAddress(depositBlocks, miningAddress, nil)
	require.NoError(t, err)
	return depositAddress.EncodeAddress(), depositTxID.String()
}

func isWalletTransactionNotFound(err error) bool {
	var rpcError *btcjson.RPCError
	return errors.As(err, &rpcError) && rpcError.Code == btcjson.ErrRPCInvalidAddressOrKey
}
