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

// nolint:funlen
func TestGetQuoteUseCase_Run(t *testing.T) {
	const (
		toAddress        = "mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe"
		rskRefundAddress = "0x79568c2989232dCa1840087D73d403602364c0D4"
		lbcAddress       = "0x1234"
		lpRskAddress     = "0x12ab"
		lpBtcAddress     = "address"
	)
	rsk := new(mocks.RootstockRpcServerMock)
	rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50000000), nil)
	rsk.On("GetHeight", test.AnyCtx).Return(uint64(100), nil)
	feeCollector := new(mocks.FeeCollectorMock)
	feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
	bridge := new(mocks.BridgeMock)
	lbc := new(mocks.LbcMock)
	lbc.On("GetAddress").Return(lbcAddress)
	lbc.On("HashPegoutQuote", mock.Anything).Return("0x9876543210", nil)
	pegoutQuoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	pegoutQuoteRepository.On("InsertQuote", test.AnyCtx, mock.Anything, mock.Anything).Return(nil)
	lp := new(mocks.ProviderMock)
	lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
	lp.On("GeneralConfiguration", test.AnyCtx).Return(getGeneralConfiguration())
	lp.On("RskAddress").Return(lpRskAddress)
	lp.On("BtcAddress").Return(lpBtcAddress)
	btcWallet := new(mocks.BtcWalletMock)
	btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000000000000000), nil)
	btc := new(mocks.BtcRpcMock)
	btc.On("ValidateAddress", mock.Anything).Return(nil)
	feeCollectorAddress := "feeCollectorAddress"
	contracts := blockchain.RskContracts{Lbc: lbc, FeeCollector: feeCollector, Bridge: bridge}
	rpc := blockchain.Rpc{Btc: btc, Rsk: rsk}
	useCase := pegout.NewGetQuoteUseCase(rpc, contracts, pegoutQuoteRepository, lp, lp, btcWallet, feeCollectorAddress)
	request := pegout.NewQuoteRequest(toAddress, entities.NewWei(1000000000000000000), rskRefundAddress)
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
	assert.Equal(t, toAddress, result.PegoutQuote.DepositAddress)
	assert.Equal(t, toAddress, result.PegoutQuote.BtcRefundAddress)
	assert.Equal(t, entities.NewWei(1000000000000000000), result.PegoutQuote.Value)
	assert.Equal(t, entities.NewWei(200), result.PegoutQuote.CallFee)
	assert.Equal(t, uint64(20), result.PegoutQuote.PenaltyFee)
	assert.Equal(t, "0x1234", result.PegoutQuote.LbcAddress)
	assert.NotEmpty(t, result.PegoutQuote.Nonce)
	assert.NotEmpty(t, result.PegoutQuote.AgreementTimestamp)
	assert.Zero(t, result.PegoutQuote.ProductFeeAmount)
	assert.Equal(t, uint16(10), result.PegoutQuote.DepositConfirmations)
	assert.Equal(t, uint16(10), result.PegoutQuote.TransferConfirmations)
	assert.Equal(t, uint32(60000), result.PegoutQuote.TransferTime)
	assert.Equal(t, 60000+result.PegoutQuote.AgreementTimestamp, result.PegoutQuote.DepositDateLimit)
	assert.Equal(t, 600+result.PegoutQuote.AgreementTimestamp, result.PegoutQuote.ExpireDate)
	assert.Equal(t, uint32(70100), result.PegoutQuote.ExpireBlock)
	assert.Equal(t, entities.NewWei(1000000000000000), result.PegoutQuote.GasFee)
	assert.Equal(t, rskRefundAddress, result.PegoutQuote.RskRefundAddress)
	assert.Equal(t, lbcAddress, result.PegoutQuote.LbcAddress)
	assert.Equal(t, lpBtcAddress, result.PegoutQuote.LpBtcAddress)
	assert.Equal(t, lpRskAddress, result.PegoutQuote.LpRskAddress)
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
				return pegout.NewQuoteRequest(wrongAddress, entities.NewWei(100000000000000000), "0x79568c2989232dCa1840087D73d403602364c0D4")
			},
			Result: blockchain.BtcAddressInvalidNetworkError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", wrongAddress).Return(blockchain.BtcAddressNotSupportedError).Once()
				btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest(wrongAddress, entities.NewWei(100000000000000000), "0x79568c2989232dCa1840087D73d403602364c0D4")
			}, Result: blockchain.BtcAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", nil, "anything")
			}, Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil)
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(100000000000000000), "0x79568c2989232dCa1840087D73d403602364c0D41")
			}, Result: usecases.RskAddressNotSupportedError,
		},
		{
			Value: func(btc *mocks.BtcRpcMock, lp *mocks.ProviderMock) pegout.QuoteRequest {
				btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
				btc.On("ValidateAddress", mock.Anything).Return(nil).Once()
				lp.On("PegoutConfiguration", test.AnyCtx).Return(getPegoutConfiguration())
				return pegout.NewQuoteRequest("mvL2bVzGUeC9oqVyQWJ4PxQspFzKgjzAqe", entities.NewWei(99999999999999999), "0x79568c2989232dCa1840087D73d403602364c0D4")
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
				rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(0), assert.AnError)
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
				rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", test.AnyCtx).Return(uint64(0), assert.AnError)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
				rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", test.AnyCtx).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), assert.AnError)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
				rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", test.AnyCtx).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				lbc.On("GetAddress").Return("0x1234")
				lbc.On("HashPegoutQuote", mock.Anything).Return("0x9876543210", nil)
				pegoutQuoteRepository.On("InsertQuote", test.AnyCtx, mock.Anything, mock.Anything).Return(assert.AnError)
				btcWallet.On("EstimateTxFees", mock.Anything, mock.Anything).Return(entities.NewWei(1000), nil)
			},
		},
		{
			Value: func(rsk *mocks.RootstockRpcServerMock, feeCollector *mocks.FeeCollectorMock, bridge *mocks.BridgeMock,
				lbc *mocks.LbcMock, lp *mocks.ProviderMock, btcWallet *mocks.BtcWalletMock, pegoutQuoteRepository *mocks.PegoutQuoteRepositoryMock) {
				rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", test.AnyCtx).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				lbc.On("GetAddress").Return("0x1234")
				lbc.On("HashPegoutQuote", mock.Anything).Return("", assert.AnError)
				pegoutQuoteRepository.On("InsertQuote", test.AnyCtx, mock.Anything, mock.Anything).Return(nil)
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
				rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", test.AnyCtx).Return(uint64(100), nil)
				feeCollector.On("DaoFeePercentage").Return(uint64(0), nil)
				lbc.On("GetAddress").Return("0x1234")
				lbc.On("HashPegoutQuote", mock.Anything).Return("0x2134", nil)
				pegoutQuoteRepository.On("InsertQuote", test.AnyCtx, mock.Anything, mock.Anything).Return(nil)
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
				rsk.On("GasPrice", test.AnyCtx).Return(entities.NewWei(50000000), nil)
				rsk.On("GetHeight", test.AnyCtx).Return(uint64(100), nil)
				rsk.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(entities.NewWei(0), assert.AnError)
				feeCollector.On("DaoFeePercentage").Return(uint64(12), nil)
				lbc.On("GetAddress").Return("0x1234")
				lbc.On("HashPegoutQuote", mock.Anything).Return("0x4321", nil)
				pegoutQuoteRepository.On("InsertQuote", test.AnyCtx, mock.Anything, mock.Anything).Return(nil)
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
		ExpireTime:     600,
		PenaltyFee:     entities.NewWei(20),
		CallFee:        entities.NewWei(200),
		MaxValue:       entities.NewUWei(10000000000000000000),
		MinValue:       entities.NewWei(100000000000000000),
		ExpireBlocks:   70000,
	}

}

func getGeneralConfiguration() lpEntity.GeneralConfiguration {
	return lpEntity.GeneralConfiguration{RskConfirmations: map[int]uint16{1: 10}, BtcConfirmations: map[int]uint16{1: 10}}
}
