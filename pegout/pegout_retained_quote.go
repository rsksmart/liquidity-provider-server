package pegout

import "github.com/rsksmart/liquidity-provider/types"

type RetainedQuote struct {
	QuoteHash          string        `json:"quoteHash" db:"quote_hash"`
	DepositAddr        string        `json:"depositAddr" db:"deposit_addr"`
	Signature          string        `json:"signature" db:"signature"`
	ReqLiq             uint64        `json:"reqLiq" db:"req_liq"`
	State              types.RQState `json:"state" db:"state"`
	DepositBlockNumber uint64        `json:"depositBlockNumber" db:"deposit_block_number" bson:"deposit_block_number"`
}
