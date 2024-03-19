package pkg

import "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"

type ProviderDetail struct {
	Fee                   uint64 `json:"fee"  required:""`
	MinTransactionValue   uint64 `json:"minTransactionValue"  required:""`
	MaxTransactionValue   uint64 `json:"maxTransactionValue"  required:""`
	RequiredConfirmations uint16 `json:"requiredConfirmations"  required:""`
}

type ProviderDetailResponse struct {
	SiteKey string         `json:"siteKey" required:""`
	Pegin   ProviderDetail `json:"pegin" required:""`
	Pegout  ProviderDetail `json:"pegout" required:""`
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
	Configuration *liquidity_provider.PeginConfiguration `json:"configuration" validate:"required"`
}

type PegoutConfigurationRequest struct {
	Configuration *liquidity_provider.PegoutConfiguration `json:"configuration" validate:"required"`
}

type GeneralConfigurationRequest struct {
	Configuration *liquidity_provider.GeneralConfiguration `json:"configuration" validate:"required"`
}
