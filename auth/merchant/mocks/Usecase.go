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
func (_m *Usecase) Update(ctx context.Context, ar *models.NewCommandMerchant, isAdmin bool, token string) error {
	ret := _m.Called(ctx, ar, isAdmin, token)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.NewCommandMerchant, bool, string) error); ok {
		r0 = rf(ctx, ar, isAdmin, token)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Fetch provides a mock function with given fields: ctx, cursor, num
func (_m *Usecase) Create(ctx context.Context, ar *models.NewCommandMerchant, token string) error {
	ret := _m.Called(ctx, ar, token)

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.NewCommandMerchant, string) error); ok {
		r1 = rf(ctx, ar, token)
	} else {
		r1 = ret.Error(1)
	}

	return r1
}

// GetByID provides a mock function with given fields: ctx, id
func (_m *Usecase) Login(ctx context.Context, ar *models.Login) (*models.GetToken, error) {
	ret := _m.Called(ctx, ar)

	var r0 *models.GetToken
	if rf, ok := ret.Get(0).(func(context.Context, *models.Login) *models.GetToken); ok {
		r0 = rf(ctx, ar)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.Login) error); ok {
		r1 = rf(ctx, ar)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Usecase) AutoLoginByCMSAdmin(ctx context.Context,merchantId string,token string)(*models.GetToken,error) {
	ret := _m.Called(ctx, merchantId,token)

	var r0 *models.GetToken
	if rf, ok := ret.Get(0).(func(context.Context, string,string) *models.GetToken); ok {
		r0 = rf(ctx, merchantId,token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.GetToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string,string) error); ok {
		r1 = rf(ctx, merchantId,token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTitle provides a mock function with given fields: ctx, title
func (_m *Usecase) ValidateTokenMerchant(ctx context.Context, token string) (*models.MerchantInfoDto, error) {
	ret := _m.Called(ctx, token)

	var r0 *models.MerchantInfoDto
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.MerchantInfoDto); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MerchantInfoDto)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: _a0, _a1
func (_m *Usecase) GetMerchantInfo(ctx context.Context, token string) (*models.MerchantInfoDto, error) {
	ret := _m.Called(ctx, token)

	var r0 *models.MerchantInfoDto
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.MerchantInfoDto); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MerchantInfoDto)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Usecase) Count(ctx context.Context) (*models.Count, error) {
	ret := _m.Called(ctx)

	var r0 *models.Count
	if rf, ok := ret.Get(0).(func(context.Context) *models.Count); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Count)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Usecase) GetMerchantTransport(ctx context.Context) ([]*models.MerchantTransport, error) {
	ret := _m.Called(ctx)

	var r0 []*models.MerchantTransport
	if rf, ok := ret.Get(0).(func(context.Context) []*models.MerchantTransport); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.MerchantTransport)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Usecase) GetMerchantExperience(ctx context.Context) ([]*models.MerchantExperience, error) {
	ret := _m.Called(ctx)

	var r0 []*models.MerchantExperience
	if rf, ok := ret.Get(0).(func(context.Context) []*models.MerchantExperience); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.MerchantExperience)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Fetch provides a mock function with given fields: ctx, cursor, num
func (_m *Usecase) List(ctx context.Context, page, limit, offset int, token string,search string) (*models.MerchantWithPagination, error) {
	ret := _m.Called(ctx, page,limit,offset,token,search)

	var r0 *models.MerchantWithPagination
	if rf, ok := ret.Get(0).(func(context.Context, int,int,int,string,string) *models.MerchantWithPagination); ok {
		r0 = rf(ctx, page,limit,offset,token,search)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MerchantWithPagination)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int,int,int,string,string) error); ok {
		r1 = rf(ctx, page,limit,offset,token,search)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByTitle provides a mock function with given fields: ctx, title
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
func (_m *Usecase) SendingEmailMerchant(ctx context.Context,ar *models.NewCommandMerchantRegistrationEmail)(*models.ResponseDelete,error) {
	ret := _m.Called(ctx, ar)

	var r0 *models.ResponseDelete
	if rf, ok := ret.Get(0).(func(context.Context, *models.NewCommandMerchantRegistrationEmail) *models.ResponseDelete); ok {
		r0 = rf(ctx, ar)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResponseDelete)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.NewCommandMerchantRegistrationEmail) error); ok {
		r1 = rf(ctx, ar)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Usecase) 	SendingEmailContactUs(ctx context.Context,ar *models.NewCommandContactUs)(*models.ResponseDelete,error) {
	ret := _m.Called(ctx, ar)

	var r0 *models.ResponseDelete
	if rf, ok := ret.Get(0).(func(context.Context, *models.NewCommandContactUs) *models.ResponseDelete); ok {
		r0 = rf(ctx, ar)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ResponseDelete)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.NewCommandContactUs) error); ok {
		r1 = rf(ctx, ar)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: _a0, _a1
func (_m *Usecase) ServiceCount(ctx context.Context, token string) (*models.ServiceCount, error) {
	ret := _m.Called(ctx, token)

	var r0 *models.ServiceCount
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.ServiceCount); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.ServiceCount)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Usecase) 	GetDetailMerchantById(ctx context.Context, id string,token string)(*models.MerchantDto,error) {
	ret := _m.Called(ctx, token, id,token)

	var r0 *models.MerchantDto
	if rf, ok := ret.Get(0).(func(context.Context, string, string) *models.MerchantDto); ok {
		r0 = rf(ctx, id, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.MerchantDto)
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
