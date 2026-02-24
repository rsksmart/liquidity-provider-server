package liquidity_provider

import (
	"context"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type LowLiquidityAlertUseCase struct {
	peginProvider     liquidity_provider.PeginLiquidityProvider
	pegoutProvider    liquidity_provider.PegoutLiquidityProvider
	alertSender       alerts.AlertSender
	recipient         string
	warningThreshold  uint64
	criticalThreshold uint64
}

func NewLowLiquidityAlertUseCase(
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	alertSender alerts.AlertSender,
	recipient string,
	warningThreshold uint64,
	criticalThreshold uint64,
) *LowLiquidityAlertUseCase {
	return &LowLiquidityAlertUseCase{
		peginProvider:     peginProvider,
		pegoutProvider:    pegoutProvider,
		alertSender:       alertSender,
		recipient:         recipient,
		warningThreshold:  warningThreshold,
		criticalThreshold: criticalThreshold,
	}
}

func (useCase *LowLiquidityAlertUseCase) Run(ctx context.Context) error {
	btcLiquidity, err := useCase.pegoutProvider.AvailablePegoutLiquidity(ctx)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.LowLiquidityAlertId, err)
	}
	rbtcLiquidity, err := useCase.peginProvider.AvailablePeginLiquidity(ctx)
	if err != nil {
		return usecases.WrapUseCaseError(usecases.LowLiquidityAlertId, err)
	}

	warningWei := entities.CoinToWei(useCase.warningThreshold)
	criticalWei := entities.CoinToWei(useCase.criticalThreshold)

	useCase.checkAndAlert(ctx, "BTC", btcLiquidity, warningWei, criticalWei)
	useCase.checkAndAlert(ctx, "RBTC", rbtcLiquidity, warningWei, criticalWei)

	return nil
}

func (useCase *LowLiquidityAlertUseCase) checkAndAlert(
	ctx context.Context,
	network string,
	current *entities.Wei,
	warningThreshold, criticalThreshold *entities.Wei,
) {
	if current.Cmp(criticalThreshold) < 0 {
		useCase.sendAlert(ctx, network, current, criticalThreshold, alerts.AlertSubjectHotWalletLowLiquidityCritical)
		return
	}
	if current.Cmp(warningThreshold) < 0 {
		useCase.sendAlert(ctx, network, current, warningThreshold, alerts.AlertSubjectHotWalletLowLiquidityWarning)
	}
}

func (useCase *LowLiquidityAlertUseCase) sendAlert(
	ctx context.Context,
	network string,
	current, threshold *entities.Wei,
	subject string,
) {
	body := fmt.Sprintf("Network: %s | Current: %s | Threshold: %s",
		network,
		current.ToRbtc().Text('f', 18),
		threshold.ToRbtc().Text('f', 18),
	)
	if err := useCase.alertSender.SendAlert(ctx, subject, body, []string{useCase.recipient}); err != nil {
		log.Error("Error sending low liquidity alert: ", err)
	}
}
