package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/exclusion_service"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// promoHandler  represent the httphandler for promo
type exclusionServicesHandler struct {
	exclusionServicesUsecase exclusion_service.Usecase
}

// NewpromoHandler will initialize the promos/ resources endpoint
func NewExclusionServicesHandler(e *echo.Echo, exclusionServicesUsecase exclusion_service.Usecase) {
	handler := &exclusionServicesHandler{
		exclusionServicesUsecase : exclusionServicesUsecase,
	}
	e.GET("master/exclusion_service", handler.List)
}

func (a *exclusionServicesHandler) List(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	//if token == "" {
	//	return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	//}

	qpage := c.QueryParam("page")
	qperPage := c.QueryParam("size")
	var limit = 20
	var page = 1
	var offset = 0

	page, _ = strconv.Atoi(qpage)
	limit, _ = strconv.Atoi(qperPage)
	offset = (page - 1) * limit

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	art, err := a.exclusionServicesUsecase.List(ctx, page, limit,offset,token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, art)

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