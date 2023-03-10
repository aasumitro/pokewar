// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ICreateBulkRepository is an autogenerated mock type for the ICreateBulkRepository type
type ICreateBulkRepository[T interface{}] struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, param
func (_m *ICreateBulkRepository[T]) Create(ctx context.Context, param []*T) error {
	ret := _m.Called(ctx, param)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []*T) error); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewICreateBulkRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewICreateBulkRepository creates a new instance of ICreateBulkRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewICreateBulkRepository[T interface{}](t mockConstructorTestingTNewICreateBulkRepository) *ICreateBulkRepository[T] {
	mock := &ICreateBulkRepository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
