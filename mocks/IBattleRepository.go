// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/aasumitro/pokewar/domain"
	mock "github.com/stretchr/testify/mock"
)

// IBattleRepository is an autogenerated mock type for the IBattleRepository type
type IBattleRepository struct {
	mock.Mock
}

// All provides a mock function with given fields: ctx, args
func (_m *IBattleRepository) All(ctx context.Context, args ...string) ([]*domain.Battle, error) {
	_va := make([]interface{}, len(args))
	for _i := range args {
		_va[_i] = args[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*domain.Battle
	if rf, ok := ret.Get(0).(func(context.Context, ...string) []*domain.Battle); ok {
		r0 = rf(ctx, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Battle)
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

// Count provides a mock function with given fields: ctx
func (_m *IBattleRepository) Count(ctx context.Context) int {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	return r0
}

// Create provides a mock function with given fields: ctx, param
func (_m *IBattleRepository) Create(ctx context.Context, param *domain.Battle) error {
	ret := _m.Called(ctx, param)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Battle) error); ok {
		r0 = rf(ctx, param)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdatePlayer provides a mock function with given fields: ctx, id
func (_m *IBattleRepository) UpdatePlayer(ctx context.Context, id int) (int64, error) {
	ret := _m.Called(ctx, id)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, int) int64); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIBattleRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIBattleRepository creates a new instance of IBattleRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIBattleRepository(t mockConstructorTestingTNewIBattleRepository) *IBattleRepository {
	mock := &IBattleRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
