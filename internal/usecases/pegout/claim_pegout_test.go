package pegout_test

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	claimRequestHash = "aabbccddeeff00112233445566778899aabbccddeeff00112233445566778899"
	claimQuoteHash   = "1122334455667788990011223344556677889900112233445566778899001122"
	claimLpAddress   = "0x7c4890a0f1d4bbf2c669ac2d1effa185c505359b"
	claimEncodedBtc  = "bcrt1qdest"
	claimTxHash      = "0xclaimtx"
	claimDateLimit   = uint32(1_700_000_000)
)

func claimEscrowQuote() quote.PegoutQuote {
	return quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          blockchain.RskZeroAddress,
		BtcRefundAddress:      "6f0011223344556677889900112233445566778899",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "6f0011223344556677889900112233445566778899",
		CallFee:               entities.NewWei(10_000_000_000_000_000),
		PenaltyFee:            entities.NewWei(1),
		Nonce:                 1,
		DepositAddress:        "6f0011223344556677889900112233445566778899",
		Value:                 entities.NewWei(1_000_000),
		AgreementTimestamp:    1,
		DepositDateLimit:      claimDateLimit,
		DepositConfirmations:  0,
		TransferConfirmations: 1,
		TransferTime:          600,
		ExpireDate:            2000000000,
		ExpireBlock:           999999,
		GasFee:                entities.NewWei(1),
	}
}

func claimCompletedQuote() quote.PegoutQuote {
	completed := claimEscrowQuote()
	completed.LpRskAddress = claimLpAddress
	return completed
}

func isCompletedClaimQuote(q quote.PegoutQuote) bool {
	return q.LpRskAddress == claimLpAddress && q.DepositDateLimit == claimDateLimit
}

type claimFixtures struct {
	escrow    *mocks.PegOutEscrowContractMock
	pegout    *mocks.PegoutContractMock
	repo      *mocks.PegoutQuoteRepositoryMock
	btcWallet *mocks.BitcoinWalletMock
	btcRpc    *mocks.BtcRpcMock
	rskRpc    *mocks.RootstockRpcServerMock
	lp        *mocks.ProviderMock
	signer    *mocks.TransactionSignerMock
	eventBus  *mocks.EventBusMock
	useCase   *pegout.ClaimPegOutUseCase
	calls     []string
}

func newClaimFixtures() *claimFixtures {
	f := &claimFixtures{
		escrow:    &mocks.PegOutEscrowContractMock{},
		pegout:    &mocks.PegoutContractMock{},
		repo:      &mocks.PegoutQuoteRepositoryMock{},
		btcWallet: &mocks.BitcoinWalletMock{},
		btcRpc:    &mocks.BtcRpcMock{},
		rskRpc:    &mocks.RootstockRpcServerMock{},
		lp:        &mocks.ProviderMock{},
		signer:    &mocks.TransactionSignerMock{},
		eventBus:  &mocks.EventBusMock{},
	}
	f.useCase = pegout.NewClaimPegOutUseCase(
		blockchain.RskContracts{PegOut: f.pegout, PegOutEscrow: f.escrow},
		blockchain.Rpc{Rsk: f.rskRpc, Btc: f.btcRpc},
		f.btcWallet,
		f.lp,
		f.repo,
		f.eventBus,
		&sync.Mutex{},
	)
	f.lp.On("RskAddress").Return(claimLpAddress)
	return f
}

func (f *claimFixtures) record(name string) func(mock.Arguments) {
	return func(mock.Arguments) { f.calls = append(f.calls, name) }
}

func (f *claimFixtures) run() (pegout.ClaimOutcome, error) {
	return f.useCase.Run(context.Background(), blockchain.PegOutRequested{RequestHash: claimRequestHash, Amount: entities.NewWei(1_000_000)})
}

func (f *claimFixtures) expectPreChecks() {
	f.pegout.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	f.repo.On("GetRetainedQuote", mock.Anything, claimRequestHash).Return(nil, nil).Once()
	f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()
	f.escrow.EXPECT().RestrictedUntil(claimLpAddress).Return(uint64(0), nil).Once()
	f.rskRpc.EXPECT().GetHeight(mock.Anything).Return(uint64(100), nil).Once()
}

func (f *claimFixtures) expectSigned() {
	f.expectSignedQuote(claimEscrowQuote())
}

func (f *claimFixtures) expectSignedQuote(q quote.PegoutQuote) {
	f.expectPreChecks()
	f.escrow.EXPECT().GetPegOutQuote(claimRequestHash).Return(q, nil).Once()
	f.btcRpc.On("EncodeAddress", mock.Anything).Return(claimEncodedBtc, nil).Times(3)
	f.btcWallet.On("GetBalance").Return(entities.NewWei(10_000_000), nil).Once()
	f.repo.On("GetRetainedQuoteByState", mock.Anything,
		quote.PegoutStateClaimPending, quote.PegoutStateClaimed, quote.PegoutStateWaitingForDepositConfirmations,
	).Return([]quote.RetainedPegoutQuote{}, nil).Once()
	eip712Hash := [32]byte{9, 8, 7}
	f.pegout.EXPECT().HashPegoutQuoteEIP712(mock.Anything).Return(eip712Hash, nil).Once()
	f.lp.On("GetSigner").Return(f.signer).Once()
	f.signer.On("SignBytes", eip712Hash[:]).Return(make([]byte, 65), nil).Once()
}

func (f *claimFixtures) expectProfitable() {
	f.expectSigned()
	f.escrow.EXPECT().EstimateClaimPegOut(claimRequestHash, mock.Anything).Return(entities.NewWei(100_000), nil).Once()
	f.btcWallet.On("EstimateTxFees", claimEncodedBtc, entities.NewWei(1_000_000)).
		Return(blockchain.BtcFeeEstimation{Value: entities.NewWei(1)}, nil).Once()
	f.rskRpc.EXPECT().GasPrice(mock.Anything).Return(entities.NewWei(1), nil).Once()
}

func (f *claimFixtures) expectQuoteHash() {
	f.pegout.EXPECT().HashPegoutQuote(mock.MatchedBy(isCompletedClaimQuote)).Return(claimQuoteHash, nil).Once()
}

func (f *claimFixtures) expectPendingRecorded() {
	f.expectProfitable()
	f.expectQuoteHash()
	f.repo.On("GetRetainedQuote", mock.Anything, claimQuoteHash).Return(nil, nil).Once()
	f.pegout.EXPECT().GetAddress().Return("0xpegout").Once()
	f.repo.On("InsertQuote", mock.Anything, mock.MatchedBy(func(c quote.CreatedPegoutQuote) bool {
		return c.Hash == claimQuoteHash && isCompletedClaimQuote(c.Quote)
	})).Return(nil).Once().Run(f.record("InsertQuote"))
	f.repo.On("InsertRetainedQuote", mock.Anything, mock.MatchedBy(func(r quote.RetainedPegoutQuote) bool {
		return r.QuoteHash == claimQuoteHash && r.State == quote.PegoutStateClaimPending
	})).Return(nil).Once().Run(f.record("InsertRetainedQuote"))
}

func (f *claimFixtures) expectClaim(receipt blockchain.TransactionReceipt, err error) {
	f.escrow.EXPECT().ClaimPegOut(mock.Anything, claimRequestHash, mock.Anything).
		Run(func(blockchain.TransactionConfig, string, []byte) { f.calls = append(f.calls, "ClaimPegOut") }).
		Return(receipt, err).Once()
}

func (f *claimFixtures) expectPromote(txHash string) {
	f.repo.On("UpdateRetainedQuote", mock.Anything, mock.MatchedBy(func(r quote.RetainedPegoutQuote) bool {
		return r.QuoteHash == claimQuoteHash && r.State == quote.PegoutStateClaimed && r.UserRskTxHash == txHash
	})).Return(nil).Once().Run(f.record("UpdateRetainedQuote"))
	f.eventBus.On("Publish", mock.MatchedBy(func(e quote.ClaimedPegoutQuoteEvent) bool {
		return e.RetainedQuote.QuoteHash == claimQuoteHash && e.RetainedQuote.State == quote.PegoutStateClaimed && isCompletedClaimQuote(e.Quote)
	})).Return().Once()
}

func (f *claimFixtures) expectDelete() {
	f.repo.On("DeleteQuotes", mock.Anything, []string{claimQuoteHash}).Return(uint(1), nil).Once()
}

func (f *claimFixtures) expectBlockTime(timestamp int64, err error) {
	f.rskRpc.EXPECT().GetBlockByNumber(mock.Anything, (*big.Int)(nil)).
		Run(func(context.Context, *big.Int) { f.calls = append(f.calls, "GetBlockByNumber") }).
		Return(blockchain.BlockInfo{Timestamp: time.Unix(timestamp, 0)}, err).Once()
}

func TestClaimPegOutUseCase_NoEscrow(t *testing.T) {
	useCase := pegout.NewClaimPegOutUseCase(blockchain.RskContracts{}, blockchain.Rpc{}, nil, nil, nil, nil, &sync.Mutex{})
	outcome, err := useCase.Run(context.Background(), blockchain.PegOutRequested{RequestHash: claimRequestHash})
	require.NoError(t, err)
	assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	require.NoError(t, useCase.ReconcilePendingClaims(context.Background()))
}

func TestClaimPegOutUseCase_ClaimsUnderCompletedQuoteHash(t *testing.T) {
	f := newClaimFixtures()
	f.expectPendingRecorded()
	f.expectClaim(blockchain.TransactionReceipt{TransactionHash: claimTxHash}, nil)
	f.expectPromote(claimTxHash)

	checkLog := test.LogContains(t, pegout.LogClaimPegoutSuccess(claimRequestHash, claimQuoteHash, claimTxHash))
	outcome, err := f.run()
	require.NoError(t, err)
	assert.Equal(t, pegout.ClaimOutcomeClaimed, outcome)
	assert.True(t, checkLog())
	assert.Equal(t, []string{"InsertQuote", "InsertRetainedQuote", "ClaimPegOut", "UpdateRetainedQuote"}, f.calls)
	f.escrow.AssertNotCalled(t, "GetPegOutState", claimQuoteHash)
	f.repo.AssertNotCalled(t, "DeleteQuotes", mock.Anything, mock.Anything)
	f.repo.AssertExpectations(t)
	f.pegout.AssertExpectations(t)
	f.eventBus.AssertExpectations(t)
}

func TestClaimPegOutUseCase_SkipsWhenCapacityIsInsufficient(t *testing.T) {
	f := newClaimFixtures()
	f.expectPreChecks()
	f.escrow.EXPECT().GetPegOutQuote(claimRequestHash).Return(claimEscrowQuote(), nil).Once()
	f.btcRpc.On("EncodeAddress", mock.Anything).Return(claimEncodedBtc, nil).Times(3)
	f.btcWallet.On("GetBalance").Return(entities.NewWei(100), nil).Once()
	f.repo.On("GetRetainedQuoteByState", mock.Anything,
		quote.PegoutStateClaimPending, quote.PegoutStateClaimed, quote.PegoutStateWaitingForDepositConfirmations,
	).Return([]quote.RetainedPegoutQuote{}, nil).Once()

	log.SetLevel(log.DebugLevel)
	t.Cleanup(func() { log.SetLevel(log.InfoLevel) })
	checkLog := test.LogContains(t, pegout.LogClaimPegoutCapacitySkip(claimRequestHash))
	outcome, err := f.run()
	require.NoError(t, err)
	assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	assert.True(t, checkLog())
	f.escrow.AssertNotCalled(t, "EstimateClaimPegOut", mock.Anything, mock.Anything)
	f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
}

func TestClaimPegOutUseCase_SkipsWhenUnprofitable(t *testing.T) {
	f := newClaimFixtures()
	q := claimEscrowQuote()
	q.CallFee = entities.NewWei(1)
	f.expectSignedQuote(q)
	f.escrow.EXPECT().EstimateClaimPegOut(claimRequestHash, mock.Anything).Return(entities.NewWei(100_000), nil).Once()
	f.btcWallet.On("EstimateTxFees", claimEncodedBtc, q.Value).Return(blockchain.BtcFeeEstimation{Value: entities.NewWei(50_000)}, nil).Once()
	f.rskRpc.EXPECT().GasPrice(mock.Anything).Return(entities.NewWei(100), nil).Once()

	log.SetLevel(log.DebugLevel)
	t.Cleanup(func() { log.SetLevel(log.InfoLevel) })
	checkLog := test.LogContains(t, pegout.LogClaimPegoutProfitabilitySkip(claimRequestHash))
	outcome, err := f.run()
	require.NoError(t, err)
	assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	assert.True(t, checkLog())
	f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
}

func TestClaimPegOutUseCase_QuoteHashErrorStopsBeforeClaim(t *testing.T) {
	f := newClaimFixtures()
	f.expectProfitable()
	f.pegout.EXPECT().HashPegoutQuote(mock.MatchedBy(isCompletedClaimQuote)).Return("", assert.AnError).Once()

	outcome, err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
	f.repo.AssertNotCalled(t, "InsertQuote", mock.Anything, mock.Anything)
	f.repo.AssertNotCalled(t, "InsertRetainedQuote", mock.Anything, mock.Anything)
}

func TestClaimPegOutUseCase_SkipsWhenClaimAlreadyRecorded(t *testing.T) {
	f := newClaimFixtures()
	f.expectProfitable()
	f.expectQuoteHash()
	f.repo.On("GetRetainedQuote", mock.Anything, claimQuoteHash).
		Return(&quote.RetainedPegoutQuote{QuoteHash: claimQuoteHash, State: quote.PegoutStateClaimPending}, nil).Once()

	outcome, err := f.run()
	require.NoError(t, err)
	assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
	f.repo.AssertNotCalled(t, "InsertQuote", mock.Anything, mock.Anything)
	f.repo.AssertNotCalled(t, "InsertRetainedQuote", mock.Anything, mock.Anything)
}

func TestClaimPegOutUseCase_NoClaimWhenPendingRecordFails(t *testing.T) {
	f := newClaimFixtures()
	f.expectProfitable()
	f.expectQuoteHash()
	f.repo.On("GetRetainedQuote", mock.Anything, claimQuoteHash).Return(nil, nil).Once()
	f.pegout.EXPECT().GetAddress().Return("0xpegout").Once()
	f.repo.On("InsertQuote", mock.Anything, mock.Anything).Return(nil).Once()
	f.repo.On("InsertRetainedQuote", mock.Anything, mock.Anything).Return(assert.AnError).Once()
	f.expectDelete()

	outcome, err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
	f.repo.AssertExpectations(t)
}

func TestClaimPegOutUseCase_KeepsPendingClaimWhenPromotionFails(t *testing.T) {
	f := newClaimFixtures()
	f.expectPendingRecorded()
	f.expectClaim(blockchain.TransactionReceipt{TransactionHash: claimTxHash}, nil)
	f.repo.On("UpdateRetainedQuote", mock.Anything, mock.Anything).Return(assert.AnError).Once()

	outcome, err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	f.repo.AssertNotCalled(t, "DeleteQuotes", mock.Anything, mock.Anything)
	f.eventBus.AssertNotCalled(t, "Publish", mock.Anything)
}

func TestClaimPegOutUseCase_KeepsPendingClaimOnMiningTimeout(t *testing.T) {
	f := newClaimFixtures()
	f.expectPendingRecorded()
	f.expectClaim(blockchain.TransactionReceipt{}, fmt.Errorf("claim peg out error: %w", context.DeadlineExceeded))

	outcome, err := f.run()
	require.ErrorIs(t, err, context.DeadlineExceeded)
	assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	f.repo.AssertNotCalled(t, "DeleteQuotes", mock.Anything, mock.Anything)
	f.repo.AssertNotCalled(t, "UpdateRetainedQuote", mock.Anything, mock.Anything)
}

func TestClaimPegOutUseCase_DeletesPendingClaimWhenClaimFailsWithRequestOpen(t *testing.T) {
	f := newClaimFixtures()
	f.expectPendingRecorded()
	f.expectClaim(blockchain.TransactionReceipt{TransactionHash: claimTxHash}, assert.AnError)
	f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()
	f.expectDelete()

	outcome, err := f.run()
	require.ErrorIs(t, err, assert.AnError)
	assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	f.escrow.AssertNotCalled(t, "GetPegOutState", claimQuoteHash)
	f.repo.AssertExpectations(t)
}

func TestClaimPegOutUseCase_DeletesPendingClaimOnLostRace(t *testing.T) {
	f := newClaimFixtures()
	f.expectPendingRecorded()
	f.expectClaim(blockchain.TransactionReceipt{}, assert.AnError)
	f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateNone, nil).Once()
	f.escrow.EXPECT().GetPegOutState(claimQuoteHash).Return(blockchain.EscrowedPegOutStateNone, nil).Once()
	f.expectDelete()

	checkLog := test.LogContains(t, pegout.LogClaimPegoutLostRace(claimRequestHash))
	outcome, err := f.run()
	require.NoError(t, err)
	assert.Equal(t, pegout.ClaimOutcomeClosed, outcome)
	assert.True(t, checkLog())
	f.repo.AssertExpectations(t)
}

func TestClaimPegOutUseCase_PromotesWhenClaimMinedDespiteError(t *testing.T) {
	f := newClaimFixtures()
	f.expectPendingRecorded()
	f.expectClaim(blockchain.TransactionReceipt{}, assert.AnError)
	f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateNone, nil).Once()
	f.escrow.EXPECT().GetPegOutState(claimQuoteHash).Return(blockchain.EscrowedPegOutStateClaimed, nil).Once()
	f.expectPromote("")

	outcome, err := f.run()
	require.NoError(t, err)
	assert.Equal(t, pegout.ClaimOutcomeClaimed, outcome)
	f.repo.AssertNotCalled(t, "DeleteQuotes", mock.Anything, mock.Anything)
	f.eventBus.AssertExpectations(t)
}

func TestClaimPegOutUseCase_ClosedWhenRequestNoLongerOpen(t *testing.T) {
	f := newClaimFixtures()
	f.pegout.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	f.repo.On("GetRetainedQuote", mock.Anything, claimRequestHash).Return(nil, nil).Once()
	f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateNone, nil).Once()

	checkLog := test.LogContains(t, pegout.LogClaimPegoutLostRace(claimRequestHash))
	outcome, err := f.run()
	require.NoError(t, err)
	assert.Equal(t, pegout.ClaimOutcomeClosed, outcome)
	assert.True(t, checkLog())
	f.escrow.AssertNotCalled(t, "GetPegOutQuote", mock.Anything)
}

func TestClaimPegOutUseCase_ClosedWhenEstimateFindsLostRace(t *testing.T) {
	f := newClaimFixtures()
	f.expectSigned()
	f.escrow.EXPECT().EstimateClaimPegOut(claimRequestHash, mock.Anything).Return(nil, assert.AnError).Once()
	f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateClaimed, nil).Once()

	outcome, err := f.run()
	require.NoError(t, err)
	assert.Equal(t, pegout.ClaimOutcomeClosed, outcome)
	f.rskRpc.AssertNotCalled(t, "GetBlockByNumber", mock.Anything, mock.Anything)
	f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
}

func TestClaimPegOutUseCase_ClaimWindow(t *testing.T) {
	setup := func(blockTime int64, blockErr error) *claimFixtures {
		f := newClaimFixtures()
		f.expectSigned()
		f.escrow.EXPECT().EstimateClaimPegOut(claimRequestHash, mock.Anything).Return(nil, assert.AnError).Once()
		f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()
		f.expectBlockTime(blockTime, blockErr)
		return f
	}
	t.Run("closed once a block is past depositDateLimit", func(t *testing.T) {
		f := setup(int64(claimDateLimit)+1, nil)
		checkLog := test.LogContains(t, pegout.LogClaimPegoutWindowClosed(claimRequestHash))
		outcome, err := f.run()
		require.NoError(t, err)
		assert.Equal(t, pegout.ClaimOutcomeClosed, outcome)
		assert.True(t, checkLog())
	})
	t.Run("still open at depositDateLimit", func(t *testing.T) {
		f := setup(int64(claimDateLimit), nil)
		outcome, err := f.run()
		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	})
	t.Run("error reading the block", func(t *testing.T) {
		f := setup(0, assert.AnError)
		outcome, err := f.run()
		require.Error(t, err)
		assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	})
}

func TestClaimPegOutUseCase_FeeEstimation(t *testing.T) {
	setup := func(feeErr error) *claimFixtures {
		f := newClaimFixtures()
		f.expectSigned()
		f.escrow.EXPECT().EstimateClaimPegOut(claimRequestHash, mock.Anything).Return(entities.NewWei(100_000), nil).Once()
		f.btcWallet.On("EstimateTxFees", claimEncodedBtc, entities.NewWei(1_000_000)).Return(blockchain.BtcFeeEstimation{}, feeErr).Once()
		return f
	}
	t.Run("insufficient funds skips the claim", func(t *testing.T) {
		rpcErr := &btcjson.RPCError{Code: btcjson.ErrRPCWallet, Message: "Insufficient funds"}
		f := setup(fmt.Errorf("%w: %w", blockchain.BtcInsufficientFundsError, rpcErr))
		outcome, err := f.run()
		require.NoError(t, err)
		assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
		f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
	})
	t.Run("other errors are returned", func(t *testing.T) {
		f := setup(assert.AnError)
		outcome, err := f.run()
		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, pegout.ClaimOutcomeSkipped, outcome)
	})
}

func (f *claimFixtures) expectPendingClaim(blockTime int64) {
	completed := claimCompletedQuote()
	f.repo.On("GetRetainedQuoteByState", mock.Anything, quote.PegoutStateClaimPending).Return([]quote.RetainedPegoutQuote{
		{QuoteHash: claimQuoteHash, State: quote.PegoutStateClaimPending, RequiredLiquidity: entities.NewWei(1)},
	}, nil).Once()
	f.expectBlockTime(blockTime, nil)
	f.repo.On("GetQuote", mock.Anything, claimQuoteHash).Return(&completed, nil).Once()
}

func (f *claimFixtures) expectPendingState(state blockchain.EscrowedPegOutState) {
	f.escrow.EXPECT().GetPegOutState(claimQuoteHash).
		Run(func(string) { f.calls = append(f.calls, "GetPegOutState") }).
		Return(state, nil).Once()
}

func TestClaimPegOutUseCase_ReconcilePendingClaims(t *testing.T) {
	t.Run("promotes a mined claim", func(t *testing.T) {
		f := newClaimFixtures()
		f.expectPendingClaim(int64(claimDateLimit) - 10)
		f.expectPendingState(blockchain.EscrowedPegOutStateClaimed)
		f.expectPromote("")

		require.NoError(t, f.useCase.ReconcilePendingClaims(context.Background()))
		assert.Equal(t, []string{"GetBlockByNumber", "GetPegOutState", "UpdateRetainedQuote"}, f.calls)
		f.repo.AssertNotCalled(t, "DeleteQuotes", mock.Anything, mock.Anything)
		f.eventBus.AssertExpectations(t)
	})
	t.Run("keeps an unmined claim while the window is open", func(t *testing.T) {
		f := newClaimFixtures()
		f.expectPendingClaim(int64(claimDateLimit))
		f.expectPendingState(blockchain.EscrowedPegOutStateNone)

		require.NoError(t, f.useCase.ReconcilePendingClaims(context.Background()))
		f.repo.AssertNotCalled(t, "DeleteQuotes", mock.Anything, mock.Anything)
		f.repo.AssertNotCalled(t, "UpdateRetainedQuote", mock.Anything, mock.Anything)
	})
	t.Run("deletes an unmined claim once the window has closed", func(t *testing.T) {
		f := newClaimFixtures()
		f.expectPendingClaim(int64(claimDateLimit) + 1)
		f.expectPendingState(blockchain.EscrowedPegOutStateNone)
		f.expectDelete()

		require.NoError(t, f.useCase.ReconcilePendingClaims(context.Background()))
		f.repo.AssertExpectations(t)
	})
	t.Run("deletes a claim the escrow already settled", func(t *testing.T) {
		f := newClaimFixtures()
		f.expectPendingClaim(int64(claimDateLimit) - 10)
		f.expectPendingState(blockchain.EscrowedPegOutStateRefunded)
		f.expectDelete()

		require.NoError(t, f.useCase.ReconcilePendingClaims(context.Background()))
		f.repo.AssertExpectations(t)
	})
	t.Run("returns an error when the quote is missing", func(t *testing.T) {
		f := newClaimFixtures()
		f.repo.On("GetRetainedQuoteByState", mock.Anything, quote.PegoutStateClaimPending).Return([]quote.RetainedPegoutQuote{
			{QuoteHash: claimQuoteHash, State: quote.PegoutStateClaimPending},
		}, nil).Once()
		f.expectBlockTime(int64(claimDateLimit), nil)
		f.repo.On("GetQuote", mock.Anything, claimQuoteHash).Return(nil, nil).Once()

		require.Error(t, f.useCase.ReconcilePendingClaims(context.Background()))
		f.repo.AssertNotCalled(t, "DeleteQuotes", mock.Anything, mock.Anything)
	})
	t.Run("does nothing without pending claims", func(t *testing.T) {
		f := newClaimFixtures()
		f.repo.On("GetRetainedQuoteByState", mock.Anything, quote.PegoutStateClaimPending).Return([]quote.RetainedPegoutQuote{}, nil).Once()

		require.NoError(t, f.useCase.ReconcilePendingClaims(context.Background()))
		f.rskRpc.AssertNotCalled(t, "GetBlockByNumber", mock.Anything, mock.Anything)
	})
}
