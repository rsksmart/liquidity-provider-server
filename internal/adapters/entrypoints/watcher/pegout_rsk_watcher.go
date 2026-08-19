package watcher

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
)

type PegoutRskDepositWatcher struct {
	quotes                       map[string]quote.WatchedPegoutQuote
	quotesMutex                  sync.RWMutex
	getWatchedPegoutQuoteUseCase *w.GetWatchedPegoutQuoteUseCase
	sendPegoutUseCase            *pegout.SendPegoutUseCase
	rpc                          blockchain.Rpc
	ticker                       utils.Ticker
	eventBus                     entities.EventBus
	watcherStopChannel           chan bool
	currentBlock                 uint64
	currentBlockMutex            sync.RWMutex
	checkTimeout                 time.Duration
}

type PegoutRskDepositWatcherUseCases struct {
	getWatchedPegoutQuoteUseCase *w.GetWatchedPegoutQuoteUseCase
	sendPegoutUseCase            *pegout.SendPegoutUseCase
}

func NewPegoutRskDepositWatcherUseCases(
	getWatchedPegoutQuoteUseCase *w.GetWatchedPegoutQuoteUseCase,
	sendPegoutUseCase *pegout.SendPegoutUseCase,
) *PegoutRskDepositWatcherUseCases {
	return &PegoutRskDepositWatcherUseCases{
		getWatchedPegoutQuoteUseCase: getWatchedPegoutQuoteUseCase,
		sendPegoutUseCase:            sendPegoutUseCase,
	}
}

func NewPegoutRskDepositWatcher(
	watcherUseCases *PegoutRskDepositWatcherUseCases,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
	ticker utils.Ticker,
	checkTimeout time.Duration,
) *PegoutRskDepositWatcher {
	return &PegoutRskDepositWatcher{
		quotes:                       make(map[string]quote.WatchedPegoutQuote),
		getWatchedPegoutQuoteUseCase: watcherUseCases.getWatchedPegoutQuoteUseCase,
		sendPegoutUseCase:            watcherUseCases.sendPegoutUseCase,
		rpc:                          rpc,
		eventBus:                     eventBus,
		watcherStopChannel:           make(chan bool, 1),
		ticker:                       ticker,
		currentBlockMutex:            sync.RWMutex{},
		quotesMutex:                  sync.RWMutex{},
		checkTimeout:                 checkTimeout,
	}
}

func (watcher *PegoutRskDepositWatcher) Prepare(ctx context.Context) error {
	height, err := watcher.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return err
	}
	watcher.currentBlockMutex.Lock()
	watcher.currentBlock = height
	watcher.currentBlockMutex.Unlock()

	watchedQuotes, err := watcher.getWatchedPegoutQuoteUseCase.Run(ctx, quote.PegoutStateClaimed)
	if err != nil {
		return err
	}
	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	for _, watchedQuote := range watchedQuotes {
		watcher.quotes[watchedQuote.RetainedQuote.QuoteHash] = watchedQuote
	}

	log.Info(LogPegoutRskStart(watcher.currentBlock))
	return nil
}

func (watcher *PegoutRskDepositWatcher) Start() {
	claimedEventChannel := watcher.eventBus.Subscribe(quote.ClaimedPegoutQuoteEventId)

watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.onTick()
		case event := <-claimedEventChannel:
			if event != nil {
				watcher.handleClaimedPegoutQuote(event)
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PegoutRskDepositWatcher) onTick() {
	watcher.currentBlockMutex.Lock()
	watcher.quotesMutex.Lock()
	defer watcher.currentBlockMutex.Unlock()
	defer watcher.quotesMutex.Unlock()

	checkContext, checkCancel := context.WithTimeout(context.Background(), watcher.checkTimeout)
	defer checkCancel()
	height, err := watcher.rpc.Rsk.GetHeight(checkContext)
	if err != nil {
		log.Errorf(LogPegoutRskChainHeight, err)
		return
	}
	if height > watcher.currentBlock {
		watcher.checkQuotes(checkContext)
		watcher.currentBlock = height
	}
}

func (watcher *PegoutRskDepositWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug(LogPegoutRskShutdown)
}

func (watcher *PegoutRskDepositWatcher) handleClaimedPegoutQuote(event entities.Event) {
	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	parsedEvent, ok := event.(quote.ClaimedPegoutQuoteEvent)
	if !ok {
		log.Error(LogPegoutRskWrongEvent)
		return
	}
	quoteHash := parsedEvent.RetainedQuote.QuoteHash

	if _, alreadyHaveQuote := watcher.quotes[quoteHash]; alreadyHaveQuote {
		log.Info(LogPegoutRskAlreadyWatched(quoteHash))
		return
	}
	watcher.quotes[quoteHash] = quote.NewWatchedPegoutQuote(
		parsedEvent.Quote,
		parsedEvent.RetainedQuote,
		quote.PegoutCreationDataZeroValue(),
	)
}

func (watcher *PegoutRskDepositWatcher) checkQuotes(ctx context.Context) {
	for _, watchedQuote := range watcher.quotes {
		if watchedQuote.RetainedQuote.State == quote.PegoutStateClaimed {
			watcher.sendPegout(ctx, watchedQuote)
		}
	}
}

func (watcher *PegoutRskDepositWatcher) sendPegout(ctx context.Context, watchedQuote quote.WatchedPegoutQuote) {
	err := watcher.sendPegoutUseCase.Run(ctx, watchedQuote.RetainedQuote)
	if errors.Is(err, usecases.NonRecoverableError) {
		log.Error(LogPegoutRskSendError(watchedQuote.RetainedQuote.QuoteHash, err))
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
	} else if err != nil {
		log.Error(LogPegoutRskSendError(watchedQuote.RetainedQuote.QuoteHash, err))
	} else {
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
	}
}

func (watcher *PegoutRskDepositWatcher) GetWatchedQuote(quoteHash string) (quote.WatchedPegoutQuote, bool) {
	watcher.quotesMutex.RLock()
	defer watcher.quotesMutex.RUnlock()
	watchedQuote, ok := watcher.quotes[quoteHash]
	return watchedQuote, ok
}

func (watcher *PegoutRskDepositWatcher) GetCurrentBlock() uint64 {
	watcher.currentBlockMutex.RLock()
	defer watcher.currentBlockMutex.RUnlock()
	return watcher.currentBlock
}
