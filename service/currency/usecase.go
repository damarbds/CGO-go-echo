package currency

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	GetAll(ctx context.Context ,page,limit,offset int)(*models.CurrencyDtoWithPagination,error)
	GetById(ctx context.Context,id int)(*models.CurrencyDto,error)
	Create(ctx context.Context, f *models.NewCommandCurrency,token string)(*models.ResponseDelete,error)
	Update(ctx context.Context, f *models.NewCommandCurrency,token string)(*models.ResponseDelete,error)
	Delete(ctx context.Context, id int,token string)(*models.ResponseDelete,error)
}
