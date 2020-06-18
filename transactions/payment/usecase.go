package payment

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	Insert(ctx context.Context, payment *models.Transaction, token string, points float64,autoComplete bool) (string, error)
	ConfirmPayment(ctx context.Context, confirmIn *models.ConfirmPaymentIn) error
}
