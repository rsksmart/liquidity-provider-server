package pegin_test

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
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
	mu        sync.Mutex
	byKey     map[string]rootstock.PegInClaim
	insert    int
	getErr    error
	listErr   error
	insertErr error
	updateErr error
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
	if repo.insertErr != nil {
		return repo.insertErr
	}
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
	if repo.getErr != nil {
		return nil, repo.getErr
	}
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
	if repo.updateErr != nil {
		return repo.updateErr
	}
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
	if repo.listErr != nil {
		return nil, repo.listErr
	}
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

type rskAccountStub struct{ addr string }

func (s rskAccountStub) RskAddress() string { return s.addr }

func matchWei(want *entities.Wei) interface{} {
	return mock.MatchedBy(func(got *entities.Wei) bool {
		return got != nil && got.Cmp(want) == 0
	})
}

const claimEstimatedGas uint64 = 100000

type claimHarness struct {
	repo         *memoryClaimRepo
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

func newClaimHarness(t *testing.T, repo *memoryClaimRepo) *claimHarness {
	t.Helper()
	amount := entities.NewWei(1_000_000_000_000_000_000)
	fee := entities.NewWei(1_000_000_000_000_000)
	harness := &claimHarness{
		repo:         repo,
		pegin:        new(mocks.PeginContractMock),
		pause:        new(mocks.PauseRegistryContractMock),
		configs:      new(mocks.FlyoverConfigurationsContractMock),
		btc:          new(mocks.BtcRpcMock),
		rsk:          new(mocks.RootstockRpcServerMock),
		amount:       amount,
		fee:          fee,
		gasPrice:     entities.NewWei(1),
		estimatedGas: claimEstimatedGas,
		rawTx:        []byte{1, 2, 3, 4, 5},
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
		rskAccountStub{addr: test.AnyRskAddress},
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
		harness := newClaimHarness(t, newMemoryClaimRepo())
		harness.expectRefetch(2)

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		assert.Equal(t, 0, harness.repo.insert)
		assert.Empty(t, harness.repo.byKey)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		harness.pause.AssertNotCalled(t, "PauseLevel")
		harness.rsk.AssertNotCalled(t, "GetBalance", mock.Anything, mock.Anything)
		harness.btc.AssertNotCalled(t, "GetRawTransaction", mock.Anything)
	})
	t.Run("releases existing reserve", func(t *testing.T) {
		created := submittingClaim()
		created.State = rootstock.PegInClaimCandidate
		created.TxHash = ""
		harness := newClaimHarness(t, newMemoryClaimRepo(created))
		harness.expectRefetch(2)

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		assert.Equal(t, "0", harness.repo.stored().ReservedWei.String())
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
}

func TestClaimPegInUseCase_WinPersistsClaimedAndAdapterArgs(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
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

	stored := harness.repo.stored()
	assert.Equal(t, rootstock.PegInClaimClaimed, stored.State)
	assert.Equal(t, hex.EncodeToString(pegInID[:]), stored.PegInID)
	assert.Equal(t, claimRskTxHash, stored.TxHash)
	assert.Equal(t, "0", stored.ReservedWei.String())
	harness.pegin.AssertExpectations(t)
}

func TestClaimPegInUseCase_PegInAlreadyProcessedIsQuietRaceLost(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	harness.expectPassingGates(entities.NewWei(0))
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
			harness.expectPassingGates(entities.NewWei(0))
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
	harness.expectPauseNone()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.expectBuildParams()
	harness.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).Return(entities.NewWei(1), nil).Once()
	harness.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(harness.estimatedGas, nil).Once()
	harness.rsk.On("GasPrice", mock.Anything).Return(harness.gasPrice.Copy(), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	assert.Equal(t, 0, harness.repo.insert)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_HardPauseDoesNotSubmit(t *testing.T) {
	t.Run("no existing claim", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo())
		harness.expectRefetch(10)
		harness.pause.On("PauseLevel").Return(blockchain.PauseLevelHard, nil).Once()

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		assert.Equal(t, 0, harness.repo.insert)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		harness.rsk.AssertNotCalled(t, "GetBalance", mock.Anything, mock.Anything)
	})
	t.Run("releases existing reserve", func(t *testing.T) {
		created := submittingClaim()
		created.State = rootstock.PegInClaimCandidate
		created.TxHash = ""
		harness := newClaimHarness(t, newMemoryClaimRepo(created))
		harness.expectRefetch(10)
		harness.pause.On("PauseLevel").Return(blockchain.PauseLevelHard, nil).Once()

		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.NoError(t, err)
		assert.Equal(t, "0", harness.repo.stored().ReservedWei.String())
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

func (h *claimHarness) expectCanonicalSuccess(event blockchain.PegInRequestedEvent, height uint64) {
	receipt := successReceipt()
	h.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()
	h.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
	h.pegin.On("UnpackPegInRequested", receipt).Return(event, nil).Once()
	h.rsk.On("GetHeight", mock.Anything).Return(height, nil).Once()
}

func TestClaimPegInUseCase_CrashAfterMinePersistsClaimedWithoutSecondBroadcast(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	pegInID := [32]byte{0xca, 0xfe}
	harness.expectCanonicalSuccess(blockchain.PegInRequestedEvent{
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

func TestClaimPegInUseCase_ZeroFirstOutputDoesNotSubmit(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	tx := blockchain.BitcoinTransactionInformation{
		Hash:          claimDepositTxID,
		Confirmations: 10,
		Outputs:       map[string][]*entities.Wei{"other": {harness.amount.Copy()}},
	}
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(tx, nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	assert.Equal(t, 0, harness.repo.insert)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.pause.AssertNotCalled(t, "PauseLevel")
}

func TestClaimPegInUseCase_RunWithTxHashReconcilesWithoutSecondBroadcast(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	pegInID := [32]byte{0x11}
	harness.expectCanonicalSuccess(blockchain.PegInRequestedEvent{
		PegInId:    pegInID,
		RskAddress: test.AnyRskAddress,
	}, 102)

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	assert.Equal(t, rootstock.PegInClaimClaimed, harness.repo.stored().State)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.btc.AssertNotCalled(t, "GetTransactionInfo", mock.Anything)
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
	harness := newClaimHarness(t, newMemoryClaimRepo(other, nilReserved, self))
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
	assert.Equal(t, self.CreatedAt, harness.repo.stored().CreatedAt)
	harness.pegin.AssertExpectations(t)
}

func TestClaimPegInUseCase_OracleAndLiquidityErrorsDoNotSubmit(t *testing.T) {
	t.Run("confirmations oracle", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo())
		harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
		harness.configs.On("GetRequiredPegInBtcConfirmations", matchWei(harness.amount)).Return(uint64(0), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.Error(t, err)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("hard-pause oracle", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo())
		harness.expectRefetch(10)
		harness.pause.On("PauseLevel").Return(uint8(0), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.Error(t, err)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("fee oracle", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo())
		harness.expectRefetch(10)
		harness.expectPauseNone()
		harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return((*entities.Wei)(nil), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.Error(t, err)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("wallet balance error", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo())
		harness.expectRefetch(10)
		harness.expectPauseNone()
		harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
		harness.expectBuildParams()
		harness.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).Return((*entities.Wei)(nil), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.Error(t, err)
		assert.Equal(t, 0, harness.repo.insert)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("gas estimate error", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo())
		harness.expectRefetch(10)
		harness.expectPauseNone()
		harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
		harness.expectBuildParams()
		harness.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).
			Return(harness.spendableRequired(entities.NewWei(0), harness.fee), nil).Once()
		harness.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(uint64(0), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.Error(t, err)
		assert.Equal(t, 0, harness.repo.insert)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("gas price error", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo())
		harness.expectRefetch(10)
		harness.expectPauseNone()
		harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
		harness.expectBuildParams()
		harness.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).
			Return(harness.spendableRequired(entities.NewWei(0), harness.fee), nil).Once()
		harness.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(harness.estimatedGas, nil).Once()
		harness.rsk.On("GasPrice", mock.Anything).Return((*entities.Wei)(nil), assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.Error(t, err)
		assert.Equal(t, 0, harness.repo.insert)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("in-flight list error", func(t *testing.T) {
		repo := newMemoryClaimRepo()
		repo.listErr = assert.AnError
		harness := newClaimHarness(t, repo)
		harness.expectRefetch(10)
		harness.expectPauseNone()
		harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.Error(t, err)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
}

func TestClaimPegInUseCase_FeeAboveAmountStillEvaluatesZeroPayable(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	harness.expectRefetch(10)
	harness.expectPauseNone()
	highFee := new(entities.Wei).Add(harness.amount.Copy(), entities.NewWei(1))
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(highFee, nil).Once()
	harness.expectBuildParams()
	harness.rsk.On("GetBalance", mock.Anything, test.AnyRskAddress).
		Return(harness.spendableRequired(entities.NewWei(0), highFee), nil).Once()
	harness.pegin.On("EstimateRequestPegInGas", mock.Anything).Return(harness.estimatedGas, nil).Once()
	harness.rsk.On("GasPrice", mock.Anything).Return(harness.gasPrice.Copy(), nil).Once()
	harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{
			TransactionHash: claimRskTxHash,
			BlockNumber:     100,
			Status:          blockchain.SuccessfulTxStatus,
		},
		Event: blockchain.PegInRequestedEvent{PegInId: [32]byte{2}, RskAddress: test.AnyRskAddress},
	}, nil).Once()
	harness.rsk.On("GetHeight", mock.Anything).Return(uint64(102), nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.NoError(t, err)
	harness.pegin.AssertExpectations(t)
}

func TestClaimPegInUseCase_MissingEventIsRetryableNotClaimed(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{TransactionHash: claimRskTxHash, BlockNumber: 100},
	}, nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.Error(t, err)
	stored := harness.repo.stored()
	assert.Equal(t, rootstock.PegInClaimRetryableFailure, stored.State)
	assert.NotEqual(t, rootstock.PegInClaimClaimed, stored.State)
}

func TestClaimPegInUseCase_WithinReorgWindowStaysSubmitting(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
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
	stored := harness.repo.stored()
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
			harness := newClaimHarness(t, newMemoryClaimRepo())
			harness.expectRefetch(10)
			harness.expectPauseNone()
			harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
			tc.stub(harness)

			err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
			require.Error(t, err)
			assert.Equal(t, 0, harness.repo.insert)
			harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
		})
	}
}

func TestClaimPegInUseCase_GetAndBtcLookupErrors(t *testing.T) {
	t.Run("claim get", func(t *testing.T) {
		repo := newMemoryClaimRepo()
		repo.getErr = assert.AnError
		harness := newClaimHarness(t, repo)
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.Error(t, err)
		harness.btc.AssertNotCalled(t, "GetTransactionInfo", mock.Anything)
	})
	t.Run("btc tx info", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo())
		harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
		err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
		require.Error(t, err)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
}

func TestClaimPegInUseCase_ReconcileSubmittingListError(t *testing.T) {
	repo := newMemoryClaimRepo(submittingClaim())
	repo.listErr = assert.AnError
	harness := newClaimHarness(t, repo)
	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.Error(t, err)
	harness.rsk.AssertNotCalled(t, "GetTransactionReceipt", mock.Anything, mock.Anything)
}

func TestClaimPegInUseCase_EmptyTxHashDoesNotResubmit(t *testing.T) {
	claim := submittingClaim()
	claim.TxHash = ""
	harness := newClaimHarness(t, newMemoryClaimRepo(claim))

	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.NoError(t, err)
	assert.Equal(t, rootstock.PegInClaimSubmitting, harness.repo.stored().State)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	harness.rsk.AssertNotCalled(t, "GetTransactionReceipt", mock.Anything, mock.Anything)
}

func TestClaimPegInUseCase_ReceiptLookupErrorDoesNotResubmit(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).
		Return(blockchain.TransactionReceipt{}, assert.AnError).Once()

	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.Error(t, err)
	assert.Equal(t, rootstock.PegInClaimSubmitting, harness.repo.stored().State)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_StatusOneWithoutEventDoesNotResubmit(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	receipt := successReceipt()
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()
	harness.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).
		Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
	harness.pegin.On("UnpackPegInRequested", receipt).Return(blockchain.PegInRequestedEvent{}, errors.New("missing")).Once()

	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.NoError(t, err)
	assert.Equal(t, rootstock.PegInClaimSubmitting, harness.repo.stored().State)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_NonCanonicalReceiptDoesNotUnpack(t *testing.T) {
	t.Run("removed log", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
		receipt := successReceipt()
		receipt.Logs = []blockchain.TransactionLog{{Removed: true}}
		harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()

		err := harness.useCase.ReconcileSubmitting(context.Background())
		require.NoError(t, err)
		assert.Equal(t, rootstock.PegInClaimSubmitting, harness.repo.stored().State)
		harness.pegin.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("empty block hash", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
		receipt := successReceipt()
		receipt.BlockHash = ""
		harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()

		err := harness.useCase.ReconcileSubmitting(context.Background())
		require.NoError(t, err)
		harness.rsk.AssertNotCalled(t, "GetBlockByHash", mock.Anything, mock.Anything)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("unknown status", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
		receipt := successReceipt()
		receipt.Status = 99
		harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(receipt, nil).Once()
		harness.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).
			Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()

		err := harness.useCase.ReconcileSubmitting(context.Background())
		require.NoError(t, err)
		harness.pegin.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
}

func TestClaimPegInUseCase_StatusZeroIdentifyNilDoesNotResubmit(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	reverted := successReceipt()
	reverted.Status = 0
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(reverted, nil).Once()
	harness.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).
		Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
	harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
	harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return(harness.fee.Copy(), nil).Once()
	harness.expectBuildParams()
	harness.pegin.On("IdentifyRequestPegIn", mock.Anything).Return(nil).Once()

	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.Error(t, err)
	assert.Equal(t, rootstock.PegInClaimRetryableFailure, harness.repo.stored().State)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_StatusZeroIdentifyLookupErrors(t *testing.T) {
	t.Run("btc tx info", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
		reverted := successReceipt()
		reverted.Status = 0
		harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(reverted, nil).Once()
		harness.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).
			Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
		harness.btc.On("GetTransactionInfo", claimDepositTxID).
			Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()

		err := harness.useCase.ReconcileSubmitting(context.Background())
		require.Error(t, err)
		harness.pegin.AssertNotCalled(t, "IdentifyRequestPegIn", mock.Anything)
		harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
	})
	t.Run("fee oracle", func(t *testing.T) {
		harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
		reverted := successReceipt()
		reverted.Status = 0
		harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(reverted, nil).Once()
		harness.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).
			Return(blockchain.BlockInfo{Hash: claimRskBlockHash}, nil).Once()
		harness.btc.On("GetTransactionInfo", claimDepositTxID).Return(harness.payingTx(10), nil).Once()
		harness.configs.On("CalculatePegInFee", matchWei(harness.amount)).Return((*entities.Wei)(nil), assert.AnError).Once()

		err := harness.useCase.ReconcileSubmitting(context.Background())
		require.Error(t, err)
		harness.pegin.AssertNotCalled(t, "IdentifyRequestPegIn", mock.Anything)
	})
}

func TestClaimPegInUseCase_GetHeightErrorDoesNotClaim(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
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
	require.Error(t, err)
	assert.Equal(t, rootstock.PegInClaimSubmitting, harness.repo.stored().State)
	assert.NotEqual(t, rootstock.PegInClaimClaimed, harness.repo.stored().State)
}

func TestClaimPegInUseCase_BlockHashMismatchStaysSubmitting(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo(submittingClaim()))
	harness.rsk.On("GetTransactionReceipt", mock.Anything, claimRskTxHash).Return(successReceipt(), nil).Once()
	harness.rsk.On("GetBlockByHash", mock.Anything, claimRskBlockHash).
		Return(blockchain.BlockInfo{Hash: "0xotherblock"}, nil).Once()

	err := harness.useCase.ReconcileSubmitting(context.Background())
	require.NoError(t, err)
	assert.Equal(t, rootstock.PegInClaimSubmitting, harness.repo.stored().State)
	harness.pegin.AssertNotCalled(t, "UnpackPegInRequested", mock.Anything)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}

func TestClaimPegInUseCase_SubmitErrorAfterHashPersistsThenClassifies(t *testing.T) {
	harness := newClaimHarness(t, newMemoryClaimRepo())
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(
		blockchain.RequestPegInResult{
			Receipt: blockchain.TransactionReceipt{TransactionHash: claimRskTxHash, BlockNumber: 100},
		},
		blockchain.ErrAddressNotRegistered,
	).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.Error(t, err)
	stored := harness.repo.stored()
	assert.Equal(t, rootstock.PegInClaimRetryableFailure, stored.State)
	assert.Equal(t, claimRskTxHash, stored.TxHash)
	assert.NotEqual(t, rootstock.PegInClaimClaimed, stored.State)
}

func TestClaimPegInUseCase_UpdateAfterHashDoesNotClaim(t *testing.T) {
	repo := newMemoryClaimRepo()
	repo.updateErr = assert.AnError
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(blockchain.RequestPegInResult{
		Receipt: blockchain.TransactionReceipt{
			TransactionHash: claimRskTxHash,
			BlockNumber:     100,
			Status:          blockchain.SuccessfulTxStatus,
		},
		Event: blockchain.PegInRequestedEvent{PegInId: [32]byte{7}, RskAddress: test.AnyRskAddress},
	}, nil).Once()

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.Error(t, err)
	assert.Equal(t, rootstock.PegInClaimCandidate, harness.repo.stored().State)
	assert.NotEqual(t, rootstock.PegInClaimClaimed, harness.repo.stored().State)
}

func TestClaimPegInUseCase_RaceLostUpdateErrorDoesNotClaim(t *testing.T) {
	repo := newMemoryClaimRepo()
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))
	harness.pegin.On("RequestPegIn", mock.Anything).Return(
		blockchain.RequestPegInResult{},
		blockchain.ErrPegInAlreadyProcessed,
	).Once()
	repo.updateErr = assert.AnError

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.Error(t, err)
	assert.Equal(t, rootstock.PegInClaimCandidate, harness.repo.stored().State)
	assert.NotEqual(t, rootstock.PegInClaimClaimed, harness.repo.stored().State)
}

func TestClaimPegInUseCase_InsertErrorDoesNotSubmit(t *testing.T) {
	repo := newMemoryClaimRepo()
	repo.insertErr = assert.AnError
	harness := newClaimHarness(t, repo)
	harness.expectPassingGates(entities.NewWei(0))

	err := harness.useCase.Run(context.Background(), harness.entry, claimDepositTxID)
	require.Error(t, err)
	assert.Equal(t, 0, harness.repo.insert)
	harness.pegin.AssertNotCalled(t, "RequestPegIn", mock.Anything)
}
