package reports

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"time"
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
	TotalQuoteCallFees *entities.Wei
	TotalPenalizations *entities.Wei
	TotalProfit        *entities.Wei
}

func (useCase *GetRevenueReportUseCase) Run(ctx context.Context, startDate time.Time, endDate time.Time) (GetRevenueReportResult, error) {
	var peginQuotes []quote.PeginQuote
	var pegoutQuotes []quote.PegoutQuote

	peginQuotes, peginQuoteHashes, err := useCase.getPeginQuotes(ctx, startDate, endDate)
	if err != nil {
		return GetRevenueReportResult{}, usecases.WrapUseCaseError(usecases.GetRevenueReportId, err)
	}

	pegoutQuotes, pegoutQuoteHashes, err := useCase.getPegoutQuotes(ctx, startDate, endDate)
	if err != nil {
		return GetRevenueReportResult{}, usecases.WrapUseCaseError(usecases.GetRevenueReportId, err)
	}

	allQuoteHashes := append(peginQuoteHashes, pegoutQuoteHashes...)
	penalizations, err := useCase.penalizedEventRepository.GetPenalizationsByQuoteHashes(ctx, allQuoteHashes)
	if err != nil {
		return GetRevenueReportResult{}, usecases.WrapUseCaseError(usecases.GetRevenueReportId, err)
	}

	totalQuoteCallFees := entities.NewWei(0)
	totalPenalizations := entities.NewWei(0)
	totalProfit := entities.NewWei(0)

	for _, q := range peginQuotes {
		totalQuoteCallFees = totalQuoteCallFees.Add(totalQuoteCallFees, q.CallFee)
	}
	for _, q := range pegoutQuotes {
		totalQuoteCallFees = totalQuoteCallFees.Add(totalQuoteCallFees, q.CallFee)
	}
	for _, p := range penalizations {
		totalPenalizations = totalPenalizations.Add(totalPenalizations, p.Penalty)
	}

	totalProfit = totalProfit.Sub(totalQuoteCallFees, totalPenalizations)

	return GetRevenueReportResult{
		TotalQuoteCallFees: totalQuoteCallFees,
		TotalPenalizations: totalPenalizations,
		TotalProfit:        totalProfit,
	}, nil
}

func (useCase *GetRevenueReportUseCase) getPeginQuotes(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]quote.PeginQuote, []string, error) {
	peginStates := []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}
	peginRetainedQuotes, err := useCase.peginQuoteRepository.GetRetainedQuoteByState(ctx, peginStates...)

	if err != nil {
		return make([]quote.PeginQuote, 0), make([]string, 0), err
	}

	peginQuoteHashes := make([]string, 0, len(peginRetainedQuotes))
	for _, q := range peginRetainedQuotes {
		peginQuoteHashes = append(peginQuoteHashes, q.QuoteHash)
	}

	peginQuotes, err := useCase.peginQuoteRepository.GetQuotesByHashesAndDate(ctx, peginQuoteHashes, startDate, endDate)
	if err != nil {
		return make([]quote.PeginQuote, 0), make([]string, 0), err
	}
	return peginQuotes, peginQuoteHashes, nil
}

func (useCase *GetRevenueReportUseCase) getPegoutQuotes(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) ([]quote.PegoutQuote, []string, error) {
	pegoutStates := []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}
	pegoutRetainedQuotes, err := useCase.pegoutQuoteRepository.GetRetainedQuoteByState(ctx, pegoutStates...)

	if err != nil {
		return make([]quote.PegoutQuote, 0), make([]string, 0), err
	}

	pegoutQuoteHashes := make([]string, 0, len(pegoutRetainedQuotes))
	for _, q := range pegoutRetainedQuotes {
		pegoutQuoteHashes = append(pegoutQuoteHashes, q.QuoteHash)
	}

	pegoutQuotes, err := useCase.pegoutQuoteRepository.GetQuotesByHashesAndDate(ctx, pegoutQuoteHashes, startDate, endDate)

	if err != nil {
		return make([]quote.PegoutQuote, 0), make([]string, 0), err
	}

	return pegoutQuotes, pegoutQuoteHashes, nil
}
