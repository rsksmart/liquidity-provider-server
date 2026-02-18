package handlers_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/url"
	"testing"
)

// nolint:funlen
func TestNewRecommendedPeginHandler(t *testing.T) {
	const path = "/pegin/recommended"
	var anyData = []byte{0x01, 0x02, 0x03}
	queryFull := url.Values{"amount": []string{"500"}, "destination_address": []string{test.AnyRskAddress}, "data": []string{"010203"}}
	t.Run("should execute use case successfully", func(t *testing.T) {
		useCase := new(mocks.RecommendedPeginUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, entities.NewWei(500), test.AnyRskAddress, anyData).
			Return(usecases.RecommendedOperationResult{RecommendedQuoteValue: entities.NewWei(495), EstimatedCallFee: entities.NewWei(1), EstimatedGasFee: entities.NewWei(2)}, nil)
		handler := handlers.NewRecommendedPeginHandler(useCase)
		assert.HTTPSuccess(t, handler, http.MethodGet, path, queryFull)
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, queryFull, `{"recommendedQuoteValue":495,"estimatedCallFee":1,"estimatedGasFee":2}`)
	})
	t.Run("should execute use case without the destination address", func(t *testing.T) {
		query := url.Values{"amount": []string{"500"}, "data": []string{"010203"}}
		useCase := new(mocks.RecommendedPeginUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, entities.NewWei(500), "", anyData).
			Return(usecases.RecommendedOperationResult{RecommendedQuoteValue: entities.NewWei(90), EstimatedCallFee: entities.NewWei(3), EstimatedGasFee: entities.NewWei(5)}, nil)
		handler := handlers.NewRecommendedPeginHandler(useCase)
		assert.HTTPSuccess(t, handler, http.MethodGet, path, query)
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, query, `{"recommendedQuoteValue":90,"estimatedCallFee":3,"estimatedGasFee":5}`)
	})
	t.Run("should execute use case without the data param", func(t *testing.T) {
		query := url.Values{"amount": []string{"500"}, "destination_address": []string{test.AnyRskAddress}}
		useCase := new(mocks.RecommendedPeginUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, entities.NewWei(500), test.AnyRskAddress, []byte{}).
			Return(usecases.RecommendedOperationResult{RecommendedQuoteValue: entities.NewWei(90), EstimatedCallFee: entities.NewWei(3), EstimatedGasFee: entities.NewWei(5)}, nil)
		handler := handlers.NewRecommendedPeginHandler(useCase)
		assert.HTTPSuccess(t, handler, http.MethodGet, path, query)
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, query, `{"recommendedQuoteValue":90,"estimatedCallFee":3,"estimatedGasFee":5}`)
	})
	t.Run("should support 0x prefix in the data", func(t *testing.T) {
		query := url.Values{"amount": []string{"500"}, "destination_address": []string{test.AnyRskAddress}, "data": []string{"0x010203"}}
		useCase := new(mocks.RecommendedPeginUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, entities.NewWei(500), test.AnyRskAddress, anyData).
			Return(usecases.RecommendedOperationResult{RecommendedQuoteValue: entities.NewWei(495), EstimatedCallFee: entities.NewWei(1), EstimatedGasFee: entities.NewWei(2)}, nil)
		handler := handlers.NewRecommendedPeginHandler(useCase)
		assert.HTTPSuccess(t, handler, http.MethodGet, path, query)
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, query, `{"recommendedQuoteValue":495,"estimatedCallFee":1,"estimatedGasFee":2}`)
	})
	t.Run("should return 400 on invalid query params", func(t *testing.T) {
		testCases := []url.Values{
			{"amount": []string{""}},
			{"amount": []string{"1"}, "data": []string{"test"}},
		}
		for _, testCase := range testCases {
			useCase := new(mocks.RecommendedPeginUseCaseMock)
			handler := handlers.NewRecommendedPeginHandler(useCase)
			assert.HTTPStatusCode(t, handler, http.MethodGet, path, testCase, http.StatusBadRequest)
			useCase.AssertNotCalled(t, "Run")
		}
	})
	t.Run("should return 400 if recommended amount is out of limits", func(t *testing.T) {
		useCase := new(mocks.RecommendedPeginUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(usecases.RecommendedOperationResult{}, liquidity_provider.AmountOutOfRangeError)
		handler := handlers.NewRecommendedPeginHandler(useCase)
		assert.HTTPStatusCode(t, handler, http.MethodGet, path, queryFull, http.StatusBadRequest)
	})
	t.Run("should return 400 if there is no liquidity for the recommended amount", func(t *testing.T) {
		useCase := new(mocks.RecommendedPeginUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(usecases.RecommendedOperationResult{}, usecases.NoLiquidityError)
		handler := handlers.NewRecommendedPeginHandler(useCase)
		assert.HTTPStatusCode(t, handler, http.MethodGet, path, queryFull, http.StatusBadRequest)
	})
	t.Run("should return 400 if the recommended amount is below the bridge minimum", func(t *testing.T) {
		useCase := new(mocks.RecommendedPeginUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(usecases.RecommendedOperationResult{}, usecases.TxBelowMinimumError)
		handler := handlers.NewRecommendedPeginHandler(useCase)
		assert.HTTPStatusCode(t, handler, http.MethodGet, path, queryFull, http.StatusBadRequest)
	})
	t.Run("should return 500 for unkown errors", func(t *testing.T) {
		useCase := new(mocks.RecommendedPeginUseCaseMock)
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(usecases.RecommendedOperationResult{}, assert.AnError)
		handler := handlers.NewRecommendedPeginHandler(useCase)
		assert.HTTPStatusCode(t, handler, http.MethodGet, path, queryFull, http.StatusInternalServerError)
	})
}
