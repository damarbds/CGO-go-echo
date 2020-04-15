package schedule

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetScheduleByTransId(ctx context.Context,transId []*string)([]*models.ScheduleDtos,error)
	GetCountSchedule(ctx context.Context,transId []*string,date string)(int,error)
	Insert(ctx context.Context, a models.Schedule) (*string, error)
	DeleteByTransId(ctx context.Context, transId *string) error
}
