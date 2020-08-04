package merchant

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	SendingEmailMerchant(ctx context.Context,ar *models.NewCommandMerchantRegistrationEmail)(*models.ResponseDelete,error)
	SendingEmailContactUs(ctx context.Context,ar *models.NewCommandContactUs)(*models.ResponseDelete,error)
	AutoLoginByCMSAdmin(ctx context.Context,merchantId string,token string)(*models.GetToken,error)
	Update(ctx context.Context, ar *models.NewCommandMerchant, isAdmin bool,token string) error
	Create(ctx context.Context, ar *models.NewCommandMerchant, token string) error
	Login(ctx context.Context, ar *models.Login) (*models.GetToken, error)
	ValidateTokenMerchant(ctx context.Context, token string) (*models.MerchantInfoDto, error)
	GetMerchantInfo(ctx context.Context, token string) (*models.MerchantInfoDto, error)
	Count(ctx context.Context) (*models.Count, error)
	List(ctx context.Context, page, limit, offset int, token string,search string) (*models.MerchantWithPagination, error)
	Delete(ctx context.Context, id string, token string) (*models.ResponseDelete, error)
	ServiceCount(ctx context.Context, token string) (*models.ServiceCount, error)
	GetDetailMerchantById(ctx context.Context, id string,token string)(*models.MerchantDto,error)
}
