package include

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context) ([]*models.IncludeDto, error)
	GetAll(ctx context.Context ,page,limit,offset int)(*models.IncludeDtoWithPagination,error)
	GetById(ctx context.Context,id int)(*models.IncludeDto,error)
	Create(ctx context.Context, f *models.NewCommandInclude,token string)(*models.ResponseDelete,error)
	Update(ctx context.Context, f *models.NewCommandInclude,token string)(*models.ResponseDelete,error)
	Delete(ctx context.Context, id int,token string)(*models.ResponseDelete,error)
}
