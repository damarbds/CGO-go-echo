package usecase

import (
	"encoding/json"
	"github.com/auth/merchant"
	"github.com/models"
	"github.com/service/exp_availability"
	"github.com/service/experience"
	"github.com/service/schedule"
	"github.com/service/time_options"
	"github.com/service/transportation"
	"golang.org/x/net/context"
	"time"
)

type scheduleUsecase struct {
	expAvailability exp_availability.Repository
	experieceRepo 	experience.Repository
	transportationRepo transportation.Repository
	merchantUsecase    merchant.Usecase
	scheduleRepo       schedule.Repository
	timeOptionsRepo    time_options.Repository
	contextTimeout     time.Duration
}



func NewScheduleUsecase(tr transportation.Repository, mr merchant.Usecase, s schedule.Repository, tmo time_options.Repository, exp experience.Repository,expA exp_availability.Repository,timeout time.Duration) schedule.Usecase {
	return &scheduleUsecase{
		transportationRepo: tr,
		merchantUsecase:    mr,
		scheduleRepo:       s,
		timeOptionsRepo:    tmo,
		experieceRepo:exp,
		expAvailability:expA,
		contextTimeout:     timeout,
	}
}
func (s scheduleUsecase) GetScheduleByMerchantId(c context.Context, merchantId string) (*models.ScheduleDtoObj,error) {
	ctx, cancel := context.WithTimeout(c, s.contextTimeout)
	defer cancel()

	getTransportationByMerchantId ,err := s.transportationRepo.SelectIdGetByMerchantId(ctx,merchantId)
	if err != nil {
		return nil,err
	}
	var scheduleDtos models.ScheduleDto
	scheduleDtos.MerchantId = merchantId
		getScheduleByTransId ,err := s.scheduleRepo.GetScheduleByTransIds(ctx,getTransportationByMerchantId)
		if err != nil {
			return nil,err
		}
		for _,element := range getScheduleByTransId {
			schedule := models.ScheduleDateObj{
				Type:  "Transportation",
				Date:  element.DepartureDate.Format("2006-01-02"),
				Count: 0,
			}
			count ,err := s.scheduleRepo.GetCountSchedule(ctx,getTransportationByMerchantId,element.DepartureDate.Format("2006-01-02"))
			if err != nil {
				return nil,err
			}
			schedule.Count = count
			scheduleDtos.ScheduleDate = append(scheduleDtos.ScheduleDate,schedule)
		}

	getExperienceByMerchantId , err := s.experieceRepo.SelectIdGetByMerchantId(ctx,merchantId)

		getAvalibitryByExpId , err := s.expAvailability.GetByExpIds(ctx,getExperienceByMerchantId)
		if err != nil {
			return nil,err
		}
		for _,element := range getAvalibitryByExpId{
				var dates []string
				if errUnmarshal := json.Unmarshal([]byte(element.ExpAvailabilityDate), &dates); errUnmarshal != nil {
					return nil, models.ErrInternalServerError
				}
				for _,date := range dates {
					schedule := models.ScheduleDateObj{
						Type:  "Experience",
						Date:  date,
						Count: 0,
					}
					count ,err := s.expAvailability.GetCountDate(ctx,date,getExperienceByMerchantId)
					if err != nil {
						return nil,err
					}
					schedule.Count = count
					var scheduler []models.ScheduleDateObj
					for _,check := range scheduleDtos.ScheduleDate {
						if schedule.Type == check.Type && schedule.Date == check.Date {
							scheduler = append(scheduler, schedule)
						}
					}
					if len(scheduler) == 0 {
						scheduleDtos.ScheduleDate = append(scheduleDtos.ScheduleDate,schedule)
					}


			}
		}
	
		var result models.ScheduleDtoObj
		result.MerchantId = merchantId
	for _,element := range scheduleDtos.ScheduleDate {
		obj := models.ScheduleObjDate{
			Date:                element.Date,
			TransportationCount: 0,
			ExperienceCount:     0,
		}
		if element.Type == "Experience"{
			obj.ExperienceCount = element.Count
		}else {
			obj.TransportationCount = element.Count
		}
		var scheduler []models.ScheduleObjDate
		for index,check := range result.ScheduleDate {
			if obj.Date == check.Date {
				if element.Type == "Experience" && check.ExperienceCount != element.Count{
					result.ScheduleDate[index].ExperienceCount = element.Count
				}else if element.Type == "Transportation" && check.ExperienceCount != element.Count{
					result.ScheduleDate[index].TransportationCount = element.Count
				}
				scheduler = append(scheduler, obj)
			}
		}
		if len(scheduler) == 0 {
			result.ScheduleDate = append(result.ScheduleDate,obj)
		}

	}
	return &result,nil

}