package rootstock

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"math/big"
)

type feeCollectorImpl struct {
	contract *bindings.LiquidityBridgeContract
}

func NewFeeCollectorImpl(contract *bindings.LiquidityBridgeContract) blockchain.FeeCollector {
	return &feeCollectorImpl{contract: contract}
}

func (fc *feeCollectorImpl) DaoFeePercentage() (uint64, error) {
	opts := bind.CallOpts{}
	amount, err := rskRetry(func() (*big.Int, error) {
		return fc.contract.ProductFeePercentage(&opts)
	})
	if err != nil {
		return 0, err
	}
	return amount.Uint64(), nil
}
