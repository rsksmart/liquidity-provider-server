package cold_wallet_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/cold_wallet"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStaticColdWallet_Init(t *testing.T) {
	t.Run("Should init successfully with valid addresses", func(t *testing.T) {
		btcRpc := &mocks.BtcRpcMock{}
		btcRpc.On("ValidateAddress", test.AnyBtcAddress).Return(nil)
		args := cold_wallet.StaticColdWalletArgs{BtcAddress: test.AnyBtcAddress, RskAddress: test.AnyRskAddress}
		wallet := cold_wallet.NewStaticColdWallet(blockchain.Rpc{Btc: btcRpc}, args)
		err := wallet.Init()
		require.NoError(t, err)
		require.Equal(t, test.AnyBtcAddress, wallet.GetBtcAddress())
		require.Equal(t, test.AnyRskAddress, wallet.GetRskAddress())
	})
	t.Run("Should fail to init with invalid rsk address", func(t *testing.T) {
		btcRpc := &mocks.BtcRpcMock{}
		btcRpc.On("ValidateAddress", test.AnyBtcAddress).Return(nil)
		args := cold_wallet.StaticColdWalletArgs{BtcAddress: test.AnyBtcAddress, RskAddress: "not-an-address"}
		wallet := cold_wallet.NewStaticColdWallet(blockchain.Rpc{Btc: btcRpc}, args)
		err := wallet.Init()
		require.Error(t, err)
	})
	t.Run("Should fail to init with invalid btc address", func(t *testing.T) {
		btcRpc := &mocks.BtcRpcMock{}
		btcRpc.On("ValidateAddress", mock.Anything).Return(assert.AnError)
		args := cold_wallet.StaticColdWalletArgs{BtcAddress: "not-an-address", RskAddress: test.AnyRskAddress}
		wallet := cold_wallet.NewStaticColdWallet(blockchain.Rpc{Btc: btcRpc}, args)
		err := wallet.Init()
		require.Error(t, err)
	})
}
