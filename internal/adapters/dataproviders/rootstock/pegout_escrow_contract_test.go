package rootstock_test

import (
	"bytes"
	"math"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	escrowBindings "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock/bindings/pegout_escrow"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPegOutEscrowContractImpl_RestrictedUntil(t *testing.T) {
	cases := []struct {
		name     string
		returned []byte
		expected uint64
	}{
		{name: "not restricted", returned: make([]byte, 32), expected: 0},
		{name: "timed restriction", returned: common.LeftPadBytes(big.NewInt(1_700_000_000).Bytes(), 32), expected: 1_700_000_000},
		{name: "admin revoke", returned: bytes.Repeat([]byte{0xff}, 32), expected: math.MaxUint64},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			contractMock := createBoundContractMock()
			binding := escrowBindings.NewPegOutEscrowContract()
			escrow := rootstock.NewPegOutEscrowContractImpl(dummyClient, test.AnyAddress, contractMock.contract, nil, rootstock.RetryParams{}, time.Duration(1), binding, Abis)
			contractMock.caller.EXPECT().CallContract(
				mock.Anything,
				matchCallData(binding.PackRestrictedUntil(common.HexToAddress(test.AnyRskAddress))),
				mock.Anything,
			).Return(c.returned, nil).Once()

			result, err := escrow.RestrictedUntil(test.AnyRskAddress)
			require.NoError(t, err)
			assert.Equal(t, c.expected, result)
			contractMock.caller.AssertExpectations(t)
		})
	}
}
