package usecase_test

import (
	"context"
	"errors"
	//_adminUsecaseMock "github.com/auth/admin/mocks"
	_expAvailabilityRepoMock "github.com/service/exp_availability/mocks"
	_timeOptionRepoMock "github.com/service/time_options/mocks"
	//_merchantUsecaseMock "github.com/auth/merchant/mocks"
	"github.com/models"
	"github.com/service/schedule/mocks"
	"github.com/service/schedule/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)
var(
	timeoutContext = time.Second*30
	mockScheduleRepo = new(mocks.Repository)
	mockExpAvaiblityrepo = new(_expAvailabilityRepoMock.Repository)
	mockTimeOptionRepo = new(_timeOptionRepoMock.Repository)
	merchantId = "sakldjlkasdjsalkdj"
	date = "2020-09"
	mockScheduleTime = []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	mockScheduleDay = []*models.ScheduleDay{
		&models.ScheduleDay{
			Year:          2020,
			Month:         "May",
			Day:           "Thursday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
		&models.ScheduleDay{
			Year:          2020,
			Month:         "June",
			Day:           "Monday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
	}
	a = models.Schedule{
		Id:                    "asdqweqwsad",
		CreatedBy:             "Test",
		CreatedDate:           time.Now(),
		ModifiedBy:            nil,
		ModifiedDate:          nil,
		DeletedBy:             nil,
		DeletedDate:           nil,
		IsDeleted:             0,
		IsActive:              1,
		TransId:               "asdqdsadqeqweqw",
		DepartureTime:         mockScheduleTime[0].DepartureTime,
		ArrivalTime:           mockScheduleTime[0].ArrivalTime,
		Day:                   mockScheduleDay[0].Day,
		Month:                 mockScheduleDay[0].Month,
		Year:                  mockScheduleDay[0].Year,
		DepartureDate:         "2020-09-17",
		Price:                 mockScheduleDay[0].Price,
		DepartureTimeoptionId: nil,
		ArrivalTimeoptionId:   nil,
	}
)
func TestGetScheduleByMerchantId(t *testing.T) {

	//mockMerchantUsecae := new(_merchantUsecaseMock.)
	t.Run("success", func(t *testing.T) {
		dateParse := date + "-" + "01"
		layoutFormat := "2006-01-02"
		dt, _ := time.Parse(layoutFormat, dateParse)
		//year := dt.Year()
		month := dt.Month().String()

		var dates []string

		start, _ := time.Parse(layoutFormat, dateParse)
		//if errDateDob != nil {
		//	return nil, errDateDob
		//}
		dates = append(dates, start.Format("2006-01-02"))

	datess:
		start = start.AddDate(0, 0, 1)
		if start.Month().String() != month {

		} else {
			dates = append(dates, start.Format("2006-01-02"))
			goto datess
		}
		var result models.ScheduleDtoObj
		result.MerchantId = merchantId
		for _,_element := range dates{
			mockScheduleRepo.On("GetCountSchedule", mock.Anything, merchantId,_element).Return(2, nil).Once()
			mockExpAvaiblityrepo.On("GetCountDate", mock.Anything, _element,merchantId).Return(2, nil).Once()
		}

		u := usecase.NewScheduleUsecase(nil, nil,mockScheduleRepo,mockTimeOptionRepo,nil,mockExpAvaiblityrepo, timeoutContext)

		a, err := u.GetScheduleByMerchantId(context.TODO(), merchantId,date)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockScheduleRepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		dateParse := date + "-" + "01"
		layoutFormat := "2006-01-02"
		dt, _ := time.Parse(layoutFormat, dateParse)
		//year := dt.Year()
		month := dt.Month().String()

		var dates []string

		start, _ := time.Parse(layoutFormat, dateParse)
		//if errDateDob != nil {
		//	return nil, errDateDob
		//}
		dates = append(dates, start.Format("2006-01-02"))

	datess:
		start = start.AddDate(0, 0, 1)
		if start.Month().String() != month {

		} else {
			dates = append(dates, start.Format("2006-01-02"))
			goto datess
		}
		var result models.ScheduleDtoObj
		result.MerchantId = merchantId
		for _,_element := range dates{
			mockScheduleRepo.On("GetCountSchedule", mock.Anything, merchantId,_element).Return(0, errors.New("unexpected")).Once()
			mockExpAvaiblityrepo.On("GetCountDate", mock.Anything, _element,merchantId).Return(0, nil).Once()
		}
		u := usecase.NewScheduleUsecase(nil, nil,mockScheduleRepo,mockTimeOptionRepo,nil,mockExpAvaiblityrepo, timeoutContext)

		a, err := u.GetScheduleByMerchantId(context.TODO(), merchantId,date)

		assert.Error(t, err)
		assert.Nil(t, a)

		//mockScheduleRepo.AssertExpectations(t)
	})

}

func TestInsertSchedule(t *testing.T) {
	mockScheduleRepo := new(mocks.Repository)
	timeOptionMock := models.TimesOption{
		Id:           1,
		CreatedBy:    "test1",
		CreatedDate:  time.Now(),
		ModifiedBy:   nil,
		ModifiedDate: nil,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     1,
		StartTime:    "13:00:00",
		EndTime:      "16:00:00",
	}
	t.Run("success", func(t *testing.T) {
		tempMockSchedule := models.NewCommandSchedule{
			TimeObj:   []models.TimeObj{
				models.TimeObj{
					DepartureTime: "12:00:00",
					ArrivalTime:   "13:00:00",
				},
			},
			CreatedBy: "Admin",
			TransId:   "qweqweqweqwe",
			Schedule:  []models.YearObj{
				models.YearObj{
					Year:  2020,
					Month: []models.MonthObj{
						models.MonthObj{
							Month:    "June",
							DayPrice: []models.DayPriceObj{
								models.DayPriceObj{
									DepartureDate: "2021-04-18",
									Day:           "Sunday",
									AdultPrice:    10,
									ChildrenPrice: 10,
									Currency:      "IDR",
								},
							},
						},
					},
				},
			},
		}
		//token := "eyJhbGciOiJSUzI1NiIsImtpZCI6ImZhZWE3Y2Q2YWFhYjM1YmIyYmE4MjE3ZTgyNWNkODE5IiwidHlwIjoiSldUIn0.eyJuYmYiOjE1OTM2MjU3MzYsImV4cCI6MTU5NDIzMDUzNiwiaXNzIjoiaHR0cDovL2lkZW50aXR5LXNlcnZlci1jZ28taW5kb25lc2lhLmF6dXJld2Vic2l0ZXMubmV0IiwiYXVkIjpbImh0dHA6Ly9pZGVudGl0eS1zZXJ2ZXItY2dvLWluZG9uZXNpYS5henVyZXdlYnNpdGVzLm5ldC9yZXNvdXJjZXMiLCJhcGkxIiwiYXBpMiJdLCJjbGllbnRfaWQiOiJyb2NsaWVudCIsInN1YiI6Ijk3MmZlMDlmLTkzZTktNDc5OC1iNjQyLTE0ZTBhYzc3YzZiZSIsImF1dGhfdGltZSI6MTU5MzYyNTczNiwiaWRwIjoibG9jYWwiLCJuYW1lIjoiYWRtaW5DR08iLCJlbWFpbCI6ImFkbWluMTIzNEBnbWFpbC5jb20iLCJzY29wZSI6WyJjdXN0b20ucHJvZmlsZSIsIm9wZW5pZCIsImFwaTEiLCJhcGkyLnJlYWRfb25seSIsIm9mZmxpbmVfYWNjZXNzIl0sImFtciI6WyJwd2QiXX0.XkNnCV-GwYRFTjoll7Y_FeTOJ6AlPyzHnJFFErgzsVM5EPTAVfetre0jXflHe8cTJ52iEWqAB3RKYi2ckHr-9-LER0Z5L3ir7kS7d7-Rmf268ob4vlhLxFNV6QFEvpoz1JRqjo6KzIKCuBWTZV22N_Ipb6R4_geLISILfSlWmxlZxEEzqMxPUdwWdY7GqByI0qNmx93-MVMyjwdcfQENGlP5xkdmuCiFzFGAjdgezy1GqJhZ4svOYNDh5R56pZf8A3kBA20n31MvQJqDn-BE4LLmygCZMCgZQdwDitJKH1AnpuU5smcnrSXZt4xbFGIv0up517TgIBEDWabbU-8U7Q"
		//tempMockSchedule.Id = 0
		mockTimeOptionRepo.On("GetByTime", mock.Anything, tempMockSchedule.TimeObj[0].DepartureTime).Return(&timeOptionMock,nil).Once()
		mockTimeOptionRepo.On("GetByTime", mock.Anything, tempMockSchedule.TimeObj[0].ArrivalTime).Return(&timeOptionMock,nil).Once()
		//for _, year := range tempMockSchedule.Schedule {
		//	for _, month := range year.Month {
		//		for _, _ = range month.DayPrice {
		//			for _, _ = range tempMockSchedule.TimeObj {
		//
		//			}
		//		}
		//	}
		//}
		scheduleId := "dasljdaskljd"
		mockScheduleRepo.On("Insert", mock.Anything, mock.AnythingOfType("models.Schedule")).Return(&scheduleId,nil).Once()
		u := usecase.NewScheduleUsecase(nil, nil,mockScheduleRepo,mockTimeOptionRepo,nil,mockExpAvaiblityrepo, timeoutContext)

		_,err := u.InsertSchedule(context.TODO(), &tempMockSchedule)

		assert.NoError(t, err)
		//assert.Equal(t, mockSchedule.ScheduleName, tempMockSchedule.ScheduleName)
		mockScheduleRepo.AssertExpectations(t)
	})


}



