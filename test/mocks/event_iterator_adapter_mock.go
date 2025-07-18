// Code generated by mockery v2.51.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// EventIteratorAdapterMock is an autogenerated mock type for the EventIteratorAdapter type
type EventIteratorAdapterMock[T any] struct {
	mock.Mock
}

type EventIteratorAdapterMock_Expecter[T any] struct {
	mock *mock.Mock
}

func (_m *EventIteratorAdapterMock[T]) EXPECT() *EventIteratorAdapterMock_Expecter[T] {
	return &EventIteratorAdapterMock_Expecter[T]{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *EventIteratorAdapterMock[T]) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventIteratorAdapterMock_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type EventIteratorAdapterMock_Close_Call[T any] struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *EventIteratorAdapterMock_Expecter[T]) Close() *EventIteratorAdapterMock_Close_Call[T] {
	return &EventIteratorAdapterMock_Close_Call[T]{Call: _e.mock.On("Close")}
}

func (_c *EventIteratorAdapterMock_Close_Call[T]) Run(run func()) *EventIteratorAdapterMock_Close_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EventIteratorAdapterMock_Close_Call[T]) Return(_a0 error) *EventIteratorAdapterMock_Close_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventIteratorAdapterMock_Close_Call[T]) RunAndReturn(run func() error) *EventIteratorAdapterMock_Close_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Error provides a mock function with no fields
func (_m *EventIteratorAdapterMock[T]) Error() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Error")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// EventIteratorAdapterMock_Error_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Error'
type EventIteratorAdapterMock_Error_Call[T any] struct {
	*mock.Call
}

// Error is a helper method to define mock.On call
func (_e *EventIteratorAdapterMock_Expecter[T]) Error() *EventIteratorAdapterMock_Error_Call[T] {
	return &EventIteratorAdapterMock_Error_Call[T]{Call: _e.mock.On("Error")}
}

func (_c *EventIteratorAdapterMock_Error_Call[T]) Run(run func()) *EventIteratorAdapterMock_Error_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EventIteratorAdapterMock_Error_Call[T]) Return(_a0 error) *EventIteratorAdapterMock_Error_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventIteratorAdapterMock_Error_Call[T]) RunAndReturn(run func() error) *EventIteratorAdapterMock_Error_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Event provides a mock function with no fields
func (_m *EventIteratorAdapterMock[T]) Event() *T {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Event")
	}

	var r0 *T
	if rf, ok := ret.Get(0).(func() *T); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	return r0
}

// EventIteratorAdapterMock_Event_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Event'
type EventIteratorAdapterMock_Event_Call[T any] struct {
	*mock.Call
}

// Event is a helper method to define mock.On call
func (_e *EventIteratorAdapterMock_Expecter[T]) Event() *EventIteratorAdapterMock_Event_Call[T] {
	return &EventIteratorAdapterMock_Event_Call[T]{Call: _e.mock.On("Event")}
}

func (_c *EventIteratorAdapterMock_Event_Call[T]) Run(run func()) *EventIteratorAdapterMock_Event_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EventIteratorAdapterMock_Event_Call[T]) Return(_a0 *T) *EventIteratorAdapterMock_Event_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventIteratorAdapterMock_Event_Call[T]) RunAndReturn(run func() *T) *EventIteratorAdapterMock_Event_Call[T] {
	_c.Call.Return(run)
	return _c
}

// Next provides a mock function with no fields
func (_m *EventIteratorAdapterMock[T]) Next() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Next")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// EventIteratorAdapterMock_Next_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Next'
type EventIteratorAdapterMock_Next_Call[T any] struct {
	*mock.Call
}

// Next is a helper method to define mock.On call
func (_e *EventIteratorAdapterMock_Expecter[T]) Next() *EventIteratorAdapterMock_Next_Call[T] {
	return &EventIteratorAdapterMock_Next_Call[T]{Call: _e.mock.On("Next")}
}

func (_c *EventIteratorAdapterMock_Next_Call[T]) Run(run func()) *EventIteratorAdapterMock_Next_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *EventIteratorAdapterMock_Next_Call[T]) Return(_a0 bool) *EventIteratorAdapterMock_Next_Call[T] {
	_c.Call.Return(_a0)
	return _c
}

func (_c *EventIteratorAdapterMock_Next_Call[T]) RunAndReturn(run func() bool) *EventIteratorAdapterMock_Next_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewEventIteratorAdapterMock creates a new instance of EventIteratorAdapterMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewEventIteratorAdapterMock[T any](t interface {
	mock.TestingT
	Cleanup(func())
}) *EventIteratorAdapterMock[T] {
	mock := &EventIteratorAdapterMock[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
