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

const testRskAddress = "0x1234567890abcdef1234567890abcdef12345678"

func newRegistrationUseCase(contracts blockchain.RskContracts, provider lp.LiquidityProvider) *liquidity_provider.RegistrationUseCase {
	return liquidity_provider.NewRegistrationUseCase(contracts, provider, 0)
}

// ── Approved state ────────────────────────────────────────────────────────────

func TestRegistrationUseCase_Run_Approved(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateApproved, nil)
	discovery.EXPECT().GetProvider(testRskAddress).Return(lp.RegisteredLiquidityProvider{
		Id: 7, Address: testRskAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
	}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "RegisterProvider")
	collateral.AssertNotCalled(t, "PausedStatus")
	assert.Equal(t, int64(7), id)
	require.NoError(t, err)
}

func TestRegistrationUseCase_Run_Approved_GetProviderError(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateApproved, nil)
	discovery.EXPECT().GetProvider(testRskAddress).Return(lp.RegisteredLiquidityProvider{}, assert.AnError)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	assert.Equal(t, int64(0), id)
	require.Error(t, err)
}

// ── Pending state ─────────────────────────────────────────────────────────────

func TestRegistrationUseCase_Run_Pending_ThenApproved(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	// Initial state check → Pending
	discovery.On("GetRegistrationState", testRskAddress).
		Return(blockchain.RegistrationStatePending, nil).Once()
	// Poll iteration 1 → still Pending
	discovery.On("GetRegistrationState", testRskAddress).
		Return(blockchain.RegistrationStatePending, nil).Once()
	// Poll iteration 2 → Approved
	discovery.On("GetRegistrationState", testRskAddress).
		Return(blockchain.RegistrationStateApproved, nil).Once()
	discovery.EXPECT().GetProvider(testRskAddress).Return(lp.RegisteredLiquidityProvider{
		Id: 3, Address: testRskAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
	}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "RegisterProvider")
	assert.Equal(t, int64(3), id)
	require.NoError(t, err)
}

func TestRegistrationUseCase_Run_Pending_ThenRejected(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.On("GetRegistrationState", testRskAddress).
		Return(blockchain.RegistrationStatePending, nil).Once()
	discovery.On("GetRegistrationState", testRskAddress).
		Return(blockchain.RegistrationStateRejected, nil).Once()

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "RegisterProvider")
	discovery.AssertNotCalled(t, "GetProvider")
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.RegistrationRejectedError)
}

// ── Rejected state ────────────────────────────────────────────────────────────

func TestRegistrationUseCase_Run_Rejected(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateRejected, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "RegisterProvider")
	discovery.AssertNotCalled(t, "GetProvider")
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.RegistrationRejectedError)
}

// ── None / Withdrawn states ───────────────────────────────────────────────────

func TestRegistrationUseCase_Run_None_RegisterAndWaitForApproval(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.On("GetRegistrationState", testRskAddress).
		Return(blockchain.RegistrationStateNone, nil).Once()
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
	discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
	discovery.On(
		"RegisterProvider",
		mock.AnythingOfType("blockchain.TransactionConfig"),
		mock.AnythingOfType("ProviderRegistrationParams"),
	).Return(int64(1), nil)
	// poll → Approved
	discovery.On("GetRegistrationState", testRskAddress).
		Return(blockchain.RegistrationStateApproved, nil).Once()
	discovery.EXPECT().GetProvider(testRskAddress).Return(lp.RegisteredLiquidityProvider{
		Id: 1, Address: testRskAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
	}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	params := blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider)
	id, err := useCase.Run(params)

	discovery.AssertExpectations(t)
	collateral.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddCollateral")
	collateral.AssertNotCalled(t, "AddPegoutCollateral")
	assert.Equal(t, int64(1), id)
	require.NoError(t, err)
}

func TestRegistrationUseCase_Run_Withdrawn_RegisterAndWaitForApproval(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.On("GetRegistrationState", testRskAddress).
		Return(blockchain.RegistrationStateWithdrawn, nil).Once()
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
	discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
	discovery.On(
		"RegisterProvider",
		mock.AnythingOfType("blockchain.TransactionConfig"),
		mock.AnythingOfType("ProviderRegistrationParams"),
	).Return(int64(2), nil)
	discovery.On("GetRegistrationState", testRskAddress).
		Return(blockchain.RegistrationStateApproved, nil).Once()
	discovery.EXPECT().GetProvider(testRskAddress).Return(lp.RegisteredLiquidityProvider{
		Id: 2, Address: testRskAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
	}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	assert.Equal(t, int64(2), id)
	require.NoError(t, err)
}

// ── Pause check (only in None/Withdrawn path) ─────────────────────────────────

func TestRegistrationUseCase_Run_NoneOrWithdrawn_PausedDiscovery(t *testing.T) {
	for _, state := range []blockchain.RegistrationState{blockchain.RegistrationStateNone, blockchain.RegistrationStateWithdrawn} {
		discovery := new(mocks.DiscoveryContractMock)
		collateral := new(mocks.CollateralManagementContractMock)
		provider := &mocks.ProviderMock{}

		provider.On("RskAddress").Return(testRskAddress)
		discovery.On("GetRegistrationState", testRskAddress).Return(state, nil).Once()
		collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
		discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
		discovery.EXPECT().GetAddress().Return("discovery-contract")

		contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
		useCase := newRegistrationUseCase(contracts, provider)
		id, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

		discovery.AssertNotCalled(t, "RegisterProvider")
		assert.Equal(t, int64(0), id)
		require.ErrorIs(t, err, blockchain.ContractPausedError)
	}
}

func TestRegistrationUseCase_Run_NoneOrWithdrawn_PausedCollateral(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.On("GetRegistrationState", testRskAddress).Return(blockchain.RegistrationStateNone, nil).Once()
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
	collateral.EXPECT().GetAddress().Return("collateral-contract")

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertNotCalled(t, "RegisterProvider")
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, blockchain.ContractPausedError)
}

// Approved and Pending states must NOT trigger the pause check.
func TestRegistrationUseCase_Run_ApprovedOrPending_NoPauseCheck(t *testing.T) {
	for _, state := range []blockchain.RegistrationState{blockchain.RegistrationStateApproved, blockchain.RegistrationStatePending} {
		discovery := new(mocks.DiscoveryContractMock)
		collateral := new(mocks.CollateralManagementContractMock)
		provider := &mocks.ProviderMock{}

		provider.On("RskAddress").Return(testRskAddress)
		discovery.On("GetRegistrationState", testRskAddress).Return(state, nil).Once()
		if state == blockchain.RegistrationStateApproved {
			discovery.On("GetProvider", testRskAddress).Return(lp.RegisteredLiquidityProvider{Id: 1}, nil)
		} else {
			// Pending → immediately Approved on first poll
			discovery.On("GetRegistrationState", testRskAddress).Return(blockchain.RegistrationStateApproved, nil).Once()
			discovery.On("GetProvider", testRskAddress).Return(lp.RegisteredLiquidityProvider{Id: 1}, nil)
		}

		contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
		useCase := newRegistrationUseCase(contracts, provider)
		_, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))
		require.NoError(t, err)

		collateral.AssertNotCalled(t, "PausedStatus")
		discovery.AssertNotCalled(t, "PausedStatus")
	}
}

// ── Validate params (only in None/Withdrawn path) ─────────────────────────────

func TestRegistrationUseCase_Run_NoneOrWithdrawn_ValidateParams(t *testing.T) {
	invalidParams := []blockchain.ProviderRegistrationParams{
		blockchain.NewProviderRegistrationParams("", test.AnyUrl, true, lp.FullProvider),
		blockchain.NewProviderRegistrationParams("name", "", true, lp.FullProvider),
		blockchain.NewProviderRegistrationParams("name", test.AnyUrl, false, lp.FullProvider),
		blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, 5),
	}
	for _, params := range invalidParams {
		discovery := new(mocks.DiscoveryContractMock)
		collateral := new(mocks.CollateralManagementContractMock)
		provider := &mocks.ProviderMock{}

		provider.On("RskAddress").Return(testRskAddress)
		discovery.On("GetRegistrationState", testRskAddress).Return(blockchain.RegistrationStateNone, nil).Once()
		discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
		collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)

		contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
		useCase := newRegistrationUseCase(contracts, provider)
		id, err := useCase.Run(params)

		discovery.AssertNotCalled(t, "RegisterProvider")
		assert.Equal(t, int64(0), id)
		require.Error(t, err)
	}
}

// ── Collateral top-ups (None/Withdrawn path) ──────────────────────────────────

func TestRegistrationUseCase_Run_AddPeginCollateral(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.On("GetRegistrationState", testRskAddress).Return(blockchain.RegistrationStateNone, nil).Once()
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(900), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
	discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
	collateral.On("AddCollateral", entities.NewUWei(100)).Return(nil)
	discovery.On(
		"RegisterProvider",
		mock.AnythingOfType("blockchain.TransactionConfig"),
		mock.AnythingOfType("ProviderRegistrationParams"),
	).Return(int64(1), nil)
	discovery.On("GetRegistrationState", testRskAddress).Return(blockchain.RegistrationStateApproved, nil).Once()
	discovery.On("GetProvider", testRskAddress).Return(lp.RegisteredLiquidityProvider{Id: 1}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	_, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.PeginProvider))

	collateral.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddPegoutCollateral")
	require.NoError(t, err)
}

func TestRegistrationUseCase_Run_AddPegoutCollateral(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.On("GetRegistrationState", testRskAddress).Return(blockchain.RegistrationStateNone, nil).Once()
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.On("GetMinimumCollateral").Return(entities.NewWei(1000), nil)
	collateral.On("GetCollateral", mock.Anything).Return(entities.NewWei(0), nil)
	collateral.On("GetPegoutCollateral", mock.Anything).Return(entities.NewWei(900), nil)
	discovery.EXPECT().IsOperational(lp.PeginProvider, mock.Anything).Return(false, nil)
	discovery.EXPECT().IsOperational(lp.PegoutProvider, mock.Anything).Return(false, nil)
	collateral.On("AddPegoutCollateral", entities.NewUWei(100)).Return(nil)
	discovery.On(
		"RegisterProvider",
		mock.AnythingOfType("blockchain.TransactionConfig"),
		mock.AnythingOfType("ProviderRegistrationParams"),
	).Return(int64(1), nil)
	discovery.On("GetRegistrationState", testRskAddress).Return(blockchain.RegistrationStateApproved, nil).Once()
	discovery.On("GetProvider", testRskAddress).Return(lp.RegisteredLiquidityProvider{Id: 1}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	_, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.PegoutProvider))

	collateral.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddCollateral")
	require.NoError(t, err)
}

// ── GetRegistrationState error ────────────────────────────────────────────────

func TestRegistrationUseCase_Run_GetRegistrationStateError(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateNone, assert.AnError)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertNotCalled(t, "RegisterProvider")
	assert.Equal(t, int64(0), id)
	require.Error(t, err)
}
