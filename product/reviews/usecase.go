package reviews

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetReviewsByExpId(ctx context.Context, exp_id, sortBy string, rating, limit, offset int) ([]*models.ReviewDto, error)
	GetReviewsByExpIdWithPagination(ctx context.Context, page, limit, offset, rating int, sortBy, exp_id string) (*models.ReviewsWithPagination, error)
}
