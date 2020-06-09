package exclude

import (
	"context"
	"github.com/models"
)

type Repository interface {
	List(ctx context.Context) ([]*models.Exclude, error)
	Fetch(ctx context.Context,limit,offset int)([]*models.Exclude, error)
	GetById(ctx context.Context,id int)(res *models.Exclude, err error)
	GetCount(ctx context.Context)(int,error)
	Insert(ctx context.Context,a *models.Exclude)(*int,error)
	Update(ctx context.Context,a *models.Exclude)error
	Delete(ctx context.Context,id int,deletedBy string)error
}