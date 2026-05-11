package pegout

import (
	"context"
	"errors"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout/common"
	log "github.com/sirupsen/logrus"
)

type AllAtOnceHandler struct {
	quoteRepository quote.PegoutQuoteRepository
	rskWallet       blockchain.RootstockWallet
	contracts       blockchain.RskContracts
	mutex           sync.Locker
}

func NewAllAtOnceHandler(
	quoteRepository quote.PegoutQuoteRepository,
	rskWallet blockchain.RootstockWallet,
	contracts blockchain.RskContracts,
	mutex sync.Locker,
) *AllAtOnceHandler {
	return &AllAtOnceHandler{
		quoteRepository: quoteRepository,
		rskWallet:       rskWallet,
		contracts:       contracts,
		mutex:           mutex,
	}
}

func (h *AllAtOnceHandler) Execute(
	ctx context.Context,
	pegoutConfig liquidity_provider.PegoutConfiguration,
	watchedQuotes []quote.WatchedPegoutQuote,
) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	totalValue := common.CalculateTotalToPegout(watchedQuotes)
	if pegoutConfig.BridgeTransactionMin.Cmp(totalValue) > 0 {
		log.Infof(
			"Refunded pegouts total value: %v out of %v. Threshold not met yet. Skipping transaction to the bridge.",
			totalValue.ToRbtc(),
			pegoutConfig.BridgeTransactionMin.ToRbtc(),
		)
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.TxBelowMinimumError)
	}

	requiredBalance := new(entities.Wei).Add(totalValue, entities.NewWei(BridgeConversionGasLimit*BridgeConversionGasPrice))
	if err := common.CheckBalance(ctx, usecases.BridgePegoutId, h.rskWallet, requiredBalance); err != nil {
		return err
	}

	config := blockchain.NewTransactionConfig(totalValue, BridgeConversionGasLimit, entities.NewWei(BridgeConversionGasPrice))
	receipt, txErr := h.rskWallet.SendRbtc(ctx, config, h.contracts.Bridge.GetAddress())
	if txErr == nil {
		log.Debugf("%s: transaction sent to the bridge successfully (%s)", usecases.BridgePegoutId, receipt.TransactionHash)
	}

	if err := h.updateQuotes(ctx, receipt, txErr, watchedQuotes); err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}
	return nil
}

func (h *AllAtOnceHandler) updateQuotes(
	ctx context.Context,
	receipt blockchain.TransactionReceipt,
	txErr error,
	watchedQuotes []quote.WatchedPegoutQuote,
) error {
	var err error
	retainedQuotes := make([]quote.RetainedPegoutQuote, 0)
	err = errors.Join(err, txErr)
	for _, watchedQuote := range watchedQuotes {
		if receipt.TransactionHash != "" && receipt.GasUsed != nil {
			watchedQuote.RetainedQuote.AppendRebalanceAllocation(receipt.TransactionHash, receipt.GasUsed.Uint64(), receipt.GasPrice)
			watchedQuote.RetainedQuote.SetDeprecatedRefundFields()
		}
		if txErr == nil {
			watchedQuote.RetainedQuote.RemainingToRefund = entities.NewWei(0)
			watchedQuote.RetainedQuote.State = quote.PegoutStateBridgeTxSucceeded
		} else {
			watchedQuote.RetainedQuote.State = quote.PegoutStateBridgeTxFailed
		}
		retainedQuotes = append(retainedQuotes, watchedQuote.RetainedQuote)
	}
	if updateErr := h.quoteRepository.UpdateRetainedQuotes(ctx, retainedQuotes); updateErr != nil {
		err = errors.Join(err, updateErr)
	}
	return err
}
