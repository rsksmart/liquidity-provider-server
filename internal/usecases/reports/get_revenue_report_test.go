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
func TestGetRevenueReportUseCase_Run(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(1000),
				GasFee:  entities.NewWei(350000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-1",
				CallForUserGasUsed:    21000,
				CallForUserGasPrice:   entities.NewWei(5),
				RegisterPeginGasUsed:  50000,
				RegisterPeginGasPrice: entities.NewWei(4),
			},
		},
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(1500),
				GasFee:  entities.NewWei(450000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-2",
				CallForUserGasUsed:    21000,
				CallForUserGasPrice:   entities.NewWei(6),
				RegisterPeginGasUsed:  55000,
				RegisterPeginGasPrice: entities.NewWei(5),
			},
		},
	}

	pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
		{
			Quote: quote.PegoutQuote{
				CallFee: entities.NewWei(1200),
				GasFee:  entities.NewWei(280000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash:            "pegout-hash-1",
				RefundPegoutGasUsed:  40000,
				RefundPegoutGasPrice: entities.NewWei(4),
				BridgeRefundGasUsed:  30000,
				BridgeRefundGasPrice: entities.NewWei(3),
				SendPegoutBtcFee:     entities.NewWei(1000),
			},
		},
		{
			Quote: quote.PegoutQuote{
				CallFee: entities.NewWei(1800),
				GasFee:  entities.NewWei(400000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash:            "pegout-hash-2",
				RefundPegoutGasUsed:  45000,
				RefundPegoutGasPrice: entities.NewWei(5),
				BridgeRefundGasUsed:  35000,
				BridgeRefundGasPrice: entities.NewWei(4),
				SendPegoutBtcFee:     entities.NewWei(1500),
			},
		},
	}

	penalizations := []penalization.PenalizedEvent{
		{
			QuoteHash: "pegin-hash-1",
			Penalty:   entities.NewWei(50),
		},
		{
			QuoteHash: "pegout-hash-2",
			Penalty:   entities.NewWei(80),
		},
	}

	// Expected calculations:
	// TotalQuoteCallFees = 1000 + 1500 + 1200 + 1800 = 5500
	// TotalGasFeesCollected = 350000 + 450000 + 280000 + 400000 = 1480000
	// TotalGasSpent (pegin) = (21000*5 + 50000*4) + (21000*6 + 55000*5) = (105000 + 200000) + (126000 + 275000) = 305000 + 401000 = 706000
	// TotalGasSpent (pegout) = (40000*4 + 30000*3 + 1000) + (45000*5 + 35000*4 + 1500) = (160000 + 90000 + 1000) + (225000 + 140000 + 1500) = 251000 + 366500 = 617500
	// TotalGasSpent = 706000 + 617500 = 1323500
	// TotalPenalizations = 50 + 80 = 130
	// GasProfit = 1480000 - 1323500 = 156500
	// TotalProfit = 5500 + 156500 - 130 = 161870

	expectedTotalQuoteCallFees := entities.NewWei(5500)
	expectedTotalGasFeesCollected := entities.NewWei(1480000)
	expectedTotalGasSpent := entities.NewWei(1323500)
	expectedTotalPenalizations := entities.NewWei(130)
	expectedTotalProfit := entities.NewWei(161870)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return(pegoutQuotesWithRetained, nil).Once()

	allQuoteHashes := []string{"pegin-hash-1", "pegin-hash-2", "pegout-hash-1", "pegout-hash-2"}
	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, allQuoteHashes).
		Return(penalizations, nil).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)
	penalizationRepo.AssertExpectations(t)

	assert.Equal(t, expectedTotalQuoteCallFees, result.TotalQuoteCallFees, "TotalQuoteCallFees mismatch")
	assert.Equal(t, expectedTotalGasFeesCollected, result.TotalGasFeesCollected, "TotalGasFeesCollected mismatch")
	assert.Equal(t, expectedTotalGasSpent, result.TotalGasSpent, "TotalGasSpent mismatch")
	assert.Equal(t, expectedTotalPenalizations, result.TotalPenalizations, "TotalPenalizations mismatch")
	assert.Equal(t, expectedTotalProfit, result.TotalProfit, "TotalProfit mismatch")
}

// nolint:funlen
func TestGetRevenueReportUseCase_Run_OnlyPeginQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(500),
				GasFee:  entities.NewWei(300),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-1",
				CallForUserGasUsed:    25000,
				CallForUserGasPrice:   entities.NewWei(5),
				RegisterPeginGasUsed:  55000,
				RegisterPeginGasPrice: entities.NewWei(6),
			},
		},
	}

	// Expected calculations:
	// TotalQuoteCallFees = 500
	// TotalGasFeesCollected = 300
	// TotalGasSpent = (25000*5 + 55000*6) = 125000 + 330000 = 455000
	// TotalPenalizations = 0
	// GasProfit = 300 - 455000 = -454700
	// TotalProfit = 500 + (-454700) - 0 = -454200

	expectedTotalQuoteCallFees := entities.NewWei(500)
	expectedTotalGasFeesCollected := entities.NewWei(300)
	expectedTotalGasSpent := entities.NewWei(455000)
	expectedTotalPenalizations := entities.NewWei(0)
	expectedTotalProfit := entities.NewWei(-454200)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return([]quote.PegoutQuoteWithRetained{}, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegin-hash-1"}).
		Return([]penalization.PenalizedEvent{}, nil).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, expectedTotalQuoteCallFees, result.TotalQuoteCallFees)
	assert.Equal(t, expectedTotalGasFeesCollected, result.TotalGasFeesCollected)
	assert.Equal(t, expectedTotalGasSpent, result.TotalGasSpent)
	assert.Equal(t, expectedTotalPenalizations, result.TotalPenalizations)
	assert.Equal(t, expectedTotalProfit, result.TotalProfit)
}

// nolint:funlen
func TestGetRevenueReportUseCase_Run_OnlyPegoutQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
		{
			Quote: quote.PegoutQuote{
				CallFee: entities.NewWei(400),
				GasFee:  entities.NewWei(250),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash:            "pegout-hash-1",
				RefundPegoutGasUsed:  35000,
				RefundPegoutGasPrice: entities.NewWei(8),
				BridgeRefundGasUsed:  28000,
				BridgeRefundGasPrice: entities.NewWei(6),
				SendPegoutBtcFee:     entities.NewWei(4000),
			},
		},
	}

	// Expected calculations:
	// TotalQuoteCallFees = 400
	// TotalGasFeesCollected = 250
	// TotalGasSpent = (35000*8 + 28000*6 + 4000) = 280000 + 168000 + 4000 = 452000
	// TotalPenalizations = 0
	// GasProfit = 250 - 452000 = -451750
	// TotalProfit = 400 + (-451750) - 0 = -451350

	expectedTotalQuoteCallFees := entities.NewWei(400)
	expectedTotalGasFeesCollected := entities.NewWei(250)
	expectedTotalGasSpent := entities.NewWei(452000)
	expectedTotalPenalizations := entities.NewWei(0)
	expectedTotalProfit := entities.NewWei(-451350)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return([]quote.PeginQuoteWithRetained{}, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return(pegoutQuotesWithRetained, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegout-hash-1"}).
		Return([]penalization.PenalizedEvent{}, nil).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, expectedTotalQuoteCallFees, result.TotalQuoteCallFees)
	assert.Equal(t, expectedTotalGasFeesCollected, result.TotalGasFeesCollected)
	assert.Equal(t, expectedTotalGasSpent, result.TotalGasSpent)
	assert.Equal(t, expectedTotalPenalizations, result.TotalPenalizations)
	assert.Equal(t, expectedTotalProfit, result.TotalProfit)
}

// nolint:funlen
func TestGetRevenueReportUseCase_Run_PositiveProfit(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	// Scenario where gas fees collected > gas spent
	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(1000000),
				GasFee:  entities.NewWei(500000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-1",
				CallForUserGasUsed:    21000,
				CallForUserGasPrice:   entities.NewWei(5),
				RegisterPeginGasUsed:  50000,
				RegisterPeginGasPrice: entities.NewWei(4),
			},
		},
	}

	// Expected calculations:
	// TotalQuoteCallFees = 1000000
	// TotalGasFeesCollected = 500000
	// TotalGasSpent = (21000*5 + 50000*4) = 105000 + 200000 = 305000
	// TotalPenalizations = 0
	// GasProfit = 500000 - 305000 = 195000
	// TotalProfit = 1000000 + 195000 - 0 = 1195000

	expectedTotalQuoteCallFees := entities.NewWei(1000000)
	expectedTotalGasFeesCollected := entities.NewWei(500000)
	expectedTotalGasSpent := entities.NewWei(305000)
	expectedTotalPenalizations := entities.NewWei(0)
	expectedTotalProfit := entities.NewWei(1195000)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return([]quote.PegoutQuoteWithRetained{}, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegin-hash-1"}).
		Return([]penalization.PenalizedEvent{}, nil).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, expectedTotalQuoteCallFees, result.TotalQuoteCallFees)
	assert.Equal(t, expectedTotalGasFeesCollected, result.TotalGasFeesCollected)
	assert.Equal(t, expectedTotalGasSpent, result.TotalGasSpent)
	assert.Equal(t, expectedTotalPenalizations, result.TotalPenalizations)
	assert.Equal(t, expectedTotalProfit, result.TotalProfit)
}

// nolint:funlen
func TestGetRevenueReportUseCase_Run_NegativeProfit(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	// Scenario where gas spent exceeds gas collected, resulting in negative profit
	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(500),
				GasFee:  entities.NewWei(100000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-1",
				CallForUserGasUsed:    21000,
				CallForUserGasPrice:   entities.NewWei(10),
				RegisterPeginGasUsed:  50000,
				RegisterPeginGasPrice: entities.NewWei(8),
			},
		},
	}

	pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
		{
			Quote: quote.PegoutQuote{
				CallFee: entities.NewWei(600),
				GasFee:  entities.NewWei(150000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash:            "pegout-hash-1",
				RefundPegoutGasUsed:  40000,
				RefundPegoutGasPrice: entities.NewWei(9),
				BridgeRefundGasUsed:  30000,
				BridgeRefundGasPrice: entities.NewWei(7),
				SendPegoutBtcFee:     entities.NewWei(5000),
			},
		},
	}

	// Expected calculations:
	// TotalQuoteCallFees = 500 + 600 = 1100
	// TotalGasFeesCollected = 100000 + 150000 = 250000
	// TotalGasSpent (pegin) = (21000*10 + 50000*8) = 210000 + 400000 = 610000
	// TotalGasSpent (pegout) = (40000*9 + 30000*7 + 5000) = 360000 + 210000 + 5000 = 575000
	// TotalGasSpent = 610000 + 575000 = 1185000
	// TotalPenalizations = 0
	// GasProfit = 250000 - 1185000 = -935000
	// TotalProfit = 1100 + (-935000) - 0 = -933900

	expectedTotalQuoteCallFees := entities.NewWei(1100)
	expectedTotalGasFeesCollected := entities.NewWei(250000)
	expectedTotalGasSpent := entities.NewWei(1185000)
	expectedTotalPenalizations := entities.NewWei(0)
	expectedTotalProfit := entities.NewWei(-933900)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return(pegoutQuotesWithRetained, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegin-hash-1", "pegout-hash-1"}).
		Return([]penalization.PenalizedEvent{}, nil).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, expectedTotalQuoteCallFees, result.TotalQuoteCallFees)
	assert.Equal(t, expectedTotalGasFeesCollected, result.TotalGasFeesCollected)
	assert.Equal(t, expectedTotalGasSpent, result.TotalGasSpent)
	assert.Equal(t, expectedTotalPenalizations, result.TotalPenalizations)
	assert.Equal(t, expectedTotalProfit, result.TotalProfit)
	assert.Negative(t, result.TotalProfit.AsBigInt().Sign(), "Expected negative profit")
}

// nolint:funlen
func TestGetRevenueReportUseCase_Run_WithHighPenalizations(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(100000),
				GasFee:  entities.NewWei(50000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-1",
				CallForUserGasUsed:    21000,
				CallForUserGasPrice:   entities.NewWei(3),
				RegisterPeginGasUsed:  50000,
				RegisterPeginGasPrice: entities.NewWei(3),
			},
		},
	}

	penalizations := []penalization.PenalizedEvent{
		{
			QuoteHash: "pegin-hash-1",
			Penalty:   entities.NewWei(80000),
		},
	}

	// Expected calculations:
	// TotalQuoteCallFees = 100000
	// TotalGasFeesCollected = 50000
	// TotalGasSpent = (21000*3 + 50000*3) = 63000 + 150000 = 213000
	// TotalPenalizations = 80000
	// GasProfit = 50000 - 213000 = -163000
	// TotalProfit = 100000 + (-163000) - 80000 = -143000

	expectedTotalQuoteCallFees := entities.NewWei(100000)
	expectedTotalGasFeesCollected := entities.NewWei(50000)
	expectedTotalGasSpent := entities.NewWei(213000)
	expectedTotalPenalizations := entities.NewWei(80000)
	expectedTotalProfit := entities.NewWei(-143000)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return([]quote.PegoutQuoteWithRetained{}, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegin-hash-1"}).
		Return(penalizations, nil).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	assert.Equal(t, expectedTotalQuoteCallFees, result.TotalQuoteCallFees)
	assert.Equal(t, expectedTotalGasFeesCollected, result.TotalGasFeesCollected)
	assert.Equal(t, expectedTotalGasSpent, result.TotalGasSpent)
	assert.Equal(t, expectedTotalPenalizations, result.TotalPenalizations)
	assert.Equal(t, expectedTotalProfit, result.TotalProfit)
}

func TestGetRevenueReportUseCase_Run_EmptyQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return([]quote.PeginQuoteWithRetained{}, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return([]quote.PegoutQuoteWithRetained{}, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{}).
		Return([]penalization.PenalizedEvent{}, nil).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)
	penalizationRepo.AssertExpectations(t)

	assert.Equal(t, entities.NewWei(0), result.TotalQuoteCallFees)
	assert.Equal(t, entities.NewWei(0), result.TotalGasFeesCollected)
	assert.Equal(t, entities.NewWei(0), result.TotalGasSpent)
	assert.Equal(t, entities.NewWei(0), result.TotalPenalizations)
	assert.Equal(t, entities.NewWei(0), result.TotalProfit)
}

func TestGetRevenueReportUseCase_Run_ErrorFetchingPeginQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return(nil, assert.AnError).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.Error(t, err)
	peginQuoteRepo.AssertExpectations(t)
	assert.Zero(t, result.TotalQuoteCallFees)
	assert.Zero(t, result.TotalGasFeesCollected)
	assert.Zero(t, result.TotalGasSpent)
	assert.Zero(t, result.TotalPenalizations)
	assert.Zero(t, result.TotalProfit)
}

func TestGetRevenueReportUseCase_Run_ErrorFetchingPegoutQuotes(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return([]quote.PeginQuoteWithRetained{}, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return(nil, assert.AnError).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.Error(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)
	assert.Zero(t, result.TotalQuoteCallFees)
	assert.Zero(t, result.TotalGasFeesCollected)
	assert.Zero(t, result.TotalGasSpent)
	assert.Zero(t, result.TotalPenalizations)
	assert.Zero(t, result.TotalProfit)
}

func TestGetRevenueReportUseCase_Run_ErrorFetchingPenalizations(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(100),
				GasFee:  entities.NewWei(50),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-1",
				CallForUserGasUsed:    21000,
				CallForUserGasPrice:   entities.NewWei(10),
				RegisterPeginGasUsed:  50000,
				RegisterPeginGasPrice: entities.NewWei(8),
			},
		},
	}

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return([]quote.PegoutQuoteWithRetained{}, nil).Once()

	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, []string{"pegin-hash-1"}).
		Return(nil, assert.AnError).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.Error(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)
	penalizationRepo.AssertExpectations(t)
	assert.Zero(t, result.TotalQuoteCallFees)
	assert.Zero(t, result.TotalGasFeesCollected)
	assert.Zero(t, result.TotalGasSpent)
	assert.Zero(t, result.TotalPenalizations)
	assert.Zero(t, result.TotalProfit)
}

// nolint:funlen
func TestGetRevenueReportUseCase_Run_MultipleQuotesWithComplexScenario(t *testing.T) {
	ctx := context.Background()
	startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC)

	// Multiple pegin quotes with varying gas prices
	peginQuotesWithRetained := []quote.PeginQuoteWithRetained{
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(2000),
				GasFee:  entities.NewWei(1100000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-1",
				CallForUserGasUsed:    22000,
				CallForUserGasPrice:   entities.NewWei(15),
				RegisterPeginGasUsed:  52000,
				RegisterPeginGasPrice: entities.NewWei(12),
			},
		},
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(2500),
				GasFee:  entities.NewWei(1300000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-2",
				CallForUserGasUsed:    23000,
				CallForUserGasPrice:   entities.NewWei(18),
				RegisterPeginGasUsed:  58000,
				RegisterPeginGasPrice: entities.NewWei(14),
			},
		},
		{
			Quote: quote.PeginQuote{
				CallFee: entities.NewWei(3000),
				GasFee:  entities.NewWei(1600000),
			},
			RetainedQuote: quote.RetainedPeginQuote{
				QuoteHash:             "pegin-hash-3",
				CallForUserGasUsed:    24000,
				CallForUserGasPrice:   entities.NewWei(20),
				RegisterPeginGasUsed:  60000,
				RegisterPeginGasPrice: entities.NewWei(16),
			},
		},
	}

	// Multiple pegout quotes with varying gas prices and BTC fees
	pegoutQuotesWithRetained := []quote.PegoutQuoteWithRetained{
		{
			Quote: quote.PegoutQuote{
				CallFee: entities.NewWei(2200),
				GasFee:  entities.NewWei(800000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash:            "pegout-hash-1",
				RefundPegoutGasUsed:  38000,
				RefundPegoutGasPrice: entities.NewWei(11),
				BridgeRefundGasUsed:  32000,
				BridgeRefundGasPrice: entities.NewWei(9),
				SendPegoutBtcFee:     entities.NewWei(5500),
			},
		},
		{
			Quote: quote.PegoutQuote{
				CallFee: entities.NewWei(2800),
				GasFee:  entities.NewWei(1000000),
			},
			RetainedQuote: quote.RetainedPegoutQuote{
				QuoteHash:            "pegout-hash-2",
				RefundPegoutGasUsed:  42000,
				RefundPegoutGasPrice: entities.NewWei(13),
				BridgeRefundGasUsed:  36000,
				BridgeRefundGasPrice: entities.NewWei(10),
				SendPegoutBtcFee:     entities.NewWei(6500),
			},
		},
	}

	// Multiple penalizations on different quotes
	penalizations := []penalization.PenalizedEvent{
		{
			QuoteHash: "pegin-hash-1",
			Penalty:   entities.NewWei(500),
		},
		{
			QuoteHash: "pegin-hash-3",
			Penalty:   entities.NewWei(700),
		},
		{
			QuoteHash: "pegout-hash-2",
			Penalty:   entities.NewWei(800),
		},
	}

	// Expected calculations:
	// TotalQuoteCallFees = (2000 + 2500 + 3000) + (2200 + 2800) = 7500 + 5000 = 12500
	// TotalGasFeesCollected = (1100000 + 1300000 + 1600000) + (800000 + 1000000) = 4000000 + 1800000 = 5800000
	// TotalGasSpent (pegin) = (22000*15 + 52000*12) + (23000*18 + 58000*14) + (24000*20 + 60000*16)
	//                       = (330000 + 624000) + (414000 + 812000) + (480000 + 960000)
	//                       = 954000 + 1226000 + 1440000 = 3620000
	// TotalGasSpent (pegout) = (38000*11 + 32000*9 + 5500) + (42000*13 + 36000*10 + 6500)
	//                        = (418000 + 288000 + 5500) + (546000 + 360000 + 6500)
	//                        = 711500 + 912500 = 1624000
	// TotalGasSpent = 3620000 + 1624000 = 5244000
	// TotalPenalizations = 500 + 700 + 800 = 2000
	// GasProfit = 5800000 - 5244000 = 556000
	// TotalProfit = 12500 + 556000 - 2000 = 566500

	expectedTotalQuoteCallFees := entities.NewWei(12500)
	expectedTotalGasFeesCollected := entities.NewWei(5800000)
	expectedTotalGasSpent := entities.NewWei(5244000)
	expectedTotalPenalizations := entities.NewWei(2000)
	expectedTotalProfit := entities.NewWei(566500)

	peginQuoteRepo := &mocks.PeginQuoteRepositoryMock{}
	pegoutQuoteRepo := &mocks.PegoutQuoteRepositoryMock{}
	penalizationRepo := &mocks.PenalizedEventRepositoryMock{}

	peginQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}, startDate, endDate).
		Return(peginQuotesWithRetained, nil).Once()

	pegoutQuoteRepo.On("GetQuotesWithRetainedByStateAndDate", ctx, []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}, startDate, endDate).
		Return(pegoutQuotesWithRetained, nil).Once()

	allQuoteHashes := []string{"pegin-hash-1", "pegin-hash-2", "pegin-hash-3", "pegout-hash-1", "pegout-hash-2"}
	penalizationRepo.On("GetPenalizationsByQuoteHashes", ctx, allQuoteHashes).
		Return(penalizations, nil).Once()

	useCase := reports.NewGetRevenueReportUseCase(peginQuoteRepo, pegoutQuoteRepo, penalizationRepo)

	result, err := useCase.Run(ctx, startDate, endDate)

	require.NoError(t, err)
	peginQuoteRepo.AssertExpectations(t)
	pegoutQuoteRepo.AssertExpectations(t)
	penalizationRepo.AssertExpectations(t)

	assert.Equal(t, expectedTotalQuoteCallFees, result.TotalQuoteCallFees, "TotalQuoteCallFees mismatch")
	assert.Equal(t, expectedTotalGasFeesCollected, result.TotalGasFeesCollected, "TotalGasFeesCollected mismatch")
	assert.Equal(t, expectedTotalGasSpent, result.TotalGasSpent, "TotalGasSpent mismatch")
	assert.Equal(t, expectedTotalPenalizations, result.TotalPenalizations, "TotalPenalizations mismatch")
	assert.Equal(t, expectedTotalProfit, result.TotalProfit, "TotalProfit mismatch")
}
