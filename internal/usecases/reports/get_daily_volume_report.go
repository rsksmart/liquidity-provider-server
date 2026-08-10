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
)

// DailyVolumeItem holds the aggregated pegin and pegout volume for a single day
type DailyVolumeItem struct {
	Day          string        `json:"day"`
	PeginVolume  *entities.Wei `json:"peginVolume"`
	PegoutVolume *entities.Wei `json:"pegout_volume"`
	PeginCount   int           `json:"peginCount"`
	PegoutCount  int           `json:"pegoutCount"`
}

// GetDailyVolumeReportResponse is the day by day volume breakdown for the requested range
type GetDailyVolumeReportResponse struct {
	Data              []DailyVolumeItem `json:"data"`
	TotalPeginVolume  *entities.Wei     `json:"total_pegin_volume"`
	TotalPegoutVolume *entities.Wei     `json:"totalPegoutVolume"`
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
	return time.Unix(int64(timestamp), 0).UTC().Format("2006-01-02")
}

// weiOrZero converts an aggregated amount to Wei, treating a missing day as zero volume
func weiOrZero(amount *big.Int) *entities.Wei {
	if amount == nil {
		return entities.NewWei(0)
	}
	return entities.NewBigWei(amount)
}

//nolint:cyclop // aggregation needs all the branches in one place
func (useCase *GetDailyVolumeReportUseCase) Run(ctx context.Context, startTime, endTime time.Time, page, perPage int) (GetDailyVolumeReportResponse, error) {
	if perPage <= 0 {
		return GetDailyVolumeReportResponse{}, errors.New("daily volume report requires a positive page size")
	}

	peginVolume, peginCounts, peginTotal, err := useCase.collectPeginVolume(ctx, startTime, endTime, page, perPage)
	if err != nil {
		if err.Error() == "not found" {
			return GetDailyVolumeReportResponse{}, errors.New("no pegin quotes in the requested range")
		}
		return GetDailyVolumeReportResponse{}, usecases.WrapUseCaseError(usecases.GetDailyVolumeReportId, err)
	}

	pegoutVolume, pegoutCounts, pegoutTotal, err := useCase.collectPegoutVolume(ctx, startTime, endTime, page, perPage)
	if err != nil {
		return GetDailyVolumeReportResponse{}, usecases.WrapUseCaseError(usecases.GetDailyVolumeReportId, err)
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
		items = append(items, DailyVolumeItem{
			Day:          day,
			PeginVolume:  weiOrZero(peginVolume[day]),
			PegoutVolume: weiOrZero(pegoutVolume[day]),
			PeginCount:   peginCounts[day],
			PegoutCount:  pegoutCounts[day],
		})
	}

	return GetDailyVolumeReportResponse{
		Data:              items,
		TotalPeginVolume:  weiOrZero(peginTotal),
		TotalPegoutVolume: weiOrZero(pegoutTotal),
	}, nil
}

// RunForSingleDay returns the daily volume report for one specific day
func (useCase *GetDailyVolumeReportUseCase) RunForSingleDay(ctx context.Context, day time.Time) (GetDailyVolumeReportResponse, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	return useCase.Run(ctx, start, start.Add(24*time.Hour), 1, 100)
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
		volumeByDay[day].Add(volumeByDay[day], pair.Quote.Value.AsBigInt())
		volumeByDay[day].Add(volumeByDay[day], pair.Quote.CallFee.AsBigInt())
		countByDay[day]++
		total.Add(total, pair.Quote.Value.AsBigInt())
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
		day := dayBucketKey(pair.Quote.DepositDateLimit)
		if volumeByDay[day] == nil {
			volumeByDay[day] = new(big.Int)
		}
		volumeByDay[day].Add(volumeByDay[day], pair.Quote.Value.AsBigInt())
		volumeByDay[day].Add(volumeByDay[day], pair.Quote.CallFee.AsBigInt())
		countByDay[day]++
		total.Add(total, pair.Quote.Value.AsBigInt())
	}

	return volumeByDay, countByDay, total, nil
}
