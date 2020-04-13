package identityserver

import (
	"github.com/models"
)

type Usecase interface {
	UpdateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser, error)
	CreateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser, error)
	SendingEmail(r *models.SendingEmail) (*models.SendingEmail, error)
	VerifiedEmail(r *models.VerifiedEmail) (*models.VerifiedEmail, error)
	GetUserInfo(token string) (*models.GetUserInfo, error)
	GetToken(username string, password string) (*models.GetToken, error)
	UploadFileToBlob(image string, folder string) (string, error)
}
