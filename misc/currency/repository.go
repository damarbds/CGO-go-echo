package currency

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByDate(ctx context.Context,from ,to string)(*models.ExChangeRate,error)
}
