package usecase

import (
	"context"
	"github.com/identityserver"
	"time"

	"github.com/merchant"
	"github.com/models"
)

type merchantUsecase struct {
	merchantRepo   merchant.Repository
	identityServerUc identityserver.Usecase
	contextTimeout time.Duration
}

// NewmerchantUsecase will create new an merchantUsecase object representation of merchant.Usecase interface
func NewmerchantUsecase(a merchant.Repository, is identityserver.Usecase,timeout time.Duration) merchant.Usecase {
	return &merchantUsecase{
		merchantRepo:   a,
		identityServerUc:is,
		contextTimeout: timeout,
	}
}
func (m merchantUsecase) Login(ctx context.Context, ar *models.Login) (*models.GetToken, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	requestToken ,err:= m.identityServerUc.GetToken(ar.Email,ar.Password)
	if err != nil{
		return nil,err
	}
	existedMerchant, _ := m.merchantRepo.GetByMerchantEmail(ctx, ar.Email)
	if existedMerchant == nil {
		return nil,models.ErrNotFound
	}
	return  requestToken,err
}

func (m merchantUsecase) ValidateTokenMerchant(ctx context.Context, token string) (*string, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs ,err := m.identityServerUc.GetUserInfo(token)
	if err != nil{
		return nil,err
	}
	existedMerchant, _ := m.merchantRepo.GetByMerchantEmail(ctx, getInfoToIs.Email)
	if existedMerchant == nil {
		return nil,models.ErrNotFound
	}
	currentUser := getInfoToIs.Username
	return &currentUser,nil
}

func (m merchantUsecase) GetMerchantInfo(ctx context.Context, token string) (*models.MerchantInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs ,err := m.identityServerUc.GetUserInfo(token)
	if err != nil{
		return nil,err
	}
	existedMerchant, _ := m.merchantRepo.GetByMerchantEmail(ctx, getInfoToIs.Email)
	if existedMerchant == nil {
		return nil,models.ErrNotFound
	}
	merchantInfo := models.MerchantInfoDto{
		Id:            existedMerchant.Id,
		MerchantName:  existedMerchant.MerchantName,
		MerchantDesc:  existedMerchant.MerchantDesc,
		MerchantEmail: existedMerchant.MerchantEmail,
		Balance:       existedMerchant.Balance,
	}

	return &merchantInfo,nil
}


func (m merchantUsecase) Update(c context.Context, ar *models.NewCommandMerchant ,user string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

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
	}
	_ ,err:= m.identityServerUc.UpdateUser(&updateUser)
	if err != nil{
		return err
	}

	merchant := models.Merchant{}
	merchant.Id = ar.Id
	merchant.ModifiedBy = &user
	merchant.MerchantName = ar.MerchantName
	merchant.MerchantDesc = ar.MerchantDesc
	merchant.MerchantEmail = ar.MerchantEmail
	merchant.Balance = ar.Balance
	return m.merchantRepo.Update(ctx, &merchant)
}

func (m merchantUsecase) Create(c context.Context, ar *models.NewCommandMerchant,user string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	existedMerchant, _ := m.merchantRepo.GetByMerchantEmail(ctx, ar.MerchantEmail)
	if existedMerchant != nil {
		return models.ErrConflict
	}
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
	}
	isUser ,errorIs:= m.identityServerUc.CreateUser(&registerUser)
	ar.Id = isUser.Id
	if errorIs != nil{
		return errorIs
	}
	merchant := models.Merchant{}
	merchant.Id = isUser.Id
	merchant.CreatedBy = ar.MerchantEmail
	merchant.MerchantName = ar.MerchantName
	merchant.MerchantDesc = ar.MerchantDesc
	merchant.MerchantEmail = ar.MerchantEmail
	merchant.Balance = ar.Balance
	err := m.merchantRepo.Insert(ctx, &merchant)
	if err != nil {
		return err
	}


	return nil
}


/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
