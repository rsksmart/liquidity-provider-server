package pegout

import "github.com/rsksmart/liquidity-provider/types"

type RetainedQuote struct {
	QuoteHash          string        `json:"quoteHash" db:"quote_hash"`
	DepositAddr        string        `json:"depositAddr" db:"deposit_addr"`
	Signature          string        `json:"signature" db:"signature"`
	ReqLiq             uint64        `json:"reqLiq" db:"req_liq"`
	State              types.RQState `json:"state" db:"state"`
	DepositTransaction string        `json:"depositTransaction" db:"deposit_transaction" bson:"deposit_transaction"`
}
