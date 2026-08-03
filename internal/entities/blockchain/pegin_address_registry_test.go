package blockchain_test

import (
	"reflect"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
)

// S4 (fly-2513): PegInAddressRegistryContract must be a read-only port. registerAddress is the
// watchtower's job (FLY-2446), never the LPS's, so it must never appear on this interface.
// nolint:funlen
func TestPegInAddressRegistryContract_MethodSet(t *testing.T) {
	contractType := reflect.TypeOf((*blockchain.PegInAddressRegistryContract)(nil)).Elem()
	expectedMethods := []string{
		"GetAddress",
		"GetPegInAddress",
		"GetPegInAddresses",
		"IsRegistered",
		"GetRegistration",
		"GetRegistrationRoot",
		"GetAddressRegisteredEvents",
	}

	assert.Equal(t, len(expectedMethods), contractType.NumMethod(), "PegInAddressRegistryContract must expose exactly the read-only surface required by the ticket AC")

	actualMethods := make([]string, contractType.NumMethod())
	for i := 0; i < contractType.NumMethod(); i++ {
		actualMethods[i] = contractType.Method(i).Name
	}
	assert.ElementsMatch(t, expectedMethods, actualMethods)

	disallowedWriteMethods := []string{"RegisterAddress", "Register", "SetRegistration", "Write"}
	for _, disallowed := range disallowedWriteMethods {
		_, found := contractType.MethodByName(disallowed)
		assert.False(t, found, "PegInAddressRegistryContract must not expose write method %q", disallowed)
	}
}

// S5 (fly-2513): the regenerated mock must still compile against, and satisfy, the port
// interface it mocks.
func TestPegInAddressRegistryContractMock_SatisfiesInterface(t *testing.T) {
	var _ blockchain.PegInAddressRegistryContract = &mocks.PegInAddressRegistryContractMock{}
}
