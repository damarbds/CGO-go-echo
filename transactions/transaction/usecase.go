package transaction

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	CountSuccess(ctx context.Context) (*models.Count, error)
	List(ctx context.Context, startDate, endDate, search, status string, page, limit, offset int) (*models.TransactionWithPagination, error)
}
