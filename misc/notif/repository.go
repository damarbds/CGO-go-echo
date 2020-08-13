package notif

import (
	"context"
	"time"

	"github.com/models"
)

type Repository interface {
	UpdateStatusNotif(ctx context.Context ,notif models.NotificationRead,modifiedBy string,modifyDate time.Time)error
	DeleteNotificationByIds(ctx context.Context,merchantId string, ids string,deletedby string,deletedDate time.Time)error
	GetCountByMerchantID(ctx context.Context,merchantId string,notifType string)(int,error)
	GetByMerchantID(ctx context.Context, merchantId string,limit,offset int,notifType string) ([]*models.Notification, error)
	Insert(ctx context.Context ,notification models.Notification)error
}
