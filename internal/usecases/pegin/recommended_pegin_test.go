package pegin_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

// nolint:funlen
func TestRecommendedPeginUseCase_Run(t *testing.T) {
	var result usecases.RecommendedOperationResult
	amount := entities.NewWei(602180000000000000)
	data := []byte{0x01, 0x02, 0x03}
	request := pegin.NewQuoteRequest(test.AnyRskAddress, data, amount, test.AnyRskAddress)
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.EXPECT().EstimateGas(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(20000), nil)
	rsk.EXPECT().GasPrice(mock.Anything).Return(entities.NewWei(100), nil)
	peginContract := new(mocks.PeginContractMock)
	bridge := new(mocks.BridgeMock)
	bridge.On("GetFedAddress").Return(fedAddress, nil)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(200), nil)
	peginContract.On("GetAddress").Return(lbcAddress)
	peginContract.On("HashPeginQuote", mock.Anything).Return("0x0102030405", nil)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	peginQuoteRepository := new(mocks.PeginQuoteRepositoryMock)
	peginQuoteRepository.EXPECT().InsertQuote(mock.Anything, mock.Anything).Return(nil)
	lp := new(mocks.ProviderMock)
	config := getPeginConfiguration()
	config.MaxValue = entities.NewUWei(math.MaxUint64)
	lp.On("PeginConfiguration", test.AnyCtx).Return(config)
	lp.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil)
	lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration())
	lp.On("RskAddress").Return(test.AnyRskAddress)
	lp.On("BtcAddress").Return(test.AnyBtcAddress)
	btc := new(mocks.BtcRpcMock)
	btc.On("NetworkName").Return(testnetNetworkName)
	contracts := blockchain.RskContracts{PegIn: peginContract, Bridge: bridge}
	rpc := blockchain.Rpc{Btc: btc, Rsk: rsk}
	getQuoteUseCase := pegin.NewGetQuoteUseCase(rpc, contracts, peginQuoteRepository, lp, lp)
	createdQuote, err := getQuoteUseCase.Run(context.Background(), request)
	require.NoError(t, err)
	t.Run("should be consistent with get pegin quote calculation", func(t *testing.T) {
		useCase := pegin.NewRecommendedPeginUseCase(lp, contracts, rpc, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PeginQuote.Total(), test.AnyRskAddress, data)
		require.NoError(t, err)
		assert.Equal(t, createdQuote.PeginQuote.Value, result.RecommendedQuoteValue)
		assert.Equal(t, createdQuote.PeginQuote.CallFee, result.EstimatedCallFee)
		assert.Equal(t, createdQuote.PeginQuote.GasFee, result.EstimatedGasFee)
	})
	t.Run("should use zero address if no destination address is provided", func(t *testing.T) {
		useCase := pegin.NewRecommendedPeginUseCase(lp, contracts, rpc, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PeginQuote.Total(), "", data)
		require.NoError(t, err)
		assert.Equal(t, createdQuote.PeginQuote.Value, result.RecommendedQuoteValue)
		assert.Equal(t, createdQuote.PeginQuote.CallFee, result.EstimatedCallFee)
		assert.Equal(t, createdQuote.PeginQuote.GasFee, result.EstimatedGasFee)
	})
	t.Run("should validate that recommended amount is between provider limits", func(t *testing.T) {
		modifiedConfig := getPeginConfiguration()
		modifiedConfig.MaxValue = modifiedConfig.MinValue
		modifiedLimitLp := new(mocks.ProviderMock)
		modifiedLimitLp.On("PeginConfiguration", mock.Anything).Return(modifiedConfig)
		modifiedLimitLp.On("GeneralConfiguration", mock.Anything).Return(getGeneralConfiguration())
		modifiedLimitLp.On("RskAddress").Return(test.AnyRskAddress)
		modifiedLimitLp.On("BtcAddress").Return(test.AnyBtcAddress)
		modifiedLimitLp.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(nil)
		useCase := pegin.NewRecommendedPeginUseCase(modifiedLimitLp, contracts, rpc, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PeginQuote.Total(), test.AnyRskAddress, data)
		require.ErrorIs(t, err, liquidity_provider.AmountOutOfRangeError)
		assert.Empty(t, result)
	})
	t.Run("should validate liquidity is enough for recommended amount", func(t *testing.T) {
		noLiquidityLp := new(mocks.ProviderMock)
		noLiquidityLp.On("PeginConfiguration", mock.Anything).Return(config)
		noLiquidityLp.On("GeneralConfiguration", mock.Anything).Return(getGeneralConfiguration())
		noLiquidityLp.On("RskAddress").Return(test.AnyRskAddress)
		noLiquidityLp.On("BtcAddress").Return(test.AnyBtcAddress)
		noLiquidityLp.On("HasPeginLiquidity", mock.Anything, mock.Anything).Return(usecases.NoLiquidityError)
		useCase := pegin.NewRecommendedPeginUseCase(noLiquidityLp, contracts, rpc, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PeginQuote.Total(), test.AnyRskAddress, data)
		require.ErrorIs(t, err, usecases.NoLiquidityError)
		assert.Empty(t, result)
	})
	t.Run("should validate recommended amount is over bridge minimum", func(t *testing.T) {
		highMinimumBridge := new(mocks.BridgeMock)
		highMinimumBridge.On("GetMinimumLockTxValue").Return(new(entities.Wei).Add(entities.NewWei(1), createdQuote.PeginQuote.Total()), nil)
		contracts.Bridge = highMinimumBridge
		useCase := pegin.NewRecommendedPeginUseCase(lp, contracts, rpc, utils.Scale)
		result, err = useCase.Run(context.Background(), createdQuote.PeginQuote.Total(), test.AnyRskAddress, data)
		require.ErrorIs(t, err, usecases.TxBelowMinimumError)
		assert.Empty(t, result)
	})
}

func TestRecommendedPeginUseCase_Run_ErrorHandling(t *testing.T) {
	for _, errorSetup := range recommendedPeginErrorSetups() {
		peginContract := new(mocks.PeginContractMock)
		rsk := new(mocks.RootstockRpcServerMock)
		btc := new(mocks.BtcRpcMock)
		lp := new(mocks.ProviderMock)
		lp.On("PeginConfiguration", mock.Anything).Return(getPeginConfiguration())
		contracts := blockchain.RskContracts{PegIn: peginContract}
		rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
		errorSetup(peginContract, rsk)
		useCase := pegin.NewRecommendedPeginUseCase(lp, contracts, rpc, utils.Scale)
		result, err := useCase.Run(context.Background(), entities.NewWei(8000), test.AnyRskAddress, []byte{1, 2, 3})
		assert.Empty(t, result)
		require.Error(t, err)
	}
}

func recommendedPeginErrorSetups() []func(peginContract *mocks.PeginContractMock, rsk *mocks.RootstockRpcServerMock) {
	return []func(peginContract *mocks.PeginContractMock, rsk *mocks.RootstockRpcServerMock){
		func(peginContract *mocks.PeginContractMock, rsk *mocks.RootstockRpcServerMock) {
			rsk.EXPECT().EstimateGas(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
		},
		func(peginContract *mocks.PeginContractMock, rsk *mocks.RootstockRpcServerMock) {
			rsk.EXPECT().EstimateGas(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(1), nil).Once()
			rsk.EXPECT().GasPrice(mock.Anything).Return(nil, assert.AnError).Once()
		},
	}
}
