package liquidity_provider_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRegistrationUseCase_Run_AlreadyRegistered(t *testing.T) {
	lbc := &mocks.LbcMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(true, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(true, nil)
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", "url.com", true, "both")
	id, err := useCase.Run(params)
	lbc.AssertExpectations(t)
	provider.AssertExpectations(t)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run_ValidateParams(t *testing.T) {
	cases := []blockchain.ProviderRegistrationParams{
		blockchain.NewProviderRegistrationParams("", "url.com", true, "both"),
		blockchain.NewProviderRegistrationParams("name", "", true, "both"),
		blockchain.NewProviderRegistrationParams("name", "url.com", false, "both"),
		blockchain.NewProviderRegistrationParams("name", "url.com", true, "anything"),
		blockchain.NewProviderRegistrationParams("", "url.com", true, ""),
	}
	lbc := &mocks.LbcMock{}
	provider := &mocks.ProviderMock{}
	contracts := blockchain.RskContracts{Lbc: lbc}
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
	lbc := &mocks.LbcMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
	lbc.On("AddCollateral", mock.AnythingOfType("*entities.Wei")).Return(nil)
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", "url.com", true, "pegin")
	id, err := useCase.Run(params)
	lbc.AssertExpectations(t)
	provider.AssertExpectations(t)
	lbc.AssertNotCalled(t, "AddPegoutCollateral", mock.AnythingOfType("*entities.Wei"))
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run_AddPegoutCollateralIfNotOperational(t *testing.T) {
	lbc := &mocks.LbcMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
	lbc.On("AddPegoutCollateral", mock.AnythingOfType("*entities.Wei")).Return(nil)
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", "url.com", true, "pegout")
	id, err := useCase.Run(params)
	lbc.AssertExpectations(t)
	provider.AssertExpectations(t)
	lbc.AssertNotCalled(t, "AddCollateral", mock.AnythingOfType("*entities.Wei"))
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run_AddCollateralIfNotOperational(t *testing.T) {
	lbc := &mocks.LbcMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(999), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(999), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
	lbc.On("AddCollateral", mock.AnythingOfType("*entities.Wei")).Return(nil)
	lbc.On("AddPegoutCollateral", mock.AnythingOfType("*entities.Wei")).Return(nil)
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", "url.com", true, "both")
	id, err := useCase.Run(params)
	lbc.AssertExpectations(t)
	provider.AssertExpectations(t)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run(t *testing.T) {
	lbc := &mocks.LbcMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
	lbc.On(
		"RegisterProvider",
		mock.AnythingOfType("blockchain.TransactionConfig"),
		mock.AnythingOfType("ProviderRegistrationParams")).
		Return(int64(1), nil)
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", "url.com", true, "both")
	id, err := useCase.Run(params)
	lbc.AssertExpectations(t)
	lbc.AssertNotCalled(t, "AddCollateral")
	lbc.AssertNotCalled(t, "AddPegoutCollateral")
	provider.AssertExpectations(t)
	assert.Equal(t, int64(1), id)
	require.NoError(t, err)
}

func TestRegistrationUseCase_Run_ErrorHandling(t *testing.T) {
	cases := registrationUseCaseUnexpectedErrorSetups()

	for _, testCase := range cases {
		lbc := &mocks.LbcMock{}
		testCase.Value(lbc) // setup function
		provider := &mocks.ProviderMock{}
		provider.On("RskAddress").Return("rskAddress")
		contracts := blockchain.RskContracts{Lbc: lbc}
		useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
		params := blockchain.NewProviderRegistrationParams("name", "url.com", true, "both")
		id, err := useCase.Run(params)
		lbc.AssertExpectations(t)
		assert.Equal(t, int64(0), id)
		require.Error(t, err)
	}
}

// nolint:funlen
func registrationUseCaseUnexpectedErrorSetups() test.Table[func(mock *mocks.LbcMock), error] {
	return test.Table[func(mock *mocks.LbcMock), error]{
		{
			Value: func(lbc *mocks.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
				lbc.On("IsOperationalPegout", mock.Anything).Return(false, assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
				lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
				lbc.On(
					"RegisterProvider",
					mock.AnythingOfType("blockchain.TransactionConfig"),
					mock.AnythingOfType("ProviderRegistrationParams")).
					Return(int64(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
				lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
				lbc.On("AddCollateral", mock.AnythingOfType("*entities.Wei")).Return(assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LbcMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
				lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
				lbc.On("AddCollateral", mock.AnythingOfType("*entities.Wei")).Return(nil)
				lbc.On("AddPegoutCollateral", mock.AnythingOfType("*entities.Wei")).Return(assert.AnError)
			},
		},
	}
}
