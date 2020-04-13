package transportation

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context, transportation models.Transportation) (*string, error)
	Update(ctx context.Context, transportation models.Transportation) (*string, error)
	FilterSearch(ctx context.Context, query string, limit, offset int) ([]*models.TransSearch, error)
	CountFilterSearch(ctx context.Context, query string) (int, error)
	GetTransCount(ctx context.Context, merchantId string) (int, error)
}
