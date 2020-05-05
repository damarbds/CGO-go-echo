package usecase

import (
	"github.com/models"
	"github.com/service/cpc"
	"golang.org/x/net/context"
	"time"
)

type cPCUsecase struct {
	cpcRepo cpc.Repository
	contextTimeout time.Duration
}

func (c cPCUsecase) GetAllCity(ctx context.Context, page, limit, offset int) ([]*models.CityDto, error) {
	panic("implement me")
}

func (c cPCUsecase) GetCityById(ctx context.Context, id int) (models.CityDto, error) {
	panic("implement me")
}

func (c cPCUsecase) CreateCity(ctx context.Context, p *models.NewCommandCity, token string) (*models.ResponseDelete, error) {
	panic("implement me")
}

func (c cPCUsecase) UpdateCity(ctx context.Context, p *models.NewCommandCity, token string) (*models.ResponseDelete, error) {
	panic("implement me")
}

func (c cPCUsecase) DeleteCity(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	panic("implement me")
}

func (c cPCUsecase) GetAllProvince(ctx context.Context, page, limit, offset int) ([]*models.ProvinceDto, error) {
	panic("implement me")
}

func (c cPCUsecase) GetProvinceById(ctx context.Context, id int) (models.ProvinceDto, error) {
	panic("implement me")
}

func (c cPCUsecase) CreateProvince(ctx context.Context, p *models.NewCommandProvince, token string) (*models.ResponseDelete, error) {
	panic("implement me")
}

func (c cPCUsecase) UpdateProvince(ctx context.Context, p *models.NewCommandProvince, token string) (*models.ResponseDelete, error) {
	panic("implement me")
}

func (c cPCUsecase) DeleteProvince(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	panic("implement me")
}

func (c cPCUsecase) GetAllCountry(ctx context.Context, page, limit, offset int) ([]*models.CountryDto, error) {
	panic("implement me")
}

func (c cPCUsecase) GetCountryById(ctx context.Context, id int) (models.CountryDto, error) {
	panic("implement me")
}

func (c cPCUsecase) CreateCountry(ctx context.Context, p *models.NewCommandProvince, token string) (*models.ResponseDelete, error) {
	panic("implement me")
}

func (c cPCUsecase) UpdateCountry(ctx context.Context, p *models.NewCommandProvince, token string) (*models.ResponseDelete, error) {
	panic("implement me")
}

func (c cPCUsecase) DeleteCountry(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	panic("implement me")
}

// NewPromoUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewCPCUsecase(cpc cpc.Repository,timeout time.Duration) cpc.Usecase {
	return &cPCUsecase{
		cpcRepo:cpc,
		contextTimeout: timeout,
	}
}