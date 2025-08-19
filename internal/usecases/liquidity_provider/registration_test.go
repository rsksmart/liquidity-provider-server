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

func TestRegistrationUseCase_Run_RegisterAgain(t *testing.T) {
	t.Run("should not register again if already registered", func(t *testing.T) {
		t.Run("after adding pegin collateral", func(t *testing.T) {
			lbc := &mocks.LiquidityBridgeContractMock{}
			lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
			lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(900), nil)
			lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
			lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
			lbc.On("IsOperationalPegout", mock.Anything).Return(true, nil)
			lbc.On("AddCollateral", entities.NewUWei(100)).Return(nil)
			lbc.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
				Id: 1, Address: test.AnyAddress, Name: test.AnyString,
				ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
			}, nil).Once()
			provider := &mocks.ProviderMock{}
			provider.On("RskAddress").Return("rskAddress")
			contracts := blockchain.RskContracts{Lbc: lbc}
			useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
			params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, "both")
			id, err := useCase.Run(params)
			lbc.AssertExpectations(t)
			lbc.AssertNotCalled(t, "RegisterProvider")
			provider.AssertExpectations(t)
			assert.Equal(t, int64(0), id)
			require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
		})
		t.Run("after adding pegout collateral", func(t *testing.T) {
			lbc := &mocks.LiquidityBridgeContractMock{}
			lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
			lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
			lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(900), nil)
			lbc.On("IsOperationalPegin", mock.Anything).Return(true, nil)
			lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
			lbc.On("AddPegoutCollateral", entities.NewUWei(100)).Return(nil)
			lbc.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
				Id: 1, Address: test.AnyAddress, Name: test.AnyString,
				ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
			}, nil).Twice()
			provider := &mocks.ProviderMock{}
			provider.On("RskAddress").Return("rskAddress")
			contracts := blockchain.RskContracts{Lbc: lbc}
			useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
			params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, "both")
			id, err := useCase.Run(params)
			lbc.AssertExpectations(t)
			lbc.AssertNotCalled(t, "RegisterProvider")
			provider.AssertExpectations(t)
			assert.Equal(t, int64(0), id)
			require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
		})
	})
}

func TestRegistrationUseCase_Run_AlreadyRegistered(t *testing.T) {
	lbc := &mocks.LiquidityBridgeContractMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(true, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(true, nil)
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, "both")
	id, err := useCase.Run(params)
	lbc.AssertExpectations(t)
	provider.AssertExpectations(t)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run_ValidateParams(t *testing.T) {
	cases := []blockchain.ProviderRegistrationParams{
		blockchain.NewProviderRegistrationParams("", test.AnyUrl, true, "both"),
		blockchain.NewProviderRegistrationParams("name", "", true, "both"),
		blockchain.NewProviderRegistrationParams("name", test.AnyUrl, false, "both"),
		blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, "anything"),
		blockchain.NewProviderRegistrationParams("", test.AnyUrl, true, ""),
	}
	lbc := &mocks.LiquidityBridgeContractMock{}
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
	lbc := &mocks.LiquidityBridgeContractMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
	lbc.On("AddCollateral", test.AnyWei).Return(nil)
	lbc.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
		Id: 1, Address: test.AnyAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.PeginProvider,
	}, nil).Once()
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, "pegin")
	id, err := useCase.Run(params)
	lbc.AssertExpectations(t)
	provider.AssertExpectations(t)
	lbc.AssertNotCalled(t, "AddPegoutCollateral", test.AnyWei)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run_AddPegoutCollateralIfNotOperational(t *testing.T) {
	lbc := &mocks.LiquidityBridgeContractMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
	lbc.On("AddPegoutCollateral", test.AnyWei).Return(nil)
	lbc.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
		Id: 1, Address: test.AnyAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.PegoutProvider,
	}, nil).Twice()
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, "pegout")
	id, err := useCase.Run(params)
	lbc.AssertExpectations(t)
	provider.AssertExpectations(t)
	lbc.AssertNotCalled(t, "AddCollateral", test.AnyWei)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run_AddCollateralIfNotOperational(t *testing.T) {
	lbc := &mocks.LiquidityBridgeContractMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(999), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(999), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
	lbc.On("AddCollateral", test.AnyWei).Return(nil)
	lbc.On("AddPegoutCollateral", test.AnyWei).Return(nil)
	lbc.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
		Id: 1, Address: test.AnyAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
	}, nil).Twice()
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, "both")
	id, err := useCase.Run(params)
	lbc.AssertExpectations(t)
	provider.AssertExpectations(t)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.AlreadyRegisteredError)
}

func TestRegistrationUseCase_Run(t *testing.T) {
	lbc := &mocks.LiquidityBridgeContractMock{}
	lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
	lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
	lbc.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{}, lp.ProviderNotFoundError).Twice()
	lbc.On(
		"RegisterProvider",
		mock.AnythingOfType("blockchain.TransactionConfig"),
		mock.AnythingOfType("ProviderRegistrationParams")).
		Return(int64(1), nil)
	provider := &mocks.ProviderMock{}
	provider.On("RskAddress").Return("rskAddress")
	contracts := blockchain.RskContracts{Lbc: lbc}
	useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, "both")
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
		lbc := &mocks.LiquidityBridgeContractMock{}
		testCase.Value(lbc) // setup function
		provider := &mocks.ProviderMock{}
		provider.On("RskAddress").Return("rskAddress")
		contracts := blockchain.RskContracts{Lbc: lbc}
		useCase := liquidity_provider.NewRegistrationUseCase(contracts, provider)
		params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, "both")
		id, err := useCase.Run(params)
		lbc.AssertExpectations(t)
		assert.Equal(t, int64(0), id)
		require.Error(t, err)
	}
}

// nolint:funlen
func registrationUseCaseUnexpectedErrorSetups() test.Table[func(mock *mocks.LiquidityBridgeContractMock), error] {
	return test.Table[func(mock *mocks.LiquidityBridgeContractMock), error]{
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(1000), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
				lbc.On("IsOperationalPegout", mock.Anything).Return(false, assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
				lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
				lbc.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{}, lp.ProviderNotFoundError)
				lbc.On(
					"RegisterProvider",
					mock.AnythingOfType("blockchain.TransactionConfig"),
					mock.AnythingOfType("ProviderRegistrationParams")).
					Return(int64(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
				lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
				lbc.On("AddCollateral", test.AnyWei).Return(assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock) {
				lbc.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
				lbc.On("GetCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				lbc.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(10), nil)
				lbc.On("IsOperationalPegin", mock.Anything).Return(false, nil)
				lbc.On("IsOperationalPegout", mock.Anything).Return(false, nil)
				lbc.On("AddCollateral", test.AnyWei).Return(nil)
				lbc.On("GetProvider", mock.Anything).Return(lp.RegisteredLiquidityProvider{
					Id: 1, Address: test.AnyAddress, Name: test.AnyString,
					ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
				}, nil).Once()
				lbc.On("AddPegoutCollateral", test.AnyWei).Return(assert.AnError)
			},
		},
	}
}
