package blockchain_test

import (
	"reflect"
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/rsksmart/liquidity-provider-server/test/mocks"
	"github.com/stretchr/testify/assert"
)

// PegInAddressRegistryContract must be a read-only port: registerAddress belongs to a
// separate on-chain watcher process, not the liquidity provider server, so it must never
// appear on this interface.
// nolint:funlen
func TestPegInAddressRegistryContract_MethodSet(t *testing.T) {
	contractType := reflect.TypeFor[blockchain.PegInAddressRegistryContract]()
	expectedMethods := []string{
		"GetAddress",
		"GetPegInAddress",
		"GetPegInAddresses",
		"IsRegistered",
		"GetRegistration",
		"GetRegistrationRoot",
		"GetAddressRegisteredEvents",
	}

	assert.Equal(t, len(expectedMethods), contractType.NumMethod(), "PegInAddressRegistryContract must expose exactly its intended read-only surface")

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

// Guards against the mock and the interface it mocks drifting apart after either one is
// regenerated.
func TestPegInAddressRegistryContractMock_SatisfiesInterface(t *testing.T) {
	var _ blockchain.PegInAddressRegistryContract = &mocks.PegInAddressRegistryContractMock{}
}
