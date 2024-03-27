package liquidity_provider_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLiquidityStatusUseCase_Run_PublicLiquidityCheckDisabled(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	providerMock.On("GeneralConfiguration", mock.Anything).Return(lp.GeneralConfiguration{PublicLiquidityCheck: false})

	useCase := liquidity_provider.NewLiquidityStatusUseCase(blockchain.RskContracts{}, providerMock, blockchain.Rpc{Btc: new(mocks.BtcRpcMock), Rsk: new(mocks.RskRpcMock)}, &mocks.BtcWalletMock{}, providerMock.PeginLiquidityProvider, providerMock.PegoutLiquidityProvider)
	status, err := useCase.Run(context.Background())

	require.Nil(t, status)
	require.Error(t, err)
	require.Contains(t, err.Error(), "GetLiquidityStatusUseCase: balance checking is disabled")
}

func TestLiquidityStatusUseCase_Run_ErrorFetchingPeginLiquidity(t *testing.T) {
	providerMock := &mocks.ProviderMock{}
	peginProviderMock := &mocks.ProviderMock{}

	providerMock.On("GeneralConfiguration", mock.Anything).Return(lp.GeneralConfiguration{PublicLiquidityCheck: false})
	peginProviderMock.On("CalculateAvailablePeginLiquidity", mock.Anything).Return(0, errors.New("network error"))

	useCase := liquidity_provider.NewLiquidityStatusUseCase(blockchain.RskContracts{}, providerMock, blockchain.Rpc{Btc: new(mocks.BtcRpcMock), Rsk: new(mocks.RskRpcMock)}, &mocks.BtcWalletMock{}, providerMock.PeginLiquidityProvider, providerMock.PegoutLiquidityProvider)
	status, err := useCase.Run(context.Background())

	require.Nil(t, status)
	require.Error(t, err)
	require.Contains(t, err.Error(), "GetLiquidityStatusUseCase: balance checking is disabled")
}

func TestLiquidityStatusUseCase_Run_Success(t *testing.T) {
	providerMock := &mocks.ProviderMock{}

	providerMock.On("GeneralConfiguration", test.AnyCtx).Return(lp.GeneralConfiguration{PublicLiquidityCheck: true})
	providerMock.On("CalculateAvailablePeginLiquidity", test.AnyCtx).Return(entities.NewWei(10000), nil)
	providerMock.On("CalculateAvailablePegoutLiquidity", test.AnyCtx).Return(entities.NewWei(10000), nil)
	lbc := new(mocks.LbcMock)
	feeCollector := new(mocks.FeeCollectorMock)
	bridge := new(mocks.BridgeMock)
	useCase := liquidity_provider.NewLiquidityStatusUseCase(blockchain.RskContracts{Lbc: lbc, FeeCollector: feeCollector, Bridge: bridge}, providerMock, blockchain.Rpc{Btc: new(mocks.BtcRpcMock), Rsk: new(mocks.RskRpcMock)}, new(mocks.BtcWalletMock), providerMock, providerMock)
	status, err := useCase.Run(context.Background())
	require.NoError(t, err)
	require.NotNil(t, status)
	require.Equal(t, entities.NewWei(10000), status.Available.Pegin)
	require.Equal(t, entities.NewWei(10000), status.Available.Pegout)
}
