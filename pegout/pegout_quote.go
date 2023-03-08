package pegout

type Quote struct {
	LBCAddr               string `json:"lbcAddress" db:"lbc_addr" validate:"required"`
	LPRSKAddr             string `json:"liquidityProviderRskAddress" db:"lp_rsk_addr" validate:"required"`
	RSKRefundAddr         string `json:"rskRefundAddress" db:"rsk_refund_addr" validate:"required"`
	CallFee               uint64 `json:"callFee" db:"callFee" validate:"required"`
	PenaltyFee            uint64 `json:"penaltyFee" db:"penalty_fee" validate:"required"`
	Nonce                 int64  `json:"nonce" db:"nonce" validate:"required"`
	Value                 uint64 `json:"value" db:"value" validate:"required"`
	AgreementTimestamp    uint32 `json:"agreementTimestamp" db:"agreement_timestamp" validate:"required"`
	DepositDateLimit      uint32 `json:"depositDateLimit" db:"deposit_date_limit" validate:"required"`
	DepositConfirmations  uint16 `json:"depositConfirmations" db:"deposit_confirmations" validate:"required"`
	TransferConfirmations uint16 `json:"transferConfirmations" db:"transfer_confirmations" validate:"required"`
	TransferTime          uint32 `json:"transferTime" db:"transfer_time" validate:"required"`
	ExpireDate            uint32 `json:"expireDate" db:"expire_date" validate:"required"`
	ExpireBlocks          uint32 `json:"expireBlocks" db:"expire_blocks" validate:"required"`
}
