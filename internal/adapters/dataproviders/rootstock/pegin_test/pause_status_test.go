package pegin_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	bindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin"
	commitFirstBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegin_commit_first"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/rootstocktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPeginContractImpl_PausedStatus(t *testing.T) {
	contractMock := rootstocktest.NewBoundContractMock()
	peginBinding := bindings.NewPeginContract()
	contract := rootstock.NewPeginContractImpl(rootstocktest.DummyClient, test.AnyAddress, contractMock.Contract, nil, rootstock.RetryParams{}, time.Duration(1), peginBinding, commitFirstBindings.NewPeginCommitFirstContract(), rootstocktest.Abis)
	t.Run("should return pause status result", func(t *testing.T) {
		contractMock.Caller.EXPECT().CallContract(
			mock.Anything,
			rootstocktest.MatchCallData(peginBinding.PackPauseStatus()),
			mock.Anything,
		).Return(rootstocktest.MustPackPauseStatus(t, rootstocktest.GeneralPauseStatus{IsPaused: true, Reason: "test", Since: 123}), nil).Once()
		result, err := contract.PausedStatus()
		require.NoError(t, err)
		assert.Equal(t, blockchain.PauseStatus{IsPaused: true, Reason: "test", Since: 123}, result)
	})
	t.Run("should handle error checking pause status", func(t *testing.T) {
		contractMock.Caller.EXPECT().CallContract(
			mock.Anything,
			rootstocktest.MatchCallData(peginBinding.PackPauseStatus()),
			mock.Anything,
		).Return(nil, assert.AnError).Once()
		result, err := contract.PausedStatus()
		require.Error(t, err)
		assert.Empty(t, result)
	})
	contractMock.Caller.AssertExpectations(t)
}
