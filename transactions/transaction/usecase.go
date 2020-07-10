package transaction

import (
	"context"

	"github.com/models"
)

type Usecase interface {
	GetDetailTransactionSchedule(ctx context.Context,date string,transId string,expId string,token string,status string)(*models.TransactionScheduleDto,error)
	GetTransactionByDate(ctx context.Context,date string,isTransportation bool,isExperience bool,token string)([]*models.TransactionByDateDto,error)
	CountSuccess(ctx context.Context) (*models.Count, error)
	List(ctx context.Context, startDate, endDate, search, status string, page, limit, offset *int,token string,isAdmin bool,isTransportation bool,isExperience bool,isSchedule bool,tripType,paymentType,activityType string,confirmType string,class string,departureTimeStart string,departureTimeEnd string,arrivalTimeStart string,arrivalTimeEnd string,transactionId string) (*models.TransactionWithPagination, error)
	CountThisMonth(ctx context.Context) (*models.TotalTransaction, error)
}
