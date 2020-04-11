package transaction

import (
	"context"
	"github.com/models"
)

type Repository interface {
	CountSuccess(ctx context.Context) (int, error)
	Count(ctx context.Context, status string) (int, error)
	List(ctx context.Context, status string, limit, offset int) ([]*models.TransactionOut, error)
}
