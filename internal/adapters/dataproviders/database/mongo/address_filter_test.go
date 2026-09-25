package mongo_test

import (
	"testing"

	"github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	trustedAccountAddressField = "address"
	ownerAccountAddressField   = "owner_account_address"
)

func matchTrustedAccountAddressFilter(t *testing.T, address string) any {
	t.Helper()
	return mock.MatchedBy(func(filter bson.M) bool {
		return matchRskAddressEquals(t, filter, trustedAccountAddressField, address)
	})
}

func matchRetainedQuotesOwnerFilter(t *testing.T, address string, states any) any {
	t.Helper()
	return mock.MatchedBy(func(filter bson.M) bool {
		if !matchRskAddressEquals(t, filter, ownerAccountAddressField, address) {
			return false
		}
		stateFilter, ok := filter["state"].(bson.D)
		if !assert.True(t, ok) || !assert.Len(t, stateFilter, 1) {
			return false
		}
		return assert.Equal(t, "$in", stateFilter[0].Key) &&
			assert.Equal(t, states, stateFilter[0].Value)
	})
}

func matchRskAddressEquals(t *testing.T, filter bson.M, field, address string) bool {
	t.Helper()
	normalized, err := blockchain.NormalizeRskAddress(address)
	require.NoError(t, err)
	expr, ok := filter["$expr"].(bson.M)
	if !assert.True(t, ok) {
		return false
	}
	eq, ok := expr["$eq"].([]any)
	if !assert.True(t, ok) || !assert.Len(t, eq, 2) {
		return false
	}
	toLower, ok := eq[0].(bson.M)
	if !assert.True(t, ok) {
		return false
	}
	return assert.Equal(t, "$"+field, toLower["$toLower"]) &&
		assert.Equal(t, normalized, eq[1])
}
