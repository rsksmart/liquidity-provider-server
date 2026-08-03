package blockchain_test

import (
	"reflect"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
)

// FlyoverConfigurationsContract must expose exactly its two read methods.
// getPegInConfiguration and the time-locked admin writes (queueChange/applyChange) are out
// of scope for this read-only port and must never appear on this interface.
func TestFlyoverConfigurationsContract_MethodSet(t *testing.T) {
	contractType := reflect.TypeFor[blockchain.FlyoverConfigurationsContract]()
	expectedMethods := []string{
		"GetAddress",
		"CalculatePegInFee",
		"GetRequiredPegInBtcConfirmations",
	}

	assert.Equal(t, len(expectedMethods), contractType.NumMethod(), "FlyoverConfigurationsContract must expose exactly its intended read-only surface")

	actualMethods := make([]string, contractType.NumMethod())
	for i := 0; i < contractType.NumMethod(); i++ {
		actualMethods[i] = contractType.Method(i).Name
	}
	assert.ElementsMatch(t, expectedMethods, actualMethods)

	disallowedWriteMethods := []string{"QueueChange", "ApplyChange", "GetPegInConfiguration"}
	for _, disallowed := range disallowedWriteMethods {
		_, found := contractType.MethodByName(disallowed)
		assert.False(t, found, "FlyoverConfigurationsContract must not expose method %q", disallowed)
	}
}

// Guards against the mock and the interface it mocks drifting apart after either one is
// regenerated.
func TestFlyoverConfigurationsContractMock_SatisfiesInterface(t *testing.T) {
	var _ blockchain.FlyoverConfigurationsContract = &mocks.FlyoverConfigurationsContractMock{}
}
