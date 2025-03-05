// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	big "math/big"

	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// RskBridgeBindingMock is an autogenerated mock type for the RskBridgeBinding type
type RskBridgeBindingMock struct {
	mock.Mock
}

type RskBridgeBindingMock_Expecter struct {
	mock *mock.Mock
}

func (_m *RskBridgeBindingMock) EXPECT() *RskBridgeBindingMock_Expecter {
	return &RskBridgeBindingMock_Expecter{mock: &_m.Mock}
}

// GetActiveFederationCreationBlockHeight provides a mock function with given fields: opts
func (_m *RskBridgeBindingMock) GetActiveFederationCreationBlockHeight(opts *bind.CallOpts) (*big.Int, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for GetActiveFederationCreationBlockHeight")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (*big.Int, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) *big.Int); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActiveFederationCreationBlockHeight'
type RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call struct {
	*mock.Call
}

// GetActiveFederationCreationBlockHeight is a helper method to define mock.On call
//   - opts *bind.CallOpts
func (_e *RskBridgeBindingMock_Expecter) GetActiveFederationCreationBlockHeight(opts interface{}) *RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call {
	return &RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call{Call: _e.mock.On("GetActiveFederationCreationBlockHeight", opts)}
}

func (_c *RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call) Run(run func(opts *bind.CallOpts)) *RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts))
	})
	return _c
}

func (_c *RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call) Return(_a0 *big.Int, _a1 error) *RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call) RunAndReturn(run func(*bind.CallOpts) (*big.Int, error)) *RskBridgeBindingMock_GetActiveFederationCreationBlockHeight_Call {
	_c.Call.Return(run)
	return _c
}

// GetActivePowpegRedeemScript provides a mock function with given fields: opts
func (_m *RskBridgeBindingMock) GetActivePowpegRedeemScript(opts *bind.CallOpts) ([]byte, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for GetActivePowpegRedeemScript")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) ([]byte, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) []byte); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_GetActivePowpegRedeemScript_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetActivePowpegRedeemScript'
type RskBridgeBindingMock_GetActivePowpegRedeemScript_Call struct {
	*mock.Call
}

// GetActivePowpegRedeemScript is a helper method to define mock.On call
//   - opts *bind.CallOpts
func (_e *RskBridgeBindingMock_Expecter) GetActivePowpegRedeemScript(opts interface{}) *RskBridgeBindingMock_GetActivePowpegRedeemScript_Call {
	return &RskBridgeBindingMock_GetActivePowpegRedeemScript_Call{Call: _e.mock.On("GetActivePowpegRedeemScript", opts)}
}

func (_c *RskBridgeBindingMock_GetActivePowpegRedeemScript_Call) Run(run func(opts *bind.CallOpts)) *RskBridgeBindingMock_GetActivePowpegRedeemScript_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts))
	})
	return _c
}

func (_c *RskBridgeBindingMock_GetActivePowpegRedeemScript_Call) Return(_a0 []byte, _a1 error) *RskBridgeBindingMock_GetActivePowpegRedeemScript_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_GetActivePowpegRedeemScript_Call) RunAndReturn(run func(*bind.CallOpts) ([]byte, error)) *RskBridgeBindingMock_GetActivePowpegRedeemScript_Call {
	_c.Call.Return(run)
	return _c
}

// GetBtcBlockchainBestChainHeight provides a mock function with given fields: opts
func (_m *RskBridgeBindingMock) GetBtcBlockchainBestChainHeight(opts *bind.CallOpts) (*big.Int, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for GetBtcBlockchainBestChainHeight")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (*big.Int, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) *big.Int); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetBtcBlockchainBestChainHeight'
type RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call struct {
	*mock.Call
}

// GetBtcBlockchainBestChainHeight is a helper method to define mock.On call
//   - opts *bind.CallOpts
func (_e *RskBridgeBindingMock_Expecter) GetBtcBlockchainBestChainHeight(opts interface{}) *RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call {
	return &RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call{Call: _e.mock.On("GetBtcBlockchainBestChainHeight", opts)}
}

func (_c *RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call) Run(run func(opts *bind.CallOpts)) *RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts))
	})
	return _c
}

func (_c *RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call) Return(_a0 *big.Int, _a1 error) *RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call) RunAndReturn(run func(*bind.CallOpts) (*big.Int, error)) *RskBridgeBindingMock_GetBtcBlockchainBestChainHeight_Call {
	_c.Call.Return(run)
	return _c
}

// GetFederationAddress provides a mock function with given fields: opts
func (_m *RskBridgeBindingMock) GetFederationAddress(opts *bind.CallOpts) (string, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for GetFederationAddress")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (string, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) string); ok {
		r0 = rf(opts)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_GetFederationAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFederationAddress'
type RskBridgeBindingMock_GetFederationAddress_Call struct {
	*mock.Call
}

// GetFederationAddress is a helper method to define mock.On call
//   - opts *bind.CallOpts
func (_e *RskBridgeBindingMock_Expecter) GetFederationAddress(opts interface{}) *RskBridgeBindingMock_GetFederationAddress_Call {
	return &RskBridgeBindingMock_GetFederationAddress_Call{Call: _e.mock.On("GetFederationAddress", opts)}
}

func (_c *RskBridgeBindingMock_GetFederationAddress_Call) Run(run func(opts *bind.CallOpts)) *RskBridgeBindingMock_GetFederationAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts))
	})
	return _c
}

func (_c *RskBridgeBindingMock_GetFederationAddress_Call) Return(_a0 string, _a1 error) *RskBridgeBindingMock_GetFederationAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_GetFederationAddress_Call) RunAndReturn(run func(*bind.CallOpts) (string, error)) *RskBridgeBindingMock_GetFederationAddress_Call {
	_c.Call.Return(run)
	return _c
}

// GetFederationSize provides a mock function with given fields: opts
func (_m *RskBridgeBindingMock) GetFederationSize(opts *bind.CallOpts) (*big.Int, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for GetFederationSize")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (*big.Int, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) *big.Int); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_GetFederationSize_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFederationSize'
type RskBridgeBindingMock_GetFederationSize_Call struct {
	*mock.Call
}

// GetFederationSize is a helper method to define mock.On call
//   - opts *bind.CallOpts
func (_e *RskBridgeBindingMock_Expecter) GetFederationSize(opts interface{}) *RskBridgeBindingMock_GetFederationSize_Call {
	return &RskBridgeBindingMock_GetFederationSize_Call{Call: _e.mock.On("GetFederationSize", opts)}
}

func (_c *RskBridgeBindingMock_GetFederationSize_Call) Run(run func(opts *bind.CallOpts)) *RskBridgeBindingMock_GetFederationSize_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts))
	})
	return _c
}

func (_c *RskBridgeBindingMock_GetFederationSize_Call) Return(_a0 *big.Int, _a1 error) *RskBridgeBindingMock_GetFederationSize_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_GetFederationSize_Call) RunAndReturn(run func(*bind.CallOpts) (*big.Int, error)) *RskBridgeBindingMock_GetFederationSize_Call {
	_c.Call.Return(run)
	return _c
}

// GetFederationThreshold provides a mock function with given fields: opts
func (_m *RskBridgeBindingMock) GetFederationThreshold(opts *bind.CallOpts) (*big.Int, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for GetFederationThreshold")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (*big.Int, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) *big.Int); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_GetFederationThreshold_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFederationThreshold'
type RskBridgeBindingMock_GetFederationThreshold_Call struct {
	*mock.Call
}

// GetFederationThreshold is a helper method to define mock.On call
//   - opts *bind.CallOpts
func (_e *RskBridgeBindingMock_Expecter) GetFederationThreshold(opts interface{}) *RskBridgeBindingMock_GetFederationThreshold_Call {
	return &RskBridgeBindingMock_GetFederationThreshold_Call{Call: _e.mock.On("GetFederationThreshold", opts)}
}

func (_c *RskBridgeBindingMock_GetFederationThreshold_Call) Run(run func(opts *bind.CallOpts)) *RskBridgeBindingMock_GetFederationThreshold_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts))
	})
	return _c
}

func (_c *RskBridgeBindingMock_GetFederationThreshold_Call) Return(_a0 *big.Int, _a1 error) *RskBridgeBindingMock_GetFederationThreshold_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_GetFederationThreshold_Call) RunAndReturn(run func(*bind.CallOpts) (*big.Int, error)) *RskBridgeBindingMock_GetFederationThreshold_Call {
	_c.Call.Return(run)
	return _c
}

// GetFederatorPublicKeyOfType provides a mock function with given fields: opts, index, atype
func (_m *RskBridgeBindingMock) GetFederatorPublicKeyOfType(opts *bind.CallOpts, index *big.Int, atype string) ([]byte, error) {
	ret := _m.Called(opts, index, atype)

	if len(ret) == 0 {
		panic("no return value specified for GetFederatorPublicKeyOfType")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, *big.Int, string) ([]byte, error)); ok {
		return rf(opts, index, atype)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, *big.Int, string) []byte); ok {
		r0 = rf(opts, index, atype)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, *big.Int, string) error); ok {
		r1 = rf(opts, index, atype)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFederatorPublicKeyOfType'
type RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call struct {
	*mock.Call
}

// GetFederatorPublicKeyOfType is a helper method to define mock.On call
//   - opts *bind.CallOpts
//   - index *big.Int
//   - atype string
func (_e *RskBridgeBindingMock_Expecter) GetFederatorPublicKeyOfType(opts interface{}, index interface{}, atype interface{}) *RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call {
	return &RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call{Call: _e.mock.On("GetFederatorPublicKeyOfType", opts, index, atype)}
}

func (_c *RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call) Run(run func(opts *bind.CallOpts, index *big.Int, atype string)) *RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts), args[1].(*big.Int), args[2].(string))
	})
	return _c
}

func (_c *RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call) Return(_a0 []byte, _a1 error) *RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call) RunAndReturn(run func(*bind.CallOpts, *big.Int, string) ([]byte, error)) *RskBridgeBindingMock_GetFederatorPublicKeyOfType_Call {
	_c.Call.Return(run)
	return _c
}

// GetMinimumLockTxValue provides a mock function with given fields: opts
func (_m *RskBridgeBindingMock) GetMinimumLockTxValue(opts *bind.CallOpts) (*big.Int, error) {
	ret := _m.Called(opts)

	if len(ret) == 0 {
		panic("no return value specified for GetMinimumLockTxValue")
	}

	var r0 *big.Int
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (*big.Int, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) *big.Int); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*big.Int)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_GetMinimumLockTxValue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetMinimumLockTxValue'
type RskBridgeBindingMock_GetMinimumLockTxValue_Call struct {
	*mock.Call
}

// GetMinimumLockTxValue is a helper method to define mock.On call
//   - opts *bind.CallOpts
func (_e *RskBridgeBindingMock_Expecter) GetMinimumLockTxValue(opts interface{}) *RskBridgeBindingMock_GetMinimumLockTxValue_Call {
	return &RskBridgeBindingMock_GetMinimumLockTxValue_Call{Call: _e.mock.On("GetMinimumLockTxValue", opts)}
}

func (_c *RskBridgeBindingMock_GetMinimumLockTxValue_Call) Run(run func(opts *bind.CallOpts)) *RskBridgeBindingMock_GetMinimumLockTxValue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts))
	})
	return _c
}

func (_c *RskBridgeBindingMock_GetMinimumLockTxValue_Call) Return(_a0 *big.Int, _a1 error) *RskBridgeBindingMock_GetMinimumLockTxValue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_GetMinimumLockTxValue_Call) RunAndReturn(run func(*bind.CallOpts) (*big.Int, error)) *RskBridgeBindingMock_GetMinimumLockTxValue_Call {
	_c.Call.Return(run)
	return _c
}

// HasBtcBlockCoinbaseTransactionInformation provides a mock function with given fields: opts, blockHash
func (_m *RskBridgeBindingMock) HasBtcBlockCoinbaseTransactionInformation(opts *bind.CallOpts, blockHash [32]byte) (bool, error) {
	ret := _m.Called(opts, blockHash)

	if len(ret) == 0 {
		panic("no return value specified for HasBtcBlockCoinbaseTransactionInformation")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, [32]byte) (bool, error)); ok {
		return rf(opts, blockHash)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, [32]byte) bool); ok {
		r0 = rf(opts, blockHash)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, [32]byte) error); ok {
		r1 = rf(opts, blockHash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'HasBtcBlockCoinbaseTransactionInformation'
type RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call struct {
	*mock.Call
}

// HasBtcBlockCoinbaseTransactionInformation is a helper method to define mock.On call
//   - opts *bind.CallOpts
//   - blockHash [32]byte
func (_e *RskBridgeBindingMock_Expecter) HasBtcBlockCoinbaseTransactionInformation(opts interface{}, blockHash interface{}) *RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call {
	return &RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call{Call: _e.mock.On("HasBtcBlockCoinbaseTransactionInformation", opts, blockHash)}
}

func (_c *RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call) Run(run func(opts *bind.CallOpts, blockHash [32]byte)) *RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts), args[1].([32]byte))
	})
	return _c
}

func (_c *RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call) Return(_a0 bool, _a1 error) *RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call) RunAndReturn(run func(*bind.CallOpts, [32]byte) (bool, error)) *RskBridgeBindingMock_HasBtcBlockCoinbaseTransactionInformation_Call {
	_c.Call.Return(run)
	return _c
}

// IsBtcTxHashAlreadyProcessed provides a mock function with given fields: opts, hash
func (_m *RskBridgeBindingMock) IsBtcTxHashAlreadyProcessed(opts *bind.CallOpts, hash string) (bool, error) {
	ret := _m.Called(opts, hash)

	if len(ret) == 0 {
		panic("no return value specified for IsBtcTxHashAlreadyProcessed")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, string) (bool, error)); ok {
		return rf(opts, hash)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, string) bool); ok {
		r0 = rf(opts, hash)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, string) error); ok {
		r1 = rf(opts, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsBtcTxHashAlreadyProcessed'
type RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call struct {
	*mock.Call
}

// IsBtcTxHashAlreadyProcessed is a helper method to define mock.On call
//   - opts *bind.CallOpts
//   - hash string
func (_e *RskBridgeBindingMock_Expecter) IsBtcTxHashAlreadyProcessed(opts interface{}, hash interface{}) *RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call {
	return &RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call{Call: _e.mock.On("IsBtcTxHashAlreadyProcessed", opts, hash)}
}

func (_c *RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call) Run(run func(opts *bind.CallOpts, hash string)) *RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.CallOpts), args[1].(string))
	})
	return _c
}

func (_c *RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call) Return(_a0 bool, _a1 error) *RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call) RunAndReturn(run func(*bind.CallOpts, string) (bool, error)) *RskBridgeBindingMock_IsBtcTxHashAlreadyProcessed_Call {
	_c.Call.Return(run)
	return _c
}

// RegisterBtcCoinbaseTransaction provides a mock function with given fields: opts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue
func (_m *RskBridgeBindingMock) RegisterBtcCoinbaseTransaction(opts *bind.TransactOpts, btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte) (*types.Transaction, error) {
	ret := _m.Called(opts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)

	if len(ret) == 0 {
		panic("no return value specified for RegisterBtcCoinbaseTransaction")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, []byte, [32]byte, []byte, [32]byte, [32]byte) (*types.Transaction, error)); ok {
		return rf(opts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
	}
	if rf, ok := ret.Get(0).(func(*bind.TransactOpts, []byte, [32]byte, []byte, [32]byte, [32]byte) *types.Transaction); ok {
		r0 = rf(opts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(*bind.TransactOpts, []byte, [32]byte, []byte, [32]byte, [32]byte) error); ok {
		r1 = rf(opts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RegisterBtcCoinbaseTransaction'
type RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call struct {
	*mock.Call
}

// RegisterBtcCoinbaseTransaction is a helper method to define mock.On call
//   - opts *bind.TransactOpts
//   - btcTxSerialized []byte
//   - blockHash [32]byte
//   - pmtSerialized []byte
//   - witnessMerkleRoot [32]byte
//   - witnessReservedValue [32]byte
func (_e *RskBridgeBindingMock_Expecter) RegisterBtcCoinbaseTransaction(opts interface{}, btcTxSerialized interface{}, blockHash interface{}, pmtSerialized interface{}, witnessMerkleRoot interface{}, witnessReservedValue interface{}) *RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call {
	return &RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call{Call: _e.mock.On("RegisterBtcCoinbaseTransaction", opts, btcTxSerialized, blockHash, pmtSerialized, witnessMerkleRoot, witnessReservedValue)}
}

func (_c *RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call) Run(run func(opts *bind.TransactOpts, btcTxSerialized []byte, blockHash [32]byte, pmtSerialized []byte, witnessMerkleRoot [32]byte, witnessReservedValue [32]byte)) *RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*bind.TransactOpts), args[1].([]byte), args[2].([32]byte), args[3].([]byte), args[4].([32]byte), args[5].([32]byte))
	})
	return _c
}

func (_c *RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call) Return(_a0 *types.Transaction, _a1 error) *RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call) RunAndReturn(run func(*bind.TransactOpts, []byte, [32]byte, []byte, [32]byte, [32]byte) (*types.Transaction, error)) *RskBridgeBindingMock_RegisterBtcCoinbaseTransaction_Call {
	_c.Call.Return(run)
	return _c
}

// NewRskBridgeBindingMock creates a new instance of RskBridgeBindingMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRskBridgeBindingMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *RskBridgeBindingMock {
	mock := &RskBridgeBindingMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
