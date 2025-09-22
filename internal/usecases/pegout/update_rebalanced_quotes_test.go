package pegout_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/rootstock"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

// nolint:funlen
func TestUpdateBtcReleaseUseCase_Run(t *testing.T) {
	eventMock := rootstock.BatchPegOut{
		TransactionHash:    test.AnyRskAddress,
		BlockHash:          test.AnyString,
		BlockNumber:        123,
		BtcTxHash:          test.AnyHash,
		ReleaseRskTxHashes: []string{test.AnyRskAddress},
	}
	retainedQuotes := []quote.RetainedPegoutQuote{
		{QuoteHash: "q1", State: quote.PegoutStateBridgeTxSucceeded, BtcReleaseTxHash: test.AnyHash},
		{QuoteHash: "q2", State: quote.PegoutStateBridgeTxSucceeded, BtcReleaseTxHash: test.AnyHash},
		{QuoteHash: "q3", State: quote.PegoutStateBridgeTxSucceeded, BtcReleaseTxHash: test.AnyHash},
	}
	t.Run("should run successfully", func(t *testing.T) {
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		eventBus := &mocks.EventBusMock{}
		batchPegOutRepository := &mocks.BatchPegOutRepositoryMock{}

		pegoutRepository.EXPECT().GetRetainedQuotesInBatch(mock.Anything, eventMock).Return(retainedQuotes, nil).Once()
		pegoutRepository.EXPECT().UpdateRetainedQuotes(mock.Anything, mock.MatchedBy(func(quotes []quote.RetainedPegoutQuote) bool {
			for _, q := range quotes {
				if q.State != quote.PegoutStateBtcReleased || q.BtcReleaseTxHash != eventMock.TransactionHash {
					return false
				}
			}
			return true
		})).Return(nil).Once()
		batchPegOutRepository.EXPECT().UpsertBatch(mock.Anything, eventMock).Return(nil).Once()
		eventBus.On("Publish", mock.MatchedBy(func(event rootstock.BatchPegOutUpdatedEvent) bool {
			return event.Event.Id() == rootstock.BatchPegOutUpdatedEventId &&
				event.BatchPegOut.TransactionHash == eventMock.TransactionHash &&
				len(event.QuoteHashes) == len(retainedQuotes)
		}))
		useCase := pegout.NewUpdateBtcReleaseUseCase(pegoutRepository, batchPegOutRepository, eventBus)
		result, err := useCase.Run(context.Background(), eventMock)
		require.NoError(t, err)
		require.Equal(t, uint(len(retainedQuotes)), result)
		eventBus.AssertExpectations(t)
		pegoutRepository.AssertExpectations(t)
		batchPegOutRepository.AssertExpectations(t)
	})
	t.Run("should return 0 if no quotes found", func(t *testing.T) {
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		eventBus := &mocks.EventBusMock{}
		batchPegOutRepository := &mocks.BatchPegOutRepositoryMock{}
		pegoutRepository.EXPECT().GetRetainedQuotesInBatch(mock.Anything, eventMock).Return([]quote.RetainedPegoutQuote{}, nil).Once()
		useCase := pegout.NewUpdateBtcReleaseUseCase(pegoutRepository, batchPegOutRepository, eventBus)
		result, err := useCase.Run(context.Background(), eventMock)
		require.NoError(t, err)
		require.Zero(t, result)
		pegoutRepository.AssertExpectations(t)
		batchPegOutRepository.AssertNotCalled(t, "UpsertBatch")
		eventBus.AssertNotCalled(t, "Publish")
	})
	t.Run("should handle error reading from db", func(t *testing.T) {
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		eventBus := &mocks.EventBusMock{}
		batchPegOutRepository := &mocks.BatchPegOutRepositoryMock{}
		pegoutRepository.EXPECT().GetRetainedQuotesInBatch(mock.Anything, eventMock).Return(nil, assert.AnError).Once()
		useCase := pegout.NewUpdateBtcReleaseUseCase(pegoutRepository, batchPegOutRepository, eventBus)
		result, err := useCase.Run(context.Background(), eventMock)
		require.Error(t, err)
		require.Zero(t, result)
		pegoutRepository.AssertExpectations(t)
		batchPegOutRepository.AssertNotCalled(t, "UpsertBatch")
		eventBus.AssertNotCalled(t, "Publish")
	})
	t.Run("should handle error updating quotes", func(t *testing.T) {
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		eventBus := &mocks.EventBusMock{}
		batchPegOutRepository := &mocks.BatchPegOutRepositoryMock{}
		pegoutRepository.EXPECT().GetRetainedQuotesInBatch(mock.Anything, eventMock).Return(retainedQuotes, nil).Once()
		pegoutRepository.EXPECT().UpdateRetainedQuotes(mock.Anything, mock.Anything).Return(assert.AnError).Once()
		useCase := pegout.NewUpdateBtcReleaseUseCase(pegoutRepository, batchPegOutRepository, eventBus)
		result, err := useCase.Run(context.Background(), eventMock)
		require.Error(t, err)
		require.Zero(t, result)
		pegoutRepository.AssertExpectations(t)
		batchPegOutRepository.AssertNotCalled(t, "UpsertBatch")
		eventBus.AssertNotCalled(t, "Publish")
	})
	t.Run("should handle error upserting event", func(t *testing.T) {
		pegoutRepository := &mocks.PegoutQuoteRepositoryMock{}
		eventBus := &mocks.EventBusMock{}
		batchPegOutRepository := &mocks.BatchPegOutRepositoryMock{}
		pegoutRepository.EXPECT().GetRetainedQuotesInBatch(mock.Anything, eventMock).Return(retainedQuotes, nil).Once()
		pegoutRepository.EXPECT().UpdateRetainedQuotes(mock.Anything, mock.Anything).Return(nil).Once()
		batchPegOutRepository.EXPECT().UpsertBatch(mock.Anything, mock.Anything).Return(assert.AnError).Once()
		useCase := pegout.NewUpdateBtcReleaseUseCase(pegoutRepository, batchPegOutRepository, eventBus)
		result, err := useCase.Run(context.Background(), eventMock)
		require.Error(t, err)
		require.Equal(t, uint(len(retainedQuotes)), result)
		pegoutRepository.AssertExpectations(t)
		batchPegOutRepository.AssertExpectations(t)
		eventBus.AssertNotCalled(t, "Publish")
	})
}
