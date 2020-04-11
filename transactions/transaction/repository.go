package transaction

import (
	"context"

	"github.com/models"
)

type Repository interface {
	CountSuccess(ctx context.Context) (int, error)
	Count(ctx context.Context, startDate, endDate, search, status string) (int, error)
	List(ctx context.Context, startDate, endDate, search, status string, limit, offset int) ([]*models.TransactionOut, error)
	CountThisMonth(ctx context.Context) (*models.TotalTransaction, error)
}
