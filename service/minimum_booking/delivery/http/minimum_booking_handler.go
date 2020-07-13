package http

import (
	"github.com/service/minimum_booking"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// harborsHandler  represent the httphandler for harbors
type MinimumBookingHandler struct {
	MinimumBookingUsecase minimum_booking.Usecase
}

// NewharborsHandler will initialize the harborss/ resources endpoint
func NewminimumBookingHandler(e *echo.Echo, minimumBookingUsecase minimum_booking.Usecase) {
	handler := &MinimumBookingHandler{
		MinimumBookingUsecase:minimumBookingUsecase,
	}
	e.GET("master/minimum_booking", handler.GetAllMinimumBooking)
}

func (a *MinimumBookingHandler) GetAllMinimumBooking(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")

	var limit = 20
	var page = 1
	var offset = 0

	page, _ = strconv.Atoi(qpage)
	limit, _ = strconv.Atoi(qsize)
	offset = (page - 1) * limit
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	art, err := a.MinimumBookingUsecase.GetAll(ctx, page,limit,offset)
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