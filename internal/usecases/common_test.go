package usecases_test

import (
	"context"
	"errors"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	u "github.com/rsksmart/liquidity-provider-server/internal/usecases"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"slices"
	"testing"
)

type rpcMock struct {
	mock.Mock
	blockchain.RootstockRpcServer
}

func (m *rpcMock) EstimateGas(ctx context.Context, addr string, value *entities.Wei, data []byte) (*entities.Wei, error) {
	args := m.Called(ctx, addr, value, data)
	return args.Get(0).(*entities.Wei), args.Error(1)
}

type bridgeMock struct {
	mock.Mock
	blockchain.RootstockBridge
}

func (m *bridgeMock) GetMinimumLockTxValue() (*entities.Wei, error) {
	args := m.Called()
	return args.Get(0).(*entities.Wei), args.Error(1)
}

func TestCalculateDaoAmounts(t *testing.T) {
	type testArgs struct {
		value      *entities.Wei
		percentage uint64
	}
	rpc := rpcMock{}
	rpc.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(entities.NewWei(500000000000000), nil)

	ctx := context.Background()
	feeCollectorAddress := "0x1234"

	cases := test.Table[testArgs, u.DaoAmounts]{
		{
			Value:  testArgs{entities.NewWei(1000000000000000000), 0},
			Result: u.DaoAmounts{DaoFeeAmount: entities.NewWei(0), DaoGasAmount: entities.NewWei(0)},
		},
		{
			Value: testArgs{entities.NewWei(500000000000000000), 50},
			Result: u.DaoAmounts{
				DaoGasAmount: entities.NewWei(500000000000000),
				DaoFeeAmount: entities.NewWei(250000000000000000),
			},
		},
		{
			Value: testArgs{entities.NewWei(6000000000000000000), 1},
			Result: u.DaoAmounts{
				DaoGasAmount: entities.NewWei(500000000000000),
				DaoFeeAmount: entities.NewWei(60000000000000000),
			},
		},
		{
			Value: testArgs{entities.NewWei(7700000000000000000), 17},
			Result: u.DaoAmounts{
				DaoGasAmount: entities.NewWei(500000000000000),
				DaoFeeAmount: entities.NewWei(1309000000000000000),
			},
		},
	}

	test.RunTable(t, cases, func(args testArgs) u.DaoAmounts {
		amounts, _ := u.CalculateDaoAmounts(ctx, &rpc, args.value, args.percentage, feeCollectorAddress)
		return amounts
	})

}

func TestCalculateDaoAmounts_Fail(t *testing.T) {
	ctx := context.Background()
	rpc := rpcMock{}
	rpc.On("EstimateGas", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(entities.NewWei(0), errors.New("some error"))
	result, err := u.CalculateDaoAmounts(ctx, &rpc, entities.NewUWei(500000000000000), 1, "0x1234")
	require.Equal(t, u.DaoAmounts{}, result)
	require.Error(t, err)
}

func TestGetRandomInt(t *testing.T) {
	var numbers []int64
	var number int64
	var err error
	for i := 0; i < 100; i++ {
		number, err = u.GetRandomInt()
		assert.Positive(t, number)
		assert.False(t, slices.Contains(numbers, number))
		require.NoError(t, err)
		numbers = append(numbers, number)
	}
}

func TestValidateMinLockValue(t *testing.T) {
	var oneBtcInSatoshi uint64 = 1 * bitcoin.BtcToSatoshi
	var useCase u.UseCaseId = "anyUseCase"
	bridge := &bridgeMock{}
	bridge.On("GetMinimumLockTxValue").Return(entities.SatoshiToWei(oneBtcInSatoshi), nil)

	err := u.ValidateMinLockValue(useCase, bridge, entities.SatoshiToWei(oneBtcInSatoshi))
	require.NoError(t, err)

	err = u.ValidateMinLockValue(useCase, bridge, entities.SatoshiToWei(oneBtcInSatoshi+1))
	require.NoError(t, err)

	value := new(entities.Wei).Sub(entities.SatoshiToWei(oneBtcInSatoshi), entities.NewWei(1))
	err = u.ValidateMinLockValue(useCase, bridge, value)
	require.Error(t, err)
	assert.Equal(t, "anyUseCase: requested amount below bridge's min transaction value. Args: {\"minimum\":\"1000000000000000000\",\"value\":\"999999999999999999\"}", err.Error())
}
