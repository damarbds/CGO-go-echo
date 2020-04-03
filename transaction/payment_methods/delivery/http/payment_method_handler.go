package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	pm "github.com/transaction/payment_methods"
	"net/http"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

// paymentMethodHandler represent the httphandler for payment method
type paymentMethodHandler struct {
	paymentMethodUsecase pm.Usecase
}

/// NewPaymentMethodHandler will initialize the promos/ resources endpoint
func NewPaymentMethodHandler(e *echo.Echo, us pm.Usecase) {
	handler := &paymentMethodHandler{
		paymentMethodUsecase: us,
	}
	e.GET("transaction/payment-method", handler.GetPaymentMethods)
}

// GetPaymentMethods will get article by given id
func (p *paymentMethodHandler) GetPaymentMethods(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

		art, err := p.paymentMethodUsecase.Fetch(ctx)
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

