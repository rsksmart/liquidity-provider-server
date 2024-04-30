package btcclient

type ReadonlyWalletRequest struct {
	WalletName         string `json:"wallet_name"`
	DisablePrivateKeys bool   `json:"disable_private_keys"`
	Blank              bool   `json:"blank"`
	AvoidReuse         bool   `json:"avoid_reuse"`
	Descriptors        bool   `json:"descriptors"`
}
