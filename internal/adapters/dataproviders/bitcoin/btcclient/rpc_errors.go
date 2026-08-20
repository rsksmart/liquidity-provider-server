package btcclient

import (
	"errors"
	"fmt"
	"strings"

	"github.com/btcsuite/btcd/btcjson"
)

// RPCCallError annotates a bitcoind RPC failure with the method used (we can extend it later to include some context of the caller too)
type RPCCallError struct {
	Method BitcoindRPCMethod
	RPC    *btcjson.RPCError
	Cause  error
}

func (e *RPCCallError) Error() string {
	if e.RPC != nil {
		return fmt.Sprintf("%s: RPC %d: %s", e.Method, e.RPC.Code, e.RPC.Message)
	}
	return fmt.Sprintf("%s: %v", e.Method, e.Cause)
}

func (e *RPCCallError) Unwrap() error { return e.Cause }

func AsRPCError(err error) (*btcjson.RPCError, bool) {
	var rpcErr *btcjson.RPCError
	if errors.As(err, &rpcErr) {
		return rpcErr, true
	}
	return nil, false
}

func IsRPCCode(err error, code btcjson.RPCErrorCode) bool {
	rpcErr, ok := AsRPCError(err)
	return ok && rpcErr.Code == code
}

func WrapRPCError(method BitcoindRPCMethod, err error) error {
	if err == nil {
		return nil
	}
	wrapped := &RPCCallError{Method: method, Cause: err}
	if rpcErr, ok := AsRPCError(err); ok {
		wrapped.RPC = rpcErr
	}
	return wrapped
}

func IsLockUnspentAlreadyUnlocked(err error) bool {
	if err == nil {
		return false
	}
	return IsRPCCode(err, btcjson.ErrRPCInvalidParameter) &&
		strings.Contains(strings.ToLower(err.Error()), "expected unspent output")
}

func IsAddressAlreadyImported(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "already imported")
}
