package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"github.com/transactions/transaction"
)

type ResponseError struct {
	Message string `json:"message"`
}

type transactionHandler struct {
	TransUsecase transaction.Usecase
}

func NewTransactionHandler(e *echo.Echo, us transaction.Usecase) {
	handler := &transactionHandler{
		TransUsecase: us,
	}
	e.GET("/transaction/count-success", handler.CountSuccess)
}

func (t *transactionHandler) CountSuccess(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	if token == "" {
		return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := t.TransUsecase.CountSuccess(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
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
	default:
		return http.StatusInternalServerError
	}
}
