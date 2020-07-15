package usecase

import (
	"context"
	"github.com/auth/admin"
	"math"
	"strconv"
	"time"

	"github.com/models"
	"github.com/service/facilities"
)

type facilityUsecase struct {
	adminUsecase   admin.Usecase
	facilityRepo   facilities.Repository
	contextTimeout time.Duration
}


func NewFacilityUsecase(adminUsecase admin.Usecase,f facilities.Repository, timeout time.Duration) facilities.Usecase {
	return &facilityUsecase{
		adminUsecase:adminUsecase,
		facilityRepo:   f,
		contextTimeout: timeout,
	}
}

func (f facilityUsecase) List(ctx context.Context) ([]*models.FacilityDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	res, err := f.facilityRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*models.FacilityDto, len(res))
	for i, n := range res {
		results[i] = &models.FacilityDto{
			Id:           n.Id,
			FacilityName: n.FacilityName,
			FacilityIcon: *n.FacilityIcon,
			IsNumerable:  n.IsNumerable,
		}
	}

	return results, nil
}

func (f facilityUsecase) GetAll(ctx context.Context, page, limit, offset int) (*models.FacilityDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getFacilities,err := f.facilityRepo.Fetch(ctx,limit,offset)
	if err != nil {
		return nil,err
	}

	facilitiesDtos := make( []*models.FacilityDto,0)

	for _,element := range getFacilities{
		dto := models.FacilityDto{
			Id:           element.Id,
			FacilityName: element.FacilityName,
			FacilityIcon: *element.FacilityIcon,
			IsNumerable:  element.IsNumerable,
		}
		facilitiesDtos = append(facilitiesDtos,&dto)
	}
	totalRecords, _ := f.facilityRepo.GetCount(ctx)

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
		RecordPerPage: len(facilitiesDtos),
	}

	response := &models.FacilityDtoWithPagination{
		Data: facilitiesDtos,
		Meta: meta,
	}
	return response, nil
}

func (f facilityUsecase) GetById(ctx context.Context, id int) (*models.FacilityDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getById,err := f.facilityRepo.GetById(ctx,id)
	if err != nil {
		return nil,err
	}
	result := models.FacilityDto{
		Id:           getById.Id,
		FacilityName: getById.FacilityName,
		FacilityIcon: *getById.FacilityIcon,
		IsNumerable:  getById.IsNumerable,
	}
	return &result,nil

}

func (f facilityUsecase) Create(ctx context.Context, fac *models.NewCommandFacilities, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	facilities := models.Facilities{
		Id:           0,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		FacilityName: fac.FacilityName,
		IsNumerable:  fac.IsNumerable,
		FacilityIcon: &fac.FacilityIcon,
	}
	id ,err := f.facilityRepo.Insert(ctx,&facilities)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(*id),
		Message: "Success Create Facilities",
	}
	return &result,nil
}

func (f facilityUsecase) Update(ctx context.Context, fac *models.NewCommandFacilities, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	if fac.FacilityIcon == ""{
		getFacilities ,_ := f.facilityRepo.GetById(ctx,fac.Id)
		fac.FacilityIcon = *getFacilities.FacilityIcon
	}
	now := time.Now()
	facilities := models.Facilities{
		Id:           fac.Id,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &currentUser.Name,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		FacilityName: fac.FacilityName,
		IsNumerable:  fac.IsNumerable,
		FacilityIcon: &fac.FacilityIcon,
	}
	err = f.facilityRepo.Update(ctx,&facilities)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(fac.Id),
		Message: "Success Update Facilities",
	}
	return &result,nil
}

func (f facilityUsecase) Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	err = f.facilityRepo.Delete(ctx,id,currentUser.Name)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete Facilities",
	}
	return &result,nil
}
