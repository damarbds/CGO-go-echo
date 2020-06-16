package booking_exp

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	DownloadTicketTransportation(ctx context.Context,orderId string)(*string, error)
	DownloadTicketExperience(ctx context.Context,orderId string)(*string, error)
	RemainingPaymentNotification(ctx context.Context)error
	GetByGuestCount(ctx context.Context, expId string, date string, guest int) (bool, error)
	GetDetailBookingID(ctx context.Context, bookingId, bookingCode string) (*models.BookingExpDetailDto, error)
	GetDetailTransportBookingID(ctx context.Context, bookingId, bookingCode string,transId *string) (*models.BookingExpDetailDto, error)
	GetHistoryBookingByUserId(ctx context.Context, token string, monthType string, page, limit, offset int) (*models.BookingHistoryDtoWithPagination, error)
	Insert(ctx context.Context, booking *models.NewBookingExpCommand, transReturnId, scheduleReturnId, token string) ([]*models.NewBookingExpCommand, error, error)
	GetByUserID(ctx context.Context, status string, token string, page, limit, offset int) (*models.MyBookingWithPagination, error)
	GetGrowthByMerchantID(ctx context.Context, token string) ([]*models.BookingGrowthDto, error)
	CountThisMonth(ctx context.Context) (*models.Count, error)
	SendCharge(ctx context.Context, bookingCode, paymentType string) (map[string]interface{}, error)
	Verify(ctx context.Context, orderId, bookingCode string) (map[string]interface{}, error)
	XenPayment(ctx context.Context, amount float64, tokenId, authId, orderId, paymentType string) (map[string]interface{}, error)
}
