package pegin

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
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

func (useCase *GetPeginReportUseCase) Run(ctx context.Context) (GetPeginReportResult, error) {
	var err error
	var retained []quote.RetainedPeginQuote
	var minimumQuoteValue *entities.Wei
	var maximumQuoteValue *entities.Wei
	var averageQuoteValue *entities.Wei
	var totalFeesCollected *entities.Wei
	var averageFeePerQuote *entities.Wei

	retained, err = useCase.peginQuoteRepository.GetRetainedQuoteByState(ctx, quote.PeginStateRegisterPegInSucceeded)

	if err != nil {
		return GetPeginReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}
	if len(retained) == 0 {
		return GetPeginReportResult{
			NumberOfQuotes:     0,
			MinimumQuoteValue:  entities.NewWei(0),
			MaximumQuoteValue:  entities.NewWei(0),
			AverageQuoteValue:  entities.NewWei(0),
			TotalFeesCollected: entities.NewWei(0),
			AverageFeePerQuote: entities.NewWei(0),
		}, nil
	}

	quotes := make([]quote.PeginQuote, 0)
	for _, q := range retained {
		var peginQuote *quote.PeginQuote
		peginQuote, err = useCase.peginQuoteRepository.GetQuote(ctx, q.QuoteHash)
		if err != nil {
			return GetPeginReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
		}
		if peginQuote != nil {
			quotes = append(quotes, *peginQuote)
		}
	}

	minimumQuoteValue = useCase.calculateMinimumQuoteValue(quotes)
	maximumQuoteValue = useCase.calculateMaximumQuoteValue(quotes)
	averageQuoteValue, err = useCase.calculateAverageQuoteValue(quotes)
	if err != nil {
		return GetPeginReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}
	totalFeesCollected = useCase.calculateTotalFeesCollected(quotes)
	averageFeePerQuote, err = useCase.calculateAverageFeePerQuote(quotes)
	if err != nil {
		return GetPeginReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}

	return GetPeginReportResult{
		NumberOfQuotes:     len(quotes),
		MinimumQuoteValue:  minimumQuoteValue,
		MaximumQuoteValue:  maximumQuoteValue,
		AverageQuoteValue:  averageQuoteValue,
		TotalFeesCollected: totalFeesCollected,
		AverageFeePerQuote: averageFeePerQuote,
	}, nil
}

func (useCase *GetPeginReportUseCase) calculateMinimumQuoteValue(quotes []quote.PeginQuote) *entities.Wei {
	if len(quotes) == 0 {
		return entities.NewWei(0)
	}
	minimum := quotes[0].Value

	if len(quotes) > 1 {
		for _, q := range quotes[1:] {
			if q.Value.Cmp(minimum) < 0 {
				minimum = q.Value
			}
		}
	}

	return minimum
}

func (useCase *GetPeginReportUseCase) calculateMaximumQuoteValue(quotes []quote.PeginQuote) *entities.Wei {
	if len(quotes) == 0 {
		return entities.NewWei(0)
	}
	maximum := quotes[0].Value

	if len(quotes) > 1 {
		for _, q := range quotes[1:] {
			if q.Value.Cmp(maximum) > 0 {
				maximum = q.Value
			}
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
	if len(quotes) == 0 {
		return entities.NewWei(0)
	}

	totalFees := entities.NewWei(0)

	for _, q := range quotes {
		totalFees = totalFees.Add(totalFees, q.CallFee)
	}

	return totalFees
}

func (useCase *GetPeginReportUseCase) calculateAverageFeePerQuote(quotes []quote.PeginQuote) (*entities.Wei, error) {
	if len(quotes) == 0 {
		return entities.NewWei(0), nil
	}

	totalFees := useCase.calculateTotalFeesCollected(quotes)

	averageFee, err := totalFees.Div(totalFees, entities.NewWei(int64(len(quotes))))
	if err != nil {
		return entities.NewWei(0), err
	}

	return averageFee, nil
}
