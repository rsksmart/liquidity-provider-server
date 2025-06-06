package reports

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"math/big"
)

type GetAssetsReportResponse struct {
	BtcBalance    *big.Int `json:"btcBalance" validate:"required"`
	RbtcBalance   *big.Int `json:"rbtcBalance" validate:"required"`
	BtcLocked     *big.Int `json:"btcLocked" validate:"required"`
	RbtcLocked    *big.Int `json:"rbtcLocked" validate:"required"`
	BtcLiquidity  *big.Int `json:"btcLiquidity" validate:"required"`
	RbtcLiquidity *big.Int `json:"rbtcLiquidity" validate:"required"`
}

type GetAssetsReportUseCase struct {
	btcWallet        blockchain.BitcoinWallet
	rsk              blockchain.Rpc
	lp               liquidity_provider.LiquidityProvider
	peginProvider    liquidity_provider.PeginLiquidityProvider
	pegoutProvider   liquidity_provider.PegoutLiquidityProvider
	peginRepository  quote.PeginQuoteRepository
	pegoutRepository quote.PegoutQuoteRepository
}

func NewGetAssetsReportUseCase(
	wallet blockchain.BitcoinWallet,
	rsk blockchain.Rpc,
	lp liquidity_provider.LiquidityProvider,
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	peginRepository quote.PeginQuoteRepository,
	pegoutRepository quote.PegoutQuoteRepository,
) *GetAssetsReportUseCase {
	return &GetAssetsReportUseCase{
		btcWallet:        wallet,
		rsk:              rsk,
		lp:               lp,
		peginProvider:    peginProvider,
		pegoutProvider:   pegoutProvider,
		peginRepository:  peginRepository,
		pegoutRepository: pegoutRepository,
	}
}

func (useCase *GetAssetsReportUseCase) Run(ctx context.Context) (GetAssetsReportResponse, error) {
	response := GetAssetsReportResponse{
		BtcBalance:    entities.NewWei(0).AsBigInt(),
		RbtcBalance:   entities.NewWei(0).AsBigInt(),
		BtcLocked:     entities.NewWei(0).AsBigInt(),
		RbtcLocked:    entities.NewWei(0).AsBigInt(),
		BtcLiquidity:  entities.NewWei(0).AsBigInt(),
		RbtcLiquidity: entities.NewWei(0).AsBigInt(),
	}
	btcBalance, err := useCase.GetBtcBalance()
	if err != nil {
		return response, err
	}

	rbtcBalance, err := useCase.GetRBTCBalance(ctx)
	if err != nil {
		return response, err
	}

	rbtcLocked, err := useCase.GetRBTCLocked(ctx)
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

	rbtcLiquidity, err := useCase.GetRBTCLiquidity(ctx)
	if err != nil {
		return response, err
	}

	response.BtcBalance = btcBalance.AsBigInt()
	response.RbtcLiquidity = rbtcLiquidity.AsBigInt()
	response.BtcLiquidity = btcLiquidity.AsBigInt()
	response.BtcLocked = lockedBtc.AsBigInt()
	response.RbtcLocked = rbtcLocked.AsBigInt()
	response.RbtcBalance = rbtcBalance.AsBigInt()

	return response, nil
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
