package http_test

import (
	"bytes"
	"errors"
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
	"github.com/service/facilities/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	FacilityHttp "github.com/service/facilities/delivery/http"
)

var (
	imagePath      = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Facilities/8941695193938718058.jpg"
	mockFacilities = models.Facilities{
		Id:           1,
		CreatedBy:    "test",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		FacilityName: "Test Facilities 2",
		FacilityIcon: &imagePath,
	}
	tempMockFacilities = models.NewCommandFacilities{
		Id:           1,
		FacilityName: mockFacilities.FacilityName,
		FacilityIcon: *mockFacilities.FacilityIcon,
		IsNumerable:  mockFacilities.IsNumerable,
	}
)

func TestList(t *testing.T) {
	var mockFacility models.FacilityDto
	err := faker.FakeData(&mockFacility)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListFacility []*models.FacilityDto
	mockListFacility = append(mockListFacility, &mockFacility)

	mockUCase.On("List", mock.Anything).Return(mockListFacility, nil)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/facilities", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Request().Header.Add("Authorization", token)

	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.List(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestListWithoutToken(t *testing.T) {
	var mockFacility models.FacilityDto
	err := faker.FakeData(&mockFacility)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListFacility []*models.FacilityDto
	mockListFacility = append(mockListFacility, &mockFacility)

	mockUCase.On("List", mock.Anything).Return(mockListFacility, nil)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/facilities", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//c.Request().Header.Add("Authorization",token)

	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.List(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestListInvalidToken(t *testing.T) {
	var mockFacility models.FacilityDto
	err := faker.FakeData(&mockFacility)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListFacility []*models.FacilityDto
	mockListFacility = append(mockListFacility, &mockFacility)

	mockUCase.On("List", mock.Anything).Return(nil, models.ErrUnAuthorize)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/facilities", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Request().Header.Add("Authorization", token)

	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.List(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetAllFacility(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockFacilityPagination = models.FacilityDtoWithPagination{}

	err := faker.FakeData(&mockFacilityPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/facilities?page="+strconv.Itoa(mockFacilityPagination.Meta.Page)+"&size="+strconv.Itoa(mockFacilityPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(&mockFacilityPagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.GetAllFacilities(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetAllFacilityErrorInternalServer(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockFacilityPagination = models.FacilityDtoWithPagination{}

	err := faker.FakeData(&mockFacilityPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/facilities?page="+strconv.Itoa(mockFacilityPagination.Meta.Page)+"&size="+strconv.Itoa(mockFacilityPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, errors.New("Internal server Error"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.GetAllFacilities(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetDetailFacilityID(t *testing.T) {
	var mockFacility models.FacilityDto
	err := faker.FakeData(&mockFacility)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockFacility.Id)

	mockUCase.On("GetById", mock.Anything, int(num)).Return(&mockFacility, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/facilities/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/facilities/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.GetDetailFacilitiesID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetDetailFacilityIDErrorNotFound(t *testing.T) {
	var mockFacility models.FacilityDto
	err := faker.FakeData(&mockFacility)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockFacility.Id)

	mockUCase.On("GetById", mock.Anything, int(num)).Return(nil, models.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/facilities/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/facilities/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.GetDetailFacilitiesID(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateFacility(t *testing.T) {
	mockFacility := tempMockFacilities
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Facility",
	}
	tempMockFacility := mockFacility
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockFacility)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandFacilities"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockFacility.FacilityIcon, nil)
	var param = url.Values{}
	param.Set("Facility_name", tempMockFacility.FacilityName)
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
	writer.WriteField("id", strconv.Itoa(tempMockFacility.Id))
	writer.WriteField("facilities_name", tempMockFacility.FacilityName)
	writer.WriteField("is_numerable", tempMockFacility.FacilityName)
	part, err := writer.CreateFormFile("facilities_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/facilities", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/facilities")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
		IsUsecase:       mockIsUsecase,
	}
	err = handler.CreateFacilities(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreateFacilityWithoutToken(t *testing.T) {
	mockFacility := tempMockFacilities
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Facility",
	}
	tempMockFacility := mockFacility
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockFacility)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandFacilities"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockFacility.FacilityIcon, nil)
	var param = url.Values{}
	param.Set("Facility_name", tempMockFacility.FacilityName)
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
	writer.WriteField("id", strconv.Itoa(tempMockFacility.Id))
	writer.WriteField("facilities_name", tempMockFacility.FacilityName)
	writer.WriteField("is_numerable", tempMockFacility.FacilityName)
	part, err := writer.CreateFormFile("facilities_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/facilities", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/facilities")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
		IsUsecase:       mockIsUsecase,
	}
	err = handler.CreateFacilities(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateFacilityConflict(t *testing.T) {
	mockFacility := tempMockFacilities
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Facility",
	//}
	tempMockFacility := mockFacility
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockFacility)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandFacilities"), mock.AnythingOfType("string")).Return(nil, models.ErrConflict)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockFacility.FacilityIcon, nil)
	var param = url.Values{}
	param.Set("Facility_name", tempMockFacility.FacilityName)
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
	writer.WriteField("id", strconv.Itoa(tempMockFacility.Id))
	writer.WriteField("facilities_name", tempMockFacility.FacilityName)
	writer.WriteField("is_numerable", tempMockFacility.FacilityName)
	part, err := writer.CreateFormFile("facilities_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/facilities", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/facilities")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
		IsUsecase:       mockIsUsecase,
	}
	err = handler.CreateFacilities(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateFacility(t *testing.T) {
	mockFacility := tempMockFacilities
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Facility",
	}
	tempMockFacility := mockFacility
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockFacility)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockFacility.Id)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandFacilities"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockFacility.FacilityIcon, nil)

	var param = url.Values{}
	param.Set("Facility_name", tempMockFacility.FacilityName)
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
	writer.WriteField("id", strconv.Itoa(tempMockFacility.Id))
	writer.WriteField("facilities_name", tempMockFacility.FacilityName)
	writer.WriteField("is_numerable", tempMockFacility.FacilityName)
	part, err := writer.CreateFormFile("facilities_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/facilities/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/facilities/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
		IsUsecase:       mockIsUsecase,
	}
	err = handler.UpdateFacilities(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdateFacilityWithoutToken(t *testing.T) {
	mockFacility := tempMockFacilities
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Facility",
	}
	tempMockFacility := mockFacility
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockFacility)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockFacility.Id)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandFacilities"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockFacility.FacilityIcon, nil)

	var param = url.Values{}
	param.Set("Facility_name", tempMockFacility.FacilityName)
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
	writer.WriteField("id", strconv.Itoa(tempMockFacility.Id))
	writer.WriteField("facilities_name", tempMockFacility.FacilityName)
	writer.WriteField("is_numerable", tempMockFacility.FacilityName)
	part, err := writer.CreateFormFile("facilities_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/facilities/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/facilities/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
		IsUsecase:       mockIsUsecase,
	}
	err = handler.UpdateFacilities(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateFacilityBadParam(t *testing.T) {
	mockFacility := tempMockFacilities

	tempMockFacility := mockFacility
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockFacility)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockFacility.Id)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandFacilities"), mock.AnythingOfType("string")).Return(nil, models.ErrBadParamInput)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockFacility.FacilityIcon, nil)

	var param = url.Values{}
	param.Set("Facility_name", tempMockFacility.FacilityName)
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
	writer.WriteField("id", strconv.Itoa(tempMockFacility.Id))
	writer.WriteField("facilities_name", tempMockFacility.FacilityName)
	writer.WriteField("is_numerable", tempMockFacility.FacilityName)
	part, err := writer.CreateFormFile("facilities_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/facilities/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/facilities/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
		IsUsecase:       mockIsUsecase,
	}
	err = handler.UpdateFacilities(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteFacility(t *testing.T) {
	mockFacility := tempMockFacilities
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Facility",
	}
	tempMockFacility := mockFacility
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockFacility)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockFacility.Id)
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/facilities/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/facilities/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.DeleteFacilities(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteFacilityWithoutToken(t *testing.T) {
	mockFacility := tempMockFacilities
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Facility",
	}
	tempMockFacility := mockFacility
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockFacility)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockFacility.Id)
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/facilities/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/facilities/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.DeleteFacilities(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteFacilityErrorInternalServer(t *testing.T) {
	mockFacility := tempMockFacilities
	tempMockFacility := mockFacility
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockFacility)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockFacility.Id)
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil, models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/facilities/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/facilities/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := FacilityHttp.FacilityHandler{
		FacilityUsecase: mockUCase,
	}
	err = handler.DeleteFacilities(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
