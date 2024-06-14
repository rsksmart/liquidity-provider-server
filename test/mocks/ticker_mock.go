// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	time "time"

	mock "github.com/stretchr/testify/mock"
)

// TickerMock is an autogenerated mock type for the Ticker type
type TickerMock struct {
	mock.Mock
}

type TickerMock_Expecter struct {
	mock *mock.Mock
}

func (_m *TickerMock) EXPECT() *TickerMock_Expecter {
	return &TickerMock_Expecter{mock: &_m.Mock}
}

// C provides a mock function with given fields:
func (_m *TickerMock) C() <-chan time.Time {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for C")
	}

	var r0 <-chan time.Time
	if rf, ok := ret.Get(0).(func() <-chan time.Time); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(<-chan time.Time)
		}
	}

	return r0
}

// TickerMock_C_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'C'
type TickerMock_C_Call struct {
	*mock.Call
}

// C is a helper method to define mock.On call
func (_e *TickerMock_Expecter) C() *TickerMock_C_Call {
	return &TickerMock_C_Call{Call: _e.mock.On("C")}
}

func (_c *TickerMock_C_Call) Run(run func()) *TickerMock_C_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TickerMock_C_Call) Return(_a0 <-chan time.Time) *TickerMock_C_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *TickerMock_C_Call) RunAndReturn(run func() <-chan time.Time) *TickerMock_C_Call {
	_c.Call.Return(run)
	return _c
}

// Stop provides a mock function with given fields:
func (_m *TickerMock) Stop() {
	_m.Called()
}

// TickerMock_Stop_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stop'
type TickerMock_Stop_Call struct {
	*mock.Call
}

// Stop is a helper method to define mock.On call
func (_e *TickerMock_Expecter) Stop() *TickerMock_Stop_Call {
	return &TickerMock_Stop_Call{Call: _e.mock.On("Stop")}
}

func (_c *TickerMock_Stop_Call) Run(run func()) *TickerMock_Stop_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *TickerMock_Stop_Call) Return() *TickerMock_Stop_Call {
	_c.Call.Return()
	return _c
}

func (_c *TickerMock_Stop_Call) RunAndReturn(run func()) *TickerMock_Stop_Call {
	_c.Call.Return(run)
	return _c
}

// NewTickerMock creates a new instance of TickerMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTickerMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *TickerMock {
	mock := &TickerMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
