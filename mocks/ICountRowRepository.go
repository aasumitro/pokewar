// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ICountRowRepository is an autogenerated mock type for the ICountRowRepository type
type ICountRowRepository struct {
	mock.Mock
}

// Count provides a mock function with given fields: ctx
func (_m *ICountRowRepository) Count(ctx context.Context) int {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

type mockConstructorTestingTNewICountRowRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewICountRowRepository creates a new instance of ICountRowRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewICountRowRepository(t mockConstructorTestingTNewICountRowRepository) *ICountRowRepository {
	mock := &ICountRowRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
