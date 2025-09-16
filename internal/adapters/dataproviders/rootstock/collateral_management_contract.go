package rootstock

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	geth "github.com/ethereum/go-ethereum/core/types"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings"
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
	contract        CollateralManagementAdapter
	signer          TransactionSigner
	abi             *abi.ABI
	retryParams     RetryParams
	miningTimeout   time.Duration
}

func NewCollateralManagementContractImpl(
	client *RskClient,
	providerAddress string,
	contractAddress string,
	contract CollateralManagementAdapter,
	signer TransactionSigner,
	retryParams RetryParams,
	miningTimeout time.Duration,
) blockchain.CollateralManagementContract {
	contractAbi, err := bindings.IPegOutMetaData.GetAbi()
	if err != nil {
		panic(fmt.Sprintf("could not get ABI for Collateral Management contract: %v", err))
	}
	return &collateralManagementContractImpl{
		client:          client.client,
		contractAddress: contractAddress,
		providerAddress: providerAddress,
		contract:        contract,
		signer:          signer,
		retryParams:     retryParams,
		miningTimeout:   miningTimeout,
		abi:             contractAbi,
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
				return collateral.contract.Resign(opts)
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
			return collateral.contract.GetPegInCollateral(opts, parsedAddress)
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
			return collateral.contract.GetPegOutCollateral(opts, parsedAddress)
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
			return collateral.contract.GetMinCollateral(opts)
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
				return collateral.contract.AddPegInCollateral(opts)
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
				return collateral.contract.AddPegOutCollateral(opts)
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
	const (
		functionName                   = "withdrawCollateral"
		notResignedError               = "NotResigned"
		resignationDelayNotPassedError = "ResignationDelayNotMet"
	)
	var res []any
	var err error

	revert := collateral.contract.Caller().Call(&bind.CallOpts{From: collateral.signer.Address()}, &res, functionName)
	parsedRevert, err := ParseRevertReason(collateral.abi, revert)
	if err != nil && parsedRevert == nil {
		return fmt.Errorf("error parsing withdrawCollateral result: %w", err)
	} else if parsedRevert != nil && (strings.EqualFold(notResignedError, parsedRevert.Name) || strings.EqualFold(resignationDelayNotPassedError, parsedRevert.Name)) {
		return liquidity_provider.ProviderNotResignedError
	} else if parsedRevert != nil {
		return fmt.Errorf("withdrawCollateral reverted with: %s", parsedRevert.Name)
	}

	opts := &bind.TransactOpts{
		From:   collateral.signer.Address(),
		Signer: collateral.signer.Sign,
	}

	receipt, err := rskRetry(collateral.retryParams.Retries, collateral.retryParams.Sleep,
		func() (*geth.Receipt, error) {
			return awaitTx(collateral.client, collateral.miningTimeout, "WithdrawCollateral", func() (*geth.Transaction, error) {
				return collateral.contract.WithdrawCollateral(opts)
			})
		})

	if err != nil {
		return fmt.Errorf("withdraw collateral error: %w", err)
	} else if receipt == nil || receipt.Status == 0 {
		return errors.New("withdraw collateral error")
	}
	return nil
}

func (collateral *collateralManagementContractImpl) GetPenalizedEvents(ctx context.Context, fromBlock uint64, toBlock *uint64) ([]penalization.PenalizedEvent, error) {
	var lbcEvent *bindings.ICollateralManagementPenalized
	result := make([]penalization.PenalizedEvent, 0)

	rawIterator, err := collateral.contract.FilterPenalized(&bind.FilterOpts{
		Start:   fromBlock,
		End:     toBlock,
		Context: ctx,
	}, []common.Address{common.HexToAddress(collateral.providerAddress)}, nil, nil)
	iterator := collateral.contract.PenalizedEventIteratorAdapter(rawIterator)
	defer func() {
		if rawIterator != nil {
			if iteratorError := iterator.Close(); iteratorError != nil {
				log.Error("Error closing Penalized event iterator: ", err)
			}
		}
	}()
	if err != nil || rawIterator == nil {
		return nil, err
	}

	for iterator.Next() {
		lbcEvent = iterator.Event()
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
