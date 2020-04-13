package experience_payment_type

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	GetAll(ctx context.Context, page *int, size *int) ([]*models.ExperiencePaymentTypeDto, error)
}
