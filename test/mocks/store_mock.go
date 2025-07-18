// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import (
	http "net/http"

	sessions "github.com/gorilla/sessions"
	mock "github.com/stretchr/testify/mock"
)

// StoreMock is an autogenerated mock type for the Store type
type StoreMock struct {
	mock.Mock
}

type StoreMock_Expecter struct {
	mock *mock.Mock
}

func (_m *StoreMock) EXPECT() *StoreMock_Expecter {
	return &StoreMock_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: r, name
func (_m *StoreMock) Get(r *http.Request, name string) (*sessions.Session, error) {
	ret := _m.Called(r, name)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *sessions.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request, string) (*sessions.Session, error)); ok {
		return rf(r, name)
	}
	if rf, ok := ret.Get(0).(func(*http.Request, string) *sessions.Session); ok {
		r0 = rf(r, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sessions.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request, string) error); ok {
		r1 = rf(r, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreMock_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type StoreMock_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - r *http.Request
//   - name string
func (_e *StoreMock_Expecter) Get(r interface{}, name interface{}) *StoreMock_Get_Call {
	return &StoreMock_Get_Call{Call: _e.mock.On("Get", r, name)}
}

func (_c *StoreMock_Get_Call) Run(run func(r *http.Request, name string)) *StoreMock_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request), args[1].(string))
	})
	return _c
}

func (_c *StoreMock_Get_Call) Return(_a0 *sessions.Session, _a1 error) *StoreMock_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StoreMock_Get_Call) RunAndReturn(run func(*http.Request, string) (*sessions.Session, error)) *StoreMock_Get_Call {
	_c.Call.Return(run)
	return _c
}

// New provides a mock function with given fields: r, name
func (_m *StoreMock) New(r *http.Request, name string) (*sessions.Session, error) {
	ret := _m.Called(r, name)

	if len(ret) == 0 {
		panic("no return value specified for New")
	}

	var r0 *sessions.Session
	var r1 error
	if rf, ok := ret.Get(0).(func(*http.Request, string) (*sessions.Session, error)); ok {
		return rf(r, name)
	}
	if rf, ok := ret.Get(0).(func(*http.Request, string) *sessions.Session); ok {
		r0 = rf(r, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sessions.Session)
		}
	}

	if rf, ok := ret.Get(1).(func(*http.Request, string) error); ok {
		r1 = rf(r, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// StoreMock_New_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'New'
type StoreMock_New_Call struct {
	*mock.Call
}

// New is a helper method to define mock.On call
//   - r *http.Request
//   - name string
func (_e *StoreMock_Expecter) New(r interface{}, name interface{}) *StoreMock_New_Call {
	return &StoreMock_New_Call{Call: _e.mock.On("New", r, name)}
}

func (_c *StoreMock_New_Call) Run(run func(r *http.Request, name string)) *StoreMock_New_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request), args[1].(string))
	})
	return _c
}

func (_c *StoreMock_New_Call) Return(_a0 *sessions.Session, _a1 error) *StoreMock_New_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *StoreMock_New_Call) RunAndReturn(run func(*http.Request, string) (*sessions.Session, error)) *StoreMock_New_Call {
	_c.Call.Return(run)
	return _c
}

// Save provides a mock function with given fields: r, w, s
func (_m *StoreMock) Save(r *http.Request, w http.ResponseWriter, s *sessions.Session) error {
	ret := _m.Called(r, w, s)

	if len(ret) == 0 {
		panic("no return value specified for Save")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*http.Request, http.ResponseWriter, *sessions.Session) error); ok {
		r0 = rf(r, w, s)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StoreMock_Save_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Save'
type StoreMock_Save_Call struct {
	*mock.Call
}

// Save is a helper method to define mock.On call
//   - r *http.Request
//   - w http.ResponseWriter
//   - s *sessions.Session
func (_e *StoreMock_Expecter) Save(r interface{}, w interface{}, s interface{}) *StoreMock_Save_Call {
	return &StoreMock_Save_Call{Call: _e.mock.On("Save", r, w, s)}
}

func (_c *StoreMock_Save_Call) Run(run func(r *http.Request, w http.ResponseWriter, s *sessions.Session)) *StoreMock_Save_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(*http.Request), args[1].(http.ResponseWriter), args[2].(*sessions.Session))
	})
	return _c
}

func (_c *StoreMock_Save_Call) Return(_a0 error) *StoreMock_Save_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *StoreMock_Save_Call) RunAndReturn(run func(*http.Request, http.ResponseWriter, *sessions.Session) error) *StoreMock_Save_Call {
	_c.Call.Return(run)
	return _c
}

// NewStoreMock creates a new instance of StoreMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStoreMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *StoreMock {
	mock := &StoreMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
