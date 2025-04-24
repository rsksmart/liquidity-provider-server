package pegin

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"time"
)

type GetPeginReportUseCase struct {
	peginQuoteRepository quote.PeginQuoteRepository
}

func NewGetPeginReportUseCase(
	peginQuoteRepository quote.PeginQuoteRepository,
) *GetPeginReportUseCase {
	return &GetPeginReportUseCase{
		peginQuoteRepository: peginQuoteRepository,
	}
}

type GetPeginReportResult struct {
	NumberOfQuotes     int
	MinimumQuoteValue  *entities.Wei
	MaximumQuoteValue  *entities.Wei
	AverageQuoteValue  *entities.Wei
	TotalFeesCollected *entities.Wei
	AverageFeePerQuote *entities.Wei
}

func (useCase *GetPeginReportUseCase) Run(ctx context.Context, startDate time.Time, endDate time.Time) (GetPeginReportResult, error) {
	var err error
	var quotes []quote.PeginQuote
	var response = GetPeginReportResult{
		NumberOfQuotes:     0,
		MinimumQuoteValue:  entities.NewWei(0),
		MaximumQuoteValue:  entities.NewWei(0),
		AverageQuoteValue:  entities.NewWei(0),
		TotalFeesCollected: entities.NewWei(0),
		AverageFeePerQuote: entities.NewWei(0),
	}

	states := []quote.PeginState{quote.PeginStateRegisterPegInSucceeded}
	retainedQuotes, err := useCase.peginQuoteRepository.GetRetainedQuoteByState(ctx, states...)

	if err != nil {
		return GetPeginReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}

	quoteHashes := make([]string, 0, len(retainedQuotes))
	for _, q := range retainedQuotes {
		quoteHashes = append(quoteHashes, q.QuoteHash)
	}

	if len(quoteHashes) == 0 {
		return response, nil
	}

	quotes, err = useCase.peginQuoteRepository.GetQuotesByHashesAndDate(ctx, quoteHashes, startDate, endDate)

	if err != nil {
		return GetPeginReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}
	if len(quotes) == 0 {
		return response, nil
	}

	response.NumberOfQuotes = len(quotes)
	response.MinimumQuoteValue = useCase.calculateMinimumQuoteValue(quotes)
	response.MaximumQuoteValue = useCase.calculateMaximumQuoteValue(quotes)
	response.AverageQuoteValue, err = useCase.calculateAverageQuoteValue(quotes)
	if err != nil {
		return GetPeginReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}
	response.TotalFeesCollected = useCase.calculateTotalFeesCollected(quotes)
	response.AverageFeePerQuote, err = useCase.calculateAverageFeePerQuote(quotes)
	if err != nil {
		return GetPeginReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}

	return response, nil
}

func (useCase *GetPeginReportUseCase) calculateMinimumQuoteValue(quotes []quote.PeginQuote) *entities.Wei {
	if len(quotes) == 0 {
		return entities.NewWei(0)
	}
	minimum := quotes[0].Value

	for _, q := range quotes {
		if q.Value.Cmp(minimum) < 0 {
			minimum = q.Value
		}
	}

	return minimum
}

func (useCase *GetPeginReportUseCase) calculateMaximumQuoteValue(quotes []quote.PeginQuote) *entities.Wei {
	if len(quotes) == 0 {
		return entities.NewWei(0)
	}
	maximum := quotes[0].Value

	for _, q := range quotes {
		if q.Value.Cmp(maximum) > 0 {
			maximum = q.Value
		}
	}

	return maximum
}

func (useCase *GetPeginReportUseCase) calculateAverageQuoteValue(quotes []quote.PeginQuote) (*entities.Wei, error) {
	if len(quotes) == 0 {
		return entities.NewWei(0), nil
	}

	total := entities.NewWei(0)

	for _, q := range quotes {
		total = total.Add(total, q.Value)
	}

	average, err := total.Div(total, entities.NewWei(int64(len(quotes))))
	if err != nil {
		return entities.NewWei(0), err
	}

	return average, nil
}

func (useCase *GetPeginReportUseCase) calculateTotalFeesCollected(quotes []quote.PeginQuote) *entities.Wei {
	totalFees := entities.NewWei(0)

	for _, q := range quotes {
		totalFees = totalFees.Add(totalFees, q.CallFee)
	}

	return totalFees
}

func (useCase *GetPeginReportUseCase) calculateAverageFeePerQuote(quotes []quote.PeginQuote) (*entities.Wei, error) {
	totalFees := useCase.calculateTotalFeesCollected(quotes)

	averageFee, err := totalFees.Div(totalFees, entities.NewWei(int64(len(quotes))))
	if err != nil {
		return entities.NewWei(0), err
	}

	return averageFee, nil
}
