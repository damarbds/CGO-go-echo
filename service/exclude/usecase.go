package exclude

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context) ([]*models.ExcludeDto, error)
	GetAll(ctx context.Context ,page,limit,offset int)(*models.ExcludeDtoWithPagination,error)
	GetById(ctx context.Context,id int)(*models.ExcludeDto,error)
	Create(ctx context.Context, f *models.NewCommandExclude,token string)(*models.ResponseDelete,error)
	Update(ctx context.Context, f *models.NewCommandExclude,token string)(*models.ResponseDelete,error)
	Delete(ctx context.Context, id int,token string)(*models.ResponseDelete,error)
}