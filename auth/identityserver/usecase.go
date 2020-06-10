package identityserver

import (
	"github.com/models"
)

type Usecase interface {
	ForgotPassword(email string,token string)(*models.ResponseDelete,error)
	CallBackGoogle(code string)(*models.GetInfoUserGoogle,error)
	DeleteUser(userId string)error
	GetListOfRole(roleType int)([]*models.RolesPermissionIs,error)
	UpdateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser, error)
	CreateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser, error)
	SendingEmail(r *models.SendingEmail) (*models.SendingEmail, error)
	VerifiedEmail(r *models.VerifiedEmail) (*models.VerifiedEmail, error)
	GetUserInfo(token string) (*models.GetUserInfo, error)
	GetToken(username string, password string,scope string) (*models.GetToken, error)
	RefreshToken(refreshToken string) (*models.RefreshToken, error)
	UploadFileToBlob(image string, folder string) (string, error)
	UploadFilePDFToBlob(bit []byte,folder string) (string, error)
	RequestOTP(phoneNumber string)(*models.RequestOTP,error)
	RequestOTPTmp(phoneNumber string,email string)(*models.RequestOTP,error)
	SendingSMS(sms *models.SendingSMS)(*models.SendingSMS,error)
	GetDetailUserById(id string,token string,isDetail string)(*models.GetUserDetail,error)
}
