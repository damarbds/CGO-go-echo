package currency

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Fetch(ctx context.Context,limit,offset int)([]*models.Currency, error)
	GetById(ctx context.Context,id int)(res *models.Currency, err error)
	GetCount(ctx context.Context)(int,error)
	Insert(ctx context.Context,a *models.Currency)(*int,error)
	Update(ctx context.Context,a *models.Currency)error
	Delete(ctx context.Context,id int,deletedBy string)error
}
