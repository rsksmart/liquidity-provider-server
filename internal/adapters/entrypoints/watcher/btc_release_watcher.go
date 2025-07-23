package watcher

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type UpdateBtcReleaseUseCase interface {
	Run(ctx context.Context, batch rootstock.BatchPegOut) (uint, error)
}

type BtcReleaseWatcher struct {
	contracts              blockchain.RskContracts
	rpc                    blockchain.Rpc
	updateRebalanceUseCase UpdateBtcReleaseUseCase
	ticker                 utils.Ticker
	watcherStopChannel     chan struct{}
	startBlock             uint64
	pageSize               uint64
	currentBlock           uint64
	btcReleaseCheckTimeout time.Duration
	currentBlockMutex      sync.RWMutex
}

func NewBtcReleaseWatcher(
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	updateRebalanceUseCase UpdateBtcReleaseUseCase,
	ticker utils.Ticker,
	startBlock uint64,
	pageSize uint64,
	btcReleaseCheckTimeout time.Duration,
) *BtcReleaseWatcher {
	const defaultPageSize = 2000
	watcherStopChannel := make(chan struct{}, 1)
	if pageSize == 0 {
		pageSize = defaultPageSize
	}
	return &BtcReleaseWatcher{
		contracts:              contracts,
		rpc:                    rpc,
		updateRebalanceUseCase: updateRebalanceUseCase,
		ticker:                 ticker,
		watcherStopChannel:     watcherStopChannel,
		startBlock:             startBlock,
		pageSize:               pageSize,
		btcReleaseCheckTimeout: btcReleaseCheckTimeout,
		currentBlockMutex:      sync.RWMutex{},
	}
}

func (watcher *BtcReleaseWatcher) Prepare(ctx context.Context) error {
	watcher.currentBlockMutex.Lock()
	defer watcher.currentBlockMutex.Unlock()
	if watcher.startBlock != 0 {
		watcher.currentBlock = watcher.startBlock
		return nil
	}
	height, err := watcher.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return err
	}
	watcher.currentBlock = height
	return nil
}

func (watcher *BtcReleaseWatcher) Start() {
	var err error
	var newCurrent uint64
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.currentBlockMutex.Lock()
			newCurrent, err = watcher.checkBatchPegOuts()
			if err == nil {
				watcher.currentBlock = newCurrent
			} else {
				log.Errorf("Error checking BatchPegOuts in BtcReleaseWatcher: %v", err)
			}
			watcher.currentBlockMutex.Unlock()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *BtcReleaseWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug("BtcReleaseWatcher shut down")
}

func (watcher *BtcReleaseWatcher) checkBatchPegOuts() (uint64, error) {
	var toBlock uint64
	var err error
	checkContext, checkCancel := context.WithTimeout(context.Background(), watcher.btcReleaseCheckTimeout)
	defer checkCancel()

	if toBlock, err = watcher.nextBlock(checkContext); err != nil {
		return 0, err
	} else if toBlock == watcher.currentBlock {
		return toBlock, nil
	}

	batches, err := watcher.contracts.Bridge.GetBatchPegOutCreatedEvent(checkContext, watcher.currentBlock, &toBlock)
	if err != nil {
		return 0, fmt.Errorf("error fetching BatchPegOutCreated events in BtcReleaseWatcher: %w", err)
	}

	log.Infof("Checking BatchPegOuts from block %d to %d, found %d batches", watcher.currentBlock, toBlock, len(batches))
	err = watcher.updateBatches(checkContext, batches)
	if err != nil {
		return 0, err
	}

	return toBlock, nil
}

func (watcher *BtcReleaseWatcher) nextBlock(ctx context.Context) (uint64, error) {
	height, err := watcher.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting RSK height in BtcReleaseWatcher: %w", err)
	}

	var toBlock uint64
	if (watcher.currentBlock + watcher.pageSize) < height {
		toBlock = watcher.currentBlock + watcher.pageSize
	} else {
		toBlock = height
	}
	return toBlock, nil
}

func (watcher *BtcReleaseWatcher) updateBatches(ctx context.Context, batches []rootstock.BatchPegOut) error {
	var updated uint
	var err error
	for _, batch := range batches {
		log.Debugf("Processing BatchPegOut: %+v", batch)
		if updated, err = watcher.updateRebalanceUseCase.Run(ctx, batch); err != nil {
			return fmt.Errorf("error processing BatchPegOut: %w", err)
		} else if updated == 0 {
			log.Infof("No PegOuts to process in batch (%s)", batch.TransactionHash)
		} else {
			log.Infof("Successfully processed %d quotes in BatchPegOut (%s)", updated, batch.TransactionHash)
		}
	}
	return nil
}

func (watcher *BtcReleaseWatcher) CurrentBlock() uint64 {
	watcher.currentBlockMutex.RLock()
	defer watcher.currentBlockMutex.RUnlock()
	return watcher.currentBlock
}
