package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/misc/currency"

	"github.com/models"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type currencyHandler struct {
	currencyUsecase currency.Usecase
}

func NewCurrencyHandler(e *echo.Echo, cu currency.Usecase) {
	handler := &currencyHandler{currencyUsecase: cu}
	e.GET("/misc/exchange-rate", handler.ExchangeRate)
}

func (cu *currencyHandler) ExchangeRate(c echo.Context) error {
	from := c.QueryParam("from")
	to := c.QueryParam("to")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	//exchangeKey := from + "_" + to
	res, err := cu.currencyUsecase.ExchangeRatesApi(ctx, from,to)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	var rates float64
	if to == "IDR" {
		rates = res.Rates.IDR
	}else if to == "USD" {
		rates = res.Rates.USD
	}
	result := map[string]interface{}{
		"rates": rates,
	}

	return c.JSON(http.StatusOK, result)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrUnAuthorize:
		return http.StatusUnauthorized
	case models.ErrConflict:
		return http.StatusBadRequest
	case models.ErrBadParamInput:
		return http.StatusBadRequest
	case models.ValidationBookedDate:
		return http.StatusBadRequest
	case models.ValidationStatus:
		return http.StatusBadRequest
	case models.ValidationBookedBy:
		return http.StatusBadRequest
	case models.ValidationExpId:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
