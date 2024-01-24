package pkg

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
