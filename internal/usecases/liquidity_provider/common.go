package liquidity_provider

import (
	"cmp"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"slices"
)

func ValidateConfiguredProvider(
	provider liquidity_provider.LiquidityProvider,
	lbc blockchain.LiquidityBridgeContract,
) (uint64, error) {
	var err error
	var providers []liquidity_provider.RegisteredLiquidityProvider

	if providers, err = lbc.GetProviders(); err != nil {
		return 0, err
	}

	index, found := slices.BinarySearchFunc(
		providers,
		liquidity_provider.RegisteredLiquidityProvider{Address: provider.RskAddress()},
		func(a, b liquidity_provider.RegisteredLiquidityProvider) int {
			return cmp.Compare(a.Address, b.Address)
		},
	)
	if !found {
		return 0, usecases.ProviderConfigurationError
	}
	return providers[index].Id, nil
}
