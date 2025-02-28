package pkg

import (
	"math/big"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
)

type ProviderDetail struct {
	FixedFee              uint64     `json:"fixedFee" required:""`
	PercentageFee         *big.Float `json:"percentageFee" required:""`
	MinTransactionValue   uint64     `json:"minTransactionValue"  required:""`
	MaxTransactionValue   uint64     `json:"maxTransactionValue"  required:""`
	RequiredConfirmations uint16     `json:"requiredConfirmations"  required:""`
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
	TimeForDeposit uint32 `json:"timeForDeposit" validate:"required"`
	CallTime       uint32 `json:"callTime" validate:"required"`
	PenaltyFee     string `json:"penaltyFee" validate:"required,numeric,positive_string"`
	PercentageFee  string `json:"percentageFee" validate:"required,percentage_fee"`
	FixedFee       string `json:"fixedFee" validate:"required,numeric,non_negative_string"`
	MaxValue       string `json:"maxValue" validate:"required,numeric,positive_string"`
	MinValue       string `json:"minValue" validate:"required,numeric,positive_string"`
}

type PegoutConfigurationRequest struct {
	Configuration PegoutConfigurationDTO `json:"configuration" validate:"required"`
}

type PegoutConfigurationDTO struct {
	TimeForDeposit       uint32 `json:"timeForDeposit" validate:"required"`
	ExpireTime           uint32 `json:"expireTime" validate:"required"`
	PenaltyFee           string `json:"penaltyFee" validate:"required,numeric,positive_string"`
	PercentageFee        string `json:"percentageFee" validate:"required,percentage_fee"`
	FixedFee             string `json:"fixedFee" validate:"required,numeric,positive_string"`
	MaxValue             string `json:"maxValue" validate:"required,numeric,positive_string"`
	MinValue             string `json:"minValue" validate:"required,numeric,positive_string"`
	ExpireBlocks         uint64 `json:"expireBlocks" validate:"required"`
	BridgeTransactionMin string `json:"bridgeTransactionMin" validate:"required,numeric,positive_string"`
}

type GeneralConfigurationRequest struct {
	Configuration *liquidity_provider.GeneralConfiguration `json:"configuration" validate:"required"`
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

type AvailableLiquidityDTO struct {
	PeginLiquidityAmount  *big.Int `json:"peginLiquidityAmount" example:"5000000000000000000" description:"Available liquidity for PegIn operations in wei"  required:""`
	PegoutLiquidityAmount *big.Int `json:"pegoutLiquidityAmount" example:"5000000000000000000" description:"Available liquidity for PegOut operations in wei" required:""`
}

type ServerInfoDTO struct {
	Version  string `json:"version" example:"v1.0.0" description:"Server version tag"  required:""`
	Revision string `json:"revision" example:"b7bf393a2b1cedde8ee15b00780f44e6e5d2ba9d" description:"Version commit hash"  required:""`
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

	maxVal := new(big.Int)
	maxVal.SetString(dto.MaxValue, base)

	minVal := new(big.Int)
	minVal.SetString(dto.MinValue, base)

	percentageFee, _ := new(big.Float).SetString(dto.PercentageFee)

	return liquidity_provider.PeginConfiguration{
		TimeForDeposit: dto.TimeForDeposit,
		CallTime:       dto.CallTime,
		PenaltyFee:     entities.NewBigWei(penaltyFee),
		PercentageFee:  percentageFee,
		FixedFee:       entities.NewBigWei(fixedFee),
		MaxValue:       entities.NewBigWei(maxVal),
		MinValue:       entities.NewBigWei(minVal),
	}
}

func FromPegoutConfigurationDTO(dto PegoutConfigurationDTO) liquidity_provider.PegoutConfiguration {
	const base = 10
	penaltyFee := new(big.Int)
	penaltyFee.SetString(dto.PenaltyFee, base)
	maxValue := new(big.Int)
	maxValue.SetString(dto.MaxValue, base)
	minValue := new(big.Int)
	minValue.SetString(dto.MinValue, base)
	bridgeTransactionMin := new(big.Int)
	bridgeTransactionMin.SetString(dto.BridgeTransactionMin, base)

	fixedFee := new(big.Int)
	fixedFee.SetString(dto.FixedFee, base)
	percentageFee, _ := new(big.Float).SetString(dto.PercentageFee)

	return liquidity_provider.PegoutConfiguration{
		TimeForDeposit:       dto.TimeForDeposit,
		ExpireTime:           dto.ExpireTime,
		PenaltyFee:           entities.NewBigWei(penaltyFee),
		PercentageFee:        percentageFee,
		FixedFee:             entities.NewBigWei(fixedFee),
		MaxValue:             entities.NewBigWei(maxValue),
		MinValue:             entities.NewBigWei(minValue),
		ExpireBlocks:         dto.ExpireBlocks,
		BridgeTransactionMin: entities.NewBigWei(bridgeTransactionMin),
	}
}

func ToServerInfoDTO(entity liquidity_provider.ServerInfo) ServerInfoDTO {
	return ServerInfoDTO{
		Version:  entity.Version,
		Revision: entity.Revision,
	}
}
