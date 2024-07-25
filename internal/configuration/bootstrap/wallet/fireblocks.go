package wallet

import "github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"

// TODO complete with fireblocks integration

const featureUnimplemented = "feature unimplemented"

type FireBlocksWalletFactory struct{}

func NewFireBlocksFactory(args FactoryCreationArgs) (AbstractFactory, error) {
	panic(featureUnimplemented)
}

func (f FireBlocksWalletFactory) BitcoinMonitoringWallet(walletId string) (blockchain.BitcoinWallet, error) {
	panic(featureUnimplemented)
}

func (f FireBlocksWalletFactory) BitcoinPaymentWallet(walletId string) (blockchain.BitcoinWallet, error) {
	panic(featureUnimplemented)
}

func (f FireBlocksWalletFactory) RskWallet() (blockchain.RootstockWallet, error) {
	panic(featureUnimplemented)
}
