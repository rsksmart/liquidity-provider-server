// Live Bitcoin/Mongo tests. Defaults match the local regtest stack; override with BTC_* and MONGODB_*.
package pegin_registry_test

import (
	"fmt"
	"testing"
	"time"

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
