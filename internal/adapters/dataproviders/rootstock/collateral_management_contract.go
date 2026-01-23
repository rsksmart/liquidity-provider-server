package rootstock

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/collateral_management"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strings"
	"time"
)

type collateralManagementContractImpl struct {
	client          RpcClientBinding
	providerAddress string
	contractAddress string
	contract        *bind.BoundContract
	signer          TransactionSigner
	abis            *FlyoverABIs
	binding         *bindings.CollateralManagementContract
	retryParams     RetryParams
	miningTimeout   time.Duration
}

func NewCollateralManagementContractImpl(
	client *RskClient,
	providerAddress string,
	contractAddress string,
	contract *bind.BoundContract,
	signer TransactionSigner,
	binding *bindings.CollateralManagementContract,
	retryParams RetryParams,
	miningTimeout time.Duration,
	abis *FlyoverABIs,
) blockchain.CollateralManagementContract {
	return &collateralManagementContractImpl{
		client:          client.client,
		contractAddress: contractAddress,
		providerAddress: providerAddress,
		contract:        contract,
		signer:          signer,
		binding:         binding,
		retryParams:     retryParams,
		miningTimeout:   miningTimeout,
		abis:            abis,
	}
}

func (collateral *collateralManagementContractImpl) GetAddress() string {
	return collateral.contractAddress
}

func (collateral *collateralManagementContractImpl) ProviderResign() error {
	opts := &bind.TransactOpts{
		From:   collateral.signer.Address(),
		Signer: collateral.signer.Sign,
	}
	receipt, err := rskRetry(collateral.retryParams.Retries, collateral.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(collateral.client, collateral.miningTimeout, "Resign", func() (*geth.Transaction, error) {
				callData, dataErr := collateral.binding.TryPackResign()
				if dataErr != nil {
					return nil, dataErr
				}
				return bind.Transact(collateral.contract, opts, callData)
			})
		})

	if err != nil {
		return err
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("resign transaction failed")
	}
	return nil
}

func (collateral *collateralManagementContractImpl) GetCollateral(address string) (*entities.Wei, error) {
	var parsedAddress common.Address
	var err error
	opts := &bind.CallOpts{}
	if err = ParseAddress(&parsedAddress, address); err != nil {
		return nil, err
	}
	result, err := rskRetry(collateral.retryParams.Retries, collateral.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := collateral.binding.TryPackGetPegInCollateral(parsedAddress)
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(collateral.contract, opts, callData, collateral.binding.UnpackGetPegInCollateral)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(result), nil
}

func (collateral *collateralManagementContractImpl) GetPegoutCollateral(address string) (*entities.Wei, error) {
	var parsedAddress common.Address
	var err error
	opts := &bind.CallOpts{}
	if err = ParseAddress(&parsedAddress, address); err != nil {
		return nil, err
	}
	result, err := rskRetry(collateral.retryParams.Retries, collateral.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := collateral.binding.TryPackGetPegOutCollateral(parsedAddress)
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(collateral.contract, opts, callData, collateral.binding.UnpackGetPegOutCollateral)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(result), nil
}

func (collateral *collateralManagementContractImpl) GetMinimumCollateral() (*entities.Wei, error) {
	var err error
	opts := &bind.CallOpts{}
	result, err := rskRetry(collateral.retryParams.Retries, collateral.retryParams.Sleep,
		func() (*big.Int, error) {
			callData, dataErr := collateral.binding.TryPackGetMinCollateral()
			if dataErr != nil {
				return nil, dataErr
			}
			return bind.Call(collateral.contract, opts, callData, collateral.binding.UnpackGetMinCollateral)
		})
	if err != nil {
		return nil, err
	}
	return entities.NewBigWei(result), nil
}

func (collateral *collateralManagementContractImpl) AddCollateral(amount *entities.Wei) error {
	opts := &bind.TransactOpts{
		From:   collateral.signer.Address(),
		Signer: collateral.signer.Sign,
		Value:  amount.AsBigInt(),
	}

	receipt, err := rskRetry(collateral.retryParams.Retries, collateral.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(collateral.client, collateral.miningTimeout, "AddCollateral", func() (*geth.Transaction, error) {
				callData, dataErr := collateral.binding.TryPackAddPegInCollateral()
				if dataErr != nil {
					return nil, dataErr
				}
				return bind.Transact(collateral.contract, opts, callData)
			})
		})

	if err != nil {
		return fmt.Errorf("error adding collateral: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("error adding pegin collateral")
	}
	return nil
}

func (collateral *collateralManagementContractImpl) AddPegoutCollateral(amount *entities.Wei) error {
	opts := &bind.TransactOpts{
		From:   collateral.signer.Address(),
		Signer: collateral.signer.Sign,
		Value:  amount.AsBigInt(),
	}

	receipt, err := rskRetry(collateral.retryParams.Retries, collateral.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(collateral.client, collateral.miningTimeout, "AddPegoutCollateral", func() (*geth.Transaction, error) {
				callData, dataErr := collateral.binding.TryPackAddPegOutCollateral()
				if dataErr != nil {
					return nil, dataErr
				}
				return bind.Transact(collateral.contract, opts, callData)
			})
		})

	if err != nil {
		return fmt.Errorf("error adding collateral: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("error adding pegout collateral")
	}
	return nil
}

func (collateral *collateralManagementContractImpl) WithdrawCollateral() error {
	callData, dataErr := collateral.binding.TryPackWithdrawCollateral()
	if dataErr != nil {
		return dataErr
	}
	if err := collateral.validateWithdrawCollateral(callData); err != nil {
		return err
	}

	opts := &bind.TransactOpts{
		From:   collateral.signer.Address(),
		Signer: collateral.signer.Sign,
	}

	receipt, err := rskRetry(collateral.retryParams.Retries, collateral.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(collateral.client, collateral.miningTimeout, "WithdrawCollateral", func() (*geth.Transaction, error) {
				return bind.Transact(collateral.contract, opts, callData)
			})
		})

	if err != nil {
		return fmt.Errorf("withdraw collateral error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("withdraw collateral error")
	}
	return nil
}

func (collateral *collateralManagementContractImpl) validateWithdrawCollateral(callData []byte) error {
	const (
		notResignedError               = "NotResigned"
		resignationDelayNotPassedError = "ResignationDelayNotMet"
	)
	_, revert := collateral.contract.CallRaw(&bind.CallOpts{From: collateral.signer.Address()}, callData)
	parsedRevert, err := ParseRevertReason(collateral.abis.CollateralManagement, revert)
	if err != nil && parsedRevert == nil {
		return fmt.Errorf("error parsing withdrawCollateral result: %w", err)
	} else if parsedRevert != nil && (strings.EqualFold(notResignedError, parsedRevert.Name) || strings.EqualFold(resignationDelayNotPassedError, parsedRevert.Name)) {
		return liquidity_provider.ProviderNotResignedError
	} else if parsedRevert != nil {
		return fmt.Errorf("withdrawCollateral reverted with: %s", parsedRevert.Name)
	}
	return nil
}

func (collateral *collateralManagementContractImpl) GetPenalizedEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]penalization.PenalizedEvent, error) {
	var lbcEvent *bindings.CollateralManagementContractPenalized
	result := make([]penalization.PenalizedEvent, 0)

	var liquidityProviderRule []any
	var punisherRule []any
	var quoteHashRule []any
	liquidityProviderRule = append(liquidityProviderRule, common.HexToAddress(collateral.providerAddress))

	iterator, err := bind.FilterEvents(
		collateral.contract,
		&bind.FilterOpts{
			Start:   fromBlock,
			End:     toBlock,
			Context: ctx,
		},
		collateral.binding.UnpackPenalizedEvent,
		liquidityProviderRule, punisherRule, quoteHashRule,
	)

	defer func() {
		if iterator != nil {
			if iteratorError := iterator.Close(); iteratorError != nil {
				log.Error("Error closing Penalized event iterator: ", err)
			}
		}
	}()
	if err != nil || iterator == nil {
		return nil, err
	}

	for iterator.Next() {
		lbcEvent = iterator.Value()
		result = append(result, penalization.PenalizedEvent{
			LiquidityProvider: lbcEvent.LiquidityProvider.String(),
			Penalty:           entities.NewBigWei(lbcEvent.Penalty),
			QuoteHash:         hex.EncodeToString(lbcEvent.QuoteHash[:]),
		})
	}
	if err = iterator.Error(); err != nil {
		return nil, err
	}
	return result, nil
}
