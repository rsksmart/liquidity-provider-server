package rootstock

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"math/big"
	"time"
)

type peginContractImpl struct {
	client        RpcClientBinding
	address       string
	contract      *bind.BoundContract
	signer        TransactionSigner
	retryParams   RetryParams
	miningTimeout time.Duration
	binding       *bindings.PeginContract
	abis          *FlyoverABIs
}

func NewPeginContractImpl(
	client *RskClient,
	address string,
	contract *bind.BoundContract,
	signer TransactionSigner,
	retryParams RetryParams,
	miningTimeout time.Duration,
	binding *bindings.PeginContract,
	abis *FlyoverABIs,
) blockchain.PeginContract {
	return &peginContractImpl{
		client:        client.client,
		address:       address,
		contract:      contract,
		signer:        signer,
		retryParams:   retryParams,
		miningTimeout: miningTimeout,
		binding:       binding,
		abis:          abis,
	}
}

func (peginContract *peginContractImpl) GetAddress() string {
	return peginContract.address
}

func (peginContract *peginContractImpl) GetBalance(address string) (*entities.Wei, error) {
	var parsedAddress common.Address
	var err error
	if err = ParseAddress(&parsedAddress, address); err != nil {
		return nil, err
	}
	balance, err := rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := peginContract.binding.TryPackGetBalance(parsedAddress)
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(peginContract.contract, &bind.CallOpts{}, callData, peginContract.binding.UnpackGetBalance)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(balance), nil
}

// Withdraw withdraws the specified amount from the LP's balance in the pegin contract.
// It first performs a dry-run call to check for reverts before submitting the actual transaction.
func (peginContract *peginContractImpl) Withdraw(amount *entities.Wei) error {
	callData, dataErr := peginContract.binding.TryPackWithdraw(amount.AsBigInt())
	if dataErr != nil {
		return dataErr
	}
	_, revert := peginContract.contract.CallRaw(&bind.CallOpts{}, callData)
	parsedRevert, err := ParseRevertReason(peginContract.abis.Flyover, revert)
	if err != nil && parsedRevert == nil {
		return fmt.Errorf("error parsing withdraw result: %w", err)
	} else if parsedRevert != nil {
		return fmt.Errorf("withdraw reverted with: %s", parsedRevert.Name)
	}

	opts := &bind.TransactOpts{
		From:   peginContract.signer.Address(),
		Signer: peginContract.signer.Sign,
	}

	receipt, err := rskRetry(peginContract.retryParams.Retries, peginContract.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(peginContract.client, peginContract.miningTimeout, "Withdraw", func() (*geth.Transaction, error) {
				return bind.Transact(peginContract.contract, opts, callData)
			})
		})

	if err != nil {
		return fmt.Errorf("withdraw error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("withdraw error: transaction failed")
	}
	return nil
}

func (peginContract *peginContractImpl) PausedStatus() (blockchain.PauseStatus, error) {
	opts := new(bind.CallOpts)
	result, err := rskRetry(
		peginContract.retryParams.Retries,
		peginContract.retryParams.Sleep,
		func() (bindings.PauseStatusOutput, error) {
			callData, dataErr := peginContract.binding.TryPackPauseStatus()
			if dataErr != nil {
				return bindings.PauseStatusOutput{}, dataErr
			}
			return bind.Call(peginContract.contract, opts, callData, peginContract.binding.UnpackPauseStatus)
		},
	)
	if err != nil {
		return blockchain.PauseStatus{}, err
	}
	return blockchain.PauseStatus{
		IsPaused: result.IsPaused,
		Reason:   result.Reason,
		Since:    result.Since,
	}, nil
}
