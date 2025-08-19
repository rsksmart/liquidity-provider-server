package pegout_test

import (
	"context"
	"testing"
	"time"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestInitPegoutDepositCacheUseCase_Run(t *testing.T) {
	lbc := new(mocks.LiquidityBridgeContractMock)
	rsk := new(mocks.RootstockRpcServerMock)
	pegoutRepository := new(mocks.PegoutQuoteRepositoryMock)
	height := uint64(10)
	rsk.On("GetHeight", context.Background()).Return(height, nil)
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
	contracts := blockchain.RskContracts{Lbc: lbc}
	rpc := blockchain.Rpc{Rsk: rsk}
	useCase := pegout.NewInitPegoutDepositCacheUseCase(pegoutRepository, contracts, rpc)
	err := useCase.Run(context.Background(), 5)
	rsk.AssertExpectations(t)
	lbc.AssertExpectations(t)
	pegoutRepository.AssertExpectations(t)
	require.NoError(t, err)
}

func TestInitPegoutDepositCacheUseCase_Run_ErrorHandling(t *testing.T) {
	cases := test.Table[func(lbc *mocks.LiquidityBridgeContractMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, rpc *mocks.RootstockRpcServerMock), error]{
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, rpc *mocks.RootstockRpcServerMock) {
				rpc.On("GetHeight", context.Background()).Return(uint64(0), assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, rpc *mocks.RootstockRpcServerMock) {
				rpc.On("GetHeight", context.Background()).Return(uint64(10), nil)
				lbc.On("GetDepositEvents", context.Background(), mock.Anything, mock.Anything).Return(nil, assert.AnError)
			},
		},
		{
			Value: func(lbc *mocks.LiquidityBridgeContractMock, quoteRepository *mocks.PegoutQuoteRepositoryMock, rpc *mocks.RootstockRpcServerMock) {
				rpc.On("GetHeight", context.Background()).Return(uint64(10), nil)
				lbc.On("GetDepositEvents", context.Background(), uint64(5), mock.Anything).Return([]quote.PegoutDeposit{}, nil)
				quoteRepository.On("UpsertPegoutDeposits", context.Background(), mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, c := range cases {
		lbc := new(mocks.LiquidityBridgeContractMock)
		quoteRepository := new(mocks.PegoutQuoteRepositoryMock)
		rsk := new(mocks.RootstockRpcServerMock)
		c.Value(lbc, quoteRepository, rsk)
		contracts := blockchain.RskContracts{Lbc: lbc}
		rpc := blockchain.Rpc{Rsk: rsk}
		useCase := pegout.NewInitPegoutDepositCacheUseCase(quoteRepository, contracts, rpc)
		err := useCase.Run(context.Background(), 5)
		lbc.AssertExpectations(t)
		quoteRepository.AssertExpectations(t)
		rsk.AssertExpectations(t)
		require.Error(t, err)
	}
}
