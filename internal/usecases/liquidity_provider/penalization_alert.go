package liquidity_provider

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	log "github.com/sirupsen/logrus"
)

type PenalizationAlertUseCase struct {
	lbc       blockchain.LiquidityBridgeContract
	sender    entities.AlertSender
	recipient string
}

func NewPenalizationAlertUseCase(lbc blockchain.LiquidityBridgeContract, sender entities.AlertSender, recipient string) *PenalizationAlertUseCase {
	return &PenalizationAlertUseCase{lbc: lbc, sender: sender, recipient: recipient}
}

func (useCase *PenalizationAlertUseCase) Run(ctx context.Context, fromBlock, toBlock uint64) error {
	var body string
	events, err := useCase.lbc.GetPeginPunishmentEvents(ctx, fromBlock, &toBlock)
	if err != nil {
		return err
	}
	for _, event := range events {
		body = fmt.Sprintf("You were punished in %v rBTC for the quoteHash %s", event.Penalty.ToRbtc(), event.QuoteHash)
		if err = useCase.sender.SendAlert(ctx, "Pegin Punishment", body, useCase.recipient); err != nil {
			log.Error("Error sending punishment alert: ", err)
		}
	}
	return nil
}
