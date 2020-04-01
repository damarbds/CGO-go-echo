package reviews

import (
	"context"
	"github.com/models"
)

type Repository interface {
	CountRating(ctx context.Context, expID string) (int, error)
	GetByExpId(ctx context.Context, expID string) ([]*models.Review, error)
}
