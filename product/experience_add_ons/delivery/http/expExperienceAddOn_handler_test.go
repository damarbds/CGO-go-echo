package http_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bxcodec/faker"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/product/experience_add_ons/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	ExperienceAddOnDtoHttp "github.com/product/experience_add_ons/delivery/http"
)

func TestGetallexperienceAddOns(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	var mockExperienceAddOnDtoPagination []*models.ExperienceAddOnDto

	err := faker.FakeData(&mockExperienceAddOnDtoPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "product/product-add-ons?exp_id=asdasdads", strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetByExpId", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(mockExperienceAddOnDtoPagination, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := ExperienceAddOnDtoHttp.Experience_add_onsHandler{
		Experience_add_onsUsecase: mockUCase,
	}
	err = handler.GetallexperienceAddOns(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}

func TestGetallexperienceAddOnsErrorInternalServer(t *testing.T) {
	mockUCase := new(mocks.Usecase)

	var mockExperienceAddOnDtoPagination []*models.ExperienceAddOnDto

	err := faker.FakeData(&mockExperienceAddOnDtoPagination)
	assert.NoError(t, err)

	e := echo.New()
	req, err := http.NewRequest(echo.GET, "product/product-add-ons?exp_id=asdasdads", strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("GetByExpId", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string")).
		Return(nil, errors.New("UnExpected"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := ExperienceAddOnDtoHttp.Experience_add_onsHandler{
		Experience_add_onsUsecase: mockUCase,
	}
	err = handler.GetallexperienceAddOns(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
