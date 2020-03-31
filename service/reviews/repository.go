package reviews

import (
	"context"
)

type Repository interface {
	CountRating(ctx context.Context, expID string) (int, error)
}
