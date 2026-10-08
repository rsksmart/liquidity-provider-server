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

type PegInClaimResolver interface {
	Run(ctx context.Context, claim rootstock.PegInClaim) error
}

// PegInResolveWatcher resolves the claimed peg-ins once their deposit has enough confirmations
type PegInResolveWatcher struct {
	resolveUseCase     PegInClaimResolver
	getClaimsUseCase   PegInClaimsGetter
	btcRpc             blockchain.BitcoinNetwork
	eventBus           entities.EventBus
	ticker             utils.Ticker
	watchedClaims      map[string]rootstock.PegInClaim
	lastHeight         uint64
	watcherStopChannel chan struct{}
}

func NewPegInResolveWatcher(
	resolveUseCase PegInClaimResolver,
	getClaimsUseCase PegInClaimsGetter,
	btcRpc blockchain.BitcoinNetwork,
	eventBus entities.EventBus,
	ticker utils.Ticker,
) *PegInResolveWatcher {
	return &PegInResolveWatcher{
		resolveUseCase:     resolveUseCase,
		getClaimsUseCase:   getClaimsUseCase,
		btcRpc:             btcRpc,
		eventBus:           eventBus,
		ticker:             ticker,
		watchedClaims:      make(map[string]rootstock.PegInClaim),
		watcherStopChannel: make(chan struct{}, 1),
	}
}

func (watcher *PegInResolveWatcher) Prepare(ctx context.Context) error {
	claims, err := watcher.getClaimsUseCase.Run(ctx, rootstock.PegInClaimClaimed)
	if err != nil {
		return err
	}
	for _, claim := range claims {
		watcher.addClaim(claim)
	}
	return nil
}

func (watcher *PegInResolveWatcher) Start() {
	eventChannel := watcher.eventBus.Subscribe(rootstock.PegInClaimCompletedEventId)
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.checkClaims()
		case event := <-eventChannel:
			if event != nil {
				watcher.handleClaimCompletedEvent(event)
			}
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

func (watcher *PegInResolveWatcher) handleClaimCompletedEvent(event entities.Event) {
	parsedEvent, ok := event.(rootstock.PegInClaimCompletedEvent)
	if !ok {
		log.Error(LogPegInResolveWrongEvent)
		return
	}
	if parsedEvent.Claim.IsClaimed() {
		watcher.addClaim(parsedEvent.Claim)
	}
}

// TODO: fix bug where settle marks a claim as claimed with an empty pegInID
func (watcher *PegInResolveWatcher) addClaim(claim rootstock.PegInClaim) {
	watcher.watchedClaims[claim.PegInID] = claim
}

func (watcher *PegInResolveWatcher) checkClaims() {
	height, err := watcher.btcRpc.GetHeight()
	if err != nil || height == nil {
		log.Errorf(LogPegInResolveBtcHeightError, err)
		return
	}
	if height.Uint64() <= watcher.lastHeight {
		return
	}
	if watcher.resolveClaims() {
		watcher.lastHeight = height.Uint64()
	}
}

func (watcher *PegInResolveWatcher) resolveClaims() bool {
	ctx := context.Background()
	for _, claim := range watcher.watchedClaims {
		if stop := watcher.resolveClaim(ctx, claim); stop {
			return false
		}
	}
	return true
}

func (watcher *PegInResolveWatcher) resolveClaim(ctx context.Context, claim rootstock.PegInClaim) (stop bool) {
	err := watcher.resolveUseCase.Run(ctx, claim)
	switch {
	case err == nil:
		delete(watcher.watchedClaims, claim.PegInID)
	case errors.Is(err, usecases.NonRecoverableError):
		log.Error(LogPegInResolveError(claim.RskAddress, claim.DepositTxID, err))
		delete(watcher.watchedClaims, claim.PegInID)
	case errors.Is(err, usecases.NoEnoughConfirmationsError):
		log.Debug(LogPegInResolveNotReady(claim.RskAddress, claim.DepositTxID))
	case errors.Is(err, blockchain.WaitingForBridgeError): // TODO: Handle this error in the contract side
		// the use case already logged it
	case errors.Is(err, usecases.InfrastructureUnavailableError):
		log.Error(LogPegInResolveAborted(claim.RskAddress, claim.DepositTxID, err))
		return true
	case errors.Is(err, blockchain.ContractPausedError):
		return true
	default:
		log.Error(LogPegInResolveError(claim.RskAddress, claim.DepositTxID, err))
	}
	return false
}
