package liquidity_provider_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetDetailUseCase_Run(t *testing.T) {
	provider := &mocks.ProviderMock{}
	prepareDetailMock(provider)
	captchaKey := "testKey"
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, provider, provider, provider)
	result, err := useCase.Run(context.Background())
	require.NoError(t, err)
	assert.Equal(t, liquidity_provider.FullLiquidityProvider{
		SiteKey: captchaKey,
		Pegin: lp.LiquidityProviderDetail{
			Fee:                   entities.NewWei(100),
			MinTransactionValue:   entities.NewWei(1000),
			MaxTransactionValue:   entities.NewWei(10000),
			RequiredConfirmations: 10,
		},
		Pegout: lp.LiquidityProviderDetail{
			Fee:                   entities.NewWei(200),
			MinTransactionValue:   entities.NewWei(2000),
			MaxTransactionValue:   entities.NewWei(20000),
			RequiredConfirmations: 20,
		},
	}, result)
}

func TestGetDetailUseCase_Run_InvalidCaptchaKey(t *testing.T) {
	provider := &mocks.ProviderMock{}
	prepareDetailMock(provider)
	captchaKey := ""
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, provider, provider, provider)
	_, err := useCase.Run(context.Background())
	assert.Equal(t, "ProviderDetail: missing captcha key", err.Error())
}

func TestGetDetailUseCase_Run_InvalidPeginDetail(t *testing.T) {
	var nilWei *entities.Wei
	provider := &mocks.ProviderMock{}
	ctx := context.Background()
	prepareDetailMock(provider)
	config := provider.PeginConfiguration(ctx)
	config.MinValue = nilWei
	provider.On("PeginConfiguration", mock.AnythingOfType("context.backgroundCtx")).Return(config)
	captchaKey := "testKey"
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, provider, provider, provider)
	_, err := useCase.Run(ctx)
	assert.Equal(t, "ProviderDetail: Key: 'LiquidityProviderDetail.MinTransactionValue' "+
		"Error:Field validation for 'MinTransactionValue' failed on the 'required' tag", err.Error())
}

func TestGetDetailUseCase_Run_InvalidPegoutDetail(t *testing.T) {
	var nilWei *entities.Wei
	provider := &mocks.ProviderMock{}
	ctx := context.Background()
	prepareDetailMock(provider)
	config := provider.PegoutConfiguration(ctx)
	config.MinValue = nilWei
	provider.On("PegoutConfiguration", mock.AnythingOfType("context.backgroundCtx")).Return(config)
	captchaKey := "testKey"
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, provider, provider, provider)
	_, err := useCase.Run(ctx)
	assert.Equal(t, "ProviderDetail: Key: 'LiquidityProviderDetail.MinTransactionValue' "+
		"Error:Field validation for 'MinTransactionValue' failed on the 'required' tag", err.Error())
}

func prepareDetailMock(provider *mocks.ProviderMock) {
	provider.On("PeginConfiguration", test.AnyCtx).Return(lp.PeginConfiguration{
		TimeForDeposit: lp.PeginTimeForDeposit,
		CallTime:       lp.PeginCallTime,
		PenaltyFee:     entities.NewWei(lp.PeginPenaltyFee),
		CallFee:        entities.NewWei(100),
		MaxValue:       entities.NewWei(10000),
		MinValue:       entities.NewWei(1000),
	}).Once()
	provider.On("PegoutConfiguration", test.AnyCtx).Return(lp.PegoutConfiguration{
		TimeForDeposit: lp.PegoutTimeForDeposit,
		CallTime:       lp.PegoutCallTime,
		PenaltyFee:     entities.NewWei(lp.PegoutPenaltyFee),
		CallFee:        entities.NewWei(200),
		MaxValue:       entities.NewWei(20000),
		MinValue:       entities.NewWei(2000),
		ExpireBlocks:   lp.PegoutExpireBlocks,
	}).Once()
	provider.On("GeneralConfiguration", test.AnyCtx).
		Return(lp.GeneralConfiguration{
			RskConfirmations: map[int]uint16{1: 20},
			BtcConfirmations: map[int]uint16{1: 10},
		}).Once()
}
