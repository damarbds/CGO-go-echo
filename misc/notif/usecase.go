package notif

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	FCMPushNotification(ctx context.Context, a models.FCMPushNotif)(*models.ResponseDelete,error)
	GetByMerchantID(ctx context.Context, token string,page, limit, offset int) (*models.NotifWithPagination, error)
}