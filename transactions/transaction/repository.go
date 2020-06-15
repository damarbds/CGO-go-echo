package transaction

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetTransactionDownPaymentByDate(ctx context.Context)([]*models.TransactionWithBooking,error)
	GetCountByExpId(ctx context.Context, date string, expId string) (*string, error)
	GetCountByTransId(ctx context.Context, transId string) (int, error)
	GetById(ctx context.Context, id string) (*models.TransactionWMerchant, error)
	CountSuccess(ctx context.Context) (int, error)
	Count(ctx context.Context, startDate, endDate, search, status string, merchantId string) (int, error)
	List(ctx context.Context, startDate, endDate, search, status string, limit, offset *int, merchantId string,isTransportation bool,isExperience bool,isSchedule bool,tripType,paymentType,activityType string,confirmType string) ([]*models.TransactionOut, error)
	CountThisMonth(ctx context.Context) (*models.TotalTransaction, error)
	UpdateAfterPayment(ctx context.Context, status int, vaNumber string, transactionId, bookingId string) error
}
