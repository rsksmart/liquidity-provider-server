package pegout_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// nolint:funlen
func TestRecommendedPegoutUseCase_Run(t *testing.T) {
	const (
		btcAddress = "mvTRXpLSD9Y8UQc18ciWnn4zxT23U1pBbK"
		rskAddress = "0xBb519e5dCB3f98ED0c48238b42BFa3fd4d1a5E45"
	)
	var result usecases.RecommendedOperationResult

	amount := entities.NewWei(602247200000000000)
	request := pegout.NewQuoteRequest(btcAddress, amount, rskAddress)
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.EXPECT().GasPrice(mock.Anything).Return(entities.NewWei(1), nil)
	rsk.EXPECT().GetHeight(mock.Anything).Return(uint64(100), nil)
	rsk.EXPECT().EstimateGas(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(4600), nil)
	pegoutContract := new(mocks.PegoutContractMock)
	pegoutContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	bridge := new(mocks.BridgeMock)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(5), nil)
	pegoutContract.On("GetAddress").Return(rskAddress)
	pegoutContract.On("HashPegoutQuote", mock.Anything).Return("0x9876543210", nil)
	pegoutQuoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	pegoutQuoteRepository.EXPECT().InsertQuote(mock.Anything, mock.Anything).Return(nil)
	lp := new(mocks.ProviderMock)
	lp.On("PegoutConfiguration", mock.Anything).Return(getPegoutConfiguration())
	lp.On("GeneralConfiguration", mock.Anything).Return(getGeneralConfiguration())
	lp.On("RskAddress").Return(test.AnyRskAddress)
	lp.On("BtcAddress").Return(test.AnyBtcAddress)
	lp.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(nil)
	btcWallet := new(mocks.BitcoinWalletMock)
	btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(blockchain.BtcFeeEstimation{
		Value:   entities.NewWei(67250000000000),
		FeeRate: utils.NewBigFloat64(25),
	}, nil)
	btc := new(mocks.BtcRpcMock)
	btc.On("ValidateAddress", mock.Anything).Return(nil)

	contracts := blockchain.RskContracts{PegOut: pegoutContract, Bridge: bridge}
	rpc := blockchain.Rpc{Btc: btc, Rsk: rsk}
	getQuoteUseCase := pegout.NewGetQuoteUseCase(rpc, contracts, pegoutQuoteRepository, lp, lp, btcWallet)
	createdQuote, err := getQuoteUseCase.Run(context.Background(), request)
	require.NoError(t, err)
	t.Run("should be consistent with get pegout quote calculation", func(t *testing.T) {
		btc.On("GetZeroAddress", mock.Anything).Return(blockchain.BitcoinTestnetP2SHZeroAddress, nil).Once()
		useCase := pegout.NewRecommendedPegoutUseCase(lp, contracts, rpc, btcWallet, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PegoutQuote.Total(), blockchain.BtcAddressTypeP2SH)
		require.NoError(t, err)
		assert.Equal(t, createdQuote.PegoutQuote.Value, result.RecommendedQuoteValue)
		assert.Equal(t, createdQuote.PegoutQuote.CallFee, result.EstimatedCallFee)
		assert.Equal(t, createdQuote.PegoutQuote.GasFee, result.EstimatedGasFee)
	})
	t.Run("should use P2PKH if no destination type is provided", func(t *testing.T) {
		btc.On("GetZeroAddress", blockchain.BtcAddressType("p2pkh")).Return(blockchain.BitcoinTestnetP2PKHZeroAddress, nil).Once()
		useCase := pegout.NewRecommendedPegoutUseCase(lp, contracts, rpc, btcWallet, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PegoutQuote.Total(), "")
		require.NoError(t, err)
		assert.Equal(t, createdQuote.PegoutQuote.Value, result.RecommendedQuoteValue)
		assert.Equal(t, createdQuote.PegoutQuote.CallFee, result.EstimatedCallFee)
		assert.Equal(t, createdQuote.PegoutQuote.GasFee, result.EstimatedGasFee)
	})
	t.Run("should validate that recommended amount is between provider limits", func(t *testing.T) {
		btc.On("GetZeroAddress", mock.Anything).Return(blockchain.BitcoinTestnetP2PKHZeroAddress, nil).Once()
		config := getPegoutConfiguration()
		config.MaxValue = config.MinValue
		modifiedLimitLp := new(mocks.ProviderMock)
		modifiedLimitLp.On("PegoutConfiguration", mock.Anything).Return(config)
		modifiedLimitLp.On("GeneralConfiguration", mock.Anything).Return(getGeneralConfiguration())
		modifiedLimitLp.On("RskAddress").Return(test.AnyRskAddress)
		modifiedLimitLp.On("BtcAddress").Return(test.AnyBtcAddress)
		useCase := pegout.NewRecommendedPegoutUseCase(modifiedLimitLp, contracts, rpc, btcWallet, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PegoutQuote.Total(), blockchain.BtcAddressTypeP2PKH)
		require.ErrorIs(t, err, liquidity_provider.AmountOutOfRangeError)
		assert.Empty(t, result)
	})
	t.Run("should validate liquidity is enough for recommended amount", func(t *testing.T) {
		btc.On("GetZeroAddress", mock.Anything).Return(blockchain.BitcoinTestnetP2PKHZeroAddress, nil).Once()
		noLiquidityLp := new(mocks.ProviderMock)
		noLiquidityLp.On("PegoutConfiguration", mock.Anything).Return(getPegoutConfiguration())
		noLiquidityLp.On("GeneralConfiguration", mock.Anything).Return(getGeneralConfiguration())
		noLiquidityLp.On("RskAddress").Return(test.AnyRskAddress)
		noLiquidityLp.On("BtcAddress").Return(test.AnyBtcAddress)
		noLiquidityLp.On("HasPegoutLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError)
		useCase := pegout.NewRecommendedPegoutUseCase(noLiquidityLp, contracts, rpc, btcWallet, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PegoutQuote.Total(), blockchain.BtcAddressTypeP2PKH)
		require.ErrorIs(t, err, usecases.NoLiquidityError)
		assert.Empty(t, result)
	})
	t.Run("should validate recommended amount is over bridge minimum", func(t *testing.T) {
		btc.On("GetZeroAddress", mock.Anything).Return(blockchain.BitcoinTestnetP2PKHZeroAddress, nil).Once()
		highMinimumBridge := new(mocks.BridgeMock)
		highMinimumBridge.On("GetMinimumLockTxValue").Return(new(entities.Wei).Add(entities.NewWei(1), createdQuote.PegoutQuote.Total()), nil)
		contracts.Bridge = highMinimumBridge
		useCase := pegout.NewRecommendedPegoutUseCase(lp, contracts, rpc, btcWallet, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PegoutQuote.Total(), blockchain.BtcAddressTypeP2PKH)
		require.ErrorIs(t, err, usecases.TxBelowMinimumError)
		assert.Empty(t, result)
	})
}

func TestRecommendedPegoutUseCase_Run_ErrorHandling(t *testing.T) {
	for _, errorSetup := range recommendedPegoutErrorSetups() {
		pegoutContract := new(mocks.PegoutContractMock)
		btc := new(mocks.BtcRpcMock)
		btcWallet := new(mocks.BitcoinWalletMock)
		lp := new(mocks.ProviderMock)
		lp.On("PegoutConfiguration", mock.Anything).Return(getPegoutConfiguration())
		contracts := blockchain.RskContracts{PegOut: pegoutContract}
		rpc := blockchain.Rpc{Btc: btc}
		errorSetup(pegoutContract, btc, btcWallet, lp)
		useCase := pegout.NewRecommendedPegoutUseCase(lp, contracts, rpc, btcWallet, utils.Scale)
		result, err := useCase.Run(context.Background(), entities.NewWei(1), blockchain.BtcAddressTypeP2PKH)
		assert.Empty(t, result)
		require.Error(t, err)
	}
}

func recommendedPegoutErrorSetups() []func(pegoutContract *mocks.PegoutContractMock, btc *mocks.BtcRpcMock, btcWallet *mocks.BitcoinWalletMock, lp *mocks.ProviderMock) {
	return []func(pegoutContract *mocks.PegoutContractMock, btc *mocks.BtcRpcMock, btcWallet *mocks.BitcoinWalletMock, lp *mocks.ProviderMock){
		func(pegoutContract *mocks.PegoutContractMock, btc *mocks.BtcRpcMock, btcWallet *mocks.BitcoinWalletMock, lp *mocks.ProviderMock) {
			btc.On("GetZeroAddress", mock.Anything).Return("", assert.AnError)
		},
		func(pegoutContract *mocks.PegoutContractMock, btc *mocks.BtcRpcMock, btcWallet *mocks.BitcoinWalletMock, lp *mocks.ProviderMock) {
			btc.On("GetZeroAddress", mock.Anything).Return(blockchain.BitcoinTestnetP2PKHZeroAddress, nil)
			btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(blockchain.BtcFeeEstimation{}, assert.AnError)
		},
	}
}
