package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/facilities"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type facilityHandler struct {
	facilityUsecase facilities.Usecase
}

func NewFacilityHandler(e *echo.Echo, us facilities.Usecase) {
	handler := &facilityHandler{
		facilityUsecase: us,
	}
	e.GET("service/facilities", handler.List)
}

func (f *facilityHandler) List(c echo.Context) error {
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

	list, err := f.facilityUsecase.List(ctx)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, list)
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
