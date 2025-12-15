package liquidity_provider_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/alerts"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/penalization"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPenalizationAlertUseCase_Run(t *testing.T) {
	collateral := &mocks.CollateralManagementContractMock{}
	events := []penalization.PenalizedEvent{
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
	collateral.On(
		"GetPenalizedEvents",
		test.AnyCtx,
		uint64(5),
		&toBlock,
	).Return(events, nil).Once()

	repo := mocks.NewPenalizedEventRepositoryMock(t)
	repo.On("InsertPenalization", mock.Anything, mock.Anything).Return(nil)

	sender := &mocks.AlertSenderMock{}
	recipient := "recipient@test.com"

	for i := 0; i < 3; i++ {
		sender.On(
			"SendAlert",
			test.AnyCtx,
			alerts.AlertSubjectPenalization,
			fmt.Sprintf("You were punished in %v rBTC for the quoteHash %s", events[i].Penalty.ToRbtc(), events[i].QuoteHash),
			[]string{recipient},
		).Return(nil).Once()
	}

	contracts := blockchain.RskContracts{CollateralManagement: collateral}
	useCase := liquidity_provider.NewPenalizationAlertUseCase(contracts, sender, recipient, repo)
	err := useCase.Run(context.Background(), 5, 10)
	require.NoError(t, err)
	collateral.AssertExpectations(t)
	sender.AssertExpectations(t)
}

func TestPenalizationAlertUseCase_Run_GetEvents(t *testing.T) {
	collateral := &mocks.CollateralManagementContractMock{}
	sender := &mocks.AlertSenderMock{}
	collateral.On("GetPenalizedEvents", test.AnyCtx, uint64(5), mock.Anything).
		Return([]penalization.PenalizedEvent{}, assert.AnError).Once()
	contracts := blockchain.RskContracts{CollateralManagement: collateral}
	repo := mocks.NewPenalizedEventRepositoryMock(t)
	useCase := liquidity_provider.NewPenalizationAlertUseCase(contracts, sender, "recipient", repo)
	err := useCase.Run(context.Background(), 5, 10)
	collateral.AssertExpectations(t)
	require.Error(t, err)
}
