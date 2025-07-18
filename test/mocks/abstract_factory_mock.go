// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	blockchain "github.com/rsksmart/liquidity-provider-server/internal/entities/blockchain"
	mock "github.com/stretchr/testify/mock"

	rootstock "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/rootstock"
)

// AbstractFactoryMock is an autogenerated mock type for the AbstractFactory type
type AbstractFactoryMock struct {
	mock.Mock
}

type AbstractFactoryMock_Expecter struct {
	mock *mock.Mock
}

func (_m *AbstractFactoryMock) EXPECT() *AbstractFactoryMock_Expecter {
	return &AbstractFactoryMock_Expecter{mock: &_m.Mock}
}

// BitcoinMonitoringWallet provides a mock function with given fields: walletId
func (_m *AbstractFactoryMock) BitcoinMonitoringWallet(walletId string) (blockchain.BitcoinWallet, error) {
	ret := _m.Called(walletId)

	if len(ret) == 0 {
		panic("no return value specified for BitcoinMonitoringWallet")
	}

	var r0 blockchain.BitcoinWallet
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (blockchain.BitcoinWallet, error)); ok {
		return rf(walletId)
	}
	if rf, ok := ret.Get(0).(func(string) blockchain.BitcoinWallet); ok {
		r0 = rf(walletId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(blockchain.BitcoinWallet)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(walletId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AbstractFactoryMock_BitcoinMonitoringWallet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BitcoinMonitoringWallet'
type AbstractFactoryMock_BitcoinMonitoringWallet_Call struct {
	*mock.Call
}

// BitcoinMonitoringWallet is a helper method to define mock.On call
//   - walletId string
func (_e *AbstractFactoryMock_Expecter) BitcoinMonitoringWallet(walletId interface{}) *AbstractFactoryMock_BitcoinMonitoringWallet_Call {
	return &AbstractFactoryMock_BitcoinMonitoringWallet_Call{Call: _e.mock.On("BitcoinMonitoringWallet", walletId)}
}

func (_c *AbstractFactoryMock_BitcoinMonitoringWallet_Call) Run(run func(walletId string)) *AbstractFactoryMock_BitcoinMonitoringWallet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AbstractFactoryMock_BitcoinMonitoringWallet_Call) Return(_a0 blockchain.BitcoinWallet, _a1 error) *AbstractFactoryMock_BitcoinMonitoringWallet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AbstractFactoryMock_BitcoinMonitoringWallet_Call) RunAndReturn(run func(string) (blockchain.BitcoinWallet, error)) *AbstractFactoryMock_BitcoinMonitoringWallet_Call {
	_c.Call.Return(run)
	return _c
}

// BitcoinPaymentWallet provides a mock function with given fields: walletId
func (_m *AbstractFactoryMock) BitcoinPaymentWallet(walletId string) (blockchain.BitcoinWallet, error) {
	ret := _m.Called(walletId)

	if len(ret) == 0 {
		panic("no return value specified for BitcoinPaymentWallet")
	}

	var r0 blockchain.BitcoinWallet
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (blockchain.BitcoinWallet, error)); ok {
		return rf(walletId)
	}
	if rf, ok := ret.Get(0).(func(string) blockchain.BitcoinWallet); ok {
		r0 = rf(walletId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(blockchain.BitcoinWallet)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(walletId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AbstractFactoryMock_BitcoinPaymentWallet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BitcoinPaymentWallet'
type AbstractFactoryMock_BitcoinPaymentWallet_Call struct {
	*mock.Call
}

// BitcoinPaymentWallet is a helper method to define mock.On call
//   - walletId string
func (_e *AbstractFactoryMock_Expecter) BitcoinPaymentWallet(walletId interface{}) *AbstractFactoryMock_BitcoinPaymentWallet_Call {
	return &AbstractFactoryMock_BitcoinPaymentWallet_Call{Call: _e.mock.On("BitcoinPaymentWallet", walletId)}
}

func (_c *AbstractFactoryMock_BitcoinPaymentWallet_Call) Run(run func(walletId string)) *AbstractFactoryMock_BitcoinPaymentWallet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *AbstractFactoryMock_BitcoinPaymentWallet_Call) Return(_a0 blockchain.BitcoinWallet, _a1 error) *AbstractFactoryMock_BitcoinPaymentWallet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AbstractFactoryMock_BitcoinPaymentWallet_Call) RunAndReturn(run func(string) (blockchain.BitcoinWallet, error)) *AbstractFactoryMock_BitcoinPaymentWallet_Call {
	_c.Call.Return(run)
	return _c
}

// RskWallet provides a mock function with no fields
func (_m *AbstractFactoryMock) RskWallet() (rootstock.RskSignerWallet, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RskWallet")
	}

	var r0 rootstock.RskSignerWallet
	var r1 error
	if rf, ok := ret.Get(0).(func() (rootstock.RskSignerWallet, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() rootstock.RskSignerWallet); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(rootstock.RskSignerWallet)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// AbstractFactoryMock_RskWallet_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RskWallet'
type AbstractFactoryMock_RskWallet_Call struct {
	*mock.Call
}

// RskWallet is a helper method to define mock.On call
func (_e *AbstractFactoryMock_Expecter) RskWallet() *AbstractFactoryMock_RskWallet_Call {
	return &AbstractFactoryMock_RskWallet_Call{Call: _e.mock.On("RskWallet")}
}

func (_c *AbstractFactoryMock_RskWallet_Call) Run(run func()) *AbstractFactoryMock_RskWallet_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *AbstractFactoryMock_RskWallet_Call) Return(_a0 rootstock.RskSignerWallet, _a1 error) *AbstractFactoryMock_RskWallet_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *AbstractFactoryMock_RskWallet_Call) RunAndReturn(run func() (rootstock.RskSignerWallet, error)) *AbstractFactoryMock_RskWallet_Call {
	_c.Call.Return(run)
	return _c
}

// NewAbstractFactoryMock creates a new instance of AbstractFactoryMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewAbstractFactoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *AbstractFactoryMock {
	mock := &AbstractFactoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
