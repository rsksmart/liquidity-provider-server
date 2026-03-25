package watcher

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

var NodeReorgAlertBodyTemplate = "Your %s node has a blockchain reorganization depth of %d blocks, exceeding the configured maximum of %d."

type reorgChainState struct {
	lastTipHash       string
	lastTipHeight     uint64
	knownHashByHeight map[uint64]string
}

type NodeReorgCheckUseCase struct {
	rpc             blockchain.Rpc
	alertSender     alerts.AlertSender
	alertRecipient  string
	eventBus        entities.EventBus
	maxDepthByNode  map[entities.NodeType]uint64
	alertCooldown   time.Duration
	lastAlertByNode map[entities.NodeType]time.Time
	bitcoinState    *reorgChainState
	rootstockState  *reorgChainState
}

func NewNodeReorgCheckUseCase(
	rpc blockchain.Rpc,
	alertSender alerts.AlertSender,
	alertRecipient string,
	eventBus entities.EventBus,
	bitcoinMaxReorgDepth uint64,
	rootstockMaxReorgDepth uint64,
	alertCooldown time.Duration,
) *NodeReorgCheckUseCase {
	return &NodeReorgCheckUseCase{
		rpc:            rpc,
		alertSender:    alertSender,
		alertRecipient: alertRecipient,
		eventBus:       eventBus,
		maxDepthByNode: map[entities.NodeType]uint64{
			entities.NodeTypeBitcoin:   bitcoinMaxReorgDepth,
			entities.NodeTypeRootstock: rootstockMaxReorgDepth,
		},
		alertCooldown:   alertCooldown,
		lastAlertByNode: make(map[entities.NodeType]time.Time),
	}
}

func reorgHistoryWindow(maxDepth uint64) uint64 {
	window := maxDepth + 10
	if window < 12 {
		window = 12
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
	if time.Since(useCase.lastAlertByNode[nodeType]) < useCase.alertCooldown {
		return nil
	}
	body := fmt.Sprintf(NodeReorgAlertBodyTemplate, nodeType, reorgDepth, maxDepth)
	if alertErr := useCase.alertSender.SendAlert(ctx, alerts.AlertSubjectNodeReorg, body, []string{useCase.alertRecipient}); alertErr != nil {
		log.Errorf("NodeReorgCheckUseCase[%s]: error sending reorg alert: %v", nodeType, alertErr)
		return useCase.reorgError(nodeType, alertErr)
	}
	useCase.lastAlertByNode[nodeType] = time.Now()
	useCase.eventBus.Publish(blockchain.NodeReorgAlertSentEvent{
		BaseEvent:     entities.NewBaseEvent(blockchain.NodeReorgAlertSentEventId),
		NodeType:      nodeType,
		DetectedDepth: reorgDepth,
	})
	return nil
}

func (useCase *NodeReorgCheckUseCase) runBitcoin(ctx context.Context) error {
	btc := useCase.rpc.Btc
	maxDepth := useCase.maxDepthByNode[entities.NodeTypeBitcoin]
	hw := reorgHistoryWindow(maxDepth)
	nodeType := entities.NodeTypeBitcoin

	info, err := btc.GetBlockchainInfo()
	if err != nil {
		log.Errorf("NodeReorgCheckUseCase[bitcoin]: error getting blockchain info: %v", err)
		return useCase.reorgError(nodeType, err)
	}

	tipHash := info.BestBlockHash
	tipHeight := uint64(info.ValidatedBlocks.Int64())

	if useCase.bitcoinState == nil {
		useCase.bitcoinState = &reorgChainState{knownHashByHeight: make(map[uint64]string)}
	}
	st := useCase.bitcoinState

	if st.lastTipHash == "" {
		return useCase.initBitcoinState(st, btc, tipHash, tipHeight, maxDepth, hw)
	}

	if tipHash == st.lastTipHash {
		useCase.publishReorgCheck(nodeType, 0, maxDepth, false)
		return nil
	}

	parentHash, err := useCase.getBitcoinParentHash(btc, tipHash)
	if err != nil {
		return err
	}

	if parentHash == st.lastTipHash {
		return useCase.advanceBitcoinTip(st, btc, tipHash, tipHeight, maxDepth, hw)
	}

	reorgDepth := bitcoinReorgDepth(btc, st, tipHeight, hw, maxDepth)
	above := reorgDepth > maxDepth
	useCase.publishReorgCheck(nodeType, reorgDepth, maxDepth, above)

	if above {
		if alertErr := useCase.handleReorgAlert(ctx, nodeType, reorgDepth, maxDepth); alertErr != nil {
			return alertErr
		}
	}

	st.lastTipHash = tipHash
	st.lastTipHeight = tipHeight
	if err = refreshBitcoinWindow(st, btc, tipHeight, hw); err != nil {
		return useCase.reorgError(nodeType, err)
	}
	return nil
}

func (useCase *NodeReorgCheckUseCase) initBitcoinState(
	st *reorgChainState,
	btc blockchain.BitcoinNetwork,
	tipHash string,
	tipHeight uint64,
	maxDepth uint64,
	hw uint64,
) error {
	if err := refreshBitcoinWindow(st, btc, tipHeight, hw); err != nil {
		return useCase.reorgError(entities.NodeTypeBitcoin, err)
	}
	st.lastTipHash = tipHash
	st.lastTipHeight = tipHeight
	useCase.publishReorgCheck(entities.NodeTypeBitcoin, 0, maxDepth, false)
	return nil
}

func (useCase *NodeReorgCheckUseCase) getBitcoinParentHash(btc blockchain.BitcoinNetwork, tipHash string) (string, error) {
	hdr, err := btc.GetBlockHeaderVerbose(tipHash)
	if err != nil {
		log.Errorf("NodeReorgCheckUseCase[bitcoin]: error getting block header: %v", err)
		return "", useCase.reorgError(entities.NodeTypeBitcoin, err)
	}
	return hdr.PreviousHash, nil
}

func (useCase *NodeReorgCheckUseCase) advanceBitcoinTip(
	st *reorgChainState,
	btc blockchain.BitcoinNetwork,
	tipHash string,
	tipHeight uint64,
	maxDepth uint64,
	hw uint64,
) error {
	st.lastTipHash = tipHash
	st.lastTipHeight = tipHeight
	if err := refreshBitcoinWindow(st, btc, tipHeight, hw); err != nil {
		return useCase.reorgError(entities.NodeTypeBitcoin, err)
	}
	useCase.publishReorgCheck(entities.NodeTypeBitcoin, 0, maxDepth, false)
	return nil
}

func refreshBitcoinWindow(st *reorgChainState, btc blockchain.BitcoinNetwork, tipHeight uint64, hw uint64) error {
	st.knownHashByHeight = make(map[uint64]string)
	for i := uint64(0); i < hw; i++ {
		if tipHeight < i {
			break
		}
		h := tipHeight - i
		hash, err := btc.GetBlockHashAtHeight(int64(h))
		if err != nil {
			return err
		}
		st.knownHashByHeight[h] = hash
	}
	return nil
}

func bitcoinReorgDepth(
	btc blockchain.BitcoinNetwork,
	st *reorgChainState,
	currentTipHeight uint64,
	hw uint64,
	maxDepth uint64,
) uint64 {
	minTip := st.lastTipHeight
	if currentTipHeight < minTip {
		minTip = currentTipHeight
	}
	maxTip := st.lastTipHeight
	if currentTipHeight > maxTip {
		maxTip = currentTipHeight
	}

	var lowerBound uint64
	if maxTip >= hw {
		lowerBound = maxTip - hw + 1
	}

	for h := minTip; ; {
		if h < lowerBound {
			return maxDepth + 1
		}
		curHash, err := btc.GetBlockHashAtHeight(int64(h))
		if err != nil {
			return maxDepth + 1
		}
		if oldHash, ok := st.knownHashByHeight[h]; ok && curHash == oldHash {
			return st.lastTipHeight - h
		}
		if h == 0 {
			return maxDepth + 1
		}
		h--
	}
}

func (useCase *NodeReorgCheckUseCase) runRootstock(ctx context.Context) error {
	rsk := useCase.rpc.Rsk
	maxDepth := useCase.maxDepthByNode[entities.NodeTypeRootstock]
	hw := reorgHistoryWindow(maxDepth)
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

	if useCase.rootstockState == nil {
		useCase.rootstockState = &reorgChainState{knownHashByHeight: make(map[uint64]string)}
	}
	st := useCase.rootstockState

	if st.lastTipHash == "" {
		return useCase.initRootstockState(ctx, st, rsk, tipHash, tipHeight, maxDepth, hw)
	}

	if tipHash == st.lastTipHash {
		useCase.publishReorgCheck(nodeType, 0, maxDepth, false)
		return nil
	}

	if tipBlock.ParentHash == st.lastTipHash {
		return useCase.advanceRootstockTip(ctx, st, rsk, tipHash, tipHeight, maxDepth, hw)
	}

	reorgDepth := rootstockReorgDepth(ctx, rsk, st, tipHeight, hw, maxDepth)
	above := reorgDepth > maxDepth
	useCase.publishReorgCheck(nodeType, reorgDepth, maxDepth, above)

	if above {
		if alertErr := useCase.handleReorgAlert(ctx, nodeType, reorgDepth, maxDepth); alertErr != nil {
			return alertErr
		}
	}

	st.lastTipHash = tipHash
	st.lastTipHeight = tipHeight
	if err = refreshRootstockWindow(ctx, st, rsk, tipHeight, hw); err != nil {
		return useCase.reorgError(nodeType, err)
	}
	return nil
}

func (useCase *NodeReorgCheckUseCase) initRootstockState(
	ctx context.Context,
	st *reorgChainState,
	rsk blockchain.RootstockRpcServer,
	tipHash string,
	tipHeight uint64,
	maxDepth uint64,
	hw uint64,
) error {
	if err := refreshRootstockWindow(ctx, st, rsk, tipHeight, hw); err != nil {
		return useCase.reorgError(entities.NodeTypeRootstock, err)
	}
	st.lastTipHash = tipHash
	st.lastTipHeight = tipHeight
	useCase.publishReorgCheck(entities.NodeTypeRootstock, 0, maxDepth, false)
	return nil
}

func (useCase *NodeReorgCheckUseCase) advanceRootstockTip(
	ctx context.Context,
	st *reorgChainState,
	rsk blockchain.RootstockRpcServer,
	tipHash string,
	tipHeight uint64,
	maxDepth uint64,
	hw uint64,
) error {
	st.lastTipHash = tipHash
	st.lastTipHeight = tipHeight
	if err := refreshRootstockWindow(ctx, st, rsk, tipHeight, hw); err != nil {
		return useCase.reorgError(entities.NodeTypeRootstock, err)
	}
	useCase.publishReorgCheck(entities.NodeTypeRootstock, 0, maxDepth, false)
	return nil
}

func refreshRootstockWindow(
	ctx context.Context,
	st *reorgChainState,
	rsk blockchain.RootstockRpcServer,
	tipHeight uint64,
	hw uint64,
) error {
	st.knownHashByHeight = make(map[uint64]string)
	for i := uint64(0); i < hw; i++ {
		if tipHeight < i {
			break
		}
		h := tipHeight - i
		bi, err := rsk.GetBlockByNumber(ctx, new(big.Int).SetUint64(h))
		if err != nil {
			return err
		}
		st.knownHashByHeight[h] = bi.Hash
	}
	return nil
}

func rootstockReorgDepth(
	ctx context.Context,
	rsk blockchain.RootstockRpcServer,
	st *reorgChainState,
	currentTipHeight uint64,
	hw uint64,
	maxDepth uint64,
) uint64 {
	minTip := st.lastTipHeight
	if currentTipHeight < minTip {
		minTip = currentTipHeight
	}
	maxTip := st.lastTipHeight
	if currentTipHeight > maxTip {
		maxTip = currentTipHeight
	}

	var lowerBound uint64
	if maxTip >= hw {
		lowerBound = maxTip - hw + 1
	}

	for h := minTip; ; {
		if h < lowerBound {
			return maxDepth + 1
		}
		bi, err := rsk.GetBlockByNumber(ctx, new(big.Int).SetUint64(h))
		if err != nil {
			return maxDepth + 1
		}
		curHash := bi.Hash
		if oldHash, ok := st.knownHashByHeight[h]; ok && curHash == oldHash {
			return st.lastTipHeight - h
		}
		if h == 0 {
			return maxDepth + 1
		}
		h--
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
