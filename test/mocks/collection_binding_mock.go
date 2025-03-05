// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	mongo "go.mongodb.org/mongo-driver/mongo"

	options "go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionBindingMock is an autogenerated mock type for the CollectionBinding type
type CollectionBindingMock struct {
	mock.Mock
}

type CollectionBindingMock_Expecter struct {
	mock *mock.Mock
}

func (_m *CollectionBindingMock) EXPECT() *CollectionBindingMock_Expecter {
	return &CollectionBindingMock_Expecter{mock: &_m.Mock}
}

// BulkWrite provides a mock function with given fields: ctx, models, opts
func (_m *CollectionBindingMock) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, models)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for BulkWrite")
	}

	var r0 *mongo.BulkWriteResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []mongo.WriteModel, ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)); ok {
		return rf(ctx, models, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []mongo.WriteModel, ...*options.BulkWriteOptions) *mongo.BulkWriteResult); ok {
		r0 = rf(ctx, models, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.BulkWriteResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []mongo.WriteModel, ...*options.BulkWriteOptions) error); ok {
		r1 = rf(ctx, models, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionBindingMock_BulkWrite_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'BulkWrite'
type CollectionBindingMock_BulkWrite_Call struct {
	*mock.Call
}

// BulkWrite is a helper method to define mock.On call
//   - ctx context.Context
//   - models []mongo.WriteModel
//   - opts ...*options.BulkWriteOptions
func (_e *CollectionBindingMock_Expecter) BulkWrite(ctx interface{}, models interface{}, opts ...interface{}) *CollectionBindingMock_BulkWrite_Call {
	return &CollectionBindingMock_BulkWrite_Call{Call: _e.mock.On("BulkWrite",
		append([]interface{}{ctx, models}, opts...)...)}
}

func (_c *CollectionBindingMock_BulkWrite_Call) Run(run func(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions)) *CollectionBindingMock_BulkWrite_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.BulkWriteOptions, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*options.BulkWriteOptions)
			}
		}
		run(args[0].(context.Context), args[1].([]mongo.WriteModel), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_BulkWrite_Call) Return(_a0 *mongo.BulkWriteResult, _a1 error) *CollectionBindingMock_BulkWrite_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionBindingMock_BulkWrite_Call) RunAndReturn(run func(context.Context, []mongo.WriteModel, ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)) *CollectionBindingMock_BulkWrite_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteMany provides a mock function with given fields: ctx, filter, opts
func (_m *CollectionBindingMock) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteMany")
	}

	var r0 *mongo.DeleteResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)); ok {
		return rf(ctx, filter, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.DeleteOptions) *mongo.DeleteResult); ok {
		r0 = rf(ctx, filter, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.DeleteResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.DeleteOptions) error); ok {
		r1 = rf(ctx, filter, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionBindingMock_DeleteMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteMany'
type CollectionBindingMock_DeleteMany_Call struct {
	*mock.Call
}

// DeleteMany is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - opts ...*options.DeleteOptions
func (_e *CollectionBindingMock_Expecter) DeleteMany(ctx interface{}, filter interface{}, opts ...interface{}) *CollectionBindingMock_DeleteMany_Call {
	return &CollectionBindingMock_DeleteMany_Call{Call: _e.mock.On("DeleteMany",
		append([]interface{}{ctx, filter}, opts...)...)}
}

func (_c *CollectionBindingMock_DeleteMany_Call) Run(run func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions)) *CollectionBindingMock_DeleteMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.DeleteOptions, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*options.DeleteOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_DeleteMany_Call) Return(_a0 *mongo.DeleteResult, _a1 error) *CollectionBindingMock_DeleteMany_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionBindingMock_DeleteMany_Call) RunAndReturn(run func(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)) *CollectionBindingMock_DeleteMany_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteOne provides a mock function with given fields: ctx, filter, opts
func (_m *CollectionBindingMock) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for DeleteOne")
	}

	var r0 *mongo.DeleteResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)); ok {
		return rf(ctx, filter, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.DeleteOptions) *mongo.DeleteResult); ok {
		r0 = rf(ctx, filter, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.DeleteResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.DeleteOptions) error); ok {
		r1 = rf(ctx, filter, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionBindingMock_DeleteOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteOne'
type CollectionBindingMock_DeleteOne_Call struct {
	*mock.Call
}

// DeleteOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - opts ...*options.DeleteOptions
func (_e *CollectionBindingMock_Expecter) DeleteOne(ctx interface{}, filter interface{}, opts ...interface{}) *CollectionBindingMock_DeleteOne_Call {
	return &CollectionBindingMock_DeleteOne_Call{Call: _e.mock.On("DeleteOne",
		append([]interface{}{ctx, filter}, opts...)...)}
}

func (_c *CollectionBindingMock_DeleteOne_Call) Run(run func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions)) *CollectionBindingMock_DeleteOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.DeleteOptions, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*options.DeleteOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_DeleteOne_Call) Return(_a0 *mongo.DeleteResult, _a1 error) *CollectionBindingMock_DeleteOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionBindingMock_DeleteOne_Call) RunAndReturn(run func(context.Context, interface{}, ...*options.DeleteOptions) (*mongo.DeleteResult, error)) *CollectionBindingMock_DeleteOne_Call {
	_c.Call.Return(run)
	return _c
}

// Find provides a mock function with given fields: ctx, filter, opts
func (_m *CollectionBindingMock) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Find")
	}

	var r0 *mongo.Cursor
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)); ok {
		return rf(ctx, filter, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.FindOptions) *mongo.Cursor); ok {
		r0 = rf(ctx, filter, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.Cursor)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.FindOptions) error); ok {
		r1 = rf(ctx, filter, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionBindingMock_Find_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Find'
type CollectionBindingMock_Find_Call struct {
	*mock.Call
}

// Find is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - opts ...*options.FindOptions
func (_e *CollectionBindingMock_Expecter) Find(ctx interface{}, filter interface{}, opts ...interface{}) *CollectionBindingMock_Find_Call {
	return &CollectionBindingMock_Find_Call{Call: _e.mock.On("Find",
		append([]interface{}{ctx, filter}, opts...)...)}
}

func (_c *CollectionBindingMock_Find_Call) Run(run func(ctx context.Context, filter interface{}, opts ...*options.FindOptions)) *CollectionBindingMock_Find_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.FindOptions, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*options.FindOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_Find_Call) Return(_a0 *mongo.Cursor, _a1 error) *CollectionBindingMock_Find_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionBindingMock_Find_Call) RunAndReturn(run func(context.Context, interface{}, ...*options.FindOptions) (*mongo.Cursor, error)) *CollectionBindingMock_Find_Call {
	_c.Call.Return(run)
	return _c
}

// FindOne provides a mock function with given fields: ctx, filter, opts
func (_m *CollectionBindingMock) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for FindOne")
	}

	var r0 *mongo.SingleResult
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult); ok {
		r0 = rf(ctx, filter, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.SingleResult)
		}
	}

	return r0
}

// CollectionBindingMock_FindOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'FindOne'
type CollectionBindingMock_FindOne_Call struct {
	*mock.Call
}

// FindOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - opts ...*options.FindOneOptions
func (_e *CollectionBindingMock_Expecter) FindOne(ctx interface{}, filter interface{}, opts ...interface{}) *CollectionBindingMock_FindOne_Call {
	return &CollectionBindingMock_FindOne_Call{Call: _e.mock.On("FindOne",
		append([]interface{}{ctx, filter}, opts...)...)}
}

func (_c *CollectionBindingMock_FindOne_Call) Run(run func(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions)) *CollectionBindingMock_FindOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.FindOneOptions, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*options.FindOneOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_FindOne_Call) Return(_a0 *mongo.SingleResult) *CollectionBindingMock_FindOne_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *CollectionBindingMock_FindOne_Call) RunAndReturn(run func(context.Context, interface{}, ...*options.FindOneOptions) *mongo.SingleResult) *CollectionBindingMock_FindOne_Call {
	_c.Call.Return(run)
	return _c
}

// InsertMany provides a mock function with given fields: ctx, documents, opts
func (_m *CollectionBindingMock) InsertMany(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, documents)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for InsertMany")
	}

	var r0 *mongo.InsertManyResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []interface{}, ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)); ok {
		return rf(ctx, documents, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []interface{}, ...*options.InsertManyOptions) *mongo.InsertManyResult); ok {
		r0 = rf(ctx, documents, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.InsertManyResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []interface{}, ...*options.InsertManyOptions) error); ok {
		r1 = rf(ctx, documents, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionBindingMock_InsertMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertMany'
type CollectionBindingMock_InsertMany_Call struct {
	*mock.Call
}

// InsertMany is a helper method to define mock.On call
//   - ctx context.Context
//   - documents []interface{}
//   - opts ...*options.InsertManyOptions
func (_e *CollectionBindingMock_Expecter) InsertMany(ctx interface{}, documents interface{}, opts ...interface{}) *CollectionBindingMock_InsertMany_Call {
	return &CollectionBindingMock_InsertMany_Call{Call: _e.mock.On("InsertMany",
		append([]interface{}{ctx, documents}, opts...)...)}
}

func (_c *CollectionBindingMock_InsertMany_Call) Run(run func(ctx context.Context, documents []interface{}, opts ...*options.InsertManyOptions)) *CollectionBindingMock_InsertMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.InsertManyOptions, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*options.InsertManyOptions)
			}
		}
		run(args[0].(context.Context), args[1].([]interface{}), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_InsertMany_Call) Return(_a0 *mongo.InsertManyResult, _a1 error) *CollectionBindingMock_InsertMany_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionBindingMock_InsertMany_Call) RunAndReturn(run func(context.Context, []interface{}, ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)) *CollectionBindingMock_InsertMany_Call {
	_c.Call.Return(run)
	return _c
}

// InsertOne provides a mock function with given fields: ctx, document, opts
func (_m *CollectionBindingMock) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, document)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for InsertOne")
	}

	var r0 *mongo.InsertOneResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)); ok {
		return rf(ctx, document, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...*options.InsertOneOptions) *mongo.InsertOneResult); ok {
		r0 = rf(ctx, document, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.InsertOneResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...*options.InsertOneOptions) error); ok {
		r1 = rf(ctx, document, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionBindingMock_InsertOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'InsertOne'
type CollectionBindingMock_InsertOne_Call struct {
	*mock.Call
}

// InsertOne is a helper method to define mock.On call
//   - ctx context.Context
//   - document interface{}
//   - opts ...*options.InsertOneOptions
func (_e *CollectionBindingMock_Expecter) InsertOne(ctx interface{}, document interface{}, opts ...interface{}) *CollectionBindingMock_InsertOne_Call {
	return &CollectionBindingMock_InsertOne_Call{Call: _e.mock.On("InsertOne",
		append([]interface{}{ctx, document}, opts...)...)}
}

func (_c *CollectionBindingMock_InsertOne_Call) Run(run func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions)) *CollectionBindingMock_InsertOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.InsertOneOptions, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(*options.InsertOneOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_InsertOne_Call) Return(_a0 *mongo.InsertOneResult, _a1 error) *CollectionBindingMock_InsertOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionBindingMock_InsertOne_Call) RunAndReturn(run func(context.Context, interface{}, ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)) *CollectionBindingMock_InsertOne_Call {
	_c.Call.Return(run)
	return _c
}

// ReplaceOne provides a mock function with given fields: ctx, filter, replacement, opts
func (_m *CollectionBindingMock) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, replacement)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for ReplaceOne")
	}

	var r0 *mongo.UpdateResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.ReplaceOptions) (*mongo.UpdateResult, error)); ok {
		return rf(ctx, filter, replacement, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.ReplaceOptions) *mongo.UpdateResult); ok {
		r0 = rf(ctx, filter, replacement, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.UpdateResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.ReplaceOptions) error); ok {
		r1 = rf(ctx, filter, replacement, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionBindingMock_ReplaceOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReplaceOne'
type CollectionBindingMock_ReplaceOne_Call struct {
	*mock.Call
}

// ReplaceOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - replacement interface{}
//   - opts ...*options.ReplaceOptions
func (_e *CollectionBindingMock_Expecter) ReplaceOne(ctx interface{}, filter interface{}, replacement interface{}, opts ...interface{}) *CollectionBindingMock_ReplaceOne_Call {
	return &CollectionBindingMock_ReplaceOne_Call{Call: _e.mock.On("ReplaceOne",
		append([]interface{}{ctx, filter, replacement}, opts...)...)}
}

func (_c *CollectionBindingMock_ReplaceOne_Call) Run(run func(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions)) *CollectionBindingMock_ReplaceOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.ReplaceOptions, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(*options.ReplaceOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_ReplaceOne_Call) Return(_a0 *mongo.UpdateResult, _a1 error) *CollectionBindingMock_ReplaceOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionBindingMock_ReplaceOne_Call) RunAndReturn(run func(context.Context, interface{}, interface{}, ...*options.ReplaceOptions) (*mongo.UpdateResult, error)) *CollectionBindingMock_ReplaceOne_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateMany provides a mock function with given fields: ctx, filter, update, opts
func (_m *CollectionBindingMock) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, update)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpdateMany")
	}

	var r0 *mongo.UpdateResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)); ok {
		return rf(ctx, filter, update, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) *mongo.UpdateResult); ok {
		r0 = rf(ctx, filter, update, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.UpdateResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) error); ok {
		r1 = rf(ctx, filter, update, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionBindingMock_UpdateMany_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateMany'
type CollectionBindingMock_UpdateMany_Call struct {
	*mock.Call
}

// UpdateMany is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - update interface{}
//   - opts ...*options.UpdateOptions
func (_e *CollectionBindingMock_Expecter) UpdateMany(ctx interface{}, filter interface{}, update interface{}, opts ...interface{}) *CollectionBindingMock_UpdateMany_Call {
	return &CollectionBindingMock_UpdateMany_Call{Call: _e.mock.On("UpdateMany",
		append([]interface{}{ctx, filter, update}, opts...)...)}
}

func (_c *CollectionBindingMock_UpdateMany_Call) Run(run func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions)) *CollectionBindingMock_UpdateMany_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.UpdateOptions, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(*options.UpdateOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_UpdateMany_Call) Return(_a0 *mongo.UpdateResult, _a1 error) *CollectionBindingMock_UpdateMany_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionBindingMock_UpdateMany_Call) RunAndReturn(run func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)) *CollectionBindingMock_UpdateMany_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateOne provides a mock function with given fields: ctx, filter, update, opts
func (_m *CollectionBindingMock) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, filter, update)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOne")
	}

	var r0 *mongo.UpdateResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)); ok {
		return rf(ctx, filter, update, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) *mongo.UpdateResult); ok {
		r0 = rf(ctx, filter, update, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*mongo.UpdateResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) error); ok {
		r1 = rf(ctx, filter, update, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CollectionBindingMock_UpdateOne_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateOne'
type CollectionBindingMock_UpdateOne_Call struct {
	*mock.Call
}

// UpdateOne is a helper method to define mock.On call
//   - ctx context.Context
//   - filter interface{}
//   - update interface{}
//   - opts ...*options.UpdateOptions
func (_e *CollectionBindingMock_Expecter) UpdateOne(ctx interface{}, filter interface{}, update interface{}, opts ...interface{}) *CollectionBindingMock_UpdateOne_Call {
	return &CollectionBindingMock_UpdateOne_Call{Call: _e.mock.On("UpdateOne",
		append([]interface{}{ctx, filter, update}, opts...)...)}
}

func (_c *CollectionBindingMock_UpdateOne_Call) Run(run func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions)) *CollectionBindingMock_UpdateOne_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]*options.UpdateOptions, len(args)-3)
		for i, a := range args[3:] {
			if a != nil {
				variadicArgs[i] = a.(*options.UpdateOptions)
			}
		}
		run(args[0].(context.Context), args[1].(interface{}), args[2].(interface{}), variadicArgs...)
	})
	return _c
}

func (_c *CollectionBindingMock_UpdateOne_Call) Return(_a0 *mongo.UpdateResult, _a1 error) *CollectionBindingMock_UpdateOne_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *CollectionBindingMock_UpdateOne_Call) RunAndReturn(run func(context.Context, interface{}, interface{}, ...*options.UpdateOptions) (*mongo.UpdateResult, error)) *CollectionBindingMock_UpdateOne_Call {
	_c.Call.Return(run)
	return _c
}

// NewCollectionBindingMock creates a new instance of CollectionBindingMock. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCollectionBindingMock(t interface {
	mock.TestingT
	Cleanup(func())
}) *CollectionBindingMock {
	mock := &CollectionBindingMock{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
