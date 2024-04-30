package pegout_test

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	pegout "github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetQuoteUseCase_Run(t *testing.T) {
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
	rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
	feeCollector := new(mocks.FeeCollectorMock)
	feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
	bridge := new(mocks.BridgeMock)
	bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(1000000), nil)
	lbc := new(mocks.LbcMock)
	lbc.On("GetAddress").Return("0x1234")
	lbc.On("HashPegoutQuote", mock.Anything).Return("0x9876543210", nil)
	pegoutQuoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	pegoutQuoteRepository.On("InsertQuote", mock.AnythingOfType("context.backgroundCtx"), mock.Anything, mock.Anything).Return(nil)
	lp := new(mocks.ProviderMock)
	lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
	lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration())
	lp.On("RskAddress").Return("0x1234")
	lp.On("BtcAddress").Return("address")
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000000000000000), nil)
	btc := new(mocks.BtcRpcMock)
	btc.On("ValidateAddress", mock.Anything).Return(nil)
	feeCollectorAddress := "feeCollectorAddress"
	contracts := blockchain.RskContracts{Lbc: lbc, FeeCollector: feeCollector, Bridge: bridge}
	rpc := blockchain.Rpc{Btc: btc, Rsk: rsk}
	useCase := pegout.NewGetQuoteUseCase(rpc, contracts, pegoutQuoteRepository, lp, lp, btcWallet, feeCollectorAddress)
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
	rsk := new(mocks.RootstockRpcServerMock)
	feeCollector := new(mocks.FeeCollectorMock)
	bridge := new(mocks.BridgeMock)
	lbc := new(mocks.LbcMock)
	pegoutQuoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	btcWallet := new(mocks.BtcWalletMock)
	feeCollectorAddress := "feeCollectorAddress"
	cases := getQuoteUseCaseErrorSetups()
	for _, testCase := range cases {
		btc := new(mocks.BtcRpcMock)
		lp := new(mocks.ProviderMock)
		contracts := blockchain.RskContracts{Lbc: lbc, FeeCollector: feeCollector, Bridge: bridge}
		rpc := blockchain.Rpc{Btc: btc, Rsk: rsk}
		useCase := pegout.NewGetQuoteUseCase(rpc, contracts, pegoutQuoteRepository, lp, lp, btcWallet, feeCollectorAddress)
		result, err := useCase.Run(context.Background(), testCase.Value(btc, lp))
		assert.Equal(t, pegout.GetPegoutQuoteResult{}, result)
		require.Error(t, err)
		require.ErrorIs(t, err, testCase.Result)
	}
}

func getQuoteUseCaseErrorSetups() test.Table[func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest, error] {
	const wrongAddress = "wrong address"
	return test.Table[func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest, error]{
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", wrongAddress).Return(blockchain.BtcAddressInvalidNetworkError).Once()
				btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest(wrongAddress, entities.NewWei(100000000000000000), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			},
			Result: blockchain.BtcAddressInvalidNetworkError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", wrongAddress).Return(blockchain.BtcAddressNotSupportedError).Once()
				btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest(wrongAddress, entities.NewWei(100000000000000000), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			}, Result: blockchain.BtcAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", nil, "anything", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			}, Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(100000000000000000), "0x79568c2989232dCa1840087D73d403602364c0D41", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			}, Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", wrongAddress).Return(blockchain.BtcAddressInvalidNetworkError).Once()
				btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(100000000000000000), "0x79568c2989232dCa1840087D73d403602364c0D4", wrongAddress)
			}, Result: blockchain.BtcAddressInvalidNetworkError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", wrongAddress).Return(blockchain.BtcAddressNotSupportedError).Once()
				btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(100000000000000000), "0x79568c2989232dCa1840087D73d403602364c0D4", wrongAddress)
			}, Result: blockchain.BtcAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
				btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(99999999999999999), "0x79568c2989232dCa1840087D73d403602364c0D4", "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe")
			}, Result: lpEntity.AmountOutOfRangeError,
		},
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
		rsk := new(mocks.RootstockRpcServerMock)
		lp := new(mocks.ProviderMock)
		feeCollector := new(mocks.FeeCollectorMock)
		bridge := new(mocks.BridgeMock)
		lbc := new(mocks.LbcMock)
		pegoutQuoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		btcWallet := new(mocks.BtcWalletMock)
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
		lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration())
		lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
		btc := new(mocks.BtcRpcMock)
		btc.On("ValidateAddress", mock.Anything).Return(nil)
		contracts := blockchain.RskContracts{Lbc: lbc, FeeCollector: feeCollector, Bridge: bridge}
		rpc := blockchain.Rpc{Btc: btc, Rsk: rsk}
		useCase := pegout.NewGetQuoteUseCase(rpc, contracts, pegoutQuoteRepository, lp, lp, btcWallet, feeCollectorAddress)
		result, err := useCase.Run(context.Background(), request)
		assert.Equal(t, pegout.GetPegoutQuoteResult{}, result)
		require.Error(t, err)
	}
}

// nolint:funlen
func getQuoteUseCaseUnexpectedErrorSetups() test.Table[func(
	rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
	lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock,
	pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock,
), error] {
	return test.Table[func(
		rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
		lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock,
		pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock,
	), error]{
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(0), assert.AnError)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			},
		},
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(0), assert.AnError)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			},
		},
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(0), errors.New("Insufficient funds"))
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
			},
		},
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(0), assert.AnError)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), assert.AnError)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
				rsk.On("GasPrice", mock.AnythingOfType("context.backgroundCtx")).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", mock.AnythingOfType("context.backgroundCtx")).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				bridge.On("GetMinimumLockTxValue").Return(entities.NewWei(0), assert.AnError)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
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
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
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
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
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
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
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

func getPegoutConfiguration() lpEntity.PegoutConfiguration {
	return lpEntity.PegoutConfiguration{
		TimeForDeposit: 60000,
		CallTime:       600,
		PenaltyFee:     entities.NewWei(20),
		CallFee:        entities.NewWei(200),
		MaxValue:       entities.NewUWei(10000000000000000000),
		MinValue:       entities.NewWei(100000000000000000),
		ExpireBlocks:   60000,
	}

}

func getGeneralConfiguration() lpEntity.GeneralConfiguration {
	return lpEntity.GeneralConfiguration{RskConfirmations: map[int]uint16{1: 10}, BtcConfirmations: map[int]uint16{1: 10}}
}
