package watcher

import (
	"context"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
	"math/big"
	"sync"
)

// PeginBridgeWatcher is a watcher that checks the state of the pegin quotes and registers
// the pegin on the bridge when the conditions are met
type PeginBridgeWatcher struct {
	quotes                      map[string]quote.WatchedPeginQuote
	quotesMutex                 sync.RWMutex
	registerPeginUseCase        *pegin.RegisterPeginUseCase
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase
	contracts                   blockchain.RskContracts
	rpc                         blockchain.Rpc
	ticker                      Ticker
	eventBus                    entities.EventBus
	watcherStopChannel          chan bool
	currentBlock                *big.Int
	currentBlockMutex           sync.RWMutex
}

func NewPeginBridgeWatcher(
	registerPeginUseCase *pegin.RegisterPeginUseCase,
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
	ticker Ticker,
) *PeginBridgeWatcher {
	quotes := make(map[string]quote.WatchedPeginQuote)
	watcherStopChannel := make(chan bool, 1)
	return &PeginBridgeWatcher{
		quotes:                      quotes,
		registerPeginUseCase:        registerPeginUseCase,
		getWatchedPeginQuoteUseCase: getWatchedPeginQuoteUseCase,
		contracts:                   contracts,
		rpc:                         rpc,
		eventBus:                    eventBus,
		watcherStopChannel:          watcherStopChannel,
		ticker:                      ticker,
		quotesMutex:                 sync.RWMutex{},
		currentBlockMutex:           sync.RWMutex{},
	}
}

func (watcher *PeginBridgeWatcher) Prepare(ctx context.Context) error {
	watcher.currentBlockMutex.Lock()
	defer watcher.currentBlockMutex.Unlock()
	watcher.currentBlock = big.NewInt(0)
	watchedQuotes, err := watcher.getWatchedPeginQuoteUseCase.Run(ctx, quote.PeginStateCallForUserSucceeded)
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

func (watcher *PeginBridgeWatcher) Start() {
	eventChannel := watcher.eventBus.Subscribe(quote.CallForUserCompletedEventId)
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.currentBlockMutex.Lock()
			watcher.quotesMutex.Lock()
			if height, err := watcher.rpc.Btc.GetHeight(); err == nil && height.Cmp(watcher.currentBlock) > 0 {
				watcher.checkQuotes()
				watcher.currentBlock = height
			} else if err != nil {
				log.Error(peginBridgeWatcherLog(blockchain.BtcChainHeightErrorTemplate, err))
			}
			watcher.currentBlockMutex.Unlock()
			watcher.quotesMutex.Unlock()
		case event := <-eventChannel:
			if event != nil {
				watcher.handleCallForUserCompleted(event)
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PeginBridgeWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug(peginBridgeWatcherLog("shut down"))
}

func (watcher *PeginBridgeWatcher) handleCallForUserCompleted(event entities.Event) {
	parsedEvent, ok := event.(quote.CallForUserCompletedEvent)
	quoteHash := parsedEvent.RetainedQuote.QuoteHash
	if !ok {
		log.Error(peginBridgeWatcherLog("Trying to parse wrong event"))
		return
	}
	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	if _, alreadyHaveQuote := watcher.quotes[quoteHash]; alreadyHaveQuote {
		log.Info(peginBridgeWatcherLog("Quote %s is already watched", quoteHash))
		return
	}
	if parsedEvent.RetainedQuote.State == quote.PeginStateCallForUserSucceeded {
		watcher.quotes[quoteHash] = quote.NewWatchedPeginQuote(parsedEvent.PeginQuote, parsedEvent.RetainedQuote, parsedEvent.CreationData)
	}
}

func (watcher *PeginBridgeWatcher) checkQuotes() {
	var err error
	var tx blockchain.BitcoinTransactionInformation
	for _, watchedQuote := range watcher.quotes {
		if tx, err = watcher.rpc.Btc.GetTransactionInfo(watchedQuote.RetainedQuote.UserBtcTxHash); err != nil {
			log.Error(peginBridgeWatcherLog(blockchain.BtcTxInfoErrorTemplate, watchedQuote.RetainedQuote.UserBtcTxHash, err))
			return
		}
		if watcher.validateQuote(watchedQuote, tx) {
			watcher.registerPegin(watchedQuote)
		}
	}
}

func (watcher *PeginBridgeWatcher) registerPegin(watchedQuote quote.WatchedPeginQuote) {
	var err error
	const registerPeginErrorMsgTemplate = "Error executing register pegin on quote %s: %v"
	if err = watcher.registerPeginUseCase.Run(context.Background(), watchedQuote.RetainedQuote); errors.Is(err, usecases.NonRecoverableError) {
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
		log.Error(peginBridgeWatcherLog(registerPeginErrorMsgTemplate, watchedQuote.RetainedQuote.QuoteHash, err))
	} else if err != nil {
		log.Error(peginBridgeWatcherLog(registerPeginErrorMsgTemplate, watchedQuote.RetainedQuote.QuoteHash, err))
	} else {
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
	}
}

func (watcher *PeginBridgeWatcher) validateQuote(watchedQuote quote.WatchedPeginQuote, tx blockchain.BitcoinTransactionInformation) bool {
	return watchedQuote.RetainedQuote.State == quote.PeginStateCallForUserSucceeded &&
		tx.Confirmations >= watcher.contracts.Bridge.GetRequiredTxConfirmations()
}

func (watcher *PeginBridgeWatcher) GetWatchedQuote(quoteHash string) (quote.WatchedPeginQuote, bool) {
	watcher.quotesMutex.RLock()
	defer watcher.quotesMutex.RUnlock()
	watchedQuote, ok := watcher.quotes[quoteHash]
	return watchedQuote, ok
}

func (watcher *PeginBridgeWatcher) GetCurrentBlock() *big.Int {
	watcher.currentBlockMutex.RLock()
	defer watcher.currentBlockMutex.RUnlock()
	return watcher.currentBlock
}

func peginBridgeWatcherLog(msg string, args ...any) string {
	return fmt.Sprintf("Pegin Bridge watcher: "+msg, args...)
}
