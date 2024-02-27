package liquidity_provider_test

import (
	"errors"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateConfiguredProvider(t *testing.T) {
	lbc := &mocks.LbcMock{}
	lbc.On("GetProviders").Return([]lpEntity.RegisteredLiquidityProvider{
		{
			Id:           1,
			Address:      "0x01",
			Name:         "one",
			ApiBaseUrl:   "api1.com",
			Status:       true,
			ProviderType: "both",
		},
		{
			Id:           2,
			Address:      "0x02",
			Name:         "two",
			ApiBaseUrl:   "api2.com",
			Status:       true,
			ProviderType: "pegin",
		},
		{
			Id:           3,
			Address:      "0x03",
			Name:         "three",
			ApiBaseUrl:   "api3.com",
			Status:       true,
			ProviderType: "pegout",
		},
	}, nil)

	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("0x02")

	id, err := liquidity_provider.ValidateConfiguredProvider(provider, lbc)
	assert.Equal(t, uint64(2), id)
	require.NoError(t, err)
}

func TestValidateConfiguredProvider_Fail(t *testing.T) {
	lbc := &mocks.LbcMock{}
	var provider *mocks.ProviderMock = nil
	lbc.On("GetProviders").Return([]lpEntity.RegisteredLiquidityProvider{}, errors.New("some error")).Once()
	id, err := liquidity_provider.ValidateConfiguredProvider(provider, lbc)
	assert.Equal(t, uint64(0), id)
	require.Error(t, err)

	provider = &mocks.ProviderMock{}
	provider.On("RskAddress").Return("0x02")
	lbc.On("GetProviders").Return([]lpEntity.RegisteredLiquidityProvider{
		{
			Id:           3,
			Address:      "0x03",
			Name:         "three",
			ApiBaseUrl:   "api3.com",
			Status:       true,
			ProviderType: "pegout",
		},
	}, nil).Once()
	id, err = liquidity_provider.ValidateConfiguredProvider(provider, lbc)
	assert.Equal(t, uint64(0), id)
	require.ErrorIs(t, err, usecases.ProviderConfigurationError)
}
