package usecase

import (
	"github.com/auth/identityserver"
	"github.com/auth/user"
	"github.com/models"
	"golang.org/x/net/context"
	"math"
	"math/rand"
	"time"
)

type userUsecase struct {
	userRepo         user.Repository
	identityServerUc identityserver.Usecase
	contextTimeout   time.Duration
}



// NewuserUsecase will create new an userUsecase object representation of user.Usecase interface
func NewuserUsecase(a user.Repository, is identityserver.Usecase, timeout time.Duration) user.Usecase {
	return &userUsecase{
		userRepo:         a,
		identityServerUc: is,
		contextTimeout:   timeout,
	}
}
func (m userUsecase) RequestOTP(ctx context.Context, phoneNumber string) (*models.RequestOTP, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	requestOTP ,err := m.identityServerUc.RequestOTP(phoneNumber)
	if err != nil {
		return nil,err
	}

	getUserByPhoneNumber ,err := m.userRepo.GetByUserNumberOTP(ctx,phoneNumber,"")
	if getUserByPhoneNumber != nil {
		user := getUserByPhoneNumber
		user.VerificationCode = requestOTP.OTP

		err = m.userRepo.Update(ctx,user)
		if err != nil {
			return nil,err
		}
	}

	return requestOTP,nil
}
func (m userUsecase) List(ctx context.Context, page, limit, offset int) (*models.UserWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	list, err := m.userRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	users := make([]*models.UserInfoDto, len(list))
	for i, item := range list {
		users[i] = &models.UserInfoDto{
			Id:             item.Id,
			CreatedDate:    item.CreatedDate,
			UpdatedDate:    item.ModifiedDate,
			IsActive:       item.IsActive,
			UserEmail:      item.UserEmail,
			FullName:       item.FullName,
			PhoneNumber:    item.PhoneNumber,
			ProfilePictUrl: item.ProfilePictUrl,
			Address:        item.Address,
			Dob:            item.Dob,
			Gender:         item.Gender,
			IdType:         item.IdType,
			IdNumber:       item.IdNumber,
			ReferralCode:   item.ReferralCode,
			Points:         item.Points,
		}
	}
	totalRecords, _ := m.userRepo.Count(ctx)
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

	response := &models.UserWithPagination{
		Data: users,
		Meta: meta,
	}

	return response, nil
}

func (m userUsecase) Login(ctx context.Context, ar *models.Login) (*models.GetToken, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()
	var requestToken *models.GetToken

	if ar.Scope == "phone_number"{
		checkPhoneNumberExists ,_ := m.userRepo.GetByUserNumberOTP(ctx, ar.Email,"")
		if checkPhoneNumberExists == nil {
			return nil, models.ErrNotYetRegister
		}
		existeduser, _ := m.userRepo.GetByUserNumberOTP(ctx, ar.Email,ar.Password)
		if existeduser == nil {
			return nil, models.ErrInvalidOTP
		}
		getToken, err := m.identityServerUc.GetToken(ar.Email, ar.Password,ar.Scope)
		if err != nil {
			return nil, err
		}

		requestToken = getToken
	}else {
		getToken, err := m.identityServerUc.GetToken(ar.Email, ar.Password,ar.Scope)
		if err != nil {
			return nil, err
		}
		existeduser, _ := m.userRepo.GetByUserEmail(ctx, ar.Email)
		if existeduser == nil {
			return nil, models.ErrUsernamePassword
		}
		requestToken = getToken
	}
	return requestToken, nil
}

func (m userUsecase) ValidateTokenUser(ctx context.Context, token string) (*models.UserInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	existeduser, _ := m.userRepo.GetByUserEmail(ctx, getInfoToIs.Email)
	if existeduser == nil {
		return nil, models.ErrUnAuthorize
	}

	userInfo := &models.UserInfoDto{
		Id:             existeduser.Id,
		UserEmail:      existeduser.UserEmail,
		FullName:       existeduser.FullName,
		PhoneNumber:    existeduser.PhoneNumber,
		ProfilePictUrl: existeduser.ProfilePictUrl,
	}

	return userInfo, nil
}

func (m userUsecase) VerifiedEmail(ctx context.Context, token string, codeOTP string) (*models.UserInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	existeduser, _ := m.userRepo.GetByUserEmail(ctx, getInfoToIs.Email)
	if existeduser == nil {
		return nil, models.ErrUnAuthorize
	}
	verifiedEmail := models.VerifiedEmail{
		Email:   existeduser.UserEmail,
		CodeOTP: codeOTP,
	}
	_, error := m.identityServerUc.VerifiedEmail(&verifiedEmail)
	if error != nil {
		return nil, error
	}
	userInfo := &models.UserInfoDto{
		Id:             existeduser.Id,
		UserEmail:      existeduser.UserEmail,
		FullName:       existeduser.FullName,
		PhoneNumber:    existeduser.PhoneNumber,
		ProfilePictUrl: existeduser.ProfilePictUrl,
		ReferralCode:   existeduser.ReferralCode,
	}

	return userInfo, nil
}
func (m userUsecase) GetUserInfo(ctx context.Context, token string) (*models.UserInfoDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	getInfoToIs, err := m.identityServerUc.GetUserInfo(token)
	if err != nil {
		return nil, err
	}
	existeduser, _ := m.userRepo.GetByUserEmail(ctx, getInfoToIs.Email)
	if existeduser == nil {
		return nil, models.ErrNotFound
	}
	userInfo := models.UserInfoDto{
		Id:             existeduser.Id,
		UserEmail:      existeduser.UserEmail,
		FullName:       existeduser.FullName,
		PhoneNumber:    existeduser.PhoneNumber,
		ProfilePictUrl: existeduser.ProfilePictUrl,
		ReferralCode:   existeduser.ReferralCode,
	}

	return &userInfo, nil
}

func (m userUsecase) Update(c context.Context, ar *models.NewCommandUser, user string) error {
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
		EmailVerified: true,
		Website:       "",
		Address:       "",
	}
	_, err := m.identityServerUc.UpdateUser(&updateUser)
	if err != nil {
		return err
	}
	layoutFormat := "2006-01-02 15:04:05"
	//verificationSendDate, errDate := time.Parse(layoutFormat,ar.VerificationSendDate)
	//if errDate != nil{
	//	return errDate
	//}
	dob, errDateDob := time.Parse(layoutFormat, ar.Dob)
	if errDateDob != nil {
		return errDateDob
	}
	existeduser, _ := m.userRepo.GetByUserEmail(ctx, ar.UserEmail)
	userModel := models.User{}
	userModel.Id = existeduser.Id
	userModel.ModifiedBy = &ar.UserEmail
	userModel.UserEmail = ar.UserEmail
	userModel.FullName = ar.FullName
	userModel.PhoneNumber = ar.PhoneNumber
	userModel.VerificationSendDate = existeduser.VerificationSendDate
	userModel.VerificationCode = existeduser.VerificationCode
	userModel.ProfilePictUrl = ar.ProfilePictUrl
	userModel.Address = ar.Address
	userModel.Dob = dob
	userModel.Gender = ar.Gender
	userModel.IdType = ar.IdType
	userModel.IdNumber = ar.IdNumber
	userModel.ReferralCode = existeduser.ReferralCode
	userModel.Points = ar.Points
	return m.userRepo.Update(ctx, &userModel)
}
func generateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes, err := generateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}
func (m userUsecase) Create(c context.Context, ar *models.NewCommandUser, user string) (*models.NewCommandUser,error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	//existeduser, _ := m.userRepo.GetByUserEmail(ctx, ar.UserEmail)
	//if existeduser != nil {
	//	return models.ErrConflict
	//}
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
		OTP:           "",
		UserType:      1,
	}
	isUser, errorIs := m.identityServerUc.CreateUser(&registerUser)
	message := "Please keep it a secret, and use this OTP: " + isUser.OTP + " code to verify your email"
	email := models.SendingEmail{
		Subject: "Verified Email",
		Message: message,
		From:    "helmy@cgo.co.id",
		To:      isUser.Email,
	}
	_, errorSending := m.identityServerUc.SendingEmail(&email)
	if errorSending != nil {
		return nil,models.ErrInternalServerError
	}
	ar.Id = isUser.Id
	var dob time.Time
	if ar.Dob != "" {

		layoutFormat := "2006-01-02 15:04:05"

		dobs, errDateDob := time.Parse(layoutFormat, ar.Dob)

		if errDateDob != nil {
			return nil,errDateDob
		}
		dob = dobs
	}

	if errorIs != nil {
		return nil,errorIs
	}
	referralCode, er := generateRandomString(9)
	if er != nil {
		return nil,er
	}
	userModel := models.User{}
	userModel.Id = isUser.Id
	userModel.CreatedBy = ar.UserEmail
	userModel.UserEmail = ar.UserEmail
	userModel.FullName = ar.FullName
	userModel.PhoneNumber = ar.PhoneNumber
	userModel.VerificationSendDate = time.Now()
	userModel.VerificationCode = isUser.OTP
	userModel.ProfilePictUrl = ar.ProfilePictUrl
	userModel.Address = ar.Address
	userModel.Dob = dob
	userModel.Gender = ar.Gender
	userModel.IdType = ar.IdType
	userModel.IdNumber = ar.IdNumber
	userModel.ReferralCode = referralCode
	userModel.Points = ar.Points
	err := m.userRepo.Insert(ctx, &userModel)
	if err != nil {
		return nil, err
	}
	requestToken ,err := m.identityServerUc.GetToken(ar.UserEmail,ar.Password,"")

	ar.Token = &requestToken.AccessToken
	return ar,nil
}

func (m userUsecase) GetCreditByID(ctx context.Context, id string) (*models.UserPoint, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	point, err := m.userRepo.GetCreditByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &models.UserPoint{Points: point}, nil
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
