package exp_types

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetExpTypes(ctx context.Context) ([]*models.ExpTypeObject, error)
	GetByName(ctx context.Context,name string)(*models.ExpTypeObject,error)
}
