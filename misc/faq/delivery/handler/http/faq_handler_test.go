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
	"github.com/misc/faq/mocks"
	"github.com/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	FAQHttp "github.com/misc/faq/delivery/handler/http"
)

func TestGetByType(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockFAQPagination []*models.FAQDto

	err := faker.FakeData(&mockFAQPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "misc/faq?type="+ strconv.Itoa(1), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetByType", mock.Anything, mock.AnythingOfType("int")).Return(mockFAQPagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := FAQHttp.FaqHandler{
		FaqUsecase: mockUCase,
	}
	err = handler.GetByType(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetByTypeErrorInternalServer(t *testing.T) {
	mockUCase := new(mocks.Usecase)

	var mockFAQPagination []*models.FAQDto

	err := faker.FakeData(&mockFAQPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "misc/faq?type="+ strconv.Itoa(1), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetByType", mock.Anything, mock.AnythingOfType("int")).Return(nil, errors.New("UnExpected"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := FAQHttp.FaqHandler{
		FaqUsecase: mockUCase,
	}
	err = handler.GetByType(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
