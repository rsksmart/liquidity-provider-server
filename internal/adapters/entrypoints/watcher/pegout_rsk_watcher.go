package watcher

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
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
	expiredUseCase               *pegout.ExpiredPegoutQuoteUseCase
	sendPegoutUseCase            *pegout.SendPegoutUseCase
	updatePegoutDepositUseCase   *w.UpdatePegoutQuoteDepositUseCase
	initDepositCacheUseCase      *pegout.InitPegoutDepositCacheUseCase
	pegoutLp                     liquidity_provider.PegoutLiquidityProvider
	rpc                          blockchain.Rpc
	contracts                    blockchain.RskContracts
	ticker                       utils.Ticker
	eventBus                     entities.EventBus
	watcherStopChannel           chan bool
	currentBlock                 uint64
	cacheStartBlock              uint64
	currentBlockMutex            sync.RWMutex
	depositCheckTimeout          time.Duration
}

type PegoutRskDepositWatcherUseCases struct {
	getWatchedPegoutQuoteUseCase *w.GetWatchedPegoutQuoteUseCase
	expiredUseCase               *pegout.ExpiredPegoutQuoteUseCase
	sendPegoutUseCase            *pegout.SendPegoutUseCase
	updatePegoutDepositUseCase   *w.UpdatePegoutQuoteDepositUseCase
	initDepositCacheUseCase      *pegout.InitPegoutDepositCacheUseCase
}

func NewPegoutRskDepositWatcherUseCases(
	getWatchedPegoutQuoteUseCase *w.GetWatchedPegoutQuoteUseCase,
	expiredUseCase *pegout.ExpiredPegoutQuoteUseCase,
	sendPegoutUseCase *pegout.SendPegoutUseCase,
	updatePegoutDepositUseCase *w.UpdatePegoutQuoteDepositUseCase,
	initDepositCacheUseCase *pegout.InitPegoutDepositCacheUseCase,
) *PegoutRskDepositWatcherUseCases {
	return &PegoutRskDepositWatcherUseCases{
		getWatchedPegoutQuoteUseCase: getWatchedPegoutQuoteUseCase,
		expiredUseCase:               expiredUseCase,
		sendPegoutUseCase:            sendPegoutUseCase,
		updatePegoutDepositUseCase:   updatePegoutDepositUseCase,
		initDepositCacheUseCase:      initDepositCacheUseCase,
	}
}

func NewPegoutRskDepositWatcher(
	watcherUseCases *PegoutRskDepositWatcherUseCases,
	pegoutLp liquidity_provider.PegoutLiquidityProvider,
	rpc blockchain.Rpc,
	contracts blockchain.RskContracts,
	eventBus entities.EventBus,
	cacheStartBlock uint64,
	ticker utils.Ticker,
	depositCheckTimeout time.Duration,
) *PegoutRskDepositWatcher {
	quotes := make(map[string]quote.WatchedPegoutQuote)
	watcherStopChannel := make(chan bool, 1)
	currentBlock := cacheStartBlock
	return &PegoutRskDepositWatcher{
		quotes:                       quotes,
		getWatchedPegoutQuoteUseCase: watcherUseCases.getWatchedPegoutQuoteUseCase,
		expiredUseCase:               watcherUseCases.expiredUseCase,
		sendPegoutUseCase:            watcherUseCases.sendPegoutUseCase,
		updatePegoutDepositUseCase:   watcherUseCases.updatePegoutDepositUseCase,
		initDepositCacheUseCase:      watcherUseCases.initDepositCacheUseCase,
		pegoutLp:                     pegoutLp,
		rpc:                          rpc,
		contracts:                    contracts,
		eventBus:                     eventBus,
		watcherStopChannel:           watcherStopChannel,
		currentBlock:                 currentBlock,
		cacheStartBlock:              cacheStartBlock,
		ticker:                       ticker,
		currentBlockMutex:            sync.RWMutex{},
		quotesMutex:                  sync.RWMutex{},
		depositCheckTimeout:          depositCheckTimeout,
	}
}

func (watcher *PegoutRskDepositWatcher) Prepare(ctx context.Context) error {
	var quoteCreationBlock, height uint64
	var err error

	if watcher.cacheStartBlock != 0 {
		if err = watcher.initDepositCacheUseCase.Run(ctx, watcher.cacheStartBlock); err != nil {
			return err
		}
	} else {
		if height, err = watcher.rpc.Rsk.GetHeight(ctx); err != nil {
			return err
		}
		watcher.currentBlockMutex.Lock()
		watcher.currentBlock = height
		watcher.cacheStartBlock = height
		watcher.currentBlockMutex.Unlock()
	}

	watchedQuotes, err := watcher.getWatchedPegoutQuoteUseCase.Run(
		ctx,
		quote.PegoutStateWaitingForDeposit,
		quote.PegoutStateWaitingForDepositConfirmations,
		quote.PegoutStateClaimed,
	)
	if err != nil {
		return err
	}
	pegoutConfig := watcher.pegoutLp.PegoutConfiguration(ctx)
	watcher.currentBlockMutex.Lock()
	defer watcher.currentBlockMutex.Unlock()
	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	for _, watchedQuote := range watchedQuotes {
		quoteCreationBlock = quote.GetCreationBlock(pegoutConfig, watchedQuote.PegoutQuote)
		if watcher.currentBlock == 0 || watcher.currentBlock > quoteCreationBlock {
			watcher.currentBlock = quoteCreationBlock
		}
		watcher.quotes[watchedQuote.RetainedQuote.QuoteHash] = watchedQuote
	}

	log.Info(LogPegoutRskStart(watcher.currentBlock))
	return nil
}

func (watcher *PegoutRskDepositWatcher) Start() {
	acceptedEventChannel := watcher.eventBus.Subscribe(quote.AcceptedPegoutQuoteEventId)
	claimedEventChannel := watcher.eventBus.Subscribe(quote.ClaimedPegoutQuoteEventId)

watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.onTick()
		case event := <-acceptedEventChannel:
			if event != nil {
				watcher.handleAcceptedPegoutQuote(event)
			}
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

	checkContext, checkCancel := context.WithTimeout(context.Background(), watcher.depositCheckTimeout)
	defer checkCancel()
	height, err := watcher.rpc.Rsk.GetHeight(checkContext)
	if err != nil {
		log.Errorf(LogPegoutRskChainHeight, err)
		return
	}
	if height > watcher.currentBlock {
		watcher.checkDeposits(checkContext, watcher.currentBlock, height)
		watcher.checkQuotes(checkContext, height)
		watcher.currentBlock = height
	}
}

func (watcher *PegoutRskDepositWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug(LogPegoutRskShutdown)
}

func (watcher *PegoutRskDepositWatcher) handleAcceptedPegoutQuote(event entities.Event) {
	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	parsedEvent, ok := event.(quote.AcceptedPegoutQuoteEvent)
	quoteHash := parsedEvent.RetainedQuote.QuoteHash
	if !ok {
		log.Error(LogPegoutRskWrongEvent)
		return
	}

	if _, alreadyHaveQuote := watcher.quotes[quoteHash]; alreadyHaveQuote {
		log.Info(LogPegoutRskAlreadyWatched(quoteHash))
		return
	}
	watcher.quotes[quoteHash] = quote.NewWatchedPegoutQuote(parsedEvent.Quote, parsedEvent.RetainedQuote, parsedEvent.CreationData)
}

func (watcher *PegoutRskDepositWatcher) handleClaimedPegoutQuote(event entities.Event) {
	watcher.quotesMutex.Lock()
	defer watcher.quotesMutex.Unlock()
	parsedEvent, ok := event.(quote.ClaimedPegoutQuoteEvent)
	quoteHash := parsedEvent.RetainedQuote.QuoteHash
	if !ok {
		log.Error(LogPegoutRskWrongEvent)
		return
	}

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

func (watcher *PegoutRskDepositWatcher) checkDeposits(ctx context.Context, fromBlock, toBlock uint64) {
	var err error
	var deposits []quote.PegoutDeposit

	deposits, err = watcher.contracts.PegOut.GetDepositEvents(ctx, fromBlock, &toBlock)
	if err != nil {
		log.Error(LogPegoutRskGetDepositsError(fromBlock, toBlock, err))
		return
	}
	for _, deposit := range deposits {
		log.Info(LogPegoutRskCheckingDeposit(deposit.TxHash, deposit.QuoteHash))
		watcher.checkDeposit(ctx, deposit)
	}
}

func (watcher *PegoutRskDepositWatcher) checkDeposit(ctx context.Context, deposit quote.PegoutDeposit) {
	var newWatchedQuote quote.WatchedPegoutQuote
	var err error
	watchedQuote, ok := watcher.quotes[deposit.QuoteHash]
	if ok && deposit.IsValidForQuote(watchedQuote.PegoutQuote) && watchedQuote.RetainedQuote.State == quote.PegoutStateWaitingForDeposit {
		if newWatchedQuote, err = watcher.updatePegoutDepositUseCase.Run(ctx, watchedQuote, deposit); err == nil {
			watcher.quotes[deposit.QuoteHash] = newWatchedQuote
		} else {
			log.Error(LogPegoutRskUpdateDepositError(watchedQuote.RetainedQuote.QuoteHash, err))
		}
	}
	if ok && !deposit.IsValidForQuote(watchedQuote.PegoutQuote) {
		watcher.logRejectReason(deposit, watchedQuote)
	}
}

func (watcher *PegoutRskDepositWatcher) checkQuotes(ctx context.Context, height uint64) {
	for _, watchedQuote := range watcher.quotes {
		watcher.checkQuote(ctx, height, watchedQuote)
	}
}

func (watcher *PegoutRskDepositWatcher) checkQuote(ctx context.Context, height uint64, watchedQuote quote.WatchedPegoutQuote) {
	var err error
	var receipt blockchain.TransactionReceipt
	if watchedQuote.RetainedQuote.State == quote.PegoutStateWaitingForDeposit && watchedQuote.PegoutQuote.IsExpired() {
		if err = watcher.expiredUseCase.Run(ctx, watchedQuote.RetainedQuote); err != nil {
			log.Error(LogPegoutRskUpdateExpiredError(watchedQuote.RetainedQuote.QuoteHash, err))
			return
		} else {
			log.Info(LogPegoutRskExpired(watchedQuote.RetainedQuote.QuoteHash, watchedQuote.PegoutQuote.ExpireTime().Unix()))
			delete(watcher.quotes, watchedQuote.RetainedQuote.QuoteHash)
		}
	}

	if watchedQuote.RetainedQuote.State == quote.PegoutStateClaimed {
		watcher.sendPegout(ctx, watchedQuote)
		return
	}

	if watchedQuote.RetainedQuote.State == quote.PegoutStateWaitingForDepositConfirmations {
		if receipt, err = watcher.rpc.Rsk.GetTransactionReceipt(ctx, watchedQuote.RetainedQuote.UserRskTxHash); err != nil {
			log.Error(LogPegoutRskReceiptError(watchedQuote.RetainedQuote.QuoteHash, err))
			return
		}
		if validateDepositedPegoutQuote(watchedQuote, receipt, height) {
			watcher.sendPegout(ctx, watchedQuote)
		}
	}
}

func (watcher *PegoutRskDepositWatcher) sendPegout(ctx context.Context, watchedQuote quote.WatchedPegoutQuote) {
	var err error
	if err = watcher.sendPegoutUseCase.Run(ctx, watchedQuote.RetainedQuote); errors.Is(err, usecases.NonRecoverableError) {
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

func validateDepositedPegoutQuote(watchedQuote quote.WatchedPegoutQuote, receipt blockchain.TransactionReceipt, height uint64) bool {
	event, err := rootstock.ParseDepositEventByQuoteHash(receipt, watchedQuote.RetainedQuote.QuoteHash, watchedQuote.PegoutQuote.LbcAddress)
	if err != nil {
		log.Error(LogPegoutRskParseError(watchedQuote.RetainedQuote.QuoteHash, err))
		return false
	}
	return receipt.BlockNumber+uint64(watchedQuote.PegoutQuote.DepositConfirmations) < height &&
		watchedQuote.RetainedQuote.State == quote.PegoutStateWaitingForDepositConfirmations &&
		event.Log.Amount.Cmp(watchedQuote.PegoutQuote.Total()) >= 0
}

func (watcher *PegoutRskDepositWatcher) logRejectReason(deposit quote.PegoutDeposit, watchedQuote quote.WatchedPegoutQuote) {
	rejectReason := LogPegoutRskRejectReason(watchedQuote.RetainedQuote.QuoteHash)
	if deposit.Timestamp.After(watchedQuote.PegoutQuote.ExpireTime()) {
		depositTime := deposit.Timestamp.Unix()
		expirationTime := watchedQuote.PegoutQuote.ExpireTime().Unix()
		rejectReason += LogPegoutRskRejectExpired(expirationTime, depositTime-expirationTime, depositTime)
	}
	paidAmount := deposit.Amount
	expectedAmount := watchedQuote.PegoutQuote.Total()
	if paidAmount.Cmp(expectedAmount) < 0 {
		rejectReason += LogPegoutRskRejectAmount(paidAmount.String(), expectedAmount.String())
	}
	log.Info(LogPegoutRskPrefix + rejectReason)
}
