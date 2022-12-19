// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	domain "github.com/aasumitro/pokewar/domain"
	mock "github.com/stretchr/testify/mock"
)

// IPokeapiRESTRepository is an autogenerated mock type for the IPokeapiRESTRepository type
type IPokeapiRESTRepository struct {
	mock.Mock
}

// Pokemon provides a mock function with given fields: offset
func (_m *IPokeapiRESTRepository) Pokemon(offset int) ([]*domain.Monster, error) {
	ret := _m.Called(offset)

	var r0 []*domain.Monster
	if rf, ok := ret.Get(0).(func(int) []*domain.Monster); ok {
		r0 = rf(offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.Monster)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewIPokeapiRESTRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewIPokeapiRESTRepository creates a new instance of IPokeapiRESTRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewIPokeapiRESTRepository(t mockConstructorTestingTNewIPokeapiRESTRepository) *IPokeapiRESTRepository {
	mock := &IPokeapiRESTRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}