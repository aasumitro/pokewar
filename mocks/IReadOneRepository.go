// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/aasumitro/pokewar/domain"
	mock "github.com/stretchr/testify/mock"
)

// IReadOneRepository is an autogenerated mock type for the IReadOneRepository type
type IReadOneRepository[T interface{}] struct {
	mock.Mock
}

// Find provides a mock function with given fields: ctx, key, val
func (_m *IReadOneRepository[T]) Find(ctx context.Context, key domain.FindWith, val interface{}) (*T, error) {
	ret := _m.Called(ctx, key, val)

	var r0 *T
	if rf, ok := ret.Get(0).(func(context.Context, domain.FindWith, interface{}) *T); ok {
		r0 = rf(ctx, key, val)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, domain.FindWith, interface{}) error); ok {
		r1 = rf(ctx, key, val)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIReadOneRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIReadOneRepository creates a new instance of IReadOneRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIReadOneRepository[T interface{}](t mockConstructorTestingTNewIReadOneRepository) *IReadOneRepository[T] {
	mock := &IReadOneRepository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
