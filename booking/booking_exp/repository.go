package booking_exp

import (
	"time"

	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	GetCountByBookingDateExp(ctx context.Context,bookingDate string,expId string)(int,error)
	Insert(ctx context.Context, booking *models.BookingExp) (*models.BookingExp, error)
	GetEmailByID(ctx context.Context, bookingId string) (string, error)
	GetDetailBookingID(ctx context.Context, bookingId, bookingCode string) (*models.BookingExpJoin, error)
	UpdateStatus(ctx context.Context, bookingId string, expiredDatePayment time.Time) error
	GetBookingExpByUserID(ctx context.Context, bookingIds []*string) ([]*models.BookingExpJoin, error)
	GetBookingTransByUserID(ctx context.Context, bookingIds []*string) ([]*models.BookingExpJoin, error)
	GetBookingCountByUserID(ctx context.Context,status string, userId string)(int,error)
	GetBookingIdByUserID(ctx context.Context,status string, userId string,limit,offset int)([]*string,error)
	QueryHistoryPer30DaysExpByUserId(ctx context.Context, bookingIds []*string) ([]*models.BookingExpHistory, error)
	QueryHistoryPerMonthExpByUserId(ctx context.Context, bookingIds []*string) ([]*models.BookingExpHistory, error)
	QueryHistoryPer30DaysTransByUserId(ctx context.Context, bookingIds []*string) ([]*models.BookingExpJoin, error)
	QueryHistoryPerMonthTransByUserId(ctx context.Context, bookingIds []*string) ([] *models.BookingExpJoin, error)
	QueryCountHistoryByUserId(ctx context.Context, userId string, yearMonth string) (int, error)
	QuerySelectIdHistoryByUserId(ctx context.Context, userId string, yearMonth string,limit ,offset int) ([]*string, error)
	GetGrowthByMerchantID(ctx context.Context, merchantId string) ([]*models.BookingGrowth, error)
	CountThisMonth(ctx context.Context) (int, error)
	UpdatePaymentUrl(ctx context.Context, bookingId, paymentUrl string) error
	GetByID(ctx context.Context, bookingId string) (*models.BookingTransactionExp, error)
	CheckBookingCode(ctx context.Context, bookingCode string) bool
	GetDetailTransportBookingID(ctx context.Context, bookingId, bookingCode string) ([]*models.BookingExpJoin, error)
}
