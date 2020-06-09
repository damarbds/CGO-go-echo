package rule

import (
	"context"
	"github.com/models"
)

type Repository interface {
	List(ctx context.Context) ([]*models.Rule, error)
	Fetch(ctx context.Context,limit,offset int)([]*models.Rule, error)
	GetById(ctx context.Context,id int)(res *models.Rule, err error)
	GetCount(ctx context.Context)(int,error)
	Insert(ctx context.Context,a *models.Rule)(*int,error)
	Update(ctx context.Context,a *models.Rule)error
	Delete(ctx context.Context,id int,deletedBy string)error
}
