package pkg_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
)

func TestToAvailableLiquidityDTO(t *testing.T) {
	peginLiquidity := new(big.Int)
	peginLiquidity.SetString("1234567890987654321", 10)
	pegoutLiquidity := new(big.Int)
	pegoutLiquidity.SetString("9876543210123456789", 10)

	liquidity := liquidity_provider.AvailableLiquidity{
		PeginLiquidity:  entities.NewBigWei(peginLiquidity),
		PegoutLiquidity: entities.NewBigWei(pegoutLiquidity),
	}
	dto := pkg.ToAvailableLiquidityDTO(liquidity)
	assert.Equal(t, "1234567890987654321", dto.PeginLiquidityAmount.String())
	assert.Equal(t, "9876543210123456789", dto.PegoutLiquidityAmount.String())
}

func TestFromPeginConfigurationDTO(t *testing.T) {
	dto := pkg.PeginConfigurationDTO{
		TimeForDeposit: 10,
		CallTime:       200,
		PenaltyFee:     "3000000000000000000000",
		FixedFee:       "5000000000000000000000",
		FeePercentage:  5.443321101,
		MaxValue:       "7000000000000000000000",
		MinValue:       "6000000000000000000000",
	}
	configuration := pkg.FromPeginConfigurationDTO(dto)
	assert.Equal(t, uint32(10), configuration.TimeForDeposit)
	assert.Equal(t, uint32(200), configuration.CallTime)
	assert.Equal(t, "3000000000000000000000", configuration.PenaltyFee.AsBigInt().String())
	assert.Equal(t, "5000000000000000000000", configuration.FixedFee.AsBigInt().String())
	assert.Equal(t, "5.443321101", configuration.FeePercentage.Native().String())
	assert.Equal(t, "7000000000000000000000", configuration.MaxValue.AsBigInt().String())
	assert.Equal(t, "6000000000000000000000", configuration.MinValue.AsBigInt().String())
	test.AssertNonZeroValues(t, dto)
}

func TestFromPegoutConfigurationDTO(t *testing.T) {
	dto := pkg.PegoutConfigurationDTO{
		TimeForDeposit:       10,
		ExpireTime:           200,
		PenaltyFee:           "3000000000000000000000",
		FixedFee:             "5000000000000000000000",
		FeePercentage:        0.5123333,
		MaxValue:             "7000000000000000000000",
		MinValue:             "6000000000000000000000",
		ExpireBlocks:         20,
		BridgeTransactionMin: "8000000000000000000000",
	}
	configuration := pkg.FromPegoutConfigurationDTO(dto)
	assert.Equal(t, uint32(10), configuration.TimeForDeposit)
	assert.Equal(t, uint32(200), configuration.ExpireTime)
	assert.Equal(t, "3000000000000000000000", configuration.PenaltyFee.AsBigInt().String())
	assert.Equal(t, "5000000000000000000000", configuration.FixedFee.AsBigInt().String())
	assert.Equal(t, "0.5123333", configuration.FeePercentage.Native().String())
	assert.Equal(t, "7000000000000000000000", configuration.MaxValue.AsBigInt().String())
	assert.Equal(t, "6000000000000000000000", configuration.MinValue.AsBigInt().String())
	assert.Equal(t, uint64(20), configuration.ExpireBlocks)
	assert.Equal(t, "8000000000000000000000", configuration.BridgeTransactionMin.AsBigInt().String())
	test.AssertNonZeroValues(t, dto)
}

func TestToServerInfoDTO(t *testing.T) {
	serverInfo := liquidity_provider.ServerInfo{
		Version:  "1.0.0",
		Revision: "1234567890",
	}
	dto := pkg.ToServerInfoDTO(serverInfo)
	assert.Equal(t, "1.0.0", dto.Version)
	assert.Equal(t, "1234567890", dto.Revision)
}

// nolint:funlen
func TestLocalLiquidityProvider_ProviderDTOValidation(t *testing.T) {
	t.Run("Test FromPegoutConfigurationDTO conversion", func(t *testing.T) {
		dto := pkg.PegoutConfigurationDTO{
			TimeForDeposit:       3600,
			ExpireTime:           7200,
			PenaltyFee:           "1000000000000000",
			FixedFee:             "2000000000000000",
			FeePercentage:        1.5,
			MaxValue:             "1000000000000000000",
			MinValue:             "100000000000000000",
			ExpireBlocks:         500,
			BridgeTransactionMin: "50000000000000000",
		}
		penaltyFeeBigInt := new(big.Int)
		penaltyFeeBigInt.SetString(dto.PenaltyFee, 10)
		fixedFeeBigInt := new(big.Int)
		fixedFeeBigInt.SetString(dto.FixedFee, 10)
		maxValueBigInt := new(big.Int)
		maxValueBigInt.SetString(dto.MaxValue, 10)
		minValueBigInt := new(big.Int)
		minValueBigInt.SetString(dto.MinValue, 10)
		bridgeTransactionMinBigInt := new(big.Int)
		bridgeTransactionMinBigInt.SetString(dto.BridgeTransactionMin, 10)
		expectedConfig := liquidity_provider.PegoutConfiguration{
			TimeForDeposit:       dto.TimeForDeposit,
			ExpireTime:           dto.ExpireTime,
			PenaltyFee:           entities.NewBigWei(penaltyFeeBigInt),
			FixedFee:             entities.NewBigWei(fixedFeeBigInt),
			FeePercentage:        utils.NewBigFloat64(dto.FeePercentage),
			MaxValue:             entities.NewBigWei(maxValueBigInt),
			MinValue:             entities.NewBigWei(minValueBigInt),
			ExpireBlocks:         dto.ExpireBlocks,
			BridgeTransactionMin: entities.NewBigWei(bridgeTransactionMinBigInt),
		}
		config := pkg.FromPegoutConfigurationDTO(dto)
		assert.Equal(t, expectedConfig, config)
		test.AssertNonZeroValues(t, dto)
	})
	t.Run("Test ToPeginConfigurationDTO conversion", func(t *testing.T) {
		config := liquidity_provider.PeginConfiguration{
			TimeForDeposit: 3600,
			CallTime:       7200,
			PenaltyFee:     entities.NewWei(1000000000000000),
			FixedFee:       entities.NewWei(2000000000000000),
			FeePercentage:  utils.NewBigFloat64(1.5),
			MaxValue:       entities.NewWei(1000000000000000000),
			MinValue:       entities.NewWei(100000000000000000),
		}
		dto := pkg.ToPeginConfigurationDTO(config)
		feePercentage, _ := config.FeePercentage.Native().Float64()
		expectedDTO := pkg.PeginConfigurationDTO{
			TimeForDeposit: config.TimeForDeposit,
			CallTime:       config.CallTime,
			PenaltyFee:     config.PenaltyFee.AsBigInt().String(),
			FixedFee:       config.FixedFee.AsBigInt().String(),
			FeePercentage:  feePercentage,
			MaxValue:       config.MaxValue.AsBigInt().String(),
			MinValue:       config.MinValue.AsBigInt().String(),
		}
		assert.Equal(t, expectedDTO, dto)
		test.AssertNonZeroValues(t, dto)
	})
}

// nolint:funlen
func TestGetReportsByPeriodRequest_ValidateGetReportsByPeriodRequest(t *testing.T) {
	t.Run("valid dates", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-02",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		assert.NoError(t, err)
	})

	t.Run("invalid startDate format", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "01/01/2023",
				EndDate:   "2023-01-02",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.Error(t, err)
		assert.Equal(t, "startDate invalid date format: must be YYYY-MM-DD or ISO 8601 UTC format (ending with Z)", err.Error())
	})

	t.Run("invalid endDate format", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-01",
				EndDate:   "01/02/2023",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.Error(t, err)
		assert.Equal(t, "endDate invalid date format: must be YYYY-MM-DD or ISO 8601 UTC format (ending with Z)", err.Error())
	})

	t.Run("endDate equal to startDate", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-01",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.NoError(t, err)
	})

	t.Run("endDate before startDate", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-02",
				EndDate:   "2023-01-01",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.Error(t, err)
		assert.Equal(t, "endDate must be on or after startDate", err.Error())
	})
}

// nolint:funlen
func TestGetReportsByPeriodRequest_GetTimestamps(t *testing.T) {
	t.Run("valid dates", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-02",
			},
		}
		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		expectedStartTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		expectedEndTime := time.Date(2023, 1, 2, 23, 59, 59, 999999999, time.UTC)

		assert.Equal(t, expectedStartTime, startTime)
		assert.Equal(t, expectedEndTime, endTime)
	})

	t.Run("invalid startDate format", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "01/01/2023",
				EndDate:   "2023-01-02",
			},
		}
		startTime, endTime, err := request.GetTimestamps()
		require.Error(t, err)
		assert.True(t, startTime.IsZero())
		assert.True(t, endTime.IsZero())
	})

	t.Run("invalid endDate format", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-01",
				EndDate:   "01/02/2023",
			},
		}
		startTime, endTime, err := request.GetTimestamps()
		require.Error(t, err)
		assert.True(t, startTime.IsZero())
		assert.True(t, endTime.IsZero())
	})

	t.Run("sets time component correctly", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-05-15",
				EndDate:   "2023-06-20",
			},
		}
		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		// Start time should be at 00:00:00
		assert.Equal(t, 0, startTime.Hour())
		assert.Equal(t, 0, startTime.Minute())
		assert.Equal(t, 0, startTime.Second())
		assert.Equal(t, 0, startTime.Nanosecond())

		// End time should be at 23:59:59
		assert.Equal(t, 23, endTime.Hour())
		assert.Equal(t, 59, endTime.Minute())
		assert.Equal(t, 59, endTime.Second())
		assert.Equal(t, 999999999, endTime.Nanosecond())

		// Dates should be preserved
		assert.Equal(t, 2023, startTime.Year())
		assert.Equal(t, time.Month(5), startTime.Month())
		assert.Equal(t, 15, startTime.Day())

		assert.Equal(t, 2023, endTime.Year())
		assert.Equal(t, time.Month(6), endTime.Month())
		assert.Equal(t, 20, endTime.Day())
	})
}

// TestGetReportsByPeriodRequest_DualFormatSupport tests the dual datetime format support
// nolint:funlen
func TestGetReportsByPeriodRequest_DualFormatSupport(t *testing.T) {
	t.Run("ISO 8601 format - basic", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15T09:30:00Z",
				EndDate:   "2023-01-15T17:45:00Z",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.NoError(t, err)

		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		expectedStartTime := time.Date(2023, 1, 15, 9, 30, 0, 0, time.UTC)
		expectedEndTime := time.Date(2023, 1, 15, 17, 45, 0, 0, time.UTC)

		assert.Equal(t, expectedStartTime, startTime)
		assert.Equal(t, expectedEndTime, endTime)
	})

	t.Run("ISO 8601 format - with milliseconds", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15T09:30:00.123Z",
				EndDate:   "2023-01-15T17:45:00.999Z",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.NoError(t, err)

		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		expectedStartTime := time.Date(2023, 1, 15, 9, 30, 0, 123000000, time.UTC)
		expectedEndTime := time.Date(2023, 1, 15, 17, 45, 0, 999000000, time.UTC)

		assert.Equal(t, expectedStartTime, startTime)
		assert.Equal(t, expectedEndTime, endTime)
	})

	t.Run("ISO 8601 format - with microseconds", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15T09:30:00.123456Z",
				EndDate:   "2023-01-15T17:45:00.999999Z",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.NoError(t, err)

		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		expectedStartTime := time.Date(2023, 1, 15, 9, 30, 0, 123456000, time.UTC)
		expectedEndTime := time.Date(2023, 1, 15, 17, 45, 0, 999999000, time.UTC)

		assert.Equal(t, expectedStartTime, startTime)
		assert.Equal(t, expectedEndTime, endTime)
	})

	t.Run("ISO 8601 format - with nanoseconds", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15T09:30:00.123456789Z",
				EndDate:   "2023-01-15T17:45:00.999999999Z",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.NoError(t, err)

		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		expectedStartTime := time.Date(2023, 1, 15, 9, 30, 0, 123456789, time.UTC)
		expectedEndTime := time.Date(2023, 1, 15, 17, 45, 0, 999999999, time.UTC)

		assert.Equal(t, expectedStartTime, startTime)
		assert.Equal(t, expectedEndTime, endTime)
	})

	t.Run("unsupported format - without Z suffix", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15T09:30:00",
				EndDate:   "2023-01-15T17:45:00Z",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "startDate invalid date format")
	})

	t.Run("mixed format usage - date start, precise end", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15",
				EndDate:   "2023-01-15T23:59:59Z",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.NoError(t, err)

		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		expectedStartTime := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
		expectedEndTime := time.Date(2023, 1, 15, 23, 59, 59, 0, time.UTC)

		assert.Equal(t, expectedStartTime, startTime)
		assert.Equal(t, expectedEndTime, endTime)
	})

	t.Run("mixed format usage - precise start, date end", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15T00:00:00Z",
				EndDate:   "2023-01-15",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.NoError(t, err)

		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		expectedStartTime := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
		expectedEndTime := time.Date(2023, 1, 15, 23, 59, 59, 999999999, time.UTC)

		assert.Equal(t, expectedStartTime, startTime)
		assert.Equal(t, expectedEndTime, endTime)
	})

	t.Run("same-day query with different formats", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15",
				EndDate:   "2023-01-15T23:59:59.999Z",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.NoError(t, err) // Should now be allowed

		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		expectedStartTime := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)
		expectedEndTime := time.Date(2023, 1, 15, 23, 59, 59, 999000000, time.UTC)

		assert.Equal(t, expectedStartTime, startTime)
		assert.Equal(t, expectedEndTime, endTime)
	})

	t.Run("unsupported timezone format - RFC3339 with offset", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15T09:30:00+02:00",
				EndDate:   "2023-01-15T17:45:00Z",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "startDate invalid date format")
	})

	t.Run("YYYY-MM-DD format", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-01",
				EndDate:   "2023-01-31",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.NoError(t, err)

		startTime, endTime, err := request.GetTimestamps()
		require.NoError(t, err)

		expectedStartTime := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		expectedEndTime := time.Date(2023, 1, 31, 23, 59, 59, 999999999, time.UTC)

		assert.Equal(t, expectedStartTime, startTime)
		assert.Equal(t, expectedEndTime, endTime)
	})

	t.Run("invalid ISO 8601 format", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "2023-01-15T25:30:00Z", // Invalid hour
				EndDate:   "2023-01-15T17:45:00Z",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "startDate invalid date format")
	})

	t.Run("completely invalid format", func(t *testing.T) {
		request := pkg.GetReportsByPeriodRequest{
			DateRangeRequest: pkg.DateRangeRequest{
				StartDate: "not-a-date",
				EndDate:   "2023-01-15",
			},
		}
		err := request.ValidateGetReportsByPeriodRequest()
		require.Error(t, err)
		assert.Contains(t, err.Error(), "startDate invalid date format")
	})
}

func TestToSummaryDataDTO(t *testing.T) {
	data := reports.SummaryData{
		TotalQuotesCount:          10,
		AcceptedQuotesCount:       8,
		TotalAcceptedQuotesAmount: entities.NewWei(2500),
		PaidQuotesCount:           5,
		PaidQuotesAmount:          entities.NewWei(1500),
		RefundedQuotesCount:       2,
		TotalRefundedQuotesAmount: entities.NewWei(800),
		PenalizationsCount:        3,
		TotalPenalizationsAmount:  entities.NewWei(50),
	}
	dto := pkg.ToSummaryDataDTO(data)
	assert.Equal(t, data.TotalQuotesCount, dto.TotalQuotesCount)
	assert.Equal(t, data.AcceptedQuotesCount, dto.AcceptedQuotesCount)
	assert.Equal(t, 0, data.TotalAcceptedQuotesAmount.AsBigInt().Cmp(dto.TotalAcceptedQuotesAmount))
	assert.Equal(t, data.PaidQuotesCount, dto.PaidQuotesCount)
	assert.Equal(t, 0, data.PaidQuotesAmount.AsBigInt().Cmp(dto.PaidQuotesAmount))
	assert.Equal(t, data.RefundedQuotesCount, dto.RefundedQuotesCount)
	assert.Equal(t, 0, data.TotalRefundedQuotesAmount.AsBigInt().Cmp(dto.TotalRefundedQuotesAmount))
	assert.Equal(t, data.PenalizationsCount, dto.PenalizationsCount)
	assert.Equal(t, 0, data.TotalPenalizationsAmount.AsBigInt().Cmp(dto.TotalPenalizationsAmount))
	test.AssertNonZeroValues(t, dto)
}

func TestToSummaryResultDTO(t *testing.T) {
	pegin := reports.SummaryData{
		TotalQuotesCount:          3,
		AcceptedQuotesCount:       2,
		TotalAcceptedQuotesAmount: entities.NewWei(800),
		PaidQuotesCount:           1,
		PaidQuotesAmount:          entities.NewWei(500),
		RefundedQuotesCount:       1,
		TotalRefundedQuotesAmount: entities.NewWei(300),
		PenalizationsCount:        1,
		TotalPenalizationsAmount:  entities.NewWei(15),
	}
	pegout := reports.SummaryData{
		TotalQuotesCount:          4,
		AcceptedQuotesCount:       3,
		TotalAcceptedQuotesAmount: entities.NewWei(1200),
		PaidQuotesCount:           2,
		PaidQuotesAmount:          entities.NewWei(900),
		RefundedQuotesCount:       1,
		TotalRefundedQuotesAmount: entities.NewWei(400),
		PenalizationsCount:        1,
		TotalPenalizationsAmount:  entities.NewWei(20),
	}
	result := reports.SummaryResult{PeginSummary: pegin, PegoutSummary: pegout}
	dto := pkg.ToSummaryResultDTO(result)
	expected := pkg.SummaryResultDTO{
		PeginSummary:  pkg.ToSummaryDataDTO(pegin),
		PegoutSummary: pkg.ToSummaryDataDTO(pegout),
	}
	assert.Equal(t, expected, dto)
	test.AssertMaxZeroValues(t, dto.PeginSummary, 1)
	test.AssertNonZeroValues(t, dto.PegoutSummary)
}

func TestToTrustedAccountDTO(t *testing.T) {
	btcLockingCap := new(big.Int)
	btcLockingCap.SetString("5000000000000000000", 10)
	rbtcLockingCap := new(big.Int)
	rbtcLockingCap.SetString("7000000000000000000", 10)
	trustedAccount := liquidity_provider.TrustedAccountDetails{
		Address:        "0x1234567890abcdef",
		Name:           "Test Trusted Account",
		BtcLockingCap:  entities.NewBigWei(btcLockingCap),
		RbtcLockingCap: entities.NewBigWei(rbtcLockingCap),
	}
	dto := pkg.ToTrustedAccountDTO(trustedAccount)
	assert.Equal(t, "0x1234567890abcdef", dto.Address)
	assert.Equal(t, "Test Trusted Account", dto.Name)
	assert.Equal(t, "5000000000000000000", dto.BtcLockingCap.String())
	assert.Equal(t, "7000000000000000000", dto.RbtcLockingCap.String())
}

func TestToTrustedAccountsDTO(t *testing.T) {
	btcLockingCap1 := new(big.Int)
	btcLockingCap1.SetString("5000000000000000000", 10)
	rbtcLockingCap1 := new(big.Int)
	rbtcLockingCap1.SetString("7000000000000000000", 10)
	btcLockingCap2 := new(big.Int)
	btcLockingCap2.SetString("9000000000000000000", 10)
	rbtcLockingCap2 := new(big.Int)
	rbtcLockingCap2.SetString("3000000000000000000", 10)
	account1 := liquidity_provider.TrustedAccountDetails{
		Address:        "0x1234567890abcdef",
		Name:           "Test Trusted Account 1",
		BtcLockingCap:  entities.NewBigWei(btcLockingCap1),
		RbtcLockingCap: entities.NewBigWei(rbtcLockingCap1),
	}
	account2 := liquidity_provider.TrustedAccountDetails{
		Address:        "0xabcdef1234567890",
		Name:           "Test Trusted Account 2",
		BtcLockingCap:  entities.NewBigWei(btcLockingCap2),
		RbtcLockingCap: entities.NewBigWei(rbtcLockingCap2),
	}

	signedAccounts := []entities.Signed[liquidity_provider.TrustedAccountDetails]{
		{
			Value:     account1,
			Signature: "signature1",
			Hash:      "hash1",
		},
		{
			Value:     account2,
			Signature: "signature2",
			Hash:      "hash2",
		},
	}

	dtos := pkg.ToTrustedAccountsDTO(signedAccounts)
	assert.Len(t, dtos, 2)
	assert.Equal(t, "0x1234567890abcdef", dtos[0].Address)
	assert.Equal(t, "Test Trusted Account 1", dtos[0].Name)
	assert.Equal(t, "5000000000000000000", dtos[0].BtcLockingCap.String())
	assert.Equal(t, "7000000000000000000", dtos[0].RbtcLockingCap.String())
	assert.Equal(t, "0xabcdef1234567890", dtos[1].Address)
	assert.Equal(t, "Test Trusted Account 2", dtos[1].Name)
	assert.Equal(t, "9000000000000000000", dtos[1].BtcLockingCap.String())
	assert.Equal(t, "3000000000000000000", dtos[1].RbtcLockingCap.String())
}

func TestFromGeneralConfigurationDTO(t *testing.T) {
	t.Run("converts valid configuration", func(t *testing.T) {
		dto := pkg.GeneralConfigurationDTO{
			RskConfirmations: map[string]uint16{
				"1000000000000000000": 5,
				"2000000000000000000": 10,
			},
			BtcConfirmations: map[string]uint16{
				"3000000000000000000": 15,
				"4000000000000000000": 20,
			},
			PublicLiquidityCheck: true,
		}

		config, err := pkg.FromGeneralConfigurationDTO(dto)

		require.NoError(t, err)
		assert.Equal(t, dto.RskConfirmations, map[string]uint16(config.RskConfirmations))
		assert.Equal(t, dto.BtcConfirmations, map[string]uint16(config.BtcConfirmations))
		assert.Equal(t, dto.PublicLiquidityCheck, config.PublicLiquidityCheck)
		test.AssertNonZeroValues(t, dto)
	})

	t.Run("returns error on invalid numeric keys", func(t *testing.T) {
		invalidBtc := pkg.GeneralConfigurationDTO{
			RskConfirmations: map[string]uint16{
				"1000000000000000000": 5,
			},
			BtcConfirmations: map[string]uint16{
				"3000000000000000000": 15,
				"notanumber":          20,
			},
			PublicLiquidityCheck: true,
		}
		invalidRsk := pkg.GeneralConfigurationDTO{
			RskConfirmations: map[string]uint16{
				"1000000000000000000": 5,
				"invalid":             10,
			},
			BtcConfirmations: map[string]uint16{
				"3000000000000000000": 15,
			},
			PublicLiquidityCheck: false,
		}

		config, err := pkg.FromGeneralConfigurationDTO(invalidBtc)
		assert.Empty(t, config)
		require.ErrorContains(t, err, "cannot deserialize BTC confirmations key notanumber")

		config, err = pkg.FromGeneralConfigurationDTO(invalidRsk)
		assert.Empty(t, config)
		require.ErrorContains(t, err, "cannot deserialize RSK confirmations key invalid")
	})
}

func TestToBtcAssetLocationDTO(t *testing.T) {
	location := reports.BtcAssetLocation{
		BtcWallet:  entities.NewWei(100000000),
		Federation: entities.NewWei(200000000),
		RskWallet:  entities.NewWei(300000000),
		Lbc:        entities.NewWei(50000000),
	}

	dto := pkg.ToBtcAssetLocationDTO(location)

	assert.Equal(t, "100000000", dto.BtcWallet.String())
	assert.Equal(t, "200000000", dto.Federation.String())
	assert.Equal(t, "300000000", dto.RskWallet.String())
	assert.Equal(t, "50000000", dto.Lbc.String())
}

func TestToBtcAssetAllocationDTO(t *testing.T) {
	allocation := reports.BtcAssetAllocation{
		ReservedForUsers: entities.NewWei(150000000),
		WaitingForRefund: entities.NewWei(75000000),
		Available:        entities.NewWei(425000000),
	}

	dto := pkg.ToBtcAssetAllocationDTO(allocation)

	assert.Equal(t, "150000000", dto.ReservedForUsers.String())
	assert.Equal(t, "75000000", dto.WaitingForRefund.String())
	assert.Equal(t, "425000000", dto.Available.String())
}

func TestToBtcAssetReportDTO(t *testing.T) {
	report := reports.BtcAssetReport{
		Total: entities.NewWei(650000000),
		Location: reports.BtcAssetLocation{
			BtcWallet:  entities.NewWei(100000000),
			Federation: entities.NewWei(200000000),
			RskWallet:  entities.NewWei(300000000),
			Lbc:        entities.NewWei(50000000),
		},
		Allocation: reports.BtcAssetAllocation{
			ReservedForUsers: entities.NewWei(150000000),
			WaitingForRefund: entities.NewWei(75000000),
			Available:        entities.NewWei(425000000),
		},
	}

	dto := pkg.ToBtcAssetReportDTO(report)

	assert.Equal(t, "650000000", dto.Total.String())
	assert.Equal(t, "100000000", dto.Location.BtcWallet.String())
	assert.Equal(t, "200000000", dto.Location.Federation.String())
	assert.Equal(t, "300000000", dto.Location.RskWallet.String())
	assert.Equal(t, "50000000", dto.Location.Lbc.String())
	assert.Equal(t, "150000000", dto.Allocation.ReservedForUsers.String())
	assert.Equal(t, "75000000", dto.Allocation.WaitingForRefund.String())
	assert.Equal(t, "425000000", dto.Allocation.Available.String())
}

func TestToRbtcAssetLocationDTO(t *testing.T) {
	location := reports.RbtcAssetLocation{
		RskWallet:  entities.NewWei(500000000000000000),
		Lbc:        entities.NewWei(250000000000000000),
		Federation: entities.NewWei(100000000000000000),
	}

	dto := pkg.ToRbtcAssetLocationDTO(location)

	assert.Equal(t, "500000000000000000", dto.RskWallet.String())
	assert.Equal(t, "250000000000000000", dto.Lbc.String())
	assert.Equal(t, "100000000000000000", dto.Federation.String())
}

func TestToRbtcAssetAllocationDTO(t *testing.T) {
	allocation := reports.RbtcAssetAllocation{
		ReservedForUsers: entities.NewWei(300000000000000000),
		WaitingForRefund: entities.NewWei(150000000000000000),
		Available:        entities.NewWei(400000000000000000),
	}

	dto := pkg.ToRbtcAssetAllocationDTO(allocation)

	assert.Equal(t, "300000000000000000", dto.ReservedForUsers.String())
	assert.Equal(t, "150000000000000000", dto.WaitingForRefund.String())
	assert.Equal(t, "400000000000000000", dto.Available.String())
}

func TestToRbtcAssetReportDTO(t *testing.T) {
	report := reports.RbtcAssetReport{
		Total: entities.NewWei(850000000000000000),
		Location: reports.RbtcAssetLocation{
			RskWallet:  entities.NewWei(500000000000000000),
			Lbc:        entities.NewWei(250000000000000000),
			Federation: entities.NewWei(100000000000000000),
		},
		Allocation: reports.RbtcAssetAllocation{
			ReservedForUsers: entities.NewWei(300000000000000000),
			WaitingForRefund: entities.NewWei(150000000000000000),
			Available:        entities.NewWei(400000000000000000),
		},
	}

	dto := pkg.ToRbtcAssetReportDTO(report)

	assert.Equal(t, "850000000000000000", dto.Total.String())
	assert.Equal(t, "500000000000000000", dto.Location.RskWallet.String())
	assert.Equal(t, "250000000000000000", dto.Location.Lbc.String())
	assert.Equal(t, "100000000000000000", dto.Location.Federation.String())
	assert.Equal(t, "300000000000000000", dto.Allocation.ReservedForUsers.String())
	assert.Equal(t, "150000000000000000", dto.Allocation.WaitingForRefund.String())
	assert.Equal(t, "400000000000000000", dto.Allocation.Available.String())
}

func TestToGetAssetsReportResponse(t *testing.T) {
	result := reports.GetAssetsReportResult{
		BtcAssetReport: reports.BtcAssetReport{
			Total: entities.NewWei(650000000),
			Location: reports.BtcAssetLocation{
				BtcWallet:  entities.NewWei(100000000),
				Federation: entities.NewWei(200000000),
				RskWallet:  entities.NewWei(300000000),
				Lbc:        entities.NewWei(50000000),
			},
			Allocation: reports.BtcAssetAllocation{
				ReservedForUsers: entities.NewWei(150000000),
				WaitingForRefund: entities.NewWei(75000000),
				Available:        entities.NewWei(425000000),
			},
		},
		RbtcAssetReport: reports.RbtcAssetReport{
			Total: entities.NewWei(850000000000000000),
			Location: reports.RbtcAssetLocation{
				RskWallet:  entities.NewWei(500000000000000000),
				Lbc:        entities.NewWei(250000000000000000),
				Federation: entities.NewWei(100000000000000000),
			},
			Allocation: reports.RbtcAssetAllocation{
				ReservedForUsers: entities.NewWei(300000000000000000),
				WaitingForRefund: entities.NewWei(150000000000000000),
				Available:        entities.NewWei(400000000000000000),
			},
		},
	}

	response := pkg.ToGetAssetsReportResponse(result)

	assert.Equal(t, "650000000", response.BtcAssetReport.Total.String())
	assert.Equal(t, "100000000", response.BtcAssetReport.Location.BtcWallet.String())
	assert.Equal(t, "200000000", response.BtcAssetReport.Location.Federation.String())
	assert.Equal(t, "300000000", response.BtcAssetReport.Location.RskWallet.String())
	assert.Equal(t, "50000000", response.BtcAssetReport.Location.Lbc.String())
	assert.Equal(t, "150000000", response.BtcAssetReport.Allocation.ReservedForUsers.String())
	assert.Equal(t, "75000000", response.BtcAssetReport.Allocation.WaitingForRefund.String())
	assert.Equal(t, "425000000", response.BtcAssetReport.Allocation.Available.String())
	assert.Equal(t, "850000000000000000", response.RbtcAssetReport.Total.String())
	assert.Equal(t, "500000000000000000", response.RbtcAssetReport.Location.RskWallet.String())
	assert.Equal(t, "250000000000000000", response.RbtcAssetReport.Location.Lbc.String())
	assert.Equal(t, "100000000000000000", response.RbtcAssetReport.Location.Federation.String())
	assert.Equal(t, "300000000000000000", response.RbtcAssetReport.Allocation.ReservedForUsers.String())
	assert.Equal(t, "150000000000000000", response.RbtcAssetReport.Allocation.WaitingForRefund.String())
	assert.Equal(t, "400000000000000000", response.RbtcAssetReport.Allocation.Available.String())
}

func TestToGetAssetsReportResponse_WithLargeValues(t *testing.T) {
	btcTotal := new(big.Int)
	btcTotal.SetString("21000000000000000000", 10)
	rbtcTotal := new(big.Int)
	rbtcTotal.SetString("999999999999999999999999", 10)

	result := reports.GetAssetsReportResult{
		BtcAssetReport: reports.BtcAssetReport{
			Total: entities.NewBigWei(btcTotal),
			Location: reports.BtcAssetLocation{
				BtcWallet:  entities.NewBigWei(btcTotal),
				Federation: entities.NewWei(0),
				RskWallet:  entities.NewWei(0),
				Lbc:        entities.NewWei(0),
			},
			Allocation: reports.BtcAssetAllocation{
				ReservedForUsers: entities.NewWei(0),
				WaitingForRefund: entities.NewWei(0),
				Available:        entities.NewBigWei(btcTotal),
			},
		},
		RbtcAssetReport: reports.RbtcAssetReport{
			Total: entities.NewBigWei(rbtcTotal),
			Location: reports.RbtcAssetLocation{
				RskWallet:  entities.NewBigWei(rbtcTotal),
				Lbc:        entities.NewWei(0),
				Federation: entities.NewWei(0),
			},
			Allocation: reports.RbtcAssetAllocation{
				ReservedForUsers: entities.NewWei(0),
				WaitingForRefund: entities.NewWei(0),
				Available:        entities.NewBigWei(rbtcTotal),
			},
		},
	}

	response := pkg.ToGetAssetsReportResponse(result)

	assert.Equal(t, "21000000000000000000", response.BtcAssetReport.Total.String())
	assert.Equal(t, "21000000000000000000", response.BtcAssetReport.Location.BtcWallet.String())
	assert.Equal(t, "21000000000000000000", response.BtcAssetReport.Allocation.Available.String())
	assert.Equal(t, "999999999999999999999999", response.RbtcAssetReport.Total.String())
	assert.Equal(t, "999999999999999999999999", response.RbtcAssetReport.Location.RskWallet.String())
	assert.Equal(t, "999999999999999999999999", response.RbtcAssetReport.Allocation.Available.String())
}

func TestToGetAssetsReportResponse_WithZeroValues(t *testing.T) {
	result := reports.GetAssetsReportResult{
		BtcAssetReport: reports.BtcAssetReport{
			Total: entities.NewWei(0),
			Location: reports.BtcAssetLocation{
				BtcWallet:  entities.NewWei(0),
				Federation: entities.NewWei(0),
				RskWallet:  entities.NewWei(0),
				Lbc:        entities.NewWei(0),
			},
			Allocation: reports.BtcAssetAllocation{
				ReservedForUsers: entities.NewWei(0),
				WaitingForRefund: entities.NewWei(0),
				Available:        entities.NewWei(0),
			},
		},
		RbtcAssetReport: reports.RbtcAssetReport{
			Total: entities.NewWei(0),
			Location: reports.RbtcAssetLocation{
				RskWallet:  entities.NewWei(0),
				Lbc:        entities.NewWei(0),
				Federation: entities.NewWei(0),
			},
			Allocation: reports.RbtcAssetAllocation{
				ReservedForUsers: entities.NewWei(0),
				WaitingForRefund: entities.NewWei(0),
				Available:        entities.NewWei(0),
			},
		},
	}

	response := pkg.ToGetAssetsReportResponse(result)

	assert.Equal(t, "0", response.BtcAssetReport.Total.String())
	assert.Equal(t, "0", response.BtcAssetReport.Location.BtcWallet.String())
	assert.Equal(t, "0", response.BtcAssetReport.Location.Federation.String())
	assert.Equal(t, "0", response.BtcAssetReport.Location.RskWallet.String())
	assert.Equal(t, "0", response.BtcAssetReport.Location.Lbc.String())
	assert.Equal(t, "0", response.BtcAssetReport.Allocation.ReservedForUsers.String())
	assert.Equal(t, "0", response.BtcAssetReport.Allocation.WaitingForRefund.String())
	assert.Equal(t, "0", response.BtcAssetReport.Allocation.Available.String())
	assert.Equal(t, "0", response.RbtcAssetReport.Total.String())
	assert.Equal(t, "0", response.RbtcAssetReport.Location.RskWallet.String())
	assert.Equal(t, "0", response.RbtcAssetReport.Location.Lbc.String())
	assert.Equal(t, "0", response.RbtcAssetReport.Location.Federation.String())
	assert.Equal(t, "0", response.RbtcAssetReport.Allocation.ReservedForUsers.String())
	assert.Equal(t, "0", response.RbtcAssetReport.Allocation.WaitingForRefund.String())
	assert.Equal(t, "0", response.RbtcAssetReport.Allocation.Available.String())
}
