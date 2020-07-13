// Code generated by mockery v1.0.0
package mocks

import (
	context "context"

	"github.com/models"

	mock "github.com/stretchr/testify/mock"
)

// Usecase is an autogenerated mock type for the Usecase type
type Usecase struct {
	mock.Mock
}

// Fetch provides a mock function with given fields: ctx, cursor, num
func (_m *Usecase) GetAll(ctx context.Context, page, limit, offset int) (*models.MinimumBookingDtoWithPagination, error) {
	ret := _m.Called(ctx, page, limit, offset)

	var r0 *models.MinimumBookingDtoWithPagination
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int) *models.MinimumBookingDtoWithPagination); ok {
		r0 = rf(ctx, page, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MinimumBookingDtoWithPagination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int, int) error); ok {
		r1 = rf(ctx, page, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
