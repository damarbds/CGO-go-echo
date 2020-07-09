package promo

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	List(ctx context.Context,page, limit, offset int, search string,token string)(*models.PromoWithPagination,error)
	Update(ctx context.Context, command models.NewCommandPromo,token string)(*models.NewCommandPromo,error)
	Create(ctx context.Context, command models.NewCommandPromo,token string)(*models.NewCommandPromo,error)
	Delete(ctx context.Context, id string,token string)(*models.ResponseDelete,error)
	GetDetail(ctx context.Context, id string,token string)(*models.PromoDto,error)
	Fetch(ctx context.Context, page *int, size *int) ([]*models.PromoDto, error)
	GetByCode(ctx context.Context, code string,promoType int,merchantId string,token string) (*models.PromoDto, error)
}
