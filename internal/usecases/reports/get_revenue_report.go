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

type peginQuotesReturn struct {
	Quotes      []quote.PeginQuote
	QuoteHashes []string
}

type pegoutQuotesReturn struct {
	Quotes      []quote.PegoutQuote
	QuoteHashes []string
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
	peginResult, err := useCase.getPeginQuotes(ctx, startDate, endDate)
	if err != nil {
		return GetRevenueReportResult{}, usecases.WrapUseCaseError(usecases.GetRevenueReportId, err)
	}

	pegoutResult, err := useCase.getPegoutQuotes(ctx, startDate, endDate)
	if err != nil {
		return GetRevenueReportResult{}, usecases.WrapUseCaseError(usecases.GetRevenueReportId, err)
	}

	allQuoteHashes := append(peginResult.QuoteHashes, pegoutResult.QuoteHashes...)
	penalizations, err := useCase.penalizedEventRepository.GetPenalizationsByQuoteHashes(ctx, allQuoteHashes)
	if err != nil {
		return GetRevenueReportResult{}, usecases.WrapUseCaseError(usecases.GetRevenueReportId, err)
	}

	totalQuoteCallFees := entities.NewWei(0)
	totalPenalizations := entities.NewWei(0)
	totalProfit := entities.NewWei(0)

	for _, q := range peginResult.Quotes {
		totalQuoteCallFees = totalQuoteCallFees.Add(totalQuoteCallFees, q.CallFee)
	}
	for _, q := range pegoutResult.Quotes {
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
) (peginQuotesReturn, error) {
	peginStates := []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}
	peginRetainedQuotes, err := useCase.peginQuoteRepository.GetRetainedQuoteByState(ctx, peginStates...)

	peginQuotesToReturn := peginQuotesReturn{
		Quotes:      make([]quote.PeginQuote, 0),
		QuoteHashes: make([]string, 0),
	}

	if err != nil {
		return peginQuotesToReturn, err
	}

	peginQuoteHashes := make([]string, 0, len(peginRetainedQuotes))
	for _, q := range peginRetainedQuotes {
		peginQuoteHashes = append(peginQuoteHashes, q.QuoteHash)
	}

	peginQuotes, err := useCase.peginQuoteRepository.GetQuotesByHashesAndDate(ctx, peginQuoteHashes, startDate, endDate)
	peginQuotesToReturn.Quotes = peginQuotes
	peginQuotesToReturn.QuoteHashes = peginQuoteHashes
	if err != nil {
		return peginQuotesToReturn, err
	}
	return peginQuotesToReturn, nil
}

func (useCase *GetRevenueReportUseCase) getPegoutQuotes(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) (pegoutQuotesReturn, error) {
	pegoutStates := []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded, quote.PegoutStateBridgeTxSucceeded}
	pegoutRetainedQuotes, err := useCase.pegoutQuoteRepository.GetRetainedQuoteByState(ctx, pegoutStates...)

	pegoutQuotesToReturn := pegoutQuotesReturn{
		Quotes:      make([]quote.PegoutQuote, 0),
		QuoteHashes: make([]string, 0),
	}

	if err != nil {
		return pegoutQuotesToReturn, err
	}

	pegoutQuoteHashes := make([]string, 0, len(pegoutRetainedQuotes))
	for _, q := range pegoutRetainedQuotes {
		pegoutQuoteHashes = append(pegoutQuoteHashes, q.QuoteHash)
	}

	pegoutQuotes, err := useCase.pegoutQuoteRepository.GetQuotesByHashesAndDate(ctx, pegoutQuoteHashes, startDate, endDate)
	pegoutQuotesToReturn.Quotes = pegoutQuotes
	pegoutQuotesToReturn.QuoteHashes = pegoutQuoteHashes

	if err != nil {
		return pegoutQuotesToReturn, err
	}

	return pegoutQuotesToReturn, nil
}
