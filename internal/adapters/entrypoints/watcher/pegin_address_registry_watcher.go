package watcher

import (
	"context"
	"fmt"
	"sync"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
)

type PegInWatcher struct {
	replayUseCase      *w.ReplayRegisteredAddressesUseCase
	discoverUseCase    *w.DiscoverRegisteredAddressUseCase
	getPendingUseCase  *w.GetPendingRegisteredAddressImportsUseCase
	finalizeUseCase    *w.FinalizeRegisteredAddressImportUseCase
	btcNetwork         blockchain.BitcoinNetwork
	wallet             blockchain.BitcoinWallet
	eventBus           entities.EventBus
	ticker             utils.Ticker
	checkpoint         *rootstock.PegInWatchCheckpoint
	scanMutex          sync.Mutex
	watcherStopChannel chan struct{}
}

func NewPegInWatcher(
	replayUseCase *w.ReplayRegisteredAddressesUseCase,
	discoverUseCase *w.DiscoverRegisteredAddressUseCase,
	getPendingUseCase *w.GetPendingRegisteredAddressImportsUseCase,
	finalizeUseCase *w.FinalizeRegisteredAddressImportUseCase,
	btcNetwork blockchain.BitcoinNetwork,
	wallet blockchain.BitcoinWallet,
	eventBus entities.EventBus,
	ticker utils.Ticker,
) *PegInWatcher {
	return &PegInWatcher{
		replayUseCase:      replayUseCase,
		discoverUseCase:    discoverUseCase,
		getPendingUseCase:  getPendingUseCase,
		finalizeUseCase:    finalizeUseCase,
		btcNetwork:         btcNetwork,
		wallet:             wallet,
		eventBus:           eventBus,
		ticker:             ticker,
		watcherStopChannel: make(chan struct{}, 1),
	}
}

func (watcher *PegInWatcher) Prepare(ctx context.Context) error {
	watcher.scanMutex.Lock()
	defer watcher.scanMutex.Unlock()
	return watcher.scanWithCheckpoint(ctx, nil)
}

func (watcher *PegInWatcher) Start() {
	var reorgEvents <-chan entities.Event
	if watcher.eventBus != nil {
		reorgEvents = watcher.eventBus.Subscribe(blockchain.NodeReorgCheckEventId)
	}
	for {
		select {
		case <-watcher.ticker.C():
			watcher.scanAndLog()
		case event, open := <-reorgEvents:
			if !open {
				reorgEvents = nil
				continue
			}
			watcher.resyncAfterReorg(event)
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			return
		}
	}
}

func (watcher *PegInWatcher) scanAndLog() {
	if err := watcher.scan(context.Background()); err != nil {
		log.Errorf("PegIn address registry watcher scan failed: %v", err)
	}
}

func (watcher *PegInWatcher) resyncAfterReorg(event entities.Event) {
	reorg, ok := event.(blockchain.NodeReorgCheckEvent)
	if !ok || reorg.NodeType != entities.NodeTypeRootstock || reorg.CurrentDepth == 0 {
		return
	}
	watcher.scanAndLog()
}

func (watcher *PegInWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("PegInWatcher shut down")
}

func (watcher *PegInWatcher) scan(ctx context.Context) error {
	watcher.scanMutex.Lock()
	defer watcher.scanMutex.Unlock()
	return watcher.scanWithCheckpoint(ctx, watcher.checkpoint)
}

func (watcher *PegInWatcher) scanWithCheckpoint(
	ctx context.Context,
	verifiedCheckpoint *rootstock.PegInWatchCheckpoint,
) error {
	result, err := watcher.replayUseCase.Run(ctx, verifiedCheckpoint)
	watcher.publishReplayOutcome(result, err)
	if err != nil {
		return err
	}
	watcher.checkpoint = result.Checkpoint
	pendingImports, err := watcher.getPendingUseCase.Run(ctx)
	if err != nil {
		return err
	}
	watchesNeedingRescan := make([]*rootstock.PegInWatch, 0, len(pendingImports))
	for _, watch := range pendingImports {
		discoverResult, discoverErr := watcher.discoverUseCase.Run(
			ctx,
			blockchain.NewAddressRegisteredFromWatchEntry(watch),
		)
		if discoverErr != nil {
			return discoverErr
		}
		if discoverResult.NeedsRescan {
			watchesNeedingRescan = append(watchesNeedingRescan, discoverResult.Watch)
		}
	}
	return watcher.finalizeUseCase.Run(
		ctx,
		watchesNeedingRescan,
		watcher.rescanPending(watchesNeedingRescan),
	)
}

func (watcher *PegInWatcher) publishReplayOutcome(result w.ReplayResult, err error) {
	if err != nil {
		log.Errorf("PegIn address registry replay failed: %v", err)
		return
	}
	if watcher.eventBus == nil {
		return
	}
	switch result.Reason {
	case blockchain.PegInAddressRegistryRecoveryRootMismatch:
		mismatch := result.Mismatch
		watcher.eventBus.Publish(blockchain.PegInAddressRegistryRootMismatchEvent{
			BaseEvent:   entities.NewBaseEvent(blockchain.PegInAddressRegistryRootMismatchEventId),
			BlockNumber: mismatch.BlockNumber,
			LocalRoot:   mismatch.LocalRoot,
			ChainRoot:   mismatch.ChainRoot,
		})
	case blockchain.PegInAddressRegistryRecoveryCatchUp:
	default:
		// An empty reason means the roots already matched, so there is nothing to announce.
		return
	}
	watcher.eventBus.Publish(blockchain.PegInAddressRegistryResyncStartedEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.PegInAddressRegistryResyncStartedEventId),
		Reason:    string(result.Reason),
	})
}

func (watcher *PegInWatcher) rescanPending(
	pending []*rootstock.PegInWatch,
) error {
	if len(pending) == 0 {
		return nil
	}
	tip, err := watcher.btcNetwork.GetHeight()
	if err != nil {
		return fmt.Errorf("get BTC height for registry rescan: %w", err)
	}
	fromHeight := max(tip.Int64()-peginAddressRegistryRescanDepthBlocks, 0)
	if _, err = watcher.wallet.RescanBlockchain(fromHeight); err != nil {
		return fmt.Errorf("rescan PegIn addresses: %w", err)
	}
	return nil
}
