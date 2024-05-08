// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"
	big "math/big"

	common "github.com/ethereum/go-ethereum/common"

	ethereum "github.com/ethereum/go-ethereum"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// RpcClientBindingMock is an autogenerated mock type for the RpcClientBinding type
type RpcClientBindingMock struct {
	mock.Mock
}

type RpcClientBindingMock_Expecter struct {
	mock *mock.Mock
}

func (_m *RpcClientBindingMock) EXPECT() *RpcClientBindingMock_Expecter {
	return &RpcClientBindingMock_Expecter{mock: &_m.Mock}
}

// BalanceAt provides a mock function with given fields: ctx, account, blockNumber
func (_m *RpcClientBindingMock) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	ret := _m.Called(ctx, account, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for BalanceAt")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) (*big.Int, error)); ok {
		return rf(ctx, account, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) *big.Int); ok {
		r0 = rf(ctx, account, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int) error); ok {
		r1 = rf(ctx, account, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_BalanceAt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BalanceAt'
type RpcClientBindingMock_BalanceAt_Call struct {
	*mock.Call
}

// BalanceAt is a helper method to define mock.On call
//   - ctx context.Context
//   - account common.Address
//   - blockNumber *big.Int
func (_e *RpcClientBindingMock_Expecter) BalanceAt(ctx interface{}, account interface{}, blockNumber interface{}) *RpcClientBindingMock_BalanceAt_Call {
	return &RpcClientBindingMock_BalanceAt_Call{Call: _e.mock.On("BalanceAt", ctx, account, blockNumber)}
}

func (_c *RpcClientBindingMock_BalanceAt_Call) Run(run func(ctx context.Context, account common.Address, blockNumber *big.Int)) *RpcClientBindingMock_BalanceAt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Address), args[2].(*big.Int))
	})
	return _c
}

func (_c *RpcClientBindingMock_BalanceAt_Call) Return(_a0 *big.Int, _a1 error) *RpcClientBindingMock_BalanceAt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_BalanceAt_Call) RunAndReturn(run func(context.Context, common.Address, *big.Int) (*big.Int, error)) *RpcClientBindingMock_BalanceAt_Call {
	_c.Call.Return(run)
	return _c
}

// BlockByHash provides a mock function with given fields: ctx, hash
func (_m *RpcClientBindingMock) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for BlockByHash")
	}

	var r0 *types.Block
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*types.Block, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *types.Block); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Block)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_BlockByHash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BlockByHash'
type RpcClientBindingMock_BlockByHash_Call struct {
	*mock.Call
}

// BlockByHash is a helper method to define mock.On call
//   - ctx context.Context
//   - hash common.Hash
func (_e *RpcClientBindingMock_Expecter) BlockByHash(ctx interface{}, hash interface{}) *RpcClientBindingMock_BlockByHash_Call {
	return &RpcClientBindingMock_BlockByHash_Call{Call: _e.mock.On("BlockByHash", ctx, hash)}
}

func (_c *RpcClientBindingMock_BlockByHash_Call) Run(run func(ctx context.Context, hash common.Hash)) *RpcClientBindingMock_BlockByHash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Hash))
	})
	return _c
}

func (_c *RpcClientBindingMock_BlockByHash_Call) Return(_a0 *types.Block, _a1 error) *RpcClientBindingMock_BlockByHash_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_BlockByHash_Call) RunAndReturn(run func(context.Context, common.Hash) (*types.Block, error)) *RpcClientBindingMock_BlockByHash_Call {
	_c.Call.Return(run)
	return _c
}

// BlockNumber provides a mock function with given fields: ctx
func (_m *RpcClientBindingMock) BlockNumber(ctx context.Context) (uint64, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for BlockNumber")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (uint64, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) uint64); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_BlockNumber_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BlockNumber'
type RpcClientBindingMock_BlockNumber_Call struct {
	*mock.Call
}

// BlockNumber is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RpcClientBindingMock_Expecter) BlockNumber(ctx interface{}) *RpcClientBindingMock_BlockNumber_Call {
	return &RpcClientBindingMock_BlockNumber_Call{Call: _e.mock.On("BlockNumber", ctx)}
}

func (_c *RpcClientBindingMock_BlockNumber_Call) Run(run func(ctx context.Context)) *RpcClientBindingMock_BlockNumber_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RpcClientBindingMock_BlockNumber_Call) Return(_a0 uint64, _a1 error) *RpcClientBindingMock_BlockNumber_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_BlockNumber_Call) RunAndReturn(run func(context.Context) (uint64, error)) *RpcClientBindingMock_BlockNumber_Call {
	_c.Call.Return(run)
	return _c
}

// ChainID provides a mock function with given fields: ctx
func (_m *RpcClientBindingMock) ChainID(ctx context.Context) (*big.Int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for ChainID")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*big.Int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *big.Int); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_ChainID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ChainID'
type RpcClientBindingMock_ChainID_Call struct {
	*mock.Call
}

// ChainID is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RpcClientBindingMock_Expecter) ChainID(ctx interface{}) *RpcClientBindingMock_ChainID_Call {
	return &RpcClientBindingMock_ChainID_Call{Call: _e.mock.On("ChainID", ctx)}
}

func (_c *RpcClientBindingMock_ChainID_Call) Run(run func(ctx context.Context)) *RpcClientBindingMock_ChainID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RpcClientBindingMock_ChainID_Call) Return(_a0 *big.Int, _a1 error) *RpcClientBindingMock_ChainID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_ChainID_Call) RunAndReturn(run func(context.Context) (*big.Int, error)) *RpcClientBindingMock_ChainID_Call {
	_c.Call.Return(run)
	return _c
}

// Close provides a mock function with given fields:
func (_m *RpcClientBindingMock) Close() {
	_m.Called()
}

// RpcClientBindingMock_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type RpcClientBindingMock_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *RpcClientBindingMock_Expecter) Close() *RpcClientBindingMock_Close_Call {
	return &RpcClientBindingMock_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *RpcClientBindingMock_Close_Call) Run(run func()) *RpcClientBindingMock_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *RpcClientBindingMock_Close_Call) Return() *RpcClientBindingMock_Close_Call {
	_c.Call.Return()
	return _c
}

func (_c *RpcClientBindingMock_Close_Call) RunAndReturn(run func()) *RpcClientBindingMock_Close_Call {
	_c.Call.Return(run)
	return _c
}

// CodeAt provides a mock function with given fields: ctx, account, blockNumber
func (_m *RpcClientBindingMock) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	ret := _m.Called(ctx, account, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for CodeAt")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) ([]byte, error)); ok {
		return rf(ctx, account, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) []byte); ok {
		r0 = rf(ctx, account, blockNumber)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int) error); ok {
		r1 = rf(ctx, account, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_CodeAt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CodeAt'
type RpcClientBindingMock_CodeAt_Call struct {
	*mock.Call
}

// CodeAt is a helper method to define mock.On call
//   - ctx context.Context
//   - account common.Address
//   - blockNumber *big.Int
func (_e *RpcClientBindingMock_Expecter) CodeAt(ctx interface{}, account interface{}, blockNumber interface{}) *RpcClientBindingMock_CodeAt_Call {
	return &RpcClientBindingMock_CodeAt_Call{Call: _e.mock.On("CodeAt", ctx, account, blockNumber)}
}

func (_c *RpcClientBindingMock_CodeAt_Call) Run(run func(ctx context.Context, account common.Address, blockNumber *big.Int)) *RpcClientBindingMock_CodeAt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Address), args[2].(*big.Int))
	})
	return _c
}

func (_c *RpcClientBindingMock_CodeAt_Call) Return(_a0 []byte, _a1 error) *RpcClientBindingMock_CodeAt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_CodeAt_Call) RunAndReturn(run func(context.Context, common.Address, *big.Int) ([]byte, error)) *RpcClientBindingMock_CodeAt_Call {
	_c.Call.Return(run)
	return _c
}

// EstimateGas provides a mock function with given fields: ctx, msg
func (_m *RpcClientBindingMock) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	ret := _m.Called(ctx, msg)

	if len(ret) == 0 {
		panic("no return value specified for EstimateGas")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ethereum.CallMsg) (uint64, error)); ok {
		return rf(ctx, msg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ethereum.CallMsg) uint64); ok {
		r0 = rf(ctx, msg)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, ethereum.CallMsg) error); ok {
		r1 = rf(ctx, msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_EstimateGas_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EstimateGas'
type RpcClientBindingMock_EstimateGas_Call struct {
	*mock.Call
}

// EstimateGas is a helper method to define mock.On call
//   - ctx context.Context
//   - msg ethereum.CallMsg
func (_e *RpcClientBindingMock_Expecter) EstimateGas(ctx interface{}, msg interface{}) *RpcClientBindingMock_EstimateGas_Call {
	return &RpcClientBindingMock_EstimateGas_Call{Call: _e.mock.On("EstimateGas", ctx, msg)}
}

func (_c *RpcClientBindingMock_EstimateGas_Call) Run(run func(ctx context.Context, msg ethereum.CallMsg)) *RpcClientBindingMock_EstimateGas_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(ethereum.CallMsg))
	})
	return _c
}

func (_c *RpcClientBindingMock_EstimateGas_Call) Return(_a0 uint64, _a1 error) *RpcClientBindingMock_EstimateGas_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_EstimateGas_Call) RunAndReturn(run func(context.Context, ethereum.CallMsg) (uint64, error)) *RpcClientBindingMock_EstimateGas_Call {
	_c.Call.Return(run)
	return _c
}

// NonceAt provides a mock function with given fields: ctx, account, blockNumber
func (_m *RpcClientBindingMock) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	ret := _m.Called(ctx, account, blockNumber)

	if len(ret) == 0 {
		panic("no return value specified for NonceAt")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) (uint64, error)); ok {
		return rf(ctx, account, blockNumber)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address, *big.Int) uint64); ok {
		r0 = rf(ctx, account, blockNumber)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address, *big.Int) error); ok {
		r1 = rf(ctx, account, blockNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_NonceAt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'NonceAt'
type RpcClientBindingMock_NonceAt_Call struct {
	*mock.Call
}

// NonceAt is a helper method to define mock.On call
//   - ctx context.Context
//   - account common.Address
//   - blockNumber *big.Int
func (_e *RpcClientBindingMock_Expecter) NonceAt(ctx interface{}, account interface{}, blockNumber interface{}) *RpcClientBindingMock_NonceAt_Call {
	return &RpcClientBindingMock_NonceAt_Call{Call: _e.mock.On("NonceAt", ctx, account, blockNumber)}
}

func (_c *RpcClientBindingMock_NonceAt_Call) Run(run func(ctx context.Context, account common.Address, blockNumber *big.Int)) *RpcClientBindingMock_NonceAt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Address), args[2].(*big.Int))
	})
	return _c
}

func (_c *RpcClientBindingMock_NonceAt_Call) Return(_a0 uint64, _a1 error) *RpcClientBindingMock_NonceAt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_NonceAt_Call) RunAndReturn(run func(context.Context, common.Address, *big.Int) (uint64, error)) *RpcClientBindingMock_NonceAt_Call {
	_c.Call.Return(run)
	return _c
}

// PendingNonceAt provides a mock function with given fields: ctx, account
func (_m *RpcClientBindingMock) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	ret := _m.Called(ctx, account)

	if len(ret) == 0 {
		panic("no return value specified for PendingNonceAt")
	}

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) (uint64, error)); ok {
		return rf(ctx, account)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) uint64); ok {
		r0 = rf(ctx, account)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address) error); ok {
		r1 = rf(ctx, account)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_PendingNonceAt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PendingNonceAt'
type RpcClientBindingMock_PendingNonceAt_Call struct {
	*mock.Call
}

// PendingNonceAt is a helper method to define mock.On call
//   - ctx context.Context
//   - account common.Address
func (_e *RpcClientBindingMock_Expecter) PendingNonceAt(ctx interface{}, account interface{}) *RpcClientBindingMock_PendingNonceAt_Call {
	return &RpcClientBindingMock_PendingNonceAt_Call{Call: _e.mock.On("PendingNonceAt", ctx, account)}
}

func (_c *RpcClientBindingMock_PendingNonceAt_Call) Run(run func(ctx context.Context, account common.Address)) *RpcClientBindingMock_PendingNonceAt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Address))
	})
	return _c
}

func (_c *RpcClientBindingMock_PendingNonceAt_Call) Return(_a0 uint64, _a1 error) *RpcClientBindingMock_PendingNonceAt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_PendingNonceAt_Call) RunAndReturn(run func(context.Context, common.Address) (uint64, error)) *RpcClientBindingMock_PendingNonceAt_Call {
	_c.Call.Return(run)
	return _c
}

// SendTransaction provides a mock function with given fields: ctx, tx
func (_m *RpcClientBindingMock) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for SendTransaction")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.Transaction) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RpcClientBindingMock_SendTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendTransaction'
type RpcClientBindingMock_SendTransaction_Call struct {
	*mock.Call
}

// SendTransaction is a helper method to define mock.On call
//   - ctx context.Context
//   - tx *types.Transaction
func (_e *RpcClientBindingMock_Expecter) SendTransaction(ctx interface{}, tx interface{}) *RpcClientBindingMock_SendTransaction_Call {
	return &RpcClientBindingMock_SendTransaction_Call{Call: _e.mock.On("SendTransaction", ctx, tx)}
}

func (_c *RpcClientBindingMock_SendTransaction_Call) Run(run func(ctx context.Context, tx *types.Transaction)) *RpcClientBindingMock_SendTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*types.Transaction))
	})
	return _c
}

func (_c *RpcClientBindingMock_SendTransaction_Call) Return(_a0 error) *RpcClientBindingMock_SendTransaction_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *RpcClientBindingMock_SendTransaction_Call) RunAndReturn(run func(context.Context, *types.Transaction) error) *RpcClientBindingMock_SendTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// SuggestGasPrice provides a mock function with given fields: ctx
func (_m *RpcClientBindingMock) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for SuggestGasPrice")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*big.Int, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *big.Int); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_SuggestGasPrice_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SuggestGasPrice'
type RpcClientBindingMock_SuggestGasPrice_Call struct {
	*mock.Call
}

// SuggestGasPrice is a helper method to define mock.On call
//   - ctx context.Context
func (_e *RpcClientBindingMock_Expecter) SuggestGasPrice(ctx interface{}) *RpcClientBindingMock_SuggestGasPrice_Call {
	return &RpcClientBindingMock_SuggestGasPrice_Call{Call: _e.mock.On("SuggestGasPrice", ctx)}
}

func (_c *RpcClientBindingMock_SuggestGasPrice_Call) Run(run func(ctx context.Context)) *RpcClientBindingMock_SuggestGasPrice_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *RpcClientBindingMock_SuggestGasPrice_Call) Return(_a0 *big.Int, _a1 error) *RpcClientBindingMock_SuggestGasPrice_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_SuggestGasPrice_Call) RunAndReturn(run func(context.Context) (*big.Int, error)) *RpcClientBindingMock_SuggestGasPrice_Call {
	_c.Call.Return(run)
	return _c
}

// TransactionByHash provides a mock function with given fields: ctx, hash
func (_m *RpcClientBindingMock) TransactionByHash(ctx context.Context, hash common.Hash) (*types.Transaction, bool, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for TransactionByHash")
	}

	var r0 *types.Transaction
	var r1 bool
	var r2 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*types.Transaction, bool, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *types.Transaction); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) bool); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Get(1).(bool)
	}

	if rf, ok := ret.Get(2).(func(context.Context, common.Hash) error); ok {
		r2 = rf(ctx, hash)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// RpcClientBindingMock_TransactionByHash_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TransactionByHash'
type RpcClientBindingMock_TransactionByHash_Call struct {
	*mock.Call
}

// TransactionByHash is a helper method to define mock.On call
//   - ctx context.Context
//   - hash common.Hash
func (_e *RpcClientBindingMock_Expecter) TransactionByHash(ctx interface{}, hash interface{}) *RpcClientBindingMock_TransactionByHash_Call {
	return &RpcClientBindingMock_TransactionByHash_Call{Call: _e.mock.On("TransactionByHash", ctx, hash)}
}

func (_c *RpcClientBindingMock_TransactionByHash_Call) Run(run func(ctx context.Context, hash common.Hash)) *RpcClientBindingMock_TransactionByHash_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Hash))
	})
	return _c
}

func (_c *RpcClientBindingMock_TransactionByHash_Call) Return(tx *types.Transaction, isPending bool, err error) *RpcClientBindingMock_TransactionByHash_Call {
	_c.Call.Return(tx, isPending, err)
	return _c
}

func (_c *RpcClientBindingMock_TransactionByHash_Call) RunAndReturn(run func(context.Context, common.Hash) (*types.Transaction, bool, error)) *RpcClientBindingMock_TransactionByHash_Call {
	_c.Call.Return(run)
	return _c
}

// TransactionReceipt provides a mock function with given fields: ctx, txHash
func (_m *RpcClientBindingMock) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	ret := _m.Called(ctx, txHash)

	if len(ret) == 0 {
		panic("no return value specified for TransactionReceipt")
	}

	var r0 *types.Receipt
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) (*types.Receipt, error)); ok {
		return rf(ctx, txHash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Hash) *types.Receipt); ok {
		r0 = rf(ctx, txHash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Receipt)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Hash) error); ok {
		r1 = rf(ctx, txHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RpcClientBindingMock_TransactionReceipt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'TransactionReceipt'
type RpcClientBindingMock_TransactionReceipt_Call struct {
	*mock.Call
}

// TransactionReceipt is a helper method to define mock.On call
//   - ctx context.Context
//   - txHash common.Hash
func (_e *RpcClientBindingMock_Expecter) TransactionReceipt(ctx interface{}, txHash interface{}) *RpcClientBindingMock_TransactionReceipt_Call {
	return &RpcClientBindingMock_TransactionReceipt_Call{Call: _e.mock.On("TransactionReceipt", ctx, txHash)}
}

func (_c *RpcClientBindingMock_TransactionReceipt_Call) Run(run func(ctx context.Context, txHash common.Hash)) *RpcClientBindingMock_TransactionReceipt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(common.Hash))
	})
	return _c
}

func (_c *RpcClientBindingMock_TransactionReceipt_Call) Return(_a0 *types.Receipt, _a1 error) *RpcClientBindingMock_TransactionReceipt_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RpcClientBindingMock_TransactionReceipt_Call) RunAndReturn(run func(context.Context, common.Hash) (*types.Receipt, error)) *RpcClientBindingMock_TransactionReceipt_Call {
	_c.Call.Return(run)
	return _c
}

// NewRpcClientBindingMock creates a new instance of RpcClientBindingMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRpcClientBindingMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *RpcClientBindingMock {
	mock := &RpcClientBindingMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}