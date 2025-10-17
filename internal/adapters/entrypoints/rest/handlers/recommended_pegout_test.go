package handlers_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"testing"
)

func TestNewRecommendedPegoutHandler(t *testing.T) {
	const path = "/pegout/recommended"
	queryOnlyAmount := url.Values{"amount": []string{"100"}}
	queryFull := url.Values{"amount": []string{"500"}, "destination_type": []string{"p2pkh"}}
	t.Run("should execute use case successfully", func(t *testing.T) {
		useCase := new(mocks.RecommendedPegoutUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, entities.NewWei(500), blockchain.BtcAddressTypeP2PKH).
			Return(usecases.RecommendedOperationResult{RecommendedQuoteValue: entities.NewWei(495), EstimatedCallFee: entities.NewWei(1), EstimatedGasFee: entities.NewWei(2), EstimatedProductFee: entities.NewWei(3)}, nil)
		handler := handlers.NewRecommendedPegoutHandler(useCase)
		assert.HTTPSuccess(t, handler, http.MethodGet, path, queryFull)
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, queryFull, `{"recommendedQuoteValue":495,"estimatedCallFee":1,"estimatedGasFee":2,"estimatedProductFee":3}`)
	})
	t.Run("should execute use case without the destination type", func(t *testing.T) {
		useCase := new(mocks.RecommendedPegoutUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, entities.NewWei(100), blockchain.BtcAddressType("")).
			Return(usecases.RecommendedOperationResult{RecommendedQuoteValue: entities.NewWei(90), EstimatedCallFee: entities.NewWei(3), EstimatedGasFee: entities.NewWei(5), EstimatedProductFee: entities.NewWei(6)}, nil)
		handler := handlers.NewRecommendedPegoutHandler(useCase)
		assert.HTTPSuccess(t, handler, http.MethodGet, path, queryOnlyAmount)
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, queryOnlyAmount, `{"recommendedQuoteValue":90,"estimatedCallFee":3,"estimatedGasFee":5,"estimatedProductFee":6}`)
	})
	t.Run("shoul return 400 on invalid query params", func(t *testing.T) {
		testCases := []url.Values{
			{"amount": []string{""}, "destination_type": []string{"p2pkh"}},
			{"amount": []string{"1"}, "destination_type": []string{"test"}},
			{"amount": []string{"a"}, "destination_type": []string{"p2pkh"}},
		}
		for _, testCase := range testCases {
			useCase := new(mocks.RecommendedPegoutUseCaseMock)
			handler := handlers.NewRecommendedPegoutHandler(useCase)
			assert.HTTPStatusCode(t, handler, http.MethodGet, path, testCase, http.StatusBadRequest)
			useCase.AssertNotCalled(t, "Run")
		}
	})
	t.Run("should return 400 if recommended amount is out of limits", func(t *testing.T) {
		useCase := new(mocks.RecommendedPegoutUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything).Return(usecases.RecommendedOperationResult{}, liquidity_provider.AmountOutOfRangeError)
		handler := handlers.NewRecommendedPegoutHandler(useCase)
		assert.HTTPStatusCode(t, handler, http.MethodGet, path, queryFull, http.StatusBadRequest)
	})
	t.Run("should return 400 if there is no liquidity for the recommended amount", func(t *testing.T) {
		useCase := new(mocks.RecommendedPegoutUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything).Return(usecases.RecommendedOperationResult{}, usecases.NoLiquidityError)
		handler := handlers.NewRecommendedPegoutHandler(useCase)
		assert.HTTPStatusCode(t, handler, http.MethodGet, path, queryFull, http.StatusBadRequest)
	})
	t.Run("should return 400 if the recommended amount is below the bridge minimum", func(t *testing.T) {
		useCase := new(mocks.RecommendedPegoutUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything).Return(usecases.RecommendedOperationResult{}, usecases.TxBelowMinimumError)
		handler := handlers.NewRecommendedPegoutHandler(useCase)
		assert.HTTPStatusCode(t, handler, http.MethodGet, path, queryFull, http.StatusBadRequest)
	})
	t.Run("should return 500 for unkown errors", func(t *testing.T) {
		useCase := new(mocks.RecommendedPegoutUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything).Return(usecases.RecommendedOperationResult{}, assert.AnError)
		handler := handlers.NewRecommendedPegoutHandler(useCase)
		assert.HTTPStatusCode(t, handler, http.MethodGet, path, queryFull, http.StatusInternalServerError)
	})
}
