package notif

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	UpdateStatusNotif(ctx context.Context ,notif models.NotificationRead,token string)(*models.ResponseDelete,error)
	DeleteNotificationByIds(ctx context.Context,token string, ids string)(*models.ResponseDelete,error)
	FCMPushNotification(ctx context.Context, a models.FCMPushNotif)(*models.ResponseDelete,error)
	GetByMerchantID(ctx context.Context, token string,page, limit, offset int) (*models.NotifWithPagination, error)
}