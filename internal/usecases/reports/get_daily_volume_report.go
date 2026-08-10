package reports

import (
	"context"
	"slices"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

// The repositories treat a zero page and page size as "return the whole range". The report
// aggregates over the full range, so it must not constrain the query to a single page.
const (
	unpaginatedPage    = 0
	unpaginatedPerPage = 0
)

// DailyVolumeItem holds the aggregated pegin and pegout volume for a single day
type DailyVolumeItem struct {
	Day          string        `json:"day"`
	PeginVolume  *entities.Wei `json:"peginVolume"`
	PegoutVolume *entities.Wei `json:"pegoutVolume"`
	PeginCount   int           `json:"peginCount"`
	PegoutCount  int           `json:"pegoutCount"`
}

// GetDailyVolumeReportResult is the day by day volume breakdown for the requested range
type GetDailyVolumeReportResult struct {
	Data              []DailyVolumeItem `json:"data"`
	TotalPeginVolume  *entities.Wei     `json:"totalPeginVolume"`
	TotalPegoutVolume *entities.Wei     `json:"totalPegoutVolume"`
}

// dailyVolume is the aggregation of one side of the report, keyed by day
type dailyVolume struct {
	volumeByDay map[string]*entities.Wei
	countByDay  map[string]int
	total       *entities.Wei
}

func newDailyVolume() dailyVolume {
	return dailyVolume{
		volumeByDay: make(map[string]*entities.Wei),
		countByDay:  make(map[string]int),
		total:       entities.NewWei(0),
	}
}

// add accumulates one quote into the given day. A nil amount counts as zero, so incomplete
// records cannot panic, and the total is built from the same amounts as the per-day volumes
// so that the two always agree.
func (aggregation dailyVolume) add(day string, amounts ...*entities.Wei) {
	if aggregation.volumeByDay[day] == nil {
		aggregation.volumeByDay[day] = entities.NewWei(0)
	}
	for _, amount := range amounts {
		if amount == nil {
			continue
		}
		aggregation.volumeByDay[day].Add(aggregation.volumeByDay[day], amount)
		aggregation.total.Add(aggregation.total, amount)
	}
	aggregation.countByDay[day]++
}

type GetDailyVolumeReportUseCase struct {
	peginRepo  quote.PeginQuoteRepository
	pegoutRepo quote.PegoutQuoteRepository
}

func NewGetDailyVolumeReportUseCase(
	peginRepo quote.PeginQuoteRepository,
	pegoutRepo quote.PegoutQuoteRepository,
) *GetDailyVolumeReportUseCase {
	return &GetDailyVolumeReportUseCase{
		peginRepo:  peginRepo,
		pegoutRepo: pegoutRepo,
	}
}

func (useCase *GetDailyVolumeReportUseCase) Run(ctx context.Context, startTime, endTime time.Time) (GetDailyVolumeReportResult, error) {
	pegin, err := useCase.collectPeginVolume(ctx, startTime, endTime)
	if err != nil {
		return GetDailyVolumeReportResult{}, usecases.WrapUseCaseError(usecases.GetDailyVolumeReportId, err)
	}

	pegout, err := useCase.collectPegoutVolume(ctx, startTime, endTime)
	if err != nil {
		return GetDailyVolumeReportResult{}, usecases.WrapUseCaseError(usecases.GetDailyVolumeReportId, err)
	}

	return GetDailyVolumeReportResult{
		Data:              useCase.buildItems(pegin, pegout),
		TotalPeginVolume:  pegin.total,
		TotalPegoutVolume: pegout.total,
	}, nil
}

// RunForSingleDay returns the daily volume report for one specific day
func (useCase *GetDailyVolumeReportUseCase) RunForSingleDay(ctx context.Context, day time.Time) (GetDailyVolumeReportResult, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return useCase.Run(ctx, start, start.Add(24*time.Hour))
}

func (useCase *GetDailyVolumeReportUseCase) collectPeginVolume(ctx context.Context, startTime, endTime time.Time) (dailyVolume, error) {
	pairs, _, err := useCase.peginRepo.ListQuotesByDateRange(ctx, startTime, endTime, unpaginatedPage, unpaginatedPerPage)
	if err != nil {
		return dailyVolume{}, err
	}

	aggregation := newDailyVolume()
	for _, pair := range pairs {
		day := useCase.dayBucketKey(pair.Quote.AgreementTimestamp)
		aggregation.add(day, pair.Quote.Value, pair.Quote.CallFee)
	}

	return aggregation, nil
}

func (useCase *GetDailyVolumeReportUseCase) collectPegoutVolume(ctx context.Context, startTime, endTime time.Time) (dailyVolume, error) {
	pairs, _, err := useCase.pegoutRepo.ListQuotesByDateRange(ctx, startTime, endTime, unpaginatedPage, unpaginatedPerPage)
	if err != nil {
		return dailyVolume{}, err
	}

	aggregation := newDailyVolume()
	for _, pair := range pairs {
		// The repository filters pegout quotes by agreement timestamp, so bucketing by any
		// other field would place volume outside the range the caller asked for.
		day := useCase.dayBucketKey(pair.Quote.AgreementTimestamp)
		aggregation.add(day, pair.Quote.Value, pair.Quote.CallFee)
	}

	return aggregation, nil
}

func (useCase *GetDailyVolumeReportUseCase) buildItems(pegin, pegout dailyVolume) []DailyVolumeItem {
	days := useCase.sortedDays(pegin, pegout)

	items := make([]DailyVolumeItem, 0, len(days))
	for _, day := range days {
		items = append(items, DailyVolumeItem{
			Day:          day,
			PeginVolume:  useCase.weiOrZero(pegin.volumeByDay[day]),
			PegoutVolume: useCase.weiOrZero(pegout.volumeByDay[day]),
			PeginCount:   pegin.countByDay[day],
			PegoutCount:  pegout.countByDay[day],
		})
	}

	return items
}

// sortedDays returns every day present on either side, in chronological order. The keys use
// time.DateOnly, so ordering them lexicographically orders them by date.
func (useCase *GetDailyVolumeReportUseCase) sortedDays(pegin, pegout dailyVolume) []string {
	days := make([]string, 0, len(pegin.volumeByDay)+len(pegout.volumeByDay))
	for day := range pegin.volumeByDay {
		days = append(days, day)
	}
	for day := range pegout.volumeByDay {
		if _, alreadyAdded := pegin.volumeByDay[day]; !alreadyAdded {
			days = append(days, day)
		}
	}
	slices.Sort(days)

	return days
}

// weiOrZero reports the volume for a day, or zero when that side had no quotes that day
func (useCase *GetDailyVolumeReportUseCase) weiOrZero(volume *entities.Wei) *entities.Wei {
	if volume == nil {
		return entities.NewWei(0)
	}
	return volume
}

func (useCase *GetDailyVolumeReportUseCase) dayBucketKey(timestamp uint32) string {
	return time.Unix(int64(timestamp), 0).UTC().Format(time.DateOnly)
}
