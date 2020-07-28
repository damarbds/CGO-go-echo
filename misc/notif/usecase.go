package notif

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	GetByMerchantID(ctx context.Context, token string) ([]*models.NotifDto, error)
}