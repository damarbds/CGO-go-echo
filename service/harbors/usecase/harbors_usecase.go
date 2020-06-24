package usecase

import (
	"context"
	"github.com/auth/admin"
	"math"
	"time"

	"github.com/service/harbors"

	"github.com/models"
)

type harborsUsecase struct {
	adminUsecase 	admin.Usecase
	harborsRepo    harbors.Repository
	contextTimeout time.Duration
}


// NewharborsUsecase will create new an harborsUsecase object representation of harbors.Usecase interface
func NewharborsUsecase(adminUsecase admin.Usecase,a harbors.Repository, timeout time.Duration) harbors.Usecase {
	return &harborsUsecase{
		adminUsecase:adminUsecase,
		harborsRepo:    a,
		contextTimeout: timeout,
	}
}
func (m harborsUsecase) GetAllWithJoinCPC(c context.Context, page *int, size *int, search string,harborsType string) ([]*models.HarborsWCPCDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err := m.harborsRepo.GetAllWithJoinCPC(ctx, page, size, search,harborsType)
	if err != nil {
		return nil, err
	}
	var harborss []*models.HarborsWCPCDto
	for _, element := range res {
		harbors := models.HarborsWCPCDto{
			Id:               element.Id,
			HarborsName:      element.HarborsName,
			HarborsLongitude: element.HarborsLongitude,
			HarborsLatitude:  element.HarborsLatitude,
			HarborsImage:     element.HarborsImage,
			CityId:           element.CityId,
			City:             element.CityName,
			ProvinceId:element.ProvinceId,
			Province:         element.ProvinceName,
			Country:          element.CountryName,
		}
		harborss = append(harborss, &harbors)
	}

	return harborss, nil
}

func (m harborsUsecase) GetAll(c context.Context, page, limit, size int) (*models.HarborsDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	getAll ,err := m.harborsRepo.Fetch(ctx,limit,size)
	if err != nil {
		return nil,err
	}

	harborsDtos := make([]*models.HarborsDto,0)
	for _,element := range getAll{

		dto := models.HarborsDto{
			Id:               element.Id,
			HarborsName:      element.HarborsName,
			HarborsLongitude: element.HarborsLongitude,
			HarborsLatitude:  element.HarborsLatitude,
			HarborsImage:     element.HarborsImage,
			CityId:           element.CityId,
			HarborsType:element.HarborsType,
		}

		harborsDtos = append(harborsDtos,&dto)
	}
	totalRecords, _ := m.harborsRepo.GetCount(ctx)

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
		RecordPerPage: len(harborsDtos),
	}

	response := &models.HarborsDtoWithPagination{
		Data: harborsDtos,
		Meta: meta,
	}
	return response, nil
}

func (m harborsUsecase) GetById(c context.Context, id string) (*models.HarborsDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	getById ,err := m.harborsRepo.GetByID(ctx,id)
	if err != nil {
		return nil,err
	}


	result := models.HarborsDto{
		Id:               getById.Id,
		HarborsName:      getById.HarborsName,
		HarborsLongitude: getById.HarborsLongitude,
		HarborsLatitude:  getById.HarborsLatitude,
		HarborsImage:     getById.HarborsImage,
		CityId:           getById.CityId,
		HarborsType:getById.HarborsType,
	}

	return &result,nil
}

func (m harborsUsecase) Create(ctx context.Context, p *models.NewCommandHarbors, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	currentUser, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	harbors := models.Harbors{
		Id:               "",
		CreatedBy:        currentUser.Name,
		CreatedDate:      time.Time{},
		ModifiedBy:       nil,
		ModifiedDate:     nil,
		DeletedBy:        nil,
		DeletedDate:      nil,
		IsDeleted:        0,
		IsActive:         0,
		HarborsName:      p.HarborsName,
		HarborsLongitude: p.HarborsLongitude,
		HarborsLatitude:  p.HarborsLatitude,
		HarborsImage:     p.HarborsImage,
		CityId:           p.CityId,
		HarborsType:		&p.HarborsType,
	}
	harborsId,err := m.harborsRepo.Insert(ctx,&harbors)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      *harborsId,
		Message: "Success Create Harbors",
	}

	return &result,nil
}

func (m harborsUsecase) Update(ctx context.Context, p *models.NewCommandHarbors, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	currentUser, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	harbors := models.Harbors{
		Id:               p.Id,
		CreatedBy:        "",
		CreatedDate:      time.Time{},
		ModifiedBy:       &currentUser.Name,
		ModifiedDate:     nil,
		DeletedBy:        nil,
		DeletedDate:      nil,
		IsDeleted:        0,
		IsActive:         0,
		HarborsName:      p.HarborsName,
		HarborsLongitude: p.HarborsLongitude,
		HarborsLatitude:  p.HarborsLatitude,
		HarborsImage:     p.HarborsImage,
		CityId:           p.CityId,
		HarborsType:		&p.HarborsType,
	}
	err = m.harborsRepo.Update(ctx,&harbors)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      p.Id,
		Message: "Success Update Harbors",
	}

	return &result,nil
}

func (c harborsUsecase) Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, c.contextTimeout)
	defer cancel()
	currentUser, err := c.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	err = c.harborsRepo.Delete(ctx,id,currentUser.Name)
	if err != nil {
		return nil,err
	}
	result := models.ResponseDelete{
		Id:      id,
		Message: "Success Delete Harbors",
	}

	return &result,nil
}
/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
