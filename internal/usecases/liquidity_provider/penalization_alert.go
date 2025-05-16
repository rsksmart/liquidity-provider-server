package liquidity_provider

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type PenalizationAlertUseCase struct {
	contracts blockchain.RskContracts
	sender    entities.AlertSender
	recipient string
	repo      liquidity_provider.LiquidityProviderRepository
}

func NewPenalizationAlertUseCase(
	contracts blockchain.RskContracts,
	sender entities.AlertSender,
	recipient string,
	repo liquidity_provider.LiquidityProviderRepository,
) *PenalizationAlertUseCase {
	return &PenalizationAlertUseCase{
		contracts: contracts,
		sender:    sender,
		recipient: recipient,
		repo:      repo,
	}
}

func (useCase *PenalizationAlertUseCase) Run(ctx context.Context, fromBlock, toBlock uint64) error {
	var body string
	events, err := useCase.contracts.Lbc.GetPunishmentEvents(ctx, fromBlock, &toBlock)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.PenalizationId, err)
	}

	for _, event := range events {
		err := useCase.repo.InsertPenalization(ctx, event)
		if err != nil {
			log.Error(usecases.PenalizationId, err)
		}
		body = fmt.Sprintf("You were punished in %v rBTC for the quoteHash %s", event.Penalty.ToRbtc(), event.QuoteHash)
		if err = useCase.sender.SendAlert(ctx, "Pegin Punishment", body, []string{useCase.recipient}); err != nil {
			log.Error("Error sending punishment alert: ", err)
		}
	}
	return nil
}
