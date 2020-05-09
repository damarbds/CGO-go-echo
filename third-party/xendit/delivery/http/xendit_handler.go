package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/auth/identityserver"
	"github.com/booking/booking_exp"
	"github.com/service/experience"
	"github.com/transactions/transaction"

	"github.com/third-party/xendit"

	"github.com/labstack/echo"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

type xenditHandler struct {
	bookingRepo     booking_exp.Repository
	expRepo         experience.Repository
	transactionRepo transaction.Repository
	bookingUseCase  booking_exp.Usecase
	isUsecase       identityserver.Usecase
}

func NewXenditHandler(e *echo.Echo, br booking_exp.Repository, er experience.Repository, tr transaction.Repository, bu booking_exp.Usecase, is identityserver.Usecase) {
	handler := &xenditHandler{
		bookingRepo:     br,
		expRepo:         er,
		transactionRepo: tr,
		bookingUseCase:  bu,
		isUsecase:       is,
	}
	e.POST("/xendit/callback", handler.XenditVACallback)
}

func (x *xenditHandler) XenditVACallback(c echo.Context) error {
	var callback xendit.VACallbackRequest
	if err := c.Bind(&callback); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	booking, err := x.bookingRepo.GetByID(ctx, callback.ExternalID)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	var bookedBy []models.BookedByObj
	if booking.BookedBy != "" {
		if errUnmarshal := json.Unmarshal([]byte(booking.BookedBy), &bookedBy); errUnmarshal != nil {
			return errUnmarshal
		}
	}

	var transactionStatus int
	if booking.ExpId != nil {
		exp, err := x.expRepo.GetByID(ctx, *booking.ExpId)
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		bookingDetail, err := x.bookingUseCase.GetDetailBookingID(ctx, booking.Id, "")
		if err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
		if exp.ExpBookingType == "No Instant Booking" {
			transactionStatus = 1
		} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Down Payment" {
			transactionStatus = 5
		} else if exp.ExpBookingType == "Instant Booking" && bookingDetail.ExperiencePaymentType.Name == "Full Payment" {
			transactionStatus = 2
		}
		if err := x.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, callback.AccountNumber, "", booking.Id); err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
	} else {
		transactionStatus = 2
		if err := x.transactionRepo.UpdateAfterPayment(ctx, transactionStatus, callback.AccountNumber, "", booking.OrderId); err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
	}
	msg := "<p>This is your order id " + booking.OrderId + " and your ticket QR code " + booking.TicketQRCode + "</p>"
	pushEmail := &models.SendingEmail{
		Subject:  "E-Ticket cGO",
		Message:  msg,
		From:     "CGO Indonesia",
		To:       bookedBy[0].Email,
		FileName: "Ticket.pdf",
	}
	if _, err := x.isUsecase.SendingEmail(pushEmail); err != nil {
		return nil
	}

	return c.JSON(http.StatusOK, ResponseError{Message: "Xendit Callback Succeed"})
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
