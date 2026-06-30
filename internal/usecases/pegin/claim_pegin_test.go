package pegin_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	peginEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	testRskAddress = "0x892813507Bf3aBF2890759d2135Ec34f4909Fea5"
	testBtcTxHash  = "0a1b2c3d4e5f"
)

func newDeposit(amount int64, confirmations uint64) peginEntity.ConfirmedDeposit {
	return peginEntity.ConfirmedDeposit{
		RskAddress:    testRskAddress,
		BtcTxHash:     testBtcTxHash,
		BtcTxHashRaw:  [32]byte{0x0a, 0x1b, 0x2c, 0x3d, 0x4e, 0x5f},
		Amount:        entities.NewWei(amount),
		Confirmations: confirmations,
	}
}

// claimMocks bundles the mocks a ClaimPegInUseCase needs.
type claimMocks struct {
	peginContract  *mocks.PeginContractMock
	registry       *mocks.PegInAddressRegistryContractMock
	configurations *mocks.FlyoverConfigurationsContractMock
	rsk            *mocks.RootstockRpcServerMock
	btc            *mocks.BtcRpcMock
	lp             *mocks.ProviderMock
	mutex          *mocks.MutexMock
}

func newClaimUseCase(m claimMocks) *pegin.ClaimPegInUseCase {
	contracts := blockchain.RskContracts{
		PegIn:                 m.peginContract,
		PegInAddressRegistry:  m.registry,
		FlyoverConfigurations: m.configurations,
	}
	rpc := blockchain.Rpc{Btc: m.btc, Rsk: m.rsk}
	return pegin.NewClaimPegInUseCase(contracts, rpc, m.lp, m.lp, m.mutex)
}

func baseMocks() claimMocks {
	peginContract := new(mocks.PeginContractMock)
	peginContract.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil).Maybe()
	peginContract.EXPECT().GetAddress().Return("pegin").Maybe()
	mutex := new(mocks.MutexMock)
	mutex.On("Lock").Return().Maybe()
	mutex.On("Unlock").Return().Maybe()
	lp := new(mocks.ProviderMock)
	lp.On("RskAddress").Return(testRskAddress).Maybe()
	return claimMocks{
		peginContract:  peginContract,
		registry:       new(mocks.PegInAddressRegistryContractMock),
		configurations: new(mocks.FlyoverConfigurationsContractMock),
		rsk:            new(mocks.RootstockRpcServerMock),
		btc:            new(mocks.BtcRpcMock),
		lp:             lp,
		mutex:          mutex,
	}
}

func TestClaimPegInUseCase_Request_FrontsRbtcAndWins(t *testing.T) {
	m := baseMocks()
	deposit := newDeposit(1_000_000, 10)
	m.registry.EXPECT().IsRegistered(testRskAddress).Return(true, nil).Once()
	m.configurations.EXPECT().CalculatePegInFee(deposit.Amount).Return(entities.NewWei(10_000), nil).Once()
	m.rsk.EXPECT().GetBalance(mock.Anything, mock.Anything).Return(entities.NewWei(5_000_000), nil).Once()
	// The LP must front amount - fee = 990_000.
	m.peginContract.EXPECT().RequestPegIn(mock.MatchedBy(func(p blockchain.RequestPegInParams) bool {
		return p.RskAddress == testRskAddress && p.Value.Cmp(entities.NewWei(990_000)) == 0 && p.Amount.Cmp(deposit.Amount) == 0
	})).Return(blockchain.TransactionReceipt{TransactionHash: "0xreq"}, nil).Once()

	state, err := newClaimUseCase(m).Request(context.Background(), deposit)
	require.NoError(t, err)
	assert.Equal(t, peginEntity.WatchedAddressStateRequested, state)
	m.peginContract.AssertExpectations(t)
	m.configurations.AssertExpectations(t)
}

func TestClaimPegInUseCase_Request_TolerateAlreadyClaimed(t *testing.T) {
	m := baseMocks()
	deposit := newDeposit(1_000_000, 10)
	m.registry.EXPECT().IsRegistered(testRskAddress).Return(true, nil).Once()
	m.configurations.EXPECT().CalculatePegInFee(deposit.Amount).Return(entities.NewWei(10_000), nil).Once()
	m.rsk.EXPECT().GetBalance(mock.Anything, mock.Anything).Return(entities.NewWei(5_000_000), nil).Once()
	m.peginContract.EXPECT().RequestPegIn(mock.Anything).Return(blockchain.TransactionReceipt{}, blockchain.AlreadyClaimedError).Once()

	// First-mined-wins: the revert is benign, no error and state becomes ClaimedByOther.
	state, err := newClaimUseCase(m).Request(context.Background(), deposit)
	require.NoError(t, err)
	assert.Equal(t, peginEntity.WatchedAddressStateClaimedByOther, state)
}

func TestClaimPegInUseCase_Request_RejectsUnregistered(t *testing.T) {
	m := baseMocks()
	deposit := newDeposit(1_000_000, 10)
	m.registry.EXPECT().IsRegistered(testRskAddress).Return(false, nil).Once()

	state, err := newClaimUseCase(m).Request(context.Background(), deposit)
	require.Error(t, err)
	assert.ErrorIs(t, err, blockchain.AddressNotRegisteredError)
	assert.Equal(t, peginEntity.WatchedAddressStateWaitingForDeposit, state)
	// No claim attempted.
	m.peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_Request_RejectsInsufficientOwnLiquidity(t *testing.T) {
	m := baseMocks()
	deposit := newDeposit(1_000_000, 10)
	m.registry.EXPECT().IsRegistered(testRskAddress).Return(true, nil).Once()
	m.configurations.EXPECT().CalculatePegInFee(deposit.Amount).Return(entities.NewWei(10_000), nil).Once()
	// Wallet cannot cover the 990_000 to front.
	m.rsk.EXPECT().GetBalance(mock.Anything, mock.Anything).Return(entities.NewWei(100), nil).Once()

	state, err := newClaimUseCase(m).Request(context.Background(), deposit)
	require.Error(t, err)
	assert.ErrorIs(t, err, usecases.NoLiquidityError)
	assert.Equal(t, peginEntity.WatchedAddressStateWaitingForDeposit, state)
	m.peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_Resolve_WaitsForConfirmations(t *testing.T) {
	m := baseMocks()
	deposit := newDeposit(1_000_000, 1) // only 1 confirmation
	m.configurations.EXPECT().GetRequiredPegInConfirmations(deposit.Amount).Return(uint64(10), nil).Once()

	state, err := newClaimUseCase(m).Resolve(context.Background(), deposit)
	require.Error(t, err)
	assert.ErrorIs(t, err, blockchain.WaitingForBridgeError)
	assert.Equal(t, peginEntity.WatchedAddressStateRequested, state)
	m.peginContract.AssertNotCalled(t, "ResolvePegIn", mock.Anything)
}

func TestClaimPegInUseCase_Resolve_SettlesWhenConfirmed(t *testing.T) {
	m := baseMocks()
	deposit := newDeposit(1_000_000, 20)
	m.configurations.EXPECT().GetRequiredPegInConfirmations(deposit.Amount).Return(uint64(10), nil).Once()
	m.btc.On("GetRawTransaction", testBtcTxHash).Return([]byte{0x01, 0x02}, nil).Once()
	m.btc.On("GetPartialMerkleTree", testBtcTxHash).Return([]byte{0x03, 0x04}, nil).Once()
	m.btc.On("GetTransactionBlockInfo", testBtcTxHash).Return(blockchain.BitcoinBlockInformation{Height: big.NewInt(500), Time: time.Now()}, nil).Once()
	m.lp.On("RskAddress").Return(testRskAddress).Once()
	m.peginContract.EXPECT().ResolvePegIn(mock.MatchedBy(func(p blockchain.ResolvePegInParams) bool {
		return p.RskAddress == testRskAddress && p.Height.Cmp(big.NewInt(500)) == 0 && len(p.BtcRawTransaction) == 2
	})).Return(blockchain.TransactionReceipt{TransactionHash: "0xresolve"}, nil).Once()

	state, err := newClaimUseCase(m).Resolve(context.Background(), deposit)
	require.NoError(t, err)
	assert.Equal(t, peginEntity.WatchedAddressStateResolved, state)
	m.peginContract.AssertExpectations(t)
}

func TestClaimPegInUseCase_Resolve_BridgeRetryMapsToWaiting(t *testing.T) {
	m := baseMocks()
	deposit := newDeposit(1_000_000, 20)
	m.configurations.EXPECT().GetRequiredPegInConfirmations(deposit.Amount).Return(uint64(10), nil).Once()
	m.btc.On("GetRawTransaction", testBtcTxHash).Return([]byte{0x01}, nil).Once()
	m.btc.On("GetPartialMerkleTree", testBtcTxHash).Return([]byte{0x02}, nil).Once()
	m.btc.On("GetTransactionBlockInfo", testBtcTxHash).Return(blockchain.BitcoinBlockInformation{Height: big.NewInt(1)}, nil).Once()
	m.lp.On("RskAddress").Return(testRskAddress).Once()
	m.peginContract.EXPECT().ResolvePegIn(mock.Anything).Return(blockchain.TransactionReceipt{}, blockchain.WaitingForBridgeError).Once()

	state, err := newClaimUseCase(m).Resolve(context.Background(), deposit)
	require.Error(t, err)
	assert.ErrorIs(t, err, blockchain.WaitingForBridgeError)
	assert.Equal(t, peginEntity.WatchedAddressStateRequested, state)
}

func TestClaimPegInUseCase_NotConfigured(t *testing.T) {
	rpc := blockchain.Rpc{Btc: new(mocks.BtcRpcMock), Rsk: new(mocks.RootstockRpcServerMock)}
	// No registry / configurations contracts wired (legacy deployment).
	useCase := pegin.NewClaimPegInUseCase(blockchain.RskContracts{PegIn: new(mocks.PeginContractMock)}, rpc, new(mocks.ProviderMock), new(mocks.ProviderMock), new(mocks.MutexMock))
	state, err := useCase.Request(context.Background(), newDeposit(1_000_000, 10))
	require.Error(t, err)
	assert.Equal(t, peginEntity.WatchedAddressStateWaitingForDeposit, state)
}
