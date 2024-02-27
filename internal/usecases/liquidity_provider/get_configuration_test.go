package liquidity_provider_test

import (
	"context"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetConfigUseCase_Run(t *testing.T) {
	lpMock := &mocks.ProviderMock{}
	lpMock.On("GeneralConfiguration", test.AnyCtx).Return(lp.DefaultGeneralConfiguration())
	lpMock.On("PeginConfiguration", test.AnyCtx).Return(lp.DefaultPeginConfiguration())
	lpMock.On("PegoutConfiguration", test.AnyCtx).Return(lp.DefaultPegoutConfiguration())
	useCase := liquidity_provider.NewGetConfigUseCase(lpMock, lpMock, lpMock)
	result := useCase.Run(context.Background())
	lpMock.AssertExpectations(t)
	require.Equal(t, liquidity_provider.FullConfiguration{
		General: lp.DefaultGeneralConfiguration(),
		Pegin:   lp.DefaultPeginConfiguration(),
		Pegout:  lp.DefaultPegoutConfiguration(),
	}, result)
}
