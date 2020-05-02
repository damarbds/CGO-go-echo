package booking_exp

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetByGuestCount(ctx context.Context,expId string,date string,guest int)(bool,error)
	GetDetailBookingID(ctx context.Context, bookingId, bookingCode string) (*models.BookingExpDetailDto, error)
	GetDetailTransportBookingID(ctx context.Context, bookingId, bookingCode string) (*models.BookingExpDetailDto, error)
	GetHistoryBookingByUserId(ctx context.Context, token string, monthType string) ([]*models.BookingHistoryDto, error)
	Insert(ctx context.Context, booking *models.NewBookingExpCommand, transReturnId, scheduleReturnId, token string) ([]*models.NewBookingExpCommand, error, error)
	GetByUserID(ctx context.Context, transactionStatus, bookingStatus int, token string) ([]*models.MyBooking, error)
	GetGrowthByMerchantID(ctx context.Context, token string) ([]*models.BookingGrowthDto, error)
	CountThisMonth(ctx context.Context) (*models.Count, error)
	SendCharge(ctx context.Context, bookingCode, paymentType string) (map[string]interface{}, error)
	Verify(ctx context.Context, orderId, bookingCode string) (map[string]interface{}, error)
}
