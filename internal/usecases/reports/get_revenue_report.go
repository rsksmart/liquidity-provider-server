package reports

import (
	"context"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

type GetRevenueReportUseCase struct {
	peginQuoteRepository     quote.PeginQuoteRepository
	pegoutQuoteRepository    quote.PegoutQuoteRepository
	penalizedEventRepository penalization.PenalizedEventRepository
}

func NewGetRevenueReportUseCase(
	peginQuoteRepository quote.PeginQuoteRepository,
	pegoutQuoteRepository quote.PegoutQuoteRepository,
	penalizedEventRepository penalization.PenalizedEventRepository,
) *GetRevenueReportUseCase {
	return &GetRevenueReportUseCase{
		peginQuoteRepository:     peginQuoteRepository,
		pegoutQuoteRepository:    pegoutQuoteRepository,
		penalizedEventRepository: penalizedEventRepository,
	}
}

type GetRevenueReportResult struct {
	TotalQuoteCallFees    *entities.Wei
	TotalGasFeesCollected *entities.Wei
	TotalGasSpent         *entities.Wei
	TotalPenalizations    *entities.Wei
	TotalProfit           *entities.Wei
}

// nolint:funlen
func (useCase *GetRevenueReportUseCase) Run(ctx context.Context, startDate time.Time, endDate time.Time) (GetRevenueReportResult, error) {
	peginResult, err := useCase.getPeginQuotes(ctx, startDate, endDate)
	if err != nil {
		return GetRevenueReportResult{}, usecases.WrapUseCaseError(usecases.GetRevenueReportId, err)
	}

	pegoutResult, err := useCase.getPegoutQuotes(ctx, startDate, endDate)
	if err != nil {
		return GetRevenueReportResult{}, usecases.WrapUseCaseError(usecases.GetRevenueReportId, err)
	}

	penalizations, err := useCase.getPenalizations(ctx, peginResult, pegoutResult)
	if err != nil {
		return GetRevenueReportResult{}, usecases.WrapUseCaseError(usecases.GetRevenueReportId, err)
	}

	totalQuoteCallFees := entities.NewWei(0)
	totalGasFeesCollected := entities.NewWei(0)
	totalGasSpent := entities.NewWei(0)
	totalPenalizations := entities.NewWei(0)

	// Calculate totals from pegin quotes
	for _, quoteWithRetained := range peginResult {
		totalQuoteCallFees = totalQuoteCallFees.Add(totalQuoteCallFees, quoteWithRetained.Quote.CallFee)

		totalGasFeesCollected = totalGasFeesCollected.Add(totalGasFeesCollected, quoteWithRetained.Quote.GasFee)

		// Calculate actual gas spent
		callForUserGasSpent := entities.NewWei(0).Mul(
			entities.NewWei(int64(quoteWithRetained.RetainedQuote.CallForUserGasUsed)),
			quoteWithRetained.RetainedQuote.CallForUserGasPrice,
		)
		registerPeginGasSpent := entities.NewWei(0).Mul(
			entities.NewWei(int64(quoteWithRetained.RetainedQuote.RegisterPeginGasUsed)),
			quoteWithRetained.RetainedQuote.RegisterPeginGasPrice,
		)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, callForUserGasSpent)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, registerPeginGasSpent)
	}

	// Calculate totals from pegout quotes
	for _, pair := range pegoutResult {
		totalQuoteCallFees = totalQuoteCallFees.Add(totalQuoteCallFees, pair.Quote.CallFee)

		totalGasFeesCollected = totalGasFeesCollected.Add(totalGasFeesCollected, pair.Quote.GasFee)

		// Calculate actual gas spent (RSK gas + BTC fees)
		refundPegoutGasCost := entities.NewWei(0).Mul(
			entities.NewWei(int64(pair.RetainedQuote.RefundPegoutGasUsed)),
			pair.RetainedQuote.RefundPegoutGasPrice,
		)
		bridgeRefundGasCost := entities.NewWei(0).Mul(
			entities.NewWei(int64(pair.RetainedQuote.BridgeRefundGasUsed)),
			pair.RetainedQuote.BridgeRefundGasPrice,
		)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, refundPegoutGasCost)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, bridgeRefundGasCost)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, pair.RetainedQuote.SendPegoutBtcFee)
	}

	for _, p := range penalizations {
		totalPenalizations = totalPenalizations.Add(totalPenalizations, p.Penalty)
	}

	// Total profit: callFees + (gasCollected - gasSpent) - penalizations
	gasProfit := entities.NewWei(0).Sub(totalGasFeesCollected, totalGasSpent)
	totalProfit := entities.NewWei(0).Add(totalQuoteCallFees, gasProfit)
	totalProfit = totalProfit.Sub(totalProfit, totalPenalizations)

	return GetRevenueReportResult{
		TotalQuoteCallFees:    totalQuoteCallFees,
		TotalGasFeesCollected: totalGasFeesCollected,
		TotalGasSpent:         totalGasSpent,
		TotalPenalizations:    totalPenalizations,
		TotalProfit:           totalProfit,
	}, nil
}

func (useCase *GetRevenueReportUseCase) getPeginQuotes(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]quote.PeginQuoteWithRetained, error) {
	peginStates := []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}
	return useCase.peginQuoteRepository.GetQuotesWithRetainedByStateAndDate(ctx, peginStates, startDate, endDate)
}

func (useCase *GetRevenueReportUseCase) getPegoutQuotes(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]quote.PegoutQuoteWithRetained, error) {
	pegoutStates := []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}
	return useCase.pegoutQuoteRepository.GetQuotesWithRetainedByStateAndDate(ctx, pegoutStates, startDate, endDate)
}

func (useCase *GetRevenueReportUseCase) getPenalizations(
	ctx context.Context,
	peginResult []quote.PeginQuoteWithRetained,
	pegoutResult []quote.PegoutQuoteWithRetained,
) ([]penalization.PenalizedEvent, error) {
	allQuoteHashes := make([]string, 0, len(peginResult)+len(pegoutResult))
	for _, quoteWithRetained := range peginResult {
		allQuoteHashes = append(allQuoteHashes, quoteWithRetained.RetainedQuote.QuoteHash)
	}
	for _, pair := range pegoutResult {
		allQuoteHashes = append(allQuoteHashes, pair.RetainedQuote.QuoteHash)
	}

	return useCase.penalizedEventRepository.GetPenalizationsByQuoteHashes(ctx, allQuoteHashes)
}
