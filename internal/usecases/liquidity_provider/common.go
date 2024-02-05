package liquidity_provider

import (
	"cmp"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"slices"
)

func ValidateConfiguredProvider(
	provider entities.LiquidityProvider,
	lbc blockchain.LiquidityBridgeContract,
) (uint64, error) {
	var err error
	var providers []entities.RegisteredLiquidityProvider

	if providers, err = lbc.GetProviders(); err != nil {
		return 0, err
	}

	index, found := slices.BinarySearchFunc(
		providers,
		entities.RegisteredLiquidityProvider{Address: provider.RskAddress()},
		func(a, b entities.RegisteredLiquidityProvider) int {
			return cmp.Compare(a.Address, b.Address)
		},
	)
	if !found {
		return 0, usecases.ProviderConfigurationError
	}
	return providers[index].Id, nil
}
