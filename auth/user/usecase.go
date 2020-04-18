package user

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Update(ctx context.Context, ar *models.NewCommandUser, user string) error
	Create(ctx context.Context, ar *models.NewCommandUser, isAdmin bool,token string) (*models.NewCommandUser,error)
	ValidateTokenUser(ctx context.Context, token string) (*models.UserInfoDto, error)
	RequestOTP(ctx context.Context,phoneNumber string)(*models.RequestOTP,error)
	VerifiedEmail(ctx context.Context, token string, codeOTP string) (*models.UserInfoDto, error)
	Login(ctx context.Context, ar *models.Login) (*models.GetToken, error)
	GetUserInfo(ctx context.Context, token string) (*models.UserInfoDto, error)
	GetCreditByID(ctx context.Context, id string) (*models.UserPoint, error)
	List(ctx context.Context, page, limit, offset int,search string) (*models.UserWithPagination, error)
	GetUserDetailById(ctx context.Context,id string,token string)(*models.UserDto, error)
}
