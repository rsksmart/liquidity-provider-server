package cold_wallet

import "github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"

type StaticColdWalletArgs struct {
	BtcAddress string `json:"btcAddress"`
	RskAddress string `json:"rskAddress"`
}

type StaticColdWallet struct {
	args       StaticColdWalletArgs
	rpc        blockchain.Rpc
	btcAddress string
	rskAddress string
}

func NewStaticColdWallet(rpc blockchain.Rpc, args StaticColdWalletArgs) ColdWallet {
	return &StaticColdWallet{args: args, rpc: rpc}
}

func (w *StaticColdWallet) Init() error {
	if err := w.rpc.Btc.ValidateAddress(w.args.BtcAddress); err != nil {
		return err
	}
	if !blockchain.IsRskAddress(w.args.RskAddress) {
		return blockchain.InvalidAddressError
	}
	w.btcAddress = w.args.BtcAddress
	w.rskAddress = w.args.RskAddress
	return nil
}

func (w *StaticColdWallet) GetBtcAddress() string {
	return w.btcAddress
}

func (w *StaticColdWallet) GetRskAddress() string {
	return w.rskAddress
}
