package harbors

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Harbors, nextCursor string, err error)
	GetByID(ctx context.Context, id string) (*models.Harbors, error)
	GetAllWithJoinCPC(ctx context.Context, page *int, size *int, search string) ([]*models.HarborsWCPC, error)
	//GetByMerchantEmail(ctx context.Context, merchantEmail string) (*models.Harbors, error)
	//Update(ctx context.Context, ar *models.Harbors) error
	//Insert(ctx context.Context, a *models.Harbors) error
	//Delete(ctx context.Context, id string,deleted_by string) error
}
