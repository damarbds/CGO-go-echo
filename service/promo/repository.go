package promo

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Fetch(ctx context.Context, page *int, size *int) ([]*models.Promo, error)
	GetByCode(ctx context.Context, code string) ([]*models.Promo, error)
}
