package usecase

import (
	"github.com/auth/admin"
	"github.com/models"
	"github.com/service/promo"
	"golang.org/x/net/context"
	"math"
	"time"
)

type promoUsecase struct {
	adminUsecase 	admin.Usecase
	promoRepo      promo.Repository
	contextTimeout time.Duration
}


// NewPromoUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPromoUsecase(p promo.Repository, au admin.Usecase,timeout time.Duration) promo.Usecase {
	return &promoUsecase{
		promoRepo:      p,
		adminUsecase:au,
		contextTimeout: timeout,
	}
}

func (m promoUsecase) List(ctx context.Context, page, limit, offset int, search string,token string) (*models.PromoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	_, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	list, err := m.promoRepo.Fetch(ctx, &offset, &limit,search)
	if err != nil {
		return nil, err
	}

	promos := make([]*models.PromoDto, len(list))
	for i, item := range list {
		promos[i] = &models.PromoDto{
			Id:                     item.Id,
			PromoCode:              item.PromoCode,
			PromoName:              item.PromoName,
			PromoDesc:              item.PromoDesc,
			PromoValue:             item.PromoValue,
			PromoType:              item.PromoType,
			PromoImage:             item.PromoImage,
			StartDate:              item.StartDate,
			EndDate:                item.EndDate,
			Currency:               item.Currency,
			MaxUsage:               item.MaxUsage,
			VoucherValueOptionType: item.VoucherValueOptionType,
		}
	}
	totalRecords, _ := m.promoRepo.GetCount(ctx)
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
		RecordPerPage: len(list),
	}

	response := &models.PromoWithPagination{
		Data: promos,
		Meta: meta,
	}

	return response, nil
}
func (p promoUsecase) Update(ctx context.Context, command models.NewCommandPromo, token string) (*models.NewCommandPromo, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	currentUser ,err := p.adminUsecase.ValidateTokenAdmin(ctx,token)
	if err != nil {
		return nil,models.ErrUnAuthorize
	}
	promo := models.Promo{
		Id:                     command.Id,
		CreatedBy:              "",
		CreatedDate:            time.Now(),
		ModifiedBy:             &currentUser.Name,
		ModifiedDate:           nil,
		DeletedBy:              nil,
		DeletedDate:            nil,
		IsDeleted:              0,
		IsActive:               0,
		PromoCode:              command.PromoCode,
		PromoName:              command.PromoName,
		PromoDesc:              command.PromoDesc,
		PromoValue:             command.PromoValue,
		PromoType:              command.PromoType,
		PromoImage:             command.PromoImage,
		StartDate:              &command.StartDate,
		EndDate:                &command.EndDate,
		Currency:               &command.Currency,
		MaxUsage:               &command.MaxUsage,
		VoucherValueOptionType: &command.VoucherValueOptionType,
	}
	err = p.promoRepo.Update(ctx,&promo)
	if err != nil {
		return nil,err
	}
	return &command,nil
}

func (p promoUsecase) Create(ctx context.Context, command models.NewCommandPromo, token string) (*models.NewCommandPromo, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	currentUser ,err := p.adminUsecase.ValidateTokenAdmin(ctx,token)
	if err != nil {
		return nil,models.ErrUnAuthorize
	}
	promo := models.Promo{
		Id:                     "",
		CreatedBy:              currentUser.Name,
		CreatedDate:            time.Now(),
		ModifiedBy:             nil,
		ModifiedDate:           nil,
		DeletedBy:              nil,
		DeletedDate:            nil,
		IsDeleted:              0,
		IsActive:               0,
		PromoCode:              command.PromoCode,
		PromoName:              command.PromoName,
		PromoDesc:              command.PromoDesc,
		PromoValue:             command.PromoValue,
		PromoType:              command.PromoType,
		PromoImage:             command.PromoImage,
		StartDate:              &command.StartDate,
		EndDate:                &command.EndDate,
		Currency:               &command.Currency,
		MaxUsage:               &command.MaxUsage,
		VoucherValueOptionType: &command.VoucherValueOptionType,
	}
	id,err := p.promoRepo.Insert(ctx,&promo)
	if err != nil {
		return nil,err
	}
	command.Id = id
	return &command,nil
}

func (p promoUsecase) Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	currentUser ,err := p.adminUsecase.ValidateTokenAdmin(ctx,token)
	if err != nil {
		return nil,models.ErrUnAuthorize
	}
	error := p.promoRepo.Delete(ctx,id,currentUser.Name)
	if error != nil {
		return nil,models.ErrNotFound
	}
	result := models.ResponseDelete{
		Id:      id,
		Message: "Success Deleted",
	}
	return &result,nil
}

func (p promoUsecase) GetDetail(ctx context.Context, id string, token string) (*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	_ ,err := p.adminUsecase.ValidateTokenAdmin(ctx,token)
	if err != nil {
		return nil,models.ErrUnAuthorize
	}

	getPromoDetail ,err := p.promoRepo.GetById(ctx,id)
	if err != nil {
		return nil,models.ErrNotFound
	}
	result := models.PromoDto{
		Id:                     getPromoDetail.Id,
		PromoCode:              getPromoDetail.PromoCode,
		PromoName:              getPromoDetail.PromoName,
		PromoDesc:              getPromoDetail.PromoDesc,
		PromoValue:             getPromoDetail.PromoValue,
		PromoType:              getPromoDetail.PromoType,
		PromoImage:             getPromoDetail.PromoImage,
		StartDate:              getPromoDetail.StartDate,
		EndDate:                getPromoDetail.EndDate,
		Currency:               getPromoDetail.Currency,
		MaxUsage:               getPromoDetail.MaxUsage,
		VoucherValueOptionType: getPromoDetail.VoucherValueOptionType,
	}

	return &result,nil
}
func (p promoUsecase) Fetch(ctx context.Context, page *int, size *int) ([]*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	promoList, err := p.promoRepo.Fetch(ctx, page, size,"")
	if err != nil {
		return nil, err
	}
	var promoDto []*models.PromoDto
	for _, element := range promoList {
		resPromo := models.PromoDto{
			Id:         element.Id,
			PromoCode:  element.PromoCode,
			PromoName:  element.PromoName,
			PromoDesc:  element.PromoDesc,
			PromoValue: element.PromoValue,
			PromoType:  element.PromoType,
			PromoImage: element.PromoImage,
		}
		promoDto = append(promoDto, &resPromo)
	}

	return promoDto, nil
}

func (p promoUsecase) GetByCode(ctx context.Context, code string) (*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	promos, err := p.promoRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	promoDto := make([]*models.PromoDto, len(promos))
	for i, p := range promos {
		promoDto[i] = &models.PromoDto{
			Id:         p.Id,
			PromoCode:  p.PromoCode,
			PromoName:  p.PromoName,
			PromoDesc:  p.PromoDesc,
			PromoValue: p.PromoValue,
			PromoType:  p.PromoType,
			PromoImage: p.PromoImage,
		}
	}

	return promoDto[0], nil
}
