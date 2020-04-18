package booking_exp

import (
	"time"

	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Insert(ctx context.Context, booking *models.BookingExp) (*models.BookingExp, error)
	GetEmailByID(ctx context.Context, bookingId string) (string, error)
	GetDetailBookingID(ctx context.Context, bookingId, bookingCode string) (*models.BookingExpJoin, error)
	UpdateStatus(ctx context.Context, bookingId string, expiredDatePayment time.Time) error
	GetByUserID(ctx context.Context, transactionStatus, bookingStatus int, userId string) ([]*models.BookingExpJoin, error)
	QueryHistoryPer30DaysByUserId(ctx context.Context, userId string) ([]*models.BookingExpHistory, error)
	QueryHistoryPerMonthByUserId(ctx context.Context, userId string, yearMonth string) ([]*models.BookingExpHistory, error)
	GetGrowthByMerchantID(ctx context.Context, merchantId string) ([]*models.BookingGrowth, error)
	CountThisMonth(ctx context.Context) (int, error)
	UpdatePaymentUrl(ctx context.Context, bookingId, paymentUrl string) error
}
