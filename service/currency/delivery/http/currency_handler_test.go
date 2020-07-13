package http_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	//_mockIdentityserver "github.com/auth/identityserver/mocks"
	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/currency/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	CurrencyHttp "github.com/service/currency/delivery/http"
)

func TestGetAllCurrency(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockCurrencyPagination = models.CurrencyDtoWithPagination{}

	err := faker.FakeData(&mockCurrencyPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/currency?page="+strconv.Itoa(mockCurrencyPagination.Meta.Page)+"&size="+strconv.Itoa(mockCurrencyPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(&mockCurrencyPagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.GetAllCurrency(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetAllCurrencyErrorInternalServer(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockCurrencyPagination = models.CurrencyDtoWithPagination{}

	err := faker.FakeData(&mockCurrencyPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/currency?page="+strconv.Itoa(mockCurrencyPagination.Meta.Page)+"&size="+strconv.Itoa(mockCurrencyPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, errors.New("Internal server Error"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.GetAllCurrency(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetDetailCurrencyID(t *testing.T) {
	var mockCurrency models.CurrencyDto
	err := faker.FakeData(&mockCurrency)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockCurrency.Id)

	mockUCase.On("GetById", mock.Anything, int(num)).Return(&mockCurrency, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/currency/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/Currency/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.GetDetailCurrencyID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetDetailCurrencyIDErrorNotFound(t *testing.T) {
	var mockCurrency models.CurrencyDto
	err := faker.FakeData(&mockCurrency)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockCurrency.Id)

	mockUCase.On("GetById", mock.Anything, int(num)).Return(nil, models.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/currency/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/currency/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.GetDetailCurrencyID(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateCurrency(t *testing.T) {
	mockCurrency := models.NewCommandCurrency{
		Id:     1,
		Code:   "Test Code 1",
		Name:   "Test Name 1",
		Symbol: "Test Symbol",
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Currency",
	}
	tempMockCurrency := mockCurrency
	j, err := json.Marshal(tempMockCurrency)
	mockUCase := new(mocks.Usecase)
	//mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockCurrency)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandCurrency"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	//mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockCurrency.CurrencyIcon, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/currency", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/currency")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.CreateCurrency(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreateCurrencyWithoutToken(t *testing.T) {
	mockCurrency := models.NewCommandCurrency{
		Id:     1,
		Code:   "Test Code 1",
		Name:   "Test Name 1",
		Symbol: "Test Symbol",
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Currency",
	}
	tempMockCurrency := mockCurrency
	j, err := json.Marshal(tempMockCurrency)
	mockUCase := new(mocks.Usecase)
	//mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockCurrency)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandCurrency"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	//mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockCurrency.CurrencyIcon, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/currency", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/Currency")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.CreateCurrency(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateCurrencyConflict(t *testing.T) {
	mockCurrency := models.NewCommandCurrency{
		Id:     1,
		Code:   "Test Code 1",
		Name:   "Test Name 1",
		Symbol: "Test Symbol",
	}
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Currency",
	//}
	tempMockCurrency := mockCurrency
	j, err := json.Marshal(tempMockCurrency)
	mockUCase := new(mocks.Usecase)
	//mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockCurrency)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandCurrency"), mock.AnythingOfType("string")).Return(nil, models.ErrConflict)
	//mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockCurrency.CurrencyIcon, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/currency", strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/currency")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.CreateCurrency(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateCurrency(t *testing.T) {
	mockCurrency := models.NewCommandCurrency{
		Id:     1,
		Code:   "Test Code 1",
		Name:   "Test Name 1",
		Symbol: "Test Symbol",
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Currency",
	}
	tempMockCurrency := mockCurrency
	j, err := json.Marshal(tempMockCurrency)
	mockUCase := new(mocks.Usecase)
	//mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockCurrency)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCurrency.Id)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandCurrency"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	//mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockCurrency.CurrencyIcon, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/currency/"+strconv.Itoa(id), strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/currency/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.UpdateCurrency(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdateCurrencyWithoutToken(t *testing.T) {
	mockCurrency := models.NewCommandCurrency{
		Id:     1,
		Code:   "Test Code 1",
		Name:   "Test Name 1",
		Symbol: "Test Symbol",
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Currency",
	}
	tempMockCurrency := mockCurrency
	j, err := json.Marshal(tempMockCurrency)
	mockUCase := new(mocks.Usecase)
	//mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockCurrency)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCurrency.Id)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandCurrency"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	//mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockCurrency.CurrencyIcon, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/currency/"+strconv.Itoa(id), strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/Currency/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.UpdateCurrency(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateCurrencyBadParam(t *testing.T) {
	mockCurrency := models.NewCommandCurrency{
		Id:     1,
		Code:   "Test Code 1",
		Name:   "Test Name 1",
		Symbol: "Test Symbol",
	}
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Currency",
	//}
	tempMockCurrency := mockCurrency
	j, err := json.Marshal(tempMockCurrency)
	mockUCase := new(mocks.Usecase)
	//mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockCurrency)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCurrency.Id)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandCurrency"), mock.AnythingOfType("string")).Return(nil, models.ErrBadParamInput)
	//mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockCurrency.CurrencyIcon, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/currency/"+strconv.Itoa(id), strings.NewReader(string(j)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/currency/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.UpdateCurrency(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteCurrency(t *testing.T) {
	mockCurrency := models.NewCommandCurrency{
		Id:     1,
		Code:   "Test Code 1",
		Name:   "Test Name 1",
		Symbol: "Test Symbol",
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Currency",
	}
	tempMockCurrency := mockCurrency
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockCurrency)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCurrency.Id)
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/currency/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/currency/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.DeleteCurrency(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteCurrencyWithoutToken(t *testing.T) {
	mockCurrency := models.NewCommandCurrency{
		Id:     1,
		Code:   "Test Code 1",
		Name:   "Test Name 1",
		Symbol: "Test Symbol",
	}
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Currency",
	}
	tempMockCurrency := mockCurrency
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockCurrency)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCurrency.Id)
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/currency/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/currency/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.DeleteCurrency(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteCurrencyErrorInternalServer(t *testing.T) {
	mockCurrency := models.NewCommandCurrency{
		Id:     1,
		Code:   "Test Code 1",
		Name:   "Test Name 1",
		Symbol: "Test Symbol",
	}
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Currency",
	//}
	tempMockCurrency := mockCurrency
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockCurrency)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockCurrency.Id)
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil, models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/currency/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/currency/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := CurrencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.DeleteCurrency(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
