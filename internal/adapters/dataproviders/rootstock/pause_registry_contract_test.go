package rootstock_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	emergencyPauseBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/emergency_pause"
	pauseRegistryBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pause_registry"
	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/rsksmart/liquidity-provider-server/test/rootstocktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	pauseRegistrySelector = "1c82732d"
	pauseLevelSelector    = "e61399b8"
)

var testPauseRegistryAddress = common.HexToAddress("0x5555555555555555555555555555555555555555")

type pauseRegistryHarness struct {
	registry     blockchain.PauseRegistryContract
	contractMock rootstocktest.BoundContractMock
	client       *mocks.RpcClientBindingMock
}

func newPauseRegistryHarness() pauseRegistryHarness {
	contractMock := rootstocktest.NewBoundContractMock()
	client := &mocks.RpcClientBindingMock{}
	registry := rootstock.NewPauseRegistryContractImpl(
		rootstock.NewRskClient(client),
		contractMock.Contract,
		rootstock.RetryParams{},
		emergencyPauseBindings.NewEmergencyPauseContract(),
		pauseRegistryBindings.NewPauseRegistryContract(),
	)
	return pauseRegistryHarness{registry: registry, contractMock: contractMock, client: client}
}

func matchPauseLevelCall(t *testing.T) any {
	expectedData := rootstocktest.MustDecodeHex(t, pauseLevelSelector)
	return mock.MatchedBy(func(msg ethereum.CallMsg) bool {
		return msg.To != nil && *msg.To == testPauseRegistryAddress && bytes.Equal(msg.Data, expectedData)
	})
}

func TestNewPauseRegistryContractImpl(t *testing.T) {
	registry := rootstock.NewPauseRegistryContractImpl(
		rootstock.NewRskClient(&mocks.RpcClientBindingMock{}),
		rootstocktest.NewBoundContractMock().Contract,
		rootstock.RetryParams{Retries: 1, Sleep: 1},
		emergencyPauseBindings.NewEmergencyPauseContract(),
		pauseRegistryBindings.NewPauseRegistryContract(),
	)
	test.AssertNonZeroValues(t, registry)
}

func TestPauseRegistryContractImpl_PauseLevel(t *testing.T) {
	levels := []struct {
		name     string
		raw      int64
		expected blockchain.PauseLevel
	}{
		{name: "none", raw: 0, expected: blockchain.PauseLevelNone},
		{name: "soft", raw: 1, expected: blockchain.PauseLevelSoft},
		{name: "hard", raw: 2, expected: blockchain.PauseLevelHard},
	}
	for _, testCase := range levels {
		t.Run("reads the "+testCase.name+" level from the registry the peg-in contract uses", func(t *testing.T) {
			h := newPauseRegistryHarness()
			h.contractMock.Caller.EXPECT().CallContract(mock.Anything, rootstocktest.MatchCallData(rootstocktest.MustDecodeHex(t, pauseRegistrySelector)), mock.Anything).
				Return(common.LeftPadBytes(testPauseRegistryAddress.Bytes(), common.HashLength), nil).Once()
			h.client.EXPECT().CallContract(mock.Anything, matchPauseLevelCall(t), mock.Anything).
				Return(rootstocktest.MustPackUint256(t, big.NewInt(testCase.raw)), nil).Once()

			level, err := h.registry.PauseLevel()
			require.NoError(t, err)
			assert.Equal(t, testCase.expected, level)
			h.contractMock.Caller.AssertExpectations(t)
			h.client.AssertExpectations(t)
		})
	}

	t.Run("returns the error reading the pause registry address", func(t *testing.T) {
		h := newPauseRegistryHarness()
		h.contractMock.Caller.EXPECT().CallContract(mock.Anything, rootstocktest.MatchCallData(rootstocktest.MustDecodeHex(t, pauseRegistrySelector)), mock.Anything).
			Return(nil, assert.AnError).Once()

		level, err := h.registry.PauseLevel()
		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, blockchain.PauseLevelNone, level)
		h.client.AssertNotCalled(t, "CallContract", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("returns the error reading the pause level", func(t *testing.T) {
		h := newPauseRegistryHarness()
		h.contractMock.Caller.EXPECT().CallContract(mock.Anything, rootstocktest.MatchCallData(rootstocktest.MustDecodeHex(t, pauseRegistrySelector)), mock.Anything).
			Return(common.LeftPadBytes(testPauseRegistryAddress.Bytes(), common.HashLength), nil).Once()
		h.client.EXPECT().CallContract(mock.Anything, matchPauseLevelCall(t), mock.Anything).Return(nil, assert.AnError).Once()

		level, err := h.registry.PauseLevel()
		require.ErrorIs(t, err, assert.AnError)
		assert.Equal(t, blockchain.PauseLevelNone, level)
		h.client.AssertExpectations(t)
	})
}
