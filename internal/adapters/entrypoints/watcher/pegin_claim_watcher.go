package watcher

import (
	"context"
	"errors"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type PegInClaimRunner interface {
	Run(ctx context.Context, entry rootstock.PegInWatch, depositTxID string) error
}

type PegInClaimSettler interface {
	Run(ctx context.Context, claim rootstock.PegInClaim) error
}

type PegInClaimsGetter interface {
	Run(ctx context.Context, states ...rootstock.PegInClaimState) ([]rootstock.PegInClaim, error)
}

type PegInWatchesGetter interface {
	Run(ctx context.Context, states ...rootstock.PegInWatchState) ([]rootstock.PegInWatch, error)
}

type PegInClaimWatcherUseCases struct {
	claim      PegInClaimRunner
	settle     PegInClaimSettler
	getClaims  PegInClaimsGetter
	getWatches PegInWatchesGetter
}

func NewPegInClaimWatcherUseCases(
	claim PegInClaimRunner,
	settle PegInClaimSettler,
	getClaims PegInClaimsGetter,
	getWatches PegInWatchesGetter,
) *PegInClaimWatcherUseCases {
	return &PegInClaimWatcherUseCases{
		claim:      claim,
		settle:     settle,
		getClaims:  getClaims,
		getWatches: getWatches,
	}
}

type PegInClaimWatcher struct {
	useCases           *PegInClaimWatcherUseCases
	btcWallet          blockchain.BitcoinWallet
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
}

func NewPegInClaimWatcher(
	useCases *PegInClaimWatcherUseCases,
	btcWallet blockchain.BitcoinWallet,
	ticker utils.Ticker,
) *PegInClaimWatcher {
	return &PegInClaimWatcher{
		useCases:           useCases,
		btcWallet:          btcWallet,
		ticker:             ticker,
		watcherStopChannel: make(chan struct{}, 1),
	}
}

func (watcher *PegInClaimWatcher) Prepare(ctx context.Context) error {
	watcher.settleSubmitting(ctx)
	return nil
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
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug(LogPegInClaimShutdown)
}

func (watcher *PegInClaimWatcher) check(ctx context.Context) {
	if watcher.settleSubmitting(ctx) {
		return
	}
	entries, err := watcher.useCases.getWatches.Run(ctx, rootstock.PegInWatchImported)
	if err != nil {
		log.Errorf(LogPegInClaimListError, err)
		return
	}
	for _, entry := range entries {
		watcher.checkEntry(ctx, entry)
	}
}

func (watcher *PegInClaimWatcher) settleSubmitting(ctx context.Context) (aborted bool) {
	claims, err := watcher.useCases.getClaims.Run(ctx, rootstock.PegInClaimSubmitting)
	if err != nil {
		log.Error(LogPegInClaimSettleListError(err))
		return true
	}
	for _, claim := range claims {
		err = watcher.useCases.settle.Run(ctx, claim)
		if errors.Is(err, usecases.InfrastructureUnavailableError) {
			log.Error(LogPegInClaimSettleAborted(claim.RskAddress, claim.DepositTxID, err))
			return true
		}
		if err != nil {
			log.Error(LogPegInClaimSettleError(claim.RskAddress, claim.DepositTxID, err))
		}
	}
	return false
}

func (watcher *PegInClaimWatcher) checkEntry(
	ctx context.Context,
	entry rootstock.PegInWatch,
) {
	txs, err := watcher.btcWallet.GetTransactions(entry.BtcAddress)
	if err != nil {
		log.Error(LogPegInClaimWalletError(entry.BtcAddress, err))
		return
	}
	for _, tx := range txs {
		if tx.FirstOutputToAddress(entry.BtcAddress).Cmp(entities.NewWei(0)) > 0 {
			if err = watcher.useCases.claim.Run(ctx, entry, tx.Hash); err != nil {
				log.Error(LogPegInClaimRunError(entry.RskAddress, tx.Hash, err))
			}
		}
	}
}
