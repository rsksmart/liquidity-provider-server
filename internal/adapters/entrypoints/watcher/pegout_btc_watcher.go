package watcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
	"math/big"
	"sync"
)

type PegoutBtcTransferWatcher struct {
	quotes                       map[string]quote.WatchedPegoutQuote
	quotesMutex                  sync.RWMutex
	getWatchedPegoutQuoteUseCase *w.GetWatchedPegoutQuoteUseCase
	refundPegoutUseCase          *pegout.RefundPegoutUseCase
	rpc                          blockchain.Rpc
	ticker                       Ticker
	eventBus                     entities.EventBus
	watcherStopChannel           chan bool
	currentBlock                 *big.Int
	currentBlockMutex            sync.RWMutex
}

func NewPegoutBtcTransferWatcher(
	getWatchedPegoutQuoteUseCase *w.GetWatchedPegoutQuoteUseCase,
	refundPegoutUseCase *pegout.RefundPegoutUseCase,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
	ticker Ticker,
) *PegoutBtcTransferWatcher {
	quotes := make(map[string]quote.WatchedPegoutQuote)
	watcherStopChannel := make(chan bool, 1)
	currentBlock := big.NewInt(0)
	return &PegoutBtcTransferWatcher{
		quotes:                       quotes,
		quotesMutex:                  sync.RWMutex{},
		getWatchedPegoutQuoteUseCase: getWatchedPegoutQuoteUseCase,
		refundPegoutUseCase:          refundPegoutUseCase,
		rpc:                          rpc,
		eventBus:                     eventBus,
		watcherStopChannel:           watcherStopChannel,
		currentBlock:                 currentBlock,
		ticker:                       ticker,
		currentBlockMutex:            sync.RWMutex{},
	}
}

func (watcher *PegoutBtcTransferWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug(pegoutBtcWatcherLog("shut down"))
}

func (watcher *PegoutBtcTransferWatcher) Prepare(ctx context.Context) error {
	watcher.currentBlockMutex.Lock()
	defer watcher.currentBlockMutex.Unlock()
	watcher.currentBlock = big.NewInt(0)
	watchedQuotes, err := watcher.getWatchedPegoutQuoteUseCase.Run(ctx, quote.PegoutStateSendPegoutSucceeded)
	if err != nil {
		return err
	}
	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	for _, watchedQuote := range watchedQuotes {
		watcher.quotes[watchedQuote.RetainedQuote.QuoteHash] = watchedQuote
	}
	return nil
}

func (watcher *PegoutBtcTransferWatcher) Start() {
	eventChannel := watcher.eventBus.Subscribe(quote.PegoutBtcSentEventId)
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.quotesMutex.Lock()
			watcher.currentBlockMutex.Lock()
			if height, err := watcher.rpc.Btc.GetHeight(); err == nil && height.Cmp(watcher.currentBlock) > 0 {
				watcher.checkQuotes()
				watcher.currentBlock = height
			} else if err != nil {
				log.Error(pegoutBtcWatcherLog(blockchain.BtcChainHeightErrorTemplate, err))
			}
			watcher.quotesMutex.Unlock()
			watcher.currentBlockMutex.Unlock()
		case event := <-eventChannel:
			if event != nil {
				watcher.handleBtcSentToUserCompleted(event)
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PegoutBtcTransferWatcher) checkQuotes() {
	var err error
	var tx blockchain.BitcoinTransactionInformation
	for _, watchedQuote := range watcher.quotes {
		if tx, err = watcher.rpc.Btc.GetTransactionInfo(watchedQuote.RetainedQuote.LpBtcTxHash); err != nil {
			log.Error(pegoutBtcWatcherLog(blockchain.BtcTxInfoErrorTemplate, watchedQuote.RetainedQuote.LpBtcTxHash, err))
			return
		}
		if watcher.validateQuote(watchedQuote, tx) {
			watcher.refundPegout(watchedQuote)
		}
	}
}

func (watcher *PegoutBtcTransferWatcher) refundPegout(watchedQuote quote.WatchedPegoutQuote) {
	var err error
	const refundPegoutErrorMsgTemplate = "Error executing refund pegout on quote %s: %v"
	if err = watcher.refundPegoutUseCase.Run(context.Background(), watchedQuote.RetainedQuote); errors.Is(err, usecases.NonRecoverableError) {
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
		log.Error(pegoutBtcWatcherLog(refundPegoutErrorMsgTemplate, watchedQuote.RetainedQuote.QuoteHash, err))
	} else if err != nil {
		log.Error(pegoutBtcWatcherLog(refundPegoutErrorMsgTemplate, watchedQuote.RetainedQuote.QuoteHash, err))
	} else {
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
	}
}

func (watcher *PegoutBtcTransferWatcher) handleBtcSentToUserCompleted(event entities.Event) {
	parsedEvent, ok := event.(quote.PegoutBtcSentToUserEvent)
	quoteHash := parsedEvent.RetainedQuote.QuoteHash
	if !ok {
		log.Error(pegoutBtcWatcherLog("Trying to parse wrong event in Pegout Bridge watcher"))
		return
	}

	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	if _, alreadyHaveQuote := watcher.quotes[quoteHash]; alreadyHaveQuote {
		log.Info(pegoutBtcWatcherLog("Quote %s is already watched", quoteHash))
		return
	}
	if parsedEvent.RetainedQuote.State != quote.PegoutStateSendPegoutSucceeded || parsedEvent.RetainedQuote.LpBtcTxHash == "" {
		log.Warn(pegoutBtcWatcherLog("Quote %s doesn't have btc tx hash to watch", quoteHash))
		return
	}
	watcher.quotes[quoteHash] = quote.NewWatchedPegoutQuote(parsedEvent.PegoutQuote, parsedEvent.RetainedQuote)
}

func (watcher *PegoutBtcTransferWatcher) GetWatchedQuote(quoteHash string) (quote.WatchedPegoutQuote, bool) {
	watcher.quotesMutex.RLock()
	defer watcher.quotesMutex.RUnlock()
	watchedQuote, ok := watcher.quotes[quoteHash]
	return watchedQuote, ok
}

func (watcher *PegoutBtcTransferWatcher) GetCurrentBlock() *big.Int {
	watcher.currentBlockMutex.RLock()
	defer watcher.currentBlockMutex.RUnlock()
	return watcher.currentBlock
}

func (watcher *PegoutBtcTransferWatcher) validateQuote(watchedQuote quote.WatchedPegoutQuote, tx blockchain.BitcoinTransactionInformation) bool {
	return watchedQuote.RetainedQuote.State == quote.PegoutStateSendPegoutSucceeded &&
		tx.Confirmations >= uint64(watchedQuote.PegoutQuote.TransferConfirmations)
}

func pegoutBtcWatcherLog(format string, args ...any) string {
	return fmt.Sprintf("PegoutBtcTransferWatcher: "+format, args...)
}
