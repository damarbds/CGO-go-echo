package schedule

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetTimeByTransId(ctx context.Context,transId string)([]*models.ScheduleTime,error)
	GetYearByTransId(ctx context.Context,transId string)([]*models.ScheduleYear,error)
	GetMonthByTransId(ctx context.Context,transId string,year int)([]*models.ScheduleMonth,error)
	GetDayByTransId(ctx context.Context,transId string,year int,month string)([]*models.ScheduleDay,error)
	GetScheduleByTransIds(ctx context.Context,transId []*string,year int,month string)([]*models.ScheduleDtos,error)
	GetCountSchedule(ctx context.Context,merchantId string,date string)(int,error)
	Insert(ctx context.Context, a models.Schedule) (*string, error)
	DeleteByTransId(ctx context.Context, transId *string) error
}
