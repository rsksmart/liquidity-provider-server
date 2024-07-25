// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	ses "github.com/aws/aws-sdk-go-v2/service/ses"
	mock "github.com/stretchr/testify/mock"
)

// SesClientMock is an autogenerated mock type for the sesClient type
type SesClientMock struct {
	mock.Mock
}

type SesClientMock_Expecter struct {
	mock *mock.Mock
}

func (_m *SesClientMock) EXPECT() *SesClientMock_Expecter {
	return &SesClientMock_Expecter{mock: &_m.Mock}
}

// SendEmail provides a mock function with given fields: ctx, params, optFns
func (_m *SesClientMock) SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
	_va := make([]interface{}, len(optFns))
	for _i := range optFns {
		_va[_i] = optFns[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, params)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for SendEmail")
	}

	var r0 *ses.SendEmailOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ses.SendEmailInput, ...func(*ses.Options)) (*ses.SendEmailOutput, error)); ok {
		return rf(ctx, params, optFns...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ses.SendEmailInput, ...func(*ses.Options)) *ses.SendEmailOutput); ok {
		r0 = rf(ctx, params, optFns...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ses.SendEmailOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ses.SendEmailInput, ...func(*ses.Options)) error); ok {
		r1 = rf(ctx, params, optFns...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SesClientMock_SendEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'SendEmail'
type SesClientMock_SendEmail_Call struct {
	*mock.Call
}

// SendEmail is a helper method to define mock.On call
//   - ctx context.Context
//   - params *ses.SendEmailInput
//   - optFns ...func(*ses.Options)
func (_e *SesClientMock_Expecter) SendEmail(ctx interface{}, params interface{}, optFns ...interface{}) *SesClientMock_SendEmail_Call {
	return &SesClientMock_SendEmail_Call{Call: _e.mock.On("SendEmail",
		append([]interface{}{ctx, params}, optFns...)...)}
}

func (_c *SesClientMock_SendEmail_Call) Run(run func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options))) *SesClientMock_SendEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]func(*ses.Options), len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(func(*ses.Options))
			}
		}
		run(args[0].(context.Context), args[1].(*ses.SendEmailInput), variadicArgs...)
	})
	return _c
}

func (_c *SesClientMock_SendEmail_Call) Return(_a0 *ses.SendEmailOutput, _a1 error) *SesClientMock_SendEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *SesClientMock_SendEmail_Call) RunAndReturn(run func(context.Context, *ses.SendEmailInput, ...func(*ses.Options)) (*ses.SendEmailOutput, error)) *SesClientMock_SendEmail_Call {
	_c.Call.Return(run)
	return _c
}

// NewSesClientMock creates a new instance of SesClientMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewSesClientMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *SesClientMock {
	mock := &SesClientMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
