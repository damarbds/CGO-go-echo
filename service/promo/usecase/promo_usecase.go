package usecase

import (
	"math"
	"time"

	"github.com/auth/user"
	"github.com/booking/booking_exp"
	guuid "github.com/google/uuid"
	"github.com/service/promo_experience_transport"
	"github.com/service/promo_user"
	"github.com/transactions/transaction"

	"github.com/auth/admin"
	"github.com/models"
	"github.com/service/promo"
	"github.com/service/promo_merchant"
	"golang.org/x/net/context"
)

type promoUsecase struct {
	bookingRepo              booking_exp.Repository
	userUsecase              user.Usecase
	promoMerchant            promo_merchant.Repository
	promoUser                promo_user.Repository
	promoExperienceTransport promo_experience_transport.Repository
	adminUsecase             admin.Usecase
	promoRepo                promo.Repository
	contextTimeout           time.Duration
	transactionRepo          transaction.Repository
}

// NewPromoUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewPromoUsecase(bookingRepo booking_exp.Repository, userUsecase user.Usecase, transactionRepo transaction.Repository, pm promo_merchant.Repository, p promo.Repository, au admin.Usecase, timeout time.Duration, pu promo_user.Repository, pet promo_experience_transport.Repository) promo.Usecase {
	return &promoUsecase{
		bookingRepo:              bookingRepo,
		userUsecase:              userUsecase,
		transactionRepo:          transactionRepo,
		promoMerchant:            pm,
		promoRepo:                p,
		adminUsecase:             au,
		contextTimeout:           timeout,
		promoUser:                pu,
		promoExperienceTransport: pet,
	}
}

func (m promoUsecase) List(ctx context.Context, page, limit, offset int, search string, token string, trans bool, exp bool, merchantIds []string) (*models.PromoWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	_, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	list, err := m.promoRepo.Fetch(ctx, &offset, &limit, search, trans, exp, merchantIds, "", "")
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
			PromoProductType:   item.PromoProductType,
			StartTripPeriod:    item.StartTripPeriod,
			EndTripPeriod:      item.EndTripPeriod,
			IsAnyTripPeriod:    item.IsAnyTripPeriod,
			MaxDiscount:        item.MaxDiscount,
			HowToUse:           item.HowToUse,
			HowToGet:           item.HowToGet,
			Disclaimer:         item.Disclaimer,
			TermCondition:      item.TermCondition,
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
	if command.PromoImage == "" {
		getById, _ := p.promoRepo.GetById(ctx, command.Id)
		command.PromoImage = getById.PromoImage
	}
	now := time.Now()
	promo := models.Promo{
		Id:                 command.Id,
		CreatedBy:          "",
		CreatedDate:        time.Now(),
		ModifiedBy:         &currentUser.Name,
		ModifiedDate:       &now,
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
		PromoProductType:   command.PromoProductType,
		StartTripPeriod:    &command.StartTripPeriod,
		EndTripPeriod:      &command.EndTripPeriod,
		IsAnyTripPeriod:    &command.IsAnyTripPeriod,
		MaxDiscount:        &command.MaxDiscount,
		HowToUse:           &command.HowToUse,
		HowToGet:           &command.HowToGet,
		TermCondition:      &command.TermCondition,
		Disclaimer:         &command.Disclaimer,
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
	for _, element := range command.UserId {
		err = p.promoUser.DeleteByUserId(ctx, element, command.Id)
		promoUser := models.PromoUser{
			Id:      0,
			PromoId: command.Id,
			UserId:  element,
		}
		err := p.promoUser.Insert(ctx, promoUser)
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
		Id:                 guuid.New().String(),
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
		PromoProductType:   command.PromoProductType,
		StartTripPeriod:    &command.StartTripPeriod,
		EndTripPeriod:      &command.EndTripPeriod,
		IsAnyTripPeriod:    &command.IsAnyTripPeriod,
		MaxDiscount:        &command.MaxDiscount,
		HowToUse:           &command.HowToUse,
		HowToGet:           &command.HowToGet,
		TermCondition:      &command.TermCondition,
		Disclaimer:         &command.Disclaimer,
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
	for _, element := range command.UserId {
		promoUser := models.PromoUser{
			Id:      0,
			PromoId: id,
			UserId:  element,
		}
		err := p.promoUser.Insert(ctx, promoUser)
		if err != nil {
			return nil, err
		}
	}
	for _, element := range command.ExperienceId {
		promoExperienceTransport := models.PromoExperienceTransport{
			Id:               0,
			PromoId:          id,
			ExperienceId:     element,
			TransportationId: "",
		}
		//err := p.promoExperience.Insert(ctx, promoExperience)
		err := p.promoExperienceTransport.Insert(ctx, promoExperienceTransport)
		if err != nil {
			return nil, err
		}
	}
	for _, element := range command.TransportationId {
		promoExperienceTransport := models.PromoExperienceTransport{
			Id:               0,
			PromoId:          id,
			ExperienceId:     "",
			TransportationId: element,
		}
		//err := p.promoExperience.Insert(ctx, promoExperience)
		err := p.promoExperienceTransport.Insert(ctx, promoExperienceTransport)
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
		PromoProductType:   getPromoDetail.PromoProductType,
		StartTripPeriod:    getPromoDetail.StartTripPeriod,
		EndTripPeriod:      getPromoDetail.EndTripPeriod,
		IsAnyTripPeriod:    getPromoDetail.IsAnyTripPeriod,
		MaxDiscount:        getPromoDetail.MaxDiscount,
		HowToUse:           getPromoDetail.HowToUse,
		HowToGet:           getPromoDetail.HowToGet,
		TermCondition:      getPromoDetail.TermCondition,
		Disclaimer:         getPromoDetail.Disclaimer,
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
func (p promoUsecase) Fetch(ctx context.Context, page *int, size *int, search string, trans bool, exp bool, merchantIds []string, sortBy string, promoId string) ([]*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	//_, err := p.userUsecase.ValidateTokenUser(ctx, token)
	//if err != nil {
	//	return nil, models.ErrUnAuthorize
	//}

	promoList, err := p.promoRepo.Fetch(ctx, page, size, search, trans, exp, merchantIds, sortBy, promoId)
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
			PromoProductType:   element.PromoProductType,
			StartTripPeriod:    element.StartTripPeriod,
			EndTripPeriod:      element.EndTripPeriod,
			IsAnyTripPeriod:    element.IsAnyTripPeriod,
			MaxDiscount:        element.MaxDiscount,
			HowToUse:           element.HowToUse,
			HowToGet:           element.HowToGet,
			TermCondition:      element.TermCondition,
			Disclaimer:         element.Disclaimer,
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

func (p promoUsecase) FetchUser(ctx context.Context, page *int, size *int, token string, search string, trans bool, exp bool, merchantIds []string, sortBy string, promoId string) ([]*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()

	_, err := p.userUsecase.ValidateTokenUser(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	promoList, err := p.promoRepo.Fetch(ctx, page, size, search, trans, exp, merchantIds, sortBy, promoId)
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
			PromoProductType:   element.PromoProductType,
			StartTripPeriod:    element.StartTripPeriod,
			EndTripPeriod:      element.EndTripPeriod,
			IsAnyTripPeriod:    element.IsAnyTripPeriod,
			MaxDiscount:        element.MaxDiscount,
			HowToUse:           element.HowToUse,
			HowToGet:           element.HowToGet,
			TermCondition:      element.TermCondition,
			Disclaimer:         element.Disclaimer,
		}
		userIds := make([]string, 0)
		getPromoUser, err := p.promoUser.GetByUserId(ctx, "", element.Id)
		if err != nil {
			return nil, err
		}
		for _, element := range getPromoUser {
			userIds = append(userIds, element.UserId)
		}

		resPromo.UserId = userIds
		promoDto = append(promoDto, &resPromo)
	}

	return promoDto, nil
}

func (p promoUsecase) GetByCode(ctx context.Context, code string, promoType string, merchantId string, token string, bookingId string,isAdmin bool) (*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()
	var userId string
	if token != "" {
		if isAdmin == true{
			_, err := p.adminUsecase.ValidateTokenAdmin(ctx,token)
			if err != nil {
				return nil, models.ErrUnAuthorize
			}
			//userId = currentUser.Id
		}else {
			currentUser, err := p.userUsecase.ValidateTokenUser(ctx, token)
			if err != nil {
				return nil, models.ErrUnAuthorize
			}
			userId = currentUser.Id

		}
	}
	var expId string
	var transId string
	var bookingDate string
	var usePromoDate string
	booking, _ := p.bookingRepo.GetByBookingID(ctx, bookingId)
	if booking != nil {
		if booking.ExpId != nil {
			expId = *booking.ExpId
		}
		if booking.TransId != nil {
			transId = *booking.TransId
		}
		bookingDate = booking.BookingDate.Format("2006-01-02")
		usePromoDate = time.Now().Format("2006-01-02")
	}
	promos, err := p.promoRepo.GetByCode(ctx, code, promoType, merchantId, userId, expId,
		transId, bookingDate, usePromoDate)

	if err != nil {
		return nil, err
	}
	if userId != "" {
		countAlreadyUse, err := p.transactionRepo.GetCountTransactionByPromoId(ctx, promos[0].Id, "")
		if err != nil {
			return nil, err
		}
		var productCapacity int
		if promos[0].ProductionCapacity != nil {
			productCapacity = *promos[0].ProductionCapacity
		}
		count := productCapacity - countAlreadyUse
		if count < 1 {
			return nil, models.ErrNotFound
		}

		countAlreadyUseWithCurrentUser, err := p.transactionRepo.GetCountTransactionByPromoId(ctx, promos[0].Id, userId)
		if err != nil {
			return nil, err
		}
		var maxUsage int
		if promos[0].MaxUsage != nil {
			maxUsage = *promos[0].MaxUsage
		}
		countUsage := maxUsage - countAlreadyUseWithCurrentUser
		if countUsage < 1 {
			return nil, models.ErrNotFound
		}

	}
	//else {
	//	countAlreadyUse ,err := p.transactionRepo.GetCountTransactionByPromoId(ctx , promos[0].Id,"")
	//	if err != nil {
	//		return nil, err
	//	}
	//	var productCapacity int
	//	if promos[0].ProductionCapacity != nil {
	//		productCapacity = *promos[0].ProductionCapacity
	//	}
	//	count := productCapacity - countAlreadyUse
	//	if count < 1  {
	//		return nil,models.ErrNotFound
	//	}
	//}
	promoDto := &models.PromoDto{
		Id:                 promos[0].Id,
		PromoCode:          promos[0].PromoCode,
		PromoName:          promos[0].PromoName,
		PromoDesc:          promos[0].PromoDesc,
		PromoValue:         promos[0].PromoValue,
		PromoType:          promos[0].PromoType,
		PromoImage:         promos[0].PromoImage,
		StartDate:          promos[0].StartDate,
		EndDate:            promos[0].EndDate,
		Currency:           promos[0].CurrencyId,
		MaxUsage:           promos[0].MaxUsage,
		ProductionCapacity: promos[0].ProductionCapacity,
		PromoProductType:   promos[0].PromoProductType,
		StartTripPeriod:    promos[0].StartTripPeriod,
		EndTripPeriod:      promos[0].EndTripPeriod,
		IsAnyTripPeriod:    promos[0].IsAnyTripPeriod,
		MaxDiscount:        promos[0].MaxDiscount,
		HowToUse:           promos[0].HowToUse,
		HowToGet:           promos[0].HowToGet,
		TermCondition:      promos[0].TermCondition,
		Disclaimer:         promos[0].Disclaimer,
	}

	return promoDto, nil
}
func (p promoUsecase) GetByFilter(ctx context.Context, code string, promoType int, merchantExpId string, merchantTransportId string, token string) (*models.PromoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, p.contextTimeout)
	defer cancel()
	var userId string
	if token != "" {
		currentUser, err := p.userUsecase.ValidateTokenUser(ctx, token)
		if err != nil {
			return nil, models.ErrUnAuthorize
		}
		userId = currentUser.Id
	}
	promos, err := p.promoRepo.GetByFilter(ctx, code, &promoType, merchantExpId, merchantTransportId)
	if err != nil {
		return nil, err
	}
	if userId != "" {
		countAlreadyUse, err := p.transactionRepo.GetCountTransactionByPromoId(ctx, promos[0].Id, "")
		if err != nil {
			return nil, err
		}
		var productCapacity int
		if promos[0].ProductionCapacity != nil {
			productCapacity = *promos[0].ProductionCapacity
		}
		count := productCapacity - countAlreadyUse
		if count < 1 {
			return nil, models.ErrNotFound
		}

		countAlreadyUseWithCurrentUser, err := p.transactionRepo.GetCountTransactionByPromoId(ctx, promos[0].Id, userId)
		if err != nil {
			return nil, err
		}
		var maxUsage int
		if promos[0].MaxUsage != nil {
			maxUsage = *promos[0].MaxUsage
		}
		countUsage := maxUsage - countAlreadyUseWithCurrentUser
		if countUsage < 1 {
			return nil, models.ErrNotFound
		}

	} else {
		countAlreadyUse, err := p.transactionRepo.GetCountTransactionByPromoId(ctx, promos[0].Id, "")
		if err != nil {
			return nil, err
		}
		var productCapacity int
		if promos[0].ProductionCapacity != nil {
			productCapacity = *promos[0].ProductionCapacity
		}
		count := productCapacity - countAlreadyUse
		if count < 1 {
			return nil, models.ErrNotFound
		}
	}
	promoDto := &models.PromoDto{
		Id:                 promos[0].Id,
		PromoCode:          promos[0].PromoCode,
		PromoName:          promos[0].PromoName,
		PromoDesc:          promos[0].PromoDesc,
		PromoValue:         promos[0].PromoValue,
		PromoType:          promos[0].PromoType,
		PromoImage:         promos[0].PromoImage,
		StartDate:          promos[0].StartDate,
		EndDate:            promos[0].EndDate,
		Currency:           promos[0].CurrencyId,
		MaxUsage:           promos[0].MaxUsage,
		ProductionCapacity: promos[0].ProductionCapacity,
		PromoProductType:   promos[0].PromoProductType,
		StartTripPeriod:    promos[0].StartTripPeriod,
		EndTripPeriod:      promos[0].EndTripPeriod,
		IsAnyTripPeriod:    promos[0].IsAnyTripPeriod,
		MaxDiscount:        promos[0].MaxDiscount,
		HowToUse:           promos[0].HowToUse,
		HowToGet:           promos[0].HowToGet,
		TermCondition:      promos[0].TermCondition,
		Disclaimer:         promos[0].Disclaimer,
	}

	return promoDto, nil
}
