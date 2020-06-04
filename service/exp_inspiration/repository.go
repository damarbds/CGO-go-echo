package exp_inspiration

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetExpInspirations(ctx context.Context) ([]*models.ExpInspirationObject, error)
}
