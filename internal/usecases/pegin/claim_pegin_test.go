package pegin_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"sync"
	"sync/atomic"
	"testing"
	"time"

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
	mu    sync.Mutex
	byKey map[string]rootstock.PegInClaim
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

func claimAccount(t *testing.T) *mocks.ProviderMock {
	t.Helper()
	account := mocks.NewProviderMock(t)
	account.On("RskAddress").Return(test.AnyRskAddress).Maybe()
	return account
}

func matchWei(want *entities.Wei) interface{} {
	return mock.MatchedBy(func(got *entities.Wei) bool {
		return got != nil && got.Cmp(want) == 0
	})
}

const claimEstimatedGas uint64 = 100000

type claimHarness struct {
	repo         rootstock.PegInClaimRepository
	pegin        *mocks.PeginContractMock
	pause        *mocks.PauseRegistryContractMock
	configs      *mocks.FlyoverConfigurationsContractMock
	btc          *mocks.BtcRpcMock
	rsk          *mocks.RootstockRpcServerMock
	useCase      *pegin.ClaimPegInUseCase
	amount       *entities.Wei
	fee          *entities.Wei
	gasPrice     *entities.Wei
	estimatedGas uint64
	rawTx        []byte
	block        blockchain.BitcoinBlockInformation
	merkle       blockchain.MerkleBranch
	entry        rootstock.PegInWatch
}

func newClaimHarness(t *testing.T, repo rootstock.PegInClaimRepository) *claimHarness {
	t.Helper()
	amount := entities.NewWei(1_000_000_000_000_000_000)
	fee := entities.NewWei(1_000_000_000_000_000)
	harness := &claimHarness{
		repo:         repo,
		pegin:        mocks.NewPeginContractMock(t),
		pause:        mocks.NewPauseRegistryContractMock(t),
		configs:      mocks.NewFlyoverConfigurationsContractMock(t),
		btc:          mocks.NewBtcRpcMock(t),
		rsk:          mocks.NewRootstockRpcServerMock(t),
		amount:       amount,
		fee:          fee,
		gasPrice:     entities.NewWei(1),
		estimatedGas: claimEstimatedGas,
		rawTx:        []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0xff, 0xaa, 0xbb},
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
		blockchain.RskContracts{
			PegIn:                 harness.pegin,
			FlyoverConfigurations: harness.configs,
			PauseRegistry:         harness.pause,
		},
		blockchain.Rpc{Btc: harness.btc, Rsk: harness.rsk},
		claimAccount(t),
		&sync.Mutex{},
		2,
	)
	return harness
}

func (h *claimHarness) payingTx(confirmations uint64) blockchain.BitcoinTransactionInformation {
	return h.payingTxFor(claimDepositTxID, confirmations)
}

func (h *claimHarness) payingTxFor(depositTxID string, confirmations uint64) blockchain.BitcoinTransactionInformation {
	return blockchain.BitcoinTransactionInformation{
		Hash:          depositTxID,
		Confirmations: confirmations,
		Outputs: map[string][]*entities.Wei{
			claimBtcAddress: {h.amount.Copy(), entities.NewWei(5)},
		},
	}
}

func expectNoExistingClaim(claims *mocks.PegInClaimRepositoryMock) {
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return((*rootstock.PegInClaim)(nil), nil)
}

func expectExistingClaim(claims *mocks.PegInClaimRepositoryMock, claim rootstock.PegInClaim) {
	copied := claim
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).Return(&copied, nil)
}

func expectEmptyInFlight(claims *mocks.PegInClaimRepositoryMock) {
	claims.On("ListByStates", mock.Anything, rootstock.PegInClaimCandidate, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{}, nil).Once()
}

func expectInsertOk(claims *mocks.PegInClaimRepositoryMock) {
	claims.On("Insert", mock.Anything, mock.Anything).Return(nil).Once()
}

func expectUpdatesOk(claims *mocks.PegInClaimRepositoryMock) {
	claims.On("Update", mock.Anything, mock.Anything).Return(nil)
}

func matchSubmittingEmptyHash() interface{} {
	return mock.MatchedBy(func(claim rootstock.PegInClaim) bool {
		return claim.State == rootstock.PegInClaimSubmitting && claim.TxHash == ""
	})
}

func expectMarkerUpdate(claims *mocks.PegInClaimRepositoryMock) {
	claims.On("Update", mock.Anything, matchSubmittingEmptyHash()).Return(nil).Once()
}

func submittingEmptyHashClaim() rootstock.PegInClaim {
	created := submittingClaim()
	created.TxHash = ""
	return created
}

func candidateWithReserve() rootstock.PegInClaim {
	created := submittingClaim()
	created.State = rootstock.PegInClaimCandidate
	created.TxHash = ""
	created.UpdatedAt = time.Now().UTC().Add(-time.Hour)
	return created
}

func winningRequestResult(pegInID [32]byte) blockchain.RequestPegInResult {
	return blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{
			TransactionHash: claimRskTxHash,
			BlockNumber:     100,
			Status:          blockchain.SuccessfulTxStatus,
		},
		Event: blockchain.PegInRequestedEvent{PegInId: pegInID, RskAddress: test.AnyRskAddress},
	}
}

func expectSuccessfulGates(h *claimHarness, times int) {
	h.btc.On("GetTransactionInfo", claimDepositTxID).Return(h.payingTx(10), nil).Times(times)
	h.configs.On("GetRequiredPegInBtcConfirmations", matchWei(h.amount)).Return(uint64(6), nil).Times(times)
	h.pause.On("PauseLevel").Return(blockchain.PauseLevelNone, nil).Times(times)
	h.configs.On("CalculatePegInFee", matchWei(h.amount)).Return(h.fee.Copy(), nil).Times(times)
	h.btc.On("GetRawTransaction", claimDepositTxID).Return(h.rawTx, nil).Times(times)
	h.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(h.block, nil).Times(times)
	h.btc.On("BuildMerkleBranch", claimDepositTxID).Return(h.merkle, nil).Times(times)
	h.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).
		Return(h.spendableRequired(entities.NewWei(0), h.fee), nil).Times(times)
	h.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(h.estimatedGas, nil).Times(times)
	h.rsk.On("GasPrice", mock.Anything).Return(h.gasPrice.Copy(), nil).Times(times)
}

func (h *claimHarness) expectRefetch(confirmations uint64) {
	h.btc.On("GetTransactionInfo", claimDepositTxID).Return(h.payingTx(confirmations), nil).Once()
	h.configs.On("GetRequiredPegInBtcConfirmations", matchWei(h.amount)).Return(uint64(6), nil).Once()
}

func (h *claimHarness) gasCost() *entities.Wei {
	return new(entities.Wei).Mul(h.gasPrice.Copy(), entities.NewUWei(h.estimatedGas))
}

func (h *claimHarness) payable(fee *entities.Wei) *entities.Wei {
	required := new(entities.Wei).Sub(h.amount.Copy(), fee.Copy())
	if required.Cmp(entities.NewWei(0)) < 0 {
		return entities.NewWei(0)
	}
	return required
}

func (h *claimHarness) spendableRequired(inFlight, fee *entities.Wei) *entities.Wei {
	return new(entities.Wei).Add(new(entities.Wei).Add(h.payable(fee), h.gasCost()), inFlight)
}

func (h *claimHarness) expectPassingGates(inFlight *entities.Wei) {
	h.expectRefetch(10)
	h.pause.On("PauseLevel").Return(blockchain.PauseLevelNone, nil).Once()
	h.configs.On("CalculatePegInFee", matchWei(h.amount)).Return(h.fee.Copy(), nil).Once()
	h.expectBuildParams()
	h.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).
		Return(h.spendableRequired(inFlight, h.fee), nil).Once()
	h.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(h.estimatedGas, nil).Once()
	h.rsk.On("GasPrice", mock.Anything).Return(h.gasPrice.Copy(), nil).Once()
}

func (h *claimHarness) expectBuildParams() {
	h.btc.On("GetRawTransaction", claimDepositTxID).Return(h.rawTx, nil).Once()
	h.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(h.block, nil).Once()
	h.btc.On("BuildMerkleBranch", claimDepositTxID).Return(h.merkle, nil).Once()
}

func (h *claimHarness) expectPauseNone() {
	h.pause.On("PauseLevel").Return(blockchain.PauseLevelNone, nil).Once()
}

func TestClaimPegInUseCase_BelowConfirmationsMakesZeroContractCalls(t *testing.T) {
	t.Run("no existing claim", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newClaimHarness(t, claims)
		expectNoExistingClaim(claims)
		harness.expectRefetch(2)

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		claims.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
		claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		harness.pause.AssertNotCalled(t, "PauseLevel")
		harness.rsk.AssertNotCalled(t, "GetBalance", mock.Anything, mock.Anything)
		harness.btc.AssertNotCalled(t, "GetRawTransaction", mock.Anything)
	})
	t.Run("releases existing reserve", func(t *testing.T) {
		created := submittingClaim()
		created.State = rootstock.PegInClaimCandidate
		created.TxHash = ""
		repo := newMemoryClaimRepo(created)
		harness := newClaimHarness(t, repo)
		harness.expectRefetch(2)

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		assert.Equal(t, "0", repo.stored().ReservedWei.String())
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
}

func TestClaimPegInUseCase_WinPersistsClaimedAndAdapterArgs(t *testing.T) {
	repo := newMemoryClaimRepo()
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))

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
		Receipt: blockchain.TransactionReceipt{
			TransactionHash: claimRskTxHash,
			BlockNumber:     100,
			Status:          blockchain.SuccessfulTxStatus,
		},
		Event: blockchain.PegInRequestedEvent{PegInId: pegInID, RskAddress: test.AnyRskAddress},
	}, nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(102), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)

	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	assert.Equal(t, hex.EncodeToString(pegInID[:]), stored.PegInID)
	assert.Equal(t, claimRskTxHash, stored.TxHash)
	assert.Equal(t, "0", stored.ReservedWei.String())
	harness.pegin.AssertExpectations(t)
}

func TestClaimPegInUseCase_PegInAlreadyProcessedIsQuietRaceLost(t *testing.T) {
	repo := newMemoryClaimRepo()
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(
		blockchain.RequestPegInResult{},
		blockchain.ErrPegInAlreadyProcessed,
	).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)

	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimRaceLost, stored.State)
	assert.Equal(t, "0", stored.ReservedWei.String())
	assert.Empty(t, stored.PegInID)
}

func TestClaimPegInUseCase_StoredCandidateAlreadyProcessedIsRaceLost(t *testing.T) {
	created := candidateWithReserve()
	repo := newMemoryClaimRepo(created)
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(
		blockchain.RequestPegInResult{},
		blockchain.ErrPegInAlreadyProcessed,
	).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimRaceLost, stored.State)
	assert.Equal(t, "0", stored.ReservedWei.String())
	assert.Empty(t, stored.TxHash)
	harness.pegin.AssertNumberOfCalls(t, "RequestPegIn", 1)
}

func TestLogPegInClaimSubmittingEmptyTxHash(t *testing.T) {
	msg := pegin.LogPegInClaimSubmittingEmptyTxHash(test.AnyRskAddress, claimDepositTxID)
	assert.Contains(t, msg, "follow incident-recovery")
	assert.Contains(t, msg, test.AnyRskAddress)
	assert.Contains(t, msg, claimDepositTxID)
	assert.Contains(t, msg, "not resubmitting")
	assert.NotContains(t, msg, "already processed")
	assert.NotContains(t, msg, "INCIDENT")
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
			repo := newMemoryClaimRepo()
			harness := newClaimHarness(t, repo)
			harness.expectPassingGates(entities.NewWei(0))
			harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{}, tc.err).Once()

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			require.ErrorIs(t, err, tc.err)

			stored := repo.stored()
			assert.Equal(t, rootstock.PegInClaimRetryableFailure, stored.State)
			assert.NotEqual(t, rootstock.PegInClaimClaimed, stored.State)
			assert.Equal(t, "0", stored.ReservedWei.String())
		})
	}
}

func TestClaimPegInUseCase_InsufficientLiquidityDoesNotSubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	expectNoExistingClaim(claims)
	expectEmptyInFlight(claims)
	harness.expectRefetch(10)
	harness.expectPauseNone()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.expectBuildParams()
	harness.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).Return(entities.NewWei(1), nil).Once()
	harness.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(harness.estimatedGas, nil).Once()
	harness.rsk.On("GasPrice", mock.Anything).Return(harness.gasPrice.Copy(), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	claims.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_HardPauseDoesNotSubmit(t *testing.T) {
	t.Run("no existing claim", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newClaimHarness(t, claims)
		expectNoExistingClaim(claims)
		harness.expectRefetch(10)
		harness.pause.On("PauseLevel").Return(blockchain.PauseLevelHard, nil).Once()

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		claims.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		harness.rsk.AssertNotCalled(t, "GetBalance", mock.Anything, mock.Anything)
	})
	t.Run("releases existing reserve", func(t *testing.T) {
		created := submittingClaim()
		created.State = rootstock.PegInClaimCandidate
		created.TxHash = ""
		repo := newMemoryClaimRepo(created)
		harness := newClaimHarness(t, repo)
		harness.expectRefetch(10)
		harness.pause.On("PauseLevel").Return(blockchain.PauseLevelHard, nil).Once()

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		assert.Equal(t, "0", repo.stored().ReservedWei.String())
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
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
			claims := mocks.NewPegInClaimRepositoryMock(t)
			harness := newClaimHarness(t, claims)
			expectExistingClaim(claims, existing)

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			require.NoError(t, err)
			harness.btc.AssertNotCalled(t, "GetTransactionInfo", mock.Anything)
			harness.rsk.AssertNotCalled(t, "GetTransactionReceipt", mock.Anything, mock.Anything)
			harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
			claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
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

func TestClaimPegInUseCase_ZeroFirstOutputDoesNotSubmit(t *testing.T) {
	zeroOutput := func(amount *entities.Wei) blockchain.BitcoinTransactionInformation {
		return blockchain.BitcoinTransactionInformation{
			Hash:          claimDepositTxID,
			Confirmations: 10,
			Outputs:       map[string][]*entities.Wei{"other": {amount.Copy()}},
		}
	}
	t.Run("no existing claim", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newClaimHarness(t, claims)
		expectNoExistingClaim(claims)
		harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(zeroOutput(harness.amount), nil).Once()

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		claims.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		harness.pause.AssertNotCalled(t, "PauseLevel")
	})
	t.Run("releases existing reserve", func(t *testing.T) {
		created := candidateWithReserve()
		repo := newMemoryClaimRepo(created)
		harness := newClaimHarness(t, repo)
		harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(zeroOutput(harness.amount), nil).Once()

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		stored := repo.stored()
		assert.Equal(t, "0", stored.ReservedWei.String())
		assert.True(t, stored.UpdatedAt.After(created.UpdatedAt))
		assert.Equal(t, rootstock.PegInClaimCandidate, stored.State)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		harness.pause.AssertNotCalled(t, "PauseLevel")
	})
}

func TestClaimPegInUseCase_ReleaseReserveDoesNotOverwriteArmedOrTerminal(t *testing.T) {
	zeroOutput := blockchain.BitcoinTransactionInformation{
		Hash:          claimDepositTxID,
		Confirmations: 10,
		Outputs:       map[string][]*entities.Wei{"other": {entities.NewWei(1)}},
	}
	cases := []struct {
		name string
		row  rootstock.PegInClaim
	}{
		{name: "submitting with hash", row: submittingClaim()},
		{name: "submitting empty hash", row: submittingEmptyHashClaim()},
		{name: "race_lost", row: func() rootstock.PegInClaim {
			row := submittingClaim()
			row.State = rootstock.PegInClaimRaceLost
			row.TxHash = ""
			return row
		}()},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			claims := mocks.NewPegInClaimRepositoryMock(t)
			harness := newClaimHarness(t, claims)
			stale := candidateWithReserve()
			armed := tc.row
			claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
				Return(&stale, nil).Once()
			harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(zeroOutput, nil).Once()
			claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
				Return(&armed, nil).Once()

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			require.NoError(t, err)
			claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		})
	}
}

func TestClaimPegInUseCase_ExistingTxHashDoesNotSubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	expectExistingClaim(claims, submittingClaim())

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.btc.AssertNotCalled(t, "GetTransactionInfo", mock.Anything)
	harness.rsk.AssertNotCalled(t, "GetTransactionReceipt", mock.Anything, mock.Anything)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestClaimPegInUseCase_SaveAlreadySubmittedDoesNotSubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	peginContract := mocks.NewPeginContractMock(t)
	pause := mocks.NewPauseRegistryContractMock(t)
	configs := mocks.NewFlyoverConfigurationsContractMock(t)
	btc := mocks.NewBtcRpcMock(t)
	rsk := mocks.NewRootstockRpcServerMock(t)
	amount := entities.NewWei(1_000_000_000_000_000_000)
	fee := entities.NewWei(1_000_000_000_000_000)
	gasPrice := entities.NewWei(1)
	rawTx := []byte{0x01, 0x00, 0x00, 0x00, 0x01, 0xff, 0xaa, 0xbb}
	block := blockchain.BitcoinBlockInformation{Hash: [32]byte{9}, Height: big.NewInt(500)}
	merkle := blockchain.MerkleBranch{Path: big.NewInt(3), Hashes: [][32]byte{{11}}}
	entry := rootstock.PegInWatch{
		RskAddress: test.AnyRskAddress,
		BtcAddress: claimBtcAddress,
		State:      rootstock.PegInWatchImported,
	}
	stored := submittingClaim()
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return((*rootstock.PegInClaim)(nil), nil).Once()
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return(&stored, nil).Once()
	btc.On("GetTransactionInfo", claimDepositTxID).Return(blockchain.BitcoinTransactionInformation{
		Hash:          claimDepositTxID,
		Confirmations: 10,
		Outputs:       map[string][]*entities.Wei{claimBtcAddress: {amount.Copy()}},
	}, nil).Once()
	configs.On("GetRequiredPegInBtcConfirmations", matchWei(amount)).Return(uint64(6), nil).Once()
	pause.On("PauseLevel").Return(blockchain.PauseLevelNone, nil).Once()
	configs.On("CalculatePegInFee", matchWei(amount)).Return(fee.Copy(), nil).Once()
	claims.On("ListByStates", mock.Anything, rootstock.PegInClaimCandidate, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{}, nil).Once()
	btc.On("GetRawTransaction", claimDepositTxID).Return(rawTx, nil).Once()
	btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(block, nil).Once()
	btc.On("BuildMerkleBranch", claimDepositTxID).Return(merkle, nil).Once()
	rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).Return(entities.NewWei(1_000_000_000_000_000_000), nil).Once()
	peginContract.On("EstimateRequestPegInGas", mock.Anything).Return(claimEstimatedGas, nil).Once()
	rsk.On("GasPrice", mock.Anything).Return(gasPrice.Copy(), nil).Once()

	useCase := pegin.NewClaimPegInUseCase(
		claims,
		blockchain.RskContracts{
			PegIn:                 peginContract,
			FlyoverConfigurations: configs,
			PauseRegistry:         pause,
		},
		blockchain.Rpc{Btc: btc, Rsk: rsk},
		claimAccount(t),
		&sync.Mutex{},
		2,
	)
	err := useCase.Run(context.Background(), entry, claimDepositTxID)
	require.NoError(t, err)
	peginContract.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	claims.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestClaimPegInUseCase_InFlightReservedExcludesSelf(t *testing.T) {
	other := rootstock.PegInClaim{
		RskAddress:  "0xother",
		DepositTxID: "ff" + claimDepositTxID[2:],
		State:       rootstock.PegInClaimCandidate,
		ReservedWei: entities.NewWei(7),
	}
	nilReserved := rootstock.PegInClaim{
		RskAddress:  "0xnil",
		DepositTxID: "ee" + claimDepositTxID[2:],
		State:       rootstock.PegInClaimCandidate,
	}
	self := submittingClaim()
	self.State = rootstock.PegInClaimCandidate
	self.TxHash = ""
	self.ReservedWei = entities.NewWei(9)
	self.CreatedAt = time.Now().UTC().Add(-time.Hour)
	repo := newMemoryClaimRepo(other, nilReserved, self)
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(7))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{
			TransactionHash: claimRskTxHash,
			BlockNumber:     100,
			Status:          blockchain.SuccessfulTxStatus,
		},
		Event: blockchain.PegInRequestedEvent{PegInId: [32]byte{1}, RskAddress: test.AnyRskAddress},
	}, nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(102), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	assert.Equal(t, self.CreatedAt, repo.stored().CreatedAt)
	harness.pegin.AssertExpectations(t)
}

func TestClaimPegInUseCase_OracleErrorsDoNotSubmit(t *testing.T) {
	t.Run("confirmations oracle", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newClaimHarness(t, claims)
		expectNoExistingClaim(claims)
		harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
		harness.configs.On("GetRequiredPegInBtcConfirmations", matchWei(harness.amount)).Return(uint64(0), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("hard-pause oracle", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newClaimHarness(t, claims)
		expectNoExistingClaim(claims)
		harness.expectRefetch(10)
		harness.pause.On("PauseLevel").Return(uint8(0), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("fee oracle", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newClaimHarness(t, claims)
		expectNoExistingClaim(claims)
		harness.expectRefetch(10)
		harness.expectPauseNone()
		harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return((*entities.Wei)(nil), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
}

func (h *claimHarness) expectGatesBeforeSpendable() {
	h.expectRefetch(10)
	h.expectPauseNone()
	h.configs.On("CalculatePegInFee", matchWei(h.amount)).Return(h.fee.Copy(), nil).Once()
	h.expectBuildParams()
}

func runClaimWithoutSubmitOnUnavailable(t *testing.T, setup func(*claimHarness, *mocks.PegInClaimRepositoryMock)) {
	t.Helper()
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	expectNoExistingClaim(claims)
	setup(harness, claims)
	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	claims.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_WalletBalanceErrorDoesNotSubmit(t *testing.T) {
	runClaimWithoutSubmitOnUnavailable(t, func(harness *claimHarness, claims *mocks.PegInClaimRepositoryMock) {
		expectEmptyInFlight(claims)
		harness.expectGatesBeforeSpendable()
		harness.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(harness.estimatedGas, nil).Once()
		harness.rsk.On("GasPrice", mock.Anything).Return(harness.gasPrice.Copy(), nil).Once()
		harness.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).Return((*entities.Wei)(nil), assert.AnError).Once()
	})
}

func TestClaimPegInUseCase_GasEstimateErrorDoesNotSubmit(t *testing.T) {
	runClaimWithoutSubmitOnUnavailable(t, func(harness *claimHarness, claims *mocks.PegInClaimRepositoryMock) {
		expectEmptyInFlight(claims)
		harness.expectGatesBeforeSpendable()
		harness.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(uint64(0), assert.AnError).Once()
	})
}

func TestClaimPegInUseCase_GasPriceErrorDoesNotSubmit(t *testing.T) {
	runClaimWithoutSubmitOnUnavailable(t, func(harness *claimHarness, claims *mocks.PegInClaimRepositoryMock) {
		expectEmptyInFlight(claims)
		harness.expectGatesBeforeSpendable()
		harness.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(harness.estimatedGas, nil).Once()
		harness.rsk.On("GasPrice", mock.Anything).Return((*entities.Wei)(nil), assert.AnError).Once()
	})
}

func TestClaimPegInUseCase_InFlightListErrorDoesNotSubmit(t *testing.T) {
	runClaimWithoutSubmitOnUnavailable(t, func(harness *claimHarness, claims *mocks.PegInClaimRepositoryMock) {
		claims.On("ListByStates", mock.Anything, rootstock.PegInClaimCandidate, rootstock.PegInClaimSubmitting).
			Return([]rootstock.PegInClaim(nil), assert.AnError).Once()
		harness.expectGatesBeforeSpendable()
	})
}

func TestClaimPegInUseCase_FeeAboveAmountReleasesReserveAndStaysReevaluable(t *testing.T) {
	created := candidateWithReserve()
	repo := newMemoryClaimRepo(created)
	harness := newClaimHarness(t, repo)
	harness.expectRefetch(10)
	harness.expectPauseNone()
	highFee := new(entities.Wei).Add(harness.amount.Copy(), entities.NewWei(1))
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(highFee, nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	require.NotErrorIs(t, err, blockchain.ErrIncorrectFronting)
	stored := repo.stored()
	assert.Equal(t, "0", stored.ReservedWei.String())
	assert.True(t, stored.UpdatedAt.After(created.UpdatedAt))
	assert.Equal(t, rootstock.PegInClaimCandidate, stored.State)
	assert.Empty(t, stored.TxHash)
	assert.NotEqual(t, rootstock.PegInClaimRetryableFailure, stored.State)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.btc.AssertNotCalled(t, "GetRawTransaction", mock.Anything)
	harness.rsk.AssertNotCalled(t, "GetBalance", mock.Anything, mock.Anything)
}

func TestClaimPegInUseCase_RejectsWitnessSerializedTx(t *testing.T) {
	cases := []struct {
		name  string
		rawTx []byte
	}{
		{name: "short raw bytes", rawTx: []byte{1, 0, 0, 0, 1}},
		{name: "witness serialized marker", rawTx: []byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x01, 0xaa, 0xbb}},
	}
	for _, tc := range cases {
		t.Run(tc.name+" no existing claim", func(t *testing.T) {
			claims := mocks.NewPegInClaimRepositoryMock(t)
			harness := newClaimHarness(t, claims)
			expectNoExistingClaim(claims)
			harness.expectRefetch(10)
			harness.expectPauseNone()
			harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
			harness.btc.On("GetRawTransaction", claimDepositTxID).Return(tc.rawTx, nil).Once()

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			require.ErrorIs(t, err, blockchain.ErrWitnessSerializedTxNotAccepted)
			require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
			claims.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
			claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			harness.pegin.AssertNotCalled(t, "EstimateRequestPegInGas", mock.Anything)
			harness.pegin.AssertNotCalled(t, "SimulateRequestPegIn", mock.Anything)
			harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
			harness.btc.AssertNotCalled(t, "GetTransactionBlockInfo", mock.Anything)
			harness.btc.AssertNotCalled(t, "BuildMerkleBranch", mock.Anything)
		})
		t.Run(tc.name+" clears existing reserve", func(t *testing.T) {
			created := submittingClaim()
			created.State = rootstock.PegInClaimCandidate
			created.TxHash = ""
			claims := mocks.NewPegInClaimRepositoryMock(t)
			harness := newClaimHarness(t, claims)
			expectExistingClaim(claims, created)
			claims.On("Update", mock.Anything, mock.MatchedBy(func(claim rootstock.PegInClaim) bool {
				return claim.ReservedWei != nil && claim.ReservedWei.Cmp(entities.NewWei(0)) == 0
			})).Return(nil).Once()
			harness.expectRefetch(10)
			harness.expectPauseNone()
			harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
			harness.btc.On("GetRawTransaction", claimDepositTxID).Return(tc.rawTx, nil).Once()

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			require.ErrorIs(t, err, blockchain.ErrWitnessSerializedTxNotAccepted)
			require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
			claims.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
			harness.pegin.AssertNotCalled(t, "EstimateRequestPegInGas", mock.Anything)
			harness.pegin.AssertNotCalled(t, "SimulateRequestPegIn", mock.Anything)
			harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
			harness.btc.AssertNotCalled(t, "GetTransactionBlockInfo", mock.Anything)
			harness.btc.AssertNotCalled(t, "BuildMerkleBranch", mock.Anything)
		})
	}
}

func TestClaimPegInUseCase_SuccessfulSendWithoutEventStillSettlesByHeight(t *testing.T) {
	repo := newMemoryClaimRepo()
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{TransactionHash: claimRskTxHash, BlockNumber: 100},
	}, nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(102), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	var emptyPegInID [32]byte
	assert.Equal(t, hex.EncodeToString(emptyPegInID[:]), stored.PegInID)
	assert.Equal(t, "0", stored.ReservedWei.String())
}

func TestClaimPegInUseCase_WithinReorgWindowStaysSubmitting(t *testing.T) {
	repo := newMemoryClaimRepo()
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{
			TransactionHash: claimRskTxHash,
			BlockNumber:     100,
			Status:          blockchain.SuccessfulTxStatus,
		},
		Event: blockchain.PegInRequestedEvent{PegInId: [32]byte{3}, RskAddress: test.AnyRskAddress},
	}, nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(101), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimSubmitting, stored.State)
	assert.Equal(t, claimRskTxHash, stored.TxHash)
	assert.NotEqual(t, "0", stored.ReservedWei.String())
}

func TestClaimPegInUseCase_BuildParamsErrorIsRetryable(t *testing.T) {
	cases := []struct {
		name string
		stub func(*claimHarness)
	}{
		{
			name: "raw tx",
			stub: func(h *claimHarness) {
				h.btc.On("GetRawTransaction", claimDepositTxID).Return([]byte(nil), assert.AnError).Once()
			},
		},
		{
			name: "block info",
			stub: func(h *claimHarness) {
				h.btc.On("GetRawTransaction", claimDepositTxID).Return(h.rawTx, nil).Once()
				h.btc.On("GetTransactionBlockInfo", claimDepositTxID).
					Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()
			},
		},
		{
			name: "merkle branch",
			stub: func(h *claimHarness) {
				h.btc.On("GetRawTransaction", claimDepositTxID).Return(h.rawTx, nil).Once()
				h.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(h.block, nil).Once()
				h.btc.On("BuildMerkleBranch", claimDepositTxID).Return(blockchain.MerkleBranch{}, assert.AnError).Once()
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			claims := mocks.NewPegInClaimRepositoryMock(t)
			harness := newClaimHarness(t, claims)
			expectNoExistingClaim(claims)
			harness.expectRefetch(10)
			harness.expectPauseNone()
			harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
			tc.stub(harness)

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
			claims.AssertNotCalled(t, "Insert", mock.Anything, mock.Anything)
			harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		})
	}
}

func TestClaimPegInUseCase_GetAndBtcLookupErrors(t *testing.T) {
	t.Run("claim get", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newClaimHarness(t, claims)
		claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
			Return((*rootstock.PegInClaim)(nil), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
		harness.btc.AssertNotCalled(t, "GetTransactionInfo", mock.Anything)
	})
	t.Run("btc tx info", func(t *testing.T) {
		claims := mocks.NewPegInClaimRepositoryMock(t)
		harness := newClaimHarness(t, claims)
		expectNoExistingClaim(claims)
		harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
}

func TestClaimPegInUseCase_GetHeightErrorDoesNotClaim(t *testing.T) {
	repo := newMemoryClaimRepo()
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{
			TransactionHash: claimRskTxHash,
			BlockNumber:     100,
			Status:          blockchain.SuccessfulTxStatus,
		},
		Event: blockchain.PegInRequestedEvent{PegInId: [32]byte{5}, RskAddress: test.AnyRskAddress},
	}, nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(0), assert.AnError).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	require.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, rootstock.PegInClaimSubmitting, repo.stored().State)
	assert.NotEqual(t, rootstock.PegInClaimClaimed, repo.stored().State)
}

func TestClaimPegInUseCase_HashedSubmitErrorStaysSubmitting(t *testing.T) {
	repo := newMemoryClaimRepo()
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(
		blockchain.RequestPegInResult{
			Receipt: blockchain.TransactionReceipt{TransactionHash: claimRskTxHash, BlockNumber: 100},
		},
		blockchain.TxFailedError,
	).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimSubmitting, stored.State)
	assert.Equal(t, claimRskTxHash, stored.TxHash)
	assert.Equal(t, harness.payable(harness.fee).String(), stored.ReservedWei.String())
	harness.pegin.AssertNumberOfCalls(t, "RequestPegIn", 1)
}

func TestClaimPegInUseCase_UpdateAfterHashDoesNotClaim(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	expectNoExistingClaim(claims)
	expectEmptyInFlight(claims)
	expectInsertOk(claims)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{
			TransactionHash: claimRskTxHash,
			BlockNumber:     100,
			Status:          blockchain.SuccessfulTxStatus,
		},
		Event: blockchain.PegInRequestedEvent{PegInId: [32]byte{7}, RskAddress: test.AnyRskAddress},
	}, nil).Once()
	expectMarkerUpdate(claims)
	claims.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Times(3)

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	harness.rsk.AssertNotCalled(t, "GetHeight", mock.Anything)
	claims.AssertNumberOfCalls(t, "Update", 4)
}

func TestClaimPegInUseCase_RaceLostUpdateErrorDoesNotClaim(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	expectNoExistingClaim(claims)
	expectEmptyInFlight(claims)
	expectInsertOk(claims)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(
		blockchain.RequestPegInResult{},
		blockchain.ErrPegInAlreadyProcessed,
	).Once()
	expectMarkerUpdate(claims)
	claims.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
}

func TestClaimPegInUseCase_InsertErrorDoesNotSubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	expectNoExistingClaim(claims)
	expectEmptyInFlight(claims)
	claims.On("Insert", mock.Anything, mock.Anything).Return(assert.AnError).Once()
	harness.expectPassingGates(entities.NewWei(0))

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_FeeEqualToAmountIsAccepted(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	expectNoExistingClaim(claims)
	expectEmptyInFlight(claims)
	expectInsertOk(claims)
	expectUpdatesOk(claims)
	harness.fee = harness.amount.Copy()
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(winningRequestResult([32]byte{9}), nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(102), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	harness.pegin.AssertNumberOfCalls(t, "RequestPegIn", 1)
}

func TestClaimPegInUseCase_InsertConflictRereadsSubmittedRow(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	originalCreated := time.Date(2024, 1, 2, 3, 4, 5, 0, time.UTC)
	stored := submittingClaim()
	stored.CreatedAt = originalCreated

	harness.expectPassingGates(entities.NewWei(0))
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return((*rootstock.PegInClaim)(nil), nil).Twice()
	claims.On("ListByStates", mock.Anything, rootstock.PegInClaimCandidate, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{}, nil).Once()
	claims.On("Insert", mock.Anything, mock.Anything).Return(rootstock.ErrPegInClaimAlreadyExists).Once()
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return(&stored, nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

type insertRaceRepo struct {
	inner             *memoryClaimRepo
	mu                sync.Mutex
	gets              int
	bothInitialGets   chan struct{}
	successfulInserts int
	conflictInserts   int
	allowSecondInsert chan struct{}
}

func newInsertRaceRepo() *insertRaceRepo {
	return &insertRaceRepo{
		inner:             newMemoryClaimRepo(),
		bothInitialGets:   make(chan struct{}),
		allowSecondInsert: make(chan struct{}),
	}
}

func (repo *insertRaceRepo) Get(ctx context.Context, rskAddress, depositTxID string) (*rootstock.PegInClaim, error) {
	repo.mu.Lock()
	repo.gets++
	if repo.gets == 2 {
		close(repo.bothInitialGets)
	}
	hide := repo.conflictInserts == 0
	repo.mu.Unlock()
	if hide {
		return nil, nil
	}
	return repo.inner.Get(ctx, rskAddress, depositTxID)
}

func (repo *insertRaceRepo) Insert(ctx context.Context, claim rootstock.PegInClaim) error {
	<-repo.bothInitialGets
	repo.mu.Lock()
	if repo.successfulInserts == 0 {
		repo.successfulInserts++
		repo.mu.Unlock()
		return repo.inner.Insert(ctx, claim)
	}
	repo.mu.Unlock()
	<-repo.allowSecondInsert
	repo.mu.Lock()
	repo.conflictInserts++
	repo.mu.Unlock()
	return rootstock.ErrPegInClaimAlreadyExists
}

func (repo *insertRaceRepo) Update(ctx context.Context, claim rootstock.PegInClaim) error {
	err := repo.inner.Update(ctx, claim)
	if err == nil && claim.TxHash != "" {
		select {
		case <-repo.allowSecondInsert:
		default:
			close(repo.allowSecondInsert)
		}
	}
	return err
}

func (repo *insertRaceRepo) ListByStates(ctx context.Context, states ...rootstock.PegInClaimState) ([]rootstock.PegInClaim, error) {
	return repo.inner.ListByStates(ctx, states...)
}

func TestClaimPegInUseCase_InsertConflictPreservesCreatedAtAndTxHash(t *testing.T) {
	repo := newInsertRaceRepo()
	harness := newClaimHarness(t, repo)
	expectSuccessfulGates(harness, 2)

	requestStarted := make(chan struct{})
	var firstCreated time.Time
	harness.pegin.On("RequestPegIn", mock.Anything).
		Run(func(mock.Arguments) {
			firstCreated = repo.inner.stored().CreatedAt
			close(requestStarted)
		}).
		Return(winningRequestResult([32]byte{4}), nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(101), nil).Once()

	var wg sync.WaitGroup
	errs := make(chan error, 2)
	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	}()
	<-requestStarted
	wg.Add(1)
	go func() {
		defer wg.Done()
		errs <- harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	}()
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	stored := repo.inner.stored()
	assert.Equal(t, firstCreated, stored.CreatedAt)
	assert.Equal(t, claimRskTxHash, stored.TxHash)
	assert.Equal(t, rootstock.PegInClaimSubmitting, stored.State)
	harness.pegin.AssertNumberOfCalls(t, "RequestPegIn", 1)
}

func expectSerializedWalletMocks(
	harness *claimHarness,
	depositB string,
	required *entities.Wei,
	requestStarted, holdRequest, requestReturned chan struct{},
	balanceCalls *atomic.Int32,
) {
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.btc.On("GetTransactionInfo", depositB).Return(harness.payingTxFor(depositB, 10), nil).Once()
	harness.configs.On("GetRequiredPegInBtcConfirmations", matchWei(harness.amount)).Return(uint64(6), nil).Twice()
	harness.pause.On("PauseLevel").Return(blockchain.PauseLevelNone, nil).Twice()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Twice()
	harness.btc.On("GetRawTransaction", claimDepositTxID).Return(harness.rawTx, nil).Once()
	harness.btc.On("GetRawTransaction", depositB).Return(harness.rawTx, nil).Once()
	harness.btc.On("GetTransactionBlockInfo", claimDepositTxID).Return(harness.block, nil).Once()
	harness.btc.On("GetTransactionBlockInfo", depositB).Return(harness.block, nil).Once()
	harness.btc.On("BuildMerkleBranch", claimDepositTxID).Return(harness.merkle, nil).Once()
	harness.btc.On("BuildMerkleBranch", depositB).Return(harness.merkle, nil).Once()
	harness.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(harness.estimatedGas, nil).Twice()
	harness.rsk.On("GasPrice", mock.Anything).Return(harness.gasPrice.Copy(), nil).Twice()
	harness.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).
		Run(func(mock.Arguments) {
			if balanceCalls.Add(1) == 2 {
				<-requestReturned
			}
		}).
		Return(required, nil).Twice()
	harness.pegin.On("RequestPegIn", mock.Anything).
		Run(func(mock.Arguments) {
			close(requestStarted)
			<-holdRequest
			close(requestReturned)
		}).
		Return(winningRequestResult([32]byte{6}), nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(101), nil).Once()
}

func TestClaimPegInUseCase_SpendableDecisionSerializedUnderWalletMutex(t *testing.T) {
	depositB := "bb" + claimDepositTxID[2:]
	repo := newMemoryClaimRepo()
	harness := newClaimHarness(t, repo)
	required := harness.spendableRequired(entities.NewWei(0), harness.fee)
	requestStarted := make(chan struct{})
	holdRequest := make(chan struct{})
	requestReturned := make(chan struct{})
	var balanceCalls atomic.Int32
	expectSerializedWalletMocks(harness, depositB, required, requestStarted, holdRequest, requestReturned, &balanceCalls)

	errA := make(chan error, 1)
	errB := make(chan error, 1)
	go func() {
		errA <- harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	}()
	<-requestStarted
	go func() {
		errB <- harness.useCase.Run(context.Background(), harness.entry, depositB)
	}()
	assert.Never(t, func() bool {
		return balanceCalls.Load() > 1
	}, 80*time.Millisecond, 5*time.Millisecond)
	close(holdRequest)
	require.NoError(t, <-errA)
	require.NoError(t, <-errB)

	harness.pegin.AssertNumberOfCalls(t, "RequestPegIn", 1)
	assert.Equal(t, int32(2), balanceCalls.Load())
	accepted := entities.NewWei(0)
	repo.mu.Lock()
	for _, claim := range repo.byKey {
		if claim.ReservedWei != nil {
			accepted.Add(accepted, claim.ReservedWei)
		}
	}
	repo.mu.Unlock()
	assert.LessOrEqual(t, accepted.Cmp(required), 0)
}

func TestClaimPegInUseCase_HashPersistUsesDetachedContext(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	harness.expectPassingGates(entities.NewWei(0))
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return((*rootstock.PegInClaim)(nil), nil)
	claims.On("ListByStates", mock.Anything, rootstock.PegInClaimCandidate, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{}, nil).Once()
	claims.On("Insert", mock.Anything, mock.Anything).Return(nil).Once()
	expectMarkerUpdate(claims)

	ctx, cancel := context.WithCancel(context.Background())
	harness.pegin.On("RequestPegIn", mock.Anything).
		Run(func(mock.Arguments) { cancel() }).
		Return(winningRequestResult([32]byte{8}), nil).Once()

	var sawHashPersist bool
	var persistCanceled bool
	claims.On("Update", mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			claim, ok := args.Get(1).(rootstock.PegInClaim)
			if !ok {
				return
			}
			if sawHashPersist || claim.TxHash != claimRskTxHash {
				return
			}
			ctxArg, ok := args.Get(0).(context.Context)
			if !ok {
				return
			}
			sawHashPersist = true
			persistCanceled = ctxArg.Err() != nil
		}).
		Return(nil)
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(102), nil).Once()

	err := harness.useCase.Run(ctx, harness.entry, claimDepositTxID)
	require.NoError(t, err)
	require.True(t, sawHashPersist)
	require.False(t, persistCanceled)
}

func TestClaimPegInUseCase_HashPersistRetriesThenFails(t *testing.T) {
	cases := []struct {
		name          string
		failTimes     int
		wantErr       bool
		wantUpdates   int
		wantGetHeight bool
	}{
		{name: "success on attempt 1", failTimes: 0, wantUpdates: 3, wantGetHeight: true},
		{name: "success on attempt 2", failTimes: 1, wantUpdates: 4, wantGetHeight: true},
		{name: "success on attempt 3", failTimes: 2, wantUpdates: 5, wantGetHeight: true},
		{name: "failure after exactly three attempts", failTimes: 3, wantErr: true, wantUpdates: 4},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			claims := mocks.NewPegInClaimRepositoryMock(t)
			harness := newClaimHarness(t, claims)
			harness.expectPassingGates(entities.NewWei(0))
			claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
				Return((*rootstock.PegInClaim)(nil), nil)
			claims.On("ListByStates", mock.Anything, rootstock.PegInClaimCandidate, rootstock.PegInClaimSubmitting).
				Return([]rootstock.PegInClaim{}, nil).Once()
			claims.On("Insert", mock.Anything, mock.Anything).Return(nil).Once()
			expectMarkerUpdate(claims)
			harness.pegin.On("RequestPegIn", mock.Anything).
				Return(winningRequestResult([32]byte{2}), nil).Once()

			hashPersist := mock.MatchedBy(func(claim rootstock.PegInClaim) bool {
				return claim.TxHash == claimRskTxHash &&
					claim.State == rootstock.PegInClaimSubmitting &&
					claim.PegInID == ""
			})
			for i := 0; i < tc.failTimes; i++ {
				claims.On("Update", mock.Anything, hashPersist).Return(assert.AnError).Once()
			}
			if !tc.wantErr {
				claims.On("Update", mock.Anything, hashPersist).Return(nil).Once()
				claims.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
				harness.rsk.On("GetHeight", mock.Anything).Return(uint64(102), nil).Once()
			}

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			if tc.wantErr {
				require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
				require.ErrorIs(t, err, assert.AnError)
				harness.rsk.AssertNotCalled(t, "GetHeight", mock.Anything)
			} else {
				require.NoError(t, err)
			}
			claims.AssertNumberOfCalls(t, "Update", tc.wantUpdates)
			harness.pegin.AssertNumberOfCalls(t, "RequestPegIn", 1)
		})
	}
}

func TestClaimPegInUseCase_HashPersistJoinsOriginalSubmitCause(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	harness.expectPassingGates(entities.NewWei(0))
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return((*rootstock.PegInClaim)(nil), nil)
	claims.On("ListByStates", mock.Anything, rootstock.PegInClaimCandidate, rootstock.PegInClaimSubmitting).
		Return([]rootstock.PegInClaim{}, nil).Once()
	claims.On("Insert", mock.Anything, mock.Anything).Return(nil).Once()
	expectMarkerUpdate(claims)
	harness.pegin.On("RequestPegIn", mock.Anything).
		Return(
			blockchain.RequestPegInResult{
				Receipt: blockchain.TransactionReceipt{TransactionHash: claimRskTxHash, BlockNumber: 100},
			},
			blockchain.ErrAddressNotRegistered,
		).Once()
	claims.On("Update", mock.Anything, mock.Anything).Return(assert.AnError).Times(3)

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorIs(t, err, blockchain.ErrAddressNotRegistered)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	harness.pegin.AssertNumberOfCalls(t, "RequestPegIn", 1)
	claims.AssertNumberOfCalls(t, "Update", 4)
}

func TestClaimPegInUseCase_SubmittingEmptyHashDoesNotSubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	expectExistingClaim(claims, submittingEmptyHashClaim())

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.btc.AssertNotCalled(t, "GetTransactionInfo", mock.Anything)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestClaimPegInUseCase_CandidateEmptyHashRetriesSubmit(t *testing.T) {
	created := candidateWithReserve()
	repo := newMemoryClaimRepo(created)
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(winningRequestResult([32]byte{0x11}), nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(102), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	assert.Equal(t, claimRskTxHash, stored.TxHash)
	harness.pegin.AssertNumberOfCalls(t, "RequestPegIn", 1)
}

type failHashPersistRepo struct {
	*memoryClaimRepo
}

func (repo *failHashPersistRepo) Update(ctx context.Context, claim rootstock.PegInClaim) error {
	if claim.TxHash != "" {
		return assert.AnError
	}
	return repo.memoryClaimRepo.Update(ctx, claim)
}

func TestClaimPegInUseCase_HashPersistFailThenNextRunDoesNotSubmit(t *testing.T) {
	repo := &failHashPersistRepo{memoryClaimRepo: newMemoryClaimRepo()}
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(winningRequestResult([32]byte{0x12}), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimSubmitting, stored.State)
	assert.Empty(t, stored.TxHash)

	err = harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	harness.pegin.AssertNumberOfCalls(t, "RequestPegIn", 1)
	assert.Equal(t, rootstock.PegInClaimSubmitting, repo.stored().State)
	assert.Empty(t, repo.stored().TxHash)
}

func TestClaimPegInUseCase_MarkerUpdateErrorDoesNotSubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	expectNoExistingClaim(claims)
	expectEmptyInFlight(claims)
	expectInsertOk(claims)
	harness.expectPassingGates(entities.NewWei(0))
	claims.On("Update", mock.Anything, matchSubmittingEmptyHash()).Return(assert.AnError).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

type failMarkerRepo struct {
	*memoryClaimRepo
}

func (repo *failMarkerRepo) Update(ctx context.Context, claim rootstock.PegInClaim) error {
	if claim.State == rootstock.PegInClaimSubmitting && claim.TxHash == "" {
		return assert.AnError
	}
	return repo.memoryClaimRepo.Update(ctx, claim)
}

func TestClaimPegInUseCase_MarkerUpdateErrorLeavesCandidate(t *testing.T) {
	repo := &failMarkerRepo{memoryClaimRepo: newMemoryClaimRepo()}
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	stored := repo.stored()
	assert.Equal(t, rootstock.PegInClaimCandidate, stored.State)
	assert.Empty(t, stored.TxHash)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_SaveConflictRaceLostDoesNotSubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	stored := rootstock.PegInClaim{
		RskAddress:  test.AnyRskAddress,
		DepositTxID: claimDepositTxID,
		BtcAddress:  claimBtcAddress,
		State:       rootstock.PegInClaimRaceLost,
		ReservedWei: entities.NewWei(0),
	}
	harness.expectPassingGates(entities.NewWei(0))
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return((*rootstock.PegInClaim)(nil), nil).Twice()
	expectEmptyInFlight(claims)
	claims.On("Insert", mock.Anything, mock.Anything).Return(rootstock.ErrPegInClaimAlreadyExists).Once()
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return(&stored, nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}

func TestClaimPegInUseCase_SaveConflictSubmittingEmptyDoesNotSubmit(t *testing.T) {
	claims := mocks.NewPegInClaimRepositoryMock(t)
	harness := newClaimHarness(t, claims)
	stored := submittingEmptyHashClaim()
	harness.expectPassingGates(entities.NewWei(0))
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return((*rootstock.PegInClaim)(nil), nil).Twice()
	expectEmptyInFlight(claims)
	claims.On("Insert", mock.Anything, mock.Anything).Return(rootstock.ErrPegInClaimAlreadyExists).Once()
	claims.On("Get", mock.Anything, test.AnyRskAddress, claimDepositTxID).
		Return(&stored, nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
}
