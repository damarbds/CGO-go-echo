package schedule

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	GetScheduleByMerchantId(ctx context.Context,merchantId  string,date string)(*models.ScheduleDtoObj,error)
}
