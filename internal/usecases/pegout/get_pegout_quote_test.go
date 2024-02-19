package pegout_test

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	pegout "github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetQuoteUseCase_Run(t *testing.T) {
	rsk := new(test.RskRpcMock)
	rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
	rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
	feeCollector := new(test.FeeCollectorMock)
	feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
	bridge := new(test.BridgeMock)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000000), nil)
	lbc := new(test.LbcMock)
	lbc.On("GetAddress").Return("0x1234")
	lbc.On("HashPegoutQuote", mock.Anything).Return("0x9876543210", nil)
	pegoutQuoteRepository := new(test.PegoutQuoteRepositoryMock)
	pegoutQuoteRepository.On("InsertQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything).Return(nil)
	lp := new(test.ProviderMock)
	lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
	lp.On("GetRootstockConfirmationsForValue", mock.Anything).Return(uint16(10))
	lp.On("GetBitcoinConfirmationsForValue", mock.Anything, mock.Anything).Return(uint16(10))
	lp.On("CallFeePegout").Return(entities.NewWei(200))
	lp.On("PenaltyFeePegout").Return(entities.NewWei(20))
	lp.On("RskAddress").Return("0x1234")
	lp.On("BtcAddress").Return("address")
	lp.On("TimeForDepositPegout").Return(uint32(60000))
	lp.On("ExpireBlocksPegout").Return(uint64(60000))
	btcWallet := new(test.BtcWalletMock)
	btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000000000000000), nil)
	feeCollectorAddress := "feeCollectorAddress"
	useCase := pegout.NewGetQuoteUseCase(rsk, feeCollector, bridge, lbc, pegoutQuoteRepository, lp, lp, btcWallet, feeCollectorAddress)
	request := pegout.NewQuoteRequest(
		"mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe",
		entities.NewWei(1000000000000000000),
		"0x79568c2989232dCa1840087D73d403602364c0D4",
		"mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe",
	)
	result, err := useCase.Run(context.Background(), request)
	rsk.AssertExpectations(t)
	feeCollector.AssertExpectations(t)
	bridge.AssertExpectations(t)
	lbc.AssertExpectations(t)
	pegoutQuoteRepository.AssertExpectations(t)
	lp.AssertExpectations(t)
	btcWallet.AssertExpectations(t)
	assert.NotEmpty(t, result.Hash)
	require.NoError(t, entities.ValidateStruct(result.PegoutQuote))
	require.NoError(t, err)
}

func TestGetQuoteUseCase_Run_ValidateRequest(t *testing.T) {
	rsk := new(test.RskRpcMock)
	lp := new(test.ProviderMock)
	feeCollector := new(test.FeeCollectorMock)
	bridge := new(test.BridgeMock)
	lbc := new(test.LbcMock)
	pegoutQuoteRepository := new(test.PegoutQuoteRepositoryMock)
	btcWallet := new(test.BtcWalletMock)
	feeCollectorAddress := "feeCollectorAddress"
	useCase := pegout.NewGetQuoteUseCase(rsk, feeCollector, bridge, lbc, pegoutQuoteRepository, lp, lp, btcWallet, feeCollectorAddress)
	cases := test.Table[pegout.QuoteRequest, error]{
		{
			Value:  pegout.NewQuoteRequest("any address", entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe"),
			Result: usecases.BtcAddressNotSupportedError,
		},
		{
			Value:  pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", "any"),
			Result: usecases.BtcAddressNotSupportedError,
		},
		{
			Value:  pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", nil, "anything", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe"),
			Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value:  pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D41", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe"),
			Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value:  pegout.NewQuoteRequest("bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq", entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe"),
			Result: usecases.BtcAddressNotSupportedError,
		},
		{
			Value:  pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", "bc1qar0srrr7xfkvy5l643lydnw9re59gtzzwf5mdq"),
			Result: usecases.BtcAddressNotSupportedError,
		},
		{
			Value:  pegout.NewQuoteRequest("tb1qw508d6qejxtdg4y5r3zarvary0c5xw7kxpjzsx", entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe"),
			Result: usecases.BtcAddressNotSupportedError,
		},
		{
			Value:  pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(1), "0x79568c2989232dCa1840087D73d403602364c0D4", "tb1qw508d6qejxtdg4y5r3zarvary0c5xw7kxpjzsx"),
			Result: usecases.BtcAddressNotSupportedError,
		},
	}
	for _, testCase := range cases {
		result, err := useCase.Run(context.Background(), testCase.Value)
		assert.Equal(t, pegout.GetPegoutQuoteResult{}, result)
		require.Error(t, err)
		require.ErrorIs(t, err, testCase.Result)
	}
}

func TestGetQuoteUseCase_Run_ErrorHandling(t *testing.T) {
	cases := getQuoteUseCaseUnexpectedErrorSetups()

	request := pegout.NewQuoteRequest(
		"mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe",
		entities.NewWei(1000000000000000000),
		"0x79568c2989232dCa1840087D73d403602364c0D4",
		"mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe",
	)
	feeCollectorAddress := "feeCollectorAddress"
	for _, testCase := range cases {
		rsk := new(test.RskRpcMock)
		lp := new(test.ProviderMock)
		feeCollector := new(test.FeeCollectorMock)
		bridge := new(test.BridgeMock)
		lbc := new(test.LbcMock)
		pegoutQuoteRepository := new(test.PegoutQuoteRepositoryMock)
		btcWallet := new(test.BtcWalletMock)
		testCase.Value(rsk, feeCollector, bridge, lbc, lp, btcWallet, pegoutQuoteRepository)
		lp.On("GetRootstockConfirmationsForValue", mock.Anything).Return(uint16(10))
		lp.On("GetBitcoinConfirmationsForValue", mock.Anything, mock.Anything).Return(uint16(10))
		lp.On("CallFeePegout").Return(entities.NewWei(200))
		lp.On("PenaltyFeePegout").Return(entities.NewWei(20))
		lp.On("RskAddress").Return("0x1234")
		lp.On("BtcAddress").Return("address")
		lp.On("TimeForDepositPegout").Return(uint32(60000))
		lp.On("ExpireBlocksPegout").Return(uint64(60000))
		lbc.On("GetAddress").Return("0x1234")
		useCase := pegout.NewGetQuoteUseCase(rsk, feeCollector, bridge, lbc, pegoutQuoteRepository, lp, lp, btcWallet, feeCollectorAddress)
		result, err := useCase.Run(context.Background(), request)
		assert.Equal(t, pegout.GetPegoutQuoteResult{}, result)
		require.Error(t, err)
	}
}

// nolint:funlen
func getQuoteUseCaseUnexpectedErrorSetups() test.Table[func(
	rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
	lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock,
	pegoutQuoteRepository *test.PegoutQuoteRepositoryMock,
), error] {
	return test.Table[func(
		rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
		lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock,
		pegoutQuoteRepository *test.PegoutQuoteRepositoryMock,
	), error]{
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(assert.AnError)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(0), assert.AnError)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(0), assert.AnError)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(0), errors.New("Insufficient funds"))
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(0), assert.AnError)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), assert.AnError)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(0), assert.AnError)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000000), nil)
				lbc.On("GetAddress").Return("0x1234")
				lbc.On("HashPegoutQuote", mock.Anything).Return("0x9876543210", nil)
				pegoutQuoteRepository.On("InsertQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything).Return(assert.AnError)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000000), nil)
				lbc.On("GetAddress").Return("0x1234")
				lbc.On("HashPegoutQuote", mock.Anything).Return("", assert.AnError)
				pegoutQuoteRepository.On("InsertQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything).Return(nil)
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				lp.On("GetRootstockConfirmationsForValue", mock.Anything).Return(uint16(10))
				lp.On("GetBitcoinConfirmationsForValue", mock.Anything, mock.Anything).Return(uint16(10))
				lp.On("CallFeePegout").Return(entities.NewWei(200))
				lp.On("PenaltyFeePegout").Return(entities.NewWei(20))
				lp.On("RskAddress").Return("0x1234")
				lp.On("BtcAddress").Return("address")
				lp.On("TimeForDepositPegout").Return(uint32(60000))
				lp.On("ExpireBlocksPegout").Return(uint64(60000))
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(10), nil)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000000), nil)
				lbc.On("GetAddress").Return("0x1234")
				lbc.On("HashPegoutQuote", mock.Anything).Return("0x2134", nil)
				pegoutQuoteRepository.On("InsertQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything).Return(nil)
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				lp.On("GetRootstockConfirmationsForValue", mock.Anything).Return(uint16(0))
				lp.On("GetBitcoinConfirmationsForValue", mock.Anything, mock.Anything).Return(uint16(0))
				lp.On("CallFeePegout").Return(entities.NewWei(0))
				lp.On("PenaltyFeePegout").Return(entities.NewWei(0))
				lp.On("RskAddress").Return("")
				lp.On("BtcAddress").Return("")
				lp.On("TimeForDepositPegout").Return(uint32(0))
				lp.On("ExpireBlocksPegout").Return(uint64(0))
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(10), nil)
			},
		},
		{
			Value: func(rsk *test.RskRpcMock, feeCollector *test.FeeCollectorMock, bridge *test.BridgeMock,
				lbc *test.LbcMock, lp *test.ProviderMock, btcWallet *test.BtcWalletMock, pegoutQuoteRepository *test.PegoutQuoteRepositoryMock) {
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
				rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(entities.NewWei(0), assert.AnError)
				feeCollector.On("DaoFeePercentage").Return(uint64(12), nil)
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000000), nil)
				lbc.On("GetAddress").Return("0x1234")
				lbc.On("HashPegoutQuote", mock.Anything).Return("0x4321", nil)
				pegoutQuoteRepository.On("InsertQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything).Return(nil)
				lp.On("ValidateAmountForPegout", mock.Anything).Return(nil)
				lp.On("GetRootstockConfirmationsForValue", mock.Anything).Return(uint16(10))
				lp.On("GetBitcoinConfirmationsForValue", mock.Anything, mock.Anything).Return(uint16(10))
				lp.On("CallFeePegout").Return(entities.NewWei(200))
				lp.On("PenaltyFeePegout").Return(entities.NewWei(20))
				lp.On("RskAddress").Return("0x1234")
				lp.On("BtcAddress").Return("address")
				lp.On("TimeForDepositPegout").Return(uint32(60000))
				lp.On("ExpireBlocksPegout").Return(uint64(60000))
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(10), nil)
			},
		},
	}
}
