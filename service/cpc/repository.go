package cpc

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	FetchCity(ctx context.Context, limit,offset int) ([]*models.City, error)
	GetCityByID(ctx context.Context, id int) (*models.City, error)
	GetCountCity(ctx context.Context)(int,error)
	InsertCity(ctx context.Context, a *models.City)(*int,error)
	UpdateCity(ctx context.Context, a *models.City)error
	DeleteCity(ctx context.Context, id string,deleted_by string) error
	FetchProvince(ctx context.Context, limit,offset int) ([]*models.Province, error)
	GetProvinceByID(ctx context.Context, id int) (*models.Province, error)
	GetCountProvince(ctx context.Context)(int,error)
	InsertProvince(ctx context.Context, a *models.Province)(*int,error)
	UpdateProvince(ctx context.Context, a *models.Province)error
	DeleteProvince(ctx context.Context, id string,deleted_by string) error
	FetchCountry(ctx context.Context, limit,offset int) ([]*models.Country, error)
	GetCountryByID(ctx context.Context, id int) (*models.Country, error)
	GetCountCountry(ctx context.Context)(int,error)
	InsertCountry(ctx context.Context, a *models.Country)(*int,error)
	UpdateCountry(ctx context.Context, a *models.Country)error
	DeleteCountry(ctx context.Context, id string,deleted_by string) error
}
