package mocks

import (
	"context"

	"github.com/rsksmart/liquidity-provider-server/internal/usecases/pegin"
	"github.com/stretchr/testify/mock"
)

type GetPeginQuoteUseCaseMock struct {
	mock.Mock
}

type GetPeginQuoteUseCaseMock_Expecter struct {
	mock *mock.Mock
}

func (_m *GetPeginQuoteUseCaseMock) EXPECT() *GetPeginQuoteUseCaseMock_Expecter {
	return &GetPeginQuoteUseCaseMock_Expecter{mock: &_m.Mock}
}

func (_m *GetPeginQuoteUseCaseMock) Run(ctx context.Context, request pegin.QuoteRequest) (pegin.GetPeginQuoteResult, error) {
	ret := _m.Called(ctx, request)

	if len(ret) == 0 {
		panic("no return value specified for Run")
	}

	var r0 pegin.GetPeginQuoteResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, pegin.QuoteRequest) (pegin.GetPeginQuoteResult, error)); ok {
		return rf(ctx, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, pegin.QuoteRequest) pegin.GetPeginQuoteResult); ok {
		r0 = rf(ctx, request)
	} else {
		r0 = ret.Get(0).(pegin.GetPeginQuoteResult)
	}

	if rf, ok := ret.Get(1).(func(context.Context, pegin.QuoteRequest) error); ok {
		r1 = rf(ctx, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type GetPeginQuoteUseCaseMock_Run_Call struct {
	*mock.Call
}

func (_e *GetPeginQuoteUseCaseMock_Expecter) Run(ctx interface{}, request interface{}) *GetPeginQuoteUseCaseMock_Run_Call {
	return &GetPeginQuoteUseCaseMock_Run_Call{Call: _e.mock.On("Run", ctx, request)}
}

func (_c *GetPeginQuoteUseCaseMock_Run_Call) Run(run func(ctx context.Context, request pegin.QuoteRequest)) *GetPeginQuoteUseCaseMock_Run_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(pegin.QuoteRequest))
	})
	return _c
}

func (_c *GetPeginQuoteUseCaseMock_Run_Call) Return(_a0 pegin.GetPeginQuoteResult, _a1 error) *GetPeginQuoteUseCaseMock_Run_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GetPeginQuoteUseCaseMock_Run_Call) RunAndReturn(run func(context.Context, pegin.QuoteRequest) (pegin.GetPeginQuoteResult, error)) *GetPeginQuoteUseCaseMock_Run_Call {
	_c.Call.Return(run)
	return _c
}

func NewGetPeginQuoteUseCaseMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *GetPeginQuoteUseCaseMock {
	mock := &GetPeginQuoteUseCaseMock{}
	mock.Mock.Test(t)
	t.Cleanup(func() { mock.AssertExpectations(t) })
	return mock
}
