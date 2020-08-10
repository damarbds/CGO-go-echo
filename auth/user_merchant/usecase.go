package user_merchant

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	UpdateFCMToken(ctx context.Context,a models.TokenFCM,token string)(*models.ResponseDelete,error)
	GetUserByMerchantId(ctx context.Context,merchantId string,token string)([]*models.UserMerchantWithRole,error)
	AssignRoles(ctx context.Context,token string,isAdmin bool,aRoles *models.NewCommandAssignRoleUserMerchant)(*models.ResponseAssignRoles,error)
	GetRoles(ctx context.Context,token string,isAdmin bool)([]*models.RolesUserMerchant,error)
	Delete(ctx context.Context,userId string,token string)(*models.ResponseDelete, error)
	Update(ctx context.Context, ar *models.NewCommandUserMerchant, isAdmin bool ,token string) error
	Create(ctx context.Context, ar *models.NewCommandUserMerchant, token string) (*models.NewCommandUserMerchant,error)
	List(ctx context.Context, page, limit, offset int,search string,token string) (*models.UserMerchantWithPagination, error)
	GetUserDetailById(ctx context.Context,id string,token string)(*models.UserMerchantDto, error)
}
