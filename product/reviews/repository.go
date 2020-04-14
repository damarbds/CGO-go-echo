package reviews

import (
	"context"
	"github.com/models"
)

type Repository interface {
	CountRating(ctx context.Context, rating int, expID string) (int, error)
	GetByExpId(ctx context.Context, expID, sortBy string, rating, limit, offset int) ([]*models.Review, error)
}
