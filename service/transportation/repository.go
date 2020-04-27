package transportation

import (
	"context"
	"github.com/models"
)

type Repository interface {
	UpdateStatus(ctx context.Context, status int,id string,user string)error
	Insert(ctx context.Context, transportation models.Transportation) (*string, error)
	Update(ctx context.Context, transportation models.Transportation) (*string, error)
	FilterSearch(ctx context.Context, query string, limit, offset int) ([]*models.TransSearch, error)
	CountFilterSearch(ctx context.Context, query string) (int, error)
	GetTransCount(ctx context.Context, merchantId string) (int, error)
	SelectIdGetByMerchantId(ctx context.Context,merchantId string)([]*string,error)
}
