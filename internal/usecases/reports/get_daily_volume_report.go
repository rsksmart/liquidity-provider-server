package reports

import (
	"context"
	"errors"
	"math/big"
	"sort"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

// DailyVolumeItem holds the aggregated pegin and pegout volume for a single day
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

func dayBucketKey(timestamp uint32) string {
	return time.Unix(int64(timestamp), 0).UTC().Format(time.DateOnly)
}

// quoteVolume is the amount a single quote contributes to both its day bucket and
// the report total. CallFee is optional, so it is only added when present.
func quoteVolume(value, callFee *entities.Wei) *big.Int {
	volume := new(big.Int).Set(value.AsBigInt())
	if callFee != nil {
		volume.Add(volume, callFee.AsBigInt())
	}
	return volume
}

//nolint:cyclop // aggregation needs all the branches in one place
func (useCase *GetDailyVolumeReportUseCase) Run(ctx context.Context, startTime, endTime time.Time, page, perPage int, includeEmptyDays bool) (GetDailyVolumeReportResult, error) {
	if perPage <= 0 {
		return GetDailyVolumeReportResult{}, errors.New("daily volume report requires a positive page size")
	}

	log.Debugf("building daily volume report from %s to %s", startTime.Format(time.DateOnly), endTime.Format(time.DateOnly))

	peginVolume, peginCounts, peginTotal, err := useCase.collectPeginVolume(ctx, startTime, endTime, page, perPage)
	if err != nil {
		if err.Error() == "not found" {
			return GetDailyVolumeReportResult{}, errors.New("no pegin quotes in the requested range")
		}
		return GetDailyVolumeReportResult{}, usecases.WrapUseCaseError(usecases.GetDailyVolumeReportId, err)
	}

	pegoutVolume, pegoutCounts, pegoutTotal, err := useCase.collectPegoutVolume(ctx, startTime, endTime, page, perPage)
	if err != nil {
		return GetDailyVolumeReportResult{}, usecases.WrapUseCaseError(usecases.GetDailyVolumeReportId, err)
	}

	days := make([]string, 0, len(peginVolume)+len(pegoutVolume))
	for day := range peginVolume {
		days = append(days, day)
	}
	for day := range pegoutVolume {
		if _, alreadyAdded := peginVolume[day]; alreadyAdded {
			continue
		}
		days = append(days, day)
	}
	sort.Strings(days)

	items := make([]DailyVolumeItem, 0, len(days))
	for _, day := range days {
		peginForDay := peginVolume[day]
		pegoutForDay := pegoutVolume[day]
		if !includeEmptyDays && peginForDay == nil && pegoutForDay == nil {
			continue
		}
		items = append(items, DailyVolumeItem{
			Day:          day,
			PeginVolume:  peginForDay,
			PegoutVolume: pegoutForDay,
			PeginCount:   peginCounts[day],
			PegoutCount:  pegoutCounts[day],
		})
	}

	log.Debugf("daily volume report completed with %d days", len(items))

	return GetDailyVolumeReportResult{
		Data:              items,
		TotalPeginVolume:  peginTotal,
		TotalPegoutVolume: pegoutTotal,
	}, nil
}

// RunForSingleDay returns the daily volume report for one specific day
func (useCase *GetDailyVolumeReportUseCase) RunForSingleDay(ctx context.Context, day time.Time) (GetDailyVolumeReportResult, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return useCase.Run(ctx, start, start.Add(24*time.Hour), 1, 100, false)
}

func (useCase *GetDailyVolumeReportUseCase) collectPeginVolume(ctx context.Context, startTime, endTime time.Time, page, perPage int) (map[string]*big.Int, map[string]int, *big.Int, error) {
	pairs, _, err := useCase.peginRepo.ListQuotesByDateRange(ctx, startTime, endTime, page, perPage)
	if err != nil {
		return nil, nil, nil, err
	}

	volumeByDay := make(map[string]*big.Int)
	countByDay := make(map[string]int)
	total := new(big.Int)

	for _, pair := range pairs {
		if pair.Quote.Value == nil {
			continue
		}
		day := dayBucketKey(pair.Quote.AgreementTimestamp)
		if volumeByDay[day] == nil {
			volumeByDay[day] = new(big.Int)
		}
		volume := quoteVolume(pair.Quote.Value, pair.Quote.CallFee)
		volumeByDay[day].Add(volumeByDay[day], volume)
		countByDay[day]++
		total.Add(total, volume)
	}

	return volumeByDay, countByDay, total, nil
}

func (useCase *GetDailyVolumeReportUseCase) collectPegoutVolume(ctx context.Context, startTime, endTime time.Time, page, perPage int) (map[string]*big.Int, map[string]int, *big.Int, error) {
	pairs, _, err := useCase.pegoutRepo.ListQuotesByDateRange(ctx, startTime, endTime, page, perPage)
	if err != nil {
		return nil, nil, nil, err
	}

	volumeByDay := make(map[string]*big.Int)
	countByDay := make(map[string]int)
	total := new(big.Int)

	for _, pair := range pairs {
		if pair.Quote.Value == nil {
			continue
		}
		day := dayBucketKey(pair.Quote.AgreementTimestamp)
		if volumeByDay[day] == nil {
			volumeByDay[day] = new(big.Int)
		}
		volume := quoteVolume(pair.Quote.Value, pair.Quote.CallFee)
		volumeByDay[day].Add(volumeByDay[day], volume)
		countByDay[day]++
		total.Add(total, volume)
	}

	return volumeByDay, countByDay, total, nil
}
