package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	_adminUsecaseMock "github.com/auth/admin/mocks"
	"github.com/models"
	"github.com/service/harbors/mocks"
	"github.com/service/harbors/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	timeoutContext = time.Second * 30
)
var (
	imagePath        = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg"
	mockHarborsRepo  = new(mocks.Repository)
	mockAdminUsecase = new(_adminUsecaseMock.Usecase)
	harborsType = 1
	provinceName = "Jawa Barat"
	mockHarbors      = models.Harbors{
		Id:               "jlkjlkjlkjlkjlkj",
		CreatedBy:        "test",
		CreatedDate:      time.Now(),
		ModifiedBy:       nil,
		ModifiedDate:     nil,
		DeletedBy:        nil,
		DeletedDate:      nil,
		IsDeleted:        0,
		IsActive:         1,
		HarborsName:      "Harbors Test 1",
		HarborsLongitude: 1213,
		HarborsLatitude:  12313,
		HarborsImage:     "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		CityId:           1,
		HarborsType:      &harborsType,
	}
	mockHarborsJoin = []models.HarborsWCPC{
		models.HarborsWCPC{
			Id:               "dfgdgdgfdgfdgf",
			HarborsName:      "harbors Test 1",
			HarborsLongitude: 1231231231,
			HarborsLatitude:  43242342,
			HarborsImage:     "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
			CityId:           1,
			CityName:         "Bogor",
			ProvinceId:       1,
			ProvinceName:     &provinceName,
			CountryName:      "Indonesia",
		},
		models.HarborsWCPC{
			Id:               "dfgdgdgfdgfdgf",
			HarborsName:      "harbors Test 2",
			HarborsLongitude: 1231231231,
			HarborsLatitude:  43242342,
			HarborsImage:     "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
			CityId:           1,
			CityName:         "Bogor",
			ProvinceId:       1,
			ProvinceName:     &provinceName,
			CountryName:      "Indonesia",
		},
	}
)

func TestGetAll(t *testing.T) {

	var mockListHarbors []*models.Harbors
	mockListHarbors = append(mockListHarbors, &mockHarbors)
	page := 1
	limit := 1
	offset := 0
	count := 1
	t.Run("success", func(t *testing.T) {
		mockHarborsRepo.On("Fetch", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(mockListHarbors, nil).Once()
		mockHarborsRepo.On("GetCount", mock.Anything).Return(count, nil).Once()

		u := usecase.NewharborsUsecase(nil, mockHarborsRepo, timeoutContext)

		_, err := u.GetAll(context.TODO(), page, limit, offset)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListHarbors))

		mockHarborsRepo.AssertExpectations(t)
		mockHarborsRepo.AssertExpectations(t)
	})

}

func TestGetAllWithJoinCPC(t *testing.T) {

	var mockListHarbors []*models.HarborsWCPC
	mockListHarbors = append(mockListHarbors, &mockHarborsJoin[0])

	t.Run("success", func(t *testing.T) {
		mockHarborsRepo.On("GetAllWithJoinCPC", mock.Anything, mock.AnythingOfType("*int"), mock.AnythingOfType("*int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockListHarbors, nil).Once()

		u := usecase.NewharborsUsecase(nil, mockHarborsRepo, timeoutContext)
		page := 1
		limit := 1
		//offset := 0
		//count := 1
		list, err := u.GetAllWithJoinCPC(context.TODO(),&page,&limit,"","")

		assert.NoError(t, err)
		assert.Len(t, list, len(mockListHarbors))

		mockHarborsRepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockHarborsRepo.On("GetAllWithJoinCPC", mock.Anything, mock.AnythingOfType("*int"), mock.AnythingOfType("*int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("error unexpected")).Once()

		u := usecase.NewharborsUsecase(nil, mockHarborsRepo, timeoutContext)
		page := 1
		limit := 1
		//offset := 0
		//count := 1
		list, err := u.GetAllWithJoinCPC(context.TODO(),&page,&limit,"","")


		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockHarborsRepo.AssertExpectations(t)
	})

}

func TestGetById(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockHarborsRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(&mockHarbors, nil).Once()
		u := usecase.NewharborsUsecase(nil, mockHarborsRepo, timeoutContext)

		a, err := u.GetById(context.TODO(), mockHarbors.Id)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockHarborsRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockHarborsRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("Unexpected")).Once()

		u := usecase.NewharborsUsecase(nil, mockHarborsRepo, timeoutContext)

		a, err := u.GetById(context.TODO(), mockHarbors.Id)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockHarborsRepo.AssertExpectations(t)
	})

}

func TestCreate(t *testing.T) {

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	t.Run("success", func(t *testing.T) {
		tempMockHarbors := models.NewCommandHarbors{
			Id:               mockHarbors.Id,
			HarborsName:      mockHarbors.HarborsName,
			HarborsLongitude: mockHarbors.HarborsLongitude,
			HarborsLatitude:  mockHarbors.HarborsLatitude,
			HarborsImage:     mockHarbors.HarborsImage,
			CityId:           mockHarbors.CityId,
			HarborsType:      *mockHarbors.HarborsType,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockHarbors.Id = ""
		mockHarborsRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Harbors")).Return(&mockHarbors.Id, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()

		u := usecase.NewharborsUsecase(mockAdminUsecase, mockHarborsRepo, timeoutContext)

		_, err := u.Create(context.TODO(), &tempMockHarbors, token)

		assert.NoError(t, err)
		assert.Equal(t, mockHarbors.HarborsName, tempMockHarbors.HarborsName)
		mockHarborsRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockHarbors := models.NewCommandHarbors{
			Id:               mockHarbors.Id,
			HarborsName:      mockHarbors.HarborsName,
			HarborsLongitude: mockHarbors.HarborsLongitude,
			HarborsLatitude:  mockHarbors.HarborsLatitude,
			HarborsImage:     mockHarbors.HarborsImage,
			CityId:           mockHarbors.CityId,
			HarborsType:      *mockHarbors.HarborsType,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockHarbors.Id = ""
		mockHarborsRepo.On("Insert", mock.Anything, mock.AnythingOfType("*models.Harbors")).Return(&mockHarbors.Id, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("UnAuthorize")).Once()

		u := usecase.NewharborsUsecase(mockAdminUsecase, mockHarborsRepo, timeoutContext)

		_, err := u.Create(context.TODO(), &tempMockHarbors, token)

		assert.Error(t, err)
		//assert.Equal(t, mockHarbors.HarborsName, tempMockHarbors.HarborsName)
		//mockHarborsRepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	now := time.Now()

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	mockHarbors.ModifiedBy = &mockAdmin.Name
	mockHarbors.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {
		tempMockHarbors := models.NewCommandHarbors{
			Id:               mockHarbors.Id,
			HarborsName:      mockHarbors.HarborsName,
			HarborsLongitude: mockHarbors.HarborsLongitude,
			HarborsLatitude:  mockHarbors.HarborsLatitude,
			HarborsImage:     mockHarbors.HarborsImage,
			CityId:           mockHarbors.CityId,
			HarborsType:      *mockHarbors.HarborsType,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		//tempMockHarbors.Id = 0
		mockHarborsRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Harbors")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()

		u := usecase.NewharborsUsecase(mockAdminUsecase, mockHarborsRepo, timeoutContext)

		_, err := u.Update(context.TODO(), &tempMockHarbors, token)

		assert.NoError(t, err)
		assert.Equal(t, mockHarbors.HarborsName, tempMockHarbors.HarborsName)
		//mockHarborsRepo.AssertExpectations(t)
	})
	t.Run("success-without-image", func(t *testing.T) {
		tempMockHarbors := models.NewCommandHarbors{
			Id:               mockHarbors.Id,
			HarborsName:      mockHarbors.HarborsName,
			HarborsLongitude: mockHarbors.HarborsLongitude,
			HarborsLatitude:  mockHarbors.HarborsLatitude,
			HarborsImage:     mockHarbors.HarborsImage,
			CityId:           mockHarbors.CityId,
			HarborsType:      *mockHarbors.HarborsType,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockHarbors.Id = ""
		tempMockHarbors.HarborsImage = ""
		mockHarborsRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Harbors")).Return(nil).Once()
		mockHarborsRepo.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(&mockHarbors, nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()

		u := usecase.NewharborsUsecase(mockAdminUsecase, mockHarborsRepo, timeoutContext)

		_, err := u.Update(context.TODO(), &tempMockHarbors, token)

		assert.NoError(t, err)
		assert.Equal(t, mockHarbors.HarborsName, tempMockHarbors.HarborsName)
		//mockHarborsRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockHarbors := models.NewCommandHarbors{
			Id:               mockHarbors.Id,
			HarborsName:      mockHarbors.HarborsName,
			HarborsLongitude: mockHarbors.HarborsLongitude,
			HarborsLatitude:  mockHarbors.HarborsLatitude,
			HarborsImage:     mockHarbors.HarborsImage,
			CityId:           mockHarbors.CityId,
			HarborsType:      *mockHarbors.HarborsType,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockHarbors.Id = ""
		mockHarborsRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.Harbors")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("unAuthorize")).Once()

		u := usecase.NewharborsUsecase(mockAdminUsecase, mockHarborsRepo, timeoutContext)

		_, err := u.Update(context.TODO(), &tempMockHarbors, token)

		assert.Error(t, err)
		//assert.Equal(t, mockHarbors.HarborsName, tempMockHarbors.HarborsName)
		//mockHarborsRepo.AssertExpectations(t)
	})
}
func TestDelete(t *testing.T) {

	now := time.Now()

	mockAdmin := &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
	mockHarbors.ModifiedBy = &mockAdmin.Name
	mockHarbors.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockHarborsRepo.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin, nil).Once()

		u := usecase.NewharborsUsecase(mockAdminUsecase, mockHarborsRepo, timeoutContext)

		_, err := u.Delete(context.TODO(), mockHarbors.Id, token)

		assert.NoError(t, err)
		//assert.Equal(t, mockHarbors.HarborsName, tempMockHarbors.HarborsName)
		//mockHarborsRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockHarborsRepo.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("unAuthorize")).Once()

		u := usecase.NewharborsUsecase(mockAdminUsecase, mockHarborsRepo, timeoutContext)

		_, err := u.Delete(context.TODO(), mockHarbors.Id, token)

		assert.Error(t, err)
		//assert.Equal(t, mockHarbors.HarborsName, tempMockHarbors.HarborsName)
		//mockHarborsRepo.AssertExpectations(t)
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
