package promo

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	List(ctx context.Context,page, limit, offset int, search string,token string,trans bool,exp bool,merchantIds []string)(*models.PromoWithPagination,error)
	Update(ctx context.Context, command models.NewCommandPromo,token string)(*models.NewCommandPromo,error)
	Create(ctx context.Context, command models.NewCommandPromo,token string)(*models.NewCommandPromo,error)
	Delete(ctx context.Context, id string,token string)(*models.ResponseDelete,error)
	GetDetail(ctx context.Context, id string,token string)(*models.PromoDto,error)
	Fetch(ctx context.Context, page *int, size *int,search string,trans bool,exp bool,merchantIds []string,sortBy string,promoId string) ([]*models.PromoDto, error)
	FetchUser(ctx context.Context, page *int, size *int, token string,search string,trans bool,exp bool,merchantIds []string,sortBy string,promoId string) ([]*models.PromoDto, error)
	GetByCode(ctx context.Context, code string,promoType string,merchantId string,token string,bookingId string,isAdmin bool) (*models.PromoDto, error)
	GetByFilter(ctx context.Context, code string,promoType int,merchantExpId string, merchantTransportId,token string) (*models.PromoDto, error)
}
