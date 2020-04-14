package schedule

import "context"

type Usecase interface {
	GetScheduleByMerchantId(ctx context.Context,merchantId  string)()
}
