// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	context "context"

	quote "github.com/rsksmart/liquidity-provider-server/internal/entities/quote"
	mock "github.com/stretchr/testify/mock"
)

// PegoutQuoteRepositoryMock is an autogenerated mock type for the PegoutQuoteRepository type
type PegoutQuoteRepositoryMock struct {
	mock.Mock
}

type PegoutQuoteRepositoryMock_Expecter struct {
	mock *mock.Mock
}

func (_m *PegoutQuoteRepositoryMock) EXPECT() *PegoutQuoteRepositoryMock_Expecter {
	return &PegoutQuoteRepositoryMock_Expecter{mock: &_m.Mock}
}

// DeleteQuotes provides a mock function with given fields: ctx, quotes
func (_m *PegoutQuoteRepositoryMock) DeleteQuotes(ctx context.Context, quotes []string) (uint, error) {
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

// PegoutQuoteRepositoryMock_DeleteQuotes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteQuotes'
type PegoutQuoteRepositoryMock_DeleteQuotes_Call struct {
	*mock.Call
}

// DeleteQuotes is a helper method to define mock.On call
//   - ctx context.Context
//   - quotes []string
func (_e *PegoutQuoteRepositoryMock_Expecter) DeleteQuotes(ctx interface{}, quotes interface{}) *PegoutQuoteRepositoryMock_DeleteQuotes_Call {
	return &PegoutQuoteRepositoryMock_DeleteQuotes_Call{Call: _e.mock.On("DeleteQuotes", ctx, quotes)}
}

func (_c *PegoutQuoteRepositoryMock_DeleteQuotes_Call) Run(run func(ctx context.Context, quotes []string)) *PegoutQuoteRepositoryMock_DeleteQuotes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_DeleteQuotes_Call) Return(_a0 uint, _a1 error) *PegoutQuoteRepositoryMock_DeleteQuotes_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_DeleteQuotes_Call) RunAndReturn(run func(context.Context, []string) (uint, error)) *PegoutQuoteRepositoryMock_DeleteQuotes_Call {
	_c.Call.Return(run)
	return _c
}

// GetQuote provides a mock function with given fields: ctx, hash
func (_m *PegoutQuoteRepositoryMock) GetQuote(ctx context.Context, hash string) (*quote.PegoutQuote, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetQuote")
	}

	var r0 *quote.PegoutQuote
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*quote.PegoutQuote, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *quote.PegoutQuote); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*quote.PegoutQuote)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PegoutQuoteRepositoryMock_GetQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetQuote'
type PegoutQuoteRepositoryMock_GetQuote_Call struct {
	*mock.Call
}

// GetQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - hash string
func (_e *PegoutQuoteRepositoryMock_Expecter) GetQuote(ctx interface{}, hash interface{}) *PegoutQuoteRepositoryMock_GetQuote_Call {
	return &PegoutQuoteRepositoryMock_GetQuote_Call{Call: _e.mock.On("GetQuote", ctx, hash)}
}

func (_c *PegoutQuoteRepositoryMock_GetQuote_Call) Run(run func(ctx context.Context, hash string)) *PegoutQuoteRepositoryMock_GetQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_GetQuote_Call) Return(_a0 *quote.PegoutQuote, _a1 error) *PegoutQuoteRepositoryMock_GetQuote_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_GetQuote_Call) RunAndReturn(run func(context.Context, string) (*quote.PegoutQuote, error)) *PegoutQuoteRepositoryMock_GetQuote_Call {
	_c.Call.Return(run)
	return _c
}

// GetRetainedQuote provides a mock function with given fields: ctx, hash
func (_m *PegoutQuoteRepositoryMock) GetRetainedQuote(ctx context.Context, hash string) (*quote.RetainedPegoutQuote, error) {
	ret := _m.Called(ctx, hash)

	if len(ret) == 0 {
		panic("no return value specified for GetRetainedQuote")
	}

	var r0 *quote.RetainedPegoutQuote
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*quote.RetainedPegoutQuote, error)); ok {
		return rf(ctx, hash)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *quote.RetainedPegoutQuote); ok {
		r0 = rf(ctx, hash)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*quote.RetainedPegoutQuote)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, hash)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PegoutQuoteRepositoryMock_GetRetainedQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRetainedQuote'
type PegoutQuoteRepositoryMock_GetRetainedQuote_Call struct {
	*mock.Call
}

// GetRetainedQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - hash string
func (_e *PegoutQuoteRepositoryMock_Expecter) GetRetainedQuote(ctx interface{}, hash interface{}) *PegoutQuoteRepositoryMock_GetRetainedQuote_Call {
	return &PegoutQuoteRepositoryMock_GetRetainedQuote_Call{Call: _e.mock.On("GetRetainedQuote", ctx, hash)}
}

func (_c *PegoutQuoteRepositoryMock_GetRetainedQuote_Call) Run(run func(ctx context.Context, hash string)) *PegoutQuoteRepositoryMock_GetRetainedQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_GetRetainedQuote_Call) Return(_a0 *quote.RetainedPegoutQuote, _a1 error) *PegoutQuoteRepositoryMock_GetRetainedQuote_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_GetRetainedQuote_Call) RunAndReturn(run func(context.Context, string) (*quote.RetainedPegoutQuote, error)) *PegoutQuoteRepositoryMock_GetRetainedQuote_Call {
	_c.Call.Return(run)
	return _c
}

// GetRetainedQuoteByState provides a mock function with given fields: ctx, states
func (_m *PegoutQuoteRepositoryMock) GetRetainedQuoteByState(ctx context.Context, states ...quote.PegoutState) ([]quote.RetainedPegoutQuote, error) {
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

	var r0 []quote.RetainedPegoutQuote
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...quote.PegoutState) ([]quote.RetainedPegoutQuote, error)); ok {
		return rf(ctx, states...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...quote.PegoutState) []quote.RetainedPegoutQuote); ok {
		r0 = rf(ctx, states...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]quote.RetainedPegoutQuote)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...quote.PegoutState) error); ok {
		r1 = rf(ctx, states...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetRetainedQuoteByState'
type PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call struct {
	*mock.Call
}

// GetRetainedQuoteByState is a helper method to define mock.On call
//   - ctx context.Context
//   - states ...quote.PegoutState
func (_e *PegoutQuoteRepositoryMock_Expecter) GetRetainedQuoteByState(ctx interface{}, states ...interface{}) *PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call {
	return &PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call{Call: _e.mock.On("GetRetainedQuoteByState",
		append([]interface{}{ctx}, states...)...)}
}

func (_c *PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call) Run(run func(ctx context.Context, states ...quote.PegoutState)) *PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]quote.PegoutState, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(quote.PegoutState)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call) Return(_a0 []quote.RetainedPegoutQuote, _a1 error) *PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call) RunAndReturn(run func(context.Context, ...quote.PegoutState) ([]quote.RetainedPegoutQuote, error)) *PegoutQuoteRepositoryMock_GetRetainedQuoteByState_Call {
	_c.Call.Return(run)
	return _c
}

// InsertQuote provides a mock function with given fields: ctx, hash, _a2
func (_m *PegoutQuoteRepositoryMock) InsertQuote(ctx context.Context, hash string, _a2 quote.PegoutQuote) error {
	ret := _m.Called(ctx, hash, _a2)

	if len(ret) == 0 {
		panic("no return value specified for InsertQuote")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, quote.PegoutQuote) error); ok {
		r0 = rf(ctx, hash, _a2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PegoutQuoteRepositoryMock_InsertQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertQuote'
type PegoutQuoteRepositoryMock_InsertQuote_Call struct {
	*mock.Call
}

// InsertQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - hash string
//   - _a2 quote.PegoutQuote
func (_e *PegoutQuoteRepositoryMock_Expecter) InsertQuote(ctx interface{}, hash interface{}, _a2 interface{}) *PegoutQuoteRepositoryMock_InsertQuote_Call {
	return &PegoutQuoteRepositoryMock_InsertQuote_Call{Call: _e.mock.On("InsertQuote", ctx, hash, _a2)}
}

func (_c *PegoutQuoteRepositoryMock_InsertQuote_Call) Run(run func(ctx context.Context, hash string, _a2 quote.PegoutQuote)) *PegoutQuoteRepositoryMock_InsertQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(quote.PegoutQuote))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_InsertQuote_Call) Return(_a0 error) *PegoutQuoteRepositoryMock_InsertQuote_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_InsertQuote_Call) RunAndReturn(run func(context.Context, string, quote.PegoutQuote) error) *PegoutQuoteRepositoryMock_InsertQuote_Call {
	_c.Call.Return(run)
	return _c
}

// InsertRetainedQuote provides a mock function with given fields: ctx, _a1
func (_m *PegoutQuoteRepositoryMock) InsertRetainedQuote(ctx context.Context, _a1 quote.RetainedPegoutQuote) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for InsertRetainedQuote")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, quote.RetainedPegoutQuote) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PegoutQuoteRepositoryMock_InsertRetainedQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertRetainedQuote'
type PegoutQuoteRepositoryMock_InsertRetainedQuote_Call struct {
	*mock.Call
}

// InsertRetainedQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 quote.RetainedPegoutQuote
func (_e *PegoutQuoteRepositoryMock_Expecter) InsertRetainedQuote(ctx interface{}, _a1 interface{}) *PegoutQuoteRepositoryMock_InsertRetainedQuote_Call {
	return &PegoutQuoteRepositoryMock_InsertRetainedQuote_Call{Call: _e.mock.On("InsertRetainedQuote", ctx, _a1)}
}

func (_c *PegoutQuoteRepositoryMock_InsertRetainedQuote_Call) Run(run func(ctx context.Context, _a1 quote.RetainedPegoutQuote)) *PegoutQuoteRepositoryMock_InsertRetainedQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(quote.RetainedPegoutQuote))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_InsertRetainedQuote_Call) Return(_a0 error) *PegoutQuoteRepositoryMock_InsertRetainedQuote_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_InsertRetainedQuote_Call) RunAndReturn(run func(context.Context, quote.RetainedPegoutQuote) error) *PegoutQuoteRepositoryMock_InsertRetainedQuote_Call {
	_c.Call.Return(run)
	return _c
}

// ListPegoutDepositsByAddress provides a mock function with given fields: ctx, address
func (_m *PegoutQuoteRepositoryMock) ListPegoutDepositsByAddress(ctx context.Context, address string) ([]quote.PegoutDeposit, error) {
	ret := _m.Called(ctx, address)

	if len(ret) == 0 {
		panic("no return value specified for ListPegoutDepositsByAddress")
	}

	var r0 []quote.PegoutDeposit
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]quote.PegoutDeposit, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []quote.PegoutDeposit); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]quote.PegoutDeposit)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListPegoutDepositsByAddress'
type PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call struct {
	*mock.Call
}

// ListPegoutDepositsByAddress is a helper method to define mock.On call
//   - ctx context.Context
//   - address string
func (_e *PegoutQuoteRepositoryMock_Expecter) ListPegoutDepositsByAddress(ctx interface{}, address interface{}) *PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call {
	return &PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call{Call: _e.mock.On("ListPegoutDepositsByAddress", ctx, address)}
}

func (_c *PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call) Run(run func(ctx context.Context, address string)) *PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call) Return(_a0 []quote.PegoutDeposit, _a1 error) *PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call) RunAndReturn(run func(context.Context, string) ([]quote.PegoutDeposit, error)) *PegoutQuoteRepositoryMock_ListPegoutDepositsByAddress_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateRetainedQuote provides a mock function with given fields: ctx, _a1
func (_m *PegoutQuoteRepositoryMock) UpdateRetainedQuote(ctx context.Context, _a1 quote.RetainedPegoutQuote) error {
	ret := _m.Called(ctx, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRetainedQuote")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, quote.RetainedPegoutQuote) error); ok {
		r0 = rf(ctx, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateRetainedQuote'
type PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call struct {
	*mock.Call
}

// UpdateRetainedQuote is a helper method to define mock.On call
//   - ctx context.Context
//   - _a1 quote.RetainedPegoutQuote
func (_e *PegoutQuoteRepositoryMock_Expecter) UpdateRetainedQuote(ctx interface{}, _a1 interface{}) *PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call {
	return &PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call{Call: _e.mock.On("UpdateRetainedQuote", ctx, _a1)}
}

func (_c *PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call) Run(run func(ctx context.Context, _a1 quote.RetainedPegoutQuote)) *PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(quote.RetainedPegoutQuote))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call) Return(_a0 error) *PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call) RunAndReturn(run func(context.Context, quote.RetainedPegoutQuote) error) *PegoutQuoteRepositoryMock_UpdateRetainedQuote_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateRetainedQuotes provides a mock function with given fields: ctx, quotes
func (_m *PegoutQuoteRepositoryMock) UpdateRetainedQuotes(ctx context.Context, quotes []quote.RetainedPegoutQuote) error {
	ret := _m.Called(ctx, quotes)

	if len(ret) == 0 {
		panic("no return value specified for UpdateRetainedQuotes")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []quote.RetainedPegoutQuote) error); ok {
		r0 = rf(ctx, quotes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateRetainedQuotes'
type PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call struct {
	*mock.Call
}

// UpdateRetainedQuotes is a helper method to define mock.On call
//   - ctx context.Context
//   - quotes []quote.RetainedPegoutQuote
func (_e *PegoutQuoteRepositoryMock_Expecter) UpdateRetainedQuotes(ctx interface{}, quotes interface{}) *PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call {
	return &PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call{Call: _e.mock.On("UpdateRetainedQuotes", ctx, quotes)}
}

func (_c *PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call) Run(run func(ctx context.Context, quotes []quote.RetainedPegoutQuote)) *PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]quote.RetainedPegoutQuote))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call) Return(_a0 error) *PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call) RunAndReturn(run func(context.Context, []quote.RetainedPegoutQuote) error) *PegoutQuoteRepositoryMock_UpdateRetainedQuotes_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertPegoutDeposit provides a mock function with given fields: ctx, deposit
func (_m *PegoutQuoteRepositoryMock) UpsertPegoutDeposit(ctx context.Context, deposit quote.PegoutDeposit) error {
	ret := _m.Called(ctx, deposit)

	if len(ret) == 0 {
		panic("no return value specified for UpsertPegoutDeposit")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, quote.PegoutDeposit) error); ok {
		r0 = rf(ctx, deposit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertPegoutDeposit'
type PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call struct {
	*mock.Call
}

// UpsertPegoutDeposit is a helper method to define mock.On call
//   - ctx context.Context
//   - deposit quote.PegoutDeposit
func (_e *PegoutQuoteRepositoryMock_Expecter) UpsertPegoutDeposit(ctx interface{}, deposit interface{}) *PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call {
	return &PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call{Call: _e.mock.On("UpsertPegoutDeposit", ctx, deposit)}
}

func (_c *PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call) Run(run func(ctx context.Context, deposit quote.PegoutDeposit)) *PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(quote.PegoutDeposit))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call) Return(_a0 error) *PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call) RunAndReturn(run func(context.Context, quote.PegoutDeposit) error) *PegoutQuoteRepositoryMock_UpsertPegoutDeposit_Call {
	_c.Call.Return(run)
	return _c
}

// UpsertPegoutDeposits provides a mock function with given fields: ctx, deposits
func (_m *PegoutQuoteRepositoryMock) UpsertPegoutDeposits(ctx context.Context, deposits []quote.PegoutDeposit) error {
	ret := _m.Called(ctx, deposits)

	if len(ret) == 0 {
		panic("no return value specified for UpsertPegoutDeposits")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []quote.PegoutDeposit) error); ok {
		r0 = rf(ctx, deposits)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpsertPegoutDeposits'
type PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call struct {
	*mock.Call
}

// UpsertPegoutDeposits is a helper method to define mock.On call
//   - ctx context.Context
//   - deposits []quote.PegoutDeposit
func (_e *PegoutQuoteRepositoryMock_Expecter) UpsertPegoutDeposits(ctx interface{}, deposits interface{}) *PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call {
	return &PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call{Call: _e.mock.On("UpsertPegoutDeposits", ctx, deposits)}
}

func (_c *PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call) Run(run func(ctx context.Context, deposits []quote.PegoutDeposit)) *PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]quote.PegoutDeposit))
	})
	return _c
}

func (_c *PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call) Return(_a0 error) *PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call) RunAndReturn(run func(context.Context, []quote.PegoutDeposit) error) *PegoutQuoteRepositoryMock_UpsertPegoutDeposits_Call {
	_c.Call.Return(run)
	return _c
}

// NewPegoutQuoteRepositoryMock creates a new instance of PegoutQuoteRepositoryMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPegoutQuoteRepositoryMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *PegoutQuoteRepositoryMock {
	mock := &PegoutQuoteRepositoryMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
