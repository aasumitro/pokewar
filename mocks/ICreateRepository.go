// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ICreateRepository is an autogenerated mock type for the ICreateRepository type
type ICreateRepository[T interface{}] struct {
	mock.Mock
}

// Create provides a mock function with given fields: ctx, param
func (_m *ICreateRepository[T]) Create(ctx context.Context, param *T) error {
	ret := _m.Called(ctx, param)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *T) error); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewICreateRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewICreateRepository creates a new instance of ICreateRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewICreateRepository[T interface{}](t mockConstructorTestingTNewICreateRepository) *ICreateRepository[T] {
	mock := &ICreateRepository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}