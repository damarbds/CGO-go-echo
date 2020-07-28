// Code generated by mockery v1.0.0
package mocks

import (
	context "context"

	"github.com/models"

	mock "github.com/stretchr/testify/mock"
)

// repository is an autogenerated mock type for the repository type
type Repository struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: ctx, cursor, num
func (_m *Repository) GetByDate(ctx context.Context,from ,to string)(*models.ExChangeRate,error) {
	ret := _m.Called(ctx, from,to)

	var r0 *models.ExChangeRate
	if rf, ok := ret.Get(0).(func(context.Context, string,string) *models.ExChangeRate); ok {
		r0 = rf(ctx, from,to)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ExChangeRate)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string,string) error); ok {
		r1 = rf(ctx, from,to)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, a
func (_m *Repository) Insert(ctx context.Context, ExChangeRate *models.ExChangeRate) error {
	ret := _m.Called(ctx, ExChangeRate)

	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.ExChangeRate) error); ok {
		r1 = rf(ctx, ExChangeRate)
	} else {
		r1 = ret.Error(0)
	}

	return r1
}
