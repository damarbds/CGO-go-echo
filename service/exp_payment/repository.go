package exp_payment

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByExpID(ctx context.Context, expID string) (*models.ExperiencePayment, error)
}