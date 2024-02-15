package quote

import "github.com/rsksmart/liquidity-provider-server/internal/entities"

type AcceptedQuote struct {
	Signature      string `json:"signature"`
	DepositAddress string `json:"depositAddress"`
}

type Fees struct {
	CallFee          *entities.Wei
	GasFee           *entities.Wei
	PenaltyFee       *entities.Wei
	ProductFeeAmount uint64
}
