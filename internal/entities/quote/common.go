package quote

type AcceptedQuote struct {
	Signature      string `json:"signature"`
	DepositAddress string `json:"depositAddress"`
}
