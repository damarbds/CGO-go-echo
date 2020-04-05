package booking_exp

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetDetailBookingID(ctx context.Context, bookingId string)(*models.BookingExpDetailDto,error)
	Insert(ctx context.Context, booking *models.NewBookingExpCommand,token string) (*models.NewBookingExpCommand,error,error)
	GetByUserID(ctx context.Context, transactionStatus, bookingStatus int, token string) ([]*models.MyBooking, error)
}
