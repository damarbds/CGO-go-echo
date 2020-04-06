package experience

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Experience, nextCursor string, err error)
	SearchExp(ctx context.Context, harborID, cityID string) ([]*models.ExpSearch, error)
	GetByID(ctx context.Context, id string) (*models.ExperienceJoinForegnKey, error)
	GetByExperienceEmail(ctx context.Context, userEmail string) (*models.Experience, error)
	GetUserDiscoverPreference(ctx context.Context,page *int,size *int)([]*models.ExpUserDiscoverPreference,error)
	GetIdByHarborsId(ctx context.Context, harborsId string) ([]*string, error)
	GetIdByCityId(ctx context.Context, cityId string) ([]*string, error)
	QueryFilterSearch(ctx context.Context,query string) ([]*models.ExpSearch, error)
	GetByCategoryID(ctx context.Context, categoryId int) ([]*models.ExpSearch, error)
	//Update(ctx context.Context, ar *models.Experience) error
	//Insert(ctx context.Context, a *models.Experience) error
	Delete(ctx context.Context, id string,deleted_by string) error
	GetSuccessBookCount(ctx context.Context, merchantId string) (int, error)
	GetExpCount(ctx context.Context, merchantId string) (int, error)
}