package usecase

import (
	"context"
	"github.com/auth/admin"
	"github.com/models"
	"github.com/service/include"
	"math"
	"strconv"
	"time"
)

type includeUsecase struct {
	adminUsecase   admin.Usecase
	includeRepo   include.Repository
	contextTimeout time.Duration
}


func NewIncludeUsecase(adminUsecase admin.Usecase,f include.Repository, timeout time.Duration) include.Usecase {
	return &includeUsecase{
		adminUsecase:adminUsecase,
		includeRepo:   f,
		contextTimeout: timeout,
	}
}

func (f includeUsecase) List(ctx context.Context) ([]*models.IncludeDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	res, err := f.includeRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]*models.IncludeDto, len(res))
	for i, n := range res {
		results[i] = &models.IncludeDto{
			Id:           n.Id,
			IncludeIcon: n.IncludeIcon,
			IncludeName: n.IncludeName,
		}
	}

	return results, nil
}

func (f includeUsecase) GetAll(ctx context.Context, page, limit, offset int) (*models.IncludeDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getIncludes,err := f.includeRepo.Fetch(ctx,limit,offset)
	if err != nil {
		return nil,err
	}

	includesDtos := make([]*models.IncludeDto,0)

	for _,element := range getIncludes{
		dto := models.IncludeDto{
			Id:          element.Id,
			IncludeName: element.IncludeName,
			IncludeIcon: element.IncludeIcon,
		}
		includesDtos = append(includesDtos,&dto)
	}
	totalRecords, _ := f.includeRepo.GetCount(ctx)

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
		RecordPerPage: len(includesDtos),
	}

	response := &models.IncludeDtoWithPagination{
		Data: includesDtos,
		Meta: meta,
	}
	return response, nil
}

func (f includeUsecase) GetById(ctx context.Context, id int) (*models.IncludeDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getById,err := f.includeRepo.GetById(ctx,id)
	if err != nil {
		return nil,err
	}
	result := models.IncludeDto{
		Id:          getById.Id,
		IncludeName: getById.IncludeName,
		IncludeIcon: getById.IncludeIcon,
	}
	return &result,nil

}

func (f includeUsecase) Create(ctx context.Context, inc *models.NewCommandInclude, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	includes := models.Include{
		Id:           0,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Time{},
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		IncludeName:  inc.IncludeName,
		IncludeIcon:  inc.IncludeIcon,
	}
	id ,err := f.includeRepo.Insert(ctx,&includes)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(*id),
		Message: "Success Create Include",
	}
	return &result,nil
}

func (f includeUsecase) Update(ctx context.Context, inc *models.NewCommandInclude, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}


	if inc.IncludeIcon == ""{
		getById,_:= f.includeRepo.GetById(ctx,inc.Id)
		inc.IncludeIcon = getById.IncludeIcon
	}

	includes := models.Include{
		Id:           inc.Id,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Time{},
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		IncludeName:  inc.IncludeName,
		IncludeIcon:  inc.IncludeIcon,
	}
	err = f.includeRepo.Update(ctx,&includes)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(inc.Id),
		Message: "Success Update Include",
	}
	return &result,nil
}

func (f includeUsecase) Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	err = f.includeRepo.Delete(ctx,id,currentUser.Name)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete Include",
	}
	return &result,nil
}