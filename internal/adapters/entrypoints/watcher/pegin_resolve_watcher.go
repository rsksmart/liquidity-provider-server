package watcher

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	w "github.com/rsksmart/liquidity-provider-server/internal/usecases/watcher"
	log "github.com/sirupsen/logrus"
)

type ResolvePegInUseCase interface {
	Run(ctx context.Context, claim rootstock.PegInClaim) error
}

type PegInResolveWatcher struct {
	getWatchedUseCase  *w.GetWatchedPegInClaimsUseCase
	resolveUseCase     ResolvePegInUseCase
	rpc                blockchain.Rpc
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
	checkedHeight      uint64
}

func NewPegInResolveWatcher(
	getWatchedUseCase *w.GetWatchedPegInClaimsUseCase,
	resolveUseCase ResolvePegInUseCase,
	rpc blockchain.Rpc,
	ticker utils.Ticker,
) *PegInResolveWatcher {
	return &PegInResolveWatcher{
		getWatchedUseCase:  getWatchedUseCase,
		resolveUseCase:     resolveUseCase,
		rpc:                rpc,
		ticker:             ticker,
		watcherStopChannel: make(chan struct{}, 1),
	}
}

func (watcher *PegInResolveWatcher) Prepare(ctx context.Context) error {
	return nil
}

func (watcher *PegInResolveWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.resolveOnNewBlock(context.Background())
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PegInResolveWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug(LogPegInResolveShutdown)
}

func (watcher *PegInResolveWatcher) resolveOnNewBlock(ctx context.Context) {
	height, err := watcher.rpc.Btc.GetHeight()
	if err != nil {
		log.Errorf(LogPegInResolveBtcHeightError, err)
		return
	}
	if height.Uint64() <= watcher.checkedHeight {
		return
	}
	claims, err := watcher.getWatchedUseCase.Run(ctx, rootstock.PegInClaimClaimed)
	if err != nil {
		log.Errorf(LogPegInResolveListError, err)
		return
	}
	for _, claim := range claims {
		watcher.resolve(ctx, claim)
	}
	watcher.checkedHeight = height.Uint64()
}

func (watcher *PegInResolveWatcher) resolve(ctx context.Context, claim rootstock.PegInClaim) {
	if err := watcher.resolveUseCase.Run(ctx, claim); err != nil {
		log.Error(LogPegInResolveError(claim.RskAddress, claim.DepositTxID, err))
	}
}
