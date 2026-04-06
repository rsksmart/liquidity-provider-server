package watcher

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

const (
	reorgWindowSize            = 10
	minReorgWindow             = 12
	nodeReorgAlertBodyTemplate = "Your %s node has a blockchain reorganization depth of %d blocks, exceeding the configured maximum of %d."
)

var errHistoryWindowExceeded = errors.New("history window exceeded")

type reorgChainState struct {
	lastTipHash       string
	lastTipHeight     uint64
	knownHashByHeight map[uint64]string
}

type NodeReorgCheckUseCase struct {
	rpc            blockchain.Rpc
	alertSender    alerts.AlertSender
	alertRecipient string
	eventBus       entities.EventBus
	maxDepth       uint64
	alertCooldown  time.Duration
	btcMu          sync.Mutex
	btcLastAlert   time.Time
	btcState       *reorgChainState
	rskMu          sync.Mutex
	rskLastAlert   time.Time
	rskState       *reorgChainState
}

func NewNodeReorgCheckUseCase(
	rpc blockchain.Rpc,
	alertSender alerts.AlertSender,
	alertRecipient string,
	eventBus entities.EventBus,
	maxDepth uint64,
	alertCooldown time.Duration,
) *NodeReorgCheckUseCase {
	return &NodeReorgCheckUseCase{
		rpc:            rpc,
		alertSender:    alertSender,
		alertRecipient: alertRecipient,
		eventBus:       eventBus,
		maxDepth:       maxDepth,
		alertCooldown:  alertCooldown,
	}
}

func reorgHistoryWindow(maxDepth uint64) uint64 {
	window := maxDepth + reorgWindowSize
	if window < minReorgWindow {
		window = minReorgWindow
	}
	return window
}

func (useCase *NodeReorgCheckUseCase) Run(ctx context.Context, nodeType entities.NodeType) error {
	switch nodeType {
	case entities.NodeTypeBitcoin:
		return useCase.runBitcoin(ctx)
	case entities.NodeTypeRootstock:
		return useCase.runRootstock(ctx)
	default:
		return fmt.Errorf("unsupported node type for reorg check: %s", nodeType)
	}
}

func (useCase *NodeReorgCheckUseCase) reorgError(nodeType entities.NodeType, err error) error {
	useCase.publishReorgCheckError(nodeType)
	return usecases.WrapUseCaseError(usecases.NodeReorgAlertId, err)
}

func (useCase *NodeReorgCheckUseCase) handleReorgAlert(
	ctx context.Context,
	nodeType entities.NodeType,
	reorgDepth uint64,
	maxDepth uint64,
) error {
	var lastAlert *time.Time
	switch nodeType {
	case entities.NodeTypeBitcoin:
		lastAlert = &useCase.btcLastAlert
	case entities.NodeTypeRootstock:
		lastAlert = &useCase.rskLastAlert
	}
	if time.Since(*lastAlert) < useCase.alertCooldown {
		return nil
	}
	body := fmt.Sprintf(nodeReorgAlertBodyTemplate, nodeType, reorgDepth, maxDepth)
	if alertErr := useCase.alertSender.SendAlert(ctx, alerts.AlertSubjectNodeReorg, body, []string{useCase.alertRecipient}); alertErr != nil {
		log.Errorf("NodeReorgCheckUseCase[%s]: error sending reorg alert: %v", nodeType, alertErr)
		return usecases.WrapUseCaseError(usecases.NodeReorgAlertId, alertErr)
	}
	*lastAlert = time.Now()
	useCase.eventBus.Publish(blockchain.NodeReorgAlertSentEvent{
		BaseEvent:     entities.NewBaseEvent(blockchain.NodeReorgAlertSentEventId),
		NodeType:      nodeType,
		DetectedDepth: reorgDepth,
	})
	return nil
}

func (useCase *NodeReorgCheckUseCase) runBitcoin(ctx context.Context) error {
	useCase.btcMu.Lock()
	defer useCase.btcMu.Unlock()

	btc := useCase.rpc.Btc
	maxDepth := useCase.maxDepth
	historyWindow := reorgHistoryWindow(maxDepth)
	nodeType := entities.NodeTypeBitcoin

	info, err := btc.GetBlockchainInfo()
	if err != nil {
		log.Errorf("NodeReorgCheckUseCase[bitcoin]: error getting blockchain info: %v", err)
		return useCase.reorgError(nodeType, err)
	}

	tipHash := info.BestBlockHash
	tipHeight := uint64(info.ValidatedBlocks.Int64())

	if useCase.btcState == nil {
		useCase.btcState = &reorgChainState{knownHashByHeight: make(map[uint64]string)}
	}
	state := useCase.btcState

	if state.lastTipHash == "" || tipHash == state.lastTipHash {
		return useCase.handleBitcoinInitialOrSameTip(state, btc, tipHash, tipHeight, historyWindow, maxDepth)
	}

	parentHash, err := useCase.getBitcoinParentHash(btc, tipHash)
	if err != nil {
		return err
	}

	if parentHash == state.lastTipHash {
		return useCase.handleBitcoinTipAppend(state, btc, tipHash, tipHeight, historyWindow, maxDepth)
	}

	return useCase.handleBitcoinPotentialReorg(ctx, state, btc, tipHash, tipHeight, historyWindow, maxDepth)
}

func (useCase *NodeReorgCheckUseCase) handleBitcoinInitialOrSameTip(
	state *reorgChainState,
	btc blockchain.BitcoinNetwork,
	tipHash string,
	tipHeight uint64,
	historyWindow uint64,
	maxDepth uint64,
) error {
	if state.lastTipHash == "" {
		if initErr := useCase.initBitcoinState(state, btc, tipHash, tipHeight, historyWindow); initErr != nil {
			return initErr
		}
	}
	useCase.publishReorgCheck(entities.NodeTypeBitcoin, 0, maxDepth, false)
	return nil
}

func (useCase *NodeReorgCheckUseCase) handleBitcoinTipAppend(
	state *reorgChainState,
	btc blockchain.BitcoinNetwork,
	tipHash string,
	tipHeight uint64,
	historyWindow uint64,
	maxDepth uint64,
) error {
	useCase.advanceBitcoinTip(state, tipHash, tipHeight)
	if err := refreshBitcoinWindow(state, btc, tipHeight, historyWindow); err != nil {
		return useCase.reorgError(entities.NodeTypeBitcoin, err)
	}
	useCase.publishReorgCheck(entities.NodeTypeBitcoin, 0, maxDepth, false)
	return nil
}

func (useCase *NodeReorgCheckUseCase) handleBitcoinPotentialReorg(
	ctx context.Context,
	state *reorgChainState,
	btc blockchain.BitcoinNetwork,
	tipHash string,
	tipHeight uint64,
	historyWindow uint64,
	maxDepth uint64,
) error {
	reorgDepth, err := bitcoinChainDivergenceDepth(btc, state, tipHeight, historyWindow)
	if errors.Is(err, errHistoryWindowExceeded) {
		state.lastTipHash = tipHash
		state.lastTipHeight = tipHeight
		if refreshErr := refreshBitcoinWindow(state, btc, tipHeight, historyWindow); refreshErr != nil {
			return useCase.reorgError(entities.NodeTypeBitcoin, refreshErr)
		}
		return nil
	}
	if err != nil {
		return useCase.reorgError(entities.NodeTypeBitcoin, err)
	}

	above := reorgDepth > maxDepth
	useCase.publishReorgCheck(entities.NodeTypeBitcoin, reorgDepth, maxDepth, above)
	if above {
		if alertErr := useCase.handleReorgAlert(ctx, entities.NodeTypeBitcoin, reorgDepth, maxDepth); alertErr != nil {
			return alertErr
		}
	}

	state.lastTipHash = tipHash
	state.lastTipHeight = tipHeight
	if refreshErr := refreshBitcoinWindow(state, btc, tipHeight, historyWindow); refreshErr != nil {
		return useCase.reorgError(entities.NodeTypeBitcoin, refreshErr)
	}
	return nil
}

func (useCase *NodeReorgCheckUseCase) initBitcoinState(
	state *reorgChainState,
	btc blockchain.BitcoinNetwork,
	tipHash string,
	tipHeight uint64,
	historyWindow uint64,
) error {
	if err := refreshBitcoinWindow(state, btc, tipHeight, historyWindow); err != nil {
		return useCase.reorgError(entities.NodeTypeBitcoin, err)
	}
	state.lastTipHash = tipHash
	state.lastTipHeight = tipHeight
	return nil
}

func (useCase *NodeReorgCheckUseCase) getBitcoinParentHash(btc blockchain.BitcoinNetwork, tipHash string) (string, error) {
	header, err := btc.GetBlockHeaderVerbose(tipHash)
	if err != nil {
		log.Errorf("NodeReorgCheckUseCase[bitcoin]: error getting block header: %v", err)
		return "", useCase.reorgError(entities.NodeTypeBitcoin, err)
	}
	return header.PreviousHash, nil
}

func (useCase *NodeReorgCheckUseCase) advanceBitcoinTip(
	state *reorgChainState,
	tipHash string,
	tipHeight uint64,
) {
	state.lastTipHash = tipHash
	state.lastTipHeight = tipHeight
}

func refreshBitcoinWindow(state *reorgChainState, btc blockchain.BitcoinNetwork, tipHeight uint64, historyWindow uint64) error {
	state.knownHashByHeight = make(map[uint64]string)
	for i := uint64(0); i < historyWindow; i++ {
		if tipHeight < i {
			break
		}
		blockHeight := tipHeight - i
		hash, err := btc.GetBlockHashAtHeight(int64(blockHeight))
		if err != nil {
			return err
		}
		state.knownHashByHeight[blockHeight] = hash
	}
	return nil
}

// bitcoinChainDivergenceDepth walks backwards from the lower of the current and
// previously seen tip heights until it finds the last block height whose hash
// still matches the locally cached history window.
//
// The returned value is the distance from the previously seen tip to that last
// common block. This is the amount of chain history that diverged from the last
// observed state, which can correspond to a reorg or to another discontinuity
// between executions.
//
// If no common block is found inside the retained history window,
// errHistoryWindowExceeded is returned. Callers must treat that as "depth
// unknown, resync required" rather than as a confirmed reorg depth.
func bitcoinChainDivergenceDepth(
	btc blockchain.BitcoinNetwork,
	state *reorgChainState,
	currentTipHeight uint64,
	historyWindow uint64,
) (uint64, error) {
	minTip := state.lastTipHeight
	if currentTipHeight < minTip {
		minTip = currentTipHeight
	}
	maxTip := state.lastTipHeight
	if currentTipHeight > maxTip {
		maxTip = currentTipHeight
	}

	var lowerBound uint64
	if maxTip >= historyWindow {
		lowerBound = maxTip - historyWindow + 1
	}

	for blockHeight := minTip; ; {
		if blockHeight < lowerBound {
			return 0, errHistoryWindowExceeded
		}
		curHash, err := btc.GetBlockHashAtHeight(int64(blockHeight))
		if err != nil {
			return 0, err
		}
		if oldHash, ok := state.knownHashByHeight[blockHeight]; ok && curHash == oldHash {
			return state.lastTipHeight - blockHeight, nil
		}
		if blockHeight == 0 {
			return 0, errHistoryWindowExceeded
		}
		blockHeight--
	}
}

func (useCase *NodeReorgCheckUseCase) runRootstock(ctx context.Context) error {
	useCase.rskMu.Lock()
	defer useCase.rskMu.Unlock()

	rsk := useCase.rpc.Rsk
	maxDepth := useCase.maxDepth
	historyWindow := reorgHistoryWindow(maxDepth)
	nodeType := entities.NodeTypeRootstock

	tipHeight, err := rsk.GetHeight(ctx)
	if err != nil {
		log.Errorf("NodeReorgCheckUseCase[rootstock]: error getting chain height: %v", err)
		return useCase.reorgError(nodeType, err)
	}

	tipBlock, err := rsk.GetBlockByNumber(ctx, new(big.Int).SetUint64(tipHeight))
	if err != nil {
		log.Errorf("NodeReorgCheckUseCase[rootstock]: error getting tip block: %v", err)
		return useCase.reorgError(nodeType, err)
	}
	tipHash := tipBlock.Hash

	if useCase.rskState == nil {
		useCase.rskState = &reorgChainState{knownHashByHeight: make(map[uint64]string)}
	}
	state := useCase.rskState

	if state.lastTipHash == "" || tipHash == state.lastTipHash {
		return useCase.handleRootstockInitialOrSameTip(ctx, state, rsk, tipHash, tipHeight, historyWindow, maxDepth)
	}

	if tipBlock.ParentHash == state.lastTipHash {
		return useCase.handleRootstockTipAppend(ctx, state, rsk, tipHash, tipHeight, historyWindow, maxDepth)
	}

	return useCase.handleRootstockPotentialReorg(ctx, state, rsk, tipHash, tipHeight, historyWindow, maxDepth)
}

func (useCase *NodeReorgCheckUseCase) handleRootstockInitialOrSameTip(
	ctx context.Context,
	state *reorgChainState,
	rsk blockchain.RootstockRpcServer,
	tipHash string,
	tipHeight uint64,
	historyWindow uint64,
	maxDepth uint64,
) error {
	if state.lastTipHash == "" {
		if initErr := useCase.initRootstockState(ctx, state, rsk, tipHash, tipHeight, historyWindow); initErr != nil {
			return initErr
		}
	}
	useCase.publishReorgCheck(entities.NodeTypeRootstock, 0, maxDepth, false)
	return nil
}

func (useCase *NodeReorgCheckUseCase) handleRootstockTipAppend(
	ctx context.Context,
	state *reorgChainState,
	rsk blockchain.RootstockRpcServer,
	tipHash string,
	tipHeight uint64,
	historyWindow uint64,
	maxDepth uint64,
) error {
	useCase.advanceRootstockTip(state, tipHash, tipHeight)
	if err := refreshRootstockWindow(ctx, state, rsk, tipHeight, historyWindow); err != nil {
		return useCase.reorgError(entities.NodeTypeRootstock, err)
	}
	useCase.publishReorgCheck(entities.NodeTypeRootstock, 0, maxDepth, false)
	return nil
}

func (useCase *NodeReorgCheckUseCase) handleRootstockPotentialReorg(
	ctx context.Context,
	state *reorgChainState,
	rsk blockchain.RootstockRpcServer,
	tipHash string,
	tipHeight uint64,
	historyWindow uint64,
	maxDepth uint64,
) error {
	reorgDepth, err := rootstockChainDivergenceDepth(ctx, rsk, state, tipHeight, historyWindow)
	if errors.Is(err, errHistoryWindowExceeded) {
		state.lastTipHash = tipHash
		state.lastTipHeight = tipHeight
		if refreshErr := refreshRootstockWindow(ctx, state, rsk, tipHeight, historyWindow); refreshErr != nil {
			return useCase.reorgError(entities.NodeTypeRootstock, refreshErr)
		}
		return nil
	}
	if err != nil {
		return useCase.reorgError(entities.NodeTypeRootstock, err)
	}

	above := reorgDepth > maxDepth
	useCase.publishReorgCheck(entities.NodeTypeRootstock, reorgDepth, maxDepth, above)
	if above {
		if alertErr := useCase.handleReorgAlert(ctx, entities.NodeTypeRootstock, reorgDepth, maxDepth); alertErr != nil {
			return alertErr
		}
	}

	state.lastTipHash = tipHash
	state.lastTipHeight = tipHeight
	if refreshErr := refreshRootstockWindow(ctx, state, rsk, tipHeight, historyWindow); refreshErr != nil {
		return useCase.reorgError(entities.NodeTypeRootstock, refreshErr)
	}
	return nil
}

func (useCase *NodeReorgCheckUseCase) initRootstockState(
	ctx context.Context,
	state *reorgChainState,
	rsk blockchain.RootstockRpcServer,
	tipHash string,
	tipHeight uint64,
	historyWindow uint64,
) error {
	if err := refreshRootstockWindow(ctx, state, rsk, tipHeight, historyWindow); err != nil {
		return useCase.reorgError(entities.NodeTypeRootstock, err)
	}
	state.lastTipHash = tipHash
	state.lastTipHeight = tipHeight
	return nil
}

func (useCase *NodeReorgCheckUseCase) advanceRootstockTip(
	state *reorgChainState,
	tipHash string,
	tipHeight uint64,
) {
	state.lastTipHash = tipHash
	state.lastTipHeight = tipHeight
}

func refreshRootstockWindow(
	ctx context.Context,
	state *reorgChainState,
	rsk blockchain.RootstockRpcServer,
	tipHeight uint64,
	historyWindow uint64,
) error {
	state.knownHashByHeight = make(map[uint64]string)
	for i := uint64(0); i < historyWindow; i++ {
		if tipHeight < i {
			break
		}
		blockHeight := tipHeight - i
		bi, err := rsk.GetBlockByNumber(ctx, new(big.Int).SetUint64(blockHeight))
		if err != nil {
			return err
		}
		state.knownHashByHeight[blockHeight] = bi.Hash
	}
	return nil
}

// rootstockChainDivergenceDepth walks backwards from the lower of the current
// and previously seen tip heights until it finds the last block height whose
// hash still matches the locally cached history window.
//
// The returned value is the distance from the previously seen tip to that last
// common block. This is the amount of chain history that diverged from the last
// observed state, which can correspond to a reorg or to another discontinuity
// between executions.
//
// If no common block is found inside the retained history window,
// errHistoryWindowExceeded is returned. Callers must treat that as "depth
// unknown, resync required" rather than as a confirmed reorg depth.
func rootstockChainDivergenceDepth(
	ctx context.Context,
	rsk blockchain.RootstockRpcServer,
	state *reorgChainState,
	currentTipHeight uint64,
	historyWindow uint64,
) (uint64, error) {
	minTip := state.lastTipHeight
	if currentTipHeight < minTip {
		minTip = currentTipHeight
	}
	maxTip := state.lastTipHeight
	if currentTipHeight > maxTip {
		maxTip = currentTipHeight
	}

	var lowerBound uint64
	if maxTip >= historyWindow {
		lowerBound = maxTip - historyWindow + 1
	}

	for blockHeight := minTip; ; {
		if blockHeight < lowerBound {
			return 0, errHistoryWindowExceeded
		}
		bi, err := rsk.GetBlockByNumber(ctx, new(big.Int).SetUint64(blockHeight))
		if err != nil {
			return 0, err
		}
		curHash := bi.Hash
		if oldHash, ok := state.knownHashByHeight[blockHeight]; ok && curHash == oldHash {
			return state.lastTipHeight - blockHeight, nil
		}
		if blockHeight == 0 {
			return 0, errHistoryWindowExceeded
		}
		blockHeight--
	}
}

func (useCase *NodeReorgCheckUseCase) publishReorgCheck(
	nodeType entities.NodeType,
	currentDepth uint64,
	maxAllowed uint64,
	aboveThreshold bool,
) {
	useCase.eventBus.Publish(blockchain.NodeReorgCheckEvent{
		BaseEvent:       entities.NewBaseEvent(blockchain.NodeReorgCheckEventId),
		NodeType:        nodeType,
		CurrentDepth:    currentDepth,
		MaxAllowedDepth: maxAllowed,
		AboveThreshold:  aboveThreshold,
	})
}

func (useCase *NodeReorgCheckUseCase) publishReorgCheckError(nodeType entities.NodeType) {
	useCase.eventBus.Publish(blockchain.NodeReorgCheckErrorEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.NodeReorgCheckErrorEventId),
		NodeType:  nodeType,
	})
}
