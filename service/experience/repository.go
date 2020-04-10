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
	QueryFilterSearch(ctx context.Context,query string, limit, offset int) ([]*models.ExpSearch, error)
	GetByCategoryID(ctx context.Context, categoryId int) ([]*models.ExpSearch, error)
	Update(ctx context.Context, a *models.Experience) (*string,error)
	Insert(ctx context.Context, a *models.Experience) (*string,error)
	Delete(ctx context.Context, id string,deleted_by string) error
	GetSuccessBookCount(ctx context.Context, merchantId string) (int, error)
	GetExpCount(ctx context.Context, merchantId string) (int, error)
	GetExpPendingTransactionCount(ctx context.Context, merchantId string) (int, error)
	GetExpFailedTransactionCount(ctx context.Context, merchantId string) (int, error)
	CountFilterSearch(ctx context.Context, query string) (int, error)
}