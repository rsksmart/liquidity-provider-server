package watcher

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
	"math/big"
	"time"
)

type PeginBridgeWatcher struct {
	quotes                      map[string]w.WatchedPeginQuote
	registerPeginUseCase        *pegin.RegisterPeginUseCase
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase
	contracts                   blockchain.RskContracts
	rpc                         blockchain.Rpc
	ticker                      *time.Ticker
	eventBus                    entities.EventBus
	watcherStopChannel          chan bool
	currentBlock                *big.Int
}

func NewPeginBridgeWatcher(
	registerPeginUseCase *pegin.RegisterPeginUseCase,
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase,
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
) *PeginBridgeWatcher {
	quotes := make(map[string]w.WatchedPeginQuote)
	watcherStopChannel := make(chan bool, 1)
	return &PeginBridgeWatcher{
		quotes:                      quotes,
		registerPeginUseCase:        registerPeginUseCase,
		getWatchedPeginQuoteUseCase: getWatchedPeginQuoteUseCase,
		contracts:                   contracts,
		rpc:                         rpc,
		eventBus:                    eventBus,
		watcherStopChannel:          watcherStopChannel,
	}
}

func (watcher *PeginBridgeWatcher) Prepare(ctx context.Context) error {
	watcher.currentBlock = big.NewInt(0)
	watchedQuotes, err := watcher.getWatchedPeginQuoteUseCase.Run(ctx, quote.PeginStateCallForUserSucceeded)
	if err != nil {
		return err
	}
	for _, watchedQuote := range watchedQuotes {
		watcher.quotes[watchedQuote.RetainedQuote.QuoteHash] = watchedQuote
	}
	return nil
}

func (watcher *PeginBridgeWatcher) Start() {
	eventChannel := watcher.eventBus.Subscribe(quote.CallForUserCompletedEventId)
	watcher.ticker = time.NewTicker(peginBridgeWatcherInterval)
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C:
			if height, err := watcher.rpc.Btc.GetHeight(); err == nil && height.Cmp(watcher.currentBlock) > 0 {
				watcher.checkQuotes()
				watcher.currentBlock = height
			} else if err != nil {
				log.Error("PeginBridgeWatcher: error getting Bitcoin chain height: ", err)
			}
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
	log.Debug("PeginBridgeWatcher shut down")
}

func (watcher *PeginBridgeWatcher) handleCallForUserCompleted(event entities.Event) {
	parsedEvent, ok := event.(quote.CallForUserCompletedEvent)
	quoteHash := parsedEvent.RetainedQuote.QuoteHash
	if !ok {
		log.Error("Trying to parse wrong event in Pegin Bridge watcher")
		return
	}

	if _, alreadyHaveQuote := watcher.quotes[quoteHash]; alreadyHaveQuote {
		log.Infof("Quote %s is already watched\n", quoteHash)
		return
	}
	if parsedEvent.RetainedQuote.State == quote.PeginStateCallForUserSucceeded {
		watcher.quotes[quoteHash] = w.NewWatchedPeginQuote(parsedEvent.PeginQuote, parsedEvent.RetainedQuote)
	}
}

func (watcher *PeginBridgeWatcher) checkQuotes() {
	var err error
	var tx blockchain.BitcoinTransactionInformation
	for _, watchedQuote := range watcher.quotes {
		if tx, err = watcher.rpc.Btc.GetTransactionInfo(watchedQuote.RetainedQuote.UserBtcTxHash); err != nil {
			log.Errorf("Error getting Bitcoin transaction information %s: %v\n", watchedQuote.RetainedQuote.UserBtcTxHash, err)
			return
		}
		if watcher.validateQuote(watchedQuote, tx) {
			watcher.registerPegin(watchedQuote)
		}
	}
}

func (watcher *PeginBridgeWatcher) registerPegin(watchedQuote w.WatchedPeginQuote) {
	var err error
	if err = watcher.registerPeginUseCase.Run(context.Background(), watchedQuote.RetainedQuote); errors.Is(err, usecases.NonRecoverableError) {
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
		log.Errorf("Error executing register pegin on quote %s: %v\n", watchedQuote.RetainedQuote.QuoteHash, err)
	} else if err != nil {
		log.Errorf("Error executing register pegin on quote %s: %v\n", watchedQuote.RetainedQuote.QuoteHash, err)
	} else {
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
	}
}

func (watcher *PeginBridgeWatcher) validateQuote(watchedQuote w.WatchedPeginQuote, tx blockchain.BitcoinTransactionInformation) bool {
	return watchedQuote.RetainedQuote.State == quote.PeginStateCallForUserSucceeded &&
		tx.Confirmations >= watcher.contracts.Bridge.GetRequiredTxConfirmations()
}
