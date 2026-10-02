package watcher

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	log "github.com/sirupsen/logrus"
)

const defaultPegoutEscrowWatcherPageSize uint64 = 2000

type PegoutEscrowWatcher struct {
	contracts          blockchain.RskContracts
	rpc                blockchain.Rpc
	repository         blockchain.PegOutEscrowWatchRepository
	claimPegOutUseCase *pegout.ClaimPegOutUseCase
	ticker             utils.Ticker
	watcherStopChannel chan struct{}
	startBlock         uint64
	pageSize           uint64
	lastScannedBlock   uint64
	checkTimeout       time.Duration
	candidates         map[string]blockchain.PegOutRequested
	stateMutex         sync.RWMutex
}

type pegOutEscrowScanRange struct {
	from uint64
	to   uint64
	skip bool
}

type pegOutEscrowEvents struct {
	requested []blockchain.PegOutRequested
	claimed   []blockchain.PegOutClaimed
	cancelled []blockchain.PegOutCancelled
}

func NewPegoutEscrowWatcher(
	contracts blockchain.RskContracts,
	rpc blockchain.Rpc,
	repository blockchain.PegOutEscrowWatchRepository,
	claimPegOutUseCase *pegout.ClaimPegOutUseCase,
	ticker utils.Ticker,
	startBlock uint64,
	pageSize uint64,
	checkTimeout time.Duration,
) *PegoutEscrowWatcher {
	if pageSize == 0 {
		pageSize = defaultPegoutEscrowWatcherPageSize
	}
	return &PegoutEscrowWatcher{
		contracts:          contracts,
		rpc:                rpc,
		repository:         repository,
		claimPegOutUseCase: claimPegOutUseCase,
		ticker:             ticker,
		watcherStopChannel: make(chan struct{}, 1),
		startBlock:         startBlock,
		pageSize:           pageSize,
		checkTimeout:       checkTimeout,
		candidates:         make(map[string]blockchain.PegOutRequested),
		stateMutex:         sync.RWMutex{},
	}
}

func (watcher *PegoutEscrowWatcher) Prepare(ctx context.Context) error {
	if watcher.contracts.PegOutEscrow == nil {
		return nil
	}
	if err := watcher.loadCheckpoint(ctx); err != nil {
		return err
	}
	if err := watcher.loadCandidates(ctx); err != nil {
		return err
	}
	if err := watcher.reconcileCandidates(ctx); err != nil {
		return err
	}
	log.Info(LogPegoutEscrowStart(watcher.lastScannedBlock + 1))
	watcher.tryClaimCandidates(ctx)
	return nil
}

func (watcher *PegoutEscrowWatcher) Start() {
watcherLoop:
	for {
		select {
		case <-watcher.ticker.C():
			watcher.onTick()
		case <-watcher.watcherStopChannel:
			watcher.ticker.Stop()
			close(watcher.watcherStopChannel)
			break watcherLoop
		}
	}
}

func (watcher *PegoutEscrowWatcher) Shutdown(closeChannel chan<- bool) {
	watcher.watcherStopChannel <- struct{}{}
	closeChannel <- true
	log.Debug(LogPegoutEscrowShutdown)
}

func (watcher *PegoutEscrowWatcher) GetCandidates() []blockchain.PegOutRequested {
	watcher.stateMutex.RLock()
	defer watcher.stateMutex.RUnlock()
	candidates := make([]blockchain.PegOutRequested, 0, len(watcher.candidates))
	for _, candidate := range watcher.candidates {
		candidates = append(candidates, candidate)
	}
	return candidates
}

func (watcher *PegoutEscrowWatcher) LastScannedBlock() uint64 {
	watcher.stateMutex.RLock()
	defer watcher.stateMutex.RUnlock()
	return watcher.lastScannedBlock
}

func (watcher *PegoutEscrowWatcher) onTick() {
	if watcher.contracts.PegOutEscrow == nil {
		return
	}
	if err := watcher.checkRequests(); err != nil {
		log.Errorf(LogPegoutEscrowError, err)
		return
	}
	checkContext, checkCancel := context.WithTimeout(context.Background(), watcher.checkTimeout)
	defer checkCancel()
	watcher.tryClaimCandidates(checkContext)
}

func (watcher *PegoutEscrowWatcher) tryClaimCandidates(ctx context.Context) {
	if watcher.claimPegOutUseCase == nil {
		return
	}
	for _, candidate := range watcher.GetCandidates() {
		claimed, err := watcher.claimPegOutUseCase.Run(ctx, candidate)
		if err != nil {
			log.Error(LogPegoutEscrowClaimError(candidate.RequestHash, err))
			continue
		}
		if claimed {
			if dropErr := watcher.dropCandidate(ctx, candidate.RequestHash); dropErr != nil {
				log.Error(LogPegoutEscrowClaimError(candidate.RequestHash, dropErr))
			}
		}
	}
}

func (watcher *PegoutEscrowWatcher) loadCheckpoint(ctx context.Context) error {
	lastScanned, found, err := watcher.repository.GetCheckpoint(ctx)
	if err != nil {
		return err
	}
	if found {
		watcher.lastScannedBlock = lastScanned
		return nil
	}
	return watcher.initializeCheckpoint(ctx)
}

func (watcher *PegoutEscrowWatcher) initializeCheckpoint(ctx context.Context) error {
	if watcher.startBlock != 0 {
		watcher.lastScannedBlock = watcher.startBlock - 1
		return watcher.repository.SetCheckpoint(ctx, watcher.lastScannedBlock)
	}
	height, err := watcher.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return err
	}
	watcher.lastScannedBlock = height
	return watcher.repository.SetCheckpoint(ctx, watcher.lastScannedBlock)
}

func (watcher *PegoutEscrowWatcher) loadCandidates(ctx context.Context) error {
	stored, err := watcher.repository.ListCandidates(ctx)
	if err != nil {
		return err
	}
	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	for _, candidate := range stored {
		watcher.candidates[candidate.RequestHash] = candidate
	}
	return nil
}

func (watcher *PegoutEscrowWatcher) reconcileCandidates(ctx context.Context) error {
	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	for _, requestHash := range watcher.candidateHashes() {
		if err := watcher.reconcileCandidate(ctx, requestHash); err != nil {
			return err
		}
	}
	return nil
}

func (watcher *PegoutEscrowWatcher) candidateHashes() []string {
	hashes := make([]string, 0, len(watcher.candidates))
	for hash := range watcher.candidates {
		hashes = append(hashes, hash)
	}
	return hashes
}

func (watcher *PegoutEscrowWatcher) reconcileCandidate(ctx context.Context, requestHash string) error {
	state, err := watcher.contracts.PegOutEscrow.GetPegOutState(requestHash)
	if err != nil {
		log.Error(LogPegoutEscrowStateError(requestHash, err))
		return nil
	}
	if state == blockchain.EscrowedPegOutStateRequested {
		return nil
	}
	return watcher.dropCandidateLocked(ctx, requestHash)
}

func (watcher *PegoutEscrowWatcher) checkRequests() error {
	checkContext, checkCancel := context.WithTimeout(context.Background(), watcher.checkTimeout)
	defer checkCancel()

	scan, err := watcher.nextRange(checkContext)
	if err != nil {
		return err
	}
	if scan.skip {
		return nil
	}
	if err = watcher.applyEvents(checkContext, scan.from, scan.to); err != nil {
		return err
	}
	return watcher.commitCheckpoint(checkContext, scan.to)
}

func (watcher *PegoutEscrowWatcher) nextRange(ctx context.Context) (pegOutEscrowScanRange, error) {
	height, err := watcher.rpc.Rsk.GetHeight(ctx)
	if err != nil {
		return pegOutEscrowScanRange{}, fmt.Errorf("error getting RSK height: %w", err)
	}
	watcher.stateMutex.RLock()
	from := watcher.lastScannedBlock + 1
	pageSize := watcher.pageSize
	watcher.stateMutex.RUnlock()
	if from > height {
		return pegOutEscrowScanRange{skip: true}, nil
	}
	to := from + pageSize - 1
	if to > height {
		to = height
	}
	return pegOutEscrowScanRange{from: from, to: to}, nil
}

func (watcher *PegoutEscrowWatcher) applyEvents(ctx context.Context, from, to uint64) error {
	events, err := watcher.fetchEvents(ctx, from, to)
	if err != nil {
		return err
	}
	log.Info(LogPegoutEscrowChecking(from, to, len(events.requested)))
	return watcher.updateCandidates(ctx, events)
}

func (watcher *PegoutEscrowWatcher) fetchEvents(ctx context.Context, from, to uint64) (pegOutEscrowEvents, error) {
	requested, err := watcher.contracts.PegOutEscrow.GetPegOutRequestedEvents(ctx, from, &to)
	if err != nil {
		return pegOutEscrowEvents{}, fmt.Errorf("error fetching PegOutRequested events: %w", err)
	}
	claimed, err := watcher.contracts.PegOutEscrow.GetPegOutClaimedEvents(ctx, from, &to)
	if err != nil {
		return pegOutEscrowEvents{}, fmt.Errorf("error fetching PegOutClaimed events: %w", err)
	}
	cancelled, err := watcher.contracts.PegOutEscrow.GetPegOutCancelledEvents(ctx, from, &to)
	if err != nil {
		return pegOutEscrowEvents{}, fmt.Errorf("error fetching PegOutCancelled events: %w", err)
	}
	return pegOutEscrowEvents{requested: requested, claimed: claimed, cancelled: cancelled}, nil
}

func (watcher *PegoutEscrowWatcher) updateCandidates(ctx context.Context, events pegOutEscrowEvents) error {
	if err := watcher.addRequested(ctx, events.requested); err != nil {
		return err
	}
	return watcher.dropClosed(ctx, events)
}

func (watcher *PegoutEscrowWatcher) addRequested(ctx context.Context, requested []blockchain.PegOutRequested) error {
	for _, candidate := range requested {
		if err := watcher.addCandidate(ctx, candidate); err != nil {
			return err
		}
	}
	return nil
}

func (watcher *PegoutEscrowWatcher) dropClosed(ctx context.Context, events pegOutEscrowEvents) error {
	for _, claimed := range events.claimed {
		if err := watcher.dropCandidate(ctx, claimed.RequestHash); err != nil {
			return err
		}
	}
	for _, cancelled := range events.cancelled {
		if err := watcher.dropCandidate(ctx, cancelled.RequestHash); err != nil {
			return err
		}
	}
	return nil
}

func (watcher *PegoutEscrowWatcher) addCandidate(ctx context.Context, candidate blockchain.PegOutRequested) error {
	if err := watcher.repository.UpsertCandidate(ctx, candidate); err != nil {
		return err
	}
	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	watcher.candidates[candidate.RequestHash] = candidate
	return nil
}

func (watcher *PegoutEscrowWatcher) dropCandidate(ctx context.Context, requestHash string) error {
	watcher.stateMutex.Lock()
	defer watcher.stateMutex.Unlock()
	return watcher.dropCandidateLocked(ctx, requestHash)
}

func (watcher *PegoutEscrowWatcher) dropCandidateLocked(ctx context.Context, requestHash string) error {
	if err := watcher.repository.DeleteCandidate(ctx, requestHash); err != nil {
		return err
	}
	delete(watcher.candidates, requestHash)
	return nil
}

func (watcher *PegoutEscrowWatcher) commitCheckpoint(ctx context.Context, to uint64) error {
	if err := watcher.repository.SetCheckpoint(ctx, to); err != nil {
		return err
	}
	watcher.stateMutex.Lock()
	watcher.lastScannedBlock = to
	watcher.stateMutex.Unlock()
	return nil
}
