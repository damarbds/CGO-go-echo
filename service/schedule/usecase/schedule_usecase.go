package usecase

import (
	"github.com/auth/merchant"
	"github.com/service/schedule"
	"github.com/service/time_options"
	"github.com/service/transportation"
	"golang.org/x/net/context"
	"time"
)

type scheduleUsecase struct {
	transportationRepo transportation.Repository
	merchantUsecase    merchant.Usecase
	scheduleRepo       schedule.Repository
	timeOptionsRepo    time_options.Repository
	contextTimeout     time.Duration
}



func NewScheduleUsecase(tr transportation.Repository, mr merchant.Usecase, s schedule.Repository, tmo time_options.Repository, timeout time.Duration) schedule.Usecase {
	return &scheduleUsecase{
		transportationRepo: tr,
		merchantUsecase:    mr,
		scheduleRepo:       s,
		timeOptionsRepo:    tmo,
		contextTimeout:     timeout,
	}
}
func (s scheduleUsecase) GetScheduleByMerchantId(c context.Context, merchantId string) () {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	getTransportationByMerchantId ,err := s.transportationRepo.GetByMerchantId(ctx,merchantId)
	if err != nil {
		return
	}

	getScheduleByTransId ,err := s.scheduleRepo.GetScheduleByTransId(ctx,getTransportationByMerchantId.Id)
}