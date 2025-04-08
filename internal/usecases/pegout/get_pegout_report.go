package pegout

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"time"
)

type GetPegoutReportUseCase struct {
	pegoutQuoteRepository quote.PegoutQuoteRepository
}

func NewGetPegoutReportUseCase(
	pegoutQuoteRepository quote.PegoutQuoteRepository,
) *GetPegoutReportUseCase {
	return &GetPegoutReportUseCase{
		pegoutQuoteRepository: pegoutQuoteRepository,
	}
}

type GetPegoutReportResult struct {
	NumberOfQuotes     int
	MinimumQuoteValue  *entities.Wei
	MaximumQuoteValue  *entities.Wei
	AverageQuoteValue  *entities.Wei
	TotalFeesCollected *entities.Wei
	AverageFeePerQuote *entities.Wei
}

func (useCase *GetPegoutReportUseCase) Run(
	ctx context.Context,
	startDate time.Time,
	endDate time.Time,
) (GetPegoutReportResult, error) {
	var err error
	var quotes []quote.PegoutQuote
	var minimumQuoteValue *entities.Wei
	var maximumQuoteValue *entities.Wei
	var averageQuoteValue *entities.Wei
	var totalFeesCollected *entities.Wei
	var averageFeePerQuote *entities.Wei

	filter := quote.GetPegoutQuotesByStateFilter{
		States:    []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded},
		StartDate: uint32(startDate.Unix()),
		EndDate:   uint32(endDate.Unix()),
	}
	quotes, err = useCase.pegoutQuoteRepository.GetQuotesByState(ctx, filter)
	if err != nil {
		return GetPegoutReportResult{}, usecases.WrapUseCaseError(usecases.GetPegoutReportId, err)
	}
	if len(quotes) == 0 {
		return GetPegoutReportResult{
			NumberOfQuotes:     0,
			MinimumQuoteValue:  entities.NewWei(0),
			MaximumQuoteValue:  entities.NewWei(0),
			AverageQuoteValue:  entities.NewWei(0),
			TotalFeesCollected: entities.NewWei(0),
			AverageFeePerQuote: entities.NewWei(0),
		}, nil
	}

	minimumQuoteValue = useCase.calculateMinimumQuoteValue(quotes)
	maximumQuoteValue = useCase.calculateMaximumQuoteValue(quotes)
	averageQuoteValue, err = useCase.calculateAverageQuoteValue(quotes)
	if err != nil {
		return GetPegoutReportResult{}, usecases.WrapUseCaseError(usecases.GetPegoutReportId, err)
	}
	totalFeesCollected = useCase.calculateTotalFeesCollected(quotes)
	averageFeePerQuote, err = useCase.calculateAverageFeePerQuote(quotes)
	if err != nil {
		return GetPegoutReportResult{}, usecases.WrapUseCaseError(usecases.GetPegoutReportId, err)
	}

	return GetPegoutReportResult{
		NumberOfQuotes:     len(quotes),
		MinimumQuoteValue:  minimumQuoteValue,
		MaximumQuoteValue:  maximumQuoteValue,
		AverageQuoteValue:  averageQuoteValue,
		TotalFeesCollected: totalFeesCollected,
		AverageFeePerQuote: averageFeePerQuote,
	}, nil
}

func (useCase *GetPegoutReportUseCase) calculateMinimumQuoteValue(quotes []quote.PegoutQuote) *entities.Wei {
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

func (useCase *GetPegoutReportUseCase) calculateMaximumQuoteValue(quotes []quote.PegoutQuote) *entities.Wei {
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

func (useCase *GetPegoutReportUseCase) calculateAverageQuoteValue(quotes []quote.PegoutQuote) (*entities.Wei, error) {
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

func (useCase *GetPegoutReportUseCase) calculateTotalFeesCollected(quotes []quote.PegoutQuote) *entities.Wei {
	totalFees := entities.NewWei(0)

	for _, q := range quotes {
		totalFees = totalFees.Add(totalFees, q.CallFee)
	}

	return totalFees
}

func (useCase *GetPegoutReportUseCase) calculateAverageFeePerQuote(quotes []quote.PegoutQuote) (*entities.Wei, error) {
	totalFees := useCase.calculateTotalFeesCollected(quotes)

	averageFee, err := totalFees.Div(totalFees, entities.NewWei(int64(len(quotes))))
	if err != nil {
		return entities.NewWei(0), err
	}

	return averageFee, nil
}
