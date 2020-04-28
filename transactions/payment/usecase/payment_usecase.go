package usecase

import (
	"context"
	"time"

	"github.com/misc/notif"
	"github.com/transactions/transaction"

	"github.com/auth/user"
	"github.com/booking/booking_exp"
	"github.com/models"
	"github.com/transactions/payment"
)

type paymentUsecase struct {
	transactionRepo  transaction.Repository
	notificationRepo notif.Repository
	paymentRepo      payment.Repository
	userUsercase     user.Usecase
	bookingRepo      booking_exp.Repository
	userRepo         user.Repository
	contextTimeout   time.Duration
}

// NewPaymentUsecase will create new an paymentUsecase object representation of payment.Usecase interface
func NewPaymentUsecase(t transaction.Repository, n notif.Repository, p payment.Repository, u user.Usecase, b booking_exp.Repository, ur user.Repository, timeout time.Duration) payment.Usecase {
	return &paymentUsecase{
		transactionRepo:  t,
		notificationRepo: n,
		paymentRepo:      p,
		userUsercase:     u,
		bookingRepo:      b,
		userRepo:         ur,
		contextTimeout:   timeout,
	}
}

func (p paymentUsecase) Insert(ctx context.Context, payment *models.Transaction, token string, points float64) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	var userId string

	if payment.PaymentMethodId == "" {
		return "", models.PaymentMethodIdRequired
	}

	if payment.Currency == "" {
		payment.Currency = "IDR"
	}
	bookingCode := payment.OrderId
	if payment.BookingExpId != nil {
		bookingCode = payment.BookingExpId
	}
	createdBy, err := p.bookingRepo.GetEmailByID(ctx, *bookingCode)
	if err != nil {
		return "", err
	}
	if token != "" {
		currentUser, err := p.userUsercase.ValidateTokenUser(ctx, token)
		if err != nil {
			return "", err
		}
		createdBy = currentUser.UserEmail
		userId = currentUser.Id
	}

	newData := &models.Transaction{
		Id:                  "",
		CreatedBy:           createdBy,
		CreatedDate:         time.Now(),
		ModifiedBy:          nil,
		ModifiedDate:        nil,
		DeletedBy:           nil,
		DeletedDate:         nil,
		IsDeleted:           0,
		IsActive:            1,
		BookingType:         payment.BookingType,
		BookingExpId:        payment.BookingExpId,
		PromoId:             payment.PromoId,
		PaymentMethodId:     payment.PaymentMethodId,
		ExperiencePaymentId: payment.ExperiencePaymentId,
		Status:              payment.Status,
		TotalPrice:          payment.TotalPrice,
		Currency:            payment.Currency,
		OrderId:             payment.OrderId,
	}

	res, err := p.paymentRepo.Insert(ctx, newData)
	if err != nil {
		return "", models.ErrInternalServerError
	}

	expiredPayment := res.CreatedDate.Add(2 * time.Hour)
	err = p.bookingRepo.UpdateStatus(ctx, *bookingCode, expiredPayment)
	if err != nil {
		return "", err
	}

	if points != 0 {
		err = p.userRepo.UpdatePointByID(ctx, points, userId)
		if err != nil {
			return "", err
		}
	}

	return res.Id, nil
}

func (p paymentUsecase) ConfirmPayment(ctx context.Context, confirmIn *models.ConfirmPaymentIn) error {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	err := p.paymentRepo.ConfirmPayment(ctx, confirmIn)
	if err != nil {
		return err
	}
	getTransaction, err := p.transactionRepo.GetById(ctx, confirmIn.TransactionID)
	if err != nil {
		return err
	}
	notif := models.Notification{
		Id:           "",
		CreatedBy:    getTransaction.CreatedBy,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		MerchantId:   getTransaction.MerchantId,
		Type:         0,
		Title:        " New Order Receive: Order ID " + getTransaction.OrderIdBook,
		Desc:         "You got a booking for " + getTransaction.ExpTitle + " , booked by " + getTransaction.CreatedBy,
	}
	pushNotifErr := p.notificationRepo.Insert(ctx, notif)
	if pushNotifErr != nil {
		return nil
	}
	return nil
}
