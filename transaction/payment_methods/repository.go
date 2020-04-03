package payment_methods

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Fetch(ctx context.Context) ([]*models.PaymentMethod, error)
}
