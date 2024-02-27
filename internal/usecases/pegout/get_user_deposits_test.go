package pegout_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetUserDepositsUseCase_Run(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	deposit := quote.PegoutDeposit{
		TxHash:      "0x123456",
		QuoteHash:   "0x654321",
		Amount:      entities.NewWei(1),
		Timestamp:   time.Now(),
		BlockNumber: 6,
	}
	quoteRepository.On(
		"ListPegoutDepositsByAddress",
		mock.AnythingOfType("context.backgroundCtx"),
		"0x123456",
	).Return([]quote.PegoutDeposit{deposit}, nil)
	useCase := pegout.NewGetUserDepositsUseCase(quoteRepository)
	result, err := useCase.Run(context.Background(), "0x123456")
	require.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, result[0], deposit)
}

func TestGetUserDepositsUseCase_Run_ErrorHandling(t *testing.T) {
	quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
	quoteRepository.On(
		"ListPegoutDepositsByAddress",
		mock.AnythingOfType("context.backgroundCtx"),
		"0x123456",
	).Return(nil, assert.AnError)
	useCase := pegout.NewGetUserDepositsUseCase(quoteRepository)
	result, err := useCase.Run(context.Background(), "0x123456")
	require.Error(t, err)
	assert.Nil(t, result)
}
