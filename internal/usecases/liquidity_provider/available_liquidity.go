package liquidity_provider

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/pkg"
)

type LiquidityStatusUseCase struct {
	contracts      blockchain.RskContracts
	provider       liquidity_provider.LiquidityProvider
	peginProvider  liquidity_provider.PeginLiquidityProvider
	pegoutProvider liquidity_provider.PegoutLiquidityProvider
	rpc            blockchain.Rpc
	btcWallet      blockchain.BitcoinWallet
}

func NewLiquidityStatusUseCase(contracts blockchain.RskContracts, provider liquidity_provider.LiquidityProvider, rpc blockchain.Rpc, btcWallet blockchain.BitcoinWallet, peginProvider liquidity_provider.PeginLiquidityProvider, pegoutProvider liquidity_provider.PegoutLiquidityProvider) *LiquidityStatusUseCase {
	return &LiquidityStatusUseCase{
		contracts:      contracts,
		provider:       provider,
		peginProvider:  peginProvider,
		pegoutProvider: pegoutProvider,
		rpc:            rpc,
		btcWallet:      btcWallet,
	}
}

func (useCase *LiquidityStatusUseCase) Run(ctx context.Context) (pkg.LiquidityStatus, error) {
	if !useCase.provider.GeneralConfiguration(ctx).PublicLiquidityCheck {
		return pkg.LiquidityStatus{}, usecases.WrapUseCaseError(usecases.CheckLiquidity, usecases.PublicLiquidityCheckDisabledError)
	}
	peginLiquidity, err := useCase.peginProvider.CalculateAvailablePeginLiquidity(ctx)
	if err != nil {
		return pkg.LiquidityStatus{}, usecases.WrapUseCaseError(usecases.UseCaseId(usecases.CheckLiquidity), err)
	}
	pegoutLiquidity, err := useCase.pegoutProvider.CalculateAvailablePegoutLiquidity(ctx)
	if err != nil {
		return pkg.LiquidityStatus{}, usecases.WrapUseCaseError(usecases.UseCaseId(usecases.CheckLiquidity), err)
	}
	response := pkg.LiquidityStatus{
		Available: pkg.Available{
			Pegin:  peginLiquidity,
			Pegout: pegoutLiquidity,
		},
	}
	return response, nil
}
