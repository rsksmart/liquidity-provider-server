package reports_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

var retainedPeginQuotes = []quote.RetainedPeginQuote{
	{
		QuoteHash:         "mockPeginQuoteHash1",
		RequiredLiquidity: entities.NewWei(1500),
		State:             quote.PeginStateWaitingForDeposit,
	},
	{
		QuoteHash:         "mockPeginQuoteHash2",
		RequiredLiquidity: entities.NewWei(2500),
		State:             quote.PeginStateCallForUserSucceeded,
	},
	{
		QuoteHash:         "mockPeginQuoteHash3",
		RequiredLiquidity: entities.NewWei(3500),
		State:             quote.PeginStateWaitingForDeposit,
	},
	{
		QuoteHash:         "mockPeginQuoteHash4",
		RequiredLiquidity: entities.NewWei(4500),
		State:             quote.PeginStateCallForUserSucceeded,
	},
	{
		QuoteHash:         "mockPeginQuoteHash5",
		RequiredLiquidity: entities.NewWei(5500),
		State:             quote.PeginStateWaitingForDeposit,
	},
}

var retainedPegoutQuotes = []quote.RetainedPegoutQuote{
	{
		QuoteHash:         "mockQuoteHash1",
		RequiredLiquidity: entities.NewWei(1000),
		State:             quote.PegoutStateWaitingForDeposit,
	},
	{
		QuoteHash:         "mockQuoteHash2",
		RequiredLiquidity: entities.NewWei(2000),
		State:             quote.PegoutStateWaitingForDepositConfirmations,
	},
	{
		QuoteHash:         "mockQuoteHash3",
		RequiredLiquidity: entities.NewWei(3000),
		State:             quote.PegoutStateWaitingForDeposit,
	},
	{
		QuoteHash:         "mockQuoteHash4",
		RequiredLiquidity: entities.NewWei(4000),
		State:             quote.PegoutStateWaitingForDepositConfirmations,
	},
	{
		QuoteHash:         "mockQuoteHash5",
		RequiredLiquidity: entities.NewWei(5000),
		State:             quote.PegoutStateWaitingForDeposit,
	},
}

// nolint:funlen
func TestGetAssetsReportUseCase_Run(t *testing.T) {
	ctx := context.Background()

	rskAddress := "rskAddress"

	wallet := mocks.NewBitcoinWalletMock(t)
	wallet.On("GetBalance").Return(entities.NewWei(100000), nil)

	lpMock := &mocks.ProviderMock{}
	lpMock.On("RskAddress").Return(rskAddress)

	rsk := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{Rsk: rsk}

	rsk.On("GetBalance", ctx, rskAddress).Return(entities.NewWei(100000), nil)

	lp := new(mocks.ProviderMock)
	lp.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(67500), nil)
	lp.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(85000), nil)

	peginRepository := mocks.NewPeginQuoteRepositoryMock(t)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return(retainedPeginQuotes, nil)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateCallForUserSucceeded).Return(retainedPeginQuotes, nil)
	pegoutRepository := mocks.NewPegoutQuoteRepositoryMock(t)
	pegoutRepository.On("GetRetainedQuoteByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).Return(retainedPegoutQuotes, nil)

	lbc := &mocks.LbcMock{}
	lbc.On("GetBalance", rskAddress).Return(entities.NewWei(100000), nil)

	useCase := reports.NewGetAssetsReportUseCase(wallet, rpc, lpMock, lp, lp, peginRepository, pegoutRepository, blockchain.RskContracts{Lbc: lbc})

	result, err := useCase.Run(ctx)

	require.NoError(t, err)
	require.Equal(t, result.RbtcLiquidity, big.NewInt(67500))
	require.Equal(t, result.BtcLiquidity, big.NewInt(85000))
	require.Equal(t, result.RbtcLockedLbc, big.NewInt(100000))
	require.Equal(t, result.RbtcLockedForUsers, big.NewInt(17500))
	require.Equal(t, result.RbtcWaitingRefund, big.NewInt(17500))
	require.Equal(t, result.RbtcLiquidity, big.NewInt(67500))
	require.Equal(t, result.RbtcWalletBalance, big.NewInt(100000))
	require.Equal(t, result.BtcLockedForUsers, big.NewInt(15000))
	require.Equal(t, result.BtcLiquidity, big.NewInt(85000))
	require.Equal(t, result.BtcWalletBalance, big.NewInt(100000))
	require.Equal(t, result.BtcRebalancing, big.NewInt(0))
	lp.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	lpMock.AssertExpectations(t)
	wallet.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

func TestGetAssetsReportUseCase_Run_GetRbtcLockedLbcError(t *testing.T) {
	ctx := context.Background()

	rskAddress := "rskAddress"

	wallet := mocks.NewBitcoinWalletMock(t)

	lpMock := &mocks.ProviderMock{}
	lpMock.On("RskAddress").Return(rskAddress)

	rsk := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{Rsk: rsk}

	lp := new(mocks.ProviderMock)

	peginRepository := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepository := mocks.NewPegoutQuoteRepositoryMock(t)

	lbc := &mocks.LbcMock{}
	lbc.On("GetBalance", rskAddress).Return(nil, assert.AnError)

	useCase := reports.NewGetAssetsReportUseCase(wallet, rpc, lpMock, lp, lp, peginRepository, pegoutRepository, blockchain.RskContracts{Lbc: lbc})

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcLockedLbc, big.NewInt(0))
	require.Equal(t, result.RbtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.RbtcWaitingRefund, big.NewInt(0))
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcRebalancing, big.NewInt(0))
	lp.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	lpMock.AssertExpectations(t)
	wallet.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

func TestGetAssetsReportUseCase_Run_GetRBTCLockedError(t *testing.T) {
	ctx := context.Background()

	rskAddress := "rskAddress"

	wallet := mocks.NewBitcoinWalletMock(t)

	lpMock := &mocks.ProviderMock{}
	lpMock.On("RskAddress").Return(rskAddress)

	rsk := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{Rsk: rsk}

	lp := new(mocks.ProviderMock)

	peginRepository := mocks.NewPeginQuoteRepositoryMock(t)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return(nil, assert.AnError)
	pegoutRepository := mocks.NewPegoutQuoteRepositoryMock(t)

	lbc := &mocks.LbcMock{}
	lbc.On("GetBalance", rskAddress).Return(entities.NewWei(100000), nil)

	useCase := reports.NewGetAssetsReportUseCase(wallet, rpc, lpMock, lp, lp, peginRepository, pegoutRepository, blockchain.RskContracts{Lbc: lbc})

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcLockedLbc, big.NewInt(0))
	require.Equal(t, result.RbtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.RbtcWaitingRefund, big.NewInt(0))
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcRebalancing, big.NewInt(0))
	lp.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	lpMock.AssertExpectations(t)
	wallet.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

func TestGetAssetsReportUseCase_Run_GetRBTCWaitingForRefundError(t *testing.T) {
	ctx := context.Background()

	rskAddress := "rskAddress"

	wallet := mocks.NewBitcoinWalletMock(t)

	lpMock := &mocks.ProviderMock{}
	lpMock.On("RskAddress").Return(rskAddress)

	rsk := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{Rsk: rsk}

	lp := new(mocks.ProviderMock)

	peginRepository := mocks.NewPeginQuoteRepositoryMock(t)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return(retainedPeginQuotes, nil)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateCallForUserSucceeded).Return(nil, assert.AnError)
	pegoutRepository := mocks.NewPegoutQuoteRepositoryMock(t)

	lbc := &mocks.LbcMock{}
	lbc.On("GetBalance", rskAddress).Return(entities.NewWei(100000), nil)

	useCase := reports.NewGetAssetsReportUseCase(wallet, rpc, lpMock, lp, lp, peginRepository, pegoutRepository, blockchain.RskContracts{Lbc: lbc})

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcLockedLbc, big.NewInt(0))
	require.Equal(t, result.RbtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.RbtcWaitingRefund, big.NewInt(0))
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcRebalancing, big.NewInt(0))
	lp.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	lpMock.AssertExpectations(t)
	wallet.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

func TestGetAssetsReportUseCase_Run_GetRBTCLiquidityError(t *testing.T) {
	ctx := context.Background()

	rskAddress := "rskAddress"

	wallet := mocks.NewBitcoinWalletMock(t)

	lpMock := &mocks.ProviderMock{}
	lpMock.On("RskAddress").Return(rskAddress)

	rsk := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{Rsk: rsk}

	lp := new(mocks.ProviderMock)
	lp.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(0), assert.AnError)

	peginRepository := mocks.NewPeginQuoteRepositoryMock(t)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return(retainedPeginQuotes, nil)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateCallForUserSucceeded).Return(retainedPeginQuotes, nil)
	pegoutRepository := mocks.NewPegoutQuoteRepositoryMock(t)

	lbc := &mocks.LbcMock{}
	lbc.On("GetBalance", rskAddress).Return(entities.NewWei(100000), nil)

	useCase := reports.NewGetAssetsReportUseCase(wallet, rpc, lpMock, lp, lp, peginRepository, pegoutRepository, blockchain.RskContracts{Lbc: lbc})

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcLockedLbc, big.NewInt(0))
	require.Equal(t, result.RbtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.RbtcWaitingRefund, big.NewInt(0))
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcRebalancing, big.NewInt(0))
	lp.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	lpMock.AssertExpectations(t)
	wallet.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

func TestGetAssetsReportUseCase_Run_GetRBTCBalanceError(t *testing.T) {
	ctx := context.Background()

	rskAddress := "rskAddress"

	wallet := mocks.NewBitcoinWalletMock(t)

	lpMock := &mocks.ProviderMock{}
	lpMock.On("RskAddress").Return(rskAddress)

	rsk := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{Rsk: rsk}

	rsk.On("GetBalance", ctx, rskAddress).Return(nil, assert.AnError)

	lp := new(mocks.ProviderMock)
	lp.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(67500), nil)

	peginRepository := mocks.NewPeginQuoteRepositoryMock(t)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return(retainedPeginQuotes, nil)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateCallForUserSucceeded).Return(retainedPeginQuotes, nil)
	pegoutRepository := mocks.NewPegoutQuoteRepositoryMock(t)

	lbc := &mocks.LbcMock{}
	lbc.On("GetBalance", rskAddress).Return(entities.NewWei(100000), nil)

	useCase := reports.NewGetAssetsReportUseCase(wallet, rpc, lpMock, lp, lp, peginRepository, pegoutRepository, blockchain.RskContracts{Lbc: lbc})

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcLockedLbc, big.NewInt(0))
	require.Equal(t, result.RbtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.RbtcWaitingRefund, big.NewInt(0))
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcRebalancing, big.NewInt(0))
	lp.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	lpMock.AssertExpectations(t)
	wallet.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

func TestGetAssetsReportUseCase_Run_GetBTCLockedError(t *testing.T) {
	ctx := context.Background()

	rskAddress := "rskAddress"

	wallet := mocks.NewBitcoinWalletMock(t)

	lpMock := &mocks.ProviderMock{}
	lpMock.On("RskAddress").Return(rskAddress)

	rsk := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{Rsk: rsk}

	rsk.On("GetBalance", ctx, rskAddress).Return(entities.NewWei(100000), nil)

	lp := new(mocks.ProviderMock)
	lp.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(67500), nil)

	peginRepository := mocks.NewPeginQuoteRepositoryMock(t)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return(retainedPeginQuotes, nil)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateCallForUserSucceeded).Return(retainedPeginQuotes, nil)
	pegoutRepository := mocks.NewPegoutQuoteRepositoryMock(t)
	pegoutRepository.On("GetRetainedQuoteByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).Return(nil, assert.AnError)

	lbc := &mocks.LbcMock{}
	lbc.On("GetBalance", rskAddress).Return(entities.NewWei(100000), nil)

	useCase := reports.NewGetAssetsReportUseCase(wallet, rpc, lpMock, lp, lp, peginRepository, pegoutRepository, blockchain.RskContracts{Lbc: lbc})

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcLockedLbc, big.NewInt(0))
	require.Equal(t, result.RbtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.RbtcWaitingRefund, big.NewInt(0))
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcRebalancing, big.NewInt(0))
	lp.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	lpMock.AssertExpectations(t)
	wallet.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

func TestGetAssetsReportUseCase_Run_GetBTCLiquidityError(t *testing.T) {
	ctx := context.Background()

	rskAddress := "rskAddress"

	wallet := mocks.NewBitcoinWalletMock(t)

	lpMock := &mocks.ProviderMock{}
	lpMock.On("RskAddress").Return(rskAddress)

	rsk := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{Rsk: rsk}

	rsk.On("GetBalance", ctx, rskAddress).Return(entities.NewWei(100000), nil)

	lp := new(mocks.ProviderMock)
	lp.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(67500), nil)
	lp.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(0), assert.AnError)

	peginRepository := mocks.NewPeginQuoteRepositoryMock(t)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return(retainedPeginQuotes, nil)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateCallForUserSucceeded).Return(retainedPeginQuotes, nil)
	pegoutRepository := mocks.NewPegoutQuoteRepositoryMock(t)
	pegoutRepository.On("GetRetainedQuoteByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).Return(retainedPegoutQuotes, nil)

	lbc := &mocks.LbcMock{}
	lbc.On("GetBalance", rskAddress).Return(entities.NewWei(100000), nil)

	useCase := reports.NewGetAssetsReportUseCase(wallet, rpc, lpMock, lp, lp, peginRepository, pegoutRepository, blockchain.RskContracts{Lbc: lbc})

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcLockedLbc, big.NewInt(0))
	require.Equal(t, result.RbtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.RbtcWaitingRefund, big.NewInt(0))
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcRebalancing, big.NewInt(0))
	lp.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	lpMock.AssertExpectations(t)
	wallet.AssertExpectations(t)
	lbc.AssertExpectations(t)
}

func TestGetAssetsReportUseCase_Run_GetBtcBalanceError(t *testing.T) {
	ctx := context.Background()

	rskAddress := "rskAddress"

	wallet := mocks.NewBitcoinWalletMock(t)
	wallet.On("GetBalance").Return(nil, assert.AnError)

	lpMock := &mocks.ProviderMock{}
	lpMock.On("RskAddress").Return(rskAddress)

	rsk := new(mocks.RootstockRpcServerMock)
	rpc := blockchain.Rpc{Rsk: rsk}

	rsk.On("GetBalance", ctx, rskAddress).Return(entities.NewWei(100000), nil)

	lp := new(mocks.ProviderMock)
	lp.On("AvailablePeginLiquidity", ctx).Return(entities.NewWei(67500), nil)
	lp.On("AvailablePegoutLiquidity", ctx).Return(entities.NewWei(85000), nil)

	peginRepository := mocks.NewPeginQuoteRepositoryMock(t)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations).Return(retainedPeginQuotes, nil)
	peginRepository.On("GetRetainedQuoteByState", ctx, quote.PeginStateCallForUserSucceeded).Return(retainedPeginQuotes, nil)
	pegoutRepository := mocks.NewPegoutQuoteRepositoryMock(t)
	pegoutRepository.On("GetRetainedQuoteByState", ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations).Return(retainedPegoutQuotes, nil)

	lbc := &mocks.LbcMock{}
	lbc.On("GetBalance", rskAddress).Return(entities.NewWei(100000), nil)

	useCase := reports.NewGetAssetsReportUseCase(wallet, rpc, lpMock, lp, lp, peginRepository, pegoutRepository, blockchain.RskContracts{Lbc: lbc})

	result, err := useCase.Run(ctx)

	require.Error(t, err)
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcLockedLbc, big.NewInt(0))
	require.Equal(t, result.RbtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.RbtcWaitingRefund, big.NewInt(0))
	require.Equal(t, result.RbtcLiquidity, big.NewInt(0))
	require.Equal(t, result.RbtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcLockedForUsers, big.NewInt(0))
	require.Equal(t, result.BtcLiquidity, big.NewInt(0))
	require.Equal(t, result.BtcWalletBalance, big.NewInt(0))
	require.Equal(t, result.BtcRebalancing, big.NewInt(0))
	lp.AssertExpectations(t)
	peginRepository.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	rsk.AssertExpectations(t)
	lpMock.AssertExpectations(t)
	wallet.AssertExpectations(t)
	lbc.AssertExpectations(t)
}
