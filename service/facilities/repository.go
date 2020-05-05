package facilities

import (
	"context"

	"github.com/models"
)

type Repository interface {
	List(ctx context.Context) ([]*models.Facilities, error)
	Fetch(ctx context.Context,limit,offset int)([]*models.Facilities, error)
	GetById(ctx context.Context,id int)(res *models.Facilities, err error)
	GetCount(ctx context.Context)(int,error)
	Insert(ctx context.Context,a *models.Facilities)(*int,error)
	Update(ctx context.Context,a *models.Facilities)error
	Delete(ctx context.Context,id int,deletedBy string)error
}
