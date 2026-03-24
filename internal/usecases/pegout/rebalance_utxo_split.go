package pegout

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type UtxoSplitHandler struct {
	quoteRepository quote.PegoutQuoteRepository
	rskWallet       blockchain.RootstockWallet
	contracts       blockchain.RskContracts
}

func NewUtxoSplitHandler(
	quoteRepository quote.PegoutQuoteRepository,
	rskWallet blockchain.RootstockWallet,
	contracts blockchain.RskContracts,
) *UtxoSplitHandler {
	return &UtxoSplitHandler{
		quoteRepository: quoteRepository,
		rskWallet:       rskWallet,
		contracts:       contracts,
	}
}

func (h *UtxoSplitHandler) Execute(
	ctx context.Context,
	pegoutConfig liquidity_provider.PegoutConfiguration,
	watchedQuotes []quote.WatchedPegoutQuote,
) error {
	bridgeMin := pegoutConfig.BridgeTransactionMin
	adjustedTotal := adjustTotalForRetries(watchedQuotes)
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
	if err := checkBalance(ctx, h.rskWallet, requiredBalance); err != nil {
		return err
	}

	chunkAmounts := buildChunkAmounts(numTxs, bridgeMin, remainder)
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
	if txErr != nil {
		log.Errorf("%s: split tx %d/%d failed: %v", usecases.BridgePegoutId, chunkIndex+1, totalChunks, txErr)
		return watchedQuotes, nil
	}
	if !isValidSplitReceipt(receipt, chunkIndex, totalChunks) {
		return watchedQuotes, nil
	}
	return h.distributeChunkReceipt(ctx, chunkIndex, totalChunks, receipt, watchedQuotes)
}

func isValidSplitReceipt(receipt blockchain.TransactionReceipt, chunkIndex, totalChunks int) bool {
	if receipt.TransactionHash == "" {
		log.Errorf("%s: split tx %d/%d failed: incomplete receipt", usecases.BridgePegoutId, chunkIndex+1, totalChunks)
		return false
	}
	log.Debugf("%s: split tx %d/%d sent to the bridge successfully (%s)", usecases.BridgePegoutId, chunkIndex+1, totalChunks, receipt.TransactionHash)
	if receipt.Value == nil {
		log.Errorf("%s: split tx %d/%d failed: missing receipt value", usecases.BridgePegoutId, chunkIndex+1, totalChunks)
		return false
	}
	return true
}

func (h *UtxoSplitHandler) distributeChunkReceipt(
	ctx context.Context,
	chunkIndex, totalChunks int,
	receipt blockchain.TransactionReceipt,
	watchedQuotes []quote.WatchedPegoutQuote,
) ([]quote.WatchedPegoutQuote, error) {
	zero := entities.NewWei(0)
	chunkRemaining := receipt.Value.Copy()
	for chunkRemaining.Cmp(zero) > 0 && len(watchedQuotes) != 0 {
		watchedQuote := watchedQuotes[0]
		quoteRemaining := getRemaining(watchedQuote)
		if quoteRemaining.Cmp(zero) == 0 {
			watchedQuote.RetainedQuote.RemainingToRefund = entities.NewWei(0)
			watchedQuote.RetainedQuote.State = quote.PegoutStateBridgeTxSucceeded
			setDeprecatedRefundFields(&watchedQuote.RetainedQuote)
			if err := h.quoteRepository.UpdateRetainedQuote(ctx, watchedQuote.RetainedQuote); err != nil {
				log.Errorf("%s: split tx %d/%d failed to persist quote update: %v", usecases.BridgePegoutId, chunkIndex+1, totalChunks, err)
				return nil, err
			}
			watchedQuotes = watchedQuotes[1:]
			continue
		}
		appendAllocation(&watchedQuote.RetainedQuote, receipt)
		taken := minWei(quoteRemaining, chunkRemaining)
		quoteRemaining.Sub(quoteRemaining, taken)
		chunkRemaining.Sub(chunkRemaining, taken)
		if quoteRemaining.Cmp(zero) == 0 {
			watchedQuote.RetainedQuote.RemainingToRefund = entities.NewWei(0)
			watchedQuote.RetainedQuote.State = quote.PegoutStateBridgeTxSucceeded
		} else {
			watchedQuote.RetainedQuote.RemainingToRefund = quoteRemaining.Copy()
		}
		setDeprecatedRefundFields(&watchedQuote.RetainedQuote)
		if err := h.quoteRepository.UpdateRetainedQuote(ctx, watchedQuote.RetainedQuote); err != nil {
			log.Errorf("%s: split tx %d/%d failed to persist quote update: %v", usecases.BridgePegoutId, chunkIndex+1, totalChunks, err)
			return nil, err
		}
		if quoteRemaining.Cmp(zero) == 0 {
			watchedQuotes = watchedQuotes[1:]
			continue
		}
		watchedQuotes[0] = watchedQuote
		break
	}
	return watchedQuotes, nil
}

func buildChunkAmounts(numTxs uint64, bridgeMin, remainder *entities.Wei) []*entities.Wei {
	chunkAmounts := make([]*entities.Wei, numTxs)
	chunkAmounts[0] = new(entities.Wei).Add(bridgeMin.Copy(), remainder)
	for i := uint64(1); i < numTxs; i++ {
		chunkAmounts[i] = bridgeMin.Copy()
	}
	return chunkAmounts
}

func adjustTotalForRetries(watchedQuotes []quote.WatchedPegoutQuote) *entities.Wei {
	adjusted := entities.NewWei(0)
	for _, wq := range watchedQuotes {
		adjusted.Add(adjusted, getRemaining(wq))
	}
	return adjusted
}

func getRemaining(wq quote.WatchedPegoutQuote) *entities.Wei {
	if wq.RetainedQuote.RemainingToRefund != nil {
		return wq.RetainedQuote.RemainingToRefund.Copy()
	}
	return quoteContribution(wq)
}

func appendAllocation(retained *quote.RetainedPegoutQuote, receipt blockchain.TransactionReceipt) {
	if receipt.TransactionHash == "" {
		return
	}
	retained.BridgeRebalances = append(retained.BridgeRebalances, quote.BridgeRebalanceAllocation{
		TxHash:   receipt.TransactionHash,
		GasUsed:  receipt.GasUsed.Uint64(),
		GasPrice: receipt.GasPrice,
	})
}

func setDeprecatedRefundFields(retained *quote.RetainedPegoutQuote) {
	if len(retained.BridgeRebalances) == 0 || retained.BridgeRefundTxHash != "" {
		return
	}
	first := retained.BridgeRebalances[0]
	retained.BridgeRefundTxHash = first.TxHash
	retained.BridgeRefundGasUsed = first.GasUsed
	retained.BridgeRefundGasPrice = first.GasPrice
}

func minWei(a, b *entities.Wei) *entities.Wei {
	if a.Cmp(b) <= 0 {
		return a.Copy()
	}
	return b.Copy()
}
