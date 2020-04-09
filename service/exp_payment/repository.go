package exp_payment

import (
	"context"

	"github.com/models"
)

type Repository interface {
	GetByExpID(ctx context.Context, expID string) ([]*models.ExperiencePaymentJoinType, error)
	Insert(ctx context.Context,payment models.ExperiencePayment)error
}
