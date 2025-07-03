package reports

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"math/big"
)

type GetAssetsReportResult struct {
	RbtcLockedLbc      *big.Int `json:"RbtcLockedLbc" validate:"required"`
	RbtcLockedForUsers *big.Int `json:"RbtcLockedForUsers" validate:"required"`
	RbtcWaitingRefund  *big.Int `json:"rbtcWaitingRefund" validate:"required"`
	RbtcLiquidity      *big.Int `json:"rbtcLiquidity" validate:"required"`
	RbtcWalletBalance  *big.Int `json:"rbtcWalletBalance" validate:"required"`
	BtcLockedForUsers  *big.Int `json:"btcLockedForUsers" validate:"required"`
	BtcLiquidity       *big.Int `json:"btcLiquidity" validate:"required"`
	BtcWalletBalance   *big.Int `json:"btcWalletBalance" validate:"required"`
	BtcRebalancing     *big.Int `json:"btcRebalancing" validate:"required"`
}

type GetAssetsReportUseCase struct {
	btcWallet        blockchain.BitcoinWallet
	rsk              blockchain.Rpc
	lp               liquidity_provider.LiquidityProvider
	peginProvider    liquidity_provider.PeginLiquidityProvider
	pegoutProvider   liquidity_provider.PegoutLiquidityProvider
	peginRepository  quote.PeginQuoteRepository
	pegoutRepository quote.PegoutQuoteRepository
	contracts        blockchain.RskContracts
}

func NewGetAssetsReportUseCase(
	wallet blockchain.BitcoinWallet,
	rsk blockchain.Rpc,
	lp liquidity_provider.LiquidityProvider,
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	peginRepository quote.PeginQuoteRepository,
	pegoutRepository quote.PegoutQuoteRepository,
	contracts blockchain.RskContracts,
) *GetAssetsReportUseCase {
	return &GetAssetsReportUseCase{
		btcWallet:        wallet,
		rsk:              rsk,
		lp:               lp,
		peginProvider:    peginProvider,
		pegoutProvider:   pegoutProvider,
		peginRepository:  peginRepository,
		pegoutRepository: pegoutRepository,
		contracts:        contracts,
	}
}

func (useCase *GetAssetsReportUseCase) Run(ctx context.Context) (GetAssetsReportResult, error) {
	response := GetAssetsReportResult{
		RbtcLockedLbc:      entities.NewWei(0).AsBigInt(),
		RbtcLockedForUsers: entities.NewWei(0).AsBigInt(),
		RbtcWaitingRefund:  entities.NewWei(0).AsBigInt(),
		RbtcLiquidity:      entities.NewWei(0).AsBigInt(),
		RbtcWalletBalance:  entities.NewWei(0).AsBigInt(),
		BtcLockedForUsers:  entities.NewWei(0).AsBigInt(),
		BtcLiquidity:       entities.NewWei(0).AsBigInt(),
		BtcWalletBalance:   entities.NewWei(0).AsBigInt(),
		BtcRebalancing:     entities.NewWei(0).AsBigInt(),
	}

	rbtcLockedLbc, err := useCase.GetRbtcLockedLbc()
	if err != nil {
		return response, err
	}
	rbtcLocked, err := useCase.GetRBTCLocked(ctx)
	if err != nil {
		return response, err
	}
	rbtcWaitingRefund, err := useCase.GetRBTCWaitingForRefund(ctx)
	if err != nil {
		return response, err
	}
	rbtcLiquidity, err := useCase.GetRBTCLiquidity(ctx)
	if err != nil {
		return response, err
	}
	rbtcBalance, err := useCase.GetRBTCBalance(ctx)
	if err != nil {
		return response, err
	}
	lockedBtc, err := useCase.GetBTCLocked(ctx)
	if err != nil {
		return response, err
	}
	btcLiquidity, err := useCase.GetBTCLiquidity(ctx)
	if err != nil {
		return response, err
	}
	btcBalance, err := useCase.GetBtcBalance()
	if err != nil {
		return response, err
	}

	response.RbtcLockedLbc = rbtcLockedLbc.AsBigInt()
	response.RbtcLockedForUsers = rbtcLocked.AsBigInt()
	response.RbtcWaitingRefund = rbtcWaitingRefund.AsBigInt()
	response.RbtcLiquidity = rbtcLiquidity.AsBigInt()
	response.RbtcWalletBalance = rbtcBalance.AsBigInt()
	response.BtcLockedForUsers = lockedBtc.AsBigInt()
	response.BtcLiquidity = btcLiquidity.AsBigInt()
	response.BtcWalletBalance = btcBalance.AsBigInt()

	return response, nil
}

func (useCase *GetAssetsReportUseCase) GetRbtcLockedLbc() (*entities.Wei, error) {
	return useCase.contracts.Lbc.GetBalance(useCase.lp.RskAddress())
}

func (useCase *GetAssetsReportUseCase) GetRBTCLiquidity(ctx context.Context) (*entities.Wei, error) {
	rbtcLiquidity, err := useCase.peginProvider.AvailablePeginLiquidity(ctx)
	if err != nil {
		return nil, err
	}
	return rbtcLiquidity, nil
}

func (useCase *GetAssetsReportUseCase) GetBTCLiquidity(ctx context.Context) (*entities.Wei, error) {
	btcLiquidity, err := useCase.pegoutProvider.AvailablePegoutLiquidity(ctx)
	if err != nil {
		return nil, err
	}
	return btcLiquidity, nil
}

func (useCase *GetAssetsReportUseCase) GetBTCLocked(ctx context.Context) (*entities.Wei, error) {
	lockedPegout := entities.NewWei(0)
	quotes, err := useCase.pegoutRepository.GetRetainedQuoteByState(ctx,
		quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations,
	)
	if err != nil {
		return nil, err
	}
	for _, retainedQuote := range quotes {
		lockedPegout.Add(lockedPegout, retainedQuote.RequiredLiquidity)
	}

	return lockedPegout, nil
}

func (useCase *GetAssetsReportUseCase) GetRBTCLocked(ctx context.Context) (*entities.Wei, error) {
	lockedPegin := entities.NewWei(0)
	peginQuotes, err := useCase.peginRepository.GetRetainedQuoteByState(ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations)
	if err != nil {
		return nil, err
	}
	for _, retainedQuote := range peginQuotes {
		lockedPegin.Add(lockedPegin, retainedQuote.RequiredLiquidity)
	}

	return lockedPegin, nil
}

func (useCase *GetAssetsReportUseCase) GetRBTCWaitingForRefund(ctx context.Context) (*entities.Wei, error) {
	lockedPegin := entities.NewWei(0)
	peginQuotes, err := useCase.peginRepository.GetRetainedQuoteByState(ctx, quote.PeginStateCallForUserSucceeded)
	if err != nil {
		return nil, err
	}
	for _, retainedQuote := range peginQuotes {
		lockedPegin.Add(lockedPegin, retainedQuote.RequiredLiquidity)
	}

	return lockedPegin, nil
}

func (useCase *GetAssetsReportUseCase) GetRBTCBalance(ctx context.Context) (*entities.Wei, error) {
	lpsBalance, err := useCase.rsk.Rsk.GetBalance(ctx, useCase.lp.RskAddress())
	if err != nil {
		return nil, err
	}
	return lpsBalance, nil
}

func (useCase *GetAssetsReportUseCase) GetBtcBalance() (*entities.Wei, error) {
	btcBalance, err := useCase.btcWallet.GetBalance()
	if err != nil {
		return nil, err
	}
	return btcBalance, nil
}
