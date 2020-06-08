package include

import (
	"context"
	"github.com/models"
)

type Repository interface {
	List(ctx context.Context) ([]*models.Include, error)
	Fetch(ctx context.Context,limit,offset int)([]*models.Include, error)
	GetById(ctx context.Context,id int)(res *models.Include, err error)
	GetCount(ctx context.Context)(int,error)
	Insert(ctx context.Context,a *models.Include)(*int,error)
	Update(ctx context.Context,a *models.Include)error
	Delete(ctx context.Context,id int,deletedBy string)error
}
