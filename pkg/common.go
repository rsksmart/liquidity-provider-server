package pkg

import (
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"math/big"
)

type AcceptQuoteRequest struct {
	QuoteHash string `json:"quoteHash" required:"" validate:"required" example:"0x0" description:"QuoteHash"`
}

type AcceptAuthenticatedQuoteRequest struct {
	QuoteHash string `json:"quoteHash" required:"" validate:"required" example:"0x0" description:"QuoteHash"`
	Signature string `json:"signature" required:"" validate:"required" example:"0x0" description:"Signature from a trusted account"`
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

type RecommendedOperationDTO struct {
	RecommendedQuoteValue *big.Int `json:"recommendedQuoteValue" example:"100000" description:"Recommended quote value for the input amount" required:""`
	EstimatedCallFee      *big.Int `json:"estimatedCallFee"  example:"100000" description:"Estimated call fee if a quote is created with the recommended amount" required:""`
	EstimatedGasFee       *big.Int `json:"estimatedGasFee"  example:"100000" description:"Estimated gas fee if a quote is created with the recommended amount" required:""`
	EstimatedProductFee   *big.Int `json:"estimatedProductFee"  example:"100000" description:"Estimated product fee if a quote is created with the recommended amount" required:""`
}

func ToRecommendedOperationDTO(domain usecases.RecommendedOperationResult) RecommendedOperationDTO {
	return RecommendedOperationDTO{
		RecommendedQuoteValue: domain.RecommendedQuoteValue.AsBigInt(),
		EstimatedCallFee:      domain.EstimatedCallFee.AsBigInt(),
		EstimatedGasFee:       domain.EstimatedGasFee.AsBigInt(),
		EstimatedProductFee:   domain.EstimatedProductFee.AsBigInt(),
	}
}
