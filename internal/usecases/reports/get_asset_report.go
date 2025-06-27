package reports

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
)

type GetAssetReportResult struct {
	RbtcLockedLbc      *entities.Wei `json:"RbtcLockedLbc" validate:"required"`
	RbtcLockedForUsers *entities.Wei `json:"RbtcLockedForUsers" validate:"required"`
	RbtcWaitingRefund  *entities.Wei `json:"rbtcWaitingRefund" validate:"required"`
	RbtcLiquidity      *entities.Wei `json:"rbtcLiquidity" validate:"required"`
	RbtcWalletBalance  *entities.Wei `json:"rbtcWalletBalance" validate:"required"`
	BtcLockedForUsers  *entities.Wei `json:"btcLockedForUsers" validate:"required"`
	BtcLiquidity       *entities.Wei `json:"btcLiquidity" validate:"required"`
	BtcWalletBalance   *entities.Wei `json:"btcWalletBalance" validate:"required"`
	BtcRebalancing     *entities.Wei `json:"btcRebalancing" validate:"required"`
}

type GetAssetReportUseCase struct {
	btcWallet        blockchain.BitcoinWallet
	rsk              blockchain.Rpc
	lp               liquidity_provider.LiquidityProvider
	peginProvider    liquidity_provider.PeginLiquidityProvider
	pegoutProvider   liquidity_provider.PegoutLiquidityProvider
	peginRepository  quote.PeginQuoteRepository
	pegoutRepository quote.PegoutQuoteRepository
	contracts        blockchain.RskContracts
}

func NewGetAssetReportUseCase(
	wallet blockchain.BitcoinWallet,
	rsk blockchain.Rpc,
	lp liquidity_provider.LiquidityProvider,
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	peginRepository quote.PeginQuoteRepository,
	pegoutRepository quote.PegoutQuoteRepository,
	contracts blockchain.RskContracts,
) *GetAssetReportUseCase {
	return &GetAssetReportUseCase{
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

// Run returns dummy values for all financial metrics
func (useCase *GetAssetReportUseCase) Run(ctx context.Context) (GetAssetReportResult, error) {
	// Return dummy values directly - simulating a typical liquidity provider state
	response := GetAssetReportResult{
		RbtcLockedLbc:      entities.NewWei(1500000000000000000), // 1.5 RBTC
		RbtcLockedForUsers: entities.NewWei(2000000000000000000), // 2.0 RBTC
		RbtcWaitingRefund:  entities.NewWei(500000000000000000),  // 0.5 RBTC
		RbtcLiquidity:      entities.NewWei(5000000000000000000), // 5.0 RBTC
		RbtcWalletBalance:  entities.NewWei(3000000000000000000), // 3.0 RBTC
		BtcLockedForUsers:  entities.NewWei(1800000000000000000), // 1.8 BTC equivalent
		BtcLiquidity:       entities.NewWei(4500000000000000000), // 4.5 BTC equivalent
		BtcWalletBalance:   entities.NewWei(2800000000000000000), // 2.8 BTC equivalent
		BtcRebalancing:     entities.NewWei(100000000000000000),  // 0.1 BTC equivalent
	}

	return response, nil
}
