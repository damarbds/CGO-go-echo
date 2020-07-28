package filter_activity_type

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetJoinExpType(ctx context.Context ,expId string)([]*models.FilterActivityTypeJoin,error)
	Insert(ctx context.Context,filterActivityType *models.FilterActivityType) error
	DeleteByExpId(ctx context.Context,expId string) error
}
