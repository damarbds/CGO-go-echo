package usecase_test

import (
	"context"
	"errors"
	_adminUsecaseMock "github.com/auth/admin/mocks"
	"github.com/models"
	"github.com/service/exclude/mocks"
	"github.com/service/exclude/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)
var(
	timeoutContext = time.Second*30
)
func TestGetAll(t *testing.T) {
	mockIncludeRepo := new(mocks.Repository)
	mockInclude := 	models.Exclude{
		Id:           1,
		CreatedBy:    "test",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		ExcludeName:  "Test Include 2",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}
	var mockListInclude []*models.Exclude
	mockListInclude = append(mockListInclude, &mockInclude)
	page := 1
	limit := 1
	offset := 0
	count := 1
	t.Run("success", func(t *testing.T) {
		mockIncludeRepo.On("Fetch", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(mockListInclude, nil).Once()
		mockIncludeRepo.On("GetCount", mock.Anything).Return(count, nil).Once()

		u := usecase.NewExcludeUsecase(nil, mockIncludeRepo, timeoutContext)

		_, err := u.GetAll(context.TODO(),page,limit,offset)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListInclude))

		mockIncludeRepo.AssertExpectations(t)
		mockIncludeRepo.AssertExpectations(t)
	})

}

func TestList(t *testing.T) {
	mockIncludeRepo := new(mocks.Repository)
	mockInclude := 	models.Exclude{
		Id:           1,
		CreatedBy:    "test",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		ExcludeName:  "Test Include 2",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}
	var mockListInclude []*models.Exclude
	mockListInclude = append(mockListInclude, &mockInclude)

	t.Run("success", func(t *testing.T) {
		mockIncludeRepo.On("List", mock.Anything).Return(mockListInclude, nil).Once()

		u := usecase.NewExcludeUsecase(nil, mockIncludeRepo, timeoutContext)

		list, err := u.List(context.TODO())

		assert.NoError(t, err)
		assert.Len(t, list, len(mockListInclude))

		mockIncludeRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockIncludeRepo.On("List", mock.Anything).Return(nil,errors.New("Unexpexted Error")).Once()

		u := usecase.NewExcludeUsecase(nil, mockIncludeRepo, timeoutContext)

		list, err := u.List(context.TODO())

		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockIncludeRepo.AssertExpectations(t)
	})

}

func TestGetById(t *testing.T) {
	mockIncludeRepo := new(mocks.Repository)
	mockInclude := 	models.Exclude{
		Id:           1,
		CreatedBy:    "test",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		ExcludeName:  "Test Include 2",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}

	t.Run("success", func(t *testing.T) {
		mockIncludeRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(&mockInclude, nil).Once()
		u := usecase.NewExcludeUsecase(nil, mockIncludeRepo, timeoutContext)

		a, err := u.GetById(context.TODO(), mockInclude.Id)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockIncludeRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockIncludeRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(nil, errors.New("Unexpected")).Once()

		u := usecase.NewExcludeUsecase(nil, mockIncludeRepo, timeoutContext)

		a, err := u.GetById(context.TODO(), mockInclude.Id)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockIncludeRepo.AssertExpectations(t)
	})

}

func TestCreate(t *testing.T) {
	mockIncludeRepo := new(mocks.Repository)
	mockAdminUsecase := new(_adminUsecaseMock.Usecase)
	mockInclude := 	models.Exclude{
		Id:           0,
		CreatedBy:    "adminCGO",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		ExcludeName:  "Test Include 2",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}
	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	t.Run("success", func(t *testing.T) {
		tempMockInclude := models.NewCommandExclude{
			Id:          0,
			ExcludeName: mockInclude.ExcludeName,
			ExcludeIcon: mockInclude.ExcludeIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockInclude.Id = 0
		mockIncludeRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Exclude")).Return(&mockInclude.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewExcludeUsecase(mockAdminUsecase, mockIncludeRepo, timeoutContext)

		_,err := u.Create(context.TODO(), &tempMockInclude,token)

		assert.NoError(t, err)
		assert.Equal(t, mockInclude.ExcludeName, tempMockInclude.ExcludeName)
		mockIncludeRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockInclude := models.NewCommandExclude{
			Id:          0,
			ExcludeName: mockInclude.ExcludeName,
			ExcludeIcon: mockInclude.ExcludeIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockInclude.Id = 0
		mockIncludeRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Exclude")).Return(&mockInclude.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("UnAuthorize")).Once()

		u := usecase.NewExcludeUsecase(mockAdminUsecase, mockIncludeRepo, timeoutContext)

		_,err := u.Create(context.TODO(), &tempMockInclude,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	mockIncludeRepo := new(mocks.Repository)
	mockAdminUsecase := new(_adminUsecaseMock.Usecase)
	now := time.Now()

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	mockInclude := 	models.Exclude{
		Id:           0,
		CreatedBy:    mockAdmin.Name,
		CreatedDate:  time.Now(),
		ModifiedBy:   &mockAdmin.Name,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		ExcludeName:  "Test Include 2",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}
	t.Run("success", func(t *testing.T) {
		tempMockInclude := models.NewCommandExclude{
			Id:          0,
			ExcludeName: mockInclude.ExcludeName,
			ExcludeIcon: mockInclude.ExcludeIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockInclude.Id = 0
		mockIncludeRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Exclude")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewExcludeUsecase(mockAdminUsecase, mockIncludeRepo, timeoutContext)

		_,err := u.Update(context.TODO(), &tempMockInclude,token)

		assert.NoError(t, err)
		assert.Equal(t, mockInclude.ExcludeName, tempMockInclude.ExcludeName)
		mockIncludeRepo.AssertExpectations(t)
	})
	t.Run("success-without-image", func(t *testing.T) {
		tempMockInclude := models.NewCommandExclude{
			Id:          0,
			ExcludeName: mockInclude.ExcludeName,
			ExcludeIcon: mockInclude.ExcludeIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockInclude.Id = 0
		tempMockInclude.ExcludeIcon = ""
		mockIncludeRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Exclude")).Return(nil).Once()
		mockIncludeRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(&mockInclude,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewExcludeUsecase(mockAdminUsecase, mockIncludeRepo, timeoutContext)

		_,err := u.Update(context.TODO(), &tempMockInclude,token)

		assert.NoError(t, err)
		assert.Equal(t, mockInclude.ExcludeName, tempMockInclude.ExcludeName)
		mockIncludeRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockInclude := models.NewCommandExclude{
			Id:          0,
			ExcludeName: mockInclude.ExcludeName,
			ExcludeIcon: mockInclude.ExcludeIcon,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockInclude.Id = 0
		mockIncludeRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Exclude")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()

		u := usecase.NewExcludeUsecase(mockAdminUsecase, mockIncludeRepo, timeoutContext)

		_,err := u.Update(context.TODO(), &tempMockInclude,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})
}
func TestDelete(t *testing.T) {
	mockIncludeRepo := new(mocks.Repository)
	mockAdminUsecase := new(_adminUsecaseMock.Usecase)
	now := time.Now()

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	mockInclude := 	models.Exclude{
		Id:           1,
		CreatedBy:    mockAdmin.Name,
		CreatedDate:  time.Now(),
		ModifiedBy:   &mockAdmin.Name,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		ExcludeName:  "Test Include 2",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}
	t.Run("success", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockIncludeRepo.On("Delete", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewExcludeUsecase(mockAdminUsecase, mockIncludeRepo, timeoutContext)

		_,err := u.Delete(context.TODO(), mockInclude.Id,token)

		assert.NoError(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		mockIncludeRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockIncludeRepo.On("Delete", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()

		u := usecase.NewExcludeUsecase(mockAdminUsecase, mockIncludeRepo, timeoutContext)

		_,err := u.Delete(context.TODO(), mockInclude.Id,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})
}

//func TestDelete(t *testing.T) {
//	mockArticleRepo := new(mocks.Repository)
//	mockArticle := models.Article{
//		Title:   "Hello",
//		Content: "Content",
//	}
//
//	t.Run("success", func(t *testing.T) {
//		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(&mockArticle, nil).Once()
//
//		mockArticleRepo.On("Delete", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()
//
//		mockAuthorrepo := new(_authorMock.Repository)
//		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
//
//		err := u.Delete(context.TODO(), mockArticle.ID)
//
//		assert.NoError(t, err)
//		mockArticleRepo.AssertExpectations(t)
//		mockAuthorrepo.AssertExpectations(t)
//	})
//	t.Run("article-is-not-exist", func(t *testing.T) {
//		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, nil).Once()
//
//		mockAuthorrepo := new(_authorMock.Repository)
//		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
//
//		err := u.Delete(context.TODO(), mockArticle.ID)
//
//		assert.Error(t, err)
//		mockArticleRepo.AssertExpectations(t)
//		mockAuthorrepo.AssertExpectations(t)
//	})
//	t.Run("error-happens-in-db", func(t *testing.T) {
//		mockArticleRepo.On("GetByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil, errors.New("Unexpected Error")).Once()
//
//		mockAuthorrepo := new(_authorMock.Repository)
//		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
//
//		err := u.Delete(context.TODO(), mockArticle.ID)
//
//		assert.Error(t, err)
//		mockArticleRepo.AssertExpectations(t)
//		mockAuthorrepo.AssertExpectations(t)
//	})
//
//}

