package pegout

type Quote struct {
	LBCAddr               string `json:"lbcAddress" db:"lbc_addr"`
	LPRSKAddr             string `json:"liquidityProviderRskAddress" db:"lp_rsk_addr"`
	RSKRefundAddr         string `json:"rskRefundAddress" db:"rsk_refund_addr"`
	Fee                   uint64 `json:"fee" db:"fee"`
	PenaltyFee            uint64 `json:"penaltyFee" db:"penalty_fee"`
	Nonce                 int64  `json:"nonce" db:"nonce"`
	Value                 uint64 `json:"value" db:"value"`
	AgreementTimestamp    uint32 `json:"agreementTimestamp" db:"agreement_timestamp"`
	DepositDateLimit      uint32 `json:"depositDateLimit" db:"deposit_date_limit"`
	DepositConfirmations  uint16 `json:"depositConfirmations" db:"deposit_confirmations"`
	TransferConfirmations uint16 `json:"transferConfirmations" db:"transfer_confirmations"`
	TransferTime          uint32 `json:"transferTime" db:"transfer_time"`
	ExpireDate            uint32 `json:"expireDate" db:"expire_date"`
	ExpireBlocks          uint32 `json:"expireBlocks" db:"expire_blocks"`
}
