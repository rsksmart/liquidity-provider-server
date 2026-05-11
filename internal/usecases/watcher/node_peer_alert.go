package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

var LowPeersAlertBodyTemplate = "Your %s node has %d peers, which is below the configured minimum of %d. " +
	"Please check your node's network connectivity."

type NodePeerCheckUseCase struct {
	rpc             blockchain.Rpc
	alertSender     alerts.AlertSender
	alertRecipient  string
	eventBus        entities.EventBus
	minPeersByNode  map[entities.NodeType]uint64
	alertCooldown   time.Duration
	lastAlertByNode map[entities.NodeType]time.Time
}

func NewNodePeerCheckUseCase(
	rpc blockchain.Rpc,
	alertSender alerts.AlertSender,
	alertRecipient string,
	eventBus entities.EventBus,
	bitcoinMinPeers uint64,
	rootstockMinPeers uint64,
	alertCooldown time.Duration,
) *NodePeerCheckUseCase {
	return &NodePeerCheckUseCase{
		rpc:            rpc,
		alertSender:    alertSender,
		alertRecipient: alertRecipient,
		eventBus:       eventBus,
		minPeersByNode: map[entities.NodeType]uint64{
			entities.NodeTypeBitcoin:   bitcoinMinPeers,
			entities.NodeTypeRootstock: rootstockMinPeers,
		},
		alertCooldown:   alertCooldown,
		lastAlertByNode: make(map[entities.NodeType]time.Time),
	}
}

func (useCase *NodePeerCheckUseCase) Run(ctx context.Context, nodeType entities.NodeType) error {
	currentPeers, err := useCase.getPeerCount(ctx, nodeType)
	if err != nil {
		log.Errorf("NodePeerCheckUseCase[%s]: error getting peer count: %v", nodeType, err)
		useCase.eventBus.Publish(blockchain.NodePeerCheckErrorEvent{
			BaseEvent: entities.NewBaseEvent(blockchain.NodePeerCheckErrorEventId),
			NodeType:  nodeType,
		})
		return usecases.WrapUseCaseError(usecases.NodePeerAlertId, err)
	}

	minPeers := useCase.minPeersByNode[nodeType]
	belowThreshold := uint64(currentPeers) < minPeers
	useCase.eventBus.Publish(blockchain.NodePeerCheckEvent{
		BaseEvent:      entities.NewBaseEvent(blockchain.NodePeerCheckEventId),
		NodeType:       nodeType,
		CurrentPeers:   currentPeers,
		MinPeers:       minPeers,
		BelowThreshold: belowThreshold,
	})

	if !belowThreshold {
		return nil
	}
	log.Warnf("NodePeerCheckUseCase[%s]: peer count %d is below minimum %d", nodeType, currentPeers, minPeers)
	if time.Since(useCase.lastAlertByNode[nodeType]) < useCase.alertCooldown {
		return nil
	}

	body := fmt.Sprintf(LowPeersAlertBodyTemplate, nodeType, currentPeers, minPeers)
	if alertErr := useCase.alertSender.SendAlert(ctx, alerts.AlertSubjectLowPeers, body, []string{useCase.alertRecipient}); alertErr != nil {
		log.Errorf("NodePeerCheckUseCase[%s]: error sending low peer alert: %v", nodeType, alertErr)
		return usecases.WrapUseCaseError(usecases.NodePeerAlertId, alertErr)
	}

	useCase.lastAlertByNode[nodeType] = time.Now()
	useCase.eventBus.Publish(blockchain.NodePeerAlertSentEvent{
		BaseEvent: entities.NewBaseEvent(blockchain.NodePeerAlertSentEventId),
		NodeType:  nodeType,
	})
	return nil
}

func (useCase *NodePeerCheckUseCase) getPeerCount(ctx context.Context, nodeType entities.NodeType) (int64, error) {
	switch nodeType {
	case entities.NodeTypeBitcoin:
		return useCase.rpc.Btc.GetConnectionCount()
	case entities.NodeTypeRootstock:
		count, err := useCase.rpc.Rsk.PeerCount(ctx)
		return int64(count), err
	default:
		return 0, fmt.Errorf("unsupported node type: %s", nodeType)
	}
}
