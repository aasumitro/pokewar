// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// IUpdateRepository is an autogenerated mock type for the IUpdateRepository type
type IUpdateRepository[T interface{}] struct {
	mock.Mock
}

// Update provides a mock function with given fields: ctx, param
func (_m *IUpdateRepository[T]) Update(ctx context.Context, param *T) error {
	ret := _m.Called(ctx, param)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *T) error); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewIUpdateRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIUpdateRepository creates a new instance of IUpdateRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIUpdateRepository[T interface{}](t mockConstructorTestingTNewIUpdateRepository) *IUpdateRepository[T] {
	mock := &IUpdateRepository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
