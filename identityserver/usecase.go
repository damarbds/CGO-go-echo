package identityserver

import (
	"github.com/models"
)

type Usecase interface {
	UpdateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser,error)
	CreateUser(ar *models.RegisterAndUpdateUser) (*models.RegisterAndUpdateUser,error)
	GetUserInfo(token string) (*models.GetUserInfo, error)
	GetToken(username string, password string) (*models.GetToken, error)
}