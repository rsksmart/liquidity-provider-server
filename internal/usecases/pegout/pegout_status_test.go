package pegout_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegout"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStatusUseCase_Run(t *testing.T) {
	const quoteHash = "quoteHash"
	retainedPegoutQuote := quote.RetainedPegoutQuote{
		QuoteHash:          quoteHash,
		DepositAddress:     test.AnyAddress,
		Signature:          test.AnyString,
		RequiredLiquidity:  entities.NewWei(500),
		State:              quote.PegoutStateRefundPegOutSucceeded,
		UserRskTxHash:      "rsk tx hash",
		LpBtcTxHash:        "lp btc tx hash",
		RefundPegoutTxHash: "refund tx hash",
		BridgeRefundTxHash: "bridge tx hash",
	}
	creationData := quote.PegoutCreationData{FeeRate: utils.NewBigFloat64(100.5), FeePercentage: utils.NewBigFloat64(0.5), GasPrice: entities.NewWei(100), FixedFee: entities.NewWei(100)}
	t.Run("Get status of a pegout quote", func(t *testing.T) {
		repo := new(mocks.PegoutQuoteRepositoryMock)
		useCase := pegout.NewStatusUseCase(repo)
		repo.On("GetQuote", context.Background(), quoteHash).Return(&pegoutQuote, nil).Once()
		repo.On("GetRetainedQuote", context.Background(), quoteHash).Return(&retainedPegoutQuote, nil).Once()
		repo.EXPECT().GetPegoutCreationData(context.Background(), quoteHash).Return(creationData).Once()
		result, err := useCase.Run(context.Background(), quoteHash)
		require.NoError(t, err)
		require.Equal(t, quote.WatchedPegoutQuote{
			PegoutQuote:   pegoutQuote,
			RetainedQuote: retainedPegoutQuote,
			CreationData:  creationData,
		}, result)
	})
	t.Run("Return not found error", func(t *testing.T) {
		repo := new(mocks.PegoutQuoteRepositoryMock)
		useCase := pegout.NewStatusUseCase(repo)
		repo.On("GetQuote", context.Background(), quoteHash).Return(nil, nil).Once()
		result, err := useCase.Run(context.Background(), quoteHash)
		require.ErrorIs(t, err, usecases.QuoteNotFoundError)
		require.Empty(t, result)
	})
	t.Run("Return not accepted error", func(t *testing.T) {
		repo := new(mocks.PegoutQuoteRepositoryMock)
		useCase := pegout.NewStatusUseCase(repo)
		repo.On("GetQuote", context.Background(), quoteHash).Return(&pegoutQuote, nil).Once()
		repo.On("GetRetainedQuote", context.Background(), quoteHash).Return(nil, nil).Once()
		result, err := useCase.Run(context.Background(), quoteHash)
		require.ErrorIs(t, err, usecases.QuoteNotAcceptedError)
		require.Empty(t, result)
	})
	t.Run("Handle database errors", func(t *testing.T) {
		repo := new(mocks.PegoutQuoteRepositoryMock)
		useCase := pegout.NewStatusUseCase(repo)

		repo.On("GetQuote", context.Background(), quoteHash).Return(nil, assert.AnError).Once()
		result, err := useCase.Run(context.Background(), quoteHash)
		require.Error(t, err)
		require.Empty(t, result)

		repo.On("GetQuote", context.Background(), quoteHash).Return(&pegoutQuote, nil).Once()
		repo.On("GetRetainedQuote", context.Background(), quoteHash).Return(nil, assert.AnError).Once()
		result, err = useCase.Run(context.Background(), quoteHash)
		require.Error(t, err)
		require.Empty(t, result)
	})
}
