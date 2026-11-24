package pegin_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	commitFirstBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/rootstocktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func TestPeginContractImpl_GetBalance(t *testing.T) {
	contractMock := rootstocktest.NewBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	peginContract := rootstock.NewPeginContractImpl(rootstocktest.DummyClient, test.AnyAddress, contractMock.Contract, nil, rootstock.RetryParams{}, time.Duration(1), peginBinding, commitFirstBindings.NewPeginCommitFirstContract(), rootstocktest.Abis)
	t.Run("Success", func(t *testing.T) {
		contractMock.Caller.EXPECT().CallContract(
			mock.Anything,
			rootstocktest.MatchCallData(peginBinding.PackGetBalance(rootstocktest.ParsedAddress)),
			mock.Anything,
		).Return(rootstocktest.MustPackUint256(t, big.NewInt(600)), nil).Once()
		result, err := peginContract.GetBalance(rootstocktest.ParsedAddress.String())
		require.NoError(t, err)
		assert.Equal(t, entities.NewWei(600), result)
		contractMock.Caller.AssertExpectations(t)
	})
	t.Run("Error handling on GetBalance call error", func(t *testing.T) {
		contractMock.Caller.EXPECT().CallContract(
			mock.Anything,
			rootstocktest.MatchCallData(peginBinding.PackGetBalance(rootstocktest.ParsedAddress)),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := peginContract.GetBalance(rootstocktest.ParsedAddress.String())
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("Error handling on invalid address for getting balance", func(t *testing.T) {
		result, err := peginContract.GetBalance(test.AnyString)
		require.Error(t, err)
		assert.Nil(t, result)
	})
}
