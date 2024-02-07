package liquidity_provider_test

import (
	"context"
	"fmt"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestPenalizationAlertUseCase_Run(t *testing.T) {
	lbc := &test.LbcMock{}
	events := []entities.PunishmentEvent{
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

	sender := &test.AlertSenderMock{}
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
	assert.Nil(t, err)
	lbc.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestPenalizationAlertUseCase_Run_GetEvents(t *testing.T) {
	lbc := &test.LbcMock{}
	sender := &test.AlertSenderMock{}
	lbc.On("GetPeginPunishmentEvents", mock.AnythingOfType("context.backgroundCtx"), uint64(5), mock.Anything).
		Return([]entities.PunishmentEvent{}, assert.AnError).Once()
	useCase := liquidity_provider.NewPenalizationAlertUseCase(lbc, sender, "recipient")
	err := useCase.Run(context.Background(), 5, 10)
	lbc.AssertExpectations(t)
	assert.NotNil(t, err)
}
