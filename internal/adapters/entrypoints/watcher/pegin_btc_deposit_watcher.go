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

type PeginDepositAddressWatcher struct {
	quotes                      map[string]quote.WatchedPeginQuote
	quotesMutex                 sync.RWMutex
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase
	updatePeginDepositUseCase   *w.UpdatePeginDepositUseCase
	callForUserUseCase          *pegin.CallForUserUseCase
	expiredUseCase              *pegin.ExpiredPeginQuoteUseCase
	btcWallet                   blockchain.BitcoinWallet
	rpc                         blockchain.Rpc
	ticker                      Ticker
	eventBus                    entities.EventBus
	watcherStopChannel          chan bool
	currentBlock                *big.Int
	currentBlockMutex           sync.RWMutex
}

type PeginDepositAddressWatcherUseCases struct {
	callForUserUseCase          *pegin.CallForUserUseCase
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase
	updatePeginDepositUseCase   *w.UpdatePeginDepositUseCase
	expiredUseCase              *pegin.ExpiredPeginQuoteUseCase
}

func NewPeginDepositAddressWatcherUseCases(
	callForUserUseCase *pegin.CallForUserUseCase,
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase,
	updatePeginDepositUseCase *w.UpdatePeginDepositUseCase,
	expiredUseCase *pegin.ExpiredPeginQuoteUseCase,
) *PeginDepositAddressWatcherUseCases {
	return &PeginDepositAddressWatcherUseCases{
		callForUserUseCase:          callForUserUseCase,
		getWatchedPeginQuoteUseCase: getWatchedPeginQuoteUseCase,
		updatePeginDepositUseCase:   updatePeginDepositUseCase,
		expiredUseCase:              expiredUseCase,
	}
}

const callForUserErrorTemplate = "Error executing call for user on quote %s: %v"

func NewPeginDepositAddressWatcher(
	useCases *PeginDepositAddressWatcherUseCases,
	btcWallet blockchain.BitcoinWallet,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
	ticker Ticker,
) *PeginDepositAddressWatcher {
	quotes := make(map[string]quote.WatchedPeginQuote)
	watcherStopChannel := make(chan bool, 1)
	return &PeginDepositAddressWatcher{
		quotes:                      quotes,
		quotesMutex:                 sync.RWMutex{},
		callForUserUseCase:          useCases.callForUserUseCase,
		updatePeginDepositUseCase:   useCases.updatePeginDepositUseCase,
		getWatchedPeginQuoteUseCase: useCases.getWatchedPeginQuoteUseCase,
		expiredUseCase:              useCases.expiredUseCase,
		btcWallet:                   btcWallet,
		eventBus:                    eventBus,
		watcherStopChannel:          watcherStopChannel,
		rpc:                         rpc,
		ticker:                      ticker,
		currentBlockMutex:           sync.RWMutex{},
	}
}

func (watcher *PeginDepositAddressWatcher) Prepare(ctx context.Context) error {
	var err error
	var depositAddress string
	var watchedQuotes []quote.WatchedPeginQuote
	watcher.currentBlockMutex.Lock()
	defer watcher.currentBlockMutex.Unlock()
	watcher.currentBlock = big.NewInt(0)

	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	watchedQuotes, err = watcher.getWatchedPeginQuoteUseCase.Run(ctx, quote.PeginStateWaitingForDeposit, quote.PeginStateWaitingForDepositConfirmations)
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
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.currentBlockMutex.Lock()
			if height, err := watcher.rpc.Btc.GetHeight(); err == nil && height.Cmp(watcher.currentBlock) > 0 {
				watcher.checkQuotes(context.Background())
				watcher.currentBlock = height
			} else if err != nil {
				log.Error(peginBtcDepositWatcherLog("error getting Bitcoin chain height: %v", err))
			}
			watcher.currentBlockMutex.Unlock()
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
	log.Debug(peginBtcDepositWatcherLog("shut down"))
}

func (watcher *PeginDepositAddressWatcher) handleAcceptedPeginQuote(event entities.Event) {
	parsedEvent, ok := event.(quote.AcceptedPeginQuoteEvent)
	quoteHash := parsedEvent.RetainedQuote.QuoteHash
	if !ok {
		log.Error(peginBtcDepositWatcherLog("trying to parse wrong event"))
		return
	}

	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	if _, alreadyHaveQuote := watcher.quotes[quoteHash]; alreadyHaveQuote {
		log.Info(peginBtcDepositWatcherLog("Quote %s is already watched", quoteHash))
		return
	}

	err := watcher.btcWallet.ImportAddress(parsedEvent.RetainedQuote.DepositAddress)
	if err != nil {
		log.Error(peginBtcDepositWatcherLog("error while importing deposit address (%s): %v", parsedEvent.RetainedQuote.DepositAddress, err))
		return
	}
	watcher.quotes[quoteHash] = quote.NewWatchedPeginQuote(parsedEvent.Quote, parsedEvent.RetainedQuote, parsedEvent.CreationData)
}

func (watcher *PeginDepositAddressWatcher) checkQuotes(ctx context.Context) {
	for _, watchedQuote := range watcher.quotes {
		watcher.handleQuote(ctx, watchedQuote)
	}
}

func (watcher *PeginDepositAddressWatcher) handleQuote(ctx context.Context, watchedQuote quote.WatchedPeginQuote) {
	var err error
	quoteHash := watchedQuote.RetainedQuote.QuoteHash

	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()

	if watchedQuote.RetainedQuote.State == quote.PeginStateWaitingForDeposit {
		if watchedQuote, err = watcher.handleNotDepositedQuote(ctx, watchedQuote); err != nil {
			log.Error(peginBtcDepositWatcherLog(callForUserErrorTemplate, quoteHash, err))
			return
		}
	}

	if watchedQuote.RetainedQuote.State == quote.PeginStateWaitingForDeposit && watchedQuote.PeginQuote.IsExpired() {
		if err = watcher.expiredUseCase.Run(ctx, watchedQuote.RetainedQuote); err != nil {
			log.Error(peginBtcDepositWatcherLog("Error updating expired quote (%s): %v", quoteHash, err))
		} else {
			log.Info(peginBtcDepositWatcherLog("Quote %s expired at %d", quoteHash, watchedQuote.PeginQuote.ExpireTime().Unix()))
			delete(watcher.quotes, quoteHash)
		}
		return
	}

	if watchedQuote.RetainedQuote.State == quote.PeginStateWaitingForDepositConfirmations {
		if err = watcher.handleDepositedQuote(ctx, watchedQuote); err != nil {
			log.Error(peginBtcDepositWatcherLog(callForUserErrorTemplate, quoteHash, err))
		}
		return
	}
}

func (watcher *PeginDepositAddressWatcher) handleNotDepositedQuote(ctx context.Context, watchedQuote quote.WatchedPeginQuote) (quote.WatchedPeginQuote, error) {
	var err error
	var block blockchain.BitcoinBlockInformation
	var txs []blockchain.BitcoinTransactionInformation
	var updatedQuote quote.WatchedPeginQuote

	depositAddress := watchedQuote.RetainedQuote.DepositAddress
	if txs, err = watcher.btcWallet.GetTransactions(depositAddress); err != nil {
		return quote.WatchedPeginQuote{}, err
	}

	for _, tx := range txs {
		log.Info(peginBtcDepositWatcherLog("Checking transaction %s for quote %s", tx.Hash, watchedQuote.RetainedQuote.QuoteHash))
		if block, err = watcher.rpc.Btc.GetTransactionBlockInfo(tx.Hash); err != nil {
			return quote.WatchedPeginQuote{}, err
		}
		onTime := watchedQuote.PeginQuote.ExpireTime().After(block.Time)
		correctAmount := tx.AmountToAddress(watchedQuote.RetainedQuote.DepositAddress).Cmp(watchedQuote.PeginQuote.Total()) >= 0
		if watchedQuote.RetainedQuote.State == quote.PeginStateWaitingForDeposit && onTime && correctAmount {
			if updatedQuote, err = watcher.updatePeginDepositUseCase.Run(ctx, watchedQuote, block, tx); err != nil {
				return quote.WatchedPeginQuote{}, err
			} else {
				watcher.quotes[watchedQuote.RetainedQuote.QuoteHash] = updatedQuote
				return updatedQuote, nil
			}
		} else {
			watcher.logRejectReason(block, tx, watchedQuote)
		}
	}
	return watchedQuote, nil
}

func (watcher *PeginDepositAddressWatcher) handleDepositedQuote(ctx context.Context, watchedQuote quote.WatchedPeginQuote) error {
	tx, err := watcher.rpc.Btc.GetTransactionInfo(watchedQuote.RetainedQuote.UserBtcTxHash)
	if err != nil {
		return err
	}
	if tx.Confirmations >= uint64(watchedQuote.PeginQuote.Confirmations) {
		watcher.callForUser(ctx, watchedQuote)
	}
	return nil
}

func (watcher *PeginDepositAddressWatcher) callForUser(ctx context.Context, watchedQuote quote.WatchedPeginQuote) {
	var err error
	quoteHash := watchedQuote.RetainedQuote.QuoteHash
	if err = watcher.callForUserUseCase.Run(ctx, watchedQuote.RetainedQuote); errors.Is(err, usecases.NonRecoverableError) {
		delete(watcher.quotes, quoteHash)
		log.Error(peginBtcDepositWatcherLog(callForUserErrorTemplate, quoteHash, err))
	} else if err != nil {
		log.Error(peginBtcDepositWatcherLog(callForUserErrorTemplate, quoteHash, err))
	} else {
		delete(watcher.quotes, quoteHash)
	}
}

func (watcher *PeginDepositAddressWatcher) GetWatchedQuote(quoteHash string) (quote.WatchedPeginQuote, bool) {
	watcher.quotesMutex.RLock()
	defer watcher.quotesMutex.RUnlock()
	watchedQuote, ok := watcher.quotes[quoteHash]
	return watchedQuote, ok
}

func (watcher *PeginDepositAddressWatcher) GetCurrentBlock() *big.Int {
	watcher.currentBlockMutex.RLock()
	defer watcher.currentBlockMutex.RUnlock()
	return watcher.currentBlock
}

func (watcher *PeginDepositAddressWatcher) logRejectReason(block blockchain.BitcoinBlockInformation, tx blockchain.BitcoinTransactionInformation, watchedQuote quote.WatchedPeginQuote) {
	rejectReason := fmt.Sprintf("Rejecting quote %s for the following reason: ", watchedQuote.RetainedQuote.QuoteHash)
	if watchedQuote.PeginQuote.ExpireTime().Before(block.Time) {
		blockTime := block.Time.Unix()
		expirationTime := watchedQuote.PeginQuote.ExpireTime().Unix()
		rejectReason += fmt.Sprintf("quote expired at %d, %d seconds before its first confirmation at %d;", expirationTime, blockTime-expirationTime, blockTime)
	}
	paidAmount := tx.AmountToAddress(watchedQuote.RetainedQuote.DepositAddress)
	expectedAmount := watchedQuote.PeginQuote.Total()
	if paidAmount.Cmp(expectedAmount) < 0 {
		rejectReason += fmt.Sprintf("transaction amount %s is less than expected %s;", paidAmount.ToSatoshi().String(), expectedAmount.ToSatoshi().String())
	}
	log.Info(peginBtcDepositWatcherLog(rejectReason))
}

func peginBtcDepositWatcherLog(msg string, args ...any) string {
	return fmt.Sprintf("PeginDepositAddressWatcher: "+msg, args...)
}
