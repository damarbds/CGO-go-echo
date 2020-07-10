package usecase

import (
	"context"
	"github.com/auth/admin"
	"github.com/service/currency"
	"math"
	"strconv"
	"time"

	"github.com/models"
)

type currencyUsecase struct {
	adminUsecase   admin.Usecase
	currencyRepo   currency.Repository
	contextTimeout time.Duration
}


func NewCurrencyUsecase(adminUsecase admin.Usecase,currencyRepo currency.Repository, timeout time.Duration) currency.Usecase {
	return &currencyUsecase{
		adminUsecase:adminUsecase,
		currencyRepo:   currencyRepo,
		contextTimeout: timeout,
	}
}


func (f currencyUsecase) GetAll(ctx context.Context, page, limit, offset int) (*models.CurrencyDtoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getCurrency,err := f.currencyRepo.Fetch(ctx,limit,offset)
	if err != nil {
		return nil,err
	}

	currencyDtos := make( []*models.CurrencyDto,0)

	for _,element := range getCurrency{
		dto := models.CurrencyDto{
			Id:     element.Id,
			Code:   element.Code,
			Name:   element.Name,
			Symbol: element.Symbol,
		}
		currencyDtos = append(currencyDtos,&dto)
	}
	totalRecords, _ := f.currencyRepo.GetCount(ctx)

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
		RecordPerPage: len(currencyDtos),
	}

	response := &models.CurrencyDtoWithPagination{
		Data: currencyDtos,
		Meta: meta,
	}
	return response, nil
}

func (f currencyUsecase) GetById(ctx context.Context, id int) (*models.CurrencyDto, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	getById,err := f.currencyRepo.GetById(ctx,id)
	if err != nil {
		return nil,err
	}
	result := models.CurrencyDto{
		Id:     getById.Id,
		Code:   getById.Code,
		Name:   getById.Name,
		Symbol: getById.Symbol,
	}
	return &result,nil

}

func (f currencyUsecase) Create(ctx context.Context, fac *models.NewCommandCurrency, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	currencyMod := models.Currency{
		Id:           0,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		Code:         fac.Code,
		Name:         fac.Name,
		Symbol:       fac.Symbol,
	}
	id ,err := f.currencyRepo.Insert(ctx,&currencyMod)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(*id),
		Message: "Success Create Currency",
	}
	return &result,nil
}

func (f currencyUsecase) Update(ctx context.Context, fac *models.NewCommandCurrency, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	currencyMod := models.Currency{
		Id:           fac.Id,
		CreatedBy:    currentUser.Name,
		CreatedDate:  time.Time{},
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		Code:         fac.Code,
		Name:         fac.Name,
		Symbol:       fac.Symbol,
	}
	err = f.currencyRepo.Update(ctx,&currencyMod)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(fac.Id),
		Message: "Success Update Currency",
	}
	return &result,nil
}

func (f currencyUsecase) Delete(ctx context.Context, id int, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, f.contextTimeout)
	defer cancel()

	currentUser, err := f.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	err = f.currencyRepo.Delete(ctx,id,currentUser.Name)
	if err != nil {
		return nil,err
	}

	result := models.ResponseDelete{
		Id:      strconv.Itoa(id),
		Message: "Success Delete Currency",
	}
	return &result,nil
}
