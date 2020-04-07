package notif

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByMerchantID(ctx context.Context, merchantId string) ([]*models.Notification, error)
}
