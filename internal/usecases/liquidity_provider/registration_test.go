package liquidity_provider_test

import (
	"context"
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
	return liquidity_provider.NewRegistrationUseCase(contracts, provider)
}

// Boot state: Approved

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
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "WatchRegistrationApproval")
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
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	assert.Equal(t, int64(0), id)
	require.Error(t, err)
}

// Boot state: Rejected

func TestRegistrationUseCase_Run_Rejected(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateRejected, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "WatchRegistrationApproval")
	discovery.AssertNotCalled(t, "RegisterProvider")
	discovery.AssertNotCalled(t, "GetProvider")
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.RegistrationRejectedError)
}

// Boot state: Withdrawn — LP owner backed out; LPS stops waiting (no re-register).

func TestRegistrationUseCase_Run_Withdrawn(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateWithdrawn, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "WatchRegistrationApproval")
	discovery.AssertNotCalled(t, "RegisterProvider")
	discovery.AssertNotCalled(t, "GetProvider")
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.RegistrationWithdrawnError)
}

// Boot state: Pending (e.g., LPS restarted mid-wait) — goes straight to WatchRegistrationApproval.

func TestRegistrationUseCase_Run_Pending_ThenApproved(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStatePending, nil)
	discovery.EXPECT().WatchRegistrationApproval(mock.Anything, testRskAddress).
		Return(blockchain.RegistrationStateApproved, nil).Once()
	discovery.EXPECT().GetProvider(testRskAddress).Return(lp.RegisteredLiquidityProvider{
		Id: 3, Address: testRskAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
	}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "RegisterProvider")
	assert.Equal(t, int64(3), id)
	require.NoError(t, err)
}

func TestRegistrationUseCase_Run_Pending_ThenRejected(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStatePending, nil)
	discovery.EXPECT().WatchRegistrationApproval(mock.Anything, testRskAddress).
		Return(blockchain.RegistrationStateRejected, nil).Once()

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "RegisterProvider")
	discovery.AssertNotCalled(t, "GetProvider")
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.RegistrationRejectedError)
}

func TestRegistrationUseCase_Run_Pending_ThenWithdrawn(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStatePending, nil)
	discovery.EXPECT().WatchRegistrationApproval(mock.Anything, testRskAddress).
		Return(blockchain.RegistrationStateWithdrawn, nil).Once()

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	discovery.AssertNotCalled(t, "RegisterProvider")
	discovery.AssertNotCalled(t, "GetProvider")
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, usecases.RegistrationWithdrawnError)
}

func TestRegistrationUseCase_Run_Pending_WatchError(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStatePending, nil)
	discovery.EXPECT().WatchRegistrationApproval(mock.Anything, testRskAddress).
		Return(blockchain.RegistrationStateNone, assert.AnError).Once()

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	assert.Equal(t, int64(0), id)
	require.Error(t, err)
}

func TestRegistrationUseCase_Run_Pending_CtxCancelled(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStatePending, nil)
	discovery.EXPECT().WatchRegistrationApproval(mock.Anything, testRskAddress).
		Return(blockchain.RegistrationStateNone, context.Canceled).Once()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(ctx, blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, context.Canceled)
}

// Boot state: None — full register-then-wait flow.

func TestRegistrationUseCase_Run_None_RegisterAndWait(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateNone, nil)
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
	discovery.EXPECT().WatchRegistrationApproval(mock.Anything, testRskAddress).
		Return(blockchain.RegistrationStateApproved, nil).Once()
	discovery.EXPECT().GetProvider(testRskAddress).Return(lp.RegisteredLiquidityProvider{
		Id: 1, Address: testRskAddress, Name: test.AnyString,
		ApiBaseUrl: test.AnyUrl, Status: true, ProviderType: lp.FullProvider,
	}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertExpectations(t)
	collateral.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddCollateral")
	collateral.AssertNotCalled(t, "AddPegoutCollateral")
	assert.Equal(t, int64(1), id)
	require.NoError(t, err)
}

func TestRegistrationUseCase_Run_None_PausedDiscovery(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateNone, nil)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
	discovery.EXPECT().GetAddress().Return("discovery-contract")

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertNotCalled(t, "RegisterProvider")
	discovery.AssertNotCalled(t, "WatchRegistrationApproval")
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, blockchain.ContractPausedError)
}

func TestRegistrationUseCase_Run_None_PausedCollateral(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateNone, nil)
	discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: true, Since: 5, Reason: "test"}, nil)
	collateral.EXPECT().GetAddress().Return("collateral-contract")

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertNotCalled(t, "RegisterProvider")
	discovery.AssertNotCalled(t, "WatchRegistrationApproval")
	assert.Equal(t, int64(0), id)
	require.ErrorIs(t, err, blockchain.ContractPausedError)
}

// Approved and Pending boot states must NOT trigger the pause check (only the None path runs the register flow).
func TestRegistrationUseCase_Run_ApprovedOrPending_NoPauseCheck(t *testing.T) {
	for _, state := range []blockchain.RegistrationState{blockchain.RegistrationStateApproved, blockchain.RegistrationStatePending} {
		discovery := new(mocks.DiscoveryContractMock)
		collateral := new(mocks.CollateralManagementContractMock)
		provider := &mocks.ProviderMock{}

		provider.On("RskAddress").Return(testRskAddress)
		discovery.EXPECT().GetRegistrationState(testRskAddress).Return(state, nil)
		if state == blockchain.RegistrationStatePending {
			discovery.EXPECT().WatchRegistrationApproval(mock.Anything, testRskAddress).
				Return(blockchain.RegistrationStateApproved, nil).Once()
		}
		discovery.On("GetProvider", testRskAddress).Return(lp.RegisteredLiquidityProvider{Id: 1}, nil)

		contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
		useCase := newRegistrationUseCase(contracts, provider)
		_, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))
		require.NoError(t, err)

		collateral.AssertNotCalled(t, "PausedStatus")
		discovery.AssertNotCalled(t, "PausedStatus")
	}
}

func TestRegistrationUseCase_Run_None_ValidateParams(t *testing.T) {
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
		discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateNone, nil)
		discovery.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
		collateral.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)

		contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
		useCase := newRegistrationUseCase(contracts, provider)
		id, err := useCase.Run(context.Background(), params)

		discovery.AssertNotCalled(t, "RegisterProvider")
		discovery.AssertNotCalled(t, "WatchRegistrationApproval")
		assert.Equal(t, int64(0), id)
		require.Error(t, err)
	}
}

func TestRegistrationUseCase_Run_None_AddPeginCollateral(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateNone, nil)
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
	discovery.EXPECT().WatchRegistrationApproval(mock.Anything, testRskAddress).
		Return(blockchain.RegistrationStateApproved, nil).Once()
	discovery.On("GetProvider", testRskAddress).Return(lp.RegisteredLiquidityProvider{Id: 1}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	_, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.PeginProvider))

	collateral.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddPegoutCollateral")
	require.NoError(t, err)
}

func TestRegistrationUseCase_Run_None_AddPegoutCollateral(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	collateral := new(mocks.CollateralManagementContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateNone, nil)
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
	discovery.EXPECT().WatchRegistrationApproval(mock.Anything, testRskAddress).
		Return(blockchain.RegistrationStateApproved, nil).Once()
	discovery.On("GetProvider", testRskAddress).Return(lp.RegisteredLiquidityProvider{Id: 1}, nil)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: collateral}
	useCase := newRegistrationUseCase(contracts, provider)
	_, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.PegoutProvider))

	collateral.AssertExpectations(t)
	collateral.AssertNotCalled(t, "AddCollateral")
	require.NoError(t, err)
}

// GetRegistrationState error at boot.

func TestRegistrationUseCase_Run_GetRegistrationStateError(t *testing.T) {
	discovery := new(mocks.DiscoveryContractMock)
	provider := &mocks.ProviderMock{}

	provider.On("RskAddress").Return(testRskAddress)
	discovery.EXPECT().GetRegistrationState(testRskAddress).Return(blockchain.RegistrationStateNone, assert.AnError)

	contracts := blockchain.RskContracts{Discovery: discovery, CollateralManagement: new(mocks.CollateralManagementContractMock)}
	useCase := newRegistrationUseCase(contracts, provider)
	id, err := useCase.Run(context.Background(), blockchain.NewProviderRegistrationParams("name", test.AnyUrl, true, lp.FullProvider))

	discovery.AssertNotCalled(t, "RegisterProvider")
	discovery.AssertNotCalled(t, "WatchRegistrationApproval")
	assert.Equal(t, int64(0), id)
	require.Error(t, err)
}
