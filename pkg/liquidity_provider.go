package pkg

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"math/big"
)

type ProviderDetail struct {
	Fee                   uint64 `json:"fee"  required:""`
	MinTransactionValue   uint64 `json:"minTransactionValue"  required:""`
	MaxTransactionValue   uint64 `json:"maxTransactionValue"  required:""`
	RequiredConfirmations uint16 `json:"requiredConfirmations"  required:""`
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
	CallFee        string `json:"callFee" validate:"required,numeric,positive_string"`
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
	CallFee              string `json:"callFee" validate:"required,numeric,positive_string"`
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
	PegoutLiquidityAmount *big.Int `json:"pegoutLiquidityAmount" example:"500000000" description:"Available liquidity for PegOut operations in satoshi" required:""`
}

type ServerInfoDTO struct {
	Version  string `json:"version" example:"v1.0.0" description:"Server version tag"  required:""`
	Revision string `json:"revision" example:"b7bf393a2b1cedde8ee15b00780f44e6e5d2ba9d" description:"Version commit hash"  required:""`
}

func ToAvailableLiquidityDTO(entity liquidity_provider.AvailableLiquidity) AvailableLiquidityDTO {
	satoshis, _ := entity.PegoutLiquidity.ToSatoshi().Int(nil)
	return AvailableLiquidityDTO{
		PeginLiquidityAmount:  entity.PeginLiquidity.AsBigInt(),
		PegoutLiquidityAmount: satoshis,
	}
}

func FromPeginConfigurationDTO(dto PeginConfigurationDTO) liquidity_provider.PeginConfiguration {
	const base = 10
	penaltyFee := new(big.Int)
	penaltyFee.SetString(dto.PenaltyFee, base)
	callFee := new(big.Int)
	callFee.SetString(dto.CallFee, base)
	maxValue := new(big.Int)
	maxValue.SetString(dto.MaxValue, base)
	minValue := new(big.Int)
	minValue.SetString(dto.MinValue, base)

	return liquidity_provider.PeginConfiguration{
		TimeForDeposit: dto.TimeForDeposit,
		CallTime:       dto.CallTime,
		PenaltyFee:     entities.NewBigWei(penaltyFee),
		CallFee:        entities.NewBigWei(callFee),
		MaxValue:       entities.NewBigWei(maxValue),
		MinValue:       entities.NewBigWei(minValue),
	}
}

func FromPegoutConfigurationDTO(dto PegoutConfigurationDTO) liquidity_provider.PegoutConfiguration {
	const base = 10
	penaltyFee := new(big.Int)
	penaltyFee.SetString(dto.PenaltyFee, base)
	callFee := new(big.Int)
	callFee.SetString(dto.CallFee, base)
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
		CallFee:              entities.NewBigWei(callFee),
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
