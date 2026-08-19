package pegout_test

import (
	"context"
	"math/big"
	"testing"
	"time"

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
	claimLpAddress   = "0xlp"
)

func claimEscrowQuote() quote.PegoutQuote {
	// Pre-claim GetPegOutQuote shape: incomplete quote (lpRskAddress = 0), empty unused BTC fields.
	return quote.PegoutQuote{
		LbcAddress:            "0xabcd01",
		LpRskAddress:          blockchain.RskZeroAddress,
		BtcRefundAddress:      "",
		RskRefundAddress:      "0xabcd04",
		LpBtcAddress:          "",
		CallFee:               entities.NewWei(10_000_000_000_000_000), // large: profitable by default
		PenaltyFee:            entities.NewWei(1),
		Nonce:                 1,
		DepositAddress:        "ddeeff00112233445566778899aabbccddeeff00",
		Value:                 entities.NewWei(1_000_000),
		AgreementTimestamp:    1,
		DepositDateLimit:      2,
		DepositConfirmations:  1,
		TransferConfirmations: 1,
		TransferTime:          600,
		ExpireDate:            2000000000,
		ExpireBlock:           999999,
		GasFee:                entities.NewWei(1),
	}
}

func claimSignature65() []byte {
	sig := make([]byte, 65)
	for i := range sig {
		sig[i] = byte(i + 1)
	}
	return sig
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
	mutex     *mocks.MutexMock
	useCase   *pegout.ClaimPegOutUseCase
	candidate blockchain.PegOutRequested
}

func newClaimFixtures(t *testing.T) *claimFixtures {
	t.Helper()
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
		mutex:     &mocks.MutexMock{},
		candidate: blockchain.PegOutRequested{
			RequestHash: claimRequestHash,
			Amount:      entities.NewWei(1_000_000),
		},
	}
	f.pegout.EXPECT().PausedStatus().Return(blockchain.PauseStatus{IsPaused: false}, nil)
	f.useCase = pegout.NewClaimPegOutUseCase(
		blockchain.RskContracts{PegOut: f.pegout, PegOutEscrow: f.escrow},
		blockchain.Rpc{Rsk: f.rskRpc, Btc: f.btcRpc},
		f.btcWallet,
		f.lp,
		f.repo,
		f.eventBus,
		f.mutex,
	)
	return f
}

func (f *claimFixtures) expectPreGates() {
	f.repo.On("GetRetainedQuote", mock.Anything, claimRequestHash).Return(nil, nil).Once()
	f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()
	f.lp.On("RskAddress").Return(claimLpAddress)
	f.escrow.EXPECT().RestrictedUntil(claimLpAddress).Return(uint64(0), nil).Once()
}

func (f *claimFixtures) expectLoadQuote(q quote.PegoutQuote) {
	f.escrow.EXPECT().GetPegOutQuote(claimRequestHash).Return(q, nil).Once()
	// Only destination (DepositAddress) is non-empty on escrow quotes.
	f.btcRpc.On("EncodeAddress", mock.Anything).Return("bcrt1qdest", nil).Once()
}

func (f *claimFixtures) expectSign(q quote.PegoutQuote) {
	hash := [32]byte{9, 8, 7}
	sig := claimSignature65()
	f.pegout.EXPECT().HashPegoutQuoteEIP712(mock.MatchedBy(func(got quote.PegoutQuote) bool {
		return got.Value.Cmp(q.Value) == 0 &&
			got.CallFee.Cmp(q.CallFee) == 0 &&
			got.DepositAddress == "bcrt1qdest" &&
			got.BtcRefundAddress == "" &&
			got.LpBtcAddress == "" &&
			got.LpRskAddress == claimLpAddress
	})).Return(hash, nil).Once()
	f.lp.On("GetSigner").Return(f.signer).Once()
	f.signer.On("SignBytes", hash[:]).Return(append([]byte(nil), sig...), nil).Once()
}

func TestClaimPegOutUseCase_Step6_SkipInsufficientCapacity(t *testing.T) {
	f := newClaimFixtures(t)
	q := claimEscrowQuote()
	f.expectPreGates()
	f.expectLoadQuote(q)
	f.btcWallet.On("GetBalance").Return(entities.NewWei(100), nil).Once()
	f.repo.On("GetRetainedQuoteByState", mock.Anything, quote.PegoutStateClaimed, quote.PegoutStateWaitingForDepositConfirmations).
		Return([]quote.RetainedPegoutQuote{}, nil).Once()

	log.SetLevel(log.DebugLevel)
	t.Cleanup(func() { log.SetLevel(log.InfoLevel) })
	checkLog := test.LogContains(t, pegout.LogClaimPegoutCapacitySkip(claimRequestHash))
	claimed, err := f.useCase.Run(context.Background(), f.candidate)
	require.NoError(t, err)
	assert.False(t, claimed)
	assert.True(t, checkLog())
	f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
	f.escrow.AssertNotCalled(t, "EstimateClaimPegOut", mock.Anything, mock.Anything)
}

func TestClaimPegOutUseCase_Step6_SkipUnprofitableFeeMath(t *testing.T) {
	f := newClaimFixtures(t)
	q := claimEscrowQuote()
	q.CallFee = entities.NewWei(1) // unprofitable even with GasFee miner-cap credit
	f.expectPreGates()
	f.expectLoadQuote(q)
	f.btcWallet.On("GetBalance").Return(entities.NewWei(10_000_000), nil).Once()
	f.repo.On("GetRetainedQuoteByState", mock.Anything, quote.PegoutStateClaimed, quote.PegoutStateWaitingForDepositConfirmations).
		Return([]quote.RetainedPegoutQuote{}, nil).Once()
	f.expectSign(q)
	f.escrow.EXPECT().EstimateClaimPegOut(claimRequestHash, mock.Anything).Return(entities.NewWei(100_000), nil).Once()
	f.btcWallet.On("EstimateTxFees", "bcrt1qdest", q.Value).Return(blockchain.BtcFeeEstimation{Value: entities.NewWei(50_000)}, nil).Once()
	f.rskRpc.EXPECT().GasPrice(mock.Anything).Return(entities.NewWei(100), nil).Once()

	log.SetLevel(log.DebugLevel)
	t.Cleanup(func() { log.SetLevel(log.InfoLevel) })
	checkLog := test.LogContains(t, pegout.LogClaimPegoutProfitabilitySkip(claimRequestHash))
	claimed, err := f.useCase.Run(context.Background(), f.candidate)
	require.NoError(t, err)
	assert.False(t, claimed)
	assert.True(t, checkLog())
	f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
}

func TestClaimPegOutUseCase_Step6_ClaimTakenWhenBothGatesPass(t *testing.T) {
	f := newClaimFixtures(t)
	q := claimEscrowQuote()
	f.expectPreGates()
	f.expectLoadQuote(q)
	f.btcWallet.On("GetBalance").Return(entities.NewWei(10_000_000), nil).Once()
	f.repo.On("GetRetainedQuoteByState", mock.Anything, quote.PegoutStateClaimed, quote.PegoutStateWaitingForDepositConfirmations).
		Return([]quote.RetainedPegoutQuote{}, nil).Once()
	f.expectSign(q)
	f.escrow.EXPECT().EstimateClaimPegOut(claimRequestHash, mock.Anything).Return(entities.NewWei(100_000), nil).Once()
	f.btcWallet.On("EstimateTxFees", "bcrt1qdest", q.Value).Return(blockchain.BtcFeeEstimation{Value: entities.NewWei(1)}, nil).Once()
	f.rskRpc.EXPECT().GasPrice(mock.Anything).Return(entities.NewWei(1), nil).Once()

	f.mutex.On("Lock").Return().Once()
	f.mutex.On("Unlock").Return().Once()
	claimTx := "0xclaimtx"
	f.escrow.EXPECT().ClaimPegOut(mock.Anything, claimRequestHash, mock.Anything).Return(blockchain.TransactionReceipt{
		TransactionHash: claimTx,
	}, nil).Once()
	f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateClaimed, nil).Once()
	f.pegout.EXPECT().GetAddress().Return("0xpegout").Once()
	f.repo.On("InsertQuote", mock.Anything, mock.MatchedBy(func(c quote.CreatedPegoutQuote) bool {
		return c.Hash == claimRequestHash && c.Quote.LpRskAddress == claimLpAddress
	})).Return(nil).Once()
	f.repo.On("InsertRetainedQuote", mock.Anything, mock.MatchedBy(func(r quote.RetainedPegoutQuote) bool {
		return r.QuoteHash == claimRequestHash && r.State == quote.PegoutStateClaimed && r.UserRskTxHash == claimTx
	})).Return(nil).Once()
	f.eventBus.On("Publish", mock.Anything).Return().Once()

	checkLog := test.LogContains(t, pegout.LogClaimPegoutSuccess(claimRequestHash, claimTx))
	claimed, err := f.useCase.Run(context.Background(), f.candidate)
	require.NoError(t, err)
	assert.True(t, claimed)
	assert.True(t, checkLog())
	f.escrow.AssertExpectations(t)
	f.repo.AssertExpectations(t)
}

func TestClaimPegOutUseCase_SkipWhenRestrictedUntilFuture(t *testing.T) {
	f := newClaimFixtures(t)
	f.repo.On("GetRetainedQuote", mock.Anything, claimRequestHash).Return(nil, nil).Once()
	f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateRequested, nil).Once()
	f.lp.On("RskAddress").Return(claimLpAddress).Once()
	restrictedUntil := uint64(1_700_000_100)
	f.escrow.EXPECT().RestrictedUntil(claimLpAddress).Return(restrictedUntil, nil).Once()
	f.rskRpc.EXPECT().GetBlockByNumber(mock.Anything, (*big.Int)(nil)).Return(blockchain.BlockInfo{
		Timestamp: time.Unix(1_700_000_000, 0),
	}, nil).Once()

	log.SetLevel(log.DebugLevel)
	t.Cleanup(func() { log.SetLevel(log.InfoLevel) })
	checkLog := test.LogContains(t, pegout.LogClaimPegoutRestrictedSkip(claimRequestHash, restrictedUntil))
	claimed, err := f.useCase.Run(context.Background(), f.candidate)
	require.NoError(t, err)
	assert.False(t, claimed)
	assert.True(t, checkLog())
	f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
	f.escrow.AssertNotCalled(t, "GetPegOutQuote", mock.Anything)
}

func TestClaimPegOutUseCase_B1_LostRaceAbsorbed(t *testing.T) {
	t.Run("preflight state not REQUESTED", func(t *testing.T) {
		f := newClaimFixtures(t)
		f.repo.On("GetRetainedQuote", mock.Anything, claimRequestHash).Return(nil, nil).Once()
		f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateClaimed, nil).Once()

		checkLog := test.LogContains(t, pegout.LogClaimPegoutLostRace(claimRequestHash))
		claimed, err := f.useCase.Run(context.Background(), f.candidate)
		require.NoError(t, err)
		assert.False(t, claimed)
		assert.True(t, checkLog())
		f.escrow.AssertNotCalled(t, "ClaimPegOut", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("claim tx reverts after peer won", func(t *testing.T) {
		f := newClaimFixtures(t)
		q := claimEscrowQuote()
		f.expectPreGates()
		f.expectLoadQuote(q)
		f.btcWallet.On("GetBalance").Return(entities.NewWei(10_000_000), nil).Once()
		f.repo.On("GetRetainedQuoteByState", mock.Anything, quote.PegoutStateClaimed, quote.PegoutStateWaitingForDepositConfirmations).
			Return([]quote.RetainedPegoutQuote{}, nil).Once()
		f.expectSign(q)
		f.escrow.EXPECT().EstimateClaimPegOut(claimRequestHash, mock.Anything).Return(entities.NewWei(100_000), nil).Once()
		f.btcWallet.On("EstimateTxFees", "bcrt1qdest", q.Value).Return(blockchain.BtcFeeEstimation{Value: entities.NewWei(1)}, nil).Once()
		f.rskRpc.EXPECT().GasPrice(mock.Anything).Return(entities.NewWei(1), nil).Once()
		f.mutex.On("Lock").Return().Once()
		f.mutex.On("Unlock").Return().Once()
		f.escrow.EXPECT().ClaimPegOut(mock.Anything, claimRequestHash, mock.Anything).Return(blockchain.TransactionReceipt{}, assert.AnError).Once()
		f.escrow.EXPECT().GetPegOutState(claimRequestHash).Return(blockchain.EscrowedPegOutStateClaimed, nil).Once()

		checkLog := test.LogContains(t, pegout.LogClaimPegoutLostRace(claimRequestHash))
		claimed, err := f.useCase.Run(context.Background(), f.candidate)
		require.NoError(t, err)
		assert.False(t, claimed)
		assert.True(t, checkLog())
	})
}
