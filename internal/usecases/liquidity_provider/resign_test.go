package liquidity_provider_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResignUseCase_Run(t *testing.T) {
	collateral := &mocks.CollateralManagementContractMock{}
	discovery := &mocks.DiscoveryContractMock{}
	provider := &mocks.ProviderMock{}
	const address = "0x01"
	provider.On("RskAddress").Return(address)
	discovery.On("GetProvider", address).Return(lp.RegisteredLiquidityProvider{Id: 1, Address: address}, nil)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("ProviderResign").Return(nil).Once()
	contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
	useCase := liquidity_provider.NewResignUseCase(contracts, provider)
	err := useCase.Run()
	discovery.AssertExpectations(t)
	collateral.AssertExpectations(t)
	require.NoError(t, err)
}

func TestResignUseCase_Run_Paused(t *testing.T) {
	provider := &mocks.ProviderMock{}
	collateral := new(mocks.CollateralManagementContractMock)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
	collateral.EXPECT().GetAddress().Return("test-contract")
	contracts := blockchain.RskContracts{CollateralManagement: collateral}
	useCase := liquidity_provider.NewResignUseCase(contracts, provider)
	err := useCase.Run()
	require.ErrorIs(t, err, blockchain.ContractPausedError)
}

func TestResignUseCase_Run_NotRegistered(t *testing.T) {
	discovery := &mocks.DiscoveryContractMock{}
	provider := &mocks.ProviderMock{}
	collateral := &mocks.CollateralManagementContractMock{}
	const address = "0x01"
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	provider.On("RskAddress").Return(address)
	discovery.On("GetProvider", address).Return(lp.RegisteredLiquidityProvider{}, assert.AnError).Once()
	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := liquidity_provider.NewResignUseCase(contracts, provider)
	err := useCase.Run()
	discovery.AssertExpectations(t)
	require.ErrorIs(t, err, usecases.ProviderConfigurationError)
}

func TestResignUseCase_Run_Error(t *testing.T) {
	collateral := &mocks.CollateralManagementContractMock{}
	discovery := &mocks.DiscoveryContractMock{}
	provider := &mocks.ProviderMock{}
	const address = "0x01"
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	provider.On("RskAddress").Return(address)
	discovery.On("GetProvider", address).Return(lp.RegisteredLiquidityProvider{Id: 1, Address: address}, nil)
	collateral.On("ProviderResign").Return(assert.AnError).Once()
	contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
	useCase := liquidity_provider.NewResignUseCase(contracts, provider)
	err := useCase.Run()
	collateral.AssertExpectations(t)
	discovery.AssertExpectations(t)
	require.Error(t, err)
}
