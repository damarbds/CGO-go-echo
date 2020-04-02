package booking_exp

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Insert(ctx context.Context, booking *models.NewBookingExpCommand) (*models.NewBookingExpCommand,error,error)
}
