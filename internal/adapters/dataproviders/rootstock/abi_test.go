package rootstock_test

import (
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
	"github.com/rsksmart/liquidity-provider-server/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMustLoadFlyoverABIs(t *testing.T) {
	assert.NotPanics(t, func() {
		result := rootstock.MustLoadFlyoverABIs()
		test.AssertNonZeroValues(t, result)
	})
}
