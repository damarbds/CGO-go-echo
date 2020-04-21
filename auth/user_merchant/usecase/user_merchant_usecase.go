package usecase

import (
	"context"
	"github.com/auth/admin"
	"github.com/auth/user_merchant"
	"math"
	"time"

	"github.com/auth/identityserver"
	"github.com/auth/merchant"
	"github.com/models"
)

type userMerchantUsecase struct {
	adminUsecase     admin.Usecase
	userMerchantRepo     user_merchant.Repository
	merchantUsecase 	merchant.Usecase
	identityServerUc identityserver.Usecase
	contextTimeout   time.Duration
}

// NewmerchantUsecase will create new an merchantUsecase object representation of merchant.Usecase interface
func NewuserMerchantUsecase(a user_merchant.Repository,  m merchant.Usecase,is identityserver.Usecase, adm admin.Usecase,timeout time.Duration) user_merchant.Usecase {
	return &userMerchantUsecase{
		adminUsecase:adm,
		userMerchantRepo :a,
		merchantUsecase:m,
		identityServerUc: is,
		contextTimeout:   timeout,
	}
}


func (m userMerchantUsecase) Delete(c context.Context, userId string, token string) (*models.ResponseDelete, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	currentUserAdmin, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}
	error := m.userMerchantRepo.Delete(ctx, userId, currentUserAdmin.Name)
	if error != nil {
		response := models.ResponseDelete{
			Id:      userId,
			Message: error.Error(),
		}
		return &response, nil
	}
	response := models.ResponseDelete{
		Id:      userId,
		Message: "Deleted Success",
	}

	return &response, nil
}

func (m userMerchantUsecase) Update(c context.Context, ar *models.NewCommandUserMerchant, isAdmin bool, token string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	var currentUser string
	if isAdmin == true {
		currentUserAdmin, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
		if err != nil {
			return err
		}
		currentUser = currentUserAdmin.Name
	}else {
		currentUsers, err := m.merchantUsecase.ValidateTokenMerchant(ctx, token)
		if err != nil {
			return err
		}
		currentUser = currentUsers.MerchantEmail
	}

	updateUser := models.RegisterAndUpdateUser{
		Id:            ar.Id,
		Username:      ar.Email,
		Password:      ar.Password,
		Name:          ar.FullName,
		GivenName:     "",
		FamilyName:    "",
		Email:         ar.Email,
		EmailVerified: false,
		Website:       "",
		Address:       "",
		OTP:           "",
		UserType:      2,
		PhoneNumber:  ar.PhoneNumber,
		UserRoles:nil,
	}
	_, err := m.identityServerUc.UpdateUser(&updateUser)
	if err != nil {
		return err
	}

	merchant := models.UserMerchant{}
	merchant.Id = ar.Id
	merchant.ModifiedBy = &currentUser
	merchant.FullName = ar.FullName
	merchant.Email = ar.Email
	merchant.PhoneNumber = ar.PhoneNumber
	merchant.MerchantId = ar.MerchantId
	return m.userMerchantRepo.Update(ctx, &merchant)
}

func (m userMerchantUsecase) Create(c context.Context, ar *models.NewCommandUserMerchant, token string) (*models.NewCommandUserMerchant, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	currentUserAdmin, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil,models.ErrUnAuthorize
	}
	existedMerchant, _ := m.userMerchantRepo.GetByUserEmail(ctx, ar.Email)
	if existedMerchant != nil {
		return nil,models.ErrConflict
	}
	//var roles []string
	registerUser := models.RegisterAndUpdateUser{
		Id:            "",
		Username:      ar.Email,
		Password:      ar.Password,
		Name:          ar.FullName,
		GivenName:     "",
		FamilyName:    "",
		Email:         ar.Email,
		EmailVerified: false,
		Website:       "",
		Address:       "",
		OTP:           "",
		UserType:      2,
		PhoneNumber: ar.PhoneNumber,
		UserRoles:nil,
	}
	isUser, errorIs := m.identityServerUc.CreateUser(&registerUser)
	ar.Id = isUser.Id
	if errorIs != nil {
		return nil,errorIs
	}
	merchant := models.UserMerchant{}
	merchant.Id = isUser.Id
	merchant.CreatedBy = currentUserAdmin.Name
	merchant.FullName = ar.FullName
	merchant.Email = ar.Email
	merchant.PhoneNumber = ar.PhoneNumber
	merchant.MerchantId = ar.MerchantId
	err = m.userMerchantRepo.Insert(ctx, &merchant)
	if err != nil {
		return nil,err
	}

	return ar,nil
}

func (m userMerchantUsecase) List(ctx context.Context, page, limit, offset int, search string,token string) (*models.UserMerchantWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	_, err := m.adminUsecase.ValidateTokenAdmin(ctx, token)
	if err != nil {
		return nil, err
	}

	list, err := m.userMerchantRepo.List(ctx, limit, offset,search)
	if err != nil {
		return nil, err
	}

	merchants := make([]*models.UserMerchantInfoDto, len(list))
	for i, item := range list {
		merchants[i] = &models.UserMerchantInfoDto{
			Id:            item.Id,
			CreatedDate:   item.CreatedDate,
			UpdatedDate:   item.ModifiedDate,
			IsActive:      item.IsActive,
			FullName:item.FullName,
			Email:item.Email,
			PhoneNumber:item.PhoneNumber,
			MerchantId:item.MerchantId,
		}
	}
	totalRecords, _ := m.userMerchantRepo.Count(ctx)
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

	response := &models.UserMerchantWithPagination{
		Data: merchants,
		Meta: meta,
	}

	return response, nil
}

func (m userMerchantUsecase) GetUserDetailById(c context.Context, id string, token string) (*models.UserMerchantDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	getUserIdentity ,err := m.identityServerUc.GetDetailUserById(id,token)
	if err != nil {
		return nil,models.ErrNotFound
	}
	getMerchant , err := m.userMerchantRepo.GetByID(ctx,id)

	result := models.UserMerchantDto{
		Id:            getMerchant.Id,
		FullName:  getMerchant.FullName,
		Email: getMerchant.Email,
		Password:      getUserIdentity.Password,
		PhoneNumber:   getMerchant.PhoneNumber,
		MerchantId:getMerchant.MerchantId,
	}

	return &result,nil
}
