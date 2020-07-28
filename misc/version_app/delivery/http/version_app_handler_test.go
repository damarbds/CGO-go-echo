package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/misc/version_app/mocks"
	"github.com/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	VersionAppHttp "github.com/misc/version_app/delivery/http"
)

func TestGet(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockVersionAppPagination []*models.VersionApp

	err := faker.FakeData(&mockVersionAppPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "version", strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAllVersion", mock.Anything, mock.AnythingOfType("int")).Return(mockVersionAppPagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := VersionAppHttp.VersionAPPHandler{
		VersionAPPUsecase: mockUCase,
	}
	err = handler.Get(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetErrorInternalServer(t *testing.T) {
	mockUCase := new(mocks.Usecase)

	var mockVersionAppPagination []*models.VersionApp

	err := faker.FakeData(&mockVersionAppPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "version", strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAllVersion", mock.Anything, mock.AnythingOfType("int")).Return(nil, errors.New("UnExpected"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := VersionAppHttp.VersionAPPHandler{
		VersionAPPUsecase: mockUCase,
	}
	err = handler.Get(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
