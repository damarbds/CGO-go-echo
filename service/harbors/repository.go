package harbors

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Fetch(ctx context.Context, limit, offset int) ([]*models.Harbors, error)
	GetByID(ctx context.Context, id string) (*models.Harbors, error)
	GetCount(ctx context.Context) (int, error)
	GetAllWithJoinCPC(ctx context.Context, page *int, size *int, search string) ([]*models.HarborsWCPC, error)
	Update(ctx context.Context, ar *models.Harbors) error
	Insert(ctx context.Context, a *models.Harbors) (*string, error)
	Delete(ctx context.Context, id string, deletedBy string) error
}
