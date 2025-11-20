package liquidity_provider_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegistrationUseCase_Run_Paused(t *testing.T) {
	t.Run("should return error if discovery is paused", func(t *testing.T) {
		collateral := new(mocks.CollateralManagementContractMock)
		collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
		discovery := new(mocks.DiscoveryContractMock)
		discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
		provider := &mocks.ProviderMock{}
		contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
		useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
		params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider)
		id, err := useCase.Run(params)
		assert.Equal(t, int64(0), id)
		require.ErrorIs(t, err, blockchain.ContractPausedError)
	})
	t.Run("should return error if collateral management is paused", func(t *testing.T) {
		collateral := new(mocks.CollateralManagementContractMock)
		collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
		discovery := new(mocks.DiscoveryContractMock)
		discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
		provider := &mocks.ProviderMock{}
		contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
		useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
		params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider)
		id, err := useCase.Run(params)
		assert.Equal(t, int64(0), id)
		require.ErrorIs(t, err, blockchain.ContractPausedError)
	})
}

func TestRegistrationUseCase_Run_RegisterAgain(t *testing.T) {
	t.Run("should not register again if already registered", func(t *testing.T) {
		t.Run("after adding pegin collateral", func(t *testing.T) {
			collateral := new(mocks.CollateralManagementContractMock)
			discovery := new(mocks.DiscoveryContractMock)
			collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
			collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
			collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(900), nil)
			collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
			discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
			discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(true, nil)
			discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
			collateral.On("AddCollateral", entities.NewUWei(100)).Return(nil)
			discovery.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
				Id: 1, Address: test.AnyAddress, Name: test.AnyString,
				ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
			}, nil).Once()
			provider := &mocks.ProviderMock{}
			provider.On("RskAddress").Return("rskAddress")
			contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
			useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
			params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider)
			id, err := useCase.Run(params)
			discovery.AssertExpectations(t)
			discovery.AssertNotCalled(t, "RegisterProvider")
			collateral.AssertExpectations(t)
			provider.AssertExpectations(t)
			assert.Equal(t, int64(0), id)
			require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
		})
		t.Run("after adding pegout collateral", func(t *testing.T) {
			collateral := new(mocks.CollateralManagementContractMock)
			discovery := new(mocks.DiscoveryContractMock)
			collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
			collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
			collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
			collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(900), nil)
			discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
			discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(true, nil)
			discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
			collateral.On("AddPegoutCollateral", entities.NewUWei(100)).Return(nil)
			discovery.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
				Id: 1, Address: test.AnyAddress, Name: test.AnyString,
				ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
			}, nil).Twice()
			provider := &mocks.ProviderMock{}
			provider.On("RskAddress").Return("rskAddress")
			contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
			useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
			params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider)
			id, err := useCase.Run(params)
			discovery.AssertExpectations(t)
			discovery.AssertNotCalled(t, "RegisterProvider")
			collateral.AssertExpectations(t)
			provider.AssertExpectations(t)
			assert.Equal(t, int64(0), id)
			require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
		})
	})
}

func TestRegistrationUseCase_Run_AlreadyRegistered(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	discovery := new(mocks.DiscoveryContractMock)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(true, nil)
	discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(true, nil)
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider)
	id, err := useCase.Run(params)
	collateral.AssertExpectations(t)
	discovery.AssertExpectations(t)
	provider.AssertExpectations(t)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run_ValidateParams(t *testing.T) {
	cases := []blockchain.ProviderRegistrationParams{
		blockchain.NewProviderRegistrationParams("", test.AnyUrl, true, lp.FullProvider),
		blockchain.NewProviderRegistrationParams("name", "", true, lp.FullProvider),
		blockchain.NewProviderRegistrationParams("name", test.AnyUrl, false, lp.FullProvider),
		blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, 5),
		blockchain.NewProviderRegistrationParams("", test.AnyUrl, true, -1),
	}
	collateral := new(mocks.CollateralManagementContractMock)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	discovery := new(mocks.DiscoveryContractMock)
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	provider := &mocks.ProviderMock{}
	contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	var id int64
	var err error
	for _, c := range cases {
		id, err = useCase.Run(c)
		assert.Equal(t, int64(0), id)
		require.Error(t, err)
	}
}

func TestRegistrationUseCase_Run_AddPeginCollateralIfNotOperational(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	discovery := new(mocks.DiscoveryContractMock)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
	discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
	collateral.On("AddCollateral", test.AnyWei).Return(nil)
	discovery.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
		Id: 1, Address: test.AnyAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.PeginProvider,
	}, nil).Once()
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.PeginProvider)
	id, err := useCase.Run(params)
	collateral.AssertExpectations(t)
	discovery.AssertExpectations(t)
	provider.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddPegoutCollateral", test.AnyWei)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run_AddPegoutCollateralIfNotOperational(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	discovery := new(mocks.DiscoveryContractMock)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
	discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
	collateral.On("AddPegoutCollateral", test.AnyWei).Return(nil)
	discovery.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
		Id: 1, Address: test.AnyAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.PegoutProvider,
	}, nil).Twice()
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.PegoutProvider)
	id, err := useCase.Run(params)
	discovery.AssertExpectations(t)
	collateral.AssertExpectations(t)
	provider.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddCollateral", test.AnyWei)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run_AddCollateralIfNotOperational(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	discovery := new(mocks.DiscoveryContractMock)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(999), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(999), nil)
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
	discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
	collateral.On("AddCollateral", test.AnyWei).Return(nil)
	collateral.On("AddPegoutCollateral", test.AnyWei).Return(nil)
	discovery.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
		Id: 1, Address: test.AnyAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
	}, nil).Twice()
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider)
	id, err := useCase.Run(params)
	discovery.AssertExpectations(t)
	collateral.AssertExpectations(t)
	provider.AssertExpectations(t)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run(t *testing.T) {
	collateral := new(mocks.CollateralManagementContractMock)
	discovery := new(mocks.DiscoveryContractMock)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
	discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
	discovery.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{}, lp.ProviderNotFoundError).Twice()
	discovery.On(
		"RegisterProvider",
		mock.AnythingOfType("blockchain.TransactionConfig"),
		mock.AnythingOfType("ProviderRegistrationParams")).
		Return(int64(1), nil)
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider)
	id, err := useCase.Run(params)
	collateral.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddCollateral")
	collateral.AssertNotCalled(t, "AddPegoutCollateral")
	discovery.AssertExpectations(t)
	provider.AssertExpectations(t)
	assert.Equal(t, int64(1), id)
	require.NoError(t, err)
}

func TestRegistrationUseCase_Run_ErrorHandling(t *testing.T) {
	cases := registrationUseCaseUnexpectedErrorSetups()

	for _, testCase := range cases {
		collateral := new(mocks.CollateralManagementContractMock)
		discovery := new(mocks.DiscoveryContractMock)
		collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
		discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
		testCase.Value(collateral, discovery) // setup function
		provider := &mocks.ProviderMock{}
		provider.On("RskAddress").Return("rskAddress")
		contracts := blockchain.RskContracts{CollateralManagement: collateral, Discovery: discovery}
		useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
		params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider)
		id, err := useCase.Run(params)
		collateral.AssertExpectations(t)
		discovery.AssertExpectations(t)
		assert.Equal(t, int64(0), id)
		require.Error(t, err)
	}
}

// nolint:funlen
func registrationUseCaseUnexpectedErrorSetups() test.Table[func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock), error] {
	return test.Table[func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock), error]{
		{
			Value: func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(0), assert.AnError)
			},
		},
		{
			Value: func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), assert.AnError)
			},
		},
		{
			Value: func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), assert.AnError)
			},
		},
		{
			Value: func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, assert.AnError)
			},
		},
		{
			Value: func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
				discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, assert.AnError)
			},
		},
		{
			Value: func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), nil)
				collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), nil)
				discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
				discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
				discovery.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{}, lp.ProviderNotFoundError)
				discovery.On(
					"RegisterProvider",
					mock.AnythingOfType("blockchain.TransactionConfig"),
					mock.AnythingOfType("ProviderRegistrationParams")).
					Return(int64(0), assert.AnError)
			},
		},
		{
			Value: func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
				discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
				collateral.On("AddCollateral", test.AnyWei).Return(assert.AnError)
			},
		},
		{
			Value: func(collateral *mocks.CollateralManagementContractMock, discovery *mocks.DiscoveryContractMock) {
				collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
				discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
				collateral.On("AddCollateral", test.AnyWei).Return(nil)
				discovery.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
					Id: 1, Address: test.AnyAddress, Name: test.AnyString,
					ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
				}, nil).Once()
				collateral.On("AddPegoutCollateral", test.AnyWei).Return(assert.AnError)
			},
		},
	}
}
