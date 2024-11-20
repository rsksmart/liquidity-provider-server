// Code generated by mockery v2.46.3. DO NOT EDIT.

package mocks

import (
	mongo "github.com/rsksmart/liquidity-provider-server/internal/adapters/dataproviders/database/mongo"
	mock "github.com/stretchr/testify/mock"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

// DbBindingMock is an autogenerated mock type for the DbBinding type
type DbBindingMock struct {
	mock.Mock
}

type DbBindingMock_Expecter struct {
	mock *mock.Mock
}

func (_m *DbBindingMock) EXPECT() *DbBindingMock_Expecter {
	return &DbBindingMock_Expecter{mock: &_m.Mock}
}

// Collection provides a mock function with given fields: name, opts
func (_m *DbBindingMock) Collection(name string, opts ...*options.CollectionOptions) mongo.CollectionBinding {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, name)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Collection")
	}

	var r0 mongo.CollectionBinding
	if rf, ok := ret.Get(0).(func(string, ...*options.CollectionOptions) mongo.CollectionBinding); ok {
		r0 = rf(name, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongo.CollectionBinding)
		}
	}

	return r0
}

// DbBindingMock_Collection_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Collection'
type DbBindingMock_Collection_Call struct {
	*mock.Call
}

// Collection is a helper method to define mock.On call
//   - name string
//   - opts ...*options.CollectionOptions
func (_e *DbBindingMock_Expecter) Collection(name interface{}, opts ...interface{}) *DbBindingMock_Collection_Call {
	return &DbBindingMock_Collection_Call{Call: _e.mock.On("Collection",
		append([]interface{}{name}, opts...)...)}
}

func (_c *DbBindingMock_Collection_Call) Run(run func(name string, opts ...*options.CollectionOptions)) *DbBindingMock_Collection_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.CollectionOptions, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(*options.CollectionOptions)
			}
		}
		run(args[0].(string), variadicArgs...)
	})
	return _c
}

func (_c *DbBindingMock_Collection_Call) Return(_a0 mongo.CollectionBinding) *DbBindingMock_Collection_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *DbBindingMock_Collection_Call) RunAndReturn(run func(string, ...*options.CollectionOptions) mongo.CollectionBinding) *DbBindingMock_Collection_Call {
	_c.Call.Return(run)
	return _c
}

// NewDbBindingMock creates a new instance of DbBindingMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDbBindingMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *DbBindingMock {
	mock := &DbBindingMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
