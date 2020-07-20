package filter_activity_type

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context,filterActivityType *models.FilterActivityType) error
	DeleteByExpId(ctx context.Context,expId string) error
}
