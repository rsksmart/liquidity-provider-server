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

type UtxoSplitHandler struct {
	quoteRepository quote.PegoutQuoteRepository
	rskWallet       blockchain.RootstockWallet
	contracts       blockchain.RskContracts
	mutex           sync.Locker
}

func NewUtxoSplitHandler(
	quoteRepository quote.PegoutQuoteRepository,
	rskWallet blockchain.RootstockWallet,
	contracts blockchain.RskContracts,
	mutex sync.Locker,
) *UtxoSplitHandler {
	return &UtxoSplitHandler{
		quoteRepository: quoteRepository,
		rskWallet:       rskWallet,
		contracts:       contracts,
		mutex:           mutex,
	}
}

func (h *UtxoSplitHandler) Execute(
	ctx context.Context,
	pegoutConfig liquidity_provider.PegoutConfiguration,
	watchedQuotes []quote.WatchedPegoutQuote,
) error {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	bridgeMin := pegoutConfig.BridgeTransactionMin
	adjustedTotal := h.adjustTotalForRetries(watchedQuotes)
	if bridgeMin.Cmp(adjustedTotal) > 0 {
		log.Infof(
			"Refunded pegouts total value: %v out of %v. Threshold not met yet. Skipping transaction to the bridge.",
			adjustedTotal.ToRbtc(),
			bridgeMin.ToRbtc(),
		)
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, usecases.TxBelowMinimumError)
	}
	numTxsWei, err := new(entities.Wei).Div(adjustedTotal, bridgeMin)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}
	numTxs := numTxsWei.Uint64()
	remainder := new(entities.Wei).Sub(adjustedTotal, new(entities.Wei).Mul(numTxsWei, bridgeMin))

	gasPerTx := entities.NewWei(BridgeConversionGasLimit * BridgeConversionGasPrice)
	requiredBalance := new(entities.Wei).Add(adjustedTotal, new(entities.Wei).Mul(numTxsWei, gasPerTx))
	if err := common.CheckBalance(ctx, usecases.BridgePegoutId, h.rskWallet, requiredBalance); err != nil {
		return err
	}

	chunkAmounts := h.buildChunkAmounts(numTxs, bridgeMin, remainder)
	if err := h.sendAndPersistProgress(ctx, chunkAmounts, watchedQuotes); err != nil {
		return usecases.WrapUseCaseError(usecases.BridgePegoutId, err)
	}
	return nil
}

func (h *UtxoSplitHandler) sendAndPersistProgress(
	ctx context.Context,
	chunkAmounts []*entities.Wei,
	watchedQuotes []quote.WatchedPegoutQuote,
) error {
	bridgeAddress := h.contracts.Bridge.GetAddress()

	for i, chunkAmount := range chunkAmounts {
		updatedQuotes, err := h.sendChunkAndPersist(ctx, i, len(chunkAmounts), chunkAmount, bridgeAddress, watchedQuotes)
		if err != nil {
			return err
		}
		watchedQuotes = updatedQuotes
	}
	return nil
}

func (h *UtxoSplitHandler) sendChunkAndPersist(
	ctx context.Context,
	chunkIndex, totalChunks int,
	chunkAmount *entities.Wei,
	bridgeAddress string,
	watchedQuotes []quote.WatchedPegoutQuote,
) ([]quote.WatchedPegoutQuote, error) {
	config := blockchain.NewTransactionConfig(chunkAmount, BridgeConversionGasLimit, entities.NewWei(BridgeConversionGasPrice))
	receipt, txErr := h.rskWallet.SendRbtc(ctx, config, bridgeAddress)
	if !h.isValidSplitReceipt(receipt, txErr, chunkIndex, totalChunks) {
		return watchedQuotes, nil
	}
	return h.distributeChunkReceipt(ctx, chunkIndex, totalChunks, receipt, watchedQuotes)
}

func (h *UtxoSplitHandler) isValidSplitReceipt(receipt blockchain.TransactionReceipt, txErr error, chunkIndex, totalChunks int) bool {
	if errors.Is(txErr, blockchain.TxFailedError) {
		log.Errorf("%s: split tx %d/%d failed: %v", usecases.BridgePegoutId, chunkIndex+1, totalChunks, txErr)
		return false
	}
	if receipt.TransactionHash == "" {
		log.Errorf("%s: split tx %d/%d failed: incomplete receipt", usecases.BridgePegoutId, chunkIndex+1, totalChunks)
		return false
	}
	if receipt.Value == nil {
		log.Errorf("%s: split tx %d/%d failed: missing receipt value", usecases.BridgePegoutId, chunkIndex+1, totalChunks)
		return false
	}
	if receipt.GasUsed == nil {
		log.Errorf("%s: split tx %d/%d failed: missing receipt gas used", usecases.BridgePegoutId, chunkIndex+1, totalChunks)
		return false
	}
	log.Debugf("%s: split tx %d/%d sent to the bridge successfully (%s)", usecases.BridgePegoutId, chunkIndex+1, totalChunks, receipt.TransactionHash)
	return true
}

func (h *UtxoSplitHandler) distributeChunkReceipt(
	ctx context.Context,
	chunkIndex, totalChunks int,
	receipt blockchain.TransactionReceipt,
	watchedQuotes []quote.WatchedPegoutQuote,
) ([]quote.WatchedPegoutQuote, error) {
	var err error
	chunkRemaining := receipt.Value.Copy()

	watchedQuotes, err = h.flushFullyCoveredQuotes(ctx, chunkIndex, totalChunks, watchedQuotes)
	if err != nil {
		return nil, err
	}

	for chunkRemaining.Cmp(entities.NewWei(0)) > 0 && len(watchedQuotes) > 0 {
		watchedQuote := watchedQuotes[0]
		watchedQuotes = watchedQuotes[1:]

		watchedQuote, chunkRemaining = h.allocateChunkToQuote(watchedQuote, receipt, chunkRemaining)
		if err = h.persistQuoteUpdate(ctx, chunkIndex, totalChunks, watchedQuote); err != nil {
			return nil, err
		}

		if watchedQuote.RetainedQuote.RemainingToRefund.Cmp(entities.NewWei(0)) > 0 {
			watchedQuotes = append([]quote.WatchedPegoutQuote{watchedQuote}, watchedQuotes...)
			return watchedQuotes, nil
		}
	}
	return watchedQuotes, nil
}

func (h *UtxoSplitHandler) flushFullyCoveredQuotes(
	ctx context.Context,
	chunkIndex, totalChunks int,
	watchedQuotes []quote.WatchedPegoutQuote,
) ([]quote.WatchedPegoutQuote, error) {
	for len(watchedQuotes) > 0 && watchedQuotes[0].Remaining().Cmp(entities.NewWei(0)) == 0 {
		wq := watchedQuotes[0]
		wq.RetainedQuote.RemainingToRefund = entities.NewWei(0)
		wq.RetainedQuote.State = quote.PegoutStateBridgeTxSucceeded
		if err := h.persistQuoteUpdate(ctx, chunkIndex, totalChunks, wq); err != nil {
			return nil, err
		}
		watchedQuotes = watchedQuotes[1:]
	}
	return watchedQuotes, nil
}

func (h *UtxoSplitHandler) allocateChunkToQuote(
	wq quote.WatchedPegoutQuote,
	receipt blockchain.TransactionReceipt,
	chunkRemaining *entities.Wei,
) (quote.WatchedPegoutQuote, *entities.Wei) {
	quoteRemaining := wq.Remaining()
	wq.RetainedQuote.AppendRebalanceAllocation(receipt.TransactionHash, receipt.GasUsed.Uint64(), receipt.GasPrice)

	taken := new(entities.Wei).Min(quoteRemaining, chunkRemaining)
	quoteRemaining.Sub(quoteRemaining, taken)
	chunkRemaining.Sub(chunkRemaining, taken)

	if quoteRemaining.Cmp(entities.NewWei(0)) == 0 {
		wq.RetainedQuote.RemainingToRefund = entities.NewWei(0)
		wq.RetainedQuote.State = quote.PegoutStateBridgeTxSucceeded
	} else {
		wq.RetainedQuote.RemainingToRefund = quoteRemaining.Copy()
	}
	return wq, chunkRemaining
}

func (h *UtxoSplitHandler) persistQuoteUpdate(
	ctx context.Context,
	chunkIndex, totalChunks int,
	wq quote.WatchedPegoutQuote,
) error {
	wq.RetainedQuote.SetDeprecatedRefundFields()
	if err := h.quoteRepository.UpdateRetainedQuote(ctx, wq.RetainedQuote); err != nil {
		log.Errorf("%s: split tx %d/%d failed to persist quote update: %v", usecases.BridgePegoutId, chunkIndex+1, totalChunks, err)
		return err
	}
	return nil
}

func (h *UtxoSplitHandler) buildChunkAmounts(numTxs uint64, bridgeMin, remainder *entities.Wei) []*entities.Wei {
	chunkAmounts := make([]*entities.Wei, numTxs)
	chunkAmounts[0] = new(entities.Wei).Add(bridgeMin.Copy(), remainder)
	for i := uint64(1); i < numTxs; i++ {
		chunkAmounts[i] = bridgeMin.Copy()
	}
	return chunkAmounts
}

func (h *UtxoSplitHandler) adjustTotalForRetries(watchedQuotes []quote.WatchedPegoutQuote) *entities.Wei {
	adjusted := entities.NewWei(0)
	for _, wq := range watchedQuotes {
		adjusted.Add(adjusted, wq.Remaining())
	}
	return adjusted
}
