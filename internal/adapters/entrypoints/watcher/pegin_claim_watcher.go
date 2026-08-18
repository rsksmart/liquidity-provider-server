package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	log "github.com/sirupsen/logrus"
)

type PegInClaimRunner interface {
	Run(ctx context.Context, entry rootstock.PegInAddressRegistryWatchEntry, depositTxID string) error
	ReconcileSubmitting(ctx context.Context) error
}

type PegInClaimWatcher struct {
	runner             PegInClaimRunner
	watches            rootstock.PegInAddressRegistryWatchRepository
	btcWallet          blockchain.BitcoinWallet
	ticker             utils.Ticker
	watcherStopChannel chan bool
}

func NewPegInClaimWatcher(
	runner PegInClaimRunner,
	watches rootstock.PegInAddressRegistryWatchRepository,
	btcWallet blockchain.BitcoinWallet,
	ticker utils.Ticker,
) *PegInClaimWatcher {
	return &PegInClaimWatcher{
		runner:             runner,
		watches:            watches,
		btcWallet:          btcWallet,
		ticker:             ticker,
		watcherStopChannel: make(chan bool, 1),
	}
}

func (watcher *PegInClaimWatcher) Prepare(ctx context.Context) error {
	return watcher.runner.ReconcileSubmitting(ctx)
}

func (watcher *PegInClaimWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.check(context.Background())
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PegInClaimWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- true
	closeChannel <- true
	log.Debug(LogPegInClaimShutdown)
}

func (watcher *PegInClaimWatcher) check(ctx context.Context) {
	entries, err := watcher.watches.List(ctx)
	if err != nil {
		log.Errorf(LogPegInClaimListError, err)
		return
	}
	for _, entry := range entries {
		watcher.checkEntry(ctx, entry)
	}
}

func (watcher *PegInClaimWatcher) checkEntry(
	ctx context.Context,
	entry rootstock.PegInAddressRegistryWatchEntry,
) {
	if entry.State != rootstock.PegInAddressRegistryWatchImported {
		return
	}
	txs, err := watcher.btcWallet.GetTransactions(entry.BtcAddress)
	if err != nil {
		log.Error(LogPegInClaimWalletError(entry.BtcAddress, err))
		return
	}
	for _, tx := range txs {
		if tx.FirstOutputToAddress(entry.BtcAddress).Cmp(entities.NewWei(0)) <= 0 {
			continue
		}
		if err = watcher.runner.Run(ctx, entry, tx.Hash); err != nil {
			log.Error(LogPegInClaimRunError(entry.RskAddress, tx.Hash, err))
		}
	}
}
