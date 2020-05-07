package reviews

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context,review models.Review)(string,error)
	CountRating(ctx context.Context, rating int, expID string) (int, error)
	GetByExpId(ctx context.Context, expID, sortBy string, rating, limit, offset int,userId string) ([]*models.Review, error)
}
