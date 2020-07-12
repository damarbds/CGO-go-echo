package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	"github.com/service/minimum_booking/mocks"
	"github.com/service/minimum_booking/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	timeoutContext = time.Second * 30
)

func TestList(t *testing.T) {
	mockMinimumBookingRepo := new(mocks.Repository)
	mockMinimumBooking := models.MinimumBooking{
		Id:                   "dsadasdasdasdasdsa",
		CreatedBy:            "test",
		CreatedDate:          time.Now(),
		ModifiedBy:           nil,
		ModifiedDate:         nil,
		DeletedBy:            nil,
		DeletedDate:          nil,
		IsDeleted:            0,
		IsActive:             1,
		MinimumBookingDesc:   "Test MinimumBooking 2",
		MinimumBookingAmount: 0,
	}
	var mockListMinimumBooking []*models.MinimumBooking
	mockListMinimumBooking = append(mockListMinimumBooking, &mockMinimumBooking)
	page := 1
	limit := 1
	offset := 0
	count := 1
	t.Run("success", func(t *testing.T) {
		mockMinimumBookingRepo.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(mockListMinimumBooking, nil).Once()
		mockMinimumBookingRepo.On("GetCount", mock.Anything).Return(count, nil).Once()

		u := usecase.NewminimumBookingUsecase(mockMinimumBookingRepo, timeoutContext)

		_, err := u.GetAll(context.TODO(), page, limit, offset)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListMinimumBooking))

		mockMinimumBookingRepo.AssertExpectations(t)
		mockMinimumBookingRepo.AssertExpectations(t)
	})

}

//
//func TestList(t *testing.T) {
//	mockMinimumBookingRepo := new(mocks.Repository)
//	mockMinimumBooking := 	models.MinimumBooking{
//		Id:           1,
//		CreatedBy:    "test",
//		CreatedDate:  time.Now(),
//		ModifiedBy:   nil,
//		ModifiedDate: nil,
//		DeletedBy:    nil,
//		DeletedDate:  nil,
//		IsDeleted:    0,
//		IsActive:     1,
//		MinimumBookingName:  "Test MinimumBooking 2",
//		MinimumBookingIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/MinimumBooking/8941695193938718058.jpg",
//	}
//	var mockListMinimumBooking []*models.MinimumBooking
//	mockListMinimumBooking = append(mockListMinimumBooking, &mockMinimumBooking)
//
//	t.Run("success", func(t *testing.T) {
//		mockMinimumBookingRepo.On("List", mock.Anything).Return(mockListMinimumBooking, nil).Once()
//
//		u := usecase.NewMinimumBookingUsecase(nil, mockMinimumBookingRepo, timeoutContext)
//
//		list, err := u.List(context.TODO())
//
//		assert.NoError(t, err)
//		assert.Len(t, list, len(mockListMinimumBooking))
//
//		mockMinimumBookingRepo.AssertExpectations(t)
//	})
//
//	t.Run("error-failed", func(t *testing.T) {
//		mockMinimumBookingRepo.On("List", mock.Anything).Return(nil,errors.New("Unexpexted Error")).Once()
//
//		u := usecase.NewMinimumBookingUsecase(nil, mockMinimumBookingRepo, timeoutContext)
//
//		list, err := u.List(context.TODO())
//
//		assert.Error(t, err)
//		assert.Len(t, list, 0)
//		mockMinimumBookingRepo.AssertExpectations(t)
//	})
//
//}
////
////func TestGetById(t *testing.T) {
////	mockMinimumBookingRepo := new(mocks.Repository)
////	mockMinimumBooking := 	models.MinimumBooking{
////		Id:           1,
////		CreatedBy:    "test",
////		CreatedDate:  time.Now(),
////		ModifiedBy:   nil,
////		ModifiedDate: nil,
////		DeletedBy:    nil,
////		DeletedDate:  nil,
////		IsDeleted:    0,
////		IsActive:     1,
////		MinimumBookingName:  "Test MinimumBooking 2",
////		MinimumBookingIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/MinimumBooking/8941695193938718058.jpg",
////	}
////
////	t.Run("success", func(t *testing.T) {
////		mockMinimumBookingRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(&mockMinimumBooking, nil).Once()
////		u := usecase.NewMinimumBookingUsecase(nil, mockMinimumBookingRepo, timeoutContext)
////
////		a, err := u.GetById(context.TODO(), mockMinimumBooking.Id)
////
////		assert.NoError(t, err)
////		assert.NotNil(t, a)
////
////		mockMinimumBookingRepo.AssertExpectations(t)
////	})
////	t.Run("error-failed", func(t *testing.T) {
////		mockMinimumBookingRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(nil, errors.New("Unexpected")).Once()
////
////		u := usecase.NewMinimumBookingUsecase(nil, mockMinimumBookingRepo, timeoutContext)
////
////		a, err := u.GetById(context.TODO(), mockMinimumBooking.Id)
////
////		assert.Error(t, err)
////		assert.Nil(t, a)
////
////		mockMinimumBookingRepo.AssertExpectations(t)
////	})
////
////}
////
////func TestCreate(t *testing.T) {
////	mockMinimumBookingRepo := new(mocks.Repository)
////	mockAdminUsecase := new(_adminUsecaseMock.Usecase)
////	mockMinimumBooking := 	models.MinimumBooking{
////		Id:           0,
////		CreatedBy:    "adminCGO",
////		CreatedDate:  time.Now(),
////		ModifiedBy:   nil,
////		ModifiedDate: nil,
////		DeletedBy:    nil,
////		DeletedDate:  nil,
////		IsDeleted:    0,
////		IsActive:     1,
////		MinimumBookingName:  "Test MinimumBooking 2",
////		MinimumBookingIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/MinimumBooking/8941695193938718058.jpg",
////	}
////	mockAdmin := &models.AdminDto{
////		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
////		Name:  "adminCGO",
////		Email: "admin1234@gmail.com",
////	}
////	t.Run("success", func(t *testing.T) {
////		tempMockMinimumBooking := models.NewCommandMinimumBooking{
////			Id:          0,
////			MinimumBookingName: mockMinimumBooking.MinimumBookingName,
////			MinimumBookingIcon: mockMinimumBooking.MinimumBookingIcon,
////		}
////		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
////		tempMockMinimumBooking.Id = 0
////		mockMinimumBookingRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.MinimumBooking")).Return(&mockMinimumBooking.Id,nil).Once()
////		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()
////
////		u := usecase.NewMinimumBookingUsecase(mockAdminUsecase, mockMinimumBookingRepo, timeoutContext)
////
////		_,err := u.Create(context.TODO(), &tempMockMinimumBooking,token)
////
////		assert.NoError(t, err)
////		assert.Equal(t, mockMinimumBooking.MinimumBookingName, tempMockMinimumBooking.MinimumBookingName)
////		mockMinimumBookingRepo.AssertExpectations(t)
////	})
////	t.Run("error-unauthorize", func(t *testing.T) {
////		tempMockMinimumBooking := models.NewCommandMinimumBooking{
////			Id:          0,
////			MinimumBookingName: mockMinimumBooking.MinimumBookingName,
////			MinimumBookingIcon: mockMinimumBooking.MinimumBookingIcon,
////		}
////		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
////		tempMockMinimumBooking.Id = 0
////		mockMinimumBookingRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.MinimumBooking")).Return(&mockMinimumBooking.Id,nil).Once()
////		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("UnAuthorize")).Once()
////
////		u := usecase.NewMinimumBookingUsecase(mockAdminUsecase, mockMinimumBookingRepo, timeoutContext)
////
////		_,err := u.Create(context.TODO(), &tempMockMinimumBooking,token)
////
////		assert.Error(t, err)
////		//assert.Equal(t, mockMinimumBooking.MinimumBookingName, tempMockMinimumBooking.MinimumBookingName)
////		//mockMinimumBookingRepo.AssertExpectations(t)
////	})
////
////}
////
////func TestUpdate(t *testing.T) {
////	mockMinimumBookingRepo := new(mocks.Repository)
////	mockAdminUsecase := new(_adminUsecaseMock.Usecase)
////	now := time.Now()
////
////	mockAdmin := &models.AdminDto{
////		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
////		Name:  "adminCGO",
////		Email: "admin1234@gmail.com",
////	}
////	mockMinimumBooking := 	models.MinimumBooking{
////		Id:           0,
////		CreatedBy:    mockAdmin.Name,
////		CreatedDate:  time.Now(),
////		ModifiedBy:   &mockAdmin.Name,
////		ModifiedDate: &now,
////		DeletedBy:    nil,
////		DeletedDate:  nil,
////		IsDeleted:    0,
////		IsActive:     1,
////		MinimumBookingName:  "Test MinimumBooking 2",
////		MinimumBookingIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/MinimumBooking/8941695193938718058.jpg",
////	}
////	t.Run("success", func(t *testing.T) {
////		tempMockMinimumBooking := models.NewCommandMinimumBooking{
////			Id:          0,
////			MinimumBookingName: mockMinimumBooking.MinimumBookingName,
////			MinimumBookingIcon: mockMinimumBooking.MinimumBookingIcon,
////		}
////		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
////		tempMockMinimumBooking.Id = 0
////		mockMinimumBookingRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.MinimumBooking")).Return(nil).Once()
////		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()
////
////		u := usecase.NewMinimumBookingUsecase(mockAdminUsecase, mockMinimumBookingRepo, timeoutContext)
////
////		_,err := u.Update(context.TODO(), &tempMockMinimumBooking,token)
////
////		assert.NoError(t, err)
////		assert.Equal(t, mockMinimumBooking.MinimumBookingName, tempMockMinimumBooking.MinimumBookingName)
////		mockMinimumBookingRepo.AssertExpectations(t)
////	})
////	t.Run("success-without-image", func(t *testing.T) {
////		tempMockMinimumBooking := models.NewCommandMinimumBooking{
////			Id:          0,
////			MinimumBookingName: mockMinimumBooking.MinimumBookingName,
////			MinimumBookingIcon: mockMinimumBooking.MinimumBookingIcon,
////		}
////		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
////		tempMockMinimumBooking.Id = 0
////		tempMockMinimumBooking.MinimumBookingIcon = ""
////		mockMinimumBookingRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.MinimumBooking")).Return(nil).Once()
////		mockMinimumBookingRepo.On("GetById", mock.Anything, mock.AnythingOfType("int")).Return(&mockMinimumBooking,nil).Once()
////		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()
////
////		u := usecase.NewMinimumBookingUsecase(mockAdminUsecase, mockMinimumBookingRepo, timeoutContext)
////
////		_,err := u.Update(context.TODO(), &tempMockMinimumBooking,token)
////
////		assert.NoError(t, err)
////		assert.Equal(t, mockMinimumBooking.MinimumBookingName, tempMockMinimumBooking.MinimumBookingName)
////		mockMinimumBookingRepo.AssertExpectations(t)
////	})
////	t.Run("error-unauthorize", func(t *testing.T) {
////		tempMockMinimumBooking := models.NewCommandMinimumBooking{
////			Id:          0,
////			MinimumBookingName: mockMinimumBooking.MinimumBookingName,
////			MinimumBookingIcon: mockMinimumBooking.MinimumBookingIcon,
////		}
////		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
////		tempMockMinimumBooking.Id = 0
////		mockMinimumBookingRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.MinimumBooking")).Return(nil).Once()
////		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()
////
////		u := usecase.NewMinimumBookingUsecase(mockAdminUsecase, mockMinimumBookingRepo, timeoutContext)
////
////		_,err := u.Update(context.TODO(), &tempMockMinimumBooking,token)
////
////		assert.Error(t, err)
////		//assert.Equal(t, mockMinimumBooking.MinimumBookingName, tempMockMinimumBooking.MinimumBookingName)
////		//mockMinimumBookingRepo.AssertExpectations(t)
////	})
////}
////func TestDelete(t *testing.T) {
////	mockMinimumBookingRepo := new(mocks.Repository)
////	mockAdminUsecase := new(_adminUsecaseMock.Usecase)
////	now := time.Now()
////
////	mockAdmin := &models.AdminDto{
////		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
////		Name:  "adminCGO",
////		Email: "admin1234@gmail.com",
////	}
////	mockMinimumBooking := 	models.MinimumBooking{
////		Id:           1,
////		CreatedBy:    mockAdmin.Name,
////		CreatedDate:  time.Now(),
////		ModifiedBy:   &mockAdmin.Name,
////		ModifiedDate: &now,
////		DeletedBy:    nil,
////		DeletedDate:  nil,
////		IsDeleted:    0,
////		IsActive:     1,
////		MinimumBookingName:  "Test MinimumBooking 2",
////		MinimumBookingIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/MinimumBooking/8941695193938718058.jpg",
////	}
////	t.Run("success", func(t *testing.T) {
////
////		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
////
////		mockMinimumBookingRepo.On("Delete", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
////		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()
////
////		u := usecase.NewMinimumBookingUsecase(mockAdminUsecase, mockMinimumBookingRepo, timeoutContext)
////
////		_,err := u.Delete(context.TODO(), mockMinimumBooking.Id,token)
////
////		assert.NoError(t, err)
////		//assert.Equal(t, mockMinimumBooking.MinimumBookingName, tempMockMinimumBooking.MinimumBookingName)
////		mockMinimumBookingRepo.AssertExpectations(t)
////	})
////	t.Run("error-unauthorize", func(t *testing.T) {
////
////		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
////
////		mockMinimumBookingRepo.On("Delete", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
////		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()
////
////		u := usecase.NewMinimumBookingUsecase(mockAdminUsecase, mockMinimumBookingRepo, timeoutContext)
////
////		_,err := u.Delete(context.TODO(), mockMinimumBooking.Id,token)
////
////		assert.Error(t, err)
////		//assert.Equal(t, mockMinimumBooking.MinimumBookingName, tempMockMinimumBooking.MinimumBookingName)
////		//mockMinimumBookingRepo.AssertExpectations(t)
////	})
////}

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
