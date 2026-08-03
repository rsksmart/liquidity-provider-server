package blockchain_test

import (
	"reflect"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
)

// S4 (fly-2513): FlyoverConfigurationsContract exposes exactly the two AC-listed reads.
// getPegInConfiguration and the time-locked admin writes (queueChange/applyChange) are out
// of scope for this ticket and must never appear on this interface.
func TestFlyoverConfigurationsContract_MethodSet(t *testing.T) {
	contractType := reflect.TypeOf((*blockchain.FlyoverConfigurationsContract)(nil)).Elem()
	expectedMethods := []string{
		"GetAddress",
		"CalculatePegInFee",
		"GetRequiredPegInBtcConfirmations",
	}

	assert.Equal(t, len(expectedMethods), contractType.NumMethod(), "FlyoverConfigurationsContract must expose exactly the read-only surface required by the ticket AC")

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

// S5 (fly-2513): the regenerated mock must still compile against, and satisfy, the port
// interface it mocks.
func TestFlyoverConfigurationsContractMock_SatisfiesInterface(t *testing.T) {
	var _ blockchain.FlyoverConfigurationsContract = &mocks.FlyoverConfigurationsContractMock{}
}
