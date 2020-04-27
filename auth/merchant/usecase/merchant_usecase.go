package usecase

import (
	"context"
	"math"
	"time"

	"github.com/auth/admin"
	"github.com/auth/user_merchant"

	"github.com/auth/identityserver"
	"github.com/service/experience"
	"github.com/service/transportation"

	"github.com/auth/merchant"
	"github.com/models"
)

type merchantUsecase struct {
	userMerchantRepo user_merchant.Repository
	adminUsecase     admin.Usecase
	merchantRepo     merchant.Repository
	expRepo          experience.Repository
	transRepo        transportation.Repository
	identityServerUc identityserver.Usecase
	contextTimeout   time.Duration
}

// NewmerchantUsecase will create new an merchantUsecase object representation of merchant.Usecase interface
func NewmerchantUsecase(usm user_merchant.Repository, a merchant.Repository, ex experience.Repository, tr transportation.Repository, is identityserver.Usecase, adm admin.Usecase, timeout time.Duration) merchant.Usecase {
	return &merchantUsecase{
		userMerchantRepo: usm,
		merchantRepo:     a,
		expRepo:          ex,
		transRepo:        tr,
		identityServerUc: is,
		adminUsecase:     adm,
		contextTimeout:   timeout,
	}
}

func (m merchantUsecase) ServiceCount(ctx context.Context, token string) (*models.ServiceCount, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}

	existedMerchant, _ := m.merchantRepo.GetByMerchantEmail(ctx, getInfoToIs.Email)
	if existedMerchant == nil {
		return nil, models.ErrNotFound
	}

	expCount, err := m.expRepo.GetExpCount(ctx, existedMerchant.Id)
	if err != nil {
		return nil, err
	}

	transCount, err := m.transRepo.GetTransCount(ctx, existedMerchant.Id)
	if err != nil {
		return nil, err
	}

	response := &models.ServiceCount{
		ExpCount:   expCount,
		TransCount: transCount,
	}

	return response, nil
}

func (m merchantUsecase) List(ctx context.Context, page, limit, offset int, token string, search string) (*models.MerchantWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	_, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	list, err := m.merchantRepo.List(ctx, limit, offset, search)
	if err != nil {
		return nil, err
	}

	merchants := make([]*models.MerchantInfoDto, len(list))
	for i, item := range list {
		merchants[i] = &models.MerchantInfoDto{
			Id:            item.Id,
			CreatedDate:   item.CreatedDate,
			UpdatedDate:   item.ModifiedDate,
			IsActive:      item.IsActive,
			MerchantName:  item.MerchantName,
			MerchantDesc:  item.MerchantDesc,
			MerchantEmail: item.MerchantEmail,
			Balance:       item.Balance,
			PhoneNumber:   item.PhoneNumber,
		}
	}
	totalRecords, _ := m.merchantRepo.Count(ctx)
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

	response := &models.MerchantWithPagination{
		Data: merchants,
		Meta: meta,
	}

	return response, nil
}

func (m merchantUsecase) Count(ctx context.Context) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	count, err := m.merchantRepo.Count(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}

func (m merchantUsecase) Login(ctx context.Context, ar *models.Login) (*models.GetToken, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	requestToken, err := m.identityServerUc.GetToken(ar.Email, ar.Password, ar.Scope)
	if err != nil {
		return nil, err
	}
	existedMerchant, _ := m.userMerchantRepo.GetByUserEmail(ctx, ar.Email)
	if existedMerchant == nil {
		return nil, models.ErrNotFound
	}
	return requestToken, err
}

func (m merchantUsecase) ValidateTokenMerchant(ctx context.Context, token string) (*models.MerchantInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	existedMerchant, _ := m.userMerchantRepo.GetByUserEmail(ctx, getInfoToIs.Email)
	if existedMerchant == nil {
		return nil, models.ErrNotFound
	}
	getMerchant,_ := m.merchantRepo.GetByID(ctx,existedMerchant.MerchantId)
	merchantInfo := models.MerchantInfoDto{
		Id:            existedMerchant.MerchantId,
		MerchantName:  existedMerchant.FullName,
		MerchantDesc:  getMerchant.MerchantDesc,
		MerchantEmail: existedMerchant.Email,
		Balance:       getMerchant.Balance,
	}

	return &merchantInfo, nil
}

func (m merchantUsecase) GetMerchantInfo(ctx context.Context, token string) (*models.MerchantInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	existedMerchant, _ := m.userMerchantRepo.GetByUserEmail(ctx, getInfoToIs.Email)
	if existedMerchant == nil {
		return nil, models.ErrNotFound
	}
	getMerchant,_ := m.merchantRepo.GetByID(ctx,existedMerchant.MerchantId)
	merchantInfo := models.MerchantInfoDto{
		Id:            existedMerchant.MerchantId,
		MerchantName:  existedMerchant.FullName,
		MerchantDesc:  getMerchant.MerchantDesc,
		MerchantEmail: existedMerchant.Email,
		Balance:       getMerchant.Balance,
	}

	return &merchantInfo, nil
}

func (m merchantUsecase) Update(c context.Context, ar *models.NewCommandMerchant, isAdmin bool, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	var currentUser string
	if isAdmin == true {
		currentUserAdmin, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
		if err != nil {
			return err
		}
		currentUser = currentUserAdmin.Name
	} else {
		currentUsers, err := m.ValidateTokenMerchant(ctx, token)
		if err != nil {
			return err
		}
		currentUser = currentUsers.MerchantEmail
	}

	updateUser := models.RegisterAndUpdateUser{
		Id:            ar.Id,
		Username:      ar.MerchantEmail,
		Password:      ar.MerchantPassword,
		Name:          ar.MerchantName,
		GivenName:     "",
		FamilyName:    "",
		Email:         ar.MerchantEmail,
		EmailVerified: false,
		Website:       "",
		Address:       "",
		OTP:           "",
		UserType:      2,
		PhoneNumber:   ar.PhoneNumber,
		UserRoles:     nil,
	}
	_, err := m.identityServerUc.UpdateUser(&updateUser)
	if err != nil {
		return err
	}

	merchant := models.Merchant{}
	merchant.Id = ar.Id
	merchant.ModifiedBy = &currentUser
	merchant.MerchantName = ar.MerchantName
	merchant.MerchantDesc = ar.MerchantDesc
	merchant.MerchantEmail = ar.MerchantEmail
	merchant.Balance = ar.Balance
	merchant.PhoneNumber = &ar.PhoneNumber
	merchant.MerchantPicture = ar.MerchantPicture
	return m.merchantRepo.Update(ctx, &merchant)
}

func (m merchantUsecase) Create(c context.Context, ar *models.NewCommandMerchant, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	currentUserAdmin, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return err
	}
	existedMerchant, _ := m.merchantRepo.GetByMerchantEmail(ctx, ar.MerchantEmail)
	if existedMerchant != nil {
		return models.ErrConflict
	}
	//var roles []string
	registerUser := models.RegisterAndUpdateUser{
		Id:            "",
		Username:      ar.MerchantEmail,
		Password:      ar.MerchantPassword,
		Name:          ar.MerchantName,
		GivenName:     "",
		FamilyName:    "",
		Email:         ar.MerchantEmail,
		EmailVerified: false,
		Website:       "",
		Address:       "",
		OTP:           "",
		UserType:      2,
		PhoneNumber:   ar.PhoneNumber,
		UserRoles:     nil,
	}
	isUser, errorIs := m.identityServerUc.CreateUser(&registerUser)
	ar.Id = isUser.Id
	if errorIs != nil {
		return errorIs
	}
	merchant := models.Merchant{}
	merchant.Id = isUser.Id
	merchant.CreatedBy = currentUserAdmin.Name
	merchant.MerchantName = ar.MerchantName
	merchant.MerchantDesc = ar.MerchantDesc
	merchant.MerchantEmail = ar.MerchantEmail
	merchant.Balance = ar.Balance
	merchant.PhoneNumber = &ar.PhoneNumber
	merchant.MerchantPicture = ar.MerchantPicture
	err = m.merchantRepo.Insert(ctx, &merchant)
	if err != nil {
		return err
	}

	return nil
}
func (m merchantUsecase) Delete(c context.Context, id string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	currentUserAdmin, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	error := m.merchantRepo.Delete(ctx, id, currentUserAdmin.Name)
	_ = m.identityServerUc.DeleteUser(id)
	if error != nil {
		response := models.ResponseDelete{
			Id:      id,
			Message: error.Error(),
		}
		return &response, nil
	}
	response := models.ResponseDelete{
		Id:      id,
		Message: "Deleted Success",
	}

	return &response, nil
}

func (m merchantUsecase) GetDetailMerchantById(c context.Context, id string, token string) (*models.MerchantDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	getUserIdentity, err := m.identityServerUc.GetDetailUserById(id, token, "true")
	if err != nil {
		return nil, err
	}
	getMerchant, err := m.merchantRepo.GetByID(ctx, id)

	result := models.MerchantDto{
		Id:              getMerchant.Id,
		CreatedDate:     getMerchant.CreatedDate,
		UpdatedDate:     getMerchant.ModifiedDate,
		IsActive:        getMerchant.IsActive,
		MerchantName:    getMerchant.MerchantName,
		MerchantDesc:    getMerchant.MerchantDesc,
		MerchantEmail:   getMerchant.MerchantEmail,
		Password:        getUserIdentity.Password,
		Balance:         getMerchant.Balance,
		PhoneNumber:     getMerchant.PhoneNumber,
		MerchantPicture: getMerchant.MerchantPicture,
	}

	return &result, nil
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
