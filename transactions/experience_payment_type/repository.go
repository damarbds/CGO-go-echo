package experience_payment_type

import (
	"context"
	"github.com/models"
)

type Repostiory interface {
	GetAll(ctx context.Context, page *int, size *int) ([]*models.ExperiencePaymentType, error)
}
