package rootstock

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"math/big"
)

type feeCollectorImpl struct {
	retryParams RetryParams
	contract    LbcBinding
}

func NewFeeCollectorImpl(contract LbcBinding, retryParams RetryParams) blockchain.FeeCollector {
	return &feeCollectorImpl{contract: contract, retryParams: retryParams}
}

func (fc *feeCollectorImpl) DaoFeePercentage() (uint64, error) {
	opts := bind.CallOpts{}
	amount, err := rskRetry(fc.retryParams.Retries, fc.retryParams.Sleep,
		func() (*big.Int, error) {
			return fc.contract.ProductFeePercentage(&opts)
		})
	if err != nil {
		return 0, err
	}
	return amount.Uint64(), nil
}
