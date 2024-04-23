// Code generated by mockery v2.42.1. DO NOT EDIT.

package mocks

import (
	context "context"

	quote "github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	mock "github.com/stretchr/testify/mock"
)

// PeginQuoteRepositoryMock is an autogenerated mock type for the PeginQuoteRepository type
type PeginQuoteRepositoryMock struct {
	mock.Mock
}

type PeginQuoteRepositoryMock_Expecter struct {
	mock *mock.Mock
}

func (_m *PeginQuoteRepositoryMock) EXPECT() *PeginQuoteRepositoryMock_Expecter {
	return &PeginQuoteRepositoryMock_Expecter{mock: &_m.Mock}
}

// DeleteQuotes provides a mock function with given fields: ctx, quotes
func (_m *PeginQuoteRepositoryMock) DeleteQuotes(ctx context.Context, quotes []string) (uint, error) {
	ret := _m.Called(ctx, quotes)

	if len(ret) == 0 {
		panic("no return value specified for DeleteQuotes")
	}

	var r0 uint
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []string) (uint, error)); ok {
		return rf(ctx, quotes)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []string) uint); ok {
		r0 = rf(ctx, quotes)
	} else {
		r0 = ret.Get(0).(uint)
	}

	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(ctx, quotes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PeginQuoteRepositoryMock_DeleteQuotes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteQuotes'
type PeginQuoteRepositoryMock_DeleteQuotes_Call struct {
	*mock.Call
}

// DeleteQuotes is a helper method to define mock.On call
//   - ctx context.Context
//   - quotes []string
func (_e *PeginQuoteRepositoryMock_Expecter) DeleteQuotes(ctx interface{}, quotes interface{}) *PeginQuoteRepositoryMock_DeleteQuotes_Call {
	return &PeginQuoteRepositoryMock_DeleteQuotes_Call{Call: _e.mock.On("DeleteQuotes", ctx, quotes)}
}

func (_c *PeginQuoteRepositoryMock_DeleteQuotes_Call) Run(run func(ctx context.Context, quotes []string)) *PeginQuoteRepositoryMock_DeleteQuotes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *PeginQuoteRepositoryMock_DeleteQuotes_Call) Return(_a0 uint, _a1 error) *PeginQuoteRepositoryMock_DeleteQuotes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PeginQuoteRepositoryMock_DeleteQuotes_Call) RunAndReturn(run func(context.Context, []string) (uint, error)) *PeginQuoteRepositoryMock_DeleteQuotes_Call {
	_c.Call.Return(run)
	return _c
}

// GetQuote provides a mock function with given fields: ctx, hash
func (_m *PeginQuoteRepositoryMock) GetQuote(ctx context.Context, hash string) (*quote.PeginQuote, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetQuote")
	}

	var r0 *quote.PeginQuote
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*quote.PeginQuote, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *quote.PeginQuote); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*quote.PeginQuote)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PeginQuoteRepositoryMock_GetQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQuote'
type PeginQuoteRepositoryMock_GetQuote_Call struct {
	*mock.Call
}

// GetQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - hash string
func (_e *PeginQuoteRepositoryMock_Expecter) GetQuote(ctx interface{}, hash interface{}) *PeginQuoteRepositoryMock_GetQuote_Call {
	return &PeginQuoteRepositoryMock_GetQuote_Call{Call: _e.mock.On("GetQuote", ctx, hash)}
}

func (_c *PeginQuoteRepositoryMock_GetQuote_Call) Run(run func(ctx context.Context, hash string)) *PeginQuoteRepositoryMock_GetQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *PeginQuoteRepositoryMock_GetQuote_Call) Return(_a0 *quote.PeginQuote, _a1 error) *PeginQuoteRepositoryMock_GetQuote_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PeginQuoteRepositoryMock_GetQuote_Call) RunAndReturn(run func(context.Context, string) (*quote.PeginQuote, error)) *PeginQuoteRepositoryMock_GetQuote_Call {
	_c.Call.Return(run)
	return _c
}

// GetRetainedQuote provides a mock function with given fields: ctx, hash
func (_m *PeginQuoteRepositoryMock) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPeginQuote, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetRetainedQuote")
	}

	var r0 *quote.RetainedPeginQuote
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*quote.RetainedPeginQuote, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *quote.RetainedPeginQuote); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*quote.RetainedPeginQuote)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PeginQuoteRepositoryMock_GetRetainedQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRetainedQuote'
type PeginQuoteRepositoryMock_GetRetainedQuote_Call struct {
	*mock.Call
}

// GetRetainedQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - hash string
func (_e *PeginQuoteRepositoryMock_Expecter) GetRetainedQuote(ctx interface{}, hash interface{}) *PeginQuoteRepositoryMock_GetRetainedQuote_Call {
	return &PeginQuoteRepositoryMock_GetRetainedQuote_Call{Call: _e.mock.On("GetRetainedQuote", ctx, hash)}
}

func (_c *PeginQuoteRepositoryMock_GetRetainedQuote_Call) Run(run func(ctx context.Context, hash string)) *PeginQuoteRepositoryMock_GetRetainedQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *PeginQuoteRepositoryMock_GetRetainedQuote_Call) Return(_a0 *quote.RetainedPeginQuote, _a1 error) *PeginQuoteRepositoryMock_GetRetainedQuote_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PeginQuoteRepositoryMock_GetRetainedQuote_Call) RunAndReturn(run func(context.Context, string) (*quote.RetainedPeginQuote, error)) *PeginQuoteRepositoryMock_GetRetainedQuote_Call {
	_c.Call.Return(run)
	return _c
}

// GetRetainedQuoteByState provides a mock function with given fields: ctx, states
func (_m *PeginQuoteRepositoryMock) GetRetainedQuoteByState(ctx context.Context, states ...quote.PeginState) ([]quote.RetainedPeginQuote, error) {
	_va := make([]interface{}, len(states))
	for _i := range states {
		_va[_i] = states[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetRetainedQuoteByState")
	}

	var r0 []quote.RetainedPeginQuote
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...quote.PeginState) ([]quote.RetainedPeginQuote, error)); ok {
		return rf(ctx, states...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...quote.PeginState) []quote.RetainedPeginQuote); ok {
		r0 = rf(ctx, states...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]quote.RetainedPeginQuote)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...quote.PeginState) error); ok {
		r1 = rf(ctx, states...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRetainedQuoteByState'
type PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call struct {
	*mock.Call
}

// GetRetainedQuoteByState is a helper method to define mock.On call
//   - ctx context.Context
//   - states ...quote.PeginState
func (_e *PeginQuoteRepositoryMock_Expecter) GetRetainedQuoteByState(ctx interface{}, states ...interface{}) *PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call {
	return &PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call{Call: _e.mock.On("GetRetainedQuoteByState",
		append([]interface{}{ctx}, states...)...)}
}

func (_c *PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call) Run(run func(ctx context.Context, states ...quote.PeginState)) *PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]quote.PeginState, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(quote.PeginState)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call) Return(_a0 []quote.RetainedPeginQuote, _a1 error) *PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call) RunAndReturn(run func(context.Context, ...quote.PeginState) ([]quote.RetainedPeginQuote, error)) *PeginQuoteRepositoryMock_GetRetainedQuoteByState_Call {
	_c.Call.Return(run)
	return _c
}

// InsertQuote provides a mock function with given fields: ctx, hash, _a2
func (_m *PeginQuoteRepositoryMock) InsertQuote(ctx context.Context, hash string, _a2 quote.PeginQuote) error {
	ret := _m.Called(ctx, hash, _a2)

	if len(ret) == 0 {
		panic("no return value specified for InsertQuote")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, quote.PeginQuote) error); ok {
		r0 = rf(ctx, hash, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PeginQuoteRepositoryMock_InsertQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertQuote'
type PeginQuoteRepositoryMock_InsertQuote_Call struct {
	*mock.Call
}

// InsertQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - hash string
//   - _a2 quote.PeginQuote
func (_e *PeginQuoteRepositoryMock_Expecter) InsertQuote(ctx interface{}, hash interface{}, _a2 interface{}) *PeginQuoteRepositoryMock_InsertQuote_Call {
	return &PeginQuoteRepositoryMock_InsertQuote_Call{Call: _e.mock.On("InsertQuote", ctx, hash, _a2)}
}

func (_c *PeginQuoteRepositoryMock_InsertQuote_Call) Run(run func(ctx context.Context, hash string, _a2 quote.PeginQuote)) *PeginQuoteRepositoryMock_InsertQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(quote.PeginQuote))
	})
	return _c
}

func (_c *PeginQuoteRepositoryMock_InsertQuote_Call) Return(_a0 error) *PeginQuoteRepositoryMock_InsertQuote_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PeginQuoteRepositoryMock_InsertQuote_Call) RunAndReturn(run func(context.Context, string, quote.PeginQuote) error) *PeginQuoteRepositoryMock_InsertQuote_Call {
	_c.Call.Return(run)
	return _c
}

// InsertRetainedQuote provides a mock function with given fields: ctx, _a1
func (_m *PeginQuoteRepositoryMock) InsertRetainedQuote(ctx context.Context, _a1 quote.RetainedPeginQuote) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for InsertRetainedQuote")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, quote.RetainedPeginQuote) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PeginQuoteRepositoryMock_InsertRetainedQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertRetainedQuote'
type PeginQuoteRepositoryMock_InsertRetainedQuote_Call struct {
	*mock.Call
}

// InsertRetainedQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 quote.RetainedPeginQuote
func (_e *PeginQuoteRepositoryMock_Expecter) InsertRetainedQuote(ctx interface{}, _a1 interface{}) *PeginQuoteRepositoryMock_InsertRetainedQuote_Call {
	return &PeginQuoteRepositoryMock_InsertRetainedQuote_Call{Call: _e.mock.On("InsertRetainedQuote", ctx, _a1)}
}

func (_c *PeginQuoteRepositoryMock_InsertRetainedQuote_Call) Run(run func(ctx context.Context, _a1 quote.RetainedPeginQuote)) *PeginQuoteRepositoryMock_InsertRetainedQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(quote.RetainedPeginQuote))
	})
	return _c
}

func (_c *PeginQuoteRepositoryMock_InsertRetainedQuote_Call) Return(_a0 error) *PeginQuoteRepositoryMock_InsertRetainedQuote_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PeginQuoteRepositoryMock_InsertRetainedQuote_Call) RunAndReturn(run func(context.Context, quote.RetainedPeginQuote) error) *PeginQuoteRepositoryMock_InsertRetainedQuote_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateRetainedQuote provides a mock function with given fields: ctx, _a1
func (_m *PeginQuoteRepositoryMock) UpdateRetainedQuote(ctx context.Context, _a1 quote.RetainedPeginQuote) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRetainedQuote")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, quote.RetainedPeginQuote) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PeginQuoteRepositoryMock_UpdateRetainedQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateRetainedQuote'
type PeginQuoteRepositoryMock_UpdateRetainedQuote_Call struct {
	*mock.Call
}

// UpdateRetainedQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 quote.RetainedPeginQuote
func (_e *PeginQuoteRepositoryMock_Expecter) UpdateRetainedQuote(ctx interface{}, _a1 interface{}) *PeginQuoteRepositoryMock_UpdateRetainedQuote_Call {
	return &PeginQuoteRepositoryMock_UpdateRetainedQuote_Call{Call: _e.mock.On("UpdateRetainedQuote", ctx, _a1)}
}

func (_c *PeginQuoteRepositoryMock_UpdateRetainedQuote_Call) Run(run func(ctx context.Context, _a1 quote.RetainedPeginQuote)) *PeginQuoteRepositoryMock_UpdateRetainedQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(quote.RetainedPeginQuote))
	})
	return _c
}

func (_c *PeginQuoteRepositoryMock_UpdateRetainedQuote_Call) Return(_a0 error) *PeginQuoteRepositoryMock_UpdateRetainedQuote_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PeginQuoteRepositoryMock_UpdateRetainedQuote_Call) RunAndReturn(run func(context.Context, quote.RetainedPeginQuote) error) *PeginQuoteRepositoryMock_UpdateRetainedQuote_Call {
	_c.Call.Return(run)
	return _c
}

// NewPeginQuoteRepositoryMock creates a new instance of PeginQuoteRepositoryMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPeginQuoteRepositoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *PeginQuoteRepositoryMock {
	mock := &PeginQuoteRepositoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
