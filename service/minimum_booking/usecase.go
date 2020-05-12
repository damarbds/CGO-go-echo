package minimum_booking

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetAll(ctx context.Context,page,limit,size int)(*models.MinimumBookingDtoWithPagination,error)
}
