package schedule

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetBookingBySchedule(ctx context.Context,transId string,departureDate string,arrivalTime string ,departureTime string)([]*models.Schedule,error)
	GetTimeByTransId(ctx context.Context,transId string)([]*models.ScheduleTime,error)
	GetYearByTransId(ctx context.Context,transId string,arrivalTime string ,departureTime string)([]*models.ScheduleYear,error)
	GetMonthByTransId(ctx context.Context,transId string,year int,arrivalTime string ,departureTime string)([]*models.ScheduleMonth,error)
	GetDayByTransId(ctx context.Context,transId string,year int,month string,arrivalTime string ,departureTime string)([]*models.ScheduleDay,error)
	GetScheduleByTransIds(ctx context.Context,transId []*string,year int,month string)([]*models.ScheduleDtos,error)
	GetCountSchedule(ctx context.Context,merchantId string,date string)(int,error)
	Insert(ctx context.Context, a models.Schedule) (*string, error)
	Update(ctx context.Context, a models.Schedule) error
	DeleteByTransId(ctx context.Context, transId *string) error
}
