package user

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Update(ctx context.Context, ar *models.NewCommandUser, user string) error
	Create(ctx context.Context, ar *models.NewCommandUser, user string) error
	ValidateTokenUser(ctx context.Context, token string) (*string,error)
	Login(ctx context.Context, ar *models.Login) (*models.GetToken,error)
	GetUserInfo(ctx context.Context, token string) (*models.UserInfoDto,error)
	GetCreditByID(ctx context.Context, id string) (*models.UserPoint, error)
}
