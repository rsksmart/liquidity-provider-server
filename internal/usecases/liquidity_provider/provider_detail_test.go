package liquidity_provider_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDetailUseCase_Run(t *testing.T) {
	provider := &test.ProviderMock{}
	provider.On("CallFeePegin").Return(entities.NewWei(100))
	provider.On("MinPegin").Return(entities.NewWei(1000))
	provider.On("MaxPegin").Return(entities.NewWei(10000))
	provider.On("MaxPeginConfirmations").Return(uint16(10))
	provider.On("CallFeePegout").Return(entities.NewWei(200))
	provider.On("MinPegout").Return(entities.NewWei(2000))
	provider.On("MaxPegout").Return(entities.NewWei(20000))
	provider.On("MaxPegoutConfirmations").Return(uint16(20))
	captchaKey := "testKey"
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, provider, provider)
	result, err := useCase.Run()
	assert.Nil(t, err)
	assert.Equal(t, result, liquidity_provider.FullLiquidityProvider{
		SiteKey: captchaKey,
		Pegin: entities.LiquidityProviderDetail{
			Fee:                   entities.NewWei(100),
			MinTransactionValue:   entities.NewWei(1000),
			MaxTransactionValue:   entities.NewWei(10000),
			RequiredConfirmations: 10,
		},
		Pegout: entities.LiquidityProviderDetail{
			Fee:                   entities.NewWei(200),
			MinTransactionValue:   entities.NewWei(2000),
			MaxTransactionValue:   entities.NewWei(20000),
			RequiredConfirmations: 20,
		},
	})
}

func TestGetDetailUseCase_Run_InvalidCaptchaKey(t *testing.T) {
	provider := &test.ProviderMock{}
	provider.On("CallFeePegin").Return(entities.NewWei(100))
	provider.On("MinPegin").Return(entities.NewWei(1000))
	provider.On("MaxPegin").Return(entities.NewWei(10000))
	provider.On("MaxPeginConfirmations").Return(uint16(10))
	provider.On("CallFeePegout").Return(entities.NewWei(200))
	provider.On("MinPegout").Return(entities.NewWei(2000))
	provider.On("MaxPegout").Return(entities.NewWei(20000))
	provider.On("MaxPegoutConfirmations").Return(uint16(20))
	captchaKey := ""
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, provider, provider)
	_, err := useCase.Run()
	assert.Equal(t, "ProviderDetail: missing captcha key", err.Error())
}

func TestGetDetailUseCase_Run_InvalidPeginDetail(t *testing.T) {
	var nilWei *entities.Wei
	provider := &test.ProviderMock{}
	provider.On("CallFeePegin").Return(entities.NewWei(100))
	provider.On("MinPegin").Return(nilWei)
	provider.On("MaxPegin").Return(entities.NewWei(10000))
	provider.On("MaxPeginConfirmations").Return(uint16(10))
	provider.On("CallFeePegout").Return(entities.NewWei(200))
	provider.On("MinPegout").Return(entities.NewWei(2000))
	provider.On("MaxPegout").Return(entities.NewWei(20000))
	provider.On("MaxPegoutConfirmations").Return(uint16(20))
	captchaKey := "testKey"
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, provider, provider)
	_, err := useCase.Run()
	assert.Equal(t, "ProviderDetail: Key: 'LiquidityProviderDetail.MinTransactionValue' "+
		"Error:Field validation for 'MinTransactionValue' failed on the 'required' tag", err.Error())
}

func TestGetDetailUseCase_Run_InvalidPegoutDetail(t *testing.T) {
	var nilWei *entities.Wei
	provider := &test.ProviderMock{}
	provider.On("CallFeePegin").Return(entities.NewWei(100))
	provider.On("MinPegin").Return(entities.NewWei(1000))
	provider.On("MaxPegin").Return(entities.NewWei(10000))
	provider.On("MaxPeginConfirmations").Return(uint16(10))
	provider.On("CallFeePegout").Return(entities.NewWei(200))
	provider.On("MinPegout").Return(nilWei)
	provider.On("MaxPegout").Return(entities.NewWei(20000))
	provider.On("MaxPegoutConfirmations").Return(uint16(20))
	captchaKey := "testKey"
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, provider, provider)
	_, err := useCase.Run()
	assert.Equal(t, "ProviderDetail: Key: 'LiquidityProviderDetail.MinTransactionValue' "+
		"Error:Field validation for 'MinTransactionValue' failed on the 'required' tag", err.Error())
}
