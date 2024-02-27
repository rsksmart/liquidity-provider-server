package liquidity_provider_test

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPenalizationAlertUseCase_Run(t *testing.T) {
	lbc := &mocks.LbcMock{}
	events := []lp.PunishmentEvent{
		{
			LiquidityProvider: "0x01",
			Penalty:           entities.NewWei(100),
			QuoteHash:         "0x02",
		},
		{
			LiquidityProvider: "0x03",
			Penalty:           entities.NewWei(200),
			QuoteHash:         "0x04",
		},
		{
			LiquidityProvider: "0x05",
			Penalty:           entities.NewWei(300),
			QuoteHash:         "0x06",
		},
	}
	toBlock := uint64(10)
	lbc.On(
		"GetPeginPunishmentEvents",
		mock.AnythingOfType("context.backgroundCtx"),
		uint64(5),
		&toBlock,
	).Return(events, nil).Once()

	sender := &mocks.AlertSenderMock{}
	recipient := "recipient@test.com"

	for i := 0; i < 3; i++ {
		sender.On(
			"SendAlert",
			mock.AnythingOfType("context.backgroundCtx"),
			"Pegin Punishment",
			fmt.Sprintf("You were punished in %v rBTC for the quoteHash %s", events[i].Penalty.ToRbtc(), events[i].QuoteHash),
			[]string{recipient},
		).Return(nil).Once()
	}

	useCase := liquidity_provider.NewPenalizationAlertUseCase(lbc, sender, recipient)
	err := useCase.Run(context.Background(), 5, 10)
	require.NoError(t, err)
	lbc.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestPenalizationAlertUseCase_Run_GetEvents(t *testing.T) {
	lbc := &mocks.LbcMock{}
	sender := &mocks.AlertSenderMock{}
	lbc.On("GetPeginPunishmentEvents", mock.AnythingOfType("context.backgroundCtx"), uint64(5), mock.Anything).
		Return([]lp.PunishmentEvent{}, assert.AnError).Once()
	useCase := liquidity_provider.NewPenalizationAlertUseCase(lbc, sender, "recipient")
	err := useCase.Run(context.Background(), 5, 10)
	lbc.AssertExpectations(t)
	require.Error(t, err)
}
