// Live Bitcoin tests against the local regtest stack.
package pegin_registry_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/configuration/bootstrap/btc_bootstrap"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/integration/regtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	depositBtc    = 0.005
	depositBlocks = 6
)

func newWatchOnlyWallet(
	t *testing.T,
	stack regtest.Stack,
	btc *regtest.Bitcoin,
	prefix string,
) blockchain.BitcoinWallet {
	t.Helper()
	walletID := fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	connection, err := btc_bootstrap.BitcoinWallet(stack.Btc, walletID)
	require.NoError(t, err)
	wallet, err := bitcoin.NewWatchOnlyWallet(connection)
	require.NoError(t, err)
	// Registered after the Bitcoin helper's cleanups, so it runs first, while the node client is
	// still connected.
	t.Cleanup(func() {
		assert.NoError(t, btc.Node().UnloadWallet(&walletID))
	})
	return wallet
}

// fundConfirmedDeposit pays a fresh address and buries the payment, reproducing the order the
// registry enforces: the deposit is already confirmed by the time a registration can be observed.
func fundConfirmedDeposit(t *testing.T, btc *regtest.Bitcoin) (string, string) {
	t.Helper()
	depositAddress := btc.NewAddress(t).EncodeAddress()
	amount, err := btcutil.NewAmount(depositBtc)
	require.NoError(t, err)
	depositTxID := btc.Pay(t, depositAddress, amount)
	btc.Mine(t, depositBlocks)
	return depositAddress, depositTxID
}
