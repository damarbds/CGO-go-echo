package http

import (
	"context"
	"net/http"
	"strings"

	"github.com/booking/booking_exp"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"github.com/transactions/payment"
	"github.com/transactions/payment_methods"
	"gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

type paymentHandler struct {
	paymentUsecase    payment.Usecase
	bookingUsecase    booking_exp.Usecase
	bookingRepo       booking_exp.Repository
	paymentMethodRepo payment_methods.Repository
}

func NewPaymentHandler(e *echo.Echo, pus payment.Usecase, bus booking_exp.Usecase, bur booking_exp.Repository, pmr payment_methods.Repository) {
	handler := &paymentHandler{
		paymentUsecase:    pus,
		bookingUsecase:    bus,
		bookingRepo:       bur,
		paymentMethodRepo: pmr,
	}
	e.POST("/transaction/payments", handler.CreatePayment)
	e.PUT("/transaction/payments/confirm", handler.ConfirmPayment)
}

func isRequestValid(m *models.TransactionIn) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p *paymentHandler) ConfirmPayment(c echo.Context) error {
	//c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//token := c.Request().Header.Get("Authorization")
	//
	//if token == "" {
	//	return c.JSON(http.StatusUnauthorized, models.ErrUnAuthorize)
	//}

	cp := new(models.ConfirmPaymentIn)
	if err := c.Bind(cp); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrBadParamInput)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := p.paymentUsecase.ConfirmPayment(ctx, cp)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	response := &models.PaymentTransaction{
		Status:        http.StatusOK,
		Message:       "Confirm Payment Succeeds",
		TransactionID: cp.TransactionID,
	}

	return c.JSON(http.StatusOK, response)
}

func (p *paymentHandler) CreatePayment(c echo.Context) error {
	c.Request().Header.Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Response().Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	token := c.Request().Header.Get("Authorization")

	t := new(models.TransactionIn)
	if err := c.Bind(t); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrBadParamInput)
	}

	if ok, err := isRequestValid(t); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	var promoId *string
	if t.PromoId != "" {
		promoId = &t.PromoId
	} else {
		promoId = nil
	}

	var expPaymentId *string
	if t.ExperiencePaymentId != "" {
		expPaymentId = &t.ExperiencePaymentId
	} else {
		expPaymentId = nil
	}

	var orderId *string
	if t.OrderId != "" {
		orderId = &t.OrderId
	} else {
		orderId = nil
	}

	var bookingId *string
	if t.BookingId != "" {
		bookingId = &t.BookingId
	} else {
		bookingId = nil
	}

	tr := &models.Transaction{
		BookingType:         t.BookingType,
		BookingExpId:        bookingId,
		OrderId:             orderId,
		PromoId:             promoId,
		PaymentMethodId:     t.PaymentMethodId,
		ExperiencePaymentId: expPaymentId,
		Status:              t.Status,
		TotalPrice:          t.TotalPrice,
		Currency:            t.Currency,
		ExChangeRates:		&t.ExChangeRates,
		ExChangeCurrency:	&t.ExChangeCurrency,
	}

	_, err := p.paymentUsecase.Insert(ctx, tr, token, t.Points)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	pm, err := p.paymentMethodRepo.GetByID(ctx, tr.PaymentMethodId)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	bookingCode := t.OrderId
	if tr.BookingExpId != nil {
		bookingCode = t.BookingId
	}

	var data map[string]interface{}

	if t.PaypalOrderId != "" && *pm.MidtransPaymentCode == "paypal"  && strings.Contains(t.PaypalOrderId,"PAYID"){
		response, err := p.bookingUsecase.PaypalAutoComplete(ctx, *tr.BookingExpId)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		return c.JSON(http.StatusOK, response)
	}else if t.PaypalOrderId != "" && *pm.MidtransPaymentCode == "paypal" {
		data, err = p.bookingUsecase.Verify(ctx, t.PaypalOrderId, bookingCode)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
	} else if *pm.MidtransPaymentCode == "BRI" {
		data, err = p.bookingUsecase.XenPayment(ctx, 0, "", "", bookingCode, *pm.MidtransPaymentCode)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
	} else if *pm.MidtransPaymentCode == "cc" || (t.CcAuthId != "" && t.CcTokenId != "") {
		data, err = p.bookingUsecase.XenPayment(ctx, t.TotalPrice, t.CcTokenId, t.CcAuthId, bookingCode, *pm.MidtransPaymentCode)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
	} else {
		data, err = p.bookingUsecase.SendCharge(ctx, bookingCode, *pm.MidtransPaymentCode)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}

		if err := p.bookingRepo.UpdatePaymentUrl(ctx, bookingCode, data["redirect_url"].(string)); err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, data)
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
