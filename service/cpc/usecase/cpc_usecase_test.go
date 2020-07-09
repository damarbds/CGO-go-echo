package usecase_test

import (
	"context"
	"encoding/json"
	"errors"
	_adminUsecaseMock "github.com/auth/admin/mocks"
	"github.com/models"
	"github.com/service/cpc/mocks"
	"github.com/service/cpc/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)
var(
	timeoutContext = time.Second*30
	mockCPCRepo = new(mocks.Repository)
	mockAdminUsecase = new(_adminUsecaseMock.Usecase)
	imageTestPath = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg"
	coverPhoto = []models.CoverPhotosObj{
		models.CoverPhotosObj{
			Original:  "",
			Thumbnail: imageTestPath,
		},
	}
	cityPhotos ,_= json.Marshal(coverPhoto)
	cityPhotosJson = string(cityPhotos)
	code = 0
	mockCity = models.City{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			CityName:  "Bogor",
			CityDesc:  "Bogor adalah kota hujan",
			CityPhotos:&cityPhotosJson,
			ProvinceId:1,
		}
	mockProvince = models.Province{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ProvinceName:  "Jawa Barat",
			ProvinceNameTransportation:&imageTestPath,
			CountryId:1,
		}
	mockCountry = models.Country{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			CountryName:  "Indonesia",
			Iso:          &imageTestPath,
			Name:         &imageTestPath,
			NiceName:     &imageTestPath,
			Iso3:         &imageTestPath,
			NumCode:      &code,
			PhoneCode:    &code,
		}
	mockAdmin = &models.AdminDto{
		Id:    "972fe09f-93e9-4798-b642-14e0ac77c6be",
		Name:  "adminCGO",
		Email: "admin1234@gmail.com",
	}
)

func TestGetAllCity(t *testing.T) {
	var mockListCity []*models.City
	mockListCity = append(mockListCity, &mockCity)
	page := 1
	limit := 1
	offset := 0
	count := 1
	t.Run("success", func(t *testing.T) {
		mockCPCRepo.On("FetchCity", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(mockListCity, nil).Once()
		mockCPCRepo.On("GetCountCity", mock.Anything).Return(count, nil).Once()

		u := usecase.NewCPCUsecase(nil, mockCPCRepo, timeoutContext)

		_, err := u.GetAllCity(context.TODO(),page,limit,offset)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListInclude))

		mockCPCRepo.AssertExpectations(t)
		mockCPCRepo.AssertExpectations(t)
	})

}

func TestGetCityById(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockCPCRepo.On("GetCityByID", mock.Anything, mock.AnythingOfType("int")).Return(&mockCity, nil).Once()
		u := usecase.NewCPCUsecase(nil, mockCPCRepo, timeoutContext)

		a, err := u.GetCityById(context.TODO(), mockCity.Id)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockCPCRepo.On("GetCityByID", mock.Anything, mock.AnythingOfType("int")).Return(nil, errors.New("Unexpected")).Once()

		u := usecase.NewCPCUsecase(nil, mockCPCRepo, timeoutContext)

		a, err := u.GetCityById(context.TODO(), mockCity.Id)

		assert.Error(t, err)
		assert.Nil(t, a)

		mockCPCRepo.AssertExpectations(t)
	})

}

func TestCreateCity(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		tempMockCity := models.NewCommandCity{
			Id:         1,
			CityName:   mockCity.CityName,
			CityDesc:   mockCity.CityDesc,
			CityPhotos: coverPhoto,
			ProvinceId: 1,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCity.Id = 0
		mockCPCRepo.On("InsertCity", mock.Anything, mock.AnythingOfType("*models.City")).Return(&mockCity.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.CreateCity(context.TODO(), &tempMockCity,token)

		assert.NoError(t, err)
		assert.Equal(t, mockCity.CityName, tempMockCity.CityName)
		mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockCity := models.NewCommandCity{
			Id:         1,
			CityName:   mockCity.CityName,
			CityDesc:   mockCity.CityDesc,
			CityPhotos: coverPhoto,
			ProvinceId: 0,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCity.Id = 0
		mockCPCRepo.On("InsertCity", mock.Anything, mock.AnythingOfType("*models.City")).Return(&mockCity.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("UnAuthorize")).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.CreateCity(context.TODO(), &tempMockCity,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})

}

func TestUpdateCity(t *testing.T) {

	now := time.Now()
	mockCity.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {
		tempMockCity := models.NewCommandCity{
			Id:         1,
			CityName:   mockCity.CityName,
			CityDesc:   mockCity.CityDesc,
			CityPhotos: coverPhoto,
			ProvinceId: 1,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCity.Id = 0
		mockCPCRepo.On("UpdateCity", mock.Anything, mock.AnythingOfType("*models.City")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.UpdateCity(context.TODO(), &tempMockCity,token)

		assert.NoError(t, err)
		assert.Equal(t, mockCity.CityName, tempMockCity.CityName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("success-without-image", func(t *testing.T) {
		tempMockCity := models.NewCommandCity{
			Id:         1,
			CityName:   mockCity.CityName,
			CityDesc:   mockCity.CityDesc,
			CityPhotos: nil,
			ProvinceId: 1,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCity.Id = 0
		tempMockCity.CityPhotos = make([]models.CoverPhotosObj,0)
		mockCPCRepo.On("UpdateCity", mock.Anything, mock.AnythingOfType("*models.City")).Return(nil).Once()
		//mockCPCRepo.On("GetCityByID", mock.Anything, mock.AnythingOfType("int")).Return(&mockCity,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.UpdateCity(context.TODO(), &tempMockCity,token)

		assert.NoError(t, err)
		assert.Equal(t, mockCity.CityName, tempMockCity.CityName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockCity := models.NewCommandCity{
			Id:         1,
			CityName:   mockCity.CityName,
			CityDesc:   mockCity.CityDesc,
			CityPhotos: coverPhoto,
			ProvinceId: 1,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCity.Id = 0
		mockCPCRepo.On("UpdateCity", mock.Anything, mock.AnythingOfType("*models.City")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.UpdateCity(context.TODO(), &tempMockCity,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})
}

func TestDeleteCity(t *testing.T) {

	now := time.Now()
	mockCity.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockCPCRepo.On("DeleteCity", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.DeleteCity(context.TODO(), mockCity.Id,token)

		assert.NoError(t, err)
		//assert.Equal(t, mockCity.CityName, tempMockInclude.IncludeName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockCPCRepo.On("DeleteCity", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.DeleteCity(context.TODO(), mockCity.Id,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})
}

func TestGetAllProvince(t *testing.T) {
	var mockListProvince []*models.Province
	mockListProvince = append(mockListProvince, &mockProvince)
	page := 1
	limit := 1
	offset := 0
	count := 1
	t.Run("success", func(t *testing.T) {
		mockCPCRepo.On("FetchProvince", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(mockListProvince, nil).Once()
		mockCPCRepo.On("GetCountProvince", mock.Anything).Return(count, nil).Once()

		u := usecase.NewCPCUsecase(nil, mockCPCRepo, timeoutContext)

		_, err := u.GetAllProvince(context.TODO(),page,limit,offset)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListInclude))

		//mockCPCRepo.AssertExpectations(t)
		//mockCPCRepo.AssertExpectations(t)
	})

}

func TestGetProvinceById(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockCPCRepo.On("GetProvinceByID", mock.Anything, mock.AnythingOfType("int")).Return(&mockProvince, nil).Once()
		u := usecase.NewCPCUsecase(nil, mockCPCRepo, timeoutContext)

		a, err := u.GetProvinceById(context.TODO(), mockProvince.Id)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockCPCRepo.On("GetProvinceByID", mock.Anything, mock.AnythingOfType("int")).Return(nil, errors.New("Unexpected")).Once()

		u := usecase.NewCPCUsecase(nil, mockCPCRepo, timeoutContext)

		a, err := u.GetProvinceById(context.TODO(), mockProvince.Id)

		assert.Error(t, err)
		assert.Nil(t, a)

		//mockCPCRepo.AssertExpectations(t)
	})

}

func TestCreateProvince(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		tempMockProvince := models.NewCommandProvince{
			Id:                         1,
			ProvinceName:               mockProvince.ProvinceName,
			CountryId:                  1,
			ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockProvince.Id = 0
		mockCPCRepo.On("InsertProvince", mock.Anything, mock.AnythingOfType("*models.Province")).Return(&mockProvince.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.CreateProvince(context.TODO(), &tempMockProvince,token)

		assert.NoError(t, err)
		assert.Equal(t, mockProvince.ProvinceName, tempMockProvince.ProvinceName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockProvince := models.NewCommandProvince{
			Id:                         1,
			ProvinceName:               mockProvince.ProvinceName,
			CountryId:                  1,
			ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockProvince.Id = 0
		mockCPCRepo.On("InsertProvince", mock.Anything, mock.AnythingOfType("*models.Province")).Return(&mockProvince.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("UnAuthorize")).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.CreateProvince(context.TODO(), &tempMockProvince,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})

}

func TestUpdateProvince(t *testing.T) {

	now := time.Now()
	mockProvince.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {
		tempMockProvince := models.NewCommandProvince{
			Id:                         1,
			ProvinceName:               mockProvince.ProvinceName,
			CountryId:                  1,
			ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockProvince.Id = 0
		mockCPCRepo.On("UpdateProvince", mock.Anything, mock.AnythingOfType("*models.Province")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.UpdateProvince(context.TODO(), &tempMockProvince,token)

		assert.NoError(t, err)
		assert.Equal(t, mockProvince.ProvinceName, tempMockProvince.ProvinceName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("success-without-image", func(t *testing.T) {
		tempMockProvince := models.NewCommandProvince{
			Id:                         1,
			ProvinceName:               mockProvince.ProvinceName,
			CountryId:                  1,
			ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockProvince.Id = 0
		//tempMockProvince.CityPhotos = make([]models.CoverPhotosObj,0)
		mockCPCRepo.On("UpdateProvince", mock.Anything, mock.AnythingOfType("*models.Province")).Return(nil).Once()
		//mockCPCRepo.On("GetCityByID", mock.Anything, mock.AnythingOfType("int")).Return(&mockCity,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.UpdateProvince(context.TODO(), &tempMockProvince,token)

		assert.NoError(t, err)
		assert.Equal(t, mockProvince.ProvinceName, tempMockProvince.ProvinceName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockProvince := models.NewCommandProvince{
			Id:                         1,
			ProvinceName:               mockProvince.ProvinceName,
			CountryId:                  1,
			ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockProvince.Id = 0
		mockCPCRepo.On("UpdateProvince", mock.Anything, mock.AnythingOfType("*models.Province")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.UpdateProvince(context.TODO(), &tempMockProvince,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})
}

func TestDeleteProvince(t *testing.T) {

	now := time.Now()
	mockProvince.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockCPCRepo.On("DeleteProvince", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.DeleteProvince(context.TODO(), mockProvince.Id,token)

		assert.NoError(t, err)
		//assert.Equal(t, mockCity.CityName, tempMockInclude.IncludeName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockCPCRepo.On("DeleteProvince", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.DeleteProvince(context.TODO(), mockProvince.Id,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})
}

func TestGetAllCountry(t *testing.T) {
	var mockListCountry []*models.Country
	mockListCountry = append(mockListCountry, &mockCountry)
	page := 1
	limit := 1
	offset := 0
	count := 1
	t.Run("success", func(t *testing.T) {
		mockCPCRepo.On("FetchCountry", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(mockListCountry, nil).Once()
		mockCPCRepo.On("GetCountCountry", mock.Anything).Return(count, nil).Once()

		u := usecase.NewCPCUsecase(nil, mockCPCRepo, timeoutContext)

		_, err := u.GetAllCountry(context.TODO(),page,limit,offset)

		assert.NoError(t, err)
		//assert.Len(t, list, len(mockListInclude))

		//mockCPCRepo.AssertExpectations(t)
		//mockCPCRepo.AssertExpectations(t)
	})

}

func TestGetCountryById(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		mockCPCRepo.On("GetCountryByID", mock.Anything, mock.AnythingOfType("int")).Return(&mockCountry, nil).Once()
		u := usecase.NewCPCUsecase(nil, mockCPCRepo, timeoutContext)

		a, err := u.GetCountryById(context.TODO(), mockProvince.Id)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockCPCRepo.On("GetCountryByID", mock.Anything, mock.AnythingOfType("int")).Return(nil, errors.New("Unexpected")).Once()

		u := usecase.NewCPCUsecase(nil, mockCPCRepo, timeoutContext)

		a, err := u.GetCountryById(context.TODO(), mockProvince.Id)

		assert.Error(t, err)
		assert.Nil(t, a)

		//mockCPCRepo.AssertExpectations(t)
	})

}

func TestCreateCountry(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		tempMockCountry := models.NewCommandCountry{
			Id:          1,
			CountryName: mockCountry.CountryName,
			Iso:         mockCountry.Iso,
			Name:        mockCountry.Name,
			NiceName:    mockCountry.NiceName,
			Iso3:        mockCountry.Iso3,
			NumCode:     mockCountry.NumCode,
			PhoneCode:   mockCountry.PhoneCode,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCountry.Id = 0
		mockCPCRepo.On("InsertCountry", mock.Anything, mock.AnythingOfType("*models.Country")).Return(&mockCountry.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.CreateCountry(context.TODO(), &tempMockCountry,token)

		assert.NoError(t, err)
		assert.Equal(t, mockCountry.CountryName, tempMockCountry.CountryName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockCountry := models.NewCommandCountry{
			Id:          1,
			CountryName: mockCountry.CountryName,
			Iso:         mockCountry.Iso,
			Name:        mockCountry.Name,
			NiceName:    mockCountry.NiceName,
			Iso3:        mockCountry.Iso3,
			NumCode:     mockCountry.NumCode,
			PhoneCode:   mockCountry.PhoneCode,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCountry.Id = 0
		mockCPCRepo.On("InsertCountry", mock.Anything, mock.AnythingOfType("*models.Country")).Return(&mockCountry.Id,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("UnAuthorize")).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.CreateCountry(context.TODO(), &tempMockCountry,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})

}

func TestUpdateCountry(t *testing.T) {

	now := time.Now()
	mockCountry.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {
		tempMockCountry := models.NewCommandCountry{
			Id:          1,
			CountryName: mockCountry.CountryName,
			Iso:         mockCountry.Iso,
			Name:        mockCountry.Name,
			NiceName:    mockCountry.NiceName,
			Iso3:        mockCountry.Iso3,
			NumCode:     mockCountry.NumCode,
			PhoneCode:   mockCountry.PhoneCode,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCountry.Id = 0
		mockCPCRepo.On("UpdateCountry", mock.Anything, mock.AnythingOfType("*models.Country")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.UpdateCountry(context.TODO(), &tempMockCountry,token)

		assert.NoError(t, err)
		assert.Equal(t, mockCountry.CountryName, tempMockCountry.CountryName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("success-without-image", func(t *testing.T) {
		tempMockCountry := models.NewCommandCountry{
			Id:          1,
			CountryName: mockCountry.CountryName,
			Iso:         mockCountry.Iso,
			Name:        mockCountry.Name,
			NiceName:    mockCountry.NiceName,
			Iso3:        mockCountry.Iso3,
			NumCode:     mockCountry.NumCode,
			PhoneCode:   mockCountry.PhoneCode,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCountry.Id = 0
		//tempMockProvince.CityPhotos = make([]models.CoverPhotosObj,0)
		mockCPCRepo.On("UpdateCountry", mock.Anything, mock.AnythingOfType("*models.Country")).Return(nil).Once()
		//mockCPCRepo.On("GetCityByID", mock.Anything, mock.AnythingOfType("int")).Return(&mockCity,nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.UpdateCountry(context.TODO(), &tempMockCountry,token)

		assert.NoError(t, err)
		assert.Equal(t, mockCountry.CountryName, tempMockCountry.CountryName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {
		tempMockCountry := models.NewCommandCountry{
			Id:          1,
			CountryName: mockCountry.CountryName,
			Iso:         mockCountry.Iso,
			Name:        mockCountry.Name,
			NiceName:    mockCountry.NiceName,
			Iso3:        mockCountry.Iso3,
			NumCode:     mockCountry.NumCode,
			PhoneCode:   mockCountry.PhoneCode,
		}
		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		tempMockCountry.Id = 0
		mockCPCRepo.On("UpdateCountry", mock.Anything, mock.AnythingOfType("*models.Country")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.UpdateCountry(context.TODO(), &tempMockCountry,token)

		assert.Error(t, err)
		//assert.Equal(t, mockInclude.IncludeName, tempMockInclude.IncludeName)
		//mockIncludeRepo.AssertExpectations(t)
	})
}

func TestDeleteCountry(t *testing.T) {

	now := time.Now()
	mockCountry.ModifiedDate = &now
	t.Run("success", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockCPCRepo.On("DeleteCountry", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(mockAdmin,nil).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.DeleteCountry(context.TODO(), mockCountry.Id,token)

		assert.NoError(t, err)
		//assert.Equal(t, mockCity.CityName, tempMockInclude.IncludeName)
		//mockCPCRepo.AssertExpectations(t)
	})
	t.Run("error-unauthorize", func(t *testing.T) {

		token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

		mockCPCRepo.On("DeleteCountry", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil).Once()
		mockAdminUsecase.On("ValidateTokenAdmin", mock.Anything, mock.AnythingOfType("string")).Return(nil,errors.New("unAuthorize")).Once()

		u := usecase.NewCPCUsecase(mockAdminUsecase, mockCPCRepo, timeoutContext)

		_,err := u.DeleteCountry(context.TODO(), mockCountry.Id,token)

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

