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
	"time"
)

type PeginDepositAddressWatcher struct {
	quotes                      map[string]w.WatchedPeginQuote
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase
	callForUserUseCase          *pegin.CallForUserUseCase
	expiredUseCase              *pegin.ExpiredPeginQuoteUseCase
	btcWallet                   blockchain.BitcoinWallet
	btcRpc                      blockchain.BitcoinNetwork
	ticker                      *time.Ticker
	eventBus                    entities.EventBus
	watcherStopChannel          chan bool
	currentBlock                *big.Int
}

func NewPeginDepositAddressWatcher(
	callForUserUseCase *pegin.CallForUserUseCase,
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase,
	expiredUseCase *pegin.ExpiredPeginQuoteUseCase,
	btcWallet blockchain.BitcoinWallet,
	btcRpc blockchain.BitcoinNetwork,
	eventBus entities.EventBus,
) *PeginDepositAddressWatcher {
	quotes := make(map[string]w.WatchedPeginQuote)
	watcherStopChannel := make(chan bool, 1)
	return &PeginDepositAddressWatcher{
		quotes:                      quotes,
		callForUserUseCase:          callForUserUseCase,
		getWatchedPeginQuoteUseCase: getWatchedPeginQuoteUseCase,
		expiredUseCase:              expiredUseCase,
		btcWallet:                   btcWallet,
		eventBus:                    eventBus,
		watcherStopChannel:          watcherStopChannel,
		btcRpc:                      btcRpc,
	}
}

func (watcher *PeginDepositAddressWatcher) Prepare(ctx context.Context) error {
	var err error
	var depositAddress string
	var watchedQuotes []w.WatchedPeginQuote
	watcher.currentBlock = big.NewInt(0)
	watchedQuotes, err = watcher.getWatchedPeginQuoteUseCase.Run(ctx, quote.PeginStateWaitingForDeposit)
	if err != nil {
		return err
	}
	for _, watchedQuote := range watchedQuotes {
		depositAddress = watchedQuote.RetainedQuote.DepositAddress
		if err = watcher.btcWallet.ImportAddress(depositAddress); err != nil {
			return fmt.Errorf("error while importing deposit address (%s): %w\n", depositAddress, err)
		}
		watcher.quotes[watchedQuote.RetainedQuote.QuoteHash] = watchedQuote
	}
	return nil
}

func (watcher *PeginDepositAddressWatcher) Start() {
	eventChannel := watcher.eventBus.Subscribe(quote.AcceptedPeginQuoteEventId)
	watcher.ticker = time.NewTicker(peginDepositWatcherInterval)
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C:
			if height, err := watcher.btcRpc.GetHeight(); err == nil && height.Cmp(watcher.currentBlock) > 0 {
				watcher.checkQuotes()
				watcher.currentBlock = height
			} else if err != nil {
				log.Error("PeginDepositAddressWatcher: error getting Bitcoin chain height: ", err)
			}
		case event := <-eventChannel:
			if event != nil {
				watcher.handleAcceptedPeginQuote(event)
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PeginDepositAddressWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug("PeginDepositAddressWatcher shut down")
}

func (watcher *PeginDepositAddressWatcher) handleAcceptedPeginQuote(event entities.Event) {
	parsedEvent, ok := event.(quote.AcceptedPeginQuoteEvent)
	quoteHash := parsedEvent.RetainedQuote.QuoteHash
	if !ok {
		log.Error("Trying to parse wrong event in Pegin Btc deposit watcher")
		return
	}

	if _, alreadyHaveQuote := watcher.quotes[quoteHash]; alreadyHaveQuote {
		log.Infof("Quote %s is already watched\n", quoteHash)
		return
	}

	err := watcher.btcWallet.ImportAddress(parsedEvent.RetainedQuote.DepositAddress)
	if err != nil {
		log.Errorf("error while importing deposit address (%s): %v\n", parsedEvent.RetainedQuote.DepositAddress, err)
		return
	}
	watcher.quotes[quoteHash] = w.NewWatchedPeginQuote(parsedEvent.Quote, parsedEvent.RetainedQuote)
}

func (watcher *PeginDepositAddressWatcher) checkQuotes() {
	for _, watchedQuote := range watcher.quotes {
		watcher.handleQuote(watchedQuote)
	}
}

func (watcher *PeginDepositAddressWatcher) handleQuote(watchedQuote w.WatchedPeginQuote) {
	quoteHash := watchedQuote.RetainedQuote.QuoteHash
	depositAddress := watchedQuote.RetainedQuote.DepositAddress
	txs, err := watcher.btcWallet.GetTransactions(depositAddress)
	if err != nil {
		log.Errorf("Error executing call for user on quote %s: %v\n", quoteHash, err)
		return
	}
	for _, tx := range txs {
		if validatePeginQuote(watchedQuote, tx) {
			watcher.callForUser(watchedQuote, tx)
			return
		}
	}
	if watchedQuote.RetainedQuote.State == quote.PeginStateWaitingForDeposit && watchedQuote.PeginQuote.IsExpired() {
		if err = watcher.expiredUseCase.Run(context.Background(), watchedQuote.RetainedQuote); err != nil {
			log.Errorf("Error updating expired quote (%s): %v\n", quoteHash, err)
		} else {
			delete(watcher.quotes, quoteHash)
		}
	}
}

func (watcher *PeginDepositAddressWatcher) callForUser(watchedQuote w.WatchedPeginQuote, tx blockchain.BitcoinTransactionInformation) {
	var err error
	quoteHash := watchedQuote.RetainedQuote.QuoteHash
	if err = watcher.callForUserUseCase.Run(context.Background(), tx.Hash, watchedQuote.RetainedQuote); errors.Is(err, usecases.NonRecoverableError) {
		delete(watcher.quotes, quoteHash)
		log.Errorf("Error executing call for user on quote %s: %v\n", quoteHash, err)
	} else if err != nil {
		log.Errorf("Error executing call for user on quote %s: %v\n", quoteHash, err)
	} else {
		delete(watcher.quotes, quoteHash)
	}
}

func validatePeginQuote(watchedQuote w.WatchedPeginQuote, tx blockchain.BitcoinTransactionInformation) bool {
	return tx.Confirmations >= uint64(watchedQuote.PeginQuote.Confirmations) &&
		watchedQuote.RetainedQuote.State == quote.PeginStateWaitingForDeposit &&
		!watchedQuote.PeginQuote.IsExpired() &&
		tx.AmountToAddress(watchedQuote.RetainedQuote.DepositAddress).Cmp(watchedQuote.PeginQuote.Total()) >= 0
}
