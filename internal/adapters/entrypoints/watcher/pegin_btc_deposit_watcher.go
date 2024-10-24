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
	quotes                      map[string]quote.WatchedPeginQuote
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase
	updatePeginDepositUseCase   *w.UpdatePeginDepositUseCase
	callForUserUseCase          *pegin.CallForUserUseCase
	expiredUseCase              *pegin.ExpiredPeginQuoteUseCase
	btcWallet                   blockchain.BitcoinWallet
	rpc                         blockchain.Rpc
	ticker                      *time.Ticker
	eventBus                    entities.EventBus
	watcherStopChannel          chan bool
	currentBlock                *big.Int
}

const callForUserErrorTemplate = "Error executing call for user on quote %s: %v"

func NewPeginDepositAddressWatcher(
	callForUserUseCase *pegin.CallForUserUseCase,
	getWatchedPeginQuoteUseCase *w.GetWatchedPeginQuoteUseCase,
	updatePeginDepositUseCase *w.UpdatePeginDepositUseCase,
	expiredUseCase *pegin.ExpiredPeginQuoteUseCase,
	btcWallet blockchain.BitcoinWallet,
	rpc blockchain.Rpc,
	eventBus entities.EventBus,
) *PeginDepositAddressWatcher {
	quotes := make(map[string]quote.WatchedPeginQuote)
	watcherStopChannel := make(chan bool, 1)
	return &PeginDepositAddressWatcher{
		quotes:                      quotes,
		callForUserUseCase:          callForUserUseCase,
		updatePeginDepositUseCase:   updatePeginDepositUseCase,
		getWatchedPeginQuoteUseCase: getWatchedPeginQuoteUseCase,
		expiredUseCase:              expiredUseCase,
		btcWallet:                   btcWallet,
		eventBus:                    eventBus,
		watcherStopChannel:          watcherStopChannel,
		rpc:                         rpc,
	}
}

func (watcher *PeginDepositAddressWatcher) Prepare(ctx context.Context) error {
	var err error
	var depositAddress string
	var watchedQuotes []quote.WatchedPeginQuote
	watcher.currentBlock = big.NewInt(0)
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
	watcher.ticker = time.NewTicker(peginDepositWatcherInterval)
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C:
			if height, err := watcher.rpc.Btc.GetHeight(); err == nil && height.Cmp(watcher.currentBlock) > 0 {
				watcher.checkQuotes(context.Background())
				watcher.currentBlock = height
			} else if err != nil {
				log.Error(peginBtcDepositWatcherLog("error getting Bitcoin chain height: %v", err))
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
	log.Debug(peginBtcDepositWatcherLog("shut down"))
}

func (watcher *PeginDepositAddressWatcher) handleAcceptedPeginQuote(event entities.Event) {
	parsedEvent, ok := event.(quote.AcceptedPeginQuoteEvent)
	quoteHash := parsedEvent.RetainedQuote.QuoteHash
	if !ok {
		log.Error(peginBtcDepositWatcherLog("trying to parse wrong event"))
		return
	}

	if _, alreadyHaveQuote := watcher.quotes[quoteHash]; alreadyHaveQuote {
		log.Info(peginBtcDepositWatcherLog("Quote %s is already watched", quoteHash))
		return
	}

	err := watcher.btcWallet.ImportAddress(parsedEvent.RetainedQuote.DepositAddress)
	if err != nil {
		log.Error(peginBtcDepositWatcherLog("error while importing deposit address (%s): %v", parsedEvent.RetainedQuote.DepositAddress, err))
		return
	}
	watcher.quotes[quoteHash] = quote.NewWatchedPeginQuote(parsedEvent.Quote, parsedEvent.RetainedQuote)
}

func (watcher *PeginDepositAddressWatcher) checkQuotes(ctx context.Context) {
	for _, watchedQuote := range watcher.quotes {
		watcher.handleQuote(ctx, watchedQuote)
	}
}

func (watcher *PeginDepositAddressWatcher) handleQuote(ctx context.Context, watchedQuote quote.WatchedPeginQuote) {
	var err error
	quoteHash := watchedQuote.RetainedQuote.QuoteHash

	if watchedQuote.RetainedQuote.State == quote.PeginStateWaitingForDeposit {
		if err = watcher.handleNotDepositedQuote(ctx, watchedQuote); err != nil {
			log.Error(peginBtcDepositWatcherLog(callForUserErrorTemplate, quoteHash, err))
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

func (watcher *PeginDepositAddressWatcher) handleNotDepositedQuote(ctx context.Context, watchedQuote quote.WatchedPeginQuote) error {
	var err error
	var block blockchain.BitcoinBlockInformation
	var txs []blockchain.BitcoinTransactionInformation
	var updatedQuote quote.WatchedPeginQuote

	depositAddress := watchedQuote.RetainedQuote.DepositAddress
	if txs, err = watcher.btcWallet.GetTransactions(depositAddress); err != nil {
		return err
	}

	for _, tx := range txs {
		log.Info(peginBtcDepositWatcherLog("Checking transaction %s for quote %s", tx.Hash, watchedQuote.RetainedQuote.QuoteHash))
		if block, err = watcher.rpc.Btc.GetTransactionBlockInfo(tx.Hash); err != nil {
			return err
		}
		onTime := watchedQuote.PeginQuote.ExpireTime().After(block.Time)
		correctAmount := tx.AmountToAddress(watchedQuote.RetainedQuote.DepositAddress).Cmp(watchedQuote.PeginQuote.Total()) >= 0
		if watchedQuote.RetainedQuote.State == quote.PeginStateWaitingForDeposit && onTime && correctAmount {
			if updatedQuote, err = watcher.updatePeginDepositUseCase.Run(ctx, watchedQuote, block, tx); err != nil {
				return err
			} else {
				watcher.quotes[watchedQuote.RetainedQuote.QuoteHash] = updatedQuote
				return nil
			}
		} else {
			watcher.logRejectReason(block, tx, watchedQuote)
		}
	}
	return nil
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
