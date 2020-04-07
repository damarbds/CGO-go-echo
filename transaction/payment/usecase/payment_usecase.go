package usecase

import (
	"context"
	"github.com/auth/user"
	"github.com/booking/booking_exp"
	"github.com/models"
	"github.com/transaction/payment"
	"time"
)

type paymentUsecase struct {
	paymentRepo payment.Repository
	userUsercase user.Usecase
	bookingRepo booking_exp.Repository
	contextTimeout time.Duration
}

// NewPaymentUsecase will create new an paymentUsecase object representation of payment.Usecase interface
func NewPaymentUsecase(p payment.Repository, u user.Usecase, b booking_exp.Repository, timeout time.Duration) payment.Usecase {
	return &paymentUsecase{
		paymentRepo: p,
		userUsercase: u,
		bookingRepo: b,
		contextTimeout: timeout,
	}
}

func (p paymentUsecase) Insert(ctx context.Context, payment *models.Transaction, token string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	if payment.BookingExpId == "" {
		return "", models.BookingExpIdRequired
	}

	//if payment.PromoId == "" {
	//	return "", models.PromoIdRequired
	//}

	if payment.PaymentMethodId == "" {
		return "", models.PaymentMethodIdRequired
	}

	if payment.ExperiencePaymentId == "" {
		return "", models.ExpPaymentIdRequired
	}

	if payment.Currency == "" {
		payment.Currency = "IDR"
	}

	createdBy, err := p.bookingRepo.GetEmailByID(ctx, payment.BookingExpId)
	if err != nil {
		return "", nil
	}
	if token != "" {
		currentUser, err := p.userUsercase.ValidateTokenUser(ctx, token)
		if err != nil {
			return "", err
		}
		createdBy = currentUser.UserEmail
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
	}

	res, err := p.paymentRepo.Insert(ctx, newData)
	if err != nil {
		return "", models.ErrInternalServerError
	}
	expiredDatePayment := newData.CreatedDate.AddDate(0,0,1)
	err = p.bookingRepo.UpdateStatus(ctx, res.BookingExpId,expiredDatePayment)
	if err != nil {
		return "", err
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

	return nil
}

