package regtest

import (
	"fmt"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/btc_bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/environment"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/require"
)

const (
	defaultBtcEndpoint = "127.0.0.1:5555"
	defaultBtcNetwork  = "regtest"
	defaultBtcUser     = "test"
	defaultBtcPassword = "test"

	fundingWalletID    = "pegin-commit-first-funding"
	depositBtc         = 0.005
	coinbaseMaturity   = 101
	minimumFundingCoin = btcutil.Amount(btcutil.SatoshiPerBitcoin)
	fixedFeeRate       = btcutil.Amount(10000)
)

type BitcoinStack struct {
	Env           environment.BtcEnv
	Node          *rpcclient.Client
	Funding       *rpcclient.Client
	Rpc           blockchain.BitcoinNetwork
	MiningAddress btcutil.Address
}

func OpenBitcoin(t *testing.T) *BitcoinStack {
	t.Helper()
	env := btcHostEnv(t)
	node := newRpcClient(t, env, "")
	funding := openFundingWallet(t, env, node)
	t.Cleanup(node.Shutdown)
	t.Cleanup(funding.Shutdown)

	conn, err := btc_bootstrap.Bitcoin(env)
	require.NoError(t, err)
	rpc := bitcoin.NewBitcoindRpc(conn)

	miningAddress, err := funding.GetNewAddress("")
	require.NoError(t, err)
	balance, err := funding.GetBalance("*")
	require.NoError(t, err)
	if balance < minimumFundingCoin {
		_, err = node.GenerateToAddress(coinbaseMaturity, miningAddress, nil)
		require.NoError(t, err)
	}
	return &BitcoinStack{
		Env:           env,
		Node:          node,
		Funding:       funding,
		Rpc:           rpc,
		MiningAddress: miningAddress,
	}
}

func (stack *BitcoinStack) PayAndMine(t *testing.T, btcAddr string, blocks int64) string {
	t.Helper()
	params, err := stack.Env.GetNetworkParams()
	require.NoError(t, err)
	decoded, err := btcutil.DecodeAddress(btcAddr, params)
	require.NoError(t, err)
	amount, err := btcutil.NewAmount(depositBtc)
	require.NoError(t, err)
	depositTxID, err := stack.Funding.SendToAddress(decoded, amount)
	require.NoError(t, err)
	stack.Mine(t, blocks)
	return depositTxID.String()
}

func (stack *BitcoinStack) Mine(t *testing.T, blocks int64) {
	t.Helper()
	_, err := stack.Node.GenerateToAddress(blocks, stack.MiningAddress, nil)
	require.NoError(t, err)
}

func (stack *BitcoinStack) Height(t *testing.T) int64 {
	t.Helper()
	height, err := stack.Rpc.GetHeight()
	require.NoError(t, err)
	return height.Int64()
}

func btcHostEnv(t *testing.T) environment.BtcEnv {
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
	if btc.Endpoint == "" || btc.Endpoint == "bitcoind:5555" {
		btc.Endpoint = defaultBtcEndpoint
	}
	return btc
}

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

func openFundingWallet(t *testing.T, env environment.BtcEnv, nodeClient *rpcclient.Client) *rpcclient.Client {
	t.Helper()
	var createError error
	_, loadError := nodeClient.LoadWallet(fundingWalletID)
	if loadError != nil {
		_, createError = nodeClient.CreateWallet(fundingWalletID)
	}
	fundingClient := newRpcClient(t, env, fundingWalletID)
	_, err := fundingClient.GetWalletInfo()
	require.NoErrorf(t, err, "the funding wallet was neither loaded (%v) nor created (%v)", loadError, createError)
	require.NoError(t, fundingClient.SetTxFee(fixedFeeRate))
	return fundingClient
}
