package reports

import (
	"context"
	"math/big"
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

// DailyVolumeItem holds the aggregated pegin and pegout volume for a single day.
// Volumes stay as *big.Int to match an external report consumer (declined Wei conversion).
type DailyVolumeItem struct {
	Day          string   `json:"day"`
	PeginVolume  *big.Int `json:"peginVolume"`
	PegoutVolume *big.Int `json:"pegoutVolume"`
	PeginCount   int      `json:"peginCount"`
	PegoutCount  int      `json:"pegoutCount"`
}

// GetDailyVolumeReportResult is the day by day volume breakdown for the requested range
type GetDailyVolumeReportResult struct {
	Data              []DailyVolumeItem `json:"data"`
	TotalPeginVolume  *big.Int          `json:"totalPeginVolume"`
	TotalPegoutVolume *big.Int          `json:"totalPegoutVolume"`
}

type dailyVolume struct {
	volumeByDay map[string]*big.Int
	countByDay  map[string]int
	total       *big.Int
}

func newDailyVolume() dailyVolume {
	return dailyVolume{
		volumeByDay: make(map[string]*big.Int),
		countByDay:  make(map[string]int),
		total:       new(big.Int),
	}
}

// add accumulates one quote into the given day. A nil amount counts as zero so incomplete
// records cannot panic, and the total is built from the same amounts as the per-day volumes.
func (aggregation dailyVolume) add(day string, amounts ...*entities.Wei) {
	if aggregation.volumeByDay[day] == nil {
		aggregation.volumeByDay[day] = new(big.Int)
	}
	for _, amount := range amounts {
		if amount == nil {
			continue
		}
		aggregation.volumeByDay[day].Add(aggregation.volumeByDay[day], amount.AsBigInt())
		aggregation.total.Add(aggregation.total, amount.AsBigInt())
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

// RunForSingleDay returns the daily volume report for one specific day.
// Kept exported deliberately; it will be wired by the reporting endpoint in a follow-up.
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
		// Match the repository filter, which uses agreement_timestamp.
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
			PeginVolume:  useCase.bigIntOrZero(pegin.volumeByDay[day]),
			PegoutVolume: useCase.bigIntOrZero(pegout.volumeByDay[day]),
			PeginCount:   pegin.countByDay[day],
			PegoutCount:  pegout.countByDay[day],
		})
	}
	return items
}

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

func (useCase *GetDailyVolumeReportUseCase) bigIntOrZero(volume *big.Int) *big.Int {
	if volume == nil {
		return big.NewInt(0)
	}
	return volume
}

func (useCase *GetDailyVolumeReportUseCase) dayBucketKey(timestamp uint32) string {
	return time.Unix(int64(timestamp), 0).UTC().Format(time.DateOnly)
}
