// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	mock "github.com/stretchr/testify/mock"
)

// LbcCallerBindingMock is an autogenerated mock type for the LbcCallerBinding type
type LbcCallerBindingMock struct {
	mock.Mock
}

type LbcCallerBindingMock_Expecter struct {
	mock *mock.Mock
}

func (_m *LbcCallerBindingMock) EXPECT() *LbcCallerBindingMock_Expecter {
	return &LbcCallerBindingMock_Expecter{mock: &_m.Mock}
}

// Call provides a mock function with given fields: opts, result, method, params
func (_m *LbcCallerBindingMock) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	var _ca []interface{}
	_ca = append(_ca, opts, result, method)
	_ca = append(_ca, params...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Call")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, *[]interface{}, string, ...interface{}) error); ok {
		r0 = rf(opts, result, method, params...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// LbcCallerBindingMock_Call_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Call'
type LbcCallerBindingMock_Call_Call struct {
	*mock.Call
}

// Call is a helper method to define mock.On call
//   - opts *bind.CallOpts
//   - result *[]interface{}
//   - method string
//   - params ...interface{}
func (_e *LbcCallerBindingMock_Expecter) Call(opts interface{}, result interface{}, method interface{}, params ...interface{}) *LbcCallerBindingMock_Call_Call {
	return &LbcCallerBindingMock_Call_Call{Call: _e.mock.On("Call",
		append([]interface{}{opts, result, method}, params...)...)}
}

func (_c *LbcCallerBindingMock_Call_Call) Run(run func(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{})) *LbcCallerBindingMock_Call_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(*bind.CallOpts), args[1].(*[]interface{}), args[2].(string), variadicArgs...)
	})
	return _c
}

func (_c *LbcCallerBindingMock_Call_Call) Return(_a0 error) *LbcCallerBindingMock_Call_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *LbcCallerBindingMock_Call_Call) RunAndReturn(run func(*bind.CallOpts, *[]interface{}, string, ...interface{}) error) *LbcCallerBindingMock_Call_Call {
	_c.Call.Return(run)
	return _c
}

// NewLbcCallerBindingMock creates a new instance of LbcCallerBindingMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLbcCallerBindingMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *LbcCallerBindingMock {
	mock := &LbcCallerBindingMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}