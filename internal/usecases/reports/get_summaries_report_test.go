package reports_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// nolint:funlen
func TestSummariesUseCase_Run_ComprehensiveScenario(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	// Pegin quotes with mixed states
	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		// Quote without retained (not accepted)
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(100000),
			},
			RetainedQuote: quote.RetainedPeginQuote{}, // Empty retained
		},
		// Accepted quote waiting for deposit
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(200000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "pegin-hash-1",
				State:     quote.PeginStateWaitingForDeposit,
			},
		},
		// Paid quote - CallForUser succeeded
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(300000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "pegin-hash-2",
				State:     quote.PeginStateCallForUserSucceeded,
			},
		},
		// Paid quote - CallForUser failed
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(150000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "pegin-hash-3",
				State:     quote.PeginStateCallForUserFailed,
			},
		},
		// Refunded quote - RegisterPegIn succeeded
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(250000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "pegin-hash-4",
				State:     quote.PeginStateRegisterPegInSucceeded,
			},
		},
	}

	// Pegout quotes with mixed states
	pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
		// Quote without retained (not accepted)
		{
			Quote: quote.PegoutQuote{
				Value: entities.NewWei(80000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{}, // Empty retained
		},
		// Accepted quote waiting for deposit confirmations
		{
			Quote: quote.PegoutQuote{
				Value: entities.NewWei(180000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash: "pegout-hash-1",
				State:     quote.PegoutStateWaitingForDepositConfirmations,
			},
		},
		// Paid quote - SendPegout succeeded
		{
			Quote: quote.PegoutQuote{
				Value: entities.NewWei(220000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash: "pegout-hash-2",
				State:     quote.PegoutStateSendPegoutSucceeded,
			},
		},
		// Refunded quote - BridgeTx succeeded
		{
			Quote: quote.PegoutQuote{
				Value: entities.NewWei(350000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash: "pegout-hash-3",
				State:     quote.PegoutStateBridgeTxSucceeded,
			},
		},
		// Refunded quote - BTC released
		{
			Quote: quote.PegoutQuote{
				Value: entities.NewWei(120000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash: "pegout-hash-4",
				State:     quote.PegoutStateBtcReleased,
			},
		},
	}

	// Expected calculations for Pegin:
	// TotalQuotesCount = 5 (all quotes, including non-accepted)
	// AcceptedQuotesCount = 4 (all except the first one without retained)
	// TotalAcceptedQuotesAmount = 200000 + 300000 + 150000 + 250000 = 900000
	// PaidQuotesCount = 3 (CallForUserSucceeded, CallForUserFailed, RegisterPegInSucceeded)
	// PaidQuotesAmount = 300000 + 150000 + 250000 = 700000
	// RefundedQuotesCount = 1 (RegisterPegInSucceeded)
	// TotalRefundedQuotesAmount = 250000
	// PenalizationsCount = 1
	// TotalPenalizationsAmount = 5000

	// Expected calculations for Pegout:
	// TotalQuotesCount = 5 (all quotes, including non-accepted)
	// AcceptedQuotesCount = 4 (all except the first one without retained)
	// TotalAcceptedQuotesAmount = 180000 + 220000 + 350000 + 120000 = 870000
	// PaidQuotesCount = 3 (SendPegoutSucceeded, BridgeTxSucceeded, BtcReleased)
	// PaidQuotesAmount = 220000 + 350000 + 120000 = 690000
	// RefundedQuotesCount = 2 (BridgeTxSucceeded, BtcReleased)
	// TotalRefundedQuotesAmount = 350000 + 120000 = 470000
	// PenalizationsCount = 1
	// TotalPenalizationsAmount = 8000

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PeginState{
			quote.PeginStateWaitingForDeposit,
			quote.PeginStateWaitingForDepositConfirmations,
			quote.PeginStateTimeForDepositElapsed,
			quote.PeginStateCallForUserSucceeded,
			quote.PeginStateCallForUserFailed,
			quote.PeginStateRegisterPegInSucceeded,
			quote.PeginStateRegisterPegInFailed,
		}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PegoutState{
			quote.PegoutStateWaitingForDeposit,
			quote.PegoutStateWaitingForDepositConfirmations,
			quote.PegoutStateTimeForDepositElapsed,
			quote.PegoutStateSendPegoutSucceeded,
			quote.PegoutStateSendPegoutFailed,
			quote.PegoutStateRefundPegOutSucceeded,
			quote.PegoutStateRefundPegOutFailed,
			quote.PegoutStateBridgeTxSucceeded,
			quote.PegoutStateBridgeTxFailed,
			quote.PegoutStateBtcReleased,
		}, startDate, endDate).
		Return(pegoutQuotesWithRetained, nil).Once()

	// Pegin accepted hashes
	peginAcceptedHashes := []string{"pegin-hash-1", "pegin-hash-2", "pegin-hash-3", "pegin-hash-4"}
	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, peginAcceptedHashes).
		Return([]penalization.PenalizedEvent{{QuoteHash: "pegin-hash-2", Penalty: entities.NewWei(5000)}}, nil).Once()

	// Pegout accepted hashes
	pegoutAcceptedHashes := []string{"pegout-hash-1", "pegout-hash-2", "pegout-hash-3", "pegout-hash-4"}
	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, pegoutAcceptedHashes).
		Return([]penalization.PenalizedEvent{{QuoteHash: "pegout-hash-3", Penalty: entities.NewWei(8000)}}, nil).Once()

	useCase := reports.NewSummariesUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)
	penalizationRepo.AssertExpectations(t)

	// Assert Pegin Summary
	assert.Equal(t, int64(5), result.PeginSummary.TotalQuotesCount, "Pegin TotalQuotesCount mismatch")
	assert.Equal(t, int64(4), result.PeginSummary.AcceptedQuotesCount, "Pegin AcceptedQuotesCount mismatch")
	assert.Equal(t, entities.NewWei(900000), result.PeginSummary.TotalAcceptedQuotesAmount, "Pegin TotalAcceptedQuotesAmount mismatch")
	assert.Equal(t, int64(3), result.PeginSummary.PaidQuotesCount, "Pegin PaidQuotesCount mismatch")
	assert.Equal(t, entities.NewWei(700000), result.PeginSummary.PaidQuotesAmount, "Pegin PaidQuotesAmount mismatch")
	assert.Equal(t, int64(1), result.PeginSummary.RefundedQuotesCount, "Pegin RefundedQuotesCount mismatch")
	assert.Equal(t, entities.NewWei(250000), result.PeginSummary.TotalRefundedQuotesAmount, "Pegin TotalRefundedQuotesAmount mismatch")
	assert.Equal(t, int64(1), result.PeginSummary.PenalizationsCount, "Pegin PenalizationsCount mismatch")
	assert.Equal(t, entities.NewWei(5000), result.PeginSummary.TotalPenalizationsAmount, "Pegin TotalPenalizationsAmount mismatch")

	// Assert Pegout Summary
	assert.Equal(t, int64(5), result.PegoutSummary.TotalQuotesCount, "Pegout TotalQuotesCount mismatch")
	assert.Equal(t, int64(4), result.PegoutSummary.AcceptedQuotesCount, "Pegout AcceptedQuotesCount mismatch")
	assert.Equal(t, entities.NewWei(870000), result.PegoutSummary.TotalAcceptedQuotesAmount, "Pegout TotalAcceptedQuotesAmount mismatch")
	assert.Equal(t, int64(3), result.PegoutSummary.PaidQuotesCount, "Pegout PaidQuotesCount mismatch")
	assert.Equal(t, entities.NewWei(690000), result.PegoutSummary.PaidQuotesAmount, "Pegout PaidQuotesAmount mismatch")
	assert.Equal(t, int64(2), result.PegoutSummary.RefundedQuotesCount, "Pegout RefundedQuotesCount mismatch")
	assert.Equal(t, entities.NewWei(470000), result.PegoutSummary.TotalRefundedQuotesAmount, "Pegout TotalRefundedQuotesAmount mismatch")
	assert.Equal(t, int64(1), result.PegoutSummary.PenalizationsCount, "Pegout PenalizationsCount mismatch")
	assert.Equal(t, entities.NewWei(8000), result.PegoutSummary.TotalPenalizationsAmount, "Pegout TotalPenalizationsAmount mismatch")
}

// nolint:funlen
func TestSummariesUseCase_Run_OnlyNonAcceptedQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	// Only quotes without retained data
	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(100000),
			},
			RetainedQuote: quote.RetainedPeginQuote{}, // No retained data
		},
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(200000),
			},
			RetainedQuote: quote.RetainedPeginQuote{}, // No retained data
		},
	}

	pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
		{
			Quote: quote.PegoutQuote{
				Value: entities.NewWei(150000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{}, // No retained data
		},
	}

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PeginState{
			quote.PeginStateWaitingForDeposit,
			quote.PeginStateWaitingForDepositConfirmations,
			quote.PeginStateTimeForDepositElapsed,
			quote.PeginStateCallForUserSucceeded,
			quote.PeginStateCallForUserFailed,
			quote.PeginStateRegisterPegInSucceeded,
			quote.PeginStateRegisterPegInFailed,
		}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PegoutState{
			quote.PegoutStateWaitingForDeposit,
			quote.PegoutStateWaitingForDepositConfirmations,
			quote.PegoutStateTimeForDepositElapsed,
			quote.PegoutStateSendPegoutSucceeded,
			quote.PegoutStateSendPegoutFailed,
			quote.PegoutStateRefundPegOutSucceeded,
			quote.PegoutStateRefundPegOutFailed,
			quote.PegoutStateBridgeTxSucceeded,
			quote.PegoutStateBridgeTxFailed,
			quote.PegoutStateBtcReleased,
		}, startDate, endDate).
		Return(pegoutQuotesWithRetained, nil).Once()

	// No accepted quotes, so no penalization calls are made (early return in getPenalizationsSummary)

	useCase := reports.NewSummariesUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)

	// Pegin - only total count should be non-zero
	assert.Equal(t, int64(2), result.PeginSummary.TotalQuotesCount)
	assert.Equal(t, int64(0), result.PeginSummary.AcceptedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.TotalAcceptedQuotesAmount)
	assert.Equal(t, int64(0), result.PeginSummary.PaidQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.PaidQuotesAmount)
	assert.Equal(t, int64(0), result.PeginSummary.RefundedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.TotalRefundedQuotesAmount)
	assert.Equal(t, int64(0), result.PeginSummary.PenalizationsCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.TotalPenalizationsAmount)

	// Pegout - only total count should be non-zero
	assert.Equal(t, int64(1), result.PegoutSummary.TotalQuotesCount)
	assert.Equal(t, int64(0), result.PegoutSummary.AcceptedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.TotalAcceptedQuotesAmount)
	assert.Equal(t, int64(0), result.PegoutSummary.PaidQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.PaidQuotesAmount)
	assert.Equal(t, int64(0), result.PegoutSummary.RefundedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.TotalRefundedQuotesAmount)
	assert.Equal(t, int64(0), result.PegoutSummary.PenalizationsCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.TotalPenalizationsAmount)
}

// nolint:funlen
func TestSummariesUseCase_Run_OnlyAcceptedNotPaidQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	// Only accepted quotes in waiting states
	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(100000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "pegin-hash-1",
				State:     quote.PeginStateWaitingForDeposit,
			},
		},
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(200000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "pegin-hash-2",
				State:     quote.PeginStateWaitingForDepositConfirmations,
			},
		},
	}

	pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
		{
			Quote: quote.PegoutQuote{
				Value: entities.NewWei(150000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash: "pegout-hash-1",
				State:     quote.PegoutStateWaitingForDeposit,
			},
		},
	}

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PeginState{
			quote.PeginStateWaitingForDeposit,
			quote.PeginStateWaitingForDepositConfirmations,
			quote.PeginStateTimeForDepositElapsed,
			quote.PeginStateCallForUserSucceeded,
			quote.PeginStateCallForUserFailed,
			quote.PeginStateRegisterPegInSucceeded,
			quote.PeginStateRegisterPegInFailed,
		}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PegoutState{
			quote.PegoutStateWaitingForDeposit,
			quote.PegoutStateWaitingForDepositConfirmations,
			quote.PegoutStateTimeForDepositElapsed,
			quote.PegoutStateSendPegoutSucceeded,
			quote.PegoutStateSendPegoutFailed,
			quote.PegoutStateRefundPegOutSucceeded,
			quote.PegoutStateRefundPegOutFailed,
			quote.PegoutStateBridgeTxSucceeded,
			quote.PegoutStateBridgeTxFailed,
			quote.PegoutStateBtcReleased,
		}, startDate, endDate).
		Return(pegoutQuotesWithRetained, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegin-hash-1", "pegin-hash-2"}).
		Return([]penalization.PenalizedEvent{}, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegout-hash-1"}).
		Return([]penalization.PenalizedEvent{}, nil).Once()

	useCase := reports.NewSummariesUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)
	penalizationRepo.AssertExpectations(t)

	// Pegin - accepted but not paid
	assert.Equal(t, int64(2), result.PeginSummary.TotalQuotesCount)
	assert.Equal(t, int64(2), result.PeginSummary.AcceptedQuotesCount)
	assert.Equal(t, entities.NewWei(300000), result.PeginSummary.TotalAcceptedQuotesAmount)
	assert.Equal(t, int64(0), result.PeginSummary.PaidQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.PaidQuotesAmount)
	assert.Equal(t, int64(0), result.PeginSummary.RefundedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.TotalRefundedQuotesAmount)
	assert.Equal(t, int64(0), result.PeginSummary.PenalizationsCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.TotalPenalizationsAmount)

	// Pegout - accepted but not paid
	assert.Equal(t, int64(1), result.PegoutSummary.TotalQuotesCount)
	assert.Equal(t, int64(1), result.PegoutSummary.AcceptedQuotesCount)
	assert.Equal(t, entities.NewWei(150000), result.PegoutSummary.TotalAcceptedQuotesAmount)
	assert.Equal(t, int64(0), result.PegoutSummary.PaidQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.PaidQuotesAmount)
	assert.Equal(t, int64(0), result.PegoutSummary.RefundedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.TotalRefundedQuotesAmount)
	assert.Equal(t, int64(0), result.PegoutSummary.PenalizationsCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.TotalPenalizationsAmount)
}

// nolint:funlen
func TestSummariesUseCase_Run_WithMultiplePenalizations(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(100000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "pegin-hash-1",
				State:     quote.PeginStateCallForUserSucceeded,
			},
		},
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(200000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "pegin-hash-2",
				State:     quote.PeginStateRegisterPegInSucceeded,
			},
		},
	}

	pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{}

	penalizations := []penalization.PenalizedEvent{
		{
			QuoteHash: "pegin-hash-1",
			Penalty:   entities.NewWei(10000),
		},
		{
			QuoteHash: "pegin-hash-1", // Same quote penalized twice
			Penalty:   entities.NewWei(5000),
		},
		{
			QuoteHash: "pegin-hash-2",
			Penalty:   entities.NewWei(8000),
		},
	}

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PeginState{
			quote.PeginStateWaitingForDeposit,
			quote.PeginStateWaitingForDepositConfirmations,
			quote.PeginStateTimeForDepositElapsed,
			quote.PeginStateCallForUserSucceeded,
			quote.PeginStateCallForUserFailed,
			quote.PeginStateRegisterPegInSucceeded,
			quote.PeginStateRegisterPegInFailed,
		}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PegoutState{
			quote.PegoutStateWaitingForDeposit,
			quote.PegoutStateWaitingForDepositConfirmations,
			quote.PegoutStateTimeForDepositElapsed,
			quote.PegoutStateSendPegoutSucceeded,
			quote.PegoutStateSendPegoutFailed,
			quote.PegoutStateRefundPegOutSucceeded,
			quote.PegoutStateRefundPegOutFailed,
			quote.PegoutStateBridgeTxSucceeded,
			quote.PegoutStateBridgeTxFailed,
			quote.PegoutStateBtcReleased,
		}, startDate, endDate).
		Return(pegoutQuotesWithRetained, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegin-hash-1", "pegin-hash-2"}).
		Return(penalizations, nil).Once()

	// No pegout penalization call because there are no pegout quotes (early return in getPenalizationsSummary)

	useCase := reports.NewSummariesUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)
	penalizationRepo.AssertExpectations(t)

	// PenalizationsCount = 3 (count of penalization events)
	// TotalPenalizationsAmount = 10000 + 5000 + 8000 = 23000
	assert.Equal(t, int64(3), result.PeginSummary.PenalizationsCount)
	assert.Equal(t, entities.NewWei(23000), result.PeginSummary.TotalPenalizationsAmount)
}

// nolint:funlen
func TestSummariesUseCase_Run_EmptyQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PeginState{
			quote.PeginStateWaitingForDeposit,
			quote.PeginStateWaitingForDepositConfirmations,
			quote.PeginStateTimeForDepositElapsed,
			quote.PeginStateCallForUserSucceeded,
			quote.PeginStateCallForUserFailed,
			quote.PeginStateRegisterPegInSucceeded,
			quote.PeginStateRegisterPegInFailed,
		}, startDate, endDate).
		Return([]quote.PeginQuoteWithRetained{}, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PegoutState{
			quote.PegoutStateWaitingForDeposit,
			quote.PegoutStateWaitingForDepositConfirmations,
			quote.PegoutStateTimeForDepositElapsed,
			quote.PegoutStateSendPegoutSucceeded,
			quote.PegoutStateSendPegoutFailed,
			quote.PegoutStateRefundPegOutSucceeded,
			quote.PegoutStateRefundPegOutFailed,
			quote.PegoutStateBridgeTxSucceeded,
			quote.PegoutStateBridgeTxFailed,
			quote.PegoutStateBtcReleased,
		}, startDate, endDate).
		Return([]quote.PegoutQuoteWithRetained{}, nil).Once()

	// When there are no quotes, no penalization calls are made (early return in getPenalizationsSummary)

	useCase := reports.NewSummariesUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)

	// Everything should be zero
	assert.Equal(t, int64(0), result.PeginSummary.TotalQuotesCount)
	assert.Equal(t, int64(0), result.PeginSummary.AcceptedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.TotalAcceptedQuotesAmount)
	assert.Equal(t, int64(0), result.PeginSummary.PaidQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.PaidQuotesAmount)
	assert.Equal(t, int64(0), result.PeginSummary.RefundedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.TotalRefundedQuotesAmount)
	assert.Equal(t, int64(0), result.PeginSummary.PenalizationsCount)
	assert.Equal(t, entities.NewWei(0), result.PeginSummary.TotalPenalizationsAmount)

	assert.Equal(t, int64(0), result.PegoutSummary.TotalQuotesCount)
	assert.Equal(t, int64(0), result.PegoutSummary.AcceptedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.TotalAcceptedQuotesAmount)
	assert.Equal(t, int64(0), result.PegoutSummary.PaidQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.PaidQuotesAmount)
	assert.Equal(t, int64(0), result.PegoutSummary.RefundedQuotesCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.TotalRefundedQuotesAmount)
	assert.Equal(t, int64(0), result.PegoutSummary.PenalizationsCount)
	assert.Equal(t, entities.NewWei(0), result.PegoutSummary.TotalPenalizationsAmount)
}

func TestSummariesUseCase_Run_ErrorFetchingPeginQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PeginState{
			quote.PeginStateWaitingForDeposit,
			quote.PeginStateWaitingForDepositConfirmations,
			quote.PeginStateTimeForDepositElapsed,
			quote.PeginStateCallForUserSucceeded,
			quote.PeginStateCallForUserFailed,
			quote.PeginStateRegisterPegInSucceeded,
			quote.PeginStateRegisterPegInFailed,
		}, startDate, endDate).
		Return(nil, assert.AnError).Once()

	useCase := reports.NewSummariesUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.Error(t, err)
	peginQuoteRepo.AssertExpectations(t)
	assert.Empty(t, result)
}

func TestSummariesUseCase_Run_ErrorFetchingPegoutQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PeginState{
			quote.PeginStateWaitingForDeposit,
			quote.PeginStateWaitingForDepositConfirmations,
			quote.PeginStateTimeForDepositElapsed,
			quote.PeginStateCallForUserSucceeded,
			quote.PeginStateCallForUserFailed,
			quote.PeginStateRegisterPegInSucceeded,
			quote.PeginStateRegisterPegInFailed,
		}, startDate, endDate).
		Return([]quote.PeginQuoteWithRetained{}, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PegoutState{
			quote.PegoutStateWaitingForDeposit,
			quote.PegoutStateWaitingForDepositConfirmations,
			quote.PegoutStateTimeForDepositElapsed,
			quote.PegoutStateSendPegoutSucceeded,
			quote.PegoutStateSendPegoutFailed,
			quote.PegoutStateRefundPegOutSucceeded,
			quote.PegoutStateRefundPegOutFailed,
			quote.PegoutStateBridgeTxSucceeded,
			quote.PegoutStateBridgeTxFailed,
			quote.PegoutStateBtcReleased,
		}, startDate, endDate).
		Return(nil, assert.AnError).Once()

	// No penalization repo call because pegin aggregation fails first

	useCase := reports.NewSummariesUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.Error(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)
	assert.Empty(t, result)
}

// nolint:funlen
func TestSummariesUseCase_Run_ErrorFetchingPenalizations(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				Value: entities.NewWei(100000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash: "pegin-hash-1",
				State:     quote.PeginStateCallForUserSucceeded,
			},
		},
	}

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx,
		[]quote.PeginState{
			quote.PeginStateWaitingForDeposit,
			quote.PeginStateWaitingForDepositConfirmations,
			quote.PeginStateTimeForDepositElapsed,
			quote.PeginStateCallForUserSucceeded,
			quote.PeginStateCallForUserFailed,
			quote.PeginStateRegisterPegInSucceeded,
			quote.PeginStateRegisterPegInFailed,
		}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegin-hash-1"}).
		Return(nil, assert.AnError).Once()

	useCase := reports.NewSummariesUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.Error(t, err)
	peginQuoteRepo.AssertExpectations(t)
	penalizationRepo.AssertExpectations(t)
	assert.Empty(t, result)
}
