package schedule

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	InsertSchedule(ctx context.Context,command *models.NewCommandSchedule)(*models.NewCommandSchedule,error)
	GetScheduleByMerchantId(ctx context.Context,merchantId  string,date string)(*models.ScheduleDtoObj,error)
}
