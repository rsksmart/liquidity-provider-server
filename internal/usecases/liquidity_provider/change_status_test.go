package liquidity_provider_test

import (
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestChangeStatusUseCase_Run(t *testing.T) {
	const address = "0x02"
	discovery := &mocks.DiscoveryContractMock{}
	discovery.On("GetProvider", address).Return(lp.RegisteredLiquidityProvider{Id: 2, Address: address}, nil).Once()
	discovery.On("SetProviderStatus", uint64(2), false).Return(nil).Once()

	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return(address)

	contracts := blockchain.RskContracts{Discovery: discovery}
	err := liquidity_provider.NewChangeStatusUseCase(contracts, provider).Run(false)

	discovery.AssertExpectations(t)
	require.NoError(t, err)
}

func TestChangeStatusUseCase_Run_Fail(t *testing.T) {
	const address = "0x01"
	discovery := &mocks.DiscoveryContractMock{}
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(address).Once()
	discovery.On("GetProvider", address).Return(
		lp.RegisteredLiquidityProvider{},
		assert.AnError,
	).Once()
	contracts := blockchain.RskContracts{Discovery: discovery}
	err := liquidity_provider.NewChangeStatusUseCase(contracts, provider).Run(false)
	discovery.AssertExpectations(t)
	require.Error(t, err)

	discovery.On("GetProvider", address).Return(lp.RegisteredLiquidityProvider{Id: 1, Address: address}, nil).Once()
	provider.On("RskAddress").Return(address)
	discovery.On("SetProviderStatus", mock.Anything, mock.Anything).Return(errors.New("some error")).Once()
	err = liquidity_provider.NewChangeStatusUseCase(contracts, provider).Run(false)
	discovery.AssertExpectations(t)
	require.Error(t, err)
}
