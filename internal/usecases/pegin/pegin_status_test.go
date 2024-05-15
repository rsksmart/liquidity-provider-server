package pegin_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStatusUseCase_Run(t *testing.T) {
	const quoteHash = "quoteHash"
	retainedPeginQuote := quote.RetainedPeginQuote{
		QuoteHash:           quoteHash,
		DepositAddress:      test.AnyAddress,
		Signature:           test.AnyString,
		RequiredLiquidity:   entities.NewWei(500),
		State:               quote.PeginStateCallForUserSucceeded,
		UserBtcTxHash:       "btc tx hash",
		CallForUserTxHash:   "cfu tx hash",
		RegisterPeginTxHash: "register pegin tx hash",
	}
	t.Run("Get status of a pegin quote", func(t *testing.T) {
		repo := new(mocks.PeginQuoteRepositoryMock)
		useCase := pegin.NewStatusUseCase(repo)
		repo.On("GetQuote", context.Background(), quoteHash).Return(&testPeginQuote, nil).Once()
		repo.On("GetRetainedQuote", context.Background(), quoteHash).Return(&retainedPeginQuote, nil).Once()
		result, err := useCase.Run(context.Background(), quoteHash)
		require.NoError(t, err)
		require.Equal(t, quote.WatchedPeginQuote{
			PeginQuote:    testPeginQuote,
			RetainedQuote: retainedPeginQuote,
		}, result)
	})
	t.Run("Return not found error", func(t *testing.T) {
		repo := new(mocks.PeginQuoteRepositoryMock)
		useCase := pegin.NewStatusUseCase(repo)
		repo.On("GetQuote", context.Background(), quoteHash).Return(nil, nil).Once()
		result, err := useCase.Run(context.Background(), quoteHash)
		require.ErrorIs(t, err, usecases.QuoteNotFoundError)
		require.Empty(t, result)
	})
	t.Run("Return not accepted error", func(t *testing.T) {
		repo := new(mocks.PeginQuoteRepositoryMock)
		useCase := pegin.NewStatusUseCase(repo)
		repo.On("GetQuote", context.Background(), quoteHash).Return(&testPeginQuote, nil).Once()
		repo.On("GetRetainedQuote", context.Background(), quoteHash).Return(nil, nil).Once()
		result, err := useCase.Run(context.Background(), quoteHash)
		require.ErrorIs(t, err, usecases.QuoteNotAcceptedError)
		require.Empty(t, result)
	})
	t.Run("Handle database errors", func(t *testing.T) {
		repo := new(mocks.PeginQuoteRepositoryMock)
		useCase := pegin.NewStatusUseCase(repo)

		repo.On("GetQuote", context.Background(), quoteHash).Return(nil, assert.AnError).Once()
		result, err := useCase.Run(context.Background(), quoteHash)
		require.Error(t, err)
		require.Empty(t, result)

		repo.On("GetQuote", context.Background(), quoteHash).Return(&testPeginQuote, nil).Once()
		repo.On("GetRetainedQuote", context.Background(), quoteHash).Return(nil, assert.AnError).Once()
		result, err = useCase.Run(context.Background(), quoteHash)
		require.Error(t, err)
		require.Empty(t, result)
	})
}
