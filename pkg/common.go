package pkg

import "math/big"

type AcceptQuoteRequest struct {
	QuoteHash string `json:"quoteHash" required:"" validate:"required" example:"0x0" description:"QuoteHash"`
}

type GetCollateralResponse struct {
	Collateral *big.Int `json:"collateral" required:""`
}

type AddCollateralRequest struct {
	Amount *big.Int `json:"amount"  required:"" validate:"required" example:"100000000000" description:"Amount to add to the collateral"`
}

type AddCollateralResponse struct {
	NewCollateralBalance *big.Int `json:"newCollateralBalance" example:"100000000000" description:"New Collateral Balance"`
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
