package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/exclusion_service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ExclusionServiceHttp "github.com/service/exclusion_service/delivery/http"
)

func TestList(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockExclusionServicePagination = models.ExclusionServiceWithPagination{}

	err := faker.FakeData(&mockExclusionServicePagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/exclusion_service?page="+strconv.Itoa(mockExclusionServicePagination.Meta.Page)+"&size="+strconv.Itoa(mockExclusionServicePagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("List", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(&mockExclusionServicePagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := ExclusionServiceHttp.ExclusionServicesHandler{
		ExclusionServicesUsecase: mockUCase,
	}
	err = handler.List(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestListErrorInternalServer(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockExclusionServicePagination = models.ExclusionServiceWithPagination{}

	err := faker.FakeData(&mockExclusionServicePagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/service/exclusion_service?page="+strconv.Itoa(mockExclusionServicePagination.Meta.Page)+"&size="+strconv.Itoa(mockExclusionServicePagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("List", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil, errors.New("Error Internal Server"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := ExclusionServiceHttp.ExclusionServicesHandler{
		ExclusionServicesUsecase: mockUCase,
	}
	err = handler.List(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
