package exp_photos

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.ExpPhotos, nextCursor string, err error)
	GetByID(ctx context.Context, id string) (*models.ExpPhotos, error)
	GetByExperienceID(ctx context.Context, id string) ([]*models.ExpPhotos, error)
	Update(ctx context.Context, a *models.ExpPhotos) (*string, error)
	Insert(ctx context.Context, a *models.ExpPhotos) (*string, error)
	Deletes(ctx context.Context, ids []string, expId string, deletedBy string) error
	DeleteByExpId(ctx context.Context,expId string,deletedBy string)error
}
