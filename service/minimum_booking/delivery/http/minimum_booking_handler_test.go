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
	"github.com/service/minimum_booking/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	MinimumBookingHttp "github.com/service/minimum_booking/delivery/http"
)

func TestList(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockMinimumBookingPagination = models.MinimumBookingDtoWithPagination{}

	err := faker.FakeData(&mockMinimumBookingPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/minimum_booking?page="+strconv.Itoa(mockMinimumBookingPagination.Meta.Page)+"&size="+strconv.Itoa(mockMinimumBookingPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(&mockMinimumBookingPagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := MinimumBookingHttp.MinimumBookingHandler{
		MinimumBookingUsecase: mockUCase,
	}
	err = handler.GetAllMinimumBooking(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestListErrorInternalServer(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockMinimumBookingPagination = models.MinimumBookingDtoWithPagination{}

	err := faker.FakeData(&mockMinimumBookingPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/master/minimum_booking?page="+strconv.Itoa(mockMinimumBookingPagination.Meta.Page)+"&size="+strconv.Itoa(mockMinimumBookingPagination.Meta.RecordPerPage), strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetAll", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("int"), mock.AnythingOfType("string")).Return(nil, errors.New("Error Internal Server"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := MinimumBookingHttp.MinimumBookingHandler{
		MinimumBookingUsecase: mockUCase,
	}
	err = handler.GetAllMinimumBooking(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
