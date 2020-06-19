package harbors

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetAllWithJoinCPC(ctx context.Context, page *int, size *int, search string,harborsType string) ([]*models.HarborsWCPCDto, error)
	GetAll(ctx context.Context,page,limit,size int)(*models.HarborsDtoWithPagination,error)
	GetById(ctx context.Context,id string)(*models.HarborsDto,error)
	Create(ctx context.Context,a *models.NewCommandHarbors,token string)(*models.ResponseDelete,error)
	Update(ctx context.Context,a *models.NewCommandHarbors,token string)(*models.ResponseDelete,error)
	Delete(ctx context.Context,id string,token string)(*models.ResponseDelete,error)
}
