package cpc

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	FetchCity(ctx context.Context, cursor string, num int64) (res []*models.City, nextCursor string, err error)
	GetCityByID(ctx context.Context, id int) (*models.City, error)
	FetchProvince(ctx context.Context, cursor string, num int64) (res []*models.Province, nextCursor string, err error)
	GetProvinceByID(ctx context.Context, id int) (*models.Province, error)
	//GetByMerchantEmail(ctx context.Context, merchantEmail string) (*models.Harbors, error)
	//Update(ctx context.Context, ar *models.Harbors) error
	//Insert(ctx context.Context, a *models.Harbors) error
	//Delete(ctx context.Context, id string,deleted_by string) error
}
