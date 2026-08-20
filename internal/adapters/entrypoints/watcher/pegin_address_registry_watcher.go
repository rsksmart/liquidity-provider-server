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

type PegInAddressRegistryWatcher struct {
	replayUseCase      *w.ReplayRegisteredAddressesUseCase
	finalizeUseCase    *w.FinalizeRegisteredAddressImportUseCase
	btcNetwork         blockchain.BitcoinNetwork
	wallet             blockchain.BitcoinWallet
	eventBus           entities.EventBus
	ticker             utils.Ticker
	startBlock         uint64
	pageSize           uint64
	scanMutex          sync.Mutex
	watcherStopChannel chan struct{}
}

func NewPegInAddressRegistryWatcher(
	replayUseCase *w.ReplayRegisteredAddressesUseCase,
	finalizeUseCase *w.FinalizeRegisteredAddressImportUseCase,
	btcNetwork blockchain.BitcoinNetwork,
	wallet blockchain.BitcoinWallet,
	eventBus entities.EventBus,
	ticker utils.Ticker,
	startBlock uint64,
	pageSize uint64,
) *PegInAddressRegistryWatcher {
	return &PegInAddressRegistryWatcher{
		replayUseCase:      replayUseCase,
		finalizeUseCase:    finalizeUseCase,
		btcNetwork:         btcNetwork,
		wallet:             wallet,
		eventBus:           eventBus,
		ticker:             ticker,
		startBlock:         startBlock,
		pageSize:           pageSize,
		watcherStopChannel: make(chan struct{}, 1),
	}
}

func (watcher *PegInAddressRegistryWatcher) Prepare(context.Context) error {
	return nil
}

func (watcher *PegInAddressRegistryWatcher) Start() {
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

func (watcher *PegInAddressRegistryWatcher) scanAndLog() {
	if err := watcher.scan(context.Background()); err != nil {
		log.Errorf("PegIn address registry watcher scan failed: %v", err)
	}
}

func (watcher *PegInAddressRegistryWatcher) resyncAfterReorg(event entities.Event) {
	reorg, ok := event.(blockchain.NodeReorgCheckEvent)
	if !ok || reorg.NodeType != entities.NodeTypeRootstock || reorg.CurrentDepth == 0 {
		return
	}
	if err := watcher.resync(context.Background()); err != nil {
		log.Errorf("PegIn address registry watcher scan failed: %v", err)
	}
}

func (watcher *PegInAddressRegistryWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("PegInAddressRegistryWatcher shut down")
}

func (watcher *PegInAddressRegistryWatcher) scan(ctx context.Context) error {
	return watcher.runIntegrity(ctx, false)
}

func (watcher *PegInAddressRegistryWatcher) resync(ctx context.Context) error {
	return watcher.runIntegrity(ctx, true)
}

func (watcher *PegInAddressRegistryWatcher) runIntegrity(ctx context.Context, discardCheckpoint bool) error {
	watcher.scanMutex.Lock()
	defer watcher.scanMutex.Unlock()

	pending, err := watcher.replayUseCase.Run(ctx, discardCheckpoint, watcher.startBlock, watcher.pageSize)
	if err != nil {
		return err
	}
	return watcher.finalizeUseCase.Run(ctx, pending, watcher.rescanPending(pending))
}

func (watcher *PegInAddressRegistryWatcher) rescanPending(
	pending []*rootstock.PegInAddressRegistryWatchEntry,
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
