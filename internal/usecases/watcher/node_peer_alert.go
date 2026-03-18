package watcher

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
)

var LowPeersAlertBodyTemplate = "Your %s node has %d peers, which is below the configured minimum of %d. " +
	"Please check your node's network connectivity."

type NodePeerAlertUseCase struct {
	alertSender    alerts.AlertSender
	alertRecipient string
}

func NewNodePeerAlertUseCase(
	alertSender alerts.AlertSender,
	alertRecipient string,
) *NodePeerAlertUseCase {
	return &NodePeerAlertUseCase{
		alertSender:    alertSender,
		alertRecipient: alertRecipient,
	}
}

func (useCase *NodePeerAlertUseCase) Run(ctx context.Context, nodeType entities.NodeType, currentPeers int64, minPeers uint64) error {
	body := fmt.Sprintf(LowPeersAlertBodyTemplate, nodeType, currentPeers, minPeers)
	err := useCase.alertSender.SendAlert(ctx, alerts.AlertSubjectLowPeers, body, []string{useCase.alertRecipient})
	if err != nil {
		return usecases.WrapUseCaseError(usecases.NodePeerAlertId, err)
	}
	return nil
}
