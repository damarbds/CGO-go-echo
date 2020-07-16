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

// Delete provides a mock function with given fields: ctx, id
func (_m *Usecase) Fetch(ctx context.Context, page *int, size *int) ([]*models.PromoDto, error) {
	ret := _m.Called(ctx,page,size)

	var r0 []*models.PromoDto
	if rf, ok := ret.Get(0).(func(context.Context,*int,*int) []*models.PromoDto); ok {
		r0 = rf(ctx,page,size)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.PromoDto)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context,*int,*int) error); ok {
		r1 = rf(ctx,page,size)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Fetch provides a mock function with given fields: ctx, cursor, num
func (_m *Usecase) List(ctx context.Context,page, limit, offset int, search string,token string)(*models.PromoWithPagination,error) {
	ret := _m.Called(ctx, page, limit, offset,search,token)

	var r0 *models.PromoWithPagination
	if rf, ok := ret.Get(0).(func(context.Context, int, int, int,string,string) *models.PromoWithPagination); ok {
		r0 = rf(ctx, page, limit, offset,search,token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PromoWithPagination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int, int, int,string,string) error); ok {
		r1 = rf(ctx, page, limit, offset,search,token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Usecase) GetDetail(ctx context.Context, id string,token string)(*models.PromoDto,error) {
	ret := _m.Called(ctx, id,token)

	var r0 *models.PromoDto
	if rf, ok := ret.Get(0).(func(context.Context, string,string) *models.PromoDto); ok {
		r0 = rf(ctx, id,token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PromoDto)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string,string) error); ok {
		r1 = rf(ctx, id,token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Usecase) GetByCode(ctx context.Context, code string,promoType int,merchantId string,token string) (*models.PromoDto, error) {
	ret := _m.Called(ctx, code,promoType,merchantId,token)

	var r0 *models.PromoDto
	if rf, ok := ret.Get(0).(func(context.Context, string,int,string,string) *models.PromoDto); ok {
		r0 = rf(ctx, code,promoType,merchantId,token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.PromoDto)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string,int,string,string) error); ok {
		r1 = rf(ctx, code,promoType,merchantId,token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTitle provides a mock function with given fields: ctx, title
func (_m *Usecase) Create(ctx context.Context, f models.NewCommandPromo, token string) (*models.NewCommandPromo, error) {
	ret := _m.Called(ctx, f, token)

	var r0 *models.NewCommandPromo
	if rf, ok := ret.Get(0).(func(context.Context, models.NewCommandPromo, string) *models.NewCommandPromo); ok {
		r0 = rf(ctx, f, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.NewCommandPromo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.NewCommandPromo, string) error); ok {
		r1 = rf(ctx, f, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: _a0, _a1
func (_m *Usecase) Update(ctx context.Context, f models.NewCommandPromo, token string) (*models.NewCommandPromo, error) {
	ret := _m.Called(ctx, f, token)

	var r0 *models.NewCommandPromo
	if rf, ok := ret.Get(0).(func(context.Context, models.NewCommandPromo, string) *models.NewCommandPromo); ok {
		r0 = rf(ctx, f, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.NewCommandPromo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.NewCommandPromo, string) error); ok {
		r1 = rf(ctx, f, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: _a0, _a1
func (_m *Usecase) Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error) {
	ret := _m.Called(ctx, id, token)

	var r0 *models.ResponseDelete
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *models.ResponseDelete); ok {
		r0 = rf(ctx, id, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResponseDelete)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, id, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
