package booking_exp

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Insert(ctx context.Context, booking *models.BookingExp) (*models.BookingExp, error)
	GetEmailByID(ctx context.Context, bookingId string) (string, error)
	GetDetailBookingID(ctx context.Context, bookingId string)(*models.BookingExpJoin,error)
	UpdateStatus(ctx context.Context, bookingId string) error
	GetByUserID(ctx context.Context, transactionStatus, bookingStatus int, userId string) ([]*models.BookingExpJoin, error)
}
