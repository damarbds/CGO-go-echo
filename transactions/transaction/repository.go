package transaction

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetCountTransactionByExpIdORTransId(ctx context.Context,date string,expId string,transId string,merchantId string,status string)(int,[]string, error)
	GetTransactionByExpIdORTransId(ctx context.Context,date string,expId string,transId string,merchantId string,status string)([]*models.TransactionOut, error)
	GetTransactionByDate(ctx context.Context,date string ,isExperience bool,isTransportation bool,merchantId string)([]*models.TransactionByDate,error)
	GetCountTransactionByPromoId(ctx context.Context,promoId string,userId string)(int,error)
	GetTransactionDownPaymentByDate(ctx context.Context)([]*models.TransactionWithBooking,error)
	GetTransactionExpired(ctx context.Context) ([]*models.TransactionWithBookingExpired, error)
	GetIdTransactionByStatus(ctx context.Context,transactionStatus int)([]*string ,error)
	GetCountByExpId(ctx context.Context, date string, expId string,isTransaction bool) ([]*string, error)
	GetCountByTransId(ctx context.Context, transId string,isTransaction bool,date string) ([]*string, error)
	GetById(ctx context.Context, id string) (*models.TransactionWMerchant, error)
	GetByBookingDate(ctx context.Context, bookingDate string,transId string,expId string) ([]*models.TransactionWMerchant, error)
	CountSuccess(ctx context.Context) (int, error)
	Count(ctx context.Context, startDate, endDate, search, status string, merchantId string,isTransportation bool,isExperience bool,isSchedule bool,tripType,paymentType,activityType string,confirmType string,class string,departureTimeStart string,departureTimeEnd string,arrivalTimeStart string,arrivalTimeEnd string,transactionId string) (int, error)
	List(ctx context.Context, startDate, endDate, search, status string, limit, offset *int, merchantId string,isTransportation bool,isExperience bool,isSchedule bool,tripType,paymentType,activityType string,confirmType string,class string,departureTimeStart string,departureTimeEnd string,arrivalTimeStart string,arrivalTimeEnd string,transactionId string) ([]*models.TransactionOut, error)
	CountThisMonth(ctx context.Context) (*models.TotalTransaction, error)
	UpdateAfterPayment(ctx context.Context, status int, vaNumber string, transactionId, bookingId string) error
}
