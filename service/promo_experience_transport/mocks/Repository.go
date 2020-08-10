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
func (_m *Repository) CountByPromoId(ctx context.Context,promoId string)(int,error) {
	ret := _m.Called(ctx,promoId)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, promoId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(int)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, promoId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTitle provides a mock function with given fields: ctx, title
func (_m *Repository) GetByExperienceTransportId(ctx context.Context,expId string,transportId string,promoId string)([]*models.PromoExperienceTransport,error) {
	ret := _m.Called(ctx, expId,transportId,promoId)

	var r0 []*models.PromoExperienceTransport
	if rf, ok := ret.Get(0).(func(context.Context, string,string,string) []*models.PromoExperienceTransport); ok {
		r0 = rf(ctx, expId,transportId,promoId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.PromoExperienceTransport)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string,string,string) error); ok {
		r1 = rf(ctx, expId,transportId,promoId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, ar
func (_m *Repository) Insert(ctx context.Context,pet models.PromoExperienceTransport)error {
	ret := _m.Called(ctx, pet)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, models.PromoExperienceTransport) error); ok {
		r0 = rf(ctx, pet)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: ctx, ar
func (_m *Repository) DeleteById(ctx context.Context,serviceId string,promoId string)error {
	ret := _m.Called(ctx, serviceId,promoId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string,string) error); ok {
		r0 = rf(ctx, serviceId,promoId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
