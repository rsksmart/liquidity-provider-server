package handlers_test

import (
	"encoding/json"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/entrypoints/rest/handlers"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lpEntity "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/pkg"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewProviderDetailsHandler(t *testing.T) {
	const (
		path            = "/providers/details"
		verb            = "GET"
		captchaKey      = "captchaKey"
		captchaDisabled = true
	)

	providerMock := &mocks.ProviderMock{}
	providerMock.On("GeneralConfiguration", mock.Anything).Return(lpEntity.GeneralConfiguration{
		RskConfirmations:     map[string]uint16{"1": 10, "2": 20, "3": 50, "4": 15},
		BtcConfirmations:     map[string]uint16{"1": 15, "2": 11, "3": 14, "4": 11},
		PublicLiquidityCheck: true,
	}).Times(5)
	providerMock.On("PeginConfiguration", mock.Anything).Return(lpEntity.PeginConfiguration{
		TimeForDeposit: 300,
		CallTime:       400,
		PenaltyFee:     entities.NewWei(500),
		FixedFee:       entities.NewWei(700),
		FeePercentage:  utils.NewBigFloat64(15.77),
		MaxValue:       entities.NewWei(800),
		MinValue:       entities.NewWei(100),
	}).Times(5)
	providerMock.On("PegoutConfiguration", mock.Anything).Return(lpEntity.PegoutConfiguration{
		TimeForDeposit:       111,
		ExpireTime:           222,
		PenaltyFee:           entities.NewWei(333),
		FixedFee:             entities.NewWei(444),
		FeePercentage:        utils.NewBigFloat64(0.33),
		MaxValue:             entities.NewWei(1000),
		MinValue:             entities.NewWei(10),
		ExpireBlocks:         500,
		BridgeTransactionMin: entities.NewWei(1500),
	}).Times(5)

	t.Run("should return 200 on success", func(t *testing.T) {
		useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, captchaDisabled, true, providerMock, providerMock, providerMock)
		handler := handlers.NewProviderDetailsHandler(useCase)
		assert.HTTPSuccess(t, handler, verb, path, nil)
		assert.HTTPBodyContains(t, handler, verb, path, nil, `{"siteKey":"captchaKey","liquidityCheckEnabled":true,"usingSegwitFederation":true,"pegin":{"fee":700,"fixedFee":700,"feePercentage":15.77,"minTransactionValue":100,"maxTransactionValue":800,"requiredConfirmations":15},"pegout":{"fee":444,"fixedFee":444,"feePercentage":0.33,"minTransactionValue":10,"maxTransactionValue":1000,"requiredConfirmations":50}}`)
	})
	t.Run("should handle internal error", func(t *testing.T) {
		useCase := liquidity_provider.NewGetDetailUseCase("", false, true, providerMock, providerMock, providerMock)
		handler := handlers.NewProviderDetailsHandler(useCase)
		assert.HTTPStatusCode(t, handler, verb, path, nil, http.StatusInternalServerError)
		assert.HTTPBodyContains(t, handler, verb, path, nil, `"details":{"error":"ProviderDetail: missing captcha key"}`)
	})
	t.Run("should return deprecated fee field", func(t *testing.T) {
		var result pkg.ProviderDetailResponse
		useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, captchaDisabled, true, providerMock, providerMock, providerMock)
		handler := handlers.NewProviderDetailsHandler(useCase)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(verb, path, nil))
		require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &result))
		// disable linter to be able to check over the deprecated field
		// nolint:staticcheck
		assert.Equal(t, result.Pegout.FixedFee, result.Pegout.Fee)
		// nolint:staticcheck
		assert.Equal(t, result.Pegin.FixedFee, result.Pegin.Fee)
	})
	providerMock.AssertExpectations(t)
}
