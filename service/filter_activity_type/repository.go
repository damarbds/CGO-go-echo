package filter_activity_type

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByExpId(context context.Context, expId string) ([]*models.FilterActivityType, error)
	Insert(ctx context.Context,filterActivityType *models.FilterActivityType) error
}
