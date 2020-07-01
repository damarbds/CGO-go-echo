package currency

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context,rate *models.ExChangeRate)error
	GetByDate(ctx context.Context,from ,to string)(*models.ExChangeRate,error)
}
