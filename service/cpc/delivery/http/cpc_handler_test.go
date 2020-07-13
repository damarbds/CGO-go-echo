package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	//_adminUsecaseMock "github.com/auth/admin/mocks"
	_mockIdentityserver "github.com/auth/identityserver/mocks"
	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/cpc/mocks"
	//includeHttp "github.com/service/include/delivery/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	//"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	cpcHttp "github.com/service/cpc/delivery/http"
)
var (
	mockUCase = new(mocks.Usecase)
	mockIsUsecase = new(_mockIdentityserver.Usecase)
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
func TestListCity(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockCityPagination = models.CityDtoWithPagination{}

	err := faker.FakeData(&mockCityPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/city?page=" + strconv.Itoa(mockCityPagination.Meta.Page) + "&size=" + strconv.Itoa(mockCityPagination.Meta.RecordPerPage) , strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAllCity", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(&mockCityPagination,nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)


	handler := cpcHttp.CPCHandler{
		CPCUsecase:mockUCase,
	}
	err = handler.ListCity(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestListCityErrorInternalServer(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockCityPagination = models.CityDtoWithPagination{}

	err := faker.FakeData(&mockCityPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/city?page=" + strconv.Itoa(mockCityPagination.Meta.Page) + "&size=" + strconv.Itoa(mockCityPagination.Meta.RecordPerPage) , strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAllCity", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(nil,errors.New("Internal server Error"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)


	handler := cpcHttp.CPCHandler{
		CPCUsecase:mockUCase,
	}
	err = handler.ListCity(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetDetailCityID(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockCity models.CityDto
	err := faker.FakeData(&mockCity)
	assert.NoError(t, err)

	num := int(mockCity.Id)

	mockUCase.On("GetCityById", mock.Anything, int(num)).Return(&mockCity, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/city/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/city/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.GetDetailCityID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetDetailCityIDErrorNotFound(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockCity models.CityDto
	err := faker.FakeData(&mockCity)
	assert.NoError(t, err)

	num := int(mockCity.Id)

	mockUCase.On("GetCityById", mock.Anything, int(num)).Return(nil, models.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/city/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/city/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.GetDetailCityID(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateCity(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCity{
		Id:         1,
		CityName:   mockCity.CityName,
		CityDesc:   mockCity.CityDesc,
		CityPhotos: coverPhoto,
		ProvinceId: 1,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateCity", mock.Anything, mock.AnythingOfType("*models.NewCommandCity"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)
	//var param = url.Values{}

	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("city_name", tempMockCity.CityName)
	writer.WriteField("city_desc", tempMockCity.CityDesc)
	writer.WriteField("province_id", strconv.Itoa(tempMockCity.ProvinceId))
	part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/city", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/city")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.CreateCity(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreateCityWithoutToken(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCity{
		Id:         1,
		CityName:   mockCity.CityName,
		CityDesc:   mockCity.CityDesc,
		CityPhotos: coverPhoto,
		ProvinceId: 1,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateCity", mock.Anything, mock.AnythingOfType("*models.NewCommandCity"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)
	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("city_name", tempMockCity.CityName)
	writer.WriteField("city_desc", tempMockCity.CityDesc)
	writer.WriteField("province_id", strconv.Itoa(tempMockCity.ProvinceId))
	part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/city", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/city")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.CreateCity(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateCityConflict(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCity{
		Id:         1,
		CityName:   mockCity.CityName,
		CityDesc:   mockCity.CityDesc,
		CityPhotos: coverPhoto,
		ProvinceId: 1,
	}
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Include",
	//}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateCity", mock.Anything, mock.AnythingOfType("*models.NewCommandCity"),mock.AnythingOfType("string")).Return(nil,models.ErrConflict)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)
	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("city_name", tempMockCity.CityName)
	writer.WriteField("city_desc", tempMockCity.CityDesc)
	writer.WriteField("province_id", strconv.Itoa(tempMockCity.ProvinceId))
	part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/city", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/city")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.CreateCity(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateCity(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCity{
		Id:         1,
		CityName:   mockCity.CityName,
		CityDesc:   mockCity.CityDesc,
		CityPhotos: coverPhoto,
		ProvinceId: 1,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("UpdateCity", mock.Anything, mock.AnythingOfType("*models.NewCommandCity"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)

	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("city_name", tempMockCity.CityName)
	writer.WriteField("city_desc", tempMockCity.CityDesc)
	writer.WriteField("province_id", strconv.Itoa(tempMockCity.ProvinceId))
	part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/city/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/city/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.UpdateCity(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdateCityWithoutToken(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCity{
		Id:         1,
		CityName:   mockCity.CityName,
		CityDesc:   mockCity.CityDesc,
		CityPhotos: coverPhoto,
		ProvinceId: 1,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("UpdateCity", mock.Anything, mock.AnythingOfType("*models.NewCommandCity"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)

	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("city_name", tempMockCity.CityName)
	writer.WriteField("city_desc", tempMockCity.CityDesc)
	writer.WriteField("province_id", strconv.Itoa(tempMockCity.ProvinceId))
	part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/city/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/city/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.UpdateCity(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateCityBadParam(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCity{
		Id:         1,
		CityName:   mockCity.CityName,
		CityDesc:   mockCity.CityDesc,
		CityPhotos: coverPhoto,
		ProvinceId: 1,
	}	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("UpdateCity", mock.Anything, mock.AnythingOfType("*models.NewCommandCity"),mock.AnythingOfType("string")).Return(nil,models.ErrBadParamInput)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)

	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("city_name", tempMockCity.CityName)
	writer.WriteField("city_desc", tempMockCity.CityDesc)
	writer.WriteField("province_id", strconv.Itoa(tempMockCity.ProvinceId))
	part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/city/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/city/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.UpdateCity(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteCity(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCity{
		Id:         1,
		CityName:   mockCity.CityName,
		CityDesc:   mockCity.CityDesc,
		CityPhotos: coverPhoto,
		ProvinceId: 1,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("DeleteCity", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockReponse,nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/city/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/city/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.DeleteCity(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteCityWithoutToken(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCity{
		Id:         1,
		CityName:   mockCity.CityName,
		CityDesc:   mockCity.CityDesc,
		CityPhotos: coverPhoto,
		ProvinceId: 1,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("DeleteCity", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockReponse,nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/city/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/city/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.DeleteCity(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteCityErrorInternalServer(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCity{
		Id:         1,
		CityName:   mockCity.CityName,
		CityDesc:   mockCity.CityDesc,
		CityPhotos: coverPhoto,
		ProvinceId: 1,
	}
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Include",
	//}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("DeleteCity", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil,models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/city/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/city/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.DeleteCity(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestListProvince(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockProvincePagination = models.ProvinceDtoWithPagination{}

	err := faker.FakeData(&mockProvincePagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/province?page=" + strconv.Itoa(mockProvincePagination.Meta.Page) + "&size=" + strconv.Itoa(mockProvincePagination.Meta.RecordPerPage) , strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAllProvince", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(&mockProvincePagination,nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)


	handler := cpcHttp.CPCHandler{
		CPCUsecase:mockUCase,
	}
	err = handler.ListProvince(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestListProvinceErrorInternalServer(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockProvincePagination = models.ProvinceDtoWithPagination{}

	err := faker.FakeData(&mockProvincePagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/province?page=" + strconv.Itoa(mockProvincePagination.Meta.Page) + "&size=" + strconv.Itoa(mockProvincePagination.Meta.RecordPerPage) , strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAllProvince", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(nil,errors.New("Internal server Error"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)


	handler := cpcHttp.CPCHandler{
		CPCUsecase:mockUCase,
	}
	err = handler.ListProvince(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetDetailProvinceID(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockProvince models.ProvinceDto
	err := faker.FakeData(&mockProvince)
	assert.NoError(t, err)

	num := int(mockProvince.Id)

	mockUCase.On("GetProvinceById", mock.Anything, int(num)).Return(&mockProvince, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/province/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/province/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.GetDetailProvinceID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetDetailProvinceIDErrorNotFound(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockProvince models.ProvinceDto
	err := faker.FakeData(&mockProvince)
	assert.NoError(t, err)

	num := int(mockProvince.Id)

	mockUCase.On("GetProvinceById", mock.Anything, int(num)).Return(nil, models.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/province/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/province/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.GetDetailProvinceID(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateProvince(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandProvince{
		Id:                         1,
		ProvinceName:               mockProvince.ProvinceName,
		CountryId:                  mockProvince.CountryId,
		ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateProvince", mock.Anything, mock.AnythingOfType("*models.NewCommandProvince"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	//mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)
	//var param = url.Values{}

	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("province_name", tempMockCity.ProvinceName)
	writer.WriteField("country_id", strconv.Itoa(tempMockCity.CountryId))
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/province", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/province")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.CreateProvince(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreateProvinceWithoutToken(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandProvince{
		Id:                         1,
		ProvinceName:               mockProvince.ProvinceName,
		CountryId:                  mockProvince.CountryId,
		ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateProvince", mock.Anything, mock.AnythingOfType("*models.NewCommandProvince"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)
	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("province_name", tempMockCity.ProvinceName)
	writer.WriteField("country_id", strconv.Itoa(tempMockCity.CountryId))
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/province", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/province")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.CreateProvince(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateProvinceConflict(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandProvince{
		Id:                         1,
		ProvinceName:               mockProvince.ProvinceName,
		CountryId:                  mockProvince.CountryId,
		ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
	}
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Include",
	//}


	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateProvince", mock.Anything, mock.AnythingOfType("*models.NewCommandProvince"),mock.AnythingOfType("string")).Return(nil,models.ErrConflict)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)
	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("province_name", tempMockCity.ProvinceName)
	writer.WriteField("country_id", strconv.Itoa(tempMockCity.CountryId))
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/province", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/province")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.CreateProvince(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateProvince(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandProvince{
		Id:                         1,
		ProvinceName:               mockProvince.ProvinceName,
		CountryId:                  mockProvince.CountryId,
		ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("UpdateProvince", mock.Anything, mock.AnythingOfType("*models.NewCommandProvince"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)

	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("province_name", tempMockCity.ProvinceName)
	writer.WriteField("country_id", strconv.Itoa(tempMockCity.CountryId))
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/province/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/province/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.UpdateProvince(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdateProvinceWithoutToken(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandProvince{
		Id:                         1,
		ProvinceName:               mockProvince.ProvinceName,
		CountryId:                  mockProvince.CountryId,
		ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("UpdateProvince", mock.Anything, mock.AnythingOfType("*models.NewCommandProvince"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)

	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("province_name", tempMockCity.ProvinceName)
	writer.WriteField("country_id", strconv.Itoa(tempMockCity.CountryId))
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/province/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/province/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.UpdateProvince(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateProvinceBadParam(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandProvince{
		Id:                         1,
		ProvinceName:               mockProvince.ProvinceName,
		CountryId:                  mockProvince.CountryId,
		ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("UpdateProvince", mock.Anything, mock.AnythingOfType("*models.NewCommandProvince"),mock.AnythingOfType("string")).Return(nil,models.ErrBadParamInput)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)

	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("province_name", tempMockCity.ProvinceName)
	writer.WriteField("country_id", strconv.Itoa(tempMockCity.CountryId))
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/province/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/province/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.UpdateProvince(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteProvince(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandProvince{
		Id:                         1,
		ProvinceName:               mockProvince.ProvinceName,
		CountryId:                  mockProvince.CountryId,
		ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("DeleteProvince", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockReponse,nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/province/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/province/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.DeleteProvince(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteProvinceWithoutToken(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandProvince{
		Id:                         1,
		ProvinceName:               mockProvince.ProvinceName,
		CountryId:                  mockProvince.CountryId,
		ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("DeleteProvince", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockReponse,nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/province/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/province/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.DeleteProvince(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteProvinceErrorInternalServer(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandProvince{
		Id:                         1,
		ProvinceName:               mockProvince.ProvinceName,
		CountryId:                  mockProvince.CountryId,
		ProvinceNameTransportation: mockProvince.ProvinceNameTransportation,
	}
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Include",
	//}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("DeleteProvince", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil,models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/province/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/province/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.DeleteProvince(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestListCountry(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockProvincePagination = models.CountryDtoWithPagination{}

	err := faker.FakeData(&mockProvincePagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/country?page=" + strconv.Itoa(mockProvincePagination.Meta.Page) + "&size=" + strconv.Itoa(mockProvincePagination.Meta.RecordPerPage) , strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAllCountry", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(&mockProvincePagination,nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)


	handler := cpcHttp.CPCHandler{
		CPCUsecase:mockUCase,
	}
	err = handler.ListCountry(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestListCountryErrorInternalServer(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockProvincePagination = models.CountryDtoWithPagination{}

	err := faker.FakeData(&mockProvincePagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/country?page=" + strconv.Itoa(mockProvincePagination.Meta.Page) + "&size=" + strconv.Itoa(mockProvincePagination.Meta.RecordPerPage) , strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAllCountry", mock.Anything,mock.AnythingOfType("int"),mock.AnythingOfType("int"),mock.AnythingOfType("int")).Return(nil,errors.New("Internal server Error"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)


	handler := cpcHttp.CPCHandler{
		CPCUsecase:mockUCase,
	}
	err = handler.ListCountry(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetDetailCountryID(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockProvince models.CountryDto
	err := faker.FakeData(&mockProvince)
	assert.NoError(t, err)

	num := int(mockProvince.Id)

	mockUCase.On("GetCountryById", mock.Anything, int(num)).Return(&mockProvince, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/country/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/country/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.GetDetailCountryID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetDetailCountryIDErrorNotFound(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	var mockProvince models.CountryDto
	err := faker.FakeData(&mockProvince)
	assert.NoError(t, err)

	num := int(mockProvince.Id)

	mockUCase.On("GetCountryById", mock.Anything, int(num)).Return(nil, models.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/country/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/country/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.GetDetailCountryID(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateCountry(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCountry{
		Id:          1,
		CountryName: mockCountry.CountryName,
		Iso:         mockCountry.Iso,
		Name:        mockCountry.Name,
		NiceName:    mockCountry.NiceName,
		Iso3:        mockCountry.Iso3,
		NumCode:     mockCountry.NumCode,
		PhoneCode:   mockCountry.PhoneCode,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateCountry", mock.Anything, mock.AnythingOfType("*models.NewCommandCountry"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	//mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)
	//var param = url.Values{}

	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("country_name", tempMockCity.CountryName)
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/country", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/country")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.CreateCountry(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreateCountryWithoutToken(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCountry{
		Id:          1,
		CountryName: mockCountry.CountryName,
		Iso:         mockCountry.Iso,
		Name:        mockCountry.Name,
		NiceName:    mockCountry.NiceName,
		Iso3:        mockCountry.Iso3,
		NumCode:     mockCountry.NumCode,
		PhoneCode:   mockCountry.PhoneCode,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateCountry", mock.Anything, mock.AnythingOfType("*models.NewCommandCountry"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)
	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("country_name", tempMockCity.CountryName)
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/country", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/country")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.CreateCountry(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateCountryConflict(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCountry{
		Id:          1,
		CountryName: mockCountry.CountryName,
		Iso:         mockCountry.Iso,
		Name:        mockCountry.Name,
		NiceName:    mockCountry.NiceName,
		Iso3:        mockCountry.Iso3,
		NumCode:     mockCountry.NumCode,
		PhoneCode:   mockCountry.PhoneCode,
	}
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Include",
	//}


	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("CreateCountry", mock.Anything, mock.AnythingOfType("*models.NewCommandCountry"),mock.AnythingOfType("string")).Return(nil,models.ErrConflict)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)
	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("country_name", tempMockCity.CountryName)
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/country", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/country")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.CreateCountry(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateCountry(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCountry{
		Id:          1,
		CountryName: mockCountry.CountryName,
		Iso:         mockCountry.Iso,
		Name:        mockCountry.Name,
		NiceName:    mockCountry.NiceName,
		Iso3:        mockCountry.Iso3,
		NumCode:     mockCountry.NumCode,
		PhoneCode:   mockCountry.PhoneCode,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("UpdateCountry", mock.Anything, mock.AnythingOfType("*models.NewCommandCountry"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)

	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("country_name", tempMockCity.CountryName)
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/country/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/country/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.UpdateCountry(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdateCountryWithoutToken(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCountry{
		Id:          1,
		CountryName: mockCountry.CountryName,
		Iso:         mockCountry.Iso,
		Name:        mockCountry.Name,
		NiceName:    mockCountry.NiceName,
		Iso3:        mockCountry.Iso3,
		NumCode:     mockCountry.NumCode,
		PhoneCode:   mockCountry.PhoneCode,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("UpdateCountry", mock.Anything, mock.AnythingOfType("*models.NewCommandCountry"),mock.AnythingOfType("string")).Return(mockReponse,nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)

	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("country_name", tempMockCity.CountryName)
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/country/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/country/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.UpdateCountry(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateCountryBadParam(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCountry{
		Id:          1,
		CountryName: mockCountry.CountryName,
		Iso:         mockCountry.Iso,
		Name:        mockCountry.Name,
		NiceName:    mockCountry.NiceName,
		Iso3:        mockCountry.Iso3,
		NumCode:     mockCountry.NumCode,
		PhoneCode:   mockCountry.PhoneCode,
	}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("UpdateCountry", mock.Anything, mock.AnythingOfType("*models.NewCommandCountry"),mock.AnythingOfType("string")).Return(nil,models.ErrBadParamInput)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(imageTestPath,nil)

	//var param = url.Values{}
	//param.Set("include_name", tempMockInclude.IncludeName)
	//var payload = bytes.NewBufferString(param.Encode())

	dir, err := os.Getwd()
	path := dir + "\\pp.jpg"
	file, err := os.Open(path)
	if err != nil {
		t.Error(err)
	}

	defer file.Close()
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("country_name", tempMockCity.CountryName)
	//part, err := writer.CreateFormFile("city_image", filepath.Base(path))
	//if err != nil {
	//	writer.Close()
	//	t.Error(err)
	//}
	//io.Copy(part, file)
	//writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/country/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/country/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
		IsUsecase:mockIsUsecase,
	}
	err = handler.UpdateCountry(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteCountry(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCountry{
		Id:          1,
		CountryName: mockCountry.CountryName,
		Iso:         mockCountry.Iso,
		Name:        mockCountry.Name,
		NiceName:    mockCountry.NiceName,
		Iso3:        mockCountry.Iso3,
		NumCode:     mockCountry.NumCode,
		PhoneCode:   mockCountry.PhoneCode,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("DeleteCountry", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockReponse,nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/country/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/country/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.DeleteCountry(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteCountryWithoutToken(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCountry{
		Id:          1,
		CountryName: mockCountry.CountryName,
		Iso:         mockCountry.Iso,
		Name:        mockCountry.Name,
		NiceName:    mockCountry.NiceName,
		Iso3:        mockCountry.Iso3,
		NumCode:     mockCountry.NumCode,
		PhoneCode:   mockCountry.PhoneCode,
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Include",
	}

	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("DeleteCountry", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(mockReponse,nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/country/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/country/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.DeleteCountry(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteCountryErrorInternalServer(t *testing.T) {
	mockUCase = new(mocks.Usecase)
	tempMockCity := models.NewCommandCountry{
		Id:          1,
		CountryName: mockCountry.CountryName,
		Iso:         mockCountry.Iso,
		Name:        mockCountry.Name,
		NiceName:    mockCountry.NiceName,
		Iso3:        mockCountry.Iso3,
		NumCode:     mockCountry.NumCode,
		PhoneCode:   mockCountry.PhoneCode,
	}
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Include",
	//}
	//j, err := json.Marshal(tempMockInclude)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCity.Id)
	mockUCase.On("DeleteCountry", mock.Anything, mock.AnythingOfType("int"),mock.AnythingOfType("string")).Return(nil,models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/country/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/country/:id")
	c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := cpcHttp.CPCHandler{
		CPCUsecase: mockUCase,
	}
	err = handler.DeleteCountry(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
