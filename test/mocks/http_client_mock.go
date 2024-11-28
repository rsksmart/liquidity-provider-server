// Code generated by mockery v2.48.0. DO NOT EDIT.

package mocks

import (
	http "net/http"

	mock "github.com/stretchr/testify/mock"
)

// HttpClientMock is an autogenerated mock type for the HttpClient type
type HttpClientMock struct {
	mock.Mock
}

type HttpClientMock_Expecter struct {
	mock *mock.Mock
}

func (_m *HttpClientMock) EXPECT() *HttpClientMock_Expecter {
	return &HttpClientMock_Expecter{mock: &_m.Mock}
}

// Do provides a mock function with given fields: req
func (_m *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for Do")
	}

	var r0 *http.Response
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request) (*http.Response, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*http.Request) *http.Response); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Response)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HttpClientMock_Do_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Do'
type HttpClientMock_Do_Call struct {
	*mock.Call
}

// Do is a helper method to define mock.On call
//   - req *http.Request
func (_e *HttpClientMock_Expecter) Do(req interface{}) *HttpClientMock_Do_Call {
	return &HttpClientMock_Do_Call{Call: _e.mock.On("Do", req)}
}

func (_c *HttpClientMock_Do_Call) Run(run func(req *http.Request)) *HttpClientMock_Do_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request))
	})
	return _c
}

func (_c *HttpClientMock_Do_Call) Return(_a0 *http.Response, _a1 error) *HttpClientMock_Do_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *HttpClientMock_Do_Call) RunAndReturn(run func(*http.Request) (*http.Response, error)) *HttpClientMock_Do_Call {
	_c.Call.Return(run)
	return _c
}

// NewHttpClientMock creates a new instance of HttpClientMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewHttpClientMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *HttpClientMock {
	mock := &HttpClientMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
