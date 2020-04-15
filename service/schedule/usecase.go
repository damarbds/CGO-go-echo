package schedule

import (
	"context"
	"github.com/models"
)

type Usecase interface {
	GetScheduleByMerchantId(ctx context.Context,merchantId  string)(*models.ScheduleDtoObj,error)
}
