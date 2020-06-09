package rule

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context) ([]*models.RuleDto, error)
	GetAll(ctx context.Context ,page,limit,offset int)(*models.RuleDtoWithPagination,error)
	GetById(ctx context.Context,id int)(*models.RuleDto,error)
	Create(ctx context.Context, f *models.NewCommandRule,token string)(*models.ResponseDelete,error)
	Update(ctx context.Context, f *models.NewCommandRule,token string)(*models.ResponseDelete,error)
	Delete(ctx context.Context, id int,token string)(*models.ResponseDelete,error)
}