package pegin_test

import (
	"context"
	"encoding/hex"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// nolint:funlen,cyclop
func TestGetQuoteUseCase_Run(t *testing.T) {
	quoteHash := "0x9876543210"
	fedAddress := "fed address"
	lbcAddress := "lbc address"
	lpBtcAddress := "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6"
	userRskAddress := "0x79568c2989232dCa1840087D73d403602364c0D4"
	quoteValue := entities.NewWei(5000)
	quoteData := []byte{1}
	userBtcAddress := "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6"
	lpRskAddress := "0x4b5b6b"
	callFee := entities.NewWei(100)
	penaltyFee := entities.NewWei(50)
	timeForDeposit := uint32(600)
	callTime := uint32(600)
	confirmations := uint16(10)
	gasLimit := entities.NewWei(100)

	request := pegin.NewQuoteRequest(userRskAddress, quoteData, quoteValue, userRskAddress, userBtcAddress)
	quoteMatchFunction := mock.MatchedBy(func(q quote.PeginQuote) bool {
		return q.FedBtcAddress == fedAddress && q.LbcAddress == lbcAddress && q.LpRskAddress == lpRskAddress &&
			q.BtcRefundAddress == userBtcAddress && q.RskRefundAddress == userRskAddress && q.LpBtcAddress == lpBtcAddress &&
			q.CallFee.Cmp(callFee) == 0 && q.PenaltyFee.Cmp(penaltyFee) == 0 && q.ContractAddress == userRskAddress &&
			q.Data == hex.EncodeToString(quoteData) && q.GasLimit == uint32(gasLimit.Uint64()) && q.Value.Cmp(quoteValue) == 0 &&
			q.Nonce > 0 && q.TimeForDeposit == timeForDeposit && q.LpCallTime == callTime && q.Confirmations == confirmations &&
			q.CallOnRegister == false && q.GasFee.Cmp(entities.NewWei(10000)) == 0 && q.ProductFeeAmount == 0
	})

	rsk := new(test.RskRpcMock)
	rsk.On("EstimateGas", mock.Anything, userRskAddress, quoteValue, quoteData).Return(gasLimit, nil)
	rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(100), nil)
	feeCollector := new(test.FeeCollectorMock)
	feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
	bridge := new(test.BridgeMock)
	bridge.On("GetFedAddress").Return(fedAddress, nil)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(200), nil)
	lbc := new(test.LbcMock)
	lbc.On("GetAddress").Return(lbcAddress)
	lbc.On("HashPeginQuote", quoteMatchFunction).Return(quoteHash, nil)
	peginQuoteRepository := new(test.PeginQuoteRepositoryMock)
	peginQuoteRepository.On("InsertQuote", mock.AnythingOfType("context.backgroundCtx"), quoteHash, quoteMatchFunction).Return(nil)
	lp := new(test.ProviderMock)
	lp.On("CallFeePegin").Return(callFee)
	lp.On("PenaltyFeePegin").Return(penaltyFee)
	lp.On("ValidateAmountForPegin", quoteValue).Return(nil)
	lp.On("RskAddress").Return(lpRskAddress)
	lp.On("BtcAddress").Return(lpBtcAddress)
	lp.On("TimeForDepositPegin").Return(timeForDeposit)
	lp.On("CallTime").Return(callTime)
	lp.On("GetBitcoinConfirmationsForValue", quoteValue).Return(confirmations)
	btc := new(test.BtcRpcMock)
	btc.On("ValidateAddress", mock.Anything).Return(nil)
	useCase := pegin.NewGetQuoteUseCase(rsk, btc, feeCollector, bridge, lbc, peginQuoteRepository, lp, lp, "feeCollectorAddress")
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
	rsk := new(test.RskRpcMock)
	lp := new(test.ProviderMock)
	feeCollector := new(test.FeeCollectorMock)
	bridge := new(test.BridgeMock)
	lbc := new(test.LbcMock)
	peginQuoteRepository := new(test.PeginQuoteRepositoryMock)
	feeCollectorAddress := "feeCollectorAddress"
	cases := test.Table[func(btc *test.BtcRpcMock) pegin.QuoteRequest, error]{
		{
			Value: func(btc *test.BtcRpcMock) pegin.QuoteRequest {
				const anyAddress = "any address"
				btc.On("ValidateAddress", anyAddress).Return(blockchain.BtcAddressNotSupportedError)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D4", []byte{1}, entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", anyAddress)
			},
			Result: blockchain.BtcAddressNotSupportedError,
		},
		{
			Value: func(btc *test.BtcRpcMock) pegin.QuoteRequest {
				const anyAddress = "any address"
				btc.On("ValidateAddress", anyAddress).Return(blockchain.BtcAddressInvalidNetworkError)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D4", []byte{1}, entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", anyAddress)
			},
			Result: blockchain.BtcAddressInvalidNetworkError,
		},
		{
			Value: func(btc *test.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				return pegin.NewQuoteRequest("any", []byte{1}, entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			},
			Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *test.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D4", []byte{1}, entities.NewWei(1), "any", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			},
			Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *test.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D4", []byte{1}, entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D41", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			},
			Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *test.BtcRpcMock) pegin.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				return pegin.NewQuoteRequest("0x79568c2989232dCa1840087D73d403602364c0D41", []byte{1}, entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			},
			Result: usecases.RskAddressNotSupportedError,
		},
	}
	for _, testCase := range cases {
		btc := new(test.BtcRpcMock)
		useCase := pegin.NewGetQuoteUseCase(rsk, btc, feeCollector, bridge, lbc, peginQuoteRepository, lp, lp, feeCollectorAddress)
		result, err := useCase.Run(context.Background(), testCase.Value(btc))
		assert.Equal(t, pegin.GetPeginQuoteResult{}, result)
		require.Error(t, err)
		require.ErrorIs(t, err, testCase.Result)
	}
}

func TestGetQuoteUseCase_Run_ErrorHandling(t *testing.T) {
	userRskAddress := "0x79568c2989232dCa1840087D73d403602364c0D4"
	request := pegin.NewQuoteRequest(userRskAddress, []byte{1}, entities.NewWei(5000), userRskAddress, "mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6")

	setups := getQuoteUseCaseUnexpectedErrorSetups()

	for _, setup := range setups {
		rsk := new(test.RskRpcMock)
		feeCollector := new(test.FeeCollectorMock)
		bridge := new(test.BridgeMock)
		lbc := new(test.LbcMock)
		peginQuoteRepository := new(test.PeginQuoteRepositoryMock)
		lp := new(test.ProviderMock)
		btc := new(test.BtcRpcMock)
		btc.On("ValidateAddress", mock.Anything).Return(nil)

		setup(rsk, feeCollector, bridge, lbc, lp, peginQuoteRepository)

		useCase := pegin.NewGetQuoteUseCase(rsk, btc, feeCollector, bridge, lbc, peginQuoteRepository, lp, lp, "feeCollectorAddress")
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
	rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
	lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock,
) {
	return []func(
		rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
		lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock,
	){
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			lp.On("ValidateAmountForPegin", mock.Anything).Return(assert.AnError)
		},
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			lp.On("ValidateAmountForPegin", mock.Anything).Return(nil)
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError)
		},
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			lp.On("ValidateAmountForPegin", mock.Anything).Return(nil)
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(nil, assert.AnError)
		},
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			lp.On("ValidateAmountForPegin", mock.Anything).Return(nil)
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), assert.AnError)
		},
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			lp.On("ValidateAmountForPegin", mock.Anything).Return(nil)
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return("", assert.AnError)
		},
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			lp.On("ValidateAmountForPegin", mock.Anything).Return(nil)
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return("fed address", nil)
			bridge.On("GetMinimumLockTxValue").Return(nil, assert.AnError)
			lbc.On("GetAddress").Return("lbc address")
			lp.On("CallFeePegin").Return(entities.NewWei(100))
			lp.On("PenaltyFeePegin").Return(entities.NewWei(50))
			lp.On("RskAddress").Return("0x4b5b6b")
			lp.On("BtcAddress").Return("mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6")
			lp.On("TimeForDepositPegin").Return(uint32(600))
			lp.On("CallTime").Return(uint32(600))
			lp.On("GetBitcoinConfirmationsForValue", mock.Anything).Return(uint16(10))
		},
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			lp.On("ValidateAmountForPegin", mock.Anything).Return(nil)
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return("fed address", nil)
			bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(200), nil)
			lbc.On("HashPeginQuote", mock.Anything).Return("", assert.AnError)
			lbc.On("GetAddress").Return("lbc address")
			lp.On("CallFeePegin").Return(entities.NewWei(100))
			lp.On("PenaltyFeePegin").Return(entities.NewWei(50))
			lp.On("RskAddress").Return("0x4b5b6b")
			lp.On("BtcAddress").Return("mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6")
			lp.On("TimeForDepositPegin").Return(uint32(600))
			lp.On("CallTime").Return(uint32(600))
			lp.On("GetBitcoinConfirmationsForValue", mock.Anything).Return(uint16(10))
		},
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			lp.On("ValidateAmountForPegin", mock.Anything).Return(nil)
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return("fed address", nil)
			bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(200), nil)
			lbc.On("HashPeginQuote", mock.Anything).Return("any hash", nil)
			lbc.On("GetAddress").Return("lbc address")
			peginQuoteRepository.On("InsertQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything).Return(assert.AnError)
			lp.On("CallFeePegin").Return(entities.NewWei(100))
			lp.On("PenaltyFeePegin").Return(entities.NewWei(50))
			lp.On("RskAddress").Return("0x4b5b6b")
			lp.On("BtcAddress").Return("mnYcQxCZBbmLzNfE9BhV7E8E2u7amdz5y6")
			lp.On("TimeForDepositPegin").Return(uint32(600))
			lp.On("CallTime").Return(uint32(600))
			lp.On("GetBitcoinConfirmationsForValue", mock.Anything).Return(uint16(10))
		},
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil).Once()
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, assert.AnError).Once()
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(10), nil)
			lp.On("ValidateAmountForPegin", mock.Anything).Return(nil)
		},
		func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
			lbc *test.LbcMock, lp *test.ProviderMock, peginQuoteRepository *test.PeginQuoteRepositoryMock) {
			lp.On("ValidateAmountForPegin", mock.Anything).Return(nil)
			rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(entities.NewWei(100), nil)
			rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(10), nil)
			feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			bridge.On("GetFedAddress").Return("fed address", nil)
			lbc.On("GetAddress").Return("")
			lp.On("CallFeePegin").Return(entities.NewWei(0))
			lp.On("PenaltyFeePegin").Return(entities.NewWei(0))
			lp.On("RskAddress").Return("")
			lp.On("BtcAddress").Return("")
			lp.On("TimeForDepositPegin").Return(uint32(0))
			lp.On("CallTime").Return(uint32(0))
			lp.On("GetBitcoinConfirmationsForValue", mock.Anything).Return(uint16(0))
		},
	}
}
