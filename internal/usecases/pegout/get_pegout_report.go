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

func (useCase *GetPegoutReportUseCase) Run(ctx context.Context, startDate time.Time, endDate time.Time) (GetPegoutReportResult, error) {
	var err error
	var quotes []quote.PegoutQuote
	var minimumQuoteValue *entities.Wei
	var maximumQuoteValue *entities.Wei
	var averageQuoteValue *entities.Wei
	var totalFeesCollected *entities.Wei
	var averageFeePerQuote *entities.Wei

	states := []quote.PegoutState{quote.PegoutStateRefundPegOutSucceeded}
	retainedQuotes, err := useCase.pegoutQuoteRepository.GetRetainedQuoteByState(ctx, states...)

	if err != nil {
		return GetPegoutReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}

	quoteHashes := make([]string, 0, len(retainedQuotes))
	for _, q := range retainedQuotes {
		quoteHashes = append(quoteHashes, q.QuoteHash)
	}
	filters := []quote.QueryFilter{
		{
			Field:    "agreement_timestamp",
			Operator: "$gte",
			Value:    startDate.Unix(),
		},
		{
			Field:    "agreement_timestamp",
			Operator: "$lte",
			Value:    endDate.Unix(),
		},
	}

	if len(quoteHashes) == 0 {
		return useCase.buildReturn(0, entities.NewWei(0), entities.NewWei(0), entities.NewWei(0), entities.NewWei(0), entities.NewWei(0)), nil
	}

	quotes, err = useCase.pegoutQuoteRepository.GetQuotes(ctx, filters, quoteHashes)

	if err != nil {
		return GetPegoutReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}
	if len(quotes) == 0 {
		return useCase.buildReturn(0, entities.NewWei(0), entities.NewWei(0), entities.NewWei(0), entities.NewWei(0), entities.NewWei(0)), nil
	}

	minimumQuoteValue = useCase.calculateMinimumQuoteValue(quotes)
	maximumQuoteValue = useCase.calculateMaximumQuoteValue(quotes)
	averageQuoteValue, err = useCase.calculateAverageQuoteValue(quotes)
	if err != nil {
		return GetPegoutReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}
	totalFeesCollected = useCase.calculateTotalFeesCollected(quotes)
	averageFeePerQuote, err = useCase.calculateAverageFeePerQuote(quotes)
	if err != nil {
		return GetPegoutReportResult{}, usecases.WrapUseCaseError(usecases.GetPeginReportId, err)
	}

	return useCase.buildReturn(len(quotes), minimumQuoteValue, maximumQuoteValue, averageQuoteValue, totalFeesCollected, averageFeePerQuote), nil
}

func (useCase *GetPegoutReportUseCase) buildReturn(
	numberOfQuotes int,
	minimum, maximum, averageQuote, totalFees, averageFess *entities.Wei,
) GetPegoutReportResult {
	return GetPegoutReportResult{
		NumberOfQuotes:     numberOfQuotes,
		MinimumQuoteValue:  minimum,
		MaximumQuoteValue:  maximum,
		AverageQuoteValue:  averageQuote,
		TotalFeesCollected: totalFees,
		AverageFeePerQuote: averageFess,
	}
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
