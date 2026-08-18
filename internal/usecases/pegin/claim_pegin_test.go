package pegin_test

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	claimDepositTxID  = "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899"
	claimBtcAddress   = "bcrt1qclaimaddress0000000000000000000000000"
	claimRskTxHash    = "0xclaimsubmit"
	claimRskBlockHash = "0xclaimblock"
)

type memoryClaimRepo struct {
	mu     sync.Mutex
	byKey  map[string]rootstock.PegInClaim
	insert int
}

func newMemoryClaimRepo(existing ...rootstock.PegInClaim) *memoryClaimRepo {
	repo := &memoryClaimRepo{byKey: map[string]rootstock.PegInClaim{}}
	for _, claim := range existing {
		repo.byKey[claimKey(claim.RskAddress, claim.DepositTxID)] = claim
	}
	return repo
}

func claimKey(rskAddress, depositTxID string) string {
	return rskAddress + "|" + depositTxID
}

func (repo *memoryClaimRepo) Insert(_ context.Context, claim rootstock.PegInClaim) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	key := claimKey(claim.RskAddress, claim.DepositTxID)
	if _, exists := repo.byKey[key]; exists {
		return rootstock.ErrPegInClaimAlreadyExists
	}
	repo.byKey[key] = claim
	repo.insert++
	return nil
}

func (repo *memoryClaimRepo) Get(_ context.Context, rskAddress, depositTxID string) (*rootstock.PegInClaim, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	claim, ok := repo.byKey[claimKey(rskAddress, depositTxID)]
	if !ok {
		return nil, nil
	}
	copied := claim
	return &copied, nil
}

func (repo *memoryClaimRepo) Update(_ context.Context, claim rootstock.PegInClaim) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	key := claimKey(claim.RskAddress, claim.DepositTxID)
	if _, ok := repo.byKey[key]; !ok {
		return rootstock.ErrPegInClaimNotFound
	}
	repo.byKey[key] = claim
	return nil
}

func (repo *memoryClaimRepo) ListByStates(_ context.Context, states ...rootstock.PegInClaimState) ([]rootstock.PegInClaim, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	allowed := map[rootstock.PegInClaimState]struct{}{}
	for _, state := range states {
		allowed[state] = struct{}{}
	}
	result := make([]rootstock.PegInClaim, 0)
	for _, claim := range repo.byKey {
		if _, ok := allowed[claim.State]; ok || len(states) == 0 {
			result = append(result, claim)
		}
	}
	return result, nil
}

func (repo *memoryClaimRepo) stored() rootstock.PegInClaim {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	return repo.byKey[claimKey(test.AnyRskAddress, claimDepositTxID)]
}

type claimLiquidityMock struct {
	mock.Mock
}

func (m *claimLiquidityMock) HasClaimLiquidity(
	ctx context.Context,
	requiredLiquidity *entities.Wei,
	inFlightReserved *entities.Wei,
) error {
	args := m.Called(ctx, requiredLiquidity, inFlightReserved)
	return args.Error(0)
}

func matchWei(want *entities.Wei) interface{} {
	return mock.MatchedBy(func(got *entities.Wei) bool {
		return got != nil && got.Cmp(want) == 0
	})
}

type claimHarness struct {
	repo      *memoryClaimRepo
	pegin     *mocks.PeginContractMock
	configs   *mocks.FlyoverConfigurationsContractMock
	btc       *mocks.BtcRpcMock
	rsk       *mocks.RootstockRpcServerMock
	liquidity *claimLiquidityMock
	useCase   *pegin.ClaimPegInUseCase
	amount    *entities.Wei
	fee       *entities.Wei
	rawTx     []byte
	block     blockchain.BitcoinBlockInformation
	merkle    blockchain.MerkleBranch
	entry     rootstock.PegInWatch
}

func newClaimHarness(t *testing.T, repo *memoryClaimRepo) *claimHarness {
	t.Helper()
	amount := entities.NewWei(1_000_000_000_000_000_000)
	fee := entities.NewWei(1_000_000_000_000_000)
	harness := &claimHarness{
		repo:      repo,
		pegin:     new(mocks.PeginContractMock),
		configs:   new(mocks.FlyoverConfigurationsContractMock),
		btc:       new(mocks.BtcRpcMock),
		rsk:       new(mocks.RootstockRpcServerMock),
		liquidity: new(claimLiquidityMock),
		amount:    amount,
		fee:       fee,
		rawTx:     []byte{1, 2, 3, 4, 5},
		block: blockchain.BitcoinBlockInformation{
			Hash:   [32]byte{9, 8, 7, 6, 5, 4, 3, 2, 1},
			Height: big.NewInt(500),
		},
		merkle: blockchain.MerkleBranch{
			Path:   big.NewInt(3),
			Hashes: [][32]byte{{11}, {12}},
		},
		entry: rootstock.PegInWatch{
			RskAddress: test.AnyRskAddress,
			BtcAddress: claimBtcAddress,
			State:      rootstock.PegInWatchImported,
		},
	}
	harness.useCase = pegin.NewClaimPegInUseCase(
		harness.repo,
		blockchain.RskContracts{PegIn: harness.pegin, FlyoverConfigurations: harness.configs},
		blockchain.Rpc{Btc: harness.btc, Rsk: harness.rsk},
		harness.liquidity,
		&sync.Mutex{},
		2,
	)
	return harness
}

func (h *claimHarness) payingTx(confirmations uint64) blockchain.BitcoinTransactionInformation {
	return blockchain.BitcoinTransactionInformation{
		Hash:          claimDepositTxID,
		Confirmations: confirmations,
		Outputs: map[string][]*entities.Wei{
			claimBtcAddress: {h.amount.Copy(), entities.NewWei(5)},
		},
	}
}

func (h *claimHarness) expectRefetch(confirmations uint64) {
	h.btc.On("GetTransactionInfo", claimDepositTxID).Return(h.payingTx(confirmations), nil).Once()
	h.configs.On("GetRequiredPegInBtcConfirmations", matchWei(h.amount)).Return(uint64(6), nil).Once()
}

func TestClaimPegInUseCase_BelowConfirmationsMakesZeroContractCalls(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	harness.expectRefetch(2)

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	assert.Equal(t, 0, harness.repo.insert)
	assert.Empty(t, harness.repo.byKey)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.pegin.AssertNotCalled(t, "IsHardPaused")
	harness.liquidity.AssertNotCalled(t, "HasClaimLiquidity", mock.Anything, mock.Anything, mock.Anything)
	harness.btc.AssertNotCalled(t, "GetRawTransaction", mock.Anything)
}

func TestClaimPegInUseCase_WinPersistsClaimedAndAdapterArgs(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	harness.expectRefetch(10)
	harness.pegin.On("IsHardPaused").Return(false, nil).Once()
	required := new(entities.Wei).Sub(harness.amount.Copy(), harness.fee.Copy())
	harness.liquidity.On(
		"HasClaimLiquidity",
		mock.Anything,
		matchWei(required),
		matchWei(entities.NewWei(0)),
	).Return(nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.btc.On("GetRawTransaction", claimDepositTxID).Return(harness.rawTx, nil).Once()
	harness.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(harness.block, nil).Once()
	harness.btc.On("BuildMerkleBranch", claimDepositTxID).Return(harness.merkle, nil).Once()

	pegInID := [32]byte{0xaa, 0xbb, 0xcc}
	harness.pegin.On("RequestPegIn", mock.MatchedBy(func(params blockchain.RequestPegInParams) bool {
		return params.RskAddress == test.AnyRskAddress &&
			assert.ObjectsAreEqual(harness.rawTx, params.BitcoinRawTx) &&
			params.BtcBlockHash == harness.block.Hash &&
			params.MerkleBranchPath.Cmp(harness.merkle.Path) == 0 &&
			len(params.MerkleBranchHashes) == 2 &&
			params.Amount.Cmp(harness.amount) == 0 &&
			params.Fee.Cmp(harness.fee) == 0
	})).Return(blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{TransactionHash: claimRskTxHash, BlockNumber: 100},
		Event:   &blockchain.PegInRequestedEvent{PegInId: pegInID, RskAddress: test.AnyRskAddress},
	}, nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(102), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)

	stored := harness.repo.stored()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	assert.Equal(t, hex.EncodeToString(pegInID[:]), stored.PegInID)
	assert.Equal(t, claimRskTxHash, stored.TxHash)
	assert.Equal(t, "0", stored.ReservedWei.String())
	harness.pegin.AssertExpectations(t)
}

func TestClaimPegInUseCase_PegInAlreadyProcessedIsQuietRaceLost(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	harness.expectRefetch(10)
	harness.pegin.On("IsHardPaused").Return(false, nil).Once()
	required := new(entities.Wei).Sub(harness.amount.Copy(), harness.fee.Copy())
	harness.liquidity.On("HasClaimLiquidity", mock.Anything, matchWei(required), matchWei(entities.NewWei(0))).
		Return(nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.btc.On("GetRawTransaction", claimDepositTxID).Return(harness.rawTx, nil).Once()
	harness.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(harness.block, nil).Once()
	harness.btc.On("BuildMerkleBranch", claimDepositTxID).Return(harness.merkle, nil).Once()
	harness.pegin.On("RequestPegIn", mock.Anything).Return(
		blockchain.RequestPegInResult{},
		blockchain.ErrPegInAlreadyProcessed,
	).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)

	stored := harness.repo.stored()
	assert.Equal(t, rootstock.PegInClaimRaceLost, stored.State)
	assert.Equal(t, "0", stored.ReservedWei.String())
	assert.Empty(t, stored.PegInID)
}

func TestClaimPegInUseCase_TypedFailuresAreRetryableNotClaimed(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{name: "AddressNotRegistered", err: blockchain.ErrAddressNotRegistered},
		{name: "IncorrectFronting", err: blockchain.ErrIncorrectFronting},
		{name: "InsufficientConfirmations", err: blockchain.ErrInsufficientConfirmations},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			harness := newClaimHarness(t, newMemoryClaimRepo())
			harness.expectRefetch(10)
			harness.pegin.On("IsHardPaused").Return(false, nil).Once()
			required := new(entities.Wei).Sub(harness.amount.Copy(), harness.fee.Copy())
			harness.liquidity.On("HasClaimLiquidity", mock.Anything, matchWei(required), matchWei(entities.NewWei(0))).
				Return(nil).Once()
			harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
			harness.btc.On("GetRawTransaction", claimDepositTxID).Return(harness.rawTx, nil).Once()
			harness.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(harness.block, nil).Once()
			harness.btc.On("BuildMerkleBranch", claimDepositTxID).Return(harness.merkle, nil).Once()
			harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{}, tc.err).Once()

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			require.Error(t, err)

			stored := harness.repo.stored()
			assert.Equal(t, rootstock.PegInClaimRetryableFailure, stored.State)
			assert.NotEqual(t, rootstock.PegInClaimClaimed, stored.State)
			assert.Equal(t, "0", stored.ReservedWei.String())
		})
	}
}

func TestClaimPegInUseCase_InsufficientLiquidityDoesNotSubmit(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	harness.expectRefetch(10)
	harness.pegin.On("IsHardPaused").Return(false, nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	required := new(entities.Wei).Sub(harness.amount.Copy(), harness.fee.Copy())
	harness.liquidity.On("HasClaimLiquidity", mock.Anything, matchWei(required), matchWei(entities.NewWei(0))).
		Return(&usecases.InsufficientLiquidityError{
			Available: entities.NewWei(1),
			Required:  required,
		}).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	assert.Equal(t, 0, harness.repo.insert)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.btc.AssertNotCalled(t, "GetRawTransaction", mock.Anything)
}

func TestClaimPegInUseCase_HardPauseDoesNotSubmit(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	harness.expectRefetch(10)
	harness.pegin.On("IsHardPaused").Return(true, nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	assert.Equal(t, 0, harness.repo.insert)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.liquidity.AssertNotCalled(t, "HasClaimLiquidity", mock.Anything, mock.Anything, mock.Anything)
}

func TestClaimPegInUseCase_TerminalStatesAreNoOp(t *testing.T) {
	for _, state := range []rootstock.PegInClaimState{
		rootstock.PegInClaimClaimed,
		rootstock.PegInClaimRaceLost,
	} {
		t.Run(string(state), func(t *testing.T) {
			existing := rootstock.PegInClaim{
				RskAddress:  test.AnyRskAddress,
				DepositTxID: claimDepositTxID,
				BtcAddress:  claimBtcAddress,
				State:       state,
				TxHash:      claimRskTxHash,
				PegInID:     "already",
				ReservedWei: entities.NewWei(9),
			}
			harness := newClaimHarness(t, newMemoryClaimRepo(existing))

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			require.NoError(t, err)
			harness.btc.AssertNotCalled(t, "GetTransactionInfo", mock.Anything)
			harness.rsk.AssertNotCalled(t, "GetTransactionReceipt", mock.Anything, mock.Anything)
			harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
			assert.Equal(t, state, harness.repo.stored().State)
		})
	}
}

func submittingClaim() rootstock.PegInClaim {
	return rootstock.PegInClaim{
		RskAddress:  test.AnyRskAddress,
		DepositTxID: claimDepositTxID,
		BtcAddress:  claimBtcAddress,
		State:       rootstock.PegInClaimSubmitting,
		TxHash:      claimRskTxHash,
		ReservedWei: entities.NewWei(9),
	}
}

func successReceipt() blockchain.TransactionReceipt {
	return blockchain.TransactionReceipt{
		TransactionHash: claimRskTxHash,
		BlockHash:       claimRskBlockHash,
		BlockNumber:     100,
		Status:          blockchain.SuccessfulTxStatus,
	}
}

func (h *claimHarness) expectCanonicalSuccess(event *blockchain.PegInRequestedEvent, height uint64) {
	receipt := successReceipt()
	h.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()
	h.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
	h.pegin.On("UnpackPegInRequested", receipt).Return(event, nil).Once()
	h.rsk.On("GetHeight", mock.Anything).Return(height, nil).Once()
}

func TestClaimPegInUseCase_CrashAfterMinePersistsClaimedWithoutSecondBroadcast(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	pegInID := [32]byte{0xca, 0xfe}
	harness.expectCanonicalSuccess(&blockchain.PegInRequestedEvent{
		PegInId:    pegInID,
		RskAddress: test.AnyRskAddress,
	}, 102)

	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.NoError(t, err)

	stored := harness.repo.stored()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	assert.Equal(t, hex.EncodeToString(pegInID[:]), stored.PegInID)
	assert.Equal(t, "0", stored.ReservedWei.String())
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)

	err = harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	assert.Equal(t, rootstock.PegInClaimClaimed, harness.repo.stored().State)
}

func TestClaimPegInUseCase_MissingReceiptDoesNotResubmit(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(blockchain.TransactionReceipt{}, blockchain.ErrTransactionReceiptNotFound).Once()

	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.NoError(t, err)

	stored := harness.repo.stored()
	assert.Equal(t, rootstock.PegInClaimSubmitting, stored.State)
	assert.Equal(t, claimRskTxHash, stored.TxHash)
	assert.Equal(t, "9", stored.ReservedWei.String())
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.rsk.AssertExpectations(t)
}

func TestClaimPegInUseCase_ReorgDroppedReceiptStaysSubmitting(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(successReceipt(), nil).Once()
	harness.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).
		Return(blockchain.BlockInfo{}, errors.New("not found")).Once()

	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.NoError(t, err)

	assert.Equal(t, rootstock.PegInClaimSubmitting, harness.repo.stored().State)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.pegin.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
}

func TestClaimPegInUseCase_StatusZeroIdentifiesViaPreflightNotReceiptData(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	reverted := successReceipt()
	reverted.Status = 0
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(reverted, nil).Once()
	harness.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).
		Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.btc.On("GetRawTransaction", claimDepositTxID).Return(harness.rawTx, nil).Once()
	harness.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(harness.block, nil).Once()
	harness.btc.On("BuildMerkleBranch", claimDepositTxID).Return(harness.merkle, nil).Once()
	harness.pegin.On("IdentifyRequestPegIn", mock.Anything).Return(blockchain.ErrPegInAlreadyProcessed).Once()

	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.NoError(t, err)

	stored := harness.repo.stored()
	assert.Equal(t, rootstock.PegInClaimRaceLost, stored.State)
	assert.Equal(t, "0", stored.ReservedWei.String())
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.pegin.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
}
