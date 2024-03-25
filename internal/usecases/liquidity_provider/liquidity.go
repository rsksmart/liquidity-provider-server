package liquidity_provider

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
)

type LiquidityStatusUseCase struct {
	contracts blockchain.RskContracts
	provider  liquidity_provider.LiquidityProvider
	rpc       blockchain.Rpc
    btcWallet blockchain.BitcoinWallet
}

type LiquidityStatus struct {
	Available struct {
		Pegin  *entities.Wei
		Pegout *entities.Wei
	}
}

func NewLiquidityStatusUseCase(contracts blockchain.RskContracts, provider liquidity_provider.LiquidityProvider, rpc blockchain.Rpc,btcWallet blockchain.BitcoinWallet) *LiquidityStatusUseCase {
	return &LiquidityStatusUseCase{
		contracts: contracts,
		provider:  provider,
		rpc:       rpc,
        btcWallet: btcWallet,
	}
}

func (useCase *LiquidityStatusUseCase) Run(ctx context.Context) (*LiquidityStatus, error) {
	var err error
	var lbcBalance, lpBalance,btcBalance, totalPegin *entities.Wei
	generalConfiguration := useCase.provider.GeneralConfiguration(ctx)
	if generalConfiguration.BalanceCheck {
		if lbcBalance, err = useCase.contracts.Lbc.GetBalance(useCase.provider.RskAddress()); err != nil {
			return nil, fmt.Errorf("error fetching LBC balance: %v", err)
		}
		if lpBalance, err = useCase.rpc.Rsk.GetBalance(ctx, useCase.provider.RskAddress()); err != nil {
			return nil, fmt.Errorf("error fetching LP balance: %v", err)
		}
        if btcBalance, err = useCase.btcWallet.GetBalance(); err != nil {
			return nil, fmt.Errorf("error fetching BTC balance: %v", err)
		}
		totalPegin = new(entities.Wei).Add(lpBalance, lbcBalance)
		return &LiquidityStatus{
			Available: struct {
				Pegin  *entities.Wei
				Pegout *entities.Wei
			}{
				Pegin:  totalPegin,
				Pegout: btcBalance,
			},
		}, nil
	}
	return nil, fmt.Errorf("balance checking is disabled")
}
