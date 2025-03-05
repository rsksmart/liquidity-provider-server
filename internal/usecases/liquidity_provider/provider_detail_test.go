package liquidity_provider_test

import (
	"context"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/utils"
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
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, false, provider, provider, provider)
	result, err := useCase.Run(context.Background())
	require.NoError(t, err)
	assert.Equal(t, liquidity_provider.FullLiquidityProvider{
		SiteKey:               captchaKey,
		LiquidityCheckEnabled: true,
		Pegin: lp.LiquidityProviderDetail{
			FixedFee:              entities.NewWei(100),
			FeePercentage:         utils.NewBigFloat64(1.33),
			MinTransactionValue:   entities.NewWei(1000),
			MaxTransactionValue:   entities.NewWei(10000),
			RequiredConfirmations: 10,
		},
		Pegout: lp.LiquidityProviderDetail{
			FixedFee:              entities.NewWei(200),
			FeePercentage:         utils.NewBigFloat64(2.5),
			MinTransactionValue:   entities.NewWei(2000),
			MaxTransactionValue:   entities.NewWei(20000),
			RequiredConfirmations: 20,
		},
	}, result)
}

func TestGetDetailUseCase_Run_InvalidCaptchaKey(t *testing.T) {
	provider := &mocks.ProviderMock{}
	captchaKey := ""

	prepareDetailMock(provider)
	useCaseCaptchaEnabled := liquidity_provider.NewGetDetailUseCase(captchaKey, false, provider, provider, provider)
	_, err := useCaseCaptchaEnabled.Run(context.Background())
	assert.Equal(t, "ProviderDetail: missing captcha key", err.Error())

	prepareDetailMock(provider)
	useCaseCaptchaDisabled := liquidity_provider.NewGetDetailUseCase(captchaKey, true, provider, provider, provider)
	_, err = useCaseCaptchaDisabled.Run(context.Background())
	require.NoError(t, err)
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
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, false, provider, provider, provider)
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
	useCase := liquidity_provider.NewGetDetailUseCase(captchaKey, false, provider, provider, provider)
	_, err := useCase.Run(ctx)
	assert.Equal(t, "ProviderDetail: Key: 'LiquidityProviderDetail.MinTransactionValue' "+
		"Error:Field validation for 'MinTransactionValue' failed on the 'required' tag", err.Error())
}

func prepareDetailMock(provider *mocks.ProviderMock) {
	provider.On("PeginConfiguration", test.AnyCtx).Return(lp.PeginConfiguration{
		TimeForDeposit: lp.PeginTimeForDeposit,
		CallTime:       lp.PeginCallTime,
		PenaltyFee:     entities.NewWei(lp.PeginPenaltyFee),
		FixedFee:       entities.NewWei(100),
		FeePercentage:  utils.NewBigFloat64(1.33),
		MaxValue:       entities.NewWei(10000),
		MinValue:       entities.NewWei(1000),
	}).Once()
	provider.On("PegoutConfiguration", test.AnyCtx).Return(lp.PegoutConfiguration{
		TimeForDeposit: lp.PegoutTimeForDeposit,
		ExpireTime:     lp.PegoutExpireTime,
		PenaltyFee:     entities.NewWei(lp.PegoutPenaltyFee),
		FixedFee:       entities.NewWei(200),
		FeePercentage:  utils.NewBigFloat64(2.5),
		MaxValue:       entities.NewWei(20000),
		MinValue:       entities.NewWei(2000),
		ExpireBlocks:   lp.PegoutExpireBlocks,
	}).Once()
	provider.On("GeneralConfiguration", test.AnyCtx).
		Return(lp.GeneralConfiguration{
			PublicLiquidityCheck: true,
			RskConfirmations:     map[int]uint16{1: 20},
			BtcConfirmations:     map[int]uint16{1: 10},
		}).Once()
}
