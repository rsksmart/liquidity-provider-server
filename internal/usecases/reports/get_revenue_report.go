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
}

type revenueTotals struct {
	CallFees     *entities.Wei
	GasCollected *entities.Wei
	GasSpent     *entities.Wei
}

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

	peginTotals := useCase.calculatePeginTotals(peginResult)
	pegoutTotals := useCase.calculatePegoutTotals(pegoutResult)
	totalPenalizations := useCase.calculateTotalPenalizations(penalizations)

	totalQuoteCallFees := entities.NewWei(0).Add(peginTotals.CallFees, pegoutTotals.CallFees)
	totalGasFeesCollected := entities.NewWei(0).Add(peginTotals.GasCollected, pegoutTotals.GasCollected)
	totalGasSpent := entities.NewWei(0).Add(peginTotals.GasSpent, pegoutTotals.GasSpent)

	return GetRevenueReportResult{
		TotalQuoteCallFees:    totalQuoteCallFees,
		TotalGasFeesCollected: totalGasFeesCollected,
		TotalGasSpent:         totalGasSpent,
		TotalPenalizations:    totalPenalizations,
	}, nil
}

func (useCase *GetRevenueReportUseCase) calculatePeginTotals(
	peginResult []quote.PeginQuoteWithRetained,
) revenueTotals {
	totalCallFees := entities.NewWei(0)
	totalGasCollected := entities.NewWei(0)
	totalGasSpent := entities.NewWei(0)

	for _, quoteWithRetained := range peginResult {
		totalCallFees = totalCallFees.Add(totalCallFees, quoteWithRetained.Quote.CallFee)
		totalGasCollected = totalGasCollected.Add(totalGasCollected, quoteWithRetained.Quote.GasFee)

		// Calculate actual gas spent
		callForUserGasSpent := entities.NewWei(0).Mul(
			entities.NewUWei(quoteWithRetained.RetainedQuote.CallForUserGasUsed),
			quoteWithRetained.RetainedQuote.CallForUserGasPrice,
		)
		registerPeginGasSpent := entities.NewWei(0).Mul(
			entities.NewUWei(quoteWithRetained.RetainedQuote.RegisterPeginGasUsed),
			quoteWithRetained.RetainedQuote.RegisterPeginGasPrice,
		)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, callForUserGasSpent)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, registerPeginGasSpent)
	}

	return revenueTotals{
		CallFees:     totalCallFees,
		GasCollected: totalGasCollected,
		GasSpent:     totalGasSpent,
	}
}

func (useCase *GetRevenueReportUseCase) calculatePegoutTotals(
	pegoutResult []quote.PegoutQuoteWithRetained,
) revenueTotals {
	totalCallFees := entities.NewWei(0)
	totalGasCollected := entities.NewWei(0)
	totalGasSpent := entities.NewWei(0)

	for _, pair := range pegoutResult {
		totalCallFees = totalCallFees.Add(totalCallFees, pair.Quote.CallFee)
		totalGasCollected = totalGasCollected.Add(totalGasCollected, pair.Quote.GasFee)

		// Calculate actual gas spent (RSK gas + BTC fees)
		refundPegoutGasCost := entities.NewWei(0).Mul(
			entities.NewUWei(pair.RetainedQuote.RefundPegoutGasUsed),
			pair.RetainedQuote.RefundPegoutGasPrice,
		)
		bridgeRefundGasCost := entities.NewWei(0).Mul(
			entities.NewUWei(pair.RetainedQuote.BridgeRefundGasUsed),
			pair.RetainedQuote.BridgeRefundGasPrice,
		)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, refundPegoutGasCost)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, bridgeRefundGasCost)
		totalGasSpent = totalGasSpent.Add(totalGasSpent, pair.RetainedQuote.SendPegoutBtcFee)
	}

	return revenueTotals{
		CallFees:     totalCallFees,
		GasCollected: totalGasCollected,
		GasSpent:     totalGasSpent,
	}
}

func (useCase *GetRevenueReportUseCase) calculateTotalPenalizations(
	penalizations []penalization.PenalizedEvent,
) *entities.Wei {
	total := entities.NewWei(0)
	for _, p := range penalizations {
		total = total.Add(total, p.Penalty)
	}
	return total
}

func (useCase *GetRevenueReportUseCase) getPeginQuotes(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]quote.PeginQuoteWithRetained, error) {
	peginStates := []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}
	quotes, err := useCase.peginQuoteRepository.GetQuotesWithRetainedByStateAndDate(ctx, peginStates, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Filter out quotes without retained data (non-accepted quotes).
	// The repository's aggregation pipeline may include quotes without retained data
	// due to its design for other use cases. Revenue report only processes completed
	// transactions that have valid retained quote data.
	result := make([]quote.PeginQuoteWithRetained, 0, len(quotes))
	for _, q := range quotes {
		if q.RetainedQuote.QuoteHash != "" {
			result = append(result, q)
		}
	}
	return result, nil
}

func (useCase *GetRevenueReportUseCase) getPegoutQuotes(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]quote.PegoutQuoteWithRetained, error) {
	pegoutStates := []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded, quote.PegoutStateBtcReleased}
	quotes, err := useCase.pegoutQuoteRepository.GetQuotesWithRetainedByStateAndDate(ctx, pegoutStates, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Filter out quotes without retained data (non-accepted quotes).
	// The repository's aggregation pipeline may include quotes without retained data
	// due to its design for other use cases. Revenue report only processes completed
	// transactions that have valid retained quote data.
	result := make([]quote.PegoutQuoteWithRetained, 0, len(quotes))
	for _, q := range quotes {
		if q.RetainedQuote.QuoteHash != "" {
			result = append(result, q)
		}
	}
	return result, nil
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
