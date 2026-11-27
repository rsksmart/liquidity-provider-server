package pegin_test

import (
	"context"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	resolveTxHash                = "0x4b2d6f8a0c1e3a5b7d9f1b3d5f7a9c1e3b5d7f9a1c3e5b7d9f1a3c5e7b9d1f3a"
	resolveRegistrant            = "0x2A6F4b9eD52cb0D9C2eA1Fb3E6d8B61B43C0dB2e"
	resolveRequiredConfirmations = uint64(100)
	resolveBlockHeight           = int64(812345)
	previousResolveTxHash        = "0x9e7c5a3f1d0b8e6c4a2f0d8b6e4c2a0f8d6b4e2c0a8f6d4b2e0c8a6f4d2b0e8c"
	resolveMiningWait            = 5 * time.Minute
)

var resolvePartialMerkleTree = []byte{0x0a, 0x0b, 0x0c, 0x0d}

func resolveClaimerPayout() *entities.Wei { return entities.NewUWei(5000000000000000000) }

func resolveRegistrantFee() *entities.Wei { return entities.NewWei(100000000000000) }

func claimedPegIn() rootstock.PegInClaim {
	return rootstock.PegInClaim{
		RskAddress:  claimRskAddress,
		DepositTxID: claimDepositTxID,
		BtcAddress:  claimBtcAddress,
		PegInID:     claimPegInID,
		State:       rootstock.PegInClaimClaimed,
		TxHash:      claimTxHash,
		Amount:      claimDepositAmount(),
		Fee:         claimFee(),
		GasFee:      claimGasFee(),
		Value:       claimValue(),
		CallSuccess: true,
		CreatedAt:   time.Date(2026, 11, 20, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2026, 11, 20, 10, 5, 0, 0, time.UTC),
	}
}

func claimWithState(state rootstock.PegInClaimState) rootstock.PegInClaim {
	claim := claimedPegIn()
	claim.State = state
	return claim
}

func resolvingPegIn(sentAgo time.Duration) rootstock.PegInClaim {
	claim := claimWithState(rootstock.PegInClaimResolving)
	claim.ResolveTxHash = previousResolveTxHash
	claim.UpdatedAt = time.Now().UTC().Add(-sentAgo)
	return claim
}

func resolveParams() blockchain.ResolvePegInParams {
	return blockchain.ResolvePegInParams{
		RskAddress:        claimRskAddress,
		BtcRawTx:          claimRawTx,
		PartialMerkleTree: resolvePartialMerkleTree,
		BlockHeight:       big.NewInt(resolveBlockHeight),
	}
}

func resolveReceipt() blockchain.TransactionReceipt {
	return blockchain.TransactionReceipt{TransactionHash: resolveTxHash, BlockNumber: 456}
}

func resolvedEvent(eventPegInID string) blockchain.PegInResolvedEvent {
	return blockchain.PegInResolvedEvent{
		PegInID:       eventPegInID,
		Claimer:       claimClaimer,
		Registrant:    resolveRegistrant,
		Released:      claimDepositAmount(),
		ClaimerPayout: resolveClaimerPayout(),
		RegistrantFee: resolveRegistrantFee(),
		UserPayout:    entities.NewWei(0),
	}
}

func expectedResolvedClaim(updated rootstock.PegInClaim) rootstock.PegInClaim {
	claim := claimWithState(rootstock.PegInClaimResolved)
	claim.ResolveTxHash = resolveTxHash
	claim.ClaimerPayout = resolveClaimerPayout()
	claim.RegistrantFee = resolveRegistrantFee()
	claim.UpdatedAt = updated.UpdatedAt
	return claim
}

func expectedClaimWithState(state rootstock.PegInClaimState, updated rootstock.PegInClaim) rootstock.PegInClaim {
	claim := claimWithState(state)
	claim.UpdatedAt = updated.UpdatedAt
	return claim
}

type resolvePegInFixtures struct {
	claims      *mocks.PegInClaimRepositoryMock
	pegIn       *mocks.PeginContractMock
	bridge      *mocks.BridgeMock
	btcRpc      *mocks.BtcRpcMock
	walletMutex *sync.Mutex
	useCase     *pegin.ResolvePegInUseCase
}

func newResolvePegInFixtures() *resolvePegInFixtures {
	f := &resolvePegInFixtures{
		claims:      &mocks.PegInClaimRepositoryMock{},
		pegIn:       &mocks.PeginContractMock{},
		bridge:      &mocks.BridgeMock{},
		btcRpc:      &mocks.BtcRpcMock{},
		walletMutex: &sync.Mutex{},
	}
	f.useCase = pegin.NewResolvePegInUseCase(
		f.claims,
		blockchain.RskContracts{PegIn: f.pegIn, Bridge: f.bridge},
		blockchain.Rpc{Btc: f.btcRpc},
		f.walletMutex,
		resolveMiningWait,
	)
	return f
}

func (f *resolvePegInFixtures) run() error {
	return f.useCase.Run(context.Background(), claimedPegIn())
}

func (f *resolvePegInFixtures) assertExpectations(t *testing.T) {
	f.claims.AssertExpectations(t)
	f.pegIn.AssertExpectations(t)
	f.bridge.AssertExpectations(t)
	f.btcRpc.AssertExpectations(t)
}

func (f *resolvePegInFixtures) expectDepositConfirmations(confirmations uint64) {
	f.btcRpc.On("GetTransactionInfo", claimDepositTxID).Return(blockchain.BitcoinTransactionInformation{
		Hash:          claimDepositTxID,
		Confirmations: confirmations,
		Outputs:       map[string][]*entities.Wei{claimBtcAddress: {claimDepositAmount()}},
	}, nil).Once()
}

func (f *resolvePegInFixtures) expectRequiredConfirmations() {
	f.bridge.On("GetRequiredTxConfirmations").Return(resolveRequiredConfirmations).Once()
}

func (f *resolvePegInFixtures) expectRawTransaction() {
	f.btcRpc.On("GetRawTransaction", claimDepositTxID).Return(claimRawTx, nil).Once()
}

func (f *resolvePegInFixtures) expectPartialMerkleTree() {
	f.btcRpc.On("GetPartialMerkleTree", claimDepositTxID).Return(resolvePartialMerkleTree, nil).Once()
}

func (f *resolvePegInFixtures) expectBlock() {
	f.btcRpc.On("GetTransactionBlockInfo", claimDepositTxID).
		Return(blockchain.BitcoinBlockInformation{Hash: claimBlockHash, Height: big.NewInt(resolveBlockHeight)}, nil).Once()
}

func (f *resolvePegInFixtures) expectResolveInputs() {
	f.expectDepositConfirmations(resolveRequiredConfirmations)
	f.expectRequiredConfirmations()
	f.expectRawTransaction()
	f.expectPartialMerkleTree()
	f.expectBlock()
}

func (f *resolvePegInFixtures) expectDryRun(err error) {
	f.pegIn.EXPECT().SimulateResolvePegIn(resolveParams()).Return(err).Once()
}

func (f *resolvePegInFixtures) expectPreparedResolve() *rootstock.PegInClaim {
	f.expectResolveInputs()
	f.expectDryRun(nil)
	return f.expectUpdate(nil)
}

func (f *resolvePegInFixtures) expectSend(receipt blockchain.TransactionReceipt, err error) {
	f.pegIn.EXPECT().ResolvePegIn(resolveParams()).Return(receipt, err).Once()
}

func (f *resolvePegInFixtures) expectEvent(event blockchain.PegInResolvedEvent, err error) {
	f.pegIn.EXPECT().PegInResolvedEvent(resolveReceipt()).Return(event, err).Once()
}

func (f *resolvePegInFixtures) expectUpdate(err error) *rootstock.PegInClaim {
	updated := &rootstock.PegInClaim{}
	f.claims.EXPECT().Update(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, claim rootstock.PegInClaim) error {
			*updated = claim
			return err
		}).Once()
	return updated
}

func assertClaimUpdatedNow(t *testing.T, updated rootstock.PegInClaim) {
	assert.WithinDuration(t, time.Now().UTC(), updated.UpdatedAt, time.Minute)
}

func TestResolvePegInUseCase_Run_ResolvesAClaimAtTheBridgeDepth(t *testing.T) {
	f := newResolvePegInFixtures()
	resolving := f.expectPreparedResolve()
	f.expectSend(resolveReceipt(), nil)
	f.expectEvent(resolvedEvent(claimPegInID), nil)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInResolved(claimRskAddress, claimDepositTxID, resolveTxHash, claimPegInID))()

	require.NoError(t, f.run())
	assert.Equal(t, expectedClaimWithState(rootstock.PegInClaimResolving, *resolving), *resolving)
	assertClaimUpdatedNow(t, *resolving)
	assert.Equal(t, expectedResolvedClaim(*updated), *updated)
	assertClaimUpdatedNow(t, *updated)
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_DoesNotSendWhenTheResolvingStateCannotBeStored(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectResolveInputs()
	f.expectDryRun(nil)
	f.expectUpdate(assert.AnError)

	err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.ResolvePegInId))
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_RejectsAClaimThatIsNeitherClaimedNorResolving(t *testing.T) {
	for _, state := range []rootstock.PegInClaimState{
		rootstock.PegInClaimSubmitting,
		rootstock.PegInClaimRaceLost,
		rootstock.PegInClaimRetryableFailure,
		rootstock.PegInClaimResolved,
		rootstock.PegInClaimResolveFailed,
	} {
		t.Run(string(state), func(t *testing.T) {
			f := newResolvePegInFixtures()

			err := f.useCase.Run(context.Background(), claimWithState(state))
			require.ErrorIs(t, err, pegin.PegInClaimNotResolvableError)
			require.ErrorContains(t, err, string(usecases.ResolvePegInId))
			require.ErrorContains(t, err, string(state))
			f.assertExpectations(t)
		})
	}
}

func TestResolvePegInUseCase_Run_WaitsForTheBridgeDepth(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f := newResolvePegInFixtures()
	f.expectDepositConfirmations(resolveRequiredConfirmations - 1)
	f.expectRequiredConfirmations()
	defer test.AssertLogContains(t, pegin.LogPegInResolveWaitingForConfirmations(
		claimRskAddress, claimDepositTxID, resolveRequiredConfirmations-1, resolveRequiredConfirmations,
	))()

	require.NoError(t, f.run())
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_SkipsWhileTheDryRunSaysNotYet(t *testing.T) {
	for _, dryRunErr := range []error{blockchain.WaitingForBridgeError, blockchain.ContractPausedError} {
		t.Run(dryRunErr.Error(), func(t *testing.T) {
			log.SetLevel(log.DebugLevel)
			f := newResolvePegInFixtures()
			f.expectResolveInputs()
			f.expectDryRun(dryRunErr)
			defer test.AssertLogContains(t, pegin.LogPegInResolveNotResolvable(claimRskAddress, claimDepositTxID, dryRunErr))()

			require.NoError(t, f.run())
			f.assertExpectations(t)
		})
	}
}

func TestResolvePegInUseCase_Run_RecordsAPegInAlreadyResolvedInTheDryRun(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectResolveInputs()
	f.expectDryRun(blockchain.PegInAlreadyProcessedError)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInAlreadyResolved(claimRskAddress, claimDepositTxID))()

	require.NoError(t, f.run())
	assert.Equal(t, expectedClaimWithState(rootstock.PegInClaimResolved, *updated), *updated)
	assertClaimUpdatedNow(t, *updated)
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_MarksABridgeRejectionInTheDryRunAsFailed(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectResolveInputs()
	f.expectDryRun(blockchain.BridgeRejectedPegInError)
	updated := f.expectUpdate(nil)

	err := f.run()
	require.ErrorIs(t, err, blockchain.BridgeRejectedPegInError)
	require.ErrorIs(t, err, usecases.NonRecoverableError)
	require.ErrorContains(t, err, string(usecases.ResolvePegInId))
	assert.Equal(t, expectedClaimWithState(rootstock.PegInClaimResolveFailed, *updated), *updated)
	assertClaimUpdatedNow(t, *updated)
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_ReturnsAnUnexpectedDryRunError(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectResolveInputs()
	f.expectDryRun(assert.AnError)

	err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.ResolvePegInId))
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_RecordsAPegInResolvedMeanwhileAtMining(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectPreparedResolve()
	f.expectSend(resolveReceipt(), errors.Join(blockchain.TxFailedError, blockchain.PegInAlreadyProcessedError))
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInAlreadyResolved(claimRskAddress, claimDepositTxID))()

	require.NoError(t, f.run())
	assert.Equal(t, expectedClaimWithState(rootstock.PegInClaimResolved, *updated), *updated)
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_MarksABridgeRejectionAtMiningAsFailed(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectPreparedResolve()
	f.expectSend(resolveReceipt(), errors.Join(blockchain.TxFailedError, blockchain.BridgeRejectedPegInError))
	updated := f.expectUpdate(nil)

	err := f.run()
	require.ErrorIs(t, err, blockchain.TxFailedError)
	require.ErrorIs(t, err, blockchain.BridgeRejectedPegInError)
	require.ErrorIs(t, err, usecases.NonRecoverableError)
	assert.Equal(t, expectedClaimWithState(rootstock.PegInClaimResolveFailed, *updated), *updated)
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_KeepsTheClaimResolvingWhenTheSendDoesNotSettleIt(t *testing.T) {
	cases := []struct {
		name       string
		receipt    blockchain.TransactionReceipt
		sendErr    error
		storedHash string
	}{
		{name: "not broadcast", receipt: blockchain.TransactionReceipt{}, sendErr: assert.AnError, storedHash: ""},
		{name: "unknown outcome", receipt: blockchain.TransactionReceipt{TransactionHash: resolveTxHash}, sendErr: assert.AnError, storedHash: resolveTxHash},
		{name: "bridge not ready at mining", receipt: resolveReceipt(), sendErr: errors.Join(blockchain.TxFailedError, blockchain.WaitingForBridgeError), storedHash: resolveTxHash},
		{name: "paused at mining", receipt: resolveReceipt(), sendErr: errors.Join(blockchain.TxFailedError, blockchain.ContractPausedError), storedHash: resolveTxHash},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			f := newResolvePegInFixtures()
			f.expectPreparedResolve()
			f.expectSend(testCase.receipt, testCase.sendErr)
			updated := f.expectUpdate(nil)
			defer test.AssertNoLog(t)()

			err := f.run()
			require.ErrorIs(t, err, testCase.sendErr)
			require.ErrorContains(t, err, string(usecases.ResolvePegInId))
			require.NotErrorIs(t, err, usecases.NonRecoverableError)
			expected := expectedClaimWithState(rootstock.PegInClaimResolving, *updated)
			expected.ResolveTxHash = testCase.storedHash
			assert.Equal(t, expected, *updated)
			assertClaimUpdatedNow(t, *updated)
			f.assertExpectations(t)
		})
	}
}

func TestResolvePegInUseCase_Run_WaitsTheMiningWaitBeforeCheckingAResolvingClaim(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f := newResolvePegInFixtures()
	defer test.AssertLogContains(t, pegin.LogPegInResolvePending(claimRskAddress, claimDepositTxID, previousResolveTxHash))()

	require.NoError(t, f.useCase.Run(context.Background(), resolvingPegIn(resolveMiningWait-time.Minute)))
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_RecordsAResolvingClaimThatWasSettledMeanwhile(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectResolveInputs()
	f.expectDryRun(blockchain.PegInAlreadyProcessedError)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInAlreadyResolved(claimRskAddress, claimDepositTxID))()

	require.NoError(t, f.useCase.Run(context.Background(), resolvingPegIn(resolveMiningWait+time.Minute)))
	assert.Equal(t, expectedClaimWithState(rootstock.PegInClaimResolved, *updated), *updated)
	assertClaimUpdatedNow(t, *updated)
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_ResendsAResolveThatWasNotMinedInTheWait(t *testing.T) {
	f := newResolvePegInFixtures()
	resolving := f.expectPreparedResolve()
	f.expectSend(resolveReceipt(), nil)
	f.expectEvent(resolvedEvent(claimPegInID), nil)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInResolveResent(claimRskAddress, claimDepositTxID, previousResolveTxHash))()

	require.NoError(t, f.useCase.Run(context.Background(), resolvingPegIn(resolveMiningWait+time.Minute)))
	expectedResolving := expectedClaimWithState(rootstock.PegInClaimResolving, *resolving)
	expectedResolving.ResolveTxHash = previousResolveTxHash
	assert.Equal(t, expectedResolving, *resolving)
	assertClaimUpdatedNow(t, *resolving)
	assert.Equal(t, expectedResolvedClaim(*updated), *updated)
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_StoresTheLastBroadcastHashWhenAResendDoesNotSettle(t *testing.T) {
	cases := []struct {
		name       string
		receipt    blockchain.TransactionReceipt
		storedHash string
	}{
		{name: "not broadcast", receipt: blockchain.TransactionReceipt{}, storedHash: previousResolveTxHash},
		{name: "broadcast", receipt: blockchain.TransactionReceipt{TransactionHash: resolveTxHash}, storedHash: resolveTxHash},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			f := newResolvePegInFixtures()
			f.expectPreparedResolve()
			f.expectSend(testCase.receipt, assert.AnError)
			updated := f.expectUpdate(nil)

			err := f.useCase.Run(context.Background(), resolvingPegIn(resolveMiningWait+time.Minute))
			require.ErrorIs(t, err, assert.AnError)
			expected := expectedClaimWithState(rootstock.PegInClaimResolving, *updated)
			expected.ResolveTxHash = testCase.storedHash
			assert.Equal(t, expected, *updated)
			f.assertExpectations(t)
		})
	}
}

func TestResolvePegInUseCase_Run_StoresAResolvedClaimWithAnUnreadableEvent(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectPreparedResolve()
	f.expectSend(resolveReceipt(), nil)
	f.expectEvent(blockchain.PegInResolvedEvent{}, assert.AnError)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInResolveEventUnreadable(claimDepositTxID, resolveTxHash, assert.AnError))()

	require.NoError(t, f.run())
	resolved := expectedClaimWithState(rootstock.PegInClaimResolved, *updated)
	resolved.ResolveTxHash = resolveTxHash
	assert.Equal(t, resolved, *updated)
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_KeepsTheLocalPegInIdWhenTheEventDiffers(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectPreparedResolve()
	const eventPegInID = "00000000000000000000000000000000000000000000000000000000000000bb"
	f.expectSend(resolveReceipt(), nil)
	f.expectEvent(resolvedEvent(eventPegInID), nil)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInResolveIdMismatch(claimDepositTxID, claimPegInID, eventPegInID))()

	require.NoError(t, f.run())
	assert.Equal(t, expectedResolvedClaim(*updated), *updated)
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_HoldsTheWalletMutexOnlyForTheSend(t *testing.T) {
	f := newResolvePegInFixtures()
	f.expectResolveInputs()
	f.pegIn.EXPECT().SimulateResolvePegIn(resolveParams()).RunAndReturn(func(blockchain.ResolvePegInParams) error {
		assertMutexFree(t, f.walletMutex, "the dry run runs without the wallet mutex")
		return nil
	}).Once()
	f.pegIn.EXPECT().ResolvePegIn(resolveParams()).
		RunAndReturn(func(blockchain.ResolvePegInParams) (blockchain.TransactionReceipt, error) {
			assertMutexHeld(t, f.walletMutex, "the send runs under the wallet mutex")
			return resolveReceipt(), nil
		}).Once()
	f.expectEvent(resolvedEvent(claimPegInID), nil)
	f.claims.EXPECT().Update(mock.Anything, mock.Anything).RunAndReturn(func(context.Context, rootstock.PegInClaim) error {
		assertMutexFree(t, f.walletMutex, "the claim updates run without the wallet mutex")
		return nil
	}).Twice()

	require.NoError(t, f.run())
	f.assertExpectations(t)
}

func assertMutexFree(t *testing.T, mutex *sync.Mutex, message string) {
	acquired := mutex.TryLock()
	if acquired {
		mutex.Unlock()
	}
	assert.True(t, acquired, message)
}

func TestResolvePegInUseCase_Run_StoresTheResultAfterTheCallerContextEnds(t *testing.T) {
	f := newResolvePegInFixtures()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	f.expectPreparedResolve()
	f.pegIn.EXPECT().ResolvePegIn(resolveParams()).
		RunAndReturn(func(blockchain.ResolvePegInParams) (blockchain.TransactionReceipt, error) {
			cancel()
			return resolveReceipt(), nil
		}).Once()
	f.expectEvent(resolvedEvent(claimPegInID), nil)
	f.claims.EXPECT().Update(mock.Anything, mock.Anything).
		RunAndReturn(func(updateCtx context.Context, claim rootstock.PegInClaim) error {
			require.NoError(t, updateCtx.Err())
			assert.Equal(t, rootstock.PegInClaimResolved, claim.State)
			return nil
		}).Once()

	require.NoError(t, f.useCase.Run(ctx, claimedPegIn()))
	f.assertExpectations(t)
}

func TestResolvePegInUseCase_Run_ReturnsTheUpdateError(t *testing.T) {
	cases := []struct {
		name   string
		expect func(f *resolvePegInFixtures)
	}{
		{
			name: "after a resolve",
			expect: func(f *resolvePegInFixtures) {
				f.expectPreparedResolve()
				f.expectSend(resolveReceipt(), nil)
				f.expectEvent(resolvedEvent(claimPegInID), nil)
			},
		},
		{
			name: "after a send with an unknown outcome",
			expect: func(f *resolvePegInFixtures) {
				f.expectPreparedResolve()
				f.expectSend(resolveReceipt(), blockchain.TxFailedError)
			},
		},
		{
			name: "after a dry run that finds the peg-in resolved",
			expect: func(f *resolvePegInFixtures) {
				f.expectResolveInputs()
				f.expectDryRun(blockchain.PegInAlreadyProcessedError)
			},
		},
		{
			name: "after a dry run that finds the peg-in rejected",
			expect: func(f *resolvePegInFixtures) {
				f.expectResolveInputs()
				f.expectDryRun(blockchain.BridgeRejectedPegInError)
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			f := newResolvePegInFixtures()
			testCase.expect(f)
			f.expectUpdate(assert.AnError)

			err := f.run()
			require.ErrorIs(t, err, assert.AnError)
			require.ErrorContains(t, err, string(usecases.ResolvePegInId))
			require.NotErrorIs(t, err, usecases.NonRecoverableError)
			f.assertExpectations(t)
		})
	}
}

func TestResolvePegInUseCase_Run_ReturnsTheErrorOfEachStep(t *testing.T) {
	steps := []struct {
		name    string
		succeed func(f *resolvePegInFixtures)
		fail    func(f *resolvePegInFixtures)
	}{
		{
			name: "deposit read",
			succeed: func(f *resolvePegInFixtures) {
				f.expectDepositConfirmations(resolveRequiredConfirmations)
				f.expectRequiredConfirmations()
			},
			fail: func(f *resolvePegInFixtures) {
				f.btcRpc.On("GetTransactionInfo", claimDepositTxID).Return(blockchain.BitcoinTransactionInformation{}, assert.AnError).Once()
			},
		},
		{
			name:    "raw transaction read",
			succeed: (*resolvePegInFixtures).expectRawTransaction,
			fail: func(f *resolvePegInFixtures) {
				f.btcRpc.On("GetRawTransaction", claimDepositTxID).Return(nil, assert.AnError).Once()
			},
		},
		{
			name:    "partial merkle tree read",
			succeed: (*resolvePegInFixtures).expectPartialMerkleTree,
			fail: func(f *resolvePegInFixtures) {
				f.btcRpc.On("GetPartialMerkleTree", claimDepositTxID).Return([]byte(nil), assert.AnError).Once()
			},
		},
		{
			name:    "block read",
			succeed: (*resolvePegInFixtures).expectBlock,
			fail: func(f *resolvePegInFixtures) {
				f.btcRpc.On("GetTransactionBlockInfo", claimDepositTxID).Return(blockchain.BitcoinBlockInformation{}, assert.AnError).Once()
			},
		},
	}
	for i, step := range steps {
		t.Run(step.name, func(t *testing.T) {
			f := newResolvePegInFixtures()
			for _, previous := range steps[:i] {
				previous.succeed(f)
			}
			step.fail(f)

			err := f.run()
			require.ErrorIs(t, err, assert.AnError)
			require.ErrorContains(t, err, string(usecases.ResolvePegInId))
			f.assertExpectations(t)
		})
	}
}
