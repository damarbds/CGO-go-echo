package user_merchant

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Delete(ctx context.Context,userId string,token string)(*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandUserMerchant, isAdmin bool ,token string) error
	Create(ctx context.Context, ar *models.NewCommandUserMerchant, token string) (*models.NewCommandUserMerchant,error)
	List(ctx context.Context, page, limit, offset int,search string,token string) (*models.UserMerchantWithPagination, error)
	GetUserDetailById(ctx context.Context,id string,token string)(*models.UserMerchantDto, error)
}
