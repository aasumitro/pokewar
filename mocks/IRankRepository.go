// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/aasumitro/pokewar/domain"
	mock "github.com/stretchr/testify/mock"
)

// IRankRepository is an autogenerated mock type for the IRankRepository type
type IRankRepository struct {
	mock.Mock
}

// All provides a mock function with given fields: ctx, args
func (_m *IRankRepository) All(ctx context.Context, args ...string) ([]*domain.Rank, error) {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*domain.Rank
	if rf, ok := ret.Get(0).(func(context.Context, ...string) []*domain.Rank); ok {
		r0 = rf(ctx, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Rank)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, ...string) error); ok {
		r1 = rf(ctx, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIRankRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIRankRepository creates a new instance of IRankRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIRankRepository(t mockConstructorTestingTNewIRankRepository) *IRankRepository {
	mock := &IRankRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
