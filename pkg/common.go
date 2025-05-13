package pkg

type AcceptQuoteRequest struct {
	QuoteHash string `json:"quoteHash" required:"" validate:"required" example:"0x0" description:"QuoteHash"`
}

type AcceptQuoteFromTrustedAccountRequest struct {
	QuoteHash string `json:"quoteHash" required:"" validate:"required" example:"0x0" description:"QuoteHash"`
	Signature string `json:"signature" required:"" validate:"required" example:"0x0" description:"Signature from a trusted account"`
}

type GetCollateralResponse struct {
	Collateral uint64 `json:"collateral" required:""`
}

type AddCollateralRequest struct {
	Amount uint64 `json:"amount"  required:"" validate:"required" example:"100000000000" description:"Amount to add to the collateral"`
}

type AddCollateralResponse struct {
	NewCollateralBalance uint64 `json:"newCollateralBalance" example:"100000000000" description:"New Collateral Balance"`
}

type HealthResponse struct {
	Status   string   `json:"status" example:"ok" description:"Overall LPS Health Status" required:""`
	Services Services `json:"services" example:"{\"db\":\"ok\",\"rsk\":\"ok\",\"btc\":\"ok\"}" description:"LPS Services Status" required:""`
}

type Services struct {
	Db  string `json:"db"`
	Rsk string `json:"rsk"`
	Btc string `json:"btc"`
}
