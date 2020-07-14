package http_test

import (
	"bytes"
	"encoding/json"
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
	"github.com/service/promo/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	PromoHttp "github.com/service/promo/delivery/http"
)
var(
	date = time.Now().String()
	isAnyTripPeriod = 1
	capacity = 20
	maxUsage = 20
	desc = "ini description"
	disc float32 = 12.2
	mockPromo = models.NewCommandPromo{
		Id:                 "asdqeqrqrasdsad",
		PromoCode:          "Test1",
		PromoName:          "Test 1",
		PromoDesc:          "Test 1",
		PromoValue:         1,
		PromoType:          1,
		PromoImage:         "asdasdasdas",
		StartDate:          date,
		EndDate:            date,
		StartTripPeriod:    date,
		EndTripPeriod:      date,
		IsAnyTripPeriod:    isAnyTripPeriod,
		HowToGet:           desc,
		HowToUse:           desc,
		TermCondition:      desc,
		Disclaimer:         desc,
		MaxDiscount:        disc,
		MaxUsage:           isAnyTripPeriod,
		ProductionCapacity: isAnyTripPeriod,
		PromoProductType:   &isAnyTripPeriod,
		MerchantId: []string{
			"adsadsa", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
	}
	merchantJson,_ = json.Marshal(mockPromo.MerchantId)
	merchantString = string(merchantJson)

)
func TestGetAllPromoWithPagination(t *testing.T) {
	var mockPromo models.PromoDto
	err := faker.FakeData(&mockPromo)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListPromo []*models.PromoDto
	mockListPromo = append(mockListPromo, &mockPromo)

	mockUCase.On("Fetch", mock.Anything, mock.AnythingOfType("*int"), mock.AnythingOfType("*int")).Return(mockListPromo, nil)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/special-promo?page="+strconv.Itoa(0)+"&size="+strconv.Itoa(1), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//c.Request().Header.Add("Authorization",token)

	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.GetAllPromo(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetAllPromoWithoutPagination(t *testing.T) {
	var mockPromo models.PromoDto
	err := faker.FakeData(&mockPromo)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListPromo []*models.PromoDto
	mockListPromo = append(mockListPromo, &mockPromo)

	mockUCase.On("Fetch", mock.Anything, mock.AnythingOfType("*int"), mock.AnythingOfType("*int")).Return(mockListPromo, nil)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/special-promo", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//c.Request().Header.Add("Authorization",token)

	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.GetAllPromo(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

//func TestListInvalidToken(t *testing.T) {
//	var mockPromo models.PromoDto
//	err := faker.FakeData(&mockPromo)
//	assert.NoError(t, err)
//	mockUCase := new(mocks.Usecase)
//	var mockListPromo []*models.PromoDto
//	mockListPromo = append(mockListPromo, &mockPromo)
//
//	mockUCase.On("List", mock.Anything).Return(nil, models.ErrUnAuthorize)
//	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
//
//	e := echo.New()
//	req, err := http.NewRequest(echo.GET, "/service/Promo", strings.NewReader(""))
//	assert.NoError(t, err)
//
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	c.Request().Header.Add("Authorization", token)
//
//	handler := PromoHttp.PromoHandler{
//		PromoUsecase: mockUCase,
//	}
//	err = handler.List(c)
//	//require.NoError(t, err)
//	assert.Equal(t, http.StatusUnauthorized, rec.Code)
//	//mockUCase.AssertExpectations(t)
//}

func TestList(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockPromoPagination = models.PromoWithPagination{}

	err := faker.FakeData(&mockPromoPagination)
	assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/admin/promo?page="+strconv.Itoa(mockPromoPagination.Meta.Page)+"&size="+strconv.Itoa(mockPromoPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("List", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"),mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&mockPromoPagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Request().Header.Add("Authorization", token)

	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.List(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestListInternalServer(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockPromoPagination = models.PromoWithPagination{}

	err := faker.FakeData(&mockPromoPagination)
	assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/admin/promo?page="+strconv.Itoa(mockPromoPagination.Meta.Page)+"&size="+strconv.Itoa(mockPromoPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("List", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"),mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, errors.New("Internal server Error"))


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Request().Header.Add("Authorization", token)

	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.List(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
func TestGetPromoByCode(t *testing.T) {
	var mockPromo models.PromoDto
	err := faker.FakeData(&mockPromo)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := mockPromo.PromoCode

	mockUCase.On("GetByCode", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("int"),mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(&mockPromo, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/special-promo/"+num, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("service/special-promo/:code")
	c.SetParamNames("code")
	c.SetParamValues(num)
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.GetPromoByCode(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetDetailID(t *testing.T) {
	var mockPromo models.PromoDto
	err := faker.FakeData(&mockPromo)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := mockPromo.Id

	mockUCase.On("GetDetail", mock.Anything, num,mock.AnythingOfType("string")).Return(&mockPromo, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/admin/promo/"+num, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("admin/promo/:id")
	c.SetParamNames("id")
	c.SetParamValues(num)
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.GetDetailID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetDetailPromoIDErrorNotFound(t *testing.T) {
	var mockPromo models.PromoDto
	err := faker.FakeData(&mockPromo)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := mockPromo.Id

	mockUCase.On("GetDetail", mock.Anything, num,mock.AnythingOfType("string")).Return(nil, models.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/admin/promo/"+num, strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("admin/promo/:id")
	c.SetParamNames("id")
	c.SetParamValues(num)
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.GetDetailID(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreatePromo(t *testing.T) {

	tempMockPromo := mockPromo
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockPromo)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("models.NewCommandPromo"), mock.AnythingOfType("string")).Return(&mockPromo, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockPromo.PromoImage, nil)
	//var param = url.Values{}
	//param.Set("Promo_name", tempMockPromo.PromoName)
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
	writer.WriteField("max_usage", strconv.Itoa(tempMockPromo.MaxUsage))
	writer.WriteField("currency", strconv.Itoa(tempMockPromo.Currency))
	writer.WriteField("promo_value", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("promo_type", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("production_capacity", fmt.Sprint(tempMockPromo.ProductionCapacity))
	writer.WriteField("promo_product_type", fmt.Sprint(*tempMockPromo.PromoProductType))
	writer.WriteField("is_any_trip_period", fmt.Sprint(tempMockPromo.IsAnyTripPeriod))
	writer.WriteField("max_discount", fmt.Sprint(tempMockPromo.MaxDiscount))
	writer.WriteField("merchant_id", merchantString)
	writer.WriteField("id", tempMockPromo.Id)
	writer.WriteField("promo_code", tempMockPromo.PromoCode)
	writer.WriteField("Promo_name", tempMockPromo.PromoName)
	writer.WriteField("promo_desc", tempMockPromo.PromoDesc)
	writer.WriteField("start_date", tempMockPromo.StartDate)
	writer.WriteField("end_date", tempMockPromo.EndDate)
	writer.WriteField("start_trip_period", tempMockPromo.StartTripPeriod)
	writer.WriteField("end_trip_period", tempMockPromo.EndTripPeriod)
	writer.WriteField("disclaimer", tempMockPromo.Disclaimer)
	writer.WriteField("term_condition", tempMockPromo.TermCondition)
	writer.WriteField("how_to_get", tempMockPromo.HowToGet)
	writer.WriteField("how_to_use", tempMockPromo.HowToUse)
	part, err := writer.CreateFormFile("promo_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/admin/promo", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/admin/promo")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
		IsUsecase:    mockIsUsecase,
	}
	err = handler.CreatePromo(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreatePromoWithoutToken(t *testing.T) {

	tempMockPromo := mockPromo
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockPromo)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("models.NewCommandPromo"), mock.AnythingOfType("string")).Return(nil, models.ErrUnAuthorize)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockPromo.PromoImage, nil)
	//var param = url.Values{}
	//param.Set("Promo_name", tempMockPromo.PromoName)
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
	writer.WriteField("max_usage", strconv.Itoa(tempMockPromo.MaxUsage))
	writer.WriteField("currency", strconv.Itoa(tempMockPromo.Currency))
	writer.WriteField("promo_value", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("promo_type", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("production_capacity", fmt.Sprint(tempMockPromo.ProductionCapacity))
	writer.WriteField("promo_product_type", fmt.Sprint(*tempMockPromo.PromoProductType))
	writer.WriteField("is_any_trip_period", fmt.Sprint(tempMockPromo.IsAnyTripPeriod))
	writer.WriteField("max_discount", fmt.Sprint(tempMockPromo.MaxDiscount))
	writer.WriteField("merchant_id", merchantString)
	writer.WriteField("id", tempMockPromo.Id)
	writer.WriteField("promo_code", tempMockPromo.PromoCode)
	writer.WriteField("Promo_name", tempMockPromo.PromoName)
	writer.WriteField("promo_desc", tempMockPromo.PromoDesc)
	writer.WriteField("start_date", tempMockPromo.StartDate)
	writer.WriteField("end_date", tempMockPromo.EndDate)
	writer.WriteField("start_trip_period", tempMockPromo.StartTripPeriod)
	writer.WriteField("end_trip_period", tempMockPromo.EndTripPeriod)
	writer.WriteField("disclaimer", tempMockPromo.Disclaimer)
	writer.WriteField("term_condition", tempMockPromo.TermCondition)
	writer.WriteField("how_to_get", tempMockPromo.HowToGet)
	writer.WriteField("how_to_use", tempMockPromo.HowToUse)
	part, err := writer.CreateFormFile("promo_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()



	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/admin/promo", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/admin/promo")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
		IsUsecase:    mockIsUsecase,
	}
	err = handler.CreatePromo(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreatePromoConflict(t *testing.T) {

	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Promo",
	//}
	tempMockPromo := mockPromo
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockPromo)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("models.NewCommandPromo"), mock.AnythingOfType("string")).Return(nil, models.ErrConflict)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockPromo.PromoImage, nil)
	//var param = url.Values{}
	//param.Set("Promo_name", tempMockPromo.PromoName)
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
	writer.WriteField("max_usage", strconv.Itoa(tempMockPromo.MaxUsage))
	writer.WriteField("currency", strconv.Itoa(tempMockPromo.Currency))
	writer.WriteField("promo_value", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("promo_type", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("production_capacity", fmt.Sprint(tempMockPromo.ProductionCapacity))
	writer.WriteField("promo_product_type", fmt.Sprint(*tempMockPromo.PromoProductType))
	writer.WriteField("is_any_trip_period", fmt.Sprint(tempMockPromo.IsAnyTripPeriod))
	writer.WriteField("max_discount", fmt.Sprint(tempMockPromo.MaxDiscount))
	writer.WriteField("merchant_id", merchantString)
	writer.WriteField("id", tempMockPromo.Id)
	writer.WriteField("promo_code", tempMockPromo.PromoCode)
	writer.WriteField("Promo_name", tempMockPromo.PromoName)
	writer.WriteField("promo_desc", tempMockPromo.PromoDesc)
	writer.WriteField("start_date", tempMockPromo.StartDate)
	writer.WriteField("end_date", tempMockPromo.EndDate)
	writer.WriteField("start_trip_period", tempMockPromo.StartTripPeriod)
	writer.WriteField("end_trip_period", tempMockPromo.EndTripPeriod)
	writer.WriteField("disclaimer", tempMockPromo.Disclaimer)
	writer.WriteField("term_condition", tempMockPromo.TermCondition)
	writer.WriteField("how_to_get", tempMockPromo.HowToGet)
	writer.WriteField("how_to_use", tempMockPromo.HowToUse)
	part, err := writer.CreateFormFile("promo_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()


	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/Promo", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/Promo")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
		IsUsecase:    mockIsUsecase,
	}
	err = handler.CreatePromo(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdatePromo(t *testing.T) {

	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Promo",
	//}
	tempMockPromo := mockPromo
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockPromo)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockPromo.Id
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("models.NewCommandPromo"), mock.AnythingOfType("string")).Return(&mockPromo, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockPromo.PromoImage, nil)

	var param = url.Values{}
	param.Set("Promo_name", tempMockPromo.PromoName)
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
	writer.WriteField("max_usage", strconv.Itoa(tempMockPromo.MaxUsage))
	writer.WriteField("currency", strconv.Itoa(tempMockPromo.Currency))
	writer.WriteField("promo_value", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("promo_type", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("production_capacity", fmt.Sprint(tempMockPromo.ProductionCapacity))
	writer.WriteField("promo_product_type", fmt.Sprint(*tempMockPromo.PromoProductType))
	writer.WriteField("is_any_trip_period", fmt.Sprint(tempMockPromo.IsAnyTripPeriod))
	writer.WriteField("max_discount", fmt.Sprint(tempMockPromo.MaxDiscount))
	writer.WriteField("merchant_id", merchantString)
	writer.WriteField("id", tempMockPromo.Id)
	writer.WriteField("promo_code", tempMockPromo.PromoCode)
	writer.WriteField("Promo_name", tempMockPromo.PromoName)
	writer.WriteField("promo_desc", tempMockPromo.PromoDesc)
	writer.WriteField("start_date", tempMockPromo.StartDate)
	writer.WriteField("end_date", tempMockPromo.EndDate)
	writer.WriteField("start_trip_period", tempMockPromo.StartTripPeriod)
	writer.WriteField("end_trip_period", tempMockPromo.EndTripPeriod)
	writer.WriteField("disclaimer", tempMockPromo.Disclaimer)
	writer.WriteField("term_condition", tempMockPromo.TermCondition)
	writer.WriteField("how_to_get", tempMockPromo.HowToGet)
	writer.WriteField("how_to_use", tempMockPromo.HowToUse)
	part, err := writer.CreateFormFile("promo_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/admin/promo/"+id, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/promo/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
		IsUsecase:    mockIsUsecase,
	}
	err = handler.UpdatePromo(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdatePromoWithoutToken(t *testing.T) {
	tempMockPromo := mockPromo
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockPromo)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockPromo.Id
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("models.NewCommandPromo"), mock.AnythingOfType("string")).Return(&mockPromo, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockPromo.PromoImage, nil)

	var param = url.Values{}
	param.Set("Promo_name", tempMockPromo.PromoName)
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
	writer.WriteField("max_usage", strconv.Itoa(tempMockPromo.MaxUsage))
	writer.WriteField("currency", strconv.Itoa(tempMockPromo.Currency))
	writer.WriteField("promo_value", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("promo_type", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("production_capacity", fmt.Sprint(tempMockPromo.ProductionCapacity))
	writer.WriteField("promo_product_type", fmt.Sprint(*tempMockPromo.PromoProductType))
	writer.WriteField("is_any_trip_period", fmt.Sprint(tempMockPromo.IsAnyTripPeriod))
	writer.WriteField("max_discount", fmt.Sprint(tempMockPromo.MaxDiscount))
	writer.WriteField("merchant_id", merchantString)
	writer.WriteField("id", tempMockPromo.Id)
	writer.WriteField("promo_code", tempMockPromo.PromoCode)
	writer.WriteField("Promo_name", tempMockPromo.PromoName)
	writer.WriteField("promo_desc", tempMockPromo.PromoDesc)
	writer.WriteField("start_date", tempMockPromo.StartDate)
	writer.WriteField("end_date", tempMockPromo.EndDate)
	writer.WriteField("start_trip_period", tempMockPromo.StartTripPeriod)
	writer.WriteField("end_trip_period", tempMockPromo.EndTripPeriod)
	writer.WriteField("disclaimer", tempMockPromo.Disclaimer)
	writer.WriteField("term_condition", tempMockPromo.TermCondition)
	writer.WriteField("how_to_get", tempMockPromo.HowToGet)
	writer.WriteField("how_to_use", tempMockPromo.HowToUse)
	part, err := writer.CreateFormFile("promo_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/admin/promo/"+id, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/promo/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
		IsUsecase:    mockIsUsecase,
	}
	err = handler.UpdatePromo(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdatePromoBadParam(t *testing.T) {
	tempMockPromo := mockPromo
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockPromo)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockPromo.Id
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("models.NewCommandPromo"), mock.AnythingOfType("string")).Return(nil, models.ErrBadParamInput)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockPromo.PromoImage, nil)

	var param = url.Values{}
	param.Set("Promo_name", tempMockPromo.PromoName)
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
	writer.WriteField("max_usage", strconv.Itoa(tempMockPromo.MaxUsage))
	writer.WriteField("currency", strconv.Itoa(tempMockPromo.Currency))
	writer.WriteField("promo_value", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("promo_type", fmt.Sprint(tempMockPromo.PromoValue))
	writer.WriteField("production_capacity", fmt.Sprint(tempMockPromo.ProductionCapacity))
	writer.WriteField("promo_product_type", fmt.Sprint(*tempMockPromo.PromoProductType))
	writer.WriteField("is_any_trip_period", fmt.Sprint(tempMockPromo.IsAnyTripPeriod))
	writer.WriteField("max_discount", fmt.Sprint(tempMockPromo.MaxDiscount))
	writer.WriteField("merchant_id", merchantString)
	writer.WriteField("id", tempMockPromo.Id)
	writer.WriteField("promo_code", tempMockPromo.PromoCode)
	writer.WriteField("Promo_name", tempMockPromo.PromoName)
	writer.WriteField("promo_desc", tempMockPromo.PromoDesc)
	writer.WriteField("start_date", tempMockPromo.StartDate)
	writer.WriteField("end_date", tempMockPromo.EndDate)
	writer.WriteField("start_trip_period", tempMockPromo.StartTripPeriod)
	writer.WriteField("end_trip_period", tempMockPromo.EndTripPeriod)
	writer.WriteField("disclaimer", tempMockPromo.Disclaimer)
	writer.WriteField("term_condition", tempMockPromo.TermCondition)
	writer.WriteField("how_to_get", tempMockPromo.HowToGet)
	writer.WriteField("how_to_use", tempMockPromo.HowToUse)
	part, err := writer.CreateFormFile("promo_image", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()


	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/admin/promo/"+id, payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/promo/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
		IsUsecase:    mockIsUsecase,
	}
	err = handler.UpdatePromo(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeletePromo(t *testing.T) {

	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Promo",
	}

	tempMockPromo := mockPromo
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockPromo)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockPromo.Id
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/admin/promo/"+id, nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/promo/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.Delete(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeletePromoWithoutToken(t *testing.T) {

	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Promo",
	}
	tempMockPromo := mockPromo
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockPromo)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockPromo.Id
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/admin/promo/"+id, nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/promo/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.Delete(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeletePromoErrorInternalServer(t *testing.T) {

	tempMockPromo := mockPromo
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockPromo)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := tempMockPromo.Id
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil, models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/admin/promo/"+id, nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/admin/promo/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(id)
	handler := PromoHttp.PromoHandler{
		PromoUsecase: mockUCase,
	}
	err = handler.Delete(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
