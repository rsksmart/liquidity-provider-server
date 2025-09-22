package pkg

import (
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	reports "github.com/rsksmart/liquidity-provider-server/internal/usecases/reports"
)

type ProviderDetail struct {
	// Deprecated: Fee is deprecated, use FixedFee and FeePercentage instead
	Fee                   *big.Int `json:"fee" required:""`
	FixedFee              *big.Int `json:"fixedFee"  required:""`
	FeePercentage         float64  `json:"feePercentage"  required:""`
	MinTransactionValue   *big.Int `json:"minTransactionValue"  required:""`
	MaxTransactionValue   *big.Int `json:"maxTransactionValue"  required:""`
	RequiredConfirmations uint16   `json:"requiredConfirmations"  required:""`
}

type ProviderDetailResponse struct {
	SiteKey               string         `json:"siteKey" required:""`
	LiquidityCheckEnabled bool           `json:"liquidityCheckEnabled" required:""`
	Pegin                 ProviderDetail `json:"pegin" required:""`
	Pegout                ProviderDetail `json:"pegout" required:""`
}

type LiquidityProvider struct {
	Id           uint64 `json:"id" example:"1" description:"Provider Id"  required:""`
	Provider     string `json:"provider" example:"0x0" description:"Provider Address"  required:""`
	Name         string `json:"name" example:"New Provider" description:"Provider Name"  required:""`
	ApiBaseUrl   string `json:"apiBaseUrl" example:"https://api.example.com" description:"API base URL"  required:""`
	Status       bool   `json:"status" example:"true" description:"Provider status"  required:""`
	ProviderType string `json:"providerType" example:"pegin" description:"Provider type"  required:""`
}

type ChangeStatusRequest struct {
	Status *bool `json:"status"`
}

type PeginConfigurationRequest struct {
	Configuration PeginConfigurationDTO `json:"configuration" validate:"required"`
}

type PeginConfigurationDTO struct {
	TimeForDeposit uint32  `json:"timeForDeposit" validate:"required"`
	CallTime       uint32  `json:"callTime" validate:"required"`
	PenaltyFee     string  `json:"penaltyFee" validate:"required,numeric,positive_string"`
	FixedFee       string  `json:"fixedFee" validate:"required,numeric,min=0"`
	FeePercentage  float64 `json:"feePercentage" validate:"numeric,gte=0,lte=100,max_decimal_places=2"`
	MaxValue       string  `json:"maxValue" validate:"required,numeric,positive_string"`
	MinValue       string  `json:"minValue" validate:"required,numeric,positive_string"`
}

type PegoutConfigurationRequest struct {
	Configuration PegoutConfigurationDTO `json:"configuration" validate:"required"`
}

type PegoutConfigurationDTO struct {
	TimeForDeposit       uint32  `json:"timeForDeposit" validate:"required"`
	ExpireTime           uint32  `json:"expireTime" validate:"required"`
	PenaltyFee           string  `json:"penaltyFee" validate:"required,numeric,positive_string"`
	FixedFee             string  `json:"fixedFee" validate:"required,numeric,min=0"`
	FeePercentage        float64 `json:"feePercentage" validate:"numeric,gte=0,lte=100,max_decimal_places=2"`
	MaxValue             string  `json:"maxValue" validate:"required,numeric,positive_string"`
	MinValue             string  `json:"minValue" validate:"required,numeric,positive_string"`
	ExpireBlocks         uint64  `json:"expireBlocks" validate:"required"`
	BridgeTransactionMin string  `json:"bridgeTransactionMin" validate:"required,numeric,positive_string"`
}

type GeneralConfigurationRequest struct {
	Configuration GeneralConfigurationDTO `json:"configuration" validate:"required"`
}

type GeneralConfigurationDTO struct {
	RskConfirmations     map[string]uint16 `json:"rskConfirmations" validate:"required,confirmations_map"`
	BtcConfirmations     map[string]uint16 `json:"btcConfirmations" validate:"required,confirmations_map"`
	PublicLiquidityCheck bool              `json:"publicLiquidityCheck" validate:""`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CredentialsUpdateRequest struct {
	OldUsername string `json:"oldUsername" validate:"required"`
	OldPassword string `json:"oldPassword" validate:"required"`
	NewUsername string `json:"newUsername" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

// DateRangeRequest provides common date range functionality with validation
type DateRangeRequest struct {
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
}

type GetReportsByPeriodRequest struct {
	DateRangeRequest
}

// parseDateTime parses a date string in either YYYY-MM-DD or ISO 8601 UTC format
// Returns the parsed time, a boolean indicating if it's a date-only format, and any error
func parseDateTime(dateStr string) (parsedTime time.Time, isDateOnly bool, err error) {
	// Try date-only format first (YYYY-MM-DD)
	if t, parseErr := time.Parse(time.DateOnly, dateStr); parseErr == nil {
		return t, true, nil
	}

	// Try ISO 8601 UTC format - must end with 'Z' to ensure UTC timezone
	if strings.HasSuffix(dateStr, "Z") {
		if t, parseErr := time.Parse(time.RFC3339Nano, dateStr); parseErr == nil {
			return t.UTC(), false, nil
		}
	}

	return time.Time{}, false, errors.New("invalid date format: must be YYYY-MM-DD or ISO 8601 UTC format (ending with Z)")
}

func (r *DateRangeRequest) ValidateDateRange() error {
	startDate, _, err := parseDateTime(r.StartDate)
	if err != nil {
		return errors.New("startDate " + err.Error())
	}

	endDate, _, err := parseDateTime(r.EndDate)
	if err != nil {
		return errors.New("endDate " + err.Error())
	}

	// Allow same-day queries and require endDate to be >= startDate
	if endDate.Before(startDate) {
		return errors.New("endDate must be on or after startDate")
	}

	return nil
}

// ValidateGetReportsByPeriodRequest validates the request using the base date range validation
func (r *GetReportsByPeriodRequest) ValidateGetReportsByPeriodRequest() error {
	return r.ValidateDateRange()
}

// GetTimestamps converts date strings to time.Time objects with proper timezone handling
func (r *DateRangeRequest) GetTimestamps() (startTime, endTime time.Time, err error) {
	startTime, _, err = parseDateTime(r.StartDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	endTime, isEndDateOnly, err := parseDateTime(r.EndDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Apply defaults for date-only formats
	// startTime is already 00:00:00 for date-only formats from parseDateTime

	if isEndDateOnly {
		// For date-only end, use 23:59:59.999 UTC
		endTime = time.Date(
			endTime.Year(),
			endTime.Month(),
			endTime.Day(),
			23, 59, 59, 999999999,
			time.UTC,
		)
	}

	return startTime, endTime, nil
}

type GetRevenueReportResponse struct {
	TotalQuoteCallFees *big.Int `json:"totalQuoteCallFees" validate:"required"`
	TotalPenalizations *big.Int `json:"totalPenalizations" validate:"required"`
	TotalProfit        *big.Int `json:"totalProfit" validate:"required"`
}

type GetAssetsReportDTO struct {
	BtcBalance    *big.Int `json:"btcBalance" example:"1000000000" description:"Current balance on the bitcoin wallet" validate:"required"`
	RbtcBalance   *big.Int `json:"rbtcBalance" example:"1000000000" description:"Current balance on the RBTC wallet" validate:"required"`
	BtcLocked     *big.Int `json:"btcLocked" example:"1000000000" description:"Amount of BTC locked by quotes" validate:"required"`
	RbtcLocked    *big.Int `json:"rbtcLocked" example:"1000000000" description:"Amount of RBTC locked by quotes" validate:"required"`
	BtcLiquidity  *big.Int `json:"btcLiquidity" example:"1000000000" description:"Amount of BTC liquidity available for new quotes" validate:"required"`
	RbtcLiquidity *big.Int `json:"rbtcLiquidity" example:"1000000000" description:"Amount of RBTC liquidity available for new quotes" validate:"required"`
}

type AvailableLiquidityDTO struct {
	PeginLiquidityAmount  *big.Int `json:"peginLiquidityAmount" example:"5000000000000000000" description:"Available liquidity for PegIn operations in wei"  required:""`
	PegoutLiquidityAmount *big.Int `json:"pegoutLiquidityAmount" example:"5000000000000000000" description:"Available liquidity for PegOut operations in wei" required:""`
}

type ServerInfoDTO struct {
	Version  string `json:"version" example:"v1.0.0" description:"Server version tag"  required:""`
	Revision string `json:"revision" example:"b7bf393a2b1cedde8ee15b00780f44e6e5d2ba9d" description:"Version commit hash"  required:""`
}

type SummaryDataDTO struct {
	TotalQuotesCount          int64    `json:"totalQuotesCount"`
	AcceptedQuotesCount       int64    `json:"acceptedQuotesCount"`
	PaidQuotesCount           int64    `json:"paidQuotesCount"`
	PaidQuotesAmount          *big.Int `json:"paidQuotesAmount"`
	TotalAcceptedQuotedAmount *big.Int `json:"totalAcceptedQuotedAmount"`
	TotalFeesCollected        *big.Int `json:"totalFeesCollected"`
	RefundedQuotesCount       int64    `json:"refundedQuotesCount"`
	TotalPenaltyAmount        *big.Int `json:"totalPenaltyAmount"`
	LpEarnings                *big.Int `json:"lpEarnings"`
}

type SummaryResultDTO struct {
	PeginSummary  SummaryDataDTO `json:"peginSummary"`
	PegoutSummary SummaryDataDTO `json:"pegoutSummary"`
}

func ToSummaryDataDTO(data reports.SummaryData) SummaryDataDTO {
	return SummaryDataDTO{
		TotalQuotesCount:          data.TotalQuotesCount,
		AcceptedQuotesCount:       data.AcceptedQuotesCount,
		PaidQuotesCount:           data.PaidQuotesCount,
		PaidQuotesAmount:          data.PaidQuotesAmount.AsBigInt(),
		TotalAcceptedQuotedAmount: data.TotalAcceptedQuotedAmount.AsBigInt(),
		TotalFeesCollected:        data.TotalFeesCollected.AsBigInt(),
		RefundedQuotesCount:       data.RefundedQuotesCount,
		TotalPenaltyAmount:        data.TotalPenaltyAmount.AsBigInt(),
		LpEarnings:                data.LpEarnings.AsBigInt(),
	}
}

func ToSummaryResultDTO(result reports.SummaryResult) SummaryResultDTO {
	return SummaryResultDTO{
		PeginSummary:  ToSummaryDataDTO(result.PeginSummary),
		PegoutSummary: ToSummaryDataDTO(result.PegoutSummary),
	}
}

type TrustedAccountDTO struct {
	Address        string   `json:"address" example:"0x1234567890abcdef" description:"Trusted account address" required:""`
	Name           string   `json:"name" example:"Example Trusted Account" description:"Trusted account name" required:""`
	BtcLockingCap  *big.Int `json:"btcLockingCap" example:"5000000000000000000" description:"Bitcoin locking capacity in wei" required:""`
	RbtcLockingCap *big.Int `json:"rbtcLockingCap" example:"5000000000000000000" description:"RBTC locking capacity in wei" required:""`
}

type TrustedAccountRequest struct {
	Address        string   `json:"address" validate:"required,eth_addr"`
	Name           string   `json:"name" validate:"required,max=100,min=1,not_blank"`
	BtcLockingCap  *big.Int `json:"btcLockingCap" validate:"required,positive_integer_bigint"`
	RbtcLockingCap *big.Int `json:"rbtcLockingCap" validate:"required,positive_integer_bigint"`
}

type TrustedAccountsResponse struct {
	Accounts []TrustedAccountDTO `json:"accounts"`
}

func ToAvailableLiquidityDTO(entity liquidity_provider.AvailableLiquidity) AvailableLiquidityDTO {
	return AvailableLiquidityDTO{
		PeginLiquidityAmount:  entity.PeginLiquidity.AsBigInt(),
		PegoutLiquidityAmount: entity.PegoutLiquidity.AsBigInt(),
	}
}

func FromPeginConfigurationDTO(dto PeginConfigurationDTO) liquidity_provider.PeginConfiguration {
	const base = 10
	penaltyFee := new(big.Int)
	penaltyFee.SetString(dto.PenaltyFee, base)
	fixedFee := new(big.Int)
	fixedFee.SetString(dto.FixedFee, base)
	maxValue := new(big.Int)
	maxValue.SetString(dto.MaxValue, base)
	minValue := new(big.Int)
	minValue.SetString(dto.MinValue, base)

	return liquidity_provider.PeginConfiguration{
		TimeForDeposit: dto.TimeForDeposit,
		CallTime:       dto.CallTime,
		PenaltyFee:     entities.NewBigWei(penaltyFee),
		FixedFee:       entities.NewBigWei(fixedFee),
		FeePercentage:  utils.NewBigFloat64(dto.FeePercentage),
		MaxValue:       entities.NewBigWei(maxValue),
		MinValue:       entities.NewBigWei(minValue),
	}
}

func FromPegoutConfigurationDTO(dto PegoutConfigurationDTO) liquidity_provider.PegoutConfiguration {
	const base = 10
	penaltyFee := new(big.Int)
	penaltyFee.SetString(dto.PenaltyFee, base)
	fixedFee := new(big.Int)
	fixedFee.SetString(dto.FixedFee, base)
	maxValue := new(big.Int)
	maxValue.SetString(dto.MaxValue, base)
	minValue := new(big.Int)
	minValue.SetString(dto.MinValue, base)
	bridgeTransactionMin := new(big.Int)
	bridgeTransactionMin.SetString(dto.BridgeTransactionMin, base)

	return liquidity_provider.PegoutConfiguration{
		TimeForDeposit:       dto.TimeForDeposit,
		ExpireTime:           dto.ExpireTime,
		PenaltyFee:           entities.NewBigWei(penaltyFee),
		FixedFee:             entities.NewBigWei(fixedFee),
		FeePercentage:        utils.NewBigFloat64(dto.FeePercentage),
		MaxValue:             entities.NewBigWei(maxValue),
		MinValue:             entities.NewBigWei(minValue),
		ExpireBlocks:         dto.ExpireBlocks,
		BridgeTransactionMin: entities.NewBigWei(bridgeTransactionMin),
	}
}

func ToPeginConfigurationDTO(config liquidity_provider.PeginConfiguration) PeginConfigurationDTO {
	feePercentage, _ := config.FeePercentage.Native().Float64()

	return PeginConfigurationDTO{
		TimeForDeposit: config.TimeForDeposit,
		CallTime:       config.CallTime,
		PenaltyFee:     config.PenaltyFee.AsBigInt().String(),
		FixedFee:       config.FixedFee.AsBigInt().String(),
		FeePercentage:  feePercentage,
		MaxValue:       config.MaxValue.AsBigInt().String(),
		MinValue:       config.MinValue.AsBigInt().String(),
	}
}

func ToPegoutConfigurationDTO(config liquidity_provider.PegoutConfiguration) PegoutConfigurationDTO {
	feePercentage, _ := config.FeePercentage.Native().Float64()
	return PegoutConfigurationDTO{
		TimeForDeposit:       config.TimeForDeposit,
		ExpireTime:           config.ExpireTime,
		PenaltyFee:           config.PenaltyFee.AsBigInt().String(),
		FixedFee:             config.FixedFee.AsBigInt().String(),
		FeePercentage:        feePercentage,
		MaxValue:             config.MaxValue.AsBigInt().String(),
		MinValue:             config.MinValue.AsBigInt().String(),
		ExpireBlocks:         config.ExpireBlocks,
		BridgeTransactionMin: config.BridgeTransactionMin.AsBigInt().String(),
	}
}

func ToServerInfoDTO(entity liquidity_provider.ServerInfo) ServerInfoDTO {
	return ServerInfoDTO{
		Version:  entity.Version,
		Revision: entity.Revision,
	}
}

func ToTrustedAccountDTO(entity liquidity_provider.TrustedAccountDetails) TrustedAccountDTO {
	return TrustedAccountDTO{
		Address:        entity.Address,
		Name:           entity.Name,
		BtcLockingCap:  entity.BtcLockingCap.AsBigInt(),
		RbtcLockingCap: entity.RbtcLockingCap.AsBigInt(),
	}
}

func ToTrustedAccountsDTO(signedEntities []entities.Signed[liquidity_provider.TrustedAccountDetails]) []TrustedAccountDTO {
	result := make([]TrustedAccountDTO, len(signedEntities))
	for i, signedEntity := range signedEntities {
		result[i] = ToTrustedAccountDTO(signedEntity.Value)
	}
	return result
}

func FromGeneralConfigurationDTO(dto GeneralConfigurationDTO) (liquidity_provider.GeneralConfiguration, error) {
	bigInt := new(big.Int)
	for key := range dto.RskConfirmations {
		_, ok := bigInt.SetString(key, 10)
		if !ok {
			return liquidity_provider.GeneralConfiguration{}, fmt.Errorf("cannot deserialize RSK confirmations key %s", key)
		}
	}
	for key := range dto.BtcConfirmations {
		_, ok := bigInt.SetString(key, 10)
		if !ok {
			return liquidity_provider.GeneralConfiguration{}, fmt.Errorf("cannot deserialize BTC confirmations key %s", key)
		}
	}
	return liquidity_provider.GeneralConfiguration{
		RskConfirmations:     dto.RskConfirmations,
		BtcConfirmations:     dto.BtcConfirmations,
		PublicLiquidityCheck: dto.PublicLiquidityCheck,
	}, nil
}
