package balance_history

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context, merchantId, status string, page int, limit, offset *int, month, year string) (*models.BalanceHistoryDtoWithPagination, error)
	Create(ctx context.Context, bHistory models.NewBalanceHistoryCommand, token string) (*models.NewBalanceHistoryCommand, error)
}
