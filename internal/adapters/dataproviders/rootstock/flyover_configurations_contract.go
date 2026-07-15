package rootstock

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	configurations "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/configurations"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
)

// flyoverConfigurationsImpl is the adapter for the new FlyoverConfigurations contract (EPIC E1).
// Commit-first peg-in fees and required confirmations are read from on-chain configuration instead
// of being negotiated in a quote.
type flyoverConfigurationsImpl struct {
	client      RpcClientBinding
	address     string
	contract    *bind.BoundContract
	retryParams RetryParams
	binding     *configurations.FlyoverConfigurations
}

func NewFlyoverConfigurationsImpl(
	client *RskClient,
	address string,
	contract *bind.BoundContract,
	retryParams RetryParams,
	binding *configurations.FlyoverConfigurations,
) blockchain.FlyoverConfigurationsContract {
	return &flyoverConfigurationsImpl{
		client:      client.client,
		address:     address,
		contract:    contract,
		retryParams: retryParams,
		binding:     binding,
	}
}

func (c *flyoverConfigurationsImpl) GetAddress() string {
	return c.address
}

func (c *flyoverConfigurationsImpl) CalculatePegInFee(amount *entities.Wei) (*entities.Wei, error) {
	fee, err := rskRetry(c.retryParams.Retries, c.retryParams.Sleep, func() (*big.Int, error) {
		callData, dataErr := c.binding.TryPackCalculatePegInFee(amount.AsBigInt())
		if dataErr != nil {
			return nil, dataErr
		}
		return bind.Call(c.contract, &bind.CallOpts{}, callData, c.binding.UnpackCalculatePegInFee)
	})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(fee), nil
}

func (c *flyoverConfigurationsImpl) GetRequiredPegInConfirmations(amount *entities.Wei) (uint64, error) {
	confirmations, err := rskRetry(c.retryParams.Retries, c.retryParams.Sleep, func() (*big.Int, error) {
		callData, dataErr := c.binding.TryPackGetRequiredPegInConfirmations(amount.AsBigInt())
		if dataErr != nil {
			return nil, dataErr
		}
		return bind.Call(c.contract, &bind.CallOpts{}, callData, c.binding.UnpackGetRequiredPegInConfirmations)
	})
	if err != nil {
		return 0, err
	}
	return confirmations.Uint64(), nil
}
