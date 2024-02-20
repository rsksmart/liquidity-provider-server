package liquidity_provider

import (
	"context"
	"errors"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type OperationType string

const (
	PeginOperation  OperationType = "PegIn"
	PegoutOperation OperationType = "PegOut"
	MessageSubject  string        = "%s: Out of liquidity"
	MessageBody     string        = "You are out of liquidity to perform a %s. Please, do a deposit"
)

type CheckLiquidityUseCase struct {
	peginProvider  liquidity_provider.PeginLiquidityProvider
	pegoutProvider liquidity_provider.PegoutLiquidityProvider
	bridge         blockchain.RootstockBridge
	alertSender    entities.AlertSender
	recipient      string
}

func NewCheckLiquidityUseCase(
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	bridge blockchain.RootstockBridge,
	alertSender entities.AlertSender,
	recipient string,
) *CheckLiquidityUseCase {
	return &CheckLiquidityUseCase{
		peginProvider:  peginProvider,
		pegoutProvider: pegoutProvider,
		bridge:         bridge,
		alertSender:    alertSender,
		recipient:      recipient,
	}
}

func (useCase *CheckLiquidityUseCase) Run(ctx context.Context) error {
	minLockTxValueInSatoshi, err := useCase.bridge.GetMinimumLockTxValue()
	if err != nil {
		return usecases.WrapUseCaseError(usecases.CheckLiquidityId, err)
	}
	minLockTxValueInWei := entities.SatoshiToWei(minLockTxValueInSatoshi.Uint64())

	err = useCase.peginProvider.HasPeginLiquidity(ctx, minLockTxValueInWei)
	if errors.Is(err, usecases.NoLiquidityError) {
		if err = useCase.alertSender.SendAlert(
			ctx,
			fmt.Sprintf(MessageSubject, PeginOperation),
			fmt.Sprintf(MessageBody, PeginOperation),
			[]string{useCase.recipient},
		); err != nil {
			log.Error("Error sending notification to liquidity provider: ", err)
		}
	} else if err != nil {
		return usecases.WrapUseCaseError(usecases.CheckLiquidityId, err)
	}

	err = useCase.pegoutProvider.HasPegoutLiquidity(ctx, minLockTxValueInWei)
	if errors.Is(err, usecases.NoLiquidityError) {
		if err = useCase.alertSender.SendAlert(
			ctx,
			fmt.Sprintf(MessageSubject, PegoutOperation),
			fmt.Sprintf(MessageBody, PegoutOperation),
			[]string{useCase.recipient},
		); err != nil {
			log.Error("Error sending notification to liquidity provider: ", err)
		}
	} else if err != nil {
		return usecases.WrapUseCaseError(usecases.CheckLiquidityId, err)
	}

	return nil
}
