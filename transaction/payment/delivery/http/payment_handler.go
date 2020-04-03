package http

import (
	"context"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"github.com/transaction/payment"
	"net/http"
)

// ResponseError represent the reseponse error struct
type ResponseError struct {
	Message string `json:"message"`
}

type paymentHandler struct {
	paymentUsecase payment.Usecase
}

func NewPaymentHandler(e *echo.Echo, pus payment.Usecase) {
	handler := &paymentHandler{
		paymentUsecase: pus,
	}
	e.POST("/transaction/payments", handler.CreatePayment)
}

func (p *paymentHandler) CreatePayment(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	t := new(models.Transaction)
	if err := c.Bind(t); err != nil {
		return models.ErrBadParamInput
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := p.paymentUsecase.Insert(ctx, t, token)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	response := &models.PaymentTransaction{
		Status:        http.StatusOK,
		Message:       "Payment Transaction Succeeds",
		TransactionID: res,
	}

	return c.JSON(http.StatusOK, response)
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
	case models.BookingTypeRequired:
		return http.StatusBadRequest
	case models.BookingExpIdRequired:
		return http.StatusBadRequest
	case models.PromoIdRequired:
		return http.StatusBadRequest
	case models.PaymentMethodIdRequired:
		return http.StatusBadRequest
	case models.ExpPaymentIdRequired:
		return http.StatusBadRequest
	case models.StatusRequired:
		return http.StatusBadRequest
	case models.TotalPriceRequired:
		return http.StatusBadRequest
	case models.CurrencyRequired:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
