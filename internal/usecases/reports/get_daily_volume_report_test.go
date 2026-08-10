package reports_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	dailyVolumeRangeStart = time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
	dailyVolumeRangeEnd   = time.Date(2023, 1, 17, 0, 0, 0, 0, time.UTC)

	firstDay  = time.Date(2023, 1, 15, 10, 0, 0, 0, time.UTC)
	secondDay = time.Date(2023, 1, 16, 22, 0, 0, 0, time.UTC)
)

func peginQuoteAt(agreementTimestamp time.Time, value, callFee *entities.Wei) quote.PeginQuoteWithRetained {
	return quote.PeginQuoteWithRetained{
		Quote: quote.PeginQuote{
			AgreementTimestamp: uint32(agreementTimestamp.Unix()),
			Value:              value,
			CallFee:            callFee,
		},
	}
}

func pegoutQuoteAt(agreementTimestamp, depositDateLimit time.Time, value, callFee *entities.Wei) quote.PegoutQuoteWithRetained {
	return quote.PegoutQuoteWithRetained{
		Quote: quote.PegoutQuote{
			AgreementTimestamp: uint32(agreementTimestamp.Unix()),
			DepositDateLimit:   uint32(depositDateLimit.Unix()),
			Value:              value,
			CallFee:            callFee,
		},
	}
}

// dailyVolumeUseCase wires the use case to mocks that answer the whole range in one unpaginated
// query, which is the only call shape the report is allowed to make.
func dailyVolumeUseCase(
	t *testing.T,
	peginQuotes []quote.PeginQuoteWithRetained,
	pegoutQuotes []quote.PegoutQuoteWithRetained,
) *reports.GetDailyVolumeReportUseCase {
	t.Helper()

	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)

	peginRepo.On("ListQuotesByDateRange", mock.Anything, dailyVolumeRangeStart, dailyVolumeRangeEnd, 0, 0).
		Return(peginQuotes, len(peginQuotes), nil)
	pegoutRepo.On("ListQuotesByDateRange", mock.Anything, dailyVolumeRangeStart, dailyVolumeRangeEnd, 0, 0).
		Return(pegoutQuotes, len(pegoutQuotes), nil)

	return reports.NewGetDailyVolumeReportUseCase(peginRepo, pegoutRepo)
}

func TestGetDailyVolumeReportUseCase_Run_AggregatesVolumeByDay(t *testing.T) {
	useCase := dailyVolumeUseCase(t,
		[]quote.PeginQuoteWithRetained{
			peginQuoteAt(firstDay, entities.NewWei(100), entities.NewWei(5)),
			peginQuoteAt(firstDay, entities.NewWei(200), entities.NewWei(10)),
			peginQuoteAt(secondDay, entities.NewWei(300), entities.NewWei(15)),
		},
		[]quote.PegoutQuoteWithRetained{
			pegoutQuoteAt(secondDay, secondDay, entities.NewWei(400), entities.NewWei(20)),
		},
	)

	result, err := useCase.Run(context.Background(), dailyVolumeRangeStart, dailyVolumeRangeEnd)

	require.NoError(t, err)
	require.Len(t, result.Data, 2)

	assert.Equal(t, "2023-01-15", result.Data[0].Day)
	assert.Equal(t, "315", result.Data[0].PeginVolume.String())
	assert.Equal(t, 2, result.Data[0].PeginCount)
	assert.Equal(t, "0", result.Data[0].PegoutVolume.String())
	assert.Equal(t, 0, result.Data[0].PegoutCount)

	assert.Equal(t, "2023-01-16", result.Data[1].Day)
	assert.Equal(t, "315", result.Data[1].PeginVolume.String())
	assert.Equal(t, 1, result.Data[1].PeginCount)
	assert.Equal(t, "420", result.Data[1].PegoutVolume.String())
	assert.Equal(t, 1, result.Data[1].PegoutCount)

	assert.Equal(t, "630", result.TotalPeginVolume.String())
	assert.Equal(t, "420", result.TotalPegoutVolume.String())
}

func TestGetDailyVolumeReportUseCase_Run_EmptyRange(t *testing.T) {
	useCase := dailyVolumeUseCase(t, nil, nil)

	result, err := useCase.Run(context.Background(), dailyVolumeRangeStart, dailyVolumeRangeEnd)

	require.NoError(t, err)
	assert.Empty(t, result.Data)
	assert.Equal(t, "0", result.TotalPeginVolume.String())
	assert.Equal(t, "0", result.TotalPegoutVolume.String())
}

func TestGetDailyVolumeReportUseCase_Run_DaysWithOneSideReportZeroForTheOther(t *testing.T) {
	useCase := dailyVolumeUseCase(t,
		[]quote.PeginQuoteWithRetained{
			peginQuoteAt(firstDay, entities.NewWei(100), entities.NewWei(5)),
		},
		[]quote.PegoutQuoteWithRetained{
			pegoutQuoteAt(secondDay, secondDay, entities.NewWei(400), entities.NewWei(20)),
		},
	)

	result, err := useCase.Run(context.Background(), dailyVolumeRangeStart, dailyVolumeRangeEnd)

	require.NoError(t, err)
	require.Len(t, result.Data, 2)

	require.NotNil(t, result.Data[0].PegoutVolume)
	assert.Equal(t, "0", result.Data[0].PegoutVolume.String())

	require.NotNil(t, result.Data[1].PeginVolume)
	assert.Equal(t, "0", result.Data[1].PeginVolume.String())
}

func TestGetDailyVolumeReportUseCase_Run_NilAmountsCountAsZero(t *testing.T) {
	useCase := dailyVolumeUseCase(t,
		[]quote.PeginQuoteWithRetained{
			peginQuoteAt(firstDay, entities.NewWei(100), nil),
			peginQuoteAt(firstDay, nil, nil),
		},
		[]quote.PegoutQuoteWithRetained{
			pegoutQuoteAt(firstDay, firstDay, nil, entities.NewWei(20)),
		},
	)

	result, err := useCase.Run(context.Background(), dailyVolumeRangeStart, dailyVolumeRangeEnd)

	require.NoError(t, err)
	require.Len(t, result.Data, 1)

	assert.Equal(t, "100", result.Data[0].PeginVolume.String())
	assert.Equal(t, 2, result.Data[0].PeginCount)
	assert.Equal(t, "20", result.Data[0].PegoutVolume.String())
	assert.Equal(t, 1, result.Data[0].PegoutCount)
}

func TestGetDailyVolumeReportUseCase_Run_TotalsMatchTheSumOfDailyVolumes(t *testing.T) {
	useCase := dailyVolumeUseCase(t,
		[]quote.PeginQuoteWithRetained{
			peginQuoteAt(firstDay, entities.NewWei(100), entities.NewWei(5)),
			peginQuoteAt(secondDay, entities.NewWei(300), entities.NewWei(15)),
		},
		[]quote.PegoutQuoteWithRetained{
			pegoutQuoteAt(firstDay, firstDay, entities.NewWei(70), entities.NewWei(3)),
			pegoutQuoteAt(secondDay, secondDay, entities.NewWei(400), entities.NewWei(20)),
		},
	)

	result, err := useCase.Run(context.Background(), dailyVolumeRangeStart, dailyVolumeRangeEnd)

	require.NoError(t, err)

	peginSum := entities.NewWei(0)
	pegoutSum := entities.NewWei(0)
	for _, item := range result.Data {
		peginSum.Add(peginSum, item.PeginVolume)
		pegoutSum.Add(pegoutSum, item.PegoutVolume)
	}

	assert.Equal(t, result.TotalPeginVolume.String(), peginSum.String())
	assert.Equal(t, result.TotalPegoutVolume.String(), pegoutSum.String())
}

func TestGetDailyVolumeReportUseCase_Run_BucketsPegoutByAgreementTimestamp(t *testing.T) {
	// The repository filters on agreement timestamp, so a deposit deadline on another day
	// must not move the volume into that day's bucket.
	useCase := dailyVolumeUseCase(t, nil,
		[]quote.PegoutQuoteWithRetained{
			pegoutQuoteAt(firstDay, secondDay, entities.NewWei(400), entities.NewWei(20)),
		},
	)

	result, err := useCase.Run(context.Background(), dailyVolumeRangeStart, dailyVolumeRangeEnd)

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "2023-01-15", result.Data[0].Day)
	assert.Equal(t, "420", result.Data[0].PegoutVolume.String())
}

func TestGetDailyVolumeReportUseCase_Run_RepositoryErrors(t *testing.T) {
	repositoryError := errors.New("database unavailable")

	t.Run("pegin repository fails", func(t *testing.T) {
		peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
		pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
		peginRepo.On("ListQuotesByDateRange", mock.Anything, dailyVolumeRangeStart, dailyVolumeRangeEnd, 0, 0).
			Return(nil, 0, repositoryError)
		useCase := reports.NewGetDailyVolumeReportUseCase(peginRepo, pegoutRepo)

		result, err := useCase.Run(context.Background(), dailyVolumeRangeStart, dailyVolumeRangeEnd)

		require.ErrorIs(t, err, repositoryError)
		assert.Empty(t, result.Data)
	})

	t.Run("pegout repository fails", func(t *testing.T) {
		peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
		pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
		peginRepo.On("ListQuotesByDateRange", mock.Anything, dailyVolumeRangeStart, dailyVolumeRangeEnd, 0, 0).
			Return([]quote.PeginQuoteWithRetained{}, 0, nil)
		pegoutRepo.On("ListQuotesByDateRange", mock.Anything, dailyVolumeRangeStart, dailyVolumeRangeEnd, 0, 0).
			Return(nil, 0, repositoryError)
		useCase := reports.NewGetDailyVolumeReportUseCase(peginRepo, pegoutRepo)

		result, err := useCase.Run(context.Background(), dailyVolumeRangeStart, dailyVolumeRangeEnd)

		require.ErrorIs(t, err, repositoryError)
		assert.Empty(t, result.Data)
	})
}

func TestGetDailyVolumeReportUseCase_RunForSingleDay(t *testing.T) {
	dayStart := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	peginRepo := mocks.NewPeginQuoteRepositoryMock(t)
	pegoutRepo := mocks.NewPegoutQuoteRepositoryMock(t)
	peginRepo.On("ListQuotesByDateRange", mock.Anything, dayStart, dayEnd, 0, 0).
		Return([]quote.PeginQuoteWithRetained{
			peginQuoteAt(firstDay, entities.NewWei(100), entities.NewWei(5)),
		}, 1, nil)
	pegoutRepo.On("ListQuotesByDateRange", mock.Anything, dayStart, dayEnd, 0, 0).
		Return([]quote.PegoutQuoteWithRetained{}, 0, nil)
	useCase := reports.NewGetDailyVolumeReportUseCase(peginRepo, pegoutRepo)

	result, err := useCase.RunForSingleDay(context.Background(), firstDay)

	require.NoError(t, err)
	require.Len(t, result.Data, 1)
	assert.Equal(t, "2023-01-15", result.Data[0].Day)
	assert.Equal(t, "105", result.Data[0].PeginVolume.String())
}
