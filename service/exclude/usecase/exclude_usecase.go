package usecase

import (
	"context"
	"github.com/auth/admin"
	"github.com/models"
	"github.com/service/exclude"
	"math"
	"strconv"
	"time"
)

type excludeUsecase struct {
	adminUsecase   admin.Usecase
	excludeRepo   exclude.Repository
	contextTimeout time.Duration
}


func NewExcludeUsecase(adminUsecase admin.Usecase,f exclude.Repository, timeout time.Duration) exclude.Usecase {
	return &excludeUsecase{
		adminUsecase:adminUsecase,
		excludeRepo:   f,
		contextTimeout: timeout,
	}
}

func (f excludeUsecase) List(ctx context.Context) ([]*models.ExcludeDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	res, err := f.excludeRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*models.ExcludeDto, len(res))
	for i, n := range res {
		results[i] = &models.ExcludeDto{
			Id:           n.Id,
			ExcludeIcon: n.ExcludeIcon,
			ExcludeName: n.ExcludeName,
		}
	}

	return results, nil
}

func (f excludeUsecase) GetAll(ctx context.Context, page, limit, offset int) (*models.ExcludeDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getExcludes,err := f.excludeRepo.Fetch(ctx,limit,offset)
	if err != nil {
		return nil,err
	}

	excludesDtos := make([]*models.ExcludeDto,0)

	for _,element := range getExcludes{
		dto := models.ExcludeDto{
			Id:          element.Id,
			ExcludeName: element.ExcludeName,
			ExcludeIcon: element.ExcludeIcon,
		}
		excludesDtos = append(excludesDtos,&dto)
	}
	totalRecords, _ := f.excludeRepo.GetCount(ctx)

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
		RecordPerPage: len(excludesDtos),
	}

	response := &models.ExcludeDtoWithPagination{
		Data: excludesDtos,
		Meta: meta,
	}
	return response, nil
}

func (f excludeUsecase) GetById(ctx context.Context, id int) (*models.ExcludeDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getById,err := f.excludeRepo.GetById(ctx,id)
	if err != nil {
		return nil,err
	}
	result := models.ExcludeDto{
		Id:          getById.Id,
		ExcludeName: getById.ExcludeName,
		ExcludeIcon: getById.ExcludeIcon,
	}
	return &result,nil

}

func (f excludeUsecase) Create(ctx context.Context, exc *models.NewCommandExclude, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	excludes := models.Exclude{
		Id:           0,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		ExcludeName:  exc.ExcludeName,
		ExcludeIcon:  exc.ExcludeIcon,
	}
	id ,err := f.excludeRepo.Insert(ctx,&excludes)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(*id),
		Message: "Success Create Exclude",
	}
	return &result,nil
}

func (f excludeUsecase) Update(ctx context.Context, exc *models.NewCommandExclude, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	if exc.ExcludeIcon == ""{
		getById ,_ := f.excludeRepo.GetById(ctx,exc.Id)
		exc.ExcludeIcon = getById.ExcludeIcon
	}
	now := time.Now()
	excludes := models.Exclude{
		Id:           exc.Id,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Now(),
		ModifiedBy:   &currentUser.Name,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		ExcludeName:  exc.ExcludeName,
		ExcludeIcon:  exc.ExcludeIcon,
	}
	err = f.excludeRepo.Update(ctx,&excludes)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(exc.Id),
		Message: "Success Update Exclude",
	}
	return &result,nil
}

func (f excludeUsecase) Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	err = f.excludeRepo.Delete(ctx,id,currentUser.Name)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete Exclude",
	}
	return &result,nil
}