package http

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/auth/identityserver"

	"github.com/booking/booking_exp"
	"github.com/labstack/echo"
	"github.com/models"
	"github.com/service/experience"
	"github.com/sirupsen/logrus"
	"github.com/third-party/midtrans"
	"github.com/transactions/transaction"
)

type ResponseError struct {
	Message string `json:"message"`
}

type midtransHandler struct {
	bookingRepo     booking_exp.Repository
	expRepo         experience.Repository
	transactionRepo transaction.Repository
	isUsecase       identityserver.Usecase
}

func NewMidtransHandler(e *echo.Echo, br booking_exp.Repository, er experience.Repository, tr transaction.Repository, is identityserver.Usecase) {
	handler := &midtransHandler{
		bookingRepo:     br,
		expRepo:         er,
		transactionRepo: tr,
		isUsecase:       is,
	}
	e.POST("/midtrans/notif", handler.MidtransNotif)
}

func (m *midtransHandler) MidtransNotif(c echo.Context) error {
	var callback midtrans.MidtransCallback
	if err := c.Bind(&callback); err != nil {
		return err
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	input := []byte(callback.OrderId + callback.StatusCode + callback.GrossAmount + midtrans.MidtransServerKey)
	hash := sha512.New()
	_, err := hash.Write(input)
	if err != nil {
		return err
	}
	signatureKey := hex.EncodeToString(hash.Sum(nil))

	if callback.SignatureKey != signatureKey {
		return c.JSON(http.StatusBadGateway, ResponseError{Message: "Signature key invalid"})
	}

	booking, err := m.bookingRepo.GetByID(ctx, callback.OrderId)
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
	if callback.TransactionStatus == "capture" || callback.TransactionStatus == "settlement" {
		if booking.ExpId != nil {
			exp, err := m.expRepo.GetByID(ctx, *booking.ExpId)
			if err != nil {
				return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
			}
			if exp.ExpBookingType == "No Instant Booking" {
				transactionStatus = 1
			} else {
				transactionStatus = 2
			}
			if err := m.transactionRepo.UpdateStatus(ctx, transactionStatus, "", booking.OrderId); err != nil {
				return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
			}
		} else {
			transactionStatus = 2
			if err := m.transactionRepo.UpdateStatus(ctx, transactionStatus, "", booking.OrderId); err != nil {
				return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
			}
		}
		msg := "<p>This is your order id " + booking.OrderId + " and your ticket QR code " + booking.TicketQRCode + "</p>"
		pushEmail := &models.SendingEmail{
			Subject: "E-Ticket cGO",
			Message: msg,
			From:    "CGO Indonesia",
			To:      bookedBy[0].Email,
			FileName:"Ticket.pdf",
		}
		if _, err := m.isUsecase.SendingEmail(pushEmail); err != nil {
			return nil
		}
	}

	if callback.TransactionStatus == "expire" || callback.TransactionStatus == "deny" {
		transactionStatus = 3
		if err := m.transactionRepo.UpdateStatus(ctx, transactionStatus, "", booking.Id); err != nil {
			return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
		}
	}

	return c.JSON(http.StatusOK, ResponseError{Message: "Midtrans Notification Succeed"})
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
