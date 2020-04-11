package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/transportation"
	"github.com/sirupsen/logrus"
	"net/http"
)

type ResponseError struct {
	Message string `json:"message"`
}

type transportationHandler struct {
	TransUsecase transportation.Usecase
}

func NewTransportationHandler(e *echo.Echo, us transportation.Usecase) {
	handler := &transportationHandler{
		TransUsecase: us,
	}
	e.GET("/service/transportation/time-options", handler.List)
}

func (t *transportationHandler) List(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	result, err := t.TransUsecase.List(ctx)
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