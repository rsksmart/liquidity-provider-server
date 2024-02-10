package pegout_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestInitPegoutDepositCacheUseCase_Run(t *testing.T) {
	lbc := new(test.LbcMock)
	rpc := new(test.RskRpcMock)
	pegoutRepository := new(test.PegoutQuoteRepositoryMock)
	height := uint64(10)
	rpc.On("GetHeight", context.Background()).Return(height, nil)
	events := []quote.PegoutDeposit{
		{
			TxHash:      "0x123456",
			QuoteHash:   "0x654321",
			Amount:      entities.NewWei(1),
			Timestamp:   time.Now(),
			BlockNumber: 6,
			From:        "0xabcdef",
		},
		{
			TxHash:      "0x987654",
			QuoteHash:   "0x445566",
			Amount:      entities.NewWei(2),
			Timestamp:   time.Now(),
			BlockNumber: 7,
			From:        "0xabcdef",
		},
	}
	lbc.On("GetDepositEvents", context.Background(), uint64(5), &height).Return(events, nil)
	pegoutRepository.On("UpsertPegoutDeposits", context.Background(), events).Return(nil)
	useCase := pegout.NewInitPegoutDepositCacheUseCase(pegoutRepository, lbc, rpc)
	err := useCase.Run(context.Background(), 5)
	rpc.AssertExpectations(t)
	lbc.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	assert.Nil(t, err)
}

func TestInitPegoutDepositCacheUseCase_Run_ErrorHandling(t *testing.T) {
	cases := test.Table[func(lbc *test.LbcMock, quoteRepository *test.PegoutQuoteRepositoryMock, rpc *test.RskRpcMock), error]{
		{
			Value: func(lbc *test.LbcMock, quoteRepository *test.PegoutQuoteRepositoryMock, rpc *test.RskRpcMock) {
				rpc.On("GetHeight", context.Background()).Return(uint64(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *test.LbcMock, quoteRepository *test.PegoutQuoteRepositoryMock, rpc *test.RskRpcMock) {
				rpc.On("GetHeight", context.Background()).Return(uint64(10), nil)
				lbc.On("GetDepositEvents", context.Background(), mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
		},
		{
			Value: func(lbc *test.LbcMock, quoteRepository *test.PegoutQuoteRepositoryMock, rpc *test.RskRpcMock) {
				rpc.On("GetHeight", context.Background()).Return(uint64(10), nil)
				lbc.On("GetDepositEvents", context.Background(), uint64(5), mock.Anything).Return([]quote.PegoutDeposit{}, nil)
				quoteRepository.On("UpsertPegoutDeposits", context.Background(), mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, c := range cases {
		lbc := new(test.LbcMock)
		quoteRepository := new(test.PegoutQuoteRepositoryMock)
		rsk := new(test.RskRpcMock)
		c.Value(lbc, quoteRepository, rsk)
		useCase := pegout.NewInitPegoutDepositCacheUseCase(quoteRepository, lbc, rsk)
		err := useCase.Run(context.Background(), 5)
		lbc.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		rsk.AssertExpectations(t)
		assert.NotNil(t, err)
	}
}
