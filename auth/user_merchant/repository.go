package user_merchant

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetUserByMerchantId(ctx context.Context,merchantId string)([]*models.UserMerchant,error)
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.UserMerchant, nextCursor string, err error)
	GetByID(ctx context.Context, id string) (*models.UserMerchant, error)
	GetByUserEmail(ctx context.Context, userEmail string) (*models.UserMerchant, error)
	Update(ctx context.Context, ar *models.UserMerchant) error
	Insert(ctx context.Context, a *models.UserMerchant) error
	Delete(ctx context.Context, id string, deleted_by string) error
	Count(ctx context.Context) (int, error)
	List(ctx context.Context, limit, offset int,search string) ([]*models.UserMerchant, error)
}
