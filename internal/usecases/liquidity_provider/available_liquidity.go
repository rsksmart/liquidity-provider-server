package liquidity_provider

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type LiquidityStatusUseCase struct {
	contracts blockchain.RskContracts
	provider  liquidity_provider.LiquidityProvider
	peginProvider liquidity_provider.PeginLiquidityProvider
	rpc       blockchain.Rpc
    btcWallet blockchain.BitcoinWallet
}

func NewLiquidityStatusUseCase(contracts blockchain.RskContracts, provider liquidity_provider.LiquidityProvider, rpc blockchain.Rpc,btcWallet blockchain.BitcoinWallet, peginProvider liquidity_provider.PeginLiquidityProvider) *LiquidityStatusUseCase {
	return &LiquidityStatusUseCase{
		contracts: contracts,
		provider:  provider,
		peginProvider: peginProvider,
		rpc:       rpc,
        btcWallet: btcWallet,
	}
}

func (useCase *LiquidityStatusUseCase) Run(ctx context.Context) (*pkg.LiquidityStatus, error) {
	if !useCase.provider.GeneralConfiguration(ctx).PublicLiquidityCheck {
		return nil, fmt.Errorf("balance checking is disabled")
	}
	peginLiquidity, err := useCase.peginProvider.CalculateAvailablePeginLiquidity(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching pegin balance: %v", err)
	}
	btcBalance, err := useCase.btcWallet.GetBalance()
	if err != nil {
		return nil, fmt.Errorf("error fetching BTC balance: %v", err)
	}
	return &pkg.LiquidityStatus{
		Available: pkg.Available{
			Pegin:  peginLiquidity,
			Pegout: btcBalance,
		},
	}, nil
}

