package usecase

import (
	"math"
	"time"

	"github.com/auth/admin"
	"github.com/models"
	"github.com/service/promo"
	"github.com/service/promo_merchant"
	"golang.org/x/net/context"
)

type promoUsecase struct {
	promoMerchant  promo_merchant.Repository
	adminUsecase   admin.Usecase
	promoRepo      promo.Repository
	contextTimeout time.Duration
}

// NewPromoUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPromoUsecase(pm promo_merchant.Repository, p promo.Repository, au admin.Usecase, timeout time.Duration) promo.Usecase {
	return &promoUsecase{
		promoMerchant:  pm,
		promoRepo:      p,
		adminUsecase:   au,
		contextTimeout: timeout,
	}
}

func (m promoUsecase) List(ctx context.Context, page, limit, offset int, search string, token string) (*models.PromoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	_, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	list, err := m.promoRepo.Fetch(ctx, &offset, &limit, search)
	if err != nil {
		return nil, err
	}

	promos := make([]*models.PromoDto, len(list))
	for i, item := range list {
		promos[i] = &models.PromoDto{
			Id:                 item.Id,
			PromoCode:          item.PromoCode,
			PromoName:          item.PromoName,
			PromoDesc:          item.PromoDesc,
			PromoValue:         item.PromoValue,
			PromoType:          item.PromoType,
			PromoImage:         item.PromoImage,
			StartDate:          item.StartDate,
			EndDate:            item.EndDate,
			Currency:           item.CurrencyId,
			MaxUsage:           item.MaxUsage,
			ProductionCapacity: item.ProductionCapacity,
			//VoucherValueOptionType: item.VoucherValueOptionType,
		}
		merchantIds := make([]string, 0)
		getPromoMerchant, err := m.promoMerchant.GetByMerchantId(ctx, "", item.Id)
		if err != nil {
			return nil, err
		}
		for _, element := range getPromoMerchant {
			merchantIds = append(merchantIds, element.MerchantId)
		}
		promos[i].MerchantId = merchantIds
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

	currentUser, err := p.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}
	promo := models.Promo{
		Id:                 command.Id,
		CreatedBy:          "",
		CreatedDate:        time.Now(),
		ModifiedBy:         &currentUser.Name,
		ModifiedDate:       nil,
		DeletedBy:          nil,
		DeletedDate:        nil,
		IsDeleted:          0,
		IsActive:           0,
		PromoCode:          command.PromoCode,
		PromoName:          command.PromoName,
		PromoDesc:          command.PromoDesc,
		PromoValue:         command.PromoValue,
		PromoType:          command.PromoType,
		PromoImage:         command.PromoImage,
		StartDate:          &command.StartDate,
		EndDate:            &command.EndDate,
		CurrencyId:         &command.Currency,
		MaxUsage:           &command.MaxUsage,
		ProductionCapacity: &command.ProductionCapacity,
		//VoucherValueOptionType: &command.VoucherValueOptionType,
	}
	err = p.promoRepo.Update(ctx, &promo)
	for _, element := range command.MerchantId {

		err = p.promoMerchant.DeleteByMerchantId(ctx, element, command.Id)
		promoMerchant := models.PromoMerchant{
			Id:         0,
			PromoId:    command.Id,
			MerchantId: element,
		}
		err := p.promoMerchant.Insert(ctx, promoMerchant)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	return &command, nil
}

func (p promoUsecase) Create(ctx context.Context, command models.NewCommandPromo, token string) (*models.NewCommandPromo, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	currentUser, err := p.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}
	promo := models.Promo{
		Id:                 "",
		CreatedBy:          currentUser.Name,
		CreatedDate:        time.Now(),
		ModifiedBy:         nil,
		ModifiedDate:       nil,
		DeletedBy:          nil,
		DeletedDate:        nil,
		IsDeleted:          0,
		IsActive:           0,
		PromoCode:          command.PromoCode,
		PromoName:          command.PromoName,
		PromoDesc:          command.PromoDesc,
		PromoValue:         command.PromoValue,
		PromoType:          command.PromoType,
		PromoImage:         command.PromoImage,
		StartDate:          &command.StartDate,
		EndDate:            &command.EndDate,
		CurrencyId:         &command.Currency,
		MaxUsage:           &command.MaxUsage,
		ProductionCapacity: &command.ProductionCapacity,
	}
	id, err := p.promoRepo.Insert(ctx, &promo)

	for _, element := range command.MerchantId {
		promoMerchant := models.PromoMerchant{
			Id:         0,
			PromoId:    id,
			MerchantId: element,
		}
		err := p.promoMerchant.Insert(ctx, promoMerchant)
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}
	command.Id = id
	return &command, nil
}

func (p promoUsecase) Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	currentUser, err := p.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}
	error := p.promoRepo.Delete(ctx, id, currentUser.Name)
	if error != nil {
		return nil, models.ErrNotFound
	}
	result := models.ResponseDelete{
		Id:      id,
		Message: "Success Deleted",
	}
	return &result, nil
}

func (p promoUsecase) GetDetail(ctx context.Context, id string, token string) (*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	_, err := p.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	getPromoDetail, err := p.promoRepo.GetById(ctx, id)
	if err != nil {
		return nil, models.ErrNotFound
	}
	result := models.PromoDto{
		Id:                 getPromoDetail.Id,
		PromoCode:          getPromoDetail.PromoCode,
		PromoName:          getPromoDetail.PromoName,
		PromoDesc:          getPromoDetail.PromoDesc,
		PromoValue:         getPromoDetail.PromoValue,
		PromoType:          getPromoDetail.PromoType,
		PromoImage:         getPromoDetail.PromoImage,
		StartDate:          getPromoDetail.StartDate,
		EndDate:            getPromoDetail.EndDate,
		Currency:           getPromoDetail.CurrencyId,
		MaxUsage:           getPromoDetail.MaxUsage,
		ProductionCapacity: getPromoDetail.ProductionCapacity,

		//VoucherValueOptionType: getPromoDetail.VoucherValueOptionType,
	}
	merchantIds := make([]string, 0)
	getPromoMerchant, err := p.promoMerchant.GetByMerchantId(ctx, "", getPromoDetail.Id)
	for _, element := range getPromoMerchant {
		merchantIds = append(merchantIds, element.MerchantId)
	}
	result.MerchantId = merchantIds

	return &result, nil
}
func (p promoUsecase) Fetch(ctx context.Context, page *int, size *int) ([]*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	promoList, err := p.promoRepo.Fetch(ctx, page, size, "")
	if err != nil {
		return nil, err
	}
	var promoDto []*models.PromoDto
	for _, element := range promoList {
		resPromo := models.PromoDto{
			Id:                 element.Id,
			PromoCode:          element.PromoCode,
			PromoName:          element.PromoName,
			PromoDesc:          element.PromoDesc,
			PromoValue:         element.PromoValue,
			PromoType:          element.PromoType,
			PromoImage:         element.PromoImage,
			StartDate:          element.StartDate,
			EndDate:            element.EndDate,
			Currency:           element.CurrencyId,
			MaxUsage:           element.MaxUsage,
			ProductionCapacity: element.ProductionCapacity,
		}
		merchantIds := make([]string, 0)
		getPromoMerchant, err := p.promoMerchant.GetByMerchantId(ctx, "", element.Id)
		if err != nil {
			return nil, err
		}
		for _, element := range getPromoMerchant {
			merchantIds = append(merchantIds, element.MerchantId)
		}
		resPromo.MerchantId = merchantIds
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
