package watcher

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
)

const pegoutBridgeWatcherLogPrefix = "PegoutBridgeWatcher:"

// PegoutBridgeWatcher is a watcher that checks the state of the pegout quotes and creates a transaction
// to send the value of multiple pegout quotes to the bridge to convert the refunded RBTC to BTC when
// a threshold is reached
type PegoutBridgeWatcher struct {
	getQuotesUseCase    *w.GetWatchedPegoutQuoteUseCase
	bridgePegoutUseCase *pegout.BridgePegoutUseCase
	ticker              Ticker
	watcherStopChannel  chan struct{}
}

func NewPegoutBridgeWatcher(getQuotesUseCase *w.GetWatchedPegoutQuoteUseCase, bridgePegoutUseCase *pegout.BridgePegoutUseCase, ticker Ticker) *PegoutBridgeWatcher {
	return &PegoutBridgeWatcher{
		getQuotesUseCase:    getQuotesUseCase,
		bridgePegoutUseCase: bridgePegoutUseCase,
		watcherStopChannel:  make(chan struct{}, 1),
		ticker:              ticker,
	}
}

func (watcher *PegoutBridgeWatcher) Prepare(ctx context.Context) error { return nil }

func (watcher *PegoutBridgeWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.runUseCases()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PegoutBridgeWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug(pegoutBridgeWatcherLogPrefix + " shut down")
}

func (watcher *PegoutBridgeWatcher) runUseCases() {
	ctx := context.Background()
	quotes, err := watcher.getQuotesUseCase.Run(ctx, quote.PegoutStateRefundPegOutSucceeded)
	if err != nil {
		log.Errorf("%s error getting pegout quotes: %v", pegoutBridgeWatcherLogPrefix, err)
		return
	}
	err = watcher.bridgePegoutUseCase.Run(ctx, quotes...)
	if err != nil && !errors.Is(err, usecases.TxBelowMinimumError) {
		log.Errorf("%s error sending pegout to bridge: %v", pegoutBridgeWatcherLogPrefix, err)
		return
	}
}
