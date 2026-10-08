package pegin_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
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
	resolveDepositTxID     = "resolve deposit tx id"
	resolveTxHash          = "resolve tx hash"
	resolveRequiredConfirm = uint64(100)
)

var (
	resolveRawTx   = []byte{1, 2, 3, 4, 5}
	resolvePmt     = []byte{6, 7, 8, 9}
	resolveBlock   = blockchain.BitcoinBlockInformation{Hash: [32]byte{1}, Height: big.NewInt(900)}
	resolveReceipt = blockchain.TransactionReceipt{
		TransactionHash: resolveTxHash,
		BlockHash:       "resolve block hash",
		BlockNumber:     500,
		GasUsed:         big.NewInt(150000),
		GasPrice:        entities.NewWei(60000000),
	}
)

type resolveHarness struct {
	claims  *mocks.PegInClaimRepositoryMock
	pegin   *mocks.PeginContractMock
	pause   *mocks.PauseRegistryContractMock
	bridge  *mocks.BridgeMock
	btc     *mocks.BitcoinNetworkMock
	mutex   *mocks.LockerMock
	useCase *pegin.ResolvePegInUseCase
}

func newResolveHarness(t *testing.T) *resolveHarness {
	harness := &resolveHarness{
		claims: mocks.NewPegInClaimRepositoryMock(t),
		pegin:  mocks.NewPeginContractMock(t),
		pause:  mocks.NewPauseRegistryContractMock(t),
		bridge: mocks.NewBridgeMock(t),
		btc:    mocks.NewBitcoinNetworkMock(t),
		mutex:  new(mocks.LockerMock),
	}
	harness.useCase = pegin.NewResolvePegInUseCase(
		harness.claims,
		blockchain.RskContracts{PegIn: harness.pegin, PauseRegistry: harness.pause, Bridge: harness.bridge},
		harness.btc,
		harness.mutex,
	)
	return harness
}

func resolvableClaim() rootstock.PegInClaim {
	return rootstock.PegInClaim{
		RskAddress:    test.AnyRskAddress,
		DepositTxID:   resolveDepositTxID,
		BtcAddress:    test.AnyBtcAddress,
		State:         rootstock.PegInClaimClaimed,
		RequestTxHash: "request tx hash",
		PegInID:       "0a1b",
	}
}

func resolveParams() blockchain.ResolvePegInParams {
	return blockchain.ResolvePegInParams{
		RskAddress:            test.AnyRskAddress,
		BitcoinRawTransaction: resolveRawTx,
		PartialMerkleTree:     resolvePmt,
		BlockHeight:           resolveBlock.Height,
	}
}

func (h *resolveHarness) expectReady(pauseLevel uint8) {
	h.pause.EXPECT().PauseLevel().Return(pauseLevel, nil).Once()
	h.btc.On("GetTransactionInfo", resolveDepositTxID).
		Return(blockchain.BitcoinTransactionInformation{Hash: resolveDepositTxID, Confirmations: resolveRequiredConfirm}, nil).Once()
	h.bridge.EXPECT().GetRequiredTxConfirmations().Return(resolveRequiredConfirm).Once()
	h.btc.On("GetRawTransaction", resolveDepositTxID).Return(resolveRawTx, nil).Once()
	h.btc.On("GetPartialMerkleTree", resolveDepositTxID).Return(resolvePmt, nil).Once()
	h.btc.On("GetTransactionBlockInfo", resolveDepositTxID).Return(resolveBlock, nil).Once()
	h.mutex.On("Lock").Return().Once()
	h.mutex.On("Unlock").Return().Once()
}

func (h *resolveHarness) expectResolve(result blockchain.ResolvePegInResult, err error) {
	h.pegin.EXPECT().ResolvePegIn(resolveParams()).Return(result, err).Once()
}

func (h *resolveHarness) expectUpdate(t *testing.T, expected rootstock.PegInClaim) {
	h.claims.EXPECT().Update(test.AnyCtx, mock.MatchedBy(func(claim rootstock.PegInClaim) bool {
		claim.UpdatedAt = expected.UpdatedAt
		return assert.Equal(t, expected, claim)
	})).Return(nil).Once()
}

func (h *resolveHarness) assertLocked(t *testing.T) {
	h.mutex.AssertExpectations(t)
}

func TestResolvePegInUseCase_Run(t *testing.T) {
	harness := newResolveHarness(t)
	harness.expectReady(blockchain.PauseLevelNone)
	harness.expectResolve(blockchain.ResolvePegInResult{
		Receipt:       resolveReceipt,
		ClaimerPayout: entities.NewWei(1000),
		RegistrantFee: entities.NewWei(10),
	}, nil)
	expected := resolvableClaim()
	expected.State = rootstock.PegInClaimResolved
	expected.ResolveTxHash = resolveTxHash
	expected.ResolveGasUsed = 150000
	expected.ResolveGasPrice = entities.NewWei(60000000)
	expected.ClaimerPayout = entities.NewWei(1000)
	expected.RegistrantFee = entities.NewWei(10)
	harness.expectUpdate(t, expected)

	err := harness.useCase.Run(context.Background(), resolvableClaim())
	require.NoError(t, err)
	harness.assertLocked(t)
}

func TestResolvePegInUseCase_Run_SoftPauseDoesNotBlock(t *testing.T) {
	harness := newResolveHarness(t)
	harness.expectReady(blockchain.PauseLevelSoft)
	harness.expectResolve(blockchain.ResolvePegInResult{Receipt: resolveReceipt, ClaimerPayout: entities.NewWei(1), RegistrantFee: entities.NewWei(0)}, nil)
	harness.claims.EXPECT().Update(test.AnyCtx, mock.MatchedBy(func(claim rootstock.PegInClaim) bool {
		return claim.State == rootstock.PegInClaimResolved
	})).Return(nil).Once()

	err := harness.useCase.Run(context.Background(), resolvableClaim())
	require.NoError(t, err)
	harness.assertLocked(t)
}

func TestResolvePegInUseCase_Run_AlreadyProcessed(t *testing.T) {
	harness := newResolveHarness(t)
	harness.expectReady(blockchain.PauseLevelNone)
	harness.expectResolve(blockchain.ResolvePegInResult{}, blockchain.ErrPegInAlreadyProcessed)
	expected := resolvableClaim()
	expected.State = rootstock.PegInClaimResolved
	harness.expectUpdate(t, expected)

	err := harness.useCase.Run(context.Background(), resolvableClaim())
	require.NoError(t, err)
	harness.assertLocked(t)
}

func TestResolvePegInUseCase_Run_WaitingForBridge(t *testing.T) {
	harness := newResolveHarness(t)
	harness.expectReady(blockchain.PauseLevelNone)
	harness.expectResolve(blockchain.ResolvePegInResult{}, blockchain.WaitingForBridgeError)

	err := harness.useCase.Run(context.Background(), resolvableClaim())
	require.ErrorIs(t, err, blockchain.WaitingForBridgeError)
	require.NotErrorIs(t, err, usecases.NonRecoverableError)
	require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
	harness.claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	harness.assertLocked(t)
}

func TestResolvePegInUseCase_Run_ResolveFailed(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"bridge rejected", fmt.Errorf("%w: code %d", blockchain.ErrBridgeRejectedPegIn, -200)},
		{"not claimed", blockchain.ErrPegInNotClaimed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			harness := newResolveHarness(t)
			harness.expectReady(blockchain.PauseLevelNone)
			harness.expectResolve(blockchain.ResolvePegInResult{}, tc.err)
			expected := resolvableClaim()
			expected.State = rootstock.PegInClaimResolveFailed
			harness.expectUpdate(t, expected)

			err := harness.useCase.Run(context.Background(), resolvableClaim())
			require.ErrorIs(t, err, tc.err)
			require.ErrorIs(t, err, usecases.NonRecoverableError)
			require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
			harness.assertLocked(t)
		})
	}
}

func TestResolvePegInUseCase_Run_SentTxDidNotResolve(t *testing.T) {
	cases := []struct {
		name string
		err  error
	}{
		{"reverted", errors.New("resolve pegin error: transaction reverted (resolve tx hash)")},
		{"no resolved event", fmt.Errorf("resolve pegin error: %w", blockchain.ErrPegInNotResolved)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			harness := newResolveHarness(t)
			harness.expectReady(blockchain.PauseLevelNone)
			harness.expectResolve(blockchain.ResolvePegInResult{Receipt: resolveReceipt}, tc.err)
			expected := resolvableClaim()
			expected.ResolveTxHash = resolveTxHash
			expected.ResolveGasUsed = 150000
			expected.ResolveGasPrice = entities.NewWei(60000000)
			harness.expectUpdate(t, expected)

			err := harness.useCase.Run(context.Background(), resolvableClaim())
			require.ErrorIs(t, err, tc.err)
			require.NotErrorIs(t, err, usecases.NonRecoverableError)
			require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
			harness.assertLocked(t)
		})
	}
}

func TestResolvePegInUseCase_Run_AdapterErrorWithoutTx(t *testing.T) {
	harness := newResolveHarness(t)
	harness.expectReady(blockchain.PauseLevelNone)
	harness.expectResolve(blockchain.ResolvePegInResult{}, errors.New("resolvePegIn reverted with: EnforcedPause"))

	err := harness.useCase.Run(context.Background(), resolvableClaim())
	require.Error(t, err)
	require.NotErrorIs(t, err, usecases.NonRecoverableError)
	require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
	harness.claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	harness.assertLocked(t)
}

func TestResolvePegInUseCase_Run_NotResolvable(t *testing.T) {
	emptyPegInID := resolvableClaim()
	emptyPegInID.PegInID = ""
	cases := []struct {
		name  string
		claim rootstock.PegInClaim
	}{
		{"empty peg-in id", emptyPegInID},
	}
	for _, state := range []rootstock.PegInClaimState{
		rootstock.PegInClaimCandidate,
		rootstock.PegInClaimSubmitting,
		rootstock.PegInClaimRaceLost,
		rootstock.PegInClaimRetryableFailure,
		rootstock.PegInClaimResolved,
		rootstock.PegInClaimResolveFailed,
	} {
		claim := resolvableClaim()
		claim.State = state
		cases = append(cases, struct {
			name  string
			claim rootstock.PegInClaim
		}{string(state), claim})
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			harness := newResolveHarness(t)

			err := harness.useCase.Run(context.Background(), tc.claim)
			require.ErrorIs(t, err, usecases.WrongStateError)
			require.ErrorIs(t, err, usecases.NonRecoverableError)
			harness.pegin.AssertNotCalled(t, "ResolvePegIn", mock.Anything)
			harness.claims.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
			harness.btc.AssertNotCalled(t, "GetTransactionInfo", mock.Anything)
		})
	}
}

func TestResolvePegInUseCase_Run_HardPause(t *testing.T) {
	harness := newResolveHarness(t)
	harness.pause.EXPECT().PauseLevel().Return(blockchain.PauseLevelHard, nil).Once()
	harness.pause.EXPECT().GetAddress().Return(test.AnyRskAddress).Once()

	err := harness.useCase.Run(context.Background(), resolvableClaim())
	require.ErrorIs(t, err, blockchain.ContractPausedError)
	require.NotErrorIs(t, err, usecases.NonRecoverableError)
	require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
	harness.pegin.AssertNotCalled(t, "ResolvePegIn", mock.Anything)
	harness.btc.AssertNotCalled(t, "GetTransactionInfo", mock.Anything)
}

func TestResolvePegInUseCase_Run_PauseLevelError(t *testing.T) {
	harness := newResolveHarness(t)
	harness.pause.EXPECT().PauseLevel().Return(uint8(0), assert.AnError).Once()

	err := harness.useCase.Run(context.Background(), resolvableClaim())
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
	require.NotErrorIs(t, err, usecases.NonRecoverableError)
	harness.pegin.AssertNotCalled(t, "ResolvePegIn", mock.Anything)
}

func TestResolvePegInUseCase_Run_NotEnoughConfirmations(t *testing.T) {
	harness := newResolveHarness(t)
	harness.pause.EXPECT().PauseLevel().Return(blockchain.PauseLevelNone, nil).Once()
	harness.btc.On("GetTransactionInfo", resolveDepositTxID).
		Return(blockchain.BitcoinTransactionInformation{Hash: resolveDepositTxID, Confirmations: resolveRequiredConfirm - 1}, nil).Once()
	harness.bridge.EXPECT().GetRequiredTxConfirmations().Return(resolveRequiredConfirm).Once()

	err := harness.useCase.Run(context.Background(), resolvableClaim())
	require.ErrorIs(t, err, usecases.NoEnoughConfirmationsError)
	require.NotErrorIs(t, err, usecases.NonRecoverableError)
	require.NotErrorIs(t, err, usecases.InfrastructureUnavailableError)
	harness.btc.AssertNotCalled(t, "GetRawTransaction", mock.Anything)
	harness.pegin.AssertNotCalled(t, "ResolvePegIn", mock.Anything)
}

func TestResolvePegInUseCase_Run_BtcRpcErrors(t *testing.T) {
	cases := []struct {
		name  string
		setup func(h *resolveHarness)
	}{
		{"transaction info", func(h *resolveHarness) {
			h.btc.On("GetTransactionInfo", resolveDepositTxID).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
		}},
		{"raw transaction", func(h *resolveHarness) {
			h.expectConfirmed()
			h.btc.On("GetRawTransaction", resolveDepositTxID).Return([]byte(nil), assert.AnError).Once()
		}},
		{"partial merkle tree", func(h *resolveHarness) {
			h.expectConfirmed()
			h.btc.On("GetRawTransaction", resolveDepositTxID).Return(resolveRawTx, nil).Once()
			h.btc.On("GetPartialMerkleTree", resolveDepositTxID).Return([]byte(nil), assert.AnError).Once()
		}},
		{"block info", func(h *resolveHarness) {
			h.expectConfirmed()
			h.btc.On("GetRawTransaction", resolveDepositTxID).Return(resolveRawTx, nil).Once()
			h.btc.On("GetPartialMerkleTree", resolveDepositTxID).Return(resolvePmt, nil).Once()
			h.btc.On("GetTransactionBlockInfo", resolveDepositTxID).Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			harness := newResolveHarness(t)
			harness.pause.EXPECT().PauseLevel().Return(blockchain.PauseLevelNone, nil).Once()
			tc.setup(harness)

			err := harness.useCase.Run(context.Background(), resolvableClaim())
			require.ErrorIs(t, err, assert.AnError)
			require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
			require.NotErrorIs(t, err, usecases.NonRecoverableError)
			harness.pegin.AssertNotCalled(t, "ResolvePegIn", mock.Anything)
			harness.mutex.AssertNotCalled(t, "Lock")
		})
	}
}

func (h *resolveHarness) expectConfirmed() {
	h.btc.On("GetTransactionInfo", resolveDepositTxID).
		Return(blockchain.BitcoinTransactionInformation{Hash: resolveDepositTxID, Confirmations: resolveRequiredConfirm}, nil).Once()
	h.bridge.EXPECT().GetRequiredTxConfirmations().Return(resolveRequiredConfirm).Once()
}

func TestResolvePegInUseCase_Run_UpdateError(t *testing.T) {
	cases := []struct {
		name      string
		resultErr error
		state     rootstock.PegInClaimState
	}{
		{"after resolve", nil, rootstock.PegInClaimResolved},
		{"after resolve failed", blockchain.ErrPegInNotClaimed, rootstock.PegInClaimResolveFailed},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			harness := newResolveHarness(t)
			harness.expectReady(blockchain.PauseLevelNone)
			harness.expectResolve(blockchain.ResolvePegInResult{Receipt: resolveReceipt, ClaimerPayout: entities.NewWei(1), RegistrantFee: entities.NewWei(0)}, tc.resultErr)
			harness.claims.EXPECT().Update(test.AnyCtx, mock.MatchedBy(func(claim rootstock.PegInClaim) bool {
				return claim.State == tc.state
			})).Return(assert.AnError).Once()

			err := harness.useCase.Run(context.Background(), resolvableClaim())
			require.ErrorIs(t, err, assert.AnError)
			require.ErrorIs(t, err, usecases.InfrastructureUnavailableError)
			require.NotErrorIs(t, err, usecases.NonRecoverableError)
			harness.assertLocked(t)
		})
	}
}

func TestResolvePegInUseCase_Run_UpdateErrorHealsOnNextRun(t *testing.T) {
	harness := newResolveHarness(t)
	harness.expectReady(blockchain.PauseLevelNone)
	harness.expectResolve(blockchain.ResolvePegInResult{Receipt: resolveReceipt, ClaimerPayout: entities.NewWei(1), RegistrantFee: entities.NewWei(0)}, nil)
	harness.claims.EXPECT().Update(test.AnyCtx, mock.Anything).Return(assert.AnError).Once()

	err := harness.useCase.Run(context.Background(), resolvableClaim())
	require.Error(t, err)
	require.NotErrorIs(t, err, usecases.NonRecoverableError)

	harness.expectReady(blockchain.PauseLevelNone)
	harness.expectResolve(blockchain.ResolvePegInResult{}, blockchain.ErrPegInAlreadyProcessed)
	harness.claims.EXPECT().Update(test.AnyCtx, mock.MatchedBy(func(claim rootstock.PegInClaim) bool {
		return claim.State == rootstock.PegInClaimResolved && claim.ResolveTxHash == ""
	})).Return(nil).Once()

	err = harness.useCase.Run(context.Background(), resolvableClaim())
	require.NoError(t, err)
	harness.pegin.AssertNumberOfCalls(t, "ResolvePegIn", 2)
	harness.assertLocked(t)
}
