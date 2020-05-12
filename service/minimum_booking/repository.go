package minimum_booking

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetAll(ctx context.Context,limit,offset int)([]*models.MinimumBooking,error)
	GetCount(ctx context.Context)(int,error)
}
