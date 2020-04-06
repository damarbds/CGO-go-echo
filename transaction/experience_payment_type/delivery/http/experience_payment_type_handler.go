package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"github.com/transaction/experience_payment_type"
	"net/http"
	"strconv"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type expPaymentTypeHandler struct {
	experiencePaymentTypeUcase experience_payment_type.Usecase
}

func NewexpPaymentTypeHandlerHandler(e *echo.Echo, pus experience_payment_type.Usecase) {
	handler := &expPaymentTypeHandler{
		experiencePaymentTypeUcase: pus,
	}
	e.GET("/transaction/payments-types", handler.GetAllPaymentTypes)
}
// GetByID will get article by given id
func (a *expPaymentTypeHandler) GetAllPaymentTypes(c echo.Context) error {
	qpage := c.QueryParam("page")
	qsize := c.QueryParam("size")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	if qpage != "" && qsize != ""{
		page , _:= strconv.Atoi(qpage)
		size , _:= strconv.Atoi(qsize)
		art, err := a.experiencePaymentTypeUcase.GetAll(ctx,&page,&size)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	}else {
		art, err := a.experiencePaymentTypeUcase.GetAll(ctx,nil,nil)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, art)
	}

	return nil
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
