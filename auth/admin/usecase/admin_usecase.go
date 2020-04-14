package usecase

import (
	"context"
	"time"

	"github.com/auth/identityserver"

	"github.com/auth/admin"
	"github.com/models"
)

type adminUsecase struct {
	adminRepo        admin.Repository
	identityServerUc identityserver.Usecase
	contextTimeout   time.Duration
}

// NewadminUsecase will create new an adminUsecase object representation of admin.Usecase interface
func NewadminUsecase(a admin.Repository, is identityserver.Usecase, timeout time.Duration) admin.Usecase {
	return &adminUsecase{
		adminRepo:        a,
		identityServerUc: is,
		contextTimeout:   timeout,
	}
}
func (m adminUsecase) Login(ctx context.Context, ar *models.Login) (*models.GetToken, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	requestToken, err := m.identityServerUc.GetToken(ar.Email, ar.Password,ar.Scope)
	if err != nil {
		return nil, err
	}
	existedadmin, _ := m.adminRepo.GetByAdminEmail(ctx, ar.Email)
	if existedadmin == nil {
		return nil, models.ErrNotFound
	}
	return requestToken, err
}

func (m adminUsecase) ValidateTokenAdmin(ctx context.Context, token string) (*models.AdminDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	existedadmin, _ := m.adminRepo.GetByAdminEmail(ctx, getInfoToIs.Email)
	if existedadmin == nil {
		return nil, models.ErrNotFound
	}
	adminInfo := models.AdminDto{
		Id:    existedadmin.Id,
		Name:  existedadmin.Name,
		Email: existedadmin.Email,
	}

	return &adminInfo, nil
}

func (m adminUsecase) GetAdminInfo(ctx context.Context, token string) (*models.AdminDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	existedadmin, _ := m.adminRepo.GetByAdminEmail(ctx, getInfoToIs.Email)
	if existedadmin == nil {
		return nil, models.ErrNotFound
	}
	adminInfo := models.AdminDto{
		Id:    existedadmin.Id,
		Name:  existedadmin.Name,
		Email: existedadmin.Email,
	}

	return &adminInfo, nil
}

func (m adminUsecase) Update(c context.Context, ar *models.NewCommandAdmin, user string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	updateUser := models.RegisterAndUpdateUser{
		Id:            ar.Id,
		Username:      ar.Email,
		Password:      ar.Password,
		Name:          ar.Name,
		GivenName:     "",
		FamilyName:    "",
		Email:         ar.Email,
		EmailVerified: false,
		Website:       "",
		Address:       "",
	}
	_, err := m.identityServerUc.UpdateUser(&updateUser)
	if err != nil {
		return err
	}

	admin := models.Admin{}
	admin.Id = ar.Id
	admin.ModifiedBy = &user
	admin.Name = ar.Name
	admin.Email = ar.Email
	return m.adminRepo.Update(ctx, &admin)
}

func (m adminUsecase) Create(c context.Context, ar *models.NewCommandAdmin, user string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	existedadmin, _ := m.adminRepo.GetByAdminEmail(ctx, ar.Email)
	if existedadmin != nil {
		return models.ErrConflict
	}
	registerUser := models.RegisterAndUpdateUser{
		Id:            "",
		Username:      ar.Email,
		Password:      ar.Password,
		Name:          ar.Name,
		GivenName:     "",
		FamilyName:    "",
		Email:         ar.Email,
		EmailVerified: false,
		Website:       "",
		Address:       "",
		OTP:           "",
		UserType:      3,
	}
	isUser, errorIs := m.identityServerUc.CreateUser(&registerUser)
	ar.Id = isUser.Id
	if errorIs != nil {
		return errorIs
	}
	admin := models.Admin{}
	admin.Id = isUser.Id
	admin.CreatedBy = ar.Email
	admin.Name = ar.Name
	admin.Email = ar.Email
	err := m.adminRepo.Insert(ctx, &admin)
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
