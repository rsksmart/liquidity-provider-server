package pegin_test

import (
	"context"
	"math/big"
	"testing"
	"time"

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

const recheckOtherDepositTxID = "5e8d0c2b7a6f4e3d2c1b0a9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c3b2a1f0e9d"

type recheckResolvedFixtures struct {
	claims  *mocks.PegInClaimRepositoryMock
	pegIn   *mocks.PeginContractMock
	btcRpc  *mocks.BitcoinNetworkMock
	useCase *pegin.RecheckResolvedPegInsUseCase
}

func newRecheckResolvedFixtures(t *testing.T) *recheckResolvedFixtures {
	f := &recheckResolvedFixtures{
		claims: mocks.NewPegInClaimRepositoryMock(t),
		pegIn:  mocks.NewPeginContractMock(t),
		btcRpc: mocks.NewBitcoinNetworkMock(t),
	}
	f.useCase = pegin.NewRecheckResolvedPegInsUseCase(
		f.claims,
		blockchain.RskContracts{PegIn: f.pegIn},
		blockchain.Rpc{Btc: f.btcRpc},
	)
	return f
}

func resolvedPegIn() rootstock.PegInClaim {
	claim := claimWithState(rootstock.PegInClaimResolved)
	claim.ResolveTxHash = resolveTxHash
	claim.ClaimerPayout = resolveClaimerPayout()
	claim.RegistrantFee = resolveRegistrantFee()
	return claim
}

func reopenedPegIn(updated rootstock.PegInClaim) rootstock.PegInClaim {
	claim := claimWithState(rootstock.PegInClaimResolving)
	claim.ResolveTxHash = resolveTxHash
	claim.UpdatedAt = updated.UpdatedAt
	return claim
}

func (f *recheckResolvedFixtures) expectRecentlyResolved(claims ...rootstock.PegInClaim) {
	windowStart := mock.MatchedBy(func(since time.Time) bool {
		return since.Sub(time.Now().UTC().Add(-pegin.ResolvedPegInRecheckWindow)).Abs() < time.Minute
	})
	f.claims.EXPECT().ListByStateUpdatedSince(mock.Anything, rootstock.PegInClaimResolved, windowStart).Return(claims, nil).Once()
}

func (f *recheckResolvedFixtures) expectDryRunInputs() {
	f.btcRpc.EXPECT().GetRawTransaction(claimDepositTxID).Return(claimRawTx, nil).Once()
	f.btcRpc.EXPECT().GetPartialMerkleTree(claimDepositTxID).Return(resolvePartialMerkleTree, nil).Once()
	f.btcRpc.EXPECT().GetTransactionBlockInfo(claimDepositTxID).
		Return(blockchain.BitcoinBlockInformation{Hash: claimBlockHash, Height: big.NewInt(resolveBlockHeight)}, nil).Once()
}

func (f *recheckResolvedFixtures) expectDryRun(err error) {
	f.pegIn.EXPECT().SimulateResolvePegIn(resolveParams()).Return(err).Once()
}

func (f *recheckResolvedFixtures) expectUpdate(err error) *rootstock.PegInClaim {
	updated := &rootstock.PegInClaim{}
	f.claims.EXPECT().Update(mock.Anything, mock.Anything).
		RunAndReturn(func(_ context.Context, claim rootstock.PegInClaim) error {
			*updated = claim
			return err
		}).Once()
	return updated
}

func TestRecheckResolvedPegInsUseCase_Run_LeavesAPegInThatIsStillResolvedOnChain(t *testing.T) {
	f := newRecheckResolvedFixtures(t)
	f.expectRecentlyResolved(resolvedPegIn())
	f.expectDryRunInputs()
	f.expectDryRun(blockchain.PegInAlreadyProcessedError)
	defer test.AssertNoLog(t)()

	require.NoError(t, f.useCase.Run(context.Background()))
}

func TestRecheckResolvedPegInsUseCase_Run_ReopensAPegInThatIsUnresolvedOnChain(t *testing.T) {
	cases := []struct {
		name      string
		dryRunErr error
	}{
		{name: "dry run passes", dryRunErr: nil},
		{name: "bridge not ready", dryRunErr: blockchain.WaitingForBridgeError},
		{name: "bridge rejects", dryRunErr: blockchain.BridgeRejectedPegInError},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			f := newRecheckResolvedFixtures(t)
			f.expectRecentlyResolved(resolvedPegIn())
			f.expectDryRunInputs()
			f.expectDryRun(testCase.dryRunErr)
			updated := f.expectUpdate(nil)
			defer test.AssertLogContains(t, pegin.LogPegInResolveReopened(claimRskAddress, claimDepositTxID, resolveTxHash))()

			require.NoError(t, f.useCase.Run(context.Background()))
			assert.Equal(t, reopenedPegIn(*updated), *updated)
			assertClaimUpdatedNow(t, *updated)
		})
	}
}

func TestRecheckResolvedPegInsUseCase_Run_LeavesThePegInWhileTheContractIsPaused(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	f := newRecheckResolvedFixtures(t)
	f.expectRecentlyResolved(resolvedPegIn())
	f.expectDryRunInputs()
	f.expectDryRun(blockchain.ContractPausedError)
	defer test.AssertLogContains(t, pegin.LogPegInResolveRecheckSkipped(claimRskAddress, claimDepositTxID, blockchain.ContractPausedError))()

	require.NoError(t, f.useCase.Run(context.Background()))
}

func TestRecheckResolvedPegInsUseCase_Run_LogsAnUnexpectedDryRunError(t *testing.T) {
	f := newRecheckResolvedFixtures(t)
	f.expectRecentlyResolved(resolvedPegIn())
	f.expectDryRunInputs()
	f.expectDryRun(assert.AnError)
	defer test.AssertLogContains(t, pegin.LogPegInResolveRecheckError(claimRskAddress, claimDepositTxID, assert.AnError))()

	require.NoError(t, f.useCase.Run(context.Background()))
}

func TestRecheckResolvedPegInsUseCase_Run_LogsTheUpdateError(t *testing.T) {
	f := newRecheckResolvedFixtures(t)
	f.expectRecentlyResolved(resolvedPegIn())
	f.expectDryRunInputs()
	f.expectDryRun(nil)
	f.expectUpdate(assert.AnError)
	defer test.AssertLogContains(t, pegin.LogPegInResolveRecheckError(claimRskAddress, claimDepositTxID, assert.AnError))()

	require.NoError(t, f.useCase.Run(context.Background()))
}

func TestRecheckResolvedPegInsUseCase_Run_KeepsCheckingTheOtherPegInsWhenOneFails(t *testing.T) {
	f := newRecheckResolvedFixtures(t)
	failing := resolvedPegIn()
	failing.DepositTxID = recheckOtherDepositTxID
	f.expectRecentlyResolved(failing, resolvedPegIn())
	f.btcRpc.EXPECT().GetRawTransaction(recheckOtherDepositTxID).Return(nil, assert.AnError).Once()
	f.expectDryRunInputs()
	f.expectDryRun(nil)
	updated := f.expectUpdate(nil)
	defer test.AssertLogContains(t, pegin.LogPegInResolveRecheckError(claimRskAddress, recheckOtherDepositTxID, assert.AnError))()

	require.NoError(t, f.useCase.Run(context.Background()))
	assert.Equal(t, reopenedPegIn(*updated), *updated)
}

func TestRecheckResolvedPegInsUseCase_Run_ReturnsTheListingError(t *testing.T) {
	f := newRecheckResolvedFixtures(t)
	f.claims.EXPECT().ListByStateUpdatedSince(mock.Anything, rootstock.PegInClaimResolved, mock.Anything).Return(nil, assert.AnError).Once()

	err := f.useCase.Run(context.Background())
	require.ErrorIs(t, err, assert.AnError)
	require.ErrorContains(t, err, string(usecases.RecheckResolvedPegInsId))
}
