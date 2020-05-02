package transaction

import (
	"context"

	"github.com/models"
)

type Repository interface {
	GetCountByTransId(ctx context.Context, transId string)(int,error)
	GetById(ctx context.Context,id string)(*models.TransactionWMerchant,error)
	CountSuccess(ctx context.Context) (int, error)
	Count(ctx context.Context, startDate, endDate, search, status string,merchantId string) (int, error)
	List(ctx context.Context, startDate, endDate, search, status string, limit, offset *int,merchantId string) ([]*models.TransactionOut, error)
	CountThisMonth(ctx context.Context) (*models.TotalTransaction, error)
	UpdateStatus(ctx context.Context, status int, transactionId, bookingId string) error
}
