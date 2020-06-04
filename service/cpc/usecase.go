package cpc

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	GetAllCity(ctx context.Context,page,limit,offset int)(*models.CityDtoWithPagination,error)
	GetCityById(ctx context.Context,id int)(*models.CityDto,error)
	CreateCity(ctx context.Context,p *models.NewCommandCity,token string)(*models.ResponseDelete,error)
	UpdateCity(ctx context.Context,p *models.NewCommandCity,token string)(*models.ResponseDelete,error)
	DeleteCity(ctx context.Context,id int,token string)(*models.ResponseDelete,error)
	GetAllProvince(ctx context.Context,page,limit,offset int)(*models.ProvinceDtoWithPagination,error)
	GetProvinceById(ctx context.Context,id int)(*models.ProvinceDto,error)
	CreateProvince(ctx context.Context,p *models.NewCommandProvince,token string)(*models.ResponseDelete,error)
	UpdateProvince(ctx context.Context,p *models.NewCommandProvince,token string)(*models.ResponseDelete,error)
	DeleteProvince(ctx context.Context,id int,token string)(*models.ResponseDelete,error)
	GetAllCountry(ctx context.Context,page,limit,offset int)(*models.CountryDtoWithPagination,error)
	GetCountryById(ctx context.Context,id int)(*models.CountryDto,error)
	CreateCountry(ctx context.Context,p *models.NewCommandCountry,token string)(*models.ResponseDelete,error)
	UpdateCountry(ctx context.Context,p *models.NewCommandCountry,token string)(*models.ResponseDelete,error)
	DeleteCountry(ctx context.Context,id int,token string)(*models.ResponseDelete,error)
}
