package usecase

import (
	"github.com/auth/identityserver"
	"github.com/auth/user"
	"github.com/models"
	"golang.org/x/net/context"
	"time"
)

type userUsecase struct {
	userRepo   user.Repository
	identityServerUc identityserver.Usecase
	contextTimeout time.Duration
}

// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewuserUsecase(a user.Repository, is identityserver.Usecase,timeout time.Duration) user.Usecase {
	return &userUsecase{
		userRepo:   a,
		identityServerUc:is,
		contextTimeout: timeout,
	}
}

func (m userUsecase) Login(ctx context.Context, ar *models.Login) (*models.GetToken, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	requestToken ,err:= m.identityServerUc.GetToken(ar.Email,ar.Password)
	if err != nil{
		return nil,err
	}
	existeduser, _ := m.userRepo.GetByUserEmail(ctx, ar.Email)
	if existeduser == nil {
		return nil,models.ErrNotFound
	}
	return  requestToken,err
}

func (m userUsecase) ValidateTokenUser(ctx context.Context, token string) (*models.UserInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs ,err := m.identityServerUc.GetUserInfo(token)
	if err != nil{
		return nil,err
	}
	existeduser, _ := m.userRepo.GetByUserEmail(ctx, getInfoToIs.Email)
	if existeduser == nil {
		return nil,models.ErrUnAuthorize
	}

	userInfo := &models.UserInfoDto{
		Id:             existeduser.Id,
		UserEmail:      existeduser.UserEmail,
		FullName:       existeduser.FullName,
		PhoneNumber:    existeduser.PhoneNumber,
		ProfilePictUrl: existeduser.ProfilePictUrl,
	}

	return userInfo ,nil
}

func (m userUsecase) GetUserInfo(ctx context.Context, token string) (*models.UserInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs ,err := m.identityServerUc.GetUserInfo(token)
	if err != nil{
		return nil,err
	}
	existeduser, _ := m.userRepo.GetByUserEmail(ctx, getInfoToIs.Email)
	if existeduser == nil {
		return nil,models.ErrNotFound
	}
	userInfo := models.UserInfoDto{
		Id:             existeduser.Id,
		UserEmail:      existeduser.UserEmail,
		FullName:       existeduser.FullName,
		PhoneNumber:    existeduser.PhoneNumber,
		ProfilePictUrl: existeduser.ProfilePictUrl,
	}

	return &userInfo,nil
}

func (m userUsecase) Update(c context.Context, ar *models.NewCommandUser ,user string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	updateUser := models.RegisterAndUpdateUser{
		Id:            ar.Id,
		Username:      ar.UserEmail,
		Password:      ar.Password,
		Name:          ar.FullName,
		GivenName:     "",
		FamilyName:    "",
		Email:         ar.UserEmail,
		EmailVerified: false,
		Website:       "",
		Address:       "",
	}
	_ ,err:= m.identityServerUc.UpdateUser(&updateUser)
	if err != nil{
		return err
	}
	layoutFormat := "2006-01-02 15:04:05"
	verificationSendDate, errDate := time.Parse(layoutFormat,ar.VerificationSendDate)
	if errDate != nil{
		return errDate
	}
	dob , errDateDob := time.Parse(layoutFormat,ar.Dob)
	if errDateDob != nil{
		return errDateDob
	}

	userModel := models.User{}
	userModel.Id = ar.Id
	userModel.ModifiedBy = &ar.UserEmail
	userModel.UserEmail = ar.UserEmail
	userModel.FullName = ar.FullName
	userModel.PhoneNumber = ar.PhoneNumber
	userModel.VerificationSendDate = verificationSendDate
	userModel.VerificationCode =ar.VerificationCode
	userModel.ProfilePictUrl = ar.ProfilePictUrl
	userModel.Address = ar.Address
	userModel.Dob = dob
	userModel.Gender = ar.Gender
	userModel.IdType = ar.IdType
	userModel.IdNumber = ar.IdNumber
	userModel.ReferralCode = ar.ReferralCode
	userModel.Points = ar.Points
	return m.userRepo.Update(ctx, &userModel)
}

func (m userUsecase) Create(c context.Context, ar *models.NewCommandUser,user string) error {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	existeduser, _ := m.userRepo.GetByUserEmail(ctx, ar.UserEmail)
	if existeduser != nil {
		return models.ErrConflict
	}
	registerUser := models.RegisterAndUpdateUser{
		Id:            "",
		Username:      ar.UserEmail,
		Password:      ar.Password,
		Name:          ar.FullName,
		GivenName:     "",
		FamilyName:    "",
		Email:         ar.UserEmail,
		EmailVerified: false,
		Website:       "",
		Address:       "",
	}
	isUser ,errorIs:= m.identityServerUc.CreateUser(&registerUser)
	ar.Id = isUser.Id
	layoutFormat := "2006-01-02 15:04:05"
	verificationSendDate, errDate := time.Parse(layoutFormat,ar.VerificationSendDate)
	if errDate != nil{
		return errDate
	}
	dob , errDateDob := time.Parse(layoutFormat,ar.Dob)
	if errDateDob != nil{
		return errDateDob
	}

	if errorIs != nil{
		return errorIs
	}
	userModel := models.User{}
	userModel.Id = isUser.Id
	userModel.CreatedBy = ar.UserEmail
	userModel.UserEmail = ar.UserEmail
	userModel.FullName = ar.FullName
	userModel.PhoneNumber = ar.PhoneNumber
	userModel.VerificationSendDate = verificationSendDate
	userModel.VerificationCode =ar.VerificationCode
	userModel.ProfilePictUrl = ar.ProfilePictUrl
	userModel.Address = ar.Address
	userModel.Dob = dob
	userModel.Gender = ar.Gender
	userModel.IdType = ar.IdType
	userModel.IdNumber = ar.IdNumber
	userModel.ReferralCode = ar.ReferralCode
	userModel.Points = ar.Points
	err := m.userRepo.Insert(ctx, &userModel)
	if err != nil {
		return err
	}

	return nil
}

func (m userUsecase) GetCreditByID(ctx context.Context, id string) (*models.UserPoint, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	point, err := m.userRepo.GetCreditByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.UserPoint{Points:point}, nil
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */

