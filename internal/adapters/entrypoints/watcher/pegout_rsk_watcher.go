package watcher

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
	"time"
)

type PegoutRskDepositWatcher struct {
	quotes                       map[string]w.WatchedPegoutQuote
	getWatchedPegoutQuoteUseCase *w.GetWatchedPegoutQuoteUseCase
	expiredUseCase               *pegout.ExpiredPegoutQuoteUseCase
	sendPegoutUseCase            *pegout.SendPegoutUseCase
	updatePegoutDepositUseCase   *w.UpdatePegoutQuoteDepositUseCase
	initDepositCacheUseCase      *pegout.InitPegoutDepositCacheUseCase
	pegoutLp                     liquidity_provider.PegoutLiquidityProvider
	rskRpc                       blockchain.RootstockRpcServer
	lbc                          blockchain.LiquidityBridgeContract
	ticker                       *time.Ticker
	eventBus                     entities.EventBus
	watcherStopChannel           chan bool
	currentBlock                 uint64
	cacheStartBlock              uint64
}

func NewPegoutRskDepositWatcher(
	getWatchedPegoutQuoteUseCase *w.GetWatchedPegoutQuoteUseCase,
	expiredUseCase *pegout.ExpiredPegoutQuoteUseCase,
	sendPegoutUseCase *pegout.SendPegoutUseCase,
	updatePegoutDepositUseCase *w.UpdatePegoutQuoteDepositUseCase,
	initDepositCacheUseCase *pegout.InitPegoutDepositCacheUseCase,
	pegoutLp liquidity_provider.PegoutLiquidityProvider,
	rskRpc blockchain.RootstockRpcServer,
	lbc blockchain.LiquidityBridgeContract,
	eventBus entities.EventBus,
	cacheStartBlock uint64,
) *PegoutRskDepositWatcher {
	quotes := make(map[string]w.WatchedPegoutQuote)
	watcherStopChannel := make(chan bool, 1)
	currentBlock := cacheStartBlock
	return &PegoutRskDepositWatcher{
		quotes:                       quotes,
		getWatchedPegoutQuoteUseCase: getWatchedPegoutQuoteUseCase,
		expiredUseCase:               expiredUseCase,
		sendPegoutUseCase:            sendPegoutUseCase,
		updatePegoutDepositUseCase:   updatePegoutDepositUseCase,
		initDepositCacheUseCase:      initDepositCacheUseCase,
		pegoutLp:                     pegoutLp,
		rskRpc:                       rskRpc,
		lbc:                          lbc,
		eventBus:                     eventBus,
		watcherStopChannel:           watcherStopChannel,
		currentBlock:                 currentBlock,
		cacheStartBlock:              cacheStartBlock,
	}
}

func (watcher *PegoutRskDepositWatcher) Prepare(ctx context.Context) error {
	var quoteCreationBlock uint64
	var err error

	if watcher.cacheStartBlock != 0 {
		if err = watcher.initDepositCacheUseCase.Run(ctx, watcher.cacheStartBlock); err != nil {
			return err
		}
	}

	watchedQuotes, err := watcher.getWatchedPegoutQuoteUseCase.Run(ctx, quote.PegoutStateWaitingForDeposit, quote.PegoutStateWaitingForDepositConfirmations)
	if err != nil {
		return err
	}
	pegoutConfig := watcher.pegoutLp.PegoutConfiguration(ctx)
	for _, watchedQuote := range watchedQuotes {
		quoteCreationBlock = quote.GetCreationBlock(pegoutConfig, watchedQuote.PegoutQuote)
		if watcher.currentBlock == 0 || watcher.currentBlock > quoteCreationBlock {
			watcher.currentBlock = quoteCreationBlock
		}
		watcher.quotes[watchedQuote.RetainedQuote.QuoteHash] = watchedQuote
	}

	log.Info("Starting to watch pegout deposits from block ", watcher.currentBlock)
	return nil
}

func (watcher *PegoutRskDepositWatcher) Start() {
	var checkContext context.Context
	var checkCancel context.CancelFunc
	eventChannel := watcher.eventBus.Subscribe(quote.AcceptedPegoutQuoteEventId)
	watcher.ticker = time.NewTicker(pegoutDepositWatcherInterval)

watcherLoop:
	for {
		select {
		case <-watcher.ticker.C:
			checkContext, checkCancel = context.WithTimeout(context.Background(), 1*time.Minute)
			if height, err := watcher.rskRpc.GetHeight(checkContext); err == nil && height > watcher.currentBlock {
				watcher.checkDeposits(checkContext, watcher.currentBlock, height)
				watcher.checkQuotes(checkContext, height)
				watcher.currentBlock = height
			} else if err != nil {
				log.Error("PegoutRskDepositWatcher: error getting RSK chain height: ", err)
			}
			checkCancel()
		case event := <-eventChannel:
			if event != nil {
				watcher.handleAcceptedPegoutQuote(event)
			}
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PegoutRskDepositWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug("PegoutRskDepositWatcher shut down")
}

func (watcher *PegoutRskDepositWatcher) handleAcceptedPegoutQuote(event entities.Event) {
	parsedEvent, ok := event.(quote.AcceptedPegoutQuoteEvent)
	quoteHash := parsedEvent.RetainedQuote.QuoteHash
	if !ok {
		log.Error("Trying to parse wrong event in Pegout Rsk deposit watcher")
		return
	}

	if _, alreadyHaveQuote := watcher.quotes[quoteHash]; alreadyHaveQuote {
		log.Infof("Quote %s is already watched\n", quoteHash)
		return
	}
	watcher.quotes[quoteHash] = w.NewWatchedPegoutQuote(parsedEvent.Quote, parsedEvent.RetainedQuote)
}

func (watcher *PegoutRskDepositWatcher) checkDeposits(ctx context.Context, fromBlock, toBlock uint64) {
	var err error
	var deposits []quote.PegoutDeposit

	deposits, err = watcher.lbc.GetDepositEvents(ctx, fromBlock, &toBlock)
	if err != nil {
		log.Errorf("Error executing getting deposits in range [%d, %d] in PegoutRskDepositWatcher\n", fromBlock, toBlock)
		return
	}
	for _, deposit := range deposits {
		watcher.checkDeposit(ctx, deposit)
	}
}

func (watcher *PegoutRskDepositWatcher) checkDeposit(ctx context.Context, deposit quote.PegoutDeposit) {
	var newWatchedQuote w.WatchedPegoutQuote
	var err error
	watchedQuote, ok := watcher.quotes[deposit.QuoteHash]
	if ok && deposit.IsValidForQuote(watchedQuote.PegoutQuote) && watchedQuote.RetainedQuote.State == quote.PegoutStateWaitingForDeposit {
		if newWatchedQuote, err = watcher.updatePegoutDepositUseCase.Run(ctx, watchedQuote, deposit); err == nil {
			watcher.quotes[deposit.QuoteHash] = newWatchedQuote
		} else {
			log.Errorf("Error updating pegout deposit quote (%s): %v\n", watchedQuote.RetainedQuote.QuoteHash, err)
		}
	}
}

func (watcher *PegoutRskDepositWatcher) checkQuotes(ctx context.Context, height uint64) {
	for _, watchedQuote := range watcher.quotes {
		watcher.checkQuote(ctx, height, watchedQuote)
	}
}

func (watcher *PegoutRskDepositWatcher) checkQuote(ctx context.Context, height uint64, watchedQuote w.WatchedPegoutQuote) {
	var err error
	var receipt blockchain.TransactionReceipt
	if watchedQuote.RetainedQuote.State == quote.PegoutStateWaitingForDeposit && watchedQuote.PegoutQuote.IsExpired() {
		if err = watcher.expiredUseCase.Run(ctx, watchedQuote.RetainedQuote); err != nil {
			log.Errorf("Error updating expired quote (%s): %v\n", watchedQuote.RetainedQuote.QuoteHash, err)
			return
		} else {
			delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
		}
	}

	if watchedQuote.RetainedQuote.State == quote.PegoutStateWaitingForDepositConfirmations {
		if receipt, err = watcher.rskRpc.GetTransactionReceipt(ctx, watchedQuote.RetainedQuote.UserRskTxHash); err != nil {
			log.Errorf("Error getting pegout deposit receipt of quote %s: %v\n", watchedQuote.RetainedQuote.QuoteHash, err)
			return
		}
		if validateDepositedPegoutQuote(watchedQuote, receipt, height) {
			watcher.sendPegout(ctx, watchedQuote)
		}
	}
}

func (watcher *PegoutRskDepositWatcher) sendPegout(ctx context.Context, watchedQuote w.WatchedPegoutQuote) {
	var err error
	if err = watcher.sendPegoutUseCase.Run(ctx, watchedQuote.RetainedQuote); errors.Is(err, usecases.NonRecoverableError) {
		log.Errorf("Error sending pegout to the user (quote %s): %v\n", watchedQuote.RetainedQuote.QuoteHash, err)
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
	} else if err != nil {
		log.Errorf("Error sending pegout to the user (quote %s): %v\n", watchedQuote.RetainedQuote.QuoteHash, err)
	} else {
		delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
	}
}

func validateDepositedPegoutQuote(watchedQuote w.WatchedPegoutQuote, receipt blockchain.TransactionReceipt, height uint64) bool {
	return receipt.BlockNumber+uint64(watchedQuote.PegoutQuote.DepositConfirmations) < height &&
		watchedQuote.RetainedQuote.State == quote.PegoutStateWaitingForDepositConfirmations &&
		!watchedQuote.PegoutQuote.IsExpired() &&
		receipt.Value.Cmp(watchedQuote.PegoutQuote.Total()) >= 0
}
