package cold_wallet

type ColdWalletType string

const (
	StaticColdWalletType ColdWalletType = "StaticColdWallet"
)

type ColdWallet interface {
	Init() error
	GetBtcAddress() string
	GetRskAddress() string
	GetLabel() string
}
