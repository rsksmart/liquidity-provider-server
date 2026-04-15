package handlers_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
	t.Run("should return 400 on non-positive amount", func(t *testing.T) {
		testCases := []url.Values{
			{"amount": []string{"0"}},
			{"amount": []string{"-1000000"}},
		}
		for _, testCase := range testCases {
			useCase := new(mocks.RecommendedPeginUseCaseMock)
			handler := handlers.NewRecommendedPeginHandler(useCase)
			assert.HTTPStatusCode(t, handler, http.MethodGet, path, testCase, http.StatusBadRequest)
			assert.HTTPBodyContains(t, handler, http.MethodGet, path, testCase, "parameter amount must be greater than zero")
			useCase.AssertNotCalled(t, "Run")
		}
	})
	t.Run("should return 400 on invalid destination_address", func(t *testing.T) {
		testCases := []url.Values{
			{"amount": []string{"500"}, "destination_address": []string{"asd"}},
			{"amount": []string{"500"}, "destination_address": []string{"bc1q9ue5ls6zmzwdrhy6zucw9zwhz5zzv6qm2zn3mv"}},
			{"amount": []string{"500"}, "destination_address": []string{"0x31c1BB940B8b44bBf67a1Af40aab4eaB9268B5f"}},
			{"amount": []string{"500"}, "destination_address": []string{"0x31c1BB940B8b44bBf67a1Af40aab4eaB9268B5fb2"}},
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
	t.Run("should return 400 with structured details when effective amount is too low after fees", func(t *testing.T) {
		effectiveErr := usecases.NewEffectiveAmountTooLowError(
			entities.NewWei(5999310330477010),
			entities.NewWei(6000000000000000),
			entities.NewWei(6000689669522990),
		)
		useCase := new(mocks.RecommendedPeginUseCaseMock)
		// Return the error wrapped as the use case produces it, to exercise the errors.As unwrap path in the handler.
		useCase.EXPECT().Run(mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(usecases.RecommendedOperationResult{}, usecases.WrapUseCaseError(usecases.RecommendedPeginId, effectiveErr))
		handler := handlers.NewRecommendedPeginHandler(useCase)
		assert.HTTPStatusCode(t, handler, http.MethodGet, path, queryFull, http.StatusBadRequest)
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, queryFull, "Amount too low")
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, queryFull, `"effectiveAmount":5999310330477010`)
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, queryFull, `"minEffectiveAmount":6000000000000000`)
		assert.HTTPBodyContains(t, handler, http.MethodGet, path, queryFull, `"suggestedAmount":6000689669522990`)
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
