package notif

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	GetByMerchantID(ctx context.Context, token string,page, limit, offset int) (*models.NotifWithPagination, error)
}