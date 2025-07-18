// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"
	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// TransactionSignerMock is an autogenerated mock type for the TransactionSigner type
type TransactionSignerMock struct {
	mock.Mock
}

type TransactionSignerMock_Expecter struct {
	mock *mock.Mock
}

func (_m *TransactionSignerMock) EXPECT() *TransactionSignerMock_Expecter {
	return &TransactionSignerMock_Expecter{mock: &_m.Mock}
}

// Address provides a mock function with no fields
func (_m *TransactionSignerMock) Address() common.Address {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Address")
	}

	var r0 common.Address
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	return r0
}

// TransactionSignerMock_Address_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Address'
type TransactionSignerMock_Address_Call struct {
	*mock.Call
}

// Address is a helper method to define mock.On call
func (_e *TransactionSignerMock_Expecter) Address() *TransactionSignerMock_Address_Call {
	return &TransactionSignerMock_Address_Call{Call: _e.mock.On("Address")}
}

func (_c *TransactionSignerMock_Address_Call) Run(run func()) *TransactionSignerMock_Address_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TransactionSignerMock_Address_Call) Return(_a0 common.Address) *TransactionSignerMock_Address_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransactionSignerMock_Address_Call) RunAndReturn(run func() common.Address) *TransactionSignerMock_Address_Call {
	_c.Call.Return(run)
	return _c
}

// Sign provides a mock function with given fields: _a0, _a1
func (_m *TransactionSignerMock) Sign(_a0 common.Address, _a1 *types.Transaction) (*types.Transaction, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Sign")
	}

	var r0 *types.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Address, *types.Transaction) (*types.Transaction, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(common.Address, *types.Transaction) *types.Transaction); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Address, *types.Transaction) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionSignerMock_Sign_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Sign'
type TransactionSignerMock_Sign_Call struct {
	*mock.Call
}

// Sign is a helper method to define mock.On call
//   - _a0 common.Address
//   - _a1 *types.Transaction
func (_e *TransactionSignerMock_Expecter) Sign(_a0 interface{}, _a1 interface{}) *TransactionSignerMock_Sign_Call {
	return &TransactionSignerMock_Sign_Call{Call: _e.mock.On("Sign", _a0, _a1)}
}

func (_c *TransactionSignerMock_Sign_Call) Run(run func(_a0 common.Address, _a1 *types.Transaction)) *TransactionSignerMock_Sign_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.Address), args[1].(*types.Transaction))
	})
	return _c
}

func (_c *TransactionSignerMock_Sign_Call) Return(_a0 *types.Transaction, _a1 error) *TransactionSignerMock_Sign_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TransactionSignerMock_Sign_Call) RunAndReturn(run func(common.Address, *types.Transaction) (*types.Transaction, error)) *TransactionSignerMock_Sign_Call {
	_c.Call.Return(run)
	return _c
}

// SignBytes provides a mock function with given fields: msg
func (_m *TransactionSignerMock) SignBytes(msg []byte) ([]byte, error) {
	ret := _m.Called(msg)

	if len(ret) == 0 {
		panic("no return value specified for SignBytes")
	}

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) ([]byte, error)); ok {
		return rf(msg)
	}
	if rf, ok := ret.Get(0).(func([]byte) []byte); ok {
		r0 = rf(msg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(msg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransactionSignerMock_SignBytes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SignBytes'
type TransactionSignerMock_SignBytes_Call struct {
	*mock.Call
}

// SignBytes is a helper method to define mock.On call
//   - msg []byte
func (_e *TransactionSignerMock_Expecter) SignBytes(msg interface{}) *TransactionSignerMock_SignBytes_Call {
	return &TransactionSignerMock_SignBytes_Call{Call: _e.mock.On("SignBytes", msg)}
}

func (_c *TransactionSignerMock_SignBytes_Call) Run(run func(msg []byte)) *TransactionSignerMock_SignBytes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *TransactionSignerMock_SignBytes_Call) Return(_a0 []byte, _a1 error) *TransactionSignerMock_SignBytes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *TransactionSignerMock_SignBytes_Call) RunAndReturn(run func([]byte) ([]byte, error)) *TransactionSignerMock_SignBytes_Call {
	_c.Call.Return(run)
	return _c
}

// Validate provides a mock function with given fields: signature, hash
func (_m *TransactionSignerMock) Validate(signature string, hash string) bool {
	ret := _m.Called(signature, hash)

	if len(ret) == 0 {
		panic("no return value specified for Validate")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(signature, hash)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// TransactionSignerMock_Validate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Validate'
type TransactionSignerMock_Validate_Call struct {
	*mock.Call
}

// Validate is a helper method to define mock.On call
//   - signature string
//   - hash string
func (_e *TransactionSignerMock_Expecter) Validate(signature interface{}, hash interface{}) *TransactionSignerMock_Validate_Call {
	return &TransactionSignerMock_Validate_Call{Call: _e.mock.On("Validate", signature, hash)}
}

func (_c *TransactionSignerMock_Validate_Call) Run(run func(signature string, hash string)) *TransactionSignerMock_Validate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *TransactionSignerMock_Validate_Call) Return(_a0 bool) *TransactionSignerMock_Validate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TransactionSignerMock_Validate_Call) RunAndReturn(run func(string, string) bool) *TransactionSignerMock_Validate_Call {
	_c.Call.Return(run)
	return _c
}

// NewTransactionSignerMock creates a new instance of TransactionSignerMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionSignerMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionSignerMock {
	mock := &TransactionSignerMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
