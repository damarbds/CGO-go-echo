package usecase

import (
	"encoding/json"
	"github.com/auth/admin"
	"github.com/models"
	"github.com/service/cpc"
	"golang.org/x/net/context"
	"math"
	"strconv"
	"time"
)

type cPCUsecase struct {
	adminUsecase admin.Usecase
	cpcRepo cpc.Repository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewCPCUsecase(adminUsecase admin.Usecase,cpcRepo cpc.Repository,timeout time.Duration) cpc.Usecase {
	return &cPCUsecase{
		adminUsecase : adminUsecase,
		cpcRepo:cpcRepo,
		contextTimeout:  timeout,
	}
}
func (c cPCUsecase) GetAllCity(ctx context.Context, page, limit, offset int) (*models.CityDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	getCity,err := c.cpcRepo.FetchCity(ctx,limit,offset)
	if err != nil {
		return nil,err
	}
	var cityDtos []*models.CityDto
	for _,element := range getCity{
		cityPhotos := make([]models.CoverPhotosObj,0)
		if element.CityPhotos != nil {
			if errUnmarshal := json.Unmarshal([]byte(*element.CityPhotos), &cityPhotos); errUnmarshal != nil {
				return nil, models.ErrInternalServerError
			}
		}
		dto := models.CityDto{
			Id:         element.Id,
			CityName:   element.CityName,
			CityDesc:   element.CityDesc,
			CityPhotos: cityPhotos,
			ProvinceId: element.ProvinceId,
		}
		cityDtos = append(cityDtos,&dto)
	}
	totalRecords, _ := c.cpcRepo.GetCountCity(ctx)

	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}

	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(cityDtos),
	}

	response := &models.CityDtoWithPagination{
		Data: cityDtos,
		Meta: meta,
	}
	return response, nil
}

func (c cPCUsecase) GetCityById(ctx context.Context, id int) (*models.CityDto, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	getById ,err := c.cpcRepo.GetCityByID(ctx,id)
	if err != nil {
		return nil,err
	}
	cityPhotos := make([]models.CoverPhotosObj,0)
	if getById.CityPhotos != nil {
		if errUnmarshal := json.Unmarshal([]byte(*getById.CityPhotos), &cityPhotos); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	dto := models.CityDto{
		Id:         getById.Id,
		CityName:   getById.CityName,
		CityDesc:   getById.CityDesc,
		CityPhotos: cityPhotos,
		ProvinceId: getById.ProvinceId,
	}
	return &dto,nil
}

func (c cPCUsecase) CreateCity(ctx context.Context, p *models.NewCommandCity, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	cityPhotos ,_:= json.Marshal(p.CityPhotos)
	cityPhotosJson := string(cityPhotos)
	city := models.City{
		Id:           0,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Time{},
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		CityName:     p.CityName,
		CityDesc:     p.CityDesc,
		CityPhotos:   &cityPhotosJson,
		ProvinceId:   p.ProvinceId,
	}
	cityId,err := c.cpcRepo.InsertCity(ctx,&city)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      strconv.Itoa(*cityId),
		Message: "Success Create City",
	}

	return &result,nil
}

func (c cPCUsecase) UpdateCity(ctx context.Context, p *models.NewCommandCity, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	cityPhotos ,_:= json.Marshal(p.CityPhotos)
	cityPhotosJson := string(cityPhotos)
	city := models.City{
		Id:           p.Id,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &currentUser.Name,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		CityName:     p.CityName,
		CityDesc:     p.CityDesc,
		CityPhotos:   &cityPhotosJson,
		ProvinceId:   p.ProvinceId,
	}
	err = c.cpcRepo.UpdateCity(ctx,&city)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      strconv.Itoa(p.Id),
		Message: "Success Update City",
	}

	return &result,nil
}

func (c cPCUsecase) DeleteCity(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	err = c.cpcRepo.DeleteCity(ctx,id,currentUser.Name)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete City",
	}

	return &result,nil
}

func (c cPCUsecase) GetAllProvince(ctx context.Context, page, limit, offset int) (*models.ProvinceDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	getProvince,err := c.cpcRepo.FetchProvince(ctx,limit,offset)
	if err != nil {
		return nil,err
	}
	var provinceDtos []*models.ProvinceDto
	for _,element := range getProvince{
		dto := models.ProvinceDto{
			Id:           element.Id,
			ProvinceName: element.ProvinceName,
			CountryId:    element.CountryId,
		}
		provinceDtos = append(provinceDtos,&dto)
	}
	totalRecords, _ := c.cpcRepo.GetCountProvince(ctx)

	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}

	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(provinceDtos),
	}

	response := &models.ProvinceDtoWithPagination{
		Data: provinceDtos,
		Meta: meta,
	}
	return response, nil
}

func (c cPCUsecase) GetProvinceById(ctx context.Context, id int) (*models.ProvinceDto, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	getById ,err := c.cpcRepo.GetProvinceByID(ctx,id)
	if err != nil {
		return nil,err
	}

	dto := models.ProvinceDto{
		Id:           getById.Id,
		ProvinceName: getById.ProvinceName,
		CountryId:    getById.CountryId,
	}
	return &dto,nil
}

func (c cPCUsecase) CreateProvince(ctx context.Context, p *models.NewCommandProvince, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	province := models.Province{
		Id:           0,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Time{},
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		ProvinceName: p.ProvinceName,
		CountryId:    p.CountryId,
	}
	provinceId,err := c.cpcRepo.InsertProvince(ctx,&province)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      strconv.Itoa(*provinceId),
		Message: "Success Create Province",
	}

	return &result,nil
}

func (c cPCUsecase) UpdateProvince(ctx context.Context, p *models.NewCommandProvince, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	province := models.Province{
		Id:           p.Id,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &currentUser.Name,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		ProvinceName: p.ProvinceName,
		CountryId:    p.CountryId,
	}
	err = c.cpcRepo.UpdateProvince(ctx,&province)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      strconv.Itoa(p.Id),
		Message: "Success Update Province",
	}

	return &result,nil
}

func (c cPCUsecase) DeleteProvince(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	err = c.cpcRepo.DeleteProvince(ctx,id,currentUser.Name)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete Province",
	}

	return &result,nil
}

func (c cPCUsecase) GetAllCountry(ctx context.Context, page, limit, offset int) (*models.CountryDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	getCountry,err := c.cpcRepo.FetchCountry(ctx,limit,offset)
	if err != nil {
		return nil,err
	}
	var countryDtos []*models.CountryDto
	for _,element := range getCountry{
		dto := models.CountryDto{
			Id:          element.Id,
			CountryName: element.CountryName,
			PhoneCode: element.PhoneCode,
		}
		countryDtos = append(countryDtos,&dto)
	}
	totalRecords, _ := c.cpcRepo.GetCountCountry(ctx)

	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}

	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(countryDtos),
	}

	response := &models.CountryDtoWithPagination{
		Data: countryDtos,
		Meta: meta,
	}
	return response, nil
}

func (c cPCUsecase) GetCountryById(ctx context.Context, id int) (*models.CountryDto, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	getById ,err := c.cpcRepo.GetCountryByID(ctx,id)
	if err != nil {
		return nil,err
	}

	dto := models.CountryDto{
		Id:           	getById.Id,
		CountryName:	getById.CountryName,
		PhoneCode: 		getById.PhoneCode,
	}
	return &dto,nil
}

func (c cPCUsecase) CreateCountry(ctx context.Context, p *models.NewCommandCountry, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	country := models.Country{
		Id:           0,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Time{},
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		CountryName:p.CountryName,
		PhoneCode: p.PhoneCode,
	}
	countryId,err := c.cpcRepo.InsertCountry(ctx,&country)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      strconv.Itoa(*countryId),
		Message: "Success Create Country",
	}

	return &result,nil
}

func (c cPCUsecase) UpdateCountry(ctx context.Context, p *models.NewCommandCountry, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()

	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	country := models.Country{
		Id:           p.Id,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &currentUser.Name,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		CountryName:p.CountryName,
		PhoneCode: p.PhoneCode,
	}
	err = c.cpcRepo.UpdateCountry(ctx,&country)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      strconv.Itoa(p.Id),
		Message: "Success Update Country",
	}

	return &result,nil
}

func (c cPCUsecase) DeleteCountry(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	err = c.cpcRepo.DeleteCountry(ctx,id,currentUser.Name)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete Country",
	}

	return &result,nil
}
