package liquidity_provider_test

import (
	"context"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	lp "github.com/rsksmart/liquidity-provider-server/internal/entities/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/internal/usecases/liquidity_provider"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupMocks(publicLiquidityCheck bool, peginLiquidity *entities.Wei, peginErr error, pegoutLiquidity *entities.Wei, pegoutErr error) *liquidity_provider.LiquidityStatusUseCase {
	providerMock := &mocks.ProviderMock{}
	rpcMocks := blockchain.Rpc{Btc: new(mocks.BtcRpcMock), Rsk: new(mocks.RskRpcMock)}
	providerMock.On("GeneralConfiguration", mock.Anything).Return(lp.GeneralConfiguration{PublicLiquidityCheck: publicLiquidityCheck})
	if peginLiquidity != nil || peginErr != nil {
		providerMock.On("CalculateAvailablePeginLiquidity", mock.Anything).Return(peginLiquidity, peginErr)
	}
	if pegoutLiquidity != nil || pegoutErr != nil {
		providerMock.On("CalculateAvailablePegoutLiquidity", mock.Anything).Return(pegoutLiquidity, pegoutErr)
	}
	return liquidity_provider.NewLiquidityStatusUseCase(blockchain.RskContracts{Lbc: new(mocks.LbcMock), FeeCollector: new(mocks.FeeCollectorMock), Bridge: new(mocks.BridgeMock)}, providerMock, rpcMocks, new(mocks.BtcWalletMock), providerMock, providerMock)
}

func TestLiquidityStatusUseCase_Run_PublicLiquidityCheckDisabled(t *testing.T) {
	useCase := setupMocks(false, nil, nil, nil, nil)
	status, err := useCase.Run(context.Background())
	require.Nil(t, status)
	require.ErrorIs(t, err, usecases.PublicLiquidityCheckDisabledError)
}

func TestLiquidityStatusUseCase_Run_ErrorFetchingPeginLiquidity(t *testing.T) {
	useCase := setupMocks(true, nil, usecases.PublicLiquidityPeginCheckError, nil, nil)
	status, err := useCase.Run(context.Background())
	require.Nil(t, status)
	require.ErrorIs(t, err, usecases.PublicLiquidityPeginCheckError)
}

func TestLiquidityStatusUseCase_Run_ErrorFetchingPegoutLiquidity(t *testing.T) {
	useCase := setupMocks(true, entities.NewWei(0), nil, nil, usecases.PublicLiquidityPegoutCheckError)
	status, err := useCase.Run(context.Background())
	require.Nil(t, status)
	require.ErrorIs(t, err, usecases.PublicLiquidityPegoutCheckError)
}

func TestLiquidityStatusUseCase_Run_Success(t *testing.T) {
	useCase := setupMocks(true, entities.NewWei(10000), nil, entities.NewWei(10000), nil)
	status, err := useCase.Run(context.Background())
	require.NoError(t, err)
	require.NotNil(t, status)
	require.Equal(t, entities.NewWei(10000), status.Available.Pegin)
	require.Equal(t, entities.NewWei(10000), status.Available.Pegout)
}
