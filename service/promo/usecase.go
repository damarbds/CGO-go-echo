package promo

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Fetch(ctx context.Context, page *int , size *int) ([]*models.PromoDto, error)
	GetByCode(ctx context.Context, code string) ([]*models.PromoDto, error)
}
