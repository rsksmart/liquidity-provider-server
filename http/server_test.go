package http

import (
	"github.com/rsksmart/liquidity-provider/providers"
	"github.com/rsksmart/liquidity-provider/types"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

type LiquidityProviderMock struct {
	address string
}
func (lp LiquidityProviderMock) Address() string {
	return lp.address
}
func (lp LiquidityProviderMock) GetQuote(q types.Quote, gas uint64, gasPrice big.Int) *types.Quote {
	return nil
}
func (lp LiquidityProviderMock) SignHash(hash []byte) ([]byte, error) {
	return nil, nil
}

var providerMocks = []LiquidityProviderMock {}


func testGetProviderByAddress(t *testing.T) {
	var liquidityProviders []providers.LiquidityProvider
	for _, mock := range providerMocks {
		liquidityProviders = append(liquidityProviders, mock)
	}

	for _, tt := range liquidityProviders {
		result := getProviderByAddress(liquidityProviders, tt.Address())
		assert.EqualValues(t, tt.Address(), result.Address())
	}
}

func setup() {
	providerMocks = []LiquidityProviderMock{
		{ address: "123" },
		{ address: "12345" },
	}
}
func TestLiquidityProviderServer(t *testing.T) {
	setup()
	t.Run("new", TestGetProviderByAddress)
}

