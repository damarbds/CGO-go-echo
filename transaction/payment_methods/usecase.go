package payment_methods

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	Fetch(ctx context.Context) ([]*models.PaymentMethodObject, error)
}
