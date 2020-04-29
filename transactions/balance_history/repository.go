package balance_history

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Count(ctx context.Context, merchantId string, status string) (int, error)
	GetAll(ctx context.Context, merchantId string, status string, limit, offset *int, month, year string) ([]*models.BalanceHistory, error)
	Insert(ctx context.Context, balanceH models.BalanceHistory) (*string, error)
	Update(ctx context.Context, balanceH models.BalanceHistory)(*string,error)
	GetById(ctx context.Context,id string)(*models.BalanceHistory,error)
}
