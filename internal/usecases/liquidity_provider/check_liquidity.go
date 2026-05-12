package liquidity_provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	log "github.com/sirupsen/logrus"
)

type OperationType string

const (
	PeginOperation  OperationType = "PegIn"
	PegoutOperation OperationType = "PegOut"
	MessageBody     string        = "You are out of liquidity to perform a %s. Please, do a deposit"
)

type CheckLiquidityUseCase struct {
	peginProvider  liquidity_provider.PeginLiquidityProvider
	pegoutProvider liquidity_provider.PegoutLiquidityProvider
	contracts      blockchain.RskContracts
	alertSender    alerts.AlertSender
	recipient      string
}

func NewCheckLiquidityUseCase(
	peginProvider liquidity_provider.PeginLiquidityProvider,
	pegoutProvider liquidity_provider.PegoutLiquidityProvider,
	contracts blockchain.RskContracts,
	alertSender alerts.AlertSender,
	recipient string,
) *CheckLiquidityUseCase {
	return &CheckLiquidityUseCase{
		peginProvider:  peginProvider,
		pegoutProvider: pegoutProvider,
		contracts:      contracts,
		alertSender:    alertSender,
		recipient:      recipient,
	}
}

func (useCase *CheckLiquidityUseCase) Run(ctx context.Context) error {
	minLockTxValueInWei, err := useCase.contracts.Bridge.GetMinimumLockTxValue()
	if err != nil {
		return usecases.WrapUseCaseError(usecases.CheckLiquidityId, err)
	}

	err = useCase.peginProvider.HasPeginLiquidity(ctx, minLockTxValueInWei)
	if errors.Is(err, usecases.NoLiquidityError) {
		if err = useCase.alertSender.SendAlert(
			ctx,
			alerts.AlertSubjectPeginOutOfLiquidity,
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
			alerts.AlertSubjectPegoutOutOfLiquidity,
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
