package pegin_test

import (
	"context"
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	fedAddress = "fed address"
	lbcAddress = "lbc address"
)

// nolint:funlen,cyclop
func TestGetQuoteUseCase_Run(t *testing.T) {
	quoteHash := "0x9876543210"
	lpBtcAddress := "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6"
	userRskAddress := "0x79568c2989232dCa1840087D73d403602364c0D4"
	quoteValue := entities.NewWei(5000)
	quoteData := []byte{1}
	userBtcAddress := "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6"
	lpRskAddress := "0x4b5b6b"
	gasLimit := entities.NewWei(100)
	config := getPeginConfiguration()

	request := pegin.NewQuoteRequest(userRskAddress, quoteData, quoteValue, userRskAddress, userBtcAddress)
	quoteMatchFunction := mock.MatchedBy(func(q quote.PeginQuote) bool {
		return q.FedBtcAddress == fedAddress && q.LbcAddress == lbcAddress && q.LpRskAddress == lpRskAddress &&
			q.BtcRefundAddress == userBtcAddress && q.RskRefundAddress == userRskAddress && q.LpBtcAddress == lpBtcAddress &&
			q.CallFee.Cmp(config.CallFee) == 0 && q.PenaltyFee.Cmp(config.PenaltyFee) == 0 && q.ContractAddress == userRskAddress &&
			q.Data == hex.EncodeToString(quoteData) && q.GasLimit == uint32(gasLimit.Uint64()) && q.Value.Cmp(quoteValue) == 0 &&
			q.Nonce > 0 && q.TimeForDeposit == config.TimeForDeposit && q.LpCallTime == config.CallTime && q.Confirmations == 10 &&
			q.CallOnRegister == false && q.GasFee.Cmp(entities.NewWei(10000)) == 0 && q.ProductFeeAmount == 0
	})

	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("EstimateGas", mock.Anything, userRskAddress, quoteValue, quoteData).Return(gasLimit, nil)
	rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(100), nil)
	feeCollector := new(mocks.FeeCollectorMock)
	feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
	bridge := new(mocks.BridgeMock)
	bridge.On("GetFedAddress").Return(fedAddress, nil)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(200), nil)
	lbc := new(mocks.LbcMock)
	lbc.On("GetAddress").Return(lbcAddress)
	lbc.On("HashPeginQuote", quoteMatchFunction).Return(quoteHash, nil)
	peginQuoteRepository := new(mocks.PeginQuoteRepositoryMock)
	peginQuoteRepository.On("InsertQuote", test.AnyCtx, quoteHash, quoteMatchFunction).Return(nil)
	lp := new(mocks.ProviderMock)
	lp.On("PeginConfiguration", test.AnyCtx).Return(config)
	lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration())
	lp.On("RskAddress").Return(lpRskAddress)
	lp.On("BtcAddress").Return(lpBtcAddress)
	btc := new(mocks.BtcRpcMock)
	btc.On("ValidateAddress", mock.Anything).Return(nil)
	contracts := blockchain.RskContracts{Lbc: lbc, FeeCollector: feeCollector, Bridge: bridge}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
	useCase := pegin.NewGetQuoteUseCase(rpc, contracts, peginQuoteRepository, lp, lp, "feeCollectorAddress")
	result, err := useCase.Run(context.Background(), request)

	rsk.AssertExpectations(t)
	feeCollector.AssertExpectations(t)
	bridge.AssertExpectations(t)
	lbc.AssertExpectations(t)
	peginQuoteRepository.AssertExpectations(t)
	lp.AssertExpectations(t)

	assert.NotEmpty(t, result.Hash)
	require.NoError(t, entities.ValidateStruct(result.PeginQuote))
	require.NoError(t, err)
}

func TestGetQuoteUseCase_Run_ValidateRequest(t *testing.T) {
	rsk := new(mocks.RootstockRpcServerMock)
	lp := new(mocks.ProviderMock)
	lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration())
	lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration())
	feeCollector := new(mocks.FeeCollectorMock)
	bridge := new(mocks.BridgeMock)
	lbc := new(mocks.LbcMock)
	peginQuoteRepository := new(mocks.PeginQuoteRepositoryMock)
	cases := validateRequestTestCases()
	for _, testCase := range cases {
		btc := new(mocks.BtcRpcMock)
		contracts := blockchain.RskContracts{Lbc: lbc, FeeCollector: feeCollector, Bridge: bridge}
		rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
		useCase := pegin.NewGetQuoteUseCase(rpc, contracts, peginQuoteRepository, lp, lp, "feeCollectorAddress")
		result, err := useCase.Run(context.Background(), testCase.Value(btc))
		assert.Equal(t, pegin.GetPeginQuoteResult{}, result)
		require.Error(t, err)
		require.ErrorIs(t, err, testCase.Result)
	}
}

func validateRequestTestCases() test.Table[func(btc *mocks.BtcRpcMock) pegin.QuoteRequest, error] {
	return test.Table[func(btc *mocks.BtcRpcMock) pegin.QuoteRequest, error]{
		{
			Value: func(btc *mocks.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", test.AnyAddress).Return(blockchain.BtcAddressNotSupportedError)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D4", []byte{1}, entities.NewWei(1000), "0x79568c2989232dCa1840087D73d403602364c0D4", test.AnyAddress)
			}, Result: blockchain.BtcAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", test.AnyAddress).Return(blockchain.BtcAddressInvalidNetworkError)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D4", []byte{1}, entities.NewWei(1000), "0x79568c2989232dCa1840087D73d403602364c0D4", test.AnyAddress)
			}, Result: blockchain.BtcAddressInvalidNetworkError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				return pegin.NewQuoteRequest("any", []byte{1}, entities.NewWei(1000), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			}, Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D4", []byte{1}, entities.NewWei(1000), "any", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			}, Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D4", []byte{1}, entities.NewWei(1000), "0x79568c2989232dCa1840087D73d403602364c0D41", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			}, Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D41", []byte{1}, entities.NewWei(1000), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			}, Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D4", []byte{1}, entities.NewWei(999), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			}, Result: lpEntity.AmountOutOfRangeError,
		},
	}
}

func TestGetQuoteUseCase_Run_BridgeMinimum(t *testing.T) {
	lp := new(mocks.ProviderMock)
	peginQuoteRepository := new(mocks.PeginQuoteRepositoryMock)
	rsk := new(mocks.RootstockRpcServerMock)
	bridge := new(mocks.BridgeMock)
	feeCollector := new(mocks.FeeCollectorMock)
	btc := new(mocks.BtcRpcMock)
	lbc := new(mocks.LbcMock)
	contracts := blockchain.RskContracts{FeeCollector: feeCollector, Bridge: bridge, Lbc: lbc}
	rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}

	lbc.On("GetAddress").Return(lbcAddress).Once()
	btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(2000), nil).Once()
	lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration()).Once()
	lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration()).Once()
	lp.On("RskAddress").Return(test.AnyAddress).Once()
	lp.On("BtcAddress").Return(test.AnyAddress).Once()
	rsk.EXPECT().EstimateGas(test.AnyCtx, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil).Once()
	rsk.EXPECT().GasPrice(test.AnyCtx).Return(entities.NewWei(10), nil).Once()
	feeCollector.On("DaoFeePercentage").Return(uint64(0), nil).Once()
	bridge.On("GetFedAddress").Return(fedAddress, nil).Once()
	useCase := pegin.NewGetQuoteUseCase(rpc, contracts, peginQuoteRepository, lp, lp, test.AnyAddress)
	t.Run("Should compare bridge minimum against quote value", func(t *testing.T) {
		// we compare 1999 of the quote value with the 2000 of the minimum, so the total is higher than the minimum due to the fees
		quoteValue := entities.NewWei(1999)
		request := pegin.NewQuoteRequest(test.AnyRskAddress, []byte{1}, quoteValue, test.AnyRskAddress, test.AnyBtcAddress)
		result, err := useCase.Run(context.Background(), request)
		assert.Empty(t, result)
		require.ErrorIs(t, err, usecases.TxBelowMinimumError)
	})

	lp.AssertExpectations(t)
	rsk.AssertExpectations(t)
	bridge.AssertExpectations(t)
	feeCollector.AssertExpectations(t)
	btc.AssertExpectations(t)
	peginQuoteRepository.AssertExpectations(t)
}

func TestGetQuoteUseCase_Run_ErrorHandling(t *testing.T) {
	userRskAddress := "0x79568c2989232dCa1840087D73d403602364c0D4"
	request := pegin.NewQuoteRequest(userRskAddress, []byte{1}, entities.NewWei(5000), userRskAddress, "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6")

	setups := getQuoteUseCaseUnexpectedErrorSetups()

	for _, setup := range setups {
		rsk := new(mocks.RootstockRpcServerMock)
		feeCollector := new(mocks.FeeCollectorMock)
		bridge := new(mocks.BridgeMock)
		lbc := new(mocks.LbcMock)
		peginQuoteRepository := new(mocks.PeginQuoteRepositoryMock)
		lp := new(mocks.ProviderMock)
		btc := new(mocks.BtcRpcMock)
		btc.On("ValidateAddress", mock.Anything).Return(nil)

		setup(rsk, feeCollector, bridge, lbc, lp, peginQuoteRepository)
		contracts := blockchain.RskContracts{Lbc: lbc, FeeCollector: feeCollector, Bridge: bridge}
		rpc := blockchain.Rpc{Rsk: rsk, Btc: btc}
		useCase := pegin.NewGetQuoteUseCase(rpc, contracts, peginQuoteRepository, lp, lp, "feeCollectorAddress")
		result, err := useCase.Run(context.Background(), request)
		rsk.AssertExpectations(t)
		feeCollector.AssertExpectations(t)
		bridge.AssertExpectations(t)
		lbc.AssertExpectations(t)
		peginQuoteRepository.AssertExpectations(t)
		lp.AssertExpectations(t)
		assert.Empty(t, result)
		require.Error(t, err)
	}
}

// nolint:funlen
func getQuoteUseCaseUnexpectedErrorSetups() []func(
	rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
	lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock,
) {
	return []func(
		rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
		lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock,
	){
		func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
			lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock) {
			lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration())
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
		},
		func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
			lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock) {
			lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration())
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", test.AnyCtx).Return(nil, assert.AnError)
		},
		func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
			lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock) {
			lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration())
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), assert.AnError)
		},
		func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
			lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock) {
			lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration())
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return("", assert.AnError)
		},
		func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
			lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock) {
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return(fedAddress, nil)
			bridge.On("GetMinimumLockTxValue").Return(nil, assert.AnError)
			lbc.On("GetAddress").Return(lbcAddress)
			lp.On("RskAddress").Return("0x4b5b6b")
			lp.On("BtcAddress").Return("mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6")
			lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration())
			lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration())
		},
		func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
			lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock) {
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return(fedAddress, nil)
			bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(200), nil)
			lbc.On("HashPeginQuote", mock.Anything).Return("", assert.AnError)
			lbc.On("GetAddress").Return(lbcAddress)
			lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration())
			lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration())
			lp.On("RskAddress").Return("0x4b5b6b")
			lp.On("BtcAddress").Return("mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6")
		},
		func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
			lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock) {
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return(fedAddress, nil)
			bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(200), nil)
			lbc.On("HashPeginQuote", mock.Anything).Return("any hash", nil)
			lbc.On("GetAddress").Return(lbcAddress)
			peginQuoteRepository.On("InsertQuote", test.AnyCtx, mock.Anything, mock.Anything).Return(assert.AnError)
			lp.On("RskAddress").Return("0x4b5b6b")
			lp.On("BtcAddress").Return("mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6")
			lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration())
			lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration())
		},
		func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
			lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock) {
			lp.On("PeginConfiguration", test.AnyCtx).Return(getPeginConfiguration())
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil).Once()
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(10), nil)
		},
		func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
			lbc *mocks.LbcMock, lp *mocks.ProviderMock, peginQuoteRepository *mocks.PeginQuoteRepositoryMock) {
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return(fedAddress, nil)
			lbc.On("GetAddress").Return("")
			peginConfig := getPeginConfiguration()
			generalConfig := getGeneralConfiguration()
			peginConfig.CallFee = entities.NewWei(0)
			peginConfig.PenaltyFee = entities.NewWei(0)
			peginConfig.TimeForDeposit = 0
			peginConfig.CallTime = 0
			lp.On("PeginConfiguration", test.AnyCtx).Return(peginConfig)
			lp.On("GeneralConfiguration", test.AnyCtx).Return(generalConfig)
			lp.On("RskAddress").Return("")
			lp.On("BtcAddress").Return("")
		},
	}
}

func getPeginConfiguration() lpEntity.PeginConfiguration {
	return lpEntity.PeginConfiguration{
		TimeForDeposit: 600,
		CallTime:       600,
		PenaltyFee:     entities.NewWei(50),
		CallFee:        entities.NewWei(100),
		MaxValue:       entities.NewWei(10000),
		MinValue:       entities.NewWei(1000),
	}

}

func getGeneralConfiguration() lpEntity.GeneralConfiguration {
	return lpEntity.GeneralConfiguration{
		RskConfirmations: map[int]uint16{1: 10},
		BtcConfirmations: map[int]uint16{1: 10},
	}
}
