package notif

import (
	"context"

	"github.com/models"
)

type Repository interface {
	GetCountByMerchantID(ctx context.Context,merchantId string)(int,error)
	GetByMerchantID(ctx context.Context, merchantId string,limit,offset int) ([]*models.Notification, error)
	Insert(ctx context.Context ,notification models.Notification)error
}
