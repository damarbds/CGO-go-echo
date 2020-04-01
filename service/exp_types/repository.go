package exp_types

import (
	"context"
	model "github.com/models"
)

type Repository interface {
	GetExpTypes(ctx context.Context) ([]*model.ExpTypeObject, error)
}
