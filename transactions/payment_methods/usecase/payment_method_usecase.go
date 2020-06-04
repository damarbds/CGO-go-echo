package usecase

import (
	"context"
	"github.com/models"
	pm "github.com/transactions/payment_methods"
	"time"
)

type paymentMethodUsecase struct {
	paymentMethodRepo pm.Repository
	contextTimeout    time.Duration
}

// NewPaymentMethodUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPaymentMethodUsecase(p pm.Repository, timeout time.Duration) pm.Usecase {
	return &paymentMethodUsecase{
		paymentMethodRepo: p,
		contextTimeout:    timeout,
	}
}

func (p paymentMethodUsecase) Fetch(ctx context.Context) ([]*models.PaymentMethodObject, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	pmList, err := p.paymentMethodRepo.Fetch(ctx)
	if err != nil {
		return nil, err
	}

	pmDto := make([]*models.PaymentMethodObject, len(pmList))
	for i, pay := range pmList {
		pmDto[i] = &models.PaymentMethodObject{
			Id:   pay.Id,
			Name: pay.Name,
			Type: pay.Type,
			Desc: pay.Desc.String,
			Icon: pay.Icon,
			MidtransPaymentCode:pay.MidtransPaymentCode,
		}
	}

	return pmDto, nil
}
