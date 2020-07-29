package promo

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetCount(ctx context.Context)(int,error)
	Insert(ctx context.Context, promo *models.Promo)(string,error)
	Update(ctx context.Context, promo *models.Promo)error
	Delete(ctx context.Context,id string,deletedBy string)error
	GetById(ctx context.Context,id string)(*models.Promo,error)
	Fetch(ctx context.Context, page *int, size *int,search string) ([]*models.Promo, error)
	GetByCode(ctx context.Context, code string,promoType *int,merchantId string) ([]*models.Promo, error)
	GetByFilter(ctx context.Context, code string,promoType *int,merchantExpId string,merchantTransId string) ([]*models.Promo, error)
}
