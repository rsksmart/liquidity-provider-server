package btcclient_test

import (
	"errors"
	"testing"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/bitcoin/btcclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAsRPCError(t *testing.T) {
	rpcErr := btcjson.NewRPCError(btcjson.ErrRPCInvalidParameter, "Invalid parameter, expected unspent output")
	extracted, ok := btcclient.AsRPCError(rpcErr)
	require.True(t, ok)
	require.NotNil(t, extracted)
	assert.Equal(t, rpcErr, extracted)

	err, ok := btcclient.AsRPCError(errors.New("plain error"))
	assert.False(t, ok)
	require.Nil(t, err)
}

func TestIsRPCCode(t *testing.T) {
	err := btcjson.NewRPCError(btcjson.ErrRPCInvalidParameter, "Invalid parameter")
	assert.True(t, btcclient.IsRPCCode(err, btcjson.ErrRPCInvalidParameter))
	assert.False(t, btcclient.IsRPCCode(err, btcjson.ErrRPCMethodNotFound.Code))
}

func TestWrapRPCError(t *testing.T) {
	cause := btcjson.NewRPCError(btcjson.ErrRPCInvalidParameter, "Invalid parameter, expected unspent output")
	wrapped := btcclient.WrapRPCError(btcclient.MethodLockUnspent, cause)
	require.Error(t, wrapped)
	require.ErrorContains(t, wrapped, string(btcclient.MethodLockUnspent))
	require.ErrorContains(t, wrapped, "expected unspent output")

	var callErr *btcclient.RPCCallError
	require.ErrorAs(t, wrapped, &callErr)
	assert.Equal(t, btcclient.MethodLockUnspent, callErr.Method)
	assert.Equal(t, cause, callErr.Cause)
}

func TestIsLockUnspentAlreadyUnlocked(t *testing.T) {
	assert.False(t, btcclient.IsLockUnspentAlreadyUnlocked(nil))

	rpcErr := btcjson.NewRPCError(btcjson.ErrRPCInvalidParameter, "Invalid parameter, expected unspent output")
	assert.True(t, btcclient.IsLockUnspentAlreadyUnlocked(rpcErr))

	otherErr := btcjson.NewRPCError(btcjson.ErrRPCMethodNotFound.Code, "Method not found")
	assert.False(t, btcclient.IsLockUnspentAlreadyUnlocked(otherErr))
}
