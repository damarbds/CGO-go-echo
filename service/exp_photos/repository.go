package exp_photos

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.ExpPhotos, nextCursor string, err error)
	GetByID(ctx context.Context, id string) (*models.ExpPhotos, error)
	GetByExperienceID(ctx context.Context, id string) ([]*models.ExpPhotos, error)
	//Update(ctx context.Context, ar *models.Experience) error
	//Insert(ctx context.Context, a *models.Experience) error
	Delete(ctx context.Context, id string,deleted_by string) error
}
