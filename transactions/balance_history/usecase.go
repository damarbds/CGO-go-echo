package balance_history

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	List(ctx context.Context, merchantId, status string, page int, limit, offset *int, month, year string,token string,isAdmin bool) (*models.BalanceHistoryDtoWithPagination, error)
	Create(ctx context.Context, bHistory models.NewBalanceHistoryCommand, token string) (*models.NewBalanceHistoryCommand, error)
	UpdateAmount(ctx context.Context,amount models.NewBalanceHistoryAmountCommand,token string)(*models.ResponseDelete,error)
	ConfirmWithdraw(ctx context.Context,command models.NewBalanceHistoryConfirmCommand,token string)(*models.ResponseDelete,error)
}
