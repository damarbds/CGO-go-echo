package http_test

import (
	"errors"
	"github.com/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	currencyHttp "github.com/misc/currency/delivery/http"
	"github.com/misc/currency/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestExchangeRateIDRToUSD(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	mockcurrencyication := models.CurrencyExChangeRate{
		Date:  "2020-09-17",
		Base:  "USD",
		Rates: models.Rates{
			IDR: 12312312,
			USD: 3242,
		},
	}

	//err := faker.FakeData(&mockcurrencyicationPagination)
	//assert.NoError(t, err)
	from := "IDR"
	to := "USD"
	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/misc/exchange-rate?from="+from+"&to="+to, strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("ExchangeRatesApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockcurrencyication, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := currencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.ExchangeRate(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}
func TestExchangeRateUSDToIDR(t *testing.T) {


	mockUCase := new(mocks.Usecase)

	mockcurrencyication := models.CurrencyExChangeRate{
		Date:  "2020-09-17",
		Base:  "USD",
		Rates: models.Rates{
			IDR: 12312312,
			USD: 3242,
		},
	}


	//err := faker.FakeData(&mockcurrencyicationPagination)
	//assert.NoError(t, err)
	from := "USD"
	to := "IDR"
	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/misc/exchange-rate?from="+from+"&to="+to, strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("ExchangeRatesApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(mockcurrencyication, nil)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := currencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.ExchangeRate(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}
func TestExchangeRateErrorBadRequest(t *testing.T) {


	mockUCase := new(mocks.Usecase)


	//err := faker.FakeData(&mockcurrencyicationPagination)
	//assert.NoError(t, err)
	from := "USD"
	to := "IDR"
	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/misc/exchange-rate?from="+from+"&to="+to, strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("ExchangeRatesApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil, models.ErrBadParamInput)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := currencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.ExchangeRate(c)	//require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}
func TestExchangeRateErrorInternalServer(t *testing.T) {


	mockUCase := new(mocks.Usecase)

	//err := faker.FakeData(&mockcurrencyicationPagination)
	//assert.NoError(t, err)
	from := "USD"
	to := "IDR"
	e := echo.New()
	req, err := http.NewRequest(echo.GET, "/misc/exchange-rate?from="+from+"&to="+to, strings.NewReader(""))
	assert.NoError(t, err)
	mockUCase.On("ExchangeRatesApi", mock.Anything, mock.AnythingOfType("string"),mock.AnythingOfType("string")).Return(nil, errors.New("UnExpected"))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := currencyHttp.CurrencyHandler{
		CurrencyUsecase: mockUCase,
	}
	err = handler.ExchangeRate(c)
	//require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	//mockUCase.AssertExpectations(t)
}
func TestCreateExChange(t *testing.T) {

	mockUCase := new(mocks.Usecase)

	mockUCase.On("Insert", mock.Anything).Return(nil)

	e := echo.New()
	req, err := http.NewRequest(echo.POST, "/misc/exchange-rate", nil)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/misc/exchange-rate")
	c.Request().ParseForm()
	handler := currencyHttp.CurrencyHandler{
		CurrencyUsecase:mockUCase,
	}
	err = handler.CreateExChange(c)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockUCase.AssertExpectations(t)
}
func TestCreateExChangeConflict(t *testing.T) {
	mockUCase := new(mocks.Usecase)

	mockUCase.On("Insert", mock.Anything).Return(models.ErrConflict)

	e := echo.New()
	req, _ := http.NewRequest(echo.POST, "/misc/exchange-rate", nil)


	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetPath("/misc/exchange-rate")
	c.Request().ParseForm()
	handler := currencyHttp.CurrencyHandler{
		CurrencyUsecase:mockUCase,
	}
	_ = handler.CreateExChange(c)
	//assert.Error(t, err)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	//mockUCase.AssertExpectations(t)
}
