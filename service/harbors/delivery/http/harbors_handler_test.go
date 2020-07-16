package http_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	_mockIdentityserver "github.com/auth/identityserver/mocks"
	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/harbors/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	HarborsHttp "github.com/service/harbors/delivery/http"
)

var (

	imagePath   = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg"
	harborsType = 1
	provinceName = "Jawa Barat"
	mockHarbors = models.Harbors{
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
	tempMockHarbors = models.NewCommandHarbors{
		Id:               mockHarbors.Id,
		HarborsName:      mockHarbors.HarborsName,
		HarborsLongitude: mockHarbors.HarborsLongitude,
		HarborsLatitude:  mockHarbors.HarborsLatitude,
		HarborsImage:     mockHarbors.HarborsImage,
		CityId:           mockHarbors.CityId,
		HarborsType:      *mockHarbors.HarborsType,
	}

)

func TestGetAllHarborsWithPagination(t *testing.T) {
	var mockHarbors models.HarborsWCPCDto
	err := faker.FakeData(&mockHarbors)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListHarbors []*models.HarborsWCPCDto
	mockListHarbors = append(mockListHarbors, &mockHarbors)

	mockUCase.On("GetAllWithJoinCPC", mock.Anything, mock.AnythingOfType("*int"), mock.AnythingOfType("*int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockListHarbors, nil)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/exp-destination?page="+strconv.Itoa(1)+"&size="+strconv.Itoa(2), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//c.Request().Header.Add("Authorization", token)

	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
	}
	err = handler.GetAllHarbors(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetAllHarborsWithoutPagination(t *testing.T) {
	var mockHarbors models.HarborsWCPCDto
	err := faker.FakeData(&mockHarbors)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListHarbors []*models.HarborsWCPCDto
	mockListHarbors = append(mockListHarbors, &mockHarbors)

	mockUCase.On("GetAllWithJoinCPC", mock.Anything, mock.AnythingOfType("*int"), mock.AnythingOfType("*int"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockListHarbors, nil)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/exp-destination", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//c.Request().Header.Add("Authorization", token)

	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
	}
	err = handler.GetAllHarbors(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestListHarbors(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockHarborsPagination = models.HarborsDtoWithPagination{}

	err := faker.FakeData(&mockHarborsPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/harbors?page="+strconv.Itoa(mockHarborsPagination.Meta.Page)+"&size="+strconv.Itoa(mockHarborsPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(&mockHarborsPagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
	}
	err = handler.ListHarbors(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestListHarborsErrorInternalServer(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockHarborsPagination = models.HarborsDtoWithPagination{}

	err := faker.FakeData(&mockHarborsPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/harbors?page="+strconv.Itoa(mockHarborsPagination.Meta.Page)+"&size="+strconv.Itoa(mockHarborsPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, errors.New("Internal server Error"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
	}
	err = handler.ListHarbors(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetDetailHarborsID(t *testing.T) {
	var mockHarbors models.HarborsDto
	err := faker.FakeData(&mockHarbors)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := mockHarbors.Id

	mockUCase.On("GetById", mock.Anything, num).Return(&mockHarbors, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/harbors/" + num, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/harbors/:id")
	c.SetParamNames("id")
	c.SetParamValues(num)
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
	}
	err = handler.GetDetailHarborsID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetDetailHarborsIDErrorNotFound(t *testing.T) {
	var mockHarbors models.HarborsDto
	err := faker.FakeData(&mockHarbors)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := mockHarbors.Id

	mockUCase.On("GetById", mock.Anything, num).Return(nil, models.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/harbors/"+num, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/harbors/:id")
	c.SetParamNames("id")
	c.SetParamValues(num)
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
	}
	err = handler.GetDetailHarborsID(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateHarbors(t *testing.T) {
	mockHarbors := tempMockHarbors
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Harbors",
	}
	tempMockHarbors := mockHarbors
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockHarbors)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandHarbors"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockHarbors.HarborsImage, nil)
	//var param = url.Values{}
	//param.Set("Harbors_name", tempMockHarbors.HarborsName)
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
	writer.WriteField("id", tempMockHarbors.Id)
	writer.WriteField("harbors_name", tempMockHarbors.HarborsName)
	writer.WriteField("harbors_longitude", fmt.Sprint(tempMockHarbors.HarborsLongitude))
	writer.WriteField("harbors_latitude", fmt.Sprint(tempMockHarbors.HarborsLatitude))
	writer.WriteField("city_id", strconv.Itoa(tempMockHarbors.CityId))
	writer.WriteField("harbors_type", strconv.Itoa(tempMockHarbors.HarborsType))
	part, err := writer.CreateFormFile("harbors_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/harbors", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/Harbors")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
		IsUsecase:      mockIsUsecase,
	}
	err = handler.CreateHarbors(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreateHarborsWithoutToken(t *testing.T) {
	mockHarbors := tempMockHarbors
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Harbors",
	}
	tempMockHarbors := mockHarbors
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockHarbors)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandHarbors"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockHarbors.HarborsImage, nil)
	var param = url.Values{}
	param.Set("Harbors_name", tempMockHarbors.HarborsName)
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
	writer.WriteField("id", tempMockHarbors.Id)
	writer.WriteField("harbors_name", tempMockHarbors.HarborsName)
	writer.WriteField("harbors_longitude", fmt.Sprint(tempMockHarbors.HarborsLongitude))
	writer.WriteField("harbors_latitude", fmt.Sprint(tempMockHarbors.HarborsLatitude))
	writer.WriteField("city_id", strconv.Itoa(tempMockHarbors.CityId))
	writer.WriteField("harbors_type", strconv.Itoa(tempMockHarbors.HarborsType))
	part, err := writer.CreateFormFile("harbors_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()


	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/harbors", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/harbors")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
		IsUsecase:      mockIsUsecase,
	}
	err = handler.CreateHarbors(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateHarborsConflict(t *testing.T) {
	mockHarbors := tempMockHarbors
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Harbors",
	//}
	tempMockHarbors := mockHarbors
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockHarbors)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandHarbors"), mock.AnythingOfType("string")).Return(nil, models.ErrConflict)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockHarbors.HarborsImage, nil)
	var param = url.Values{}
	param.Set("Harbors_name", tempMockHarbors.HarborsName)
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
	writer.WriteField("id", tempMockHarbors.Id)
	writer.WriteField("harbors_name", tempMockHarbors.HarborsName)
	writer.WriteField("harbors_longitude", fmt.Sprint(tempMockHarbors.HarborsLongitude))
	writer.WriteField("harbors_latitude", fmt.Sprint(tempMockHarbors.HarborsLatitude))
	writer.WriteField("city_id", strconv.Itoa(tempMockHarbors.CityId))
	writer.WriteField("harbors_type", strconv.Itoa(tempMockHarbors.HarborsType))
	part, err := writer.CreateFormFile("harbors_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/harbors", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/harbors")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
		IsUsecase:      mockIsUsecase,
	}
	err = handler.CreateHarbors(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateHarbors(t *testing.T) {
	mockHarbors := tempMockHarbors
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Harbors",
	}
	tempMockHarbors := mockHarbors
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockHarbors)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockHarbors.Id
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandHarbors"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockHarbors.HarborsImage, nil)

	var param = url.Values{}
	param.Set("Harbors_name", tempMockHarbors.HarborsName)
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
	writer.WriteField("id", tempMockHarbors.Id)
	writer.WriteField("harbors_name", tempMockHarbors.HarborsName)
	writer.WriteField("harbors_longitude", fmt.Sprint(tempMockHarbors.HarborsLongitude))
	writer.WriteField("harbors_latitude", fmt.Sprint(tempMockHarbors.HarborsLatitude))
	writer.WriteField("city_id", strconv.Itoa(tempMockHarbors.CityId))
	writer.WriteField("harbors_type", strconv.Itoa(tempMockHarbors.HarborsType))
	part, err := writer.CreateFormFile("harbors_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/harbors/"+id, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/harbors/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
		IsUsecase:      mockIsUsecase,
	}
	err = handler.UpdateHarbors(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdateHarborsWithoutToken(t *testing.T) {
	mockHarbors := tempMockHarbors
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Harbors",
	}
	tempMockHarbors := mockHarbors
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockHarbors)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockHarbors.Id
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandHarbors"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockHarbors.HarborsImage, nil)

	var param = url.Values{}
	param.Set("Harbors_name", tempMockHarbors.HarborsName)
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
	writer.WriteField("id", tempMockHarbors.Id)
	writer.WriteField("harbors_name", tempMockHarbors.HarborsName)
	writer.WriteField("harbors_longitude", fmt.Sprint(tempMockHarbors.HarborsLongitude))
	writer.WriteField("harbors_latitude", fmt.Sprint(tempMockHarbors.HarborsLatitude))
	writer.WriteField("city_id", strconv.Itoa(tempMockHarbors.CityId))
	writer.WriteField("harbors_type", strconv.Itoa(tempMockHarbors.HarborsType))
	part, err := writer.CreateFormFile("harbors_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/harbors/"+id, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/Harbors/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
		IsUsecase:      mockIsUsecase,
	}
	err = handler.UpdateHarbors(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateHarborsBadParam(t *testing.T) {
	mockHarbors := tempMockHarbors

	tempMockHarbors := mockHarbors
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockHarbors)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockHarbors.Id
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandHarbors"), mock.AnythingOfType("string")).Return(nil, models.ErrBadParamInput)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockHarbors.HarborsImage, nil)

	var param = url.Values{}
	param.Set("Harbors_name", tempMockHarbors.HarborsName)
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
	writer.WriteField("id", tempMockHarbors.Id)
	writer.WriteField("harbors_name", tempMockHarbors.HarborsName)
	writer.WriteField("harbors_longitude", fmt.Sprint(tempMockHarbors.HarborsLongitude))
	writer.WriteField("harbors_latitude", fmt.Sprint(tempMockHarbors.HarborsLatitude))
	writer.WriteField("city_id", strconv.Itoa(tempMockHarbors.CityId))
	writer.WriteField("harbors_type", strconv.Itoa(tempMockHarbors.HarborsType))
	part, err := writer.CreateFormFile("harbors_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/harbors/"+id, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/Harbors/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
		IsUsecase:      mockIsUsecase,
	}
	err = handler.UpdateHarbors(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteHarbors(t *testing.T) {
	mockHarbors := tempMockHarbors
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Harbors",
	}
	tempMockHarbors := mockHarbors
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockHarbors)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockHarbors.Id
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/harbors/"+id, nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/harbors/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
	}
	err = handler.DeleteHarbors(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteHarborsWithoutToken(t *testing.T) {
	mockHarbors := tempMockHarbors
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Harbors",
	}
	tempMockHarbors := mockHarbors
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockHarbors)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockHarbors.Id
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/harbors/"+id, nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/Harbors/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
	}
	err = handler.DeleteHarbors(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteHarborsErrorInternalServer(t *testing.T) {
	mockHarbors := tempMockHarbors
	tempMockHarbors := mockHarbors
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockHarbors)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockHarbors.Id
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/harbors/"+id, nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/harbors/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := HarborsHttp.HarborsHandler{
		HarborsUsecase: mockUCase,
	}
	err = handler.DeleteHarbors(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
