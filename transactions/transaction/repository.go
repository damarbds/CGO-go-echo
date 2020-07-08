package transaction

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetTransactionByExpIdORTransId(ctx context.Context,date string,expId string,transId string,merchantId string)([]*models.TransactionOut, error)
	GetTransactionByDate(ctx context.Context,date string ,isExperience bool,isTransportation bool,merchantId string)([]*models.TransactionByDate,error)
	GetCountTransactionByPromoId(ctx context.Context,promoId string,userId string)(int,error)
	GetTransactionDownPaymentByDate(ctx context.Context)([]*models.TransactionWithBooking,error)
	GetIdTransactionExpired(ctx context.Context)([]*string ,error)
	GetCountByExpId(ctx context.Context, date string, expId string) (*string, error)
	GetCountByTransId(ctx context.Context, transId string) (int, error)
	GetById(ctx context.Context, id string) (*models.TransactionWMerchant, error)
	GetByBookingDate(ctx context.Context, bookingDate string,transId string,expId string) ([]*models.TransactionWMerchant, error)
	CountSuccess(ctx context.Context) (int, error)
	Count(ctx context.Context, startDate, endDate, search, status string, merchantId string,isTransportation bool,isExperience bool,isSchedule bool,tripType,paymentType,activityType string,confirmType string,class string,departureTimeStart string,departureTimeEnd string,arrivalTimeStart string,arrivalTimeEnd string) (int, error)
	List(ctx context.Context, startDate, endDate, search, status string, limit, offset *int, merchantId string,isTransportation bool,isExperience bool,isSchedule bool,tripType,paymentType,activityType string,confirmType string,class string,departureTimeStart string,departureTimeEnd string,arrivalTimeStart string,arrivalTimeEnd string) ([]*models.TransactionOut, error)
	CountThisMonth(ctx context.Context) (*models.TotalTransaction, error)
	UpdateAfterPayment(ctx context.Context, status int, vaNumber string, transactionId, bookingId string) error
}
