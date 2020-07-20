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
	"github.com/service/rule/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	RuleHttp "github.com/service/rule/delivery/http"
)

var (
	imagePath = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Rule/8941695193938718058.jpg"
	mockRule  = models.Rule{
		Id:           1,
		CreatedBy:    "test",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		RuleName:     "Test Rule 2",
		RuleIcon:     imagePath,
	}
	tempMockRule = models.NewCommandRule{
		Id:          1,
		RuleName:    mockRule.RuleName,
		RuleIcon:    mockRule.RuleIcon,
	}
)

func TestList(t *testing.T) {
	var mockRule models.RuleDto
	err := faker.FakeData(&mockRule)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListRule []*models.RuleDto
	mockListRule = append(mockListRule, &mockRule)

	mockUCase.On("List", mock.Anything).Return(mockListRule, nil)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/rule", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Request().Header.Add("Authorization", token)

	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.List(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestListWithoutToken(t *testing.T) {
	var mockRule models.RuleDto
	err := faker.FakeData(&mockRule)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListRule []*models.RuleDto
	mockListRule = append(mockListRule, &mockRule)

	mockUCase.On("List", mock.Anything).Return(mockListRule, nil)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/rule", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	//c.Request().Header.Add("Authorization",token)

	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.List(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestListInvalidToken(t *testing.T) {
	var mockRule models.RuleDto
	err := faker.FakeData(&mockRule)
	assert.NoError(t, err)
	mockUCase := new(mocks.Usecase)
	var mockListRule []*models.RuleDto
	mockListRule = append(mockListRule, &mockRule)

	mockUCase.On("List", mock.Anything).Return(nil, models.ErrUnAuthorize)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/rule", strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.Request().Header.Add("Authorization", token)

	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.List(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetAllRule(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockRulePagination = models.RuleDtoWithPagination{}

	err := faker.FakeData(&mockRulePagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/rule?page="+strconv.Itoa(mockRulePagination.Meta.Page)+"&size="+strconv.Itoa(mockRulePagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(&mockRulePagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.GetAllRule(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetAllRuleErrorInternalServer(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockRulePagination = models.RuleDtoWithPagination{}

	err := faker.FakeData(&mockRulePagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/rule?page="+strconv.Itoa(mockRulePagination.Meta.Page)+"&size="+strconv.Itoa(mockRulePagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, errors.New("Internal server Error"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.GetAllRule(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestGetDetailRuleID(t *testing.T) {
	var mockRule models.RuleDto
	err := faker.FakeData(&mockRule)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockRule.Id)

	mockUCase.On("GetById", mock.Anything, int(num)).Return(&mockRule, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/rule/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/rule/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.GetDetailRuleID(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetDetailRuleIDErrorNotFound(t *testing.T) {
	var mockRule models.RuleDto
	err := faker.FakeData(&mockRule)
	assert.NoError(t, err)

	mockUCase := new(mocks.Usecase)

	num := int(mockRule.Id)

	mockUCase.On("GetById", mock.Anything, int(num)).Return(nil, models.ErrNotFound)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/rule/"+strconv.Itoa(num), strings.NewReader(""))
	assert.NoError(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("master/rule/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(num))
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.GetDetailRuleID(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateRule(t *testing.T) {
	mockRule := tempMockRule
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Rule",
	}
	tempMockRule := mockRule
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockRule)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandRule"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockRule.RuleIcon, nil)
	var param = url.Values{}
	param.Set("rule_name", tempMockRule.RuleName)
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
	writer.WriteField("id", strconv.Itoa(tempMockRule.Id))
	writer.WriteField("rule_name", tempMockRule.RuleName)
	part, err := writer.CreateFormFile("rule_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/rule", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/rule")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
		IsUsecase:   mockIsUsecase,
	}
	err = handler.CreateRule(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestCreateRuleWithoutToken(t *testing.T) {
	mockRule := tempMockRule
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Rule",
	}
	tempMockRule := mockRule
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockRule)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandRule"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockRule.RuleIcon, nil)
	var param = url.Values{}
	param.Set("rule_name", tempMockRule.RuleName)
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
	writer.WriteField("id", strconv.Itoa(tempMockRule.Id))
	writer.WriteField("rule_name", tempMockRule.RuleName)
	part, err := writer.CreateFormFile("rule_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/rule", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/rule")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
		IsUsecase:   mockIsUsecase,
	}
	err = handler.CreateRule(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestCreateRuleConflict(t *testing.T) {
	mockRule := tempMockRule
	//mockReponse := &models.ResponseDelete{
	//	Id:      "1",
	//	Message: "Success Create Rule",
	//}
	tempMockRule := mockRule
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockRule)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"

	mockUCase.On("Create", mock.Anything, mock.AnythingOfType("*models.NewCommandRule"), mock.AnythingOfType("string")).Return(nil, models.ErrConflict)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockRule.RuleIcon, nil)
	var param = url.Values{}
	param.Set("rule_name", tempMockRule.RuleName)
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
	writer.WriteField("id", strconv.Itoa(tempMockRule.Id))
	writer.WriteField("rule_name", tempMockRule.RuleName)
	part, err := writer.CreateFormFile("rule_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/master/rule", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/master/rule")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
		IsUsecase:   mockIsUsecase,
	}
	err = handler.CreateRule(c)
	//require.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateRule(t *testing.T) {
	mockRule := tempMockRule
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Rule",
	}
	tempMockRule := mockRule
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockRule)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockRule.Id)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandRule"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockRule.RuleIcon, nil)

	var param = url.Values{}
	param.Set("rule_name", tempMockRule.RuleName)
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
	writer.WriteField("id", strconv.Itoa(tempMockRule.Id))
	writer.WriteField("rule_name", tempMockRule.RuleName)
	part, err := writer.CreateFormFile("rule_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/rule/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/rule/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
		IsUsecase:   mockIsUsecase,
	}
	err = handler.UpdateRule(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestUpdateRuleWithoutToken(t *testing.T) {
	mockRule := tempMockRule
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Rule",
	}
	tempMockRule := mockRule
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockRule)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockRule.Id)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandRule"), mock.AnythingOfType("string")).Return(mockReponse, nil)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockRule.RuleIcon, nil)

	var param = url.Values{}
	param.Set("rule_name", tempMockRule.RuleName)
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
	writer.WriteField("id", strconv.Itoa(tempMockRule.Id))
	writer.WriteField("rule_name", tempMockRule.RuleName)
	part, err := writer.CreateFormFile("rule_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/rule/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/rule/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
		IsUsecase:   mockIsUsecase,
	}
	err = handler.UpdateRule(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestUpdateRuleBadParam(t *testing.T) {
	mockRule := tempMockRule

	tempMockRule := mockRule
	mockUCase := new(mocks.Usecase)
	mockIsUsecase := new(_mockIdentityserver.Usecase)
	//j, err := json.Marshal(tempMockRule)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockRule.Id)
	mockUCase.On("Update", mock.Anything, mock.AnythingOfType("*models.NewCommandRule"), mock.AnythingOfType("string")).Return(nil, models.ErrBadParamInput)
	mockIsUsecase.On("UploadFileToBlob", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mockRule.RuleIcon, nil)

	var param = url.Values{}
	param.Set("rule_name", tempMockRule.RuleName)
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
	writer.WriteField("id", strconv.Itoa(tempMockRule.Id))
	writer.WriteField("rule_name", tempMockRule.RuleName)
	part, err := writer.CreateFormFile("rule_icon", filepath.Base(path))
	if err != nil {
		writer.Close()
		t.Error(err)
	}
	io.Copy(part, file)
	writer.Close()

	e := echo.New()
	req, err := http.NewRequest(echo.PUT, "/master/rule/"+strconv.Itoa(id), payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	assert.NoError(t, err)
	//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/rule/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
		IsUsecase:   mockIsUsecase,
	}
	err = handler.UpdateRule(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteRule(t *testing.T) {
	mockRule := tempMockRule
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Rule",
	}
	tempMockRule := mockRule
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockRule)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockRule.Id)
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/rule/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/rule/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.DeleteRule(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestDeleteRuleWithoutToken(t *testing.T) {
	mockRule := tempMockRule
	mockReponse := &models.ResponseDelete{
		Id:      "1",
		Message: "Success Create Rule",
	}
	tempMockRule := mockRule
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockRule)
	//assert.NoError(t, err)
	//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockRule.Id)
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(mockReponse, nil)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/rule/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/rule/:id")
	//c.Request().Header.Add("Authorization",token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.DeleteRule(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	//mockUCase.AssertExpectations(t)
}

func TestDeleteRuleErrorInternalServer(t *testing.T) {
	mockRule := tempMockRule
	tempMockRule := mockRule
	mockUCase := new(mocks.Usecase)

	//j, err := json.Marshal(tempMockRule)
	//assert.NoError(t, err)
	token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
	id := int(tempMockRule.Id)
	mockUCase.On("Delete", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil, models.ErrInternalServerError)

	e := echo.New()
	req, err := http.NewRequest(echo.DELETE, "/master/rule/"+strconv.Itoa(id), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/master/rule/:id")
	c.Request().Header.Add("Authorization", token)
	c.Request().ParseForm()
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(id))
	handler := RuleHttp.RuleHandler{
		RuleUsecase: mockUCase,
	}
	err = handler.DeleteRule(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
