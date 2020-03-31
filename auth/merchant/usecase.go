package merchant

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Update(ctx context.Context, ar *models.NewCommandMerchant, user string) error
	Create(ctx context.Context, ar *models.NewCommandMerchant, user string) error
	Login(ctx context.Context, ar *models.Login) (*models.GetToken,error)
	ValidateTokenMerchant(ctx context.Context, token string) (*string,error)
	GetMerchantInfo(ctx context.Context, token string) (*models.MerchantInfoDto,error)
}