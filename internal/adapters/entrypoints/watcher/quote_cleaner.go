package watcher

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
)

type QuoteCleanerWatcher struct {
	cleanUseCase       *watcher.CleanExpiredQuotesUseCase
	ticker             Ticker
	watcherStopChannel chan bool
}

func NewQuoteCleanerWatcher(
	cleanUseCase *watcher.CleanExpiredQuotesUseCase,
	ticker Ticker,
) *QuoteCleanerWatcher {
	watcherStopChannel := make(chan bool, 1)
	return &QuoteCleanerWatcher{
		cleanUseCase:       cleanUseCase,
		watcherStopChannel: watcherStopChannel,
		ticker:             ticker,
	}
}

func (watcher *QuoteCleanerWatcher) Prepare(ctx context.Context) error {
	return nil
}

func (watcher *QuoteCleanerWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.clean()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *QuoteCleanerWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug("QuoteCleanerWatcher shut down")
}

func (watcher *QuoteCleanerWatcher) clean() {
	txIds, err := watcher.cleanUseCase.Run(context.Background())
	if err != nil {
		log.Error("Error cleaning quotes: ", err)
	}
	log.Infof("Cleaned %d quotes:\n", len(txIds))
	for _, id := range txIds {
		log.Infof("Quote %s cleaned\n", id)
	}
}
