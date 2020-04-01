package experience

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Experience, nextCursor string, err error)
	SearchExp(ctx context.Context, harborID, cityID string) ([]*models.ExpSearch, error)
	GetByID(ctx context.Context, id string) (*models.Experience, error)
	GetByExperienceEmail(ctx context.Context, userEmail string) (*models.Experience, error)
	GetUserDiscoverPreference(ctx context.Context,page *int,size *int)([]*models.ExpUserDiscoverPreference,error)
	//Update(ctx context.Context, ar *models.Experience) error
	//Insert(ctx context.Context, a *models.Experience) error
	Delete(ctx context.Context, id string,deleted_by string) error
}