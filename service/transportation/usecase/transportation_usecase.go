package usecase

import (
	"encoding/json"
	"time"

	"github.com/auth/merchant"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/service/schedule"
	"github.com/service/time_options"
	"github.com/service/transportation"
	"golang.org/x/net/context"
)

type transportationUsecase struct {
	transportationRepo transportation.Repository
	merchantUsecase    merchant.Usecase
	scheduleRepo       schedule.Repository
	timeOptionsRepo    time_options.Repository
	contextTimeout     time.Duration
}

// NewPromoUsecase will create new an articleUsecase object representation of article.Usecase interface
func NewTransportationUsecase(tr transportation.Repository, mr merchant.Usecase, s schedule.Repository, tmo time_options.Repository, timeout time.Duration) transportation.Usecase {
	return &transportationUsecase{
		transportationRepo: tr,
		merchantUsecase:    mr,
		scheduleRepo:       s,
		timeOptionsRepo:    tmo,
		contextTimeout:     timeout,
	}
}

func (t transportationUsecase) TimeOptions(ctx context.Context) ([]*models.TimeOptionDto, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	list, err := t.timeOptionsRepo.TimeOptions(ctx)
	if err != nil {
		return nil, err
	}

	timeOptions := make([]*models.TimeOptionDto, len(list))
	for i, item := range list {
		timeOptions[i] = &models.TimeOptionDto{
			Id:        item.Id,
			StartTime: item.StartTime,
			EndTime:   item.EndTime,
		}
	}

	return timeOptions, nil
}

func (t transportationUsecase) CreateTransportation(c context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	currentUserMerchant, err := t.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}
	var transImages string
	var harborsSourceId string
	var harborsDestId string
	if len(newCommandTransportation.TransImages) != 0 {
		transImagesConvert, _ := json.Marshal(newCommandTransportation.TransImages)
		transImages = string(transImagesConvert)
	}
	boatDetails, _ := json.Marshal(newCommandTransportation.BoatDetails)
	harborsDestId = newCommandTransportation.DepartureRoute.HarborsIdFrom
	harborsSourceId = newCommandTransportation.DepartureRoute.HarborsIdTo

	transportation := models.Transportation{
		Id:              guuid.New().String(),
		CreatedBy:       currentUserMerchant.MerchantEmail,
		CreatedDate:     time.Time{},
		ModifiedBy:      nil,
		ModifiedDate:    nil,
		DeletedBy:       nil,
		DeletedDate:     nil,
		IsDeleted:       0,
		IsActive:        0,
		TransName:       newCommandTransportation.TransName,
		HarborsSourceId: harborsSourceId,
		HarborsDestId:   harborsDestId,
		MerchantId:      currentUserMerchant.Id,
		TransCapacity:   newCommandTransportation.TransCapacity,
		TransTitle:      newCommandTransportation.TransTitle,
		TransStatus:     newCommandTransportation.Status,
		TransImages:     transImages,
		ReturnTransId:   nil,
		BoatDetails:     string(boatDetails),
		Transcoverphoto: newCommandTransportation.Transcoverphoto,
		Class:           newCommandTransportation.Class,
	}
	if newCommandTransportation.ReturnRoute != nil {
		harborsDestReturnId := newCommandTransportation.ReturnRoute.HarborsIdFrom
		harborsSourceReturnId := newCommandTransportation.ReturnRoute.HarborsIdTo
		transportationReturn := models.Transportation{
			Id:              guuid.New().String(),
			CreatedBy:       currentUserMerchant.MerchantEmail,
			CreatedDate:     time.Time{},
			ModifiedBy:      nil,
			ModifiedDate:    nil,
			DeletedBy:       nil,
			DeletedDate:     nil,
			IsDeleted:       0,
			IsActive:        0,
			TransName:       newCommandTransportation.TransName,
			HarborsSourceId: harborsSourceReturnId,
			HarborsDestId:   harborsDestReturnId,
			MerchantId:      currentUserMerchant.Id,
			TransCapacity:   newCommandTransportation.TransCapacity,
			TransTitle:      newCommandTransportation.TransTitle,
			TransStatus:     newCommandTransportation.Status,
			TransImages:     transImages,
			ReturnTransId:   nil,
			BoatDetails:     string(boatDetails),
			Transcoverphoto: newCommandTransportation.Transcoverphoto,
			Class:           newCommandTransportation.Class,
		}
		insertTransportationReturn, err := t.transportationRepo.Insert(ctx, transportationReturn)
		if err != nil {
			return nil, err
		}

		for _, year := range newCommandTransportation.ReturnRoute.Schedule {
			for _, month := range year.Month {
				for _, day := range month.DayPrice {
					for _, times := range newCommandTransportation.DepartureRoute.Time {
						var currency int
						if day.Currency == "USD" {
							currency = 1
						} else {
							currency = 0
						}
						priceObj := models.PriceObj{
							AdultPrice:    day.AdultPrice,
							ChildrenPrice: day.ChildrenPrice,
							Currency:      currency,
						}
						departureTimeOption, err := t.timeOptionsRepo.GetByTime(ctx, times.DepartureTime)
						if err != nil {
							return nil, err
						}
						arrivalTimeOption, err := t.timeOptionsRepo.GetByTime(ctx, times.ArrivalTime)
						if err != nil {
							return nil, err
						}
						price, _ := json.Marshal(priceObj)
						schedule := models.Schedule{
							Id:                    "",
							CreatedBy:             currentUserMerchant.MerchantEmail,
							CreatedDate:           time.Time{},
							ModifiedBy:            nil,
							ModifiedDate:          nil,
							DeletedBy:             nil,
							DeletedDate:           nil,
							IsDeleted:             0,
							IsActive:              0,
							TransId:               *insertTransportationReturn,
							DepartureTime:         times.DepartureTime,
							ArrivalTime:           times.ArrivalTime,
							Day:                   day.Day,
							Month:                 month.Month,
							Year:                  year.Year,
							DepartureDate: day.DepartureDate,
							Price:                 string(price),
							DepartureTimeoptionId: &departureTimeOption.Id,
							ArrivalTimeoptionId:   &arrivalTimeOption.Id,
						}
						_, err = t.scheduleRepo.Insert(ctx, schedule)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
		transportation.ReturnTransId = insertTransportationReturn
	}
	insertTransportation, err := t.transportationRepo.Insert(ctx, transportation)
	if err != nil {
		return nil, err
	}

	for _, year := range newCommandTransportation.DepartureRoute.Schedule {
		for _, month := range year.Month {
			for _, day := range month.DayPrice {
				for _, times := range newCommandTransportation.DepartureRoute.Time {
					var currency int
					if day.Currency == "USD" {
						currency = 1
					} else {
						currency = 0
					}
					priceObj := models.PriceObj{
						AdultPrice:    day.AdultPrice,
						ChildrenPrice: day.ChildrenPrice,
						Currency:      currency,
					}
					departureTimeOption, err := t.timeOptionsRepo.GetByTime(ctx, times.DepartureTime)
					if err != nil {
						return nil, err
					}
					arrivalTimeOption, err := t.timeOptionsRepo.GetByTime(ctx, times.ArrivalTime)
					if err != nil {
						return nil, err
					}
					price, _ := json.Marshal(priceObj)
					schedule := models.Schedule{
						Id:                    "",
						CreatedBy:             currentUserMerchant.MerchantEmail,
						CreatedDate:           time.Time{},
						ModifiedBy:            nil,
						ModifiedDate:          nil,
						DeletedBy:             nil,
						DeletedDate:           nil,
						IsDeleted:             0,
						IsActive:              0,
						TransId:               *insertTransportation,
						DepartureTime:         times.DepartureTime,
						ArrivalTime:           times.ArrivalTime,
						Day:                   day.Day,
						Month:                 month.Month,
						Year:                  year.Year,
						DepartureDate: day.DepartureDate,
						Price:                 string(price),
						DepartureTimeoptionId: &departureTimeOption.Id,
						ArrivalTimeoptionId:   &arrivalTimeOption.Id,
					}
					_, err = t.scheduleRepo.Insert(ctx, schedule)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	var status string
	if newCommandTransportation.Status == 0 {
		status = "Draft"
	} else if newCommandTransportation.Status == 3 {
		status = "Publish"
	}
	response := models.ResponseCreateExperience{
		Id:      *insertTransportation,
		Message: "Success " + status,
	}

	return &response, nil

}
func (t transportationUsecase) UpdateTransportation(c context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()

	currentUserMerchant, err := t.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}
	var transImages string
	var harborsSourceId string
	var harborsDestId string
	if len(newCommandTransportation.TransImages) != 0 {
		transImagesConvert, _ := json.Marshal(newCommandTransportation.TransImages)
		transImages = string(transImagesConvert)
	}
	boatDetails, _ := json.Marshal(newCommandTransportation.BoatDetails)
	harborsDestId = newCommandTransportation.DepartureRoute.HarborsIdFrom
	harborsSourceId = newCommandTransportation.DepartureRoute.HarborsIdTo

	transportation := models.Transportation{
		Id:              newCommandTransportation.Id,
		CreatedBy:       "",
		CreatedDate:     time.Time{},
		ModifiedBy:      &currentUserMerchant.MerchantEmail,
		ModifiedDate:    &time.Time{},
		DeletedBy:       nil,
		DeletedDate:     nil,
		IsDeleted:       0,
		IsActive:        0,
		TransName:       newCommandTransportation.TransName,
		HarborsSourceId: harborsSourceId,
		HarborsDestId:   harborsDestId,
		MerchantId:      currentUserMerchant.Id,
		TransCapacity:   newCommandTransportation.TransCapacity,
		TransTitle:      newCommandTransportation.TransTitle,
		TransStatus:     newCommandTransportation.Status,
		TransImages:     transImages,
		ReturnTransId:   nil,
		BoatDetails:     string(boatDetails),
		Transcoverphoto: newCommandTransportation.Transcoverphoto,
		Class:           newCommandTransportation.Class,
	}
	if newCommandTransportation.ReturnRoute != nil {
		harborsDestReturnId := newCommandTransportation.ReturnRoute.HarborsIdFrom
		harborsSourceReturnId := newCommandTransportation.ReturnRoute.HarborsIdTo
		transportationReturn := models.Transportation{
			Id:              newCommandTransportation.ReturnRoute.Id,
			CreatedBy:       "",
			CreatedDate:     time.Time{},
			ModifiedBy:      &currentUserMerchant.MerchantEmail,
			ModifiedDate:    &time.Time{},
			DeletedBy:       nil,
			DeletedDate:     nil,
			IsDeleted:       0,
			IsActive:        0,
			TransName:       newCommandTransportation.TransName,
			HarborsSourceId: harborsSourceReturnId,
			HarborsDestId:   harborsDestReturnId,
			MerchantId:      currentUserMerchant.Id,
			TransCapacity:   newCommandTransportation.TransCapacity,
			TransTitle:      newCommandTransportation.TransTitle,
			TransStatus:     newCommandTransportation.Status,
			TransImages:     transImages,
			ReturnTransId:   nil,
			BoatDetails:     string(boatDetails),
			Transcoverphoto: newCommandTransportation.Transcoverphoto,
			Class:           newCommandTransportation.Class,
		}
		insertTransportationReturn, err := t.transportationRepo.Update(ctx, transportationReturn)
		if err != nil {
			return nil, err
		}
		errorDelete := t.scheduleRepo.DeleteByTransId(ctx, insertTransportationReturn)
		if errorDelete != nil {
			return nil, errorDelete
		}
		for _, year := range newCommandTransportation.ReturnRoute.Schedule {
			for _, month := range year.Month {
				for _, day := range month.DayPrice {
					for _, times := range newCommandTransportation.DepartureRoute.Time {
						var currency int
						if day.Currency == "USD" {
							currency = 1
						} else {
							currency = 0
						}
						priceObj := models.PriceObj{
							AdultPrice:    day.AdultPrice,
							ChildrenPrice: day.ChildrenPrice,
							Currency:      currency,
						}
						departureTimeOption, err := t.timeOptionsRepo.GetByTime(ctx, times.DepartureTime)
						if err != nil {
							return nil, err
						}
						arrivalTimeOption, err := t.timeOptionsRepo.GetByTime(ctx, times.ArrivalTime)
						if err != nil {
							return nil, err
						}
						price, _ := json.Marshal(priceObj)
						schedule := models.Schedule{
							Id:                    "",
							CreatedBy:             currentUserMerchant.MerchantEmail,
							CreatedDate:           time.Time{},
							ModifiedBy:            nil,
							ModifiedDate:          nil,
							DeletedBy:             nil,
							DeletedDate:           nil,
							IsDeleted:             0,
							IsActive:              0,
							TransId:               *insertTransportationReturn,
							DepartureTime:         times.DepartureTime,
							ArrivalTime:           times.ArrivalTime,
							Day:                   day.Day,
							Month:                 month.Month,
							Year:                  year.Year,
							DepartureDate: day.DepartureDate,
							Price:                 string(price),
							DepartureTimeoptionId: &departureTimeOption.Id,
							ArrivalTimeoptionId:   &arrivalTimeOption.Id,
						}
						_, err = t.scheduleRepo.Insert(ctx, schedule)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
		transportation.ReturnTransId = insertTransportationReturn
	}
	insertTransportation, err := t.transportationRepo.Update(ctx, transportation)
	if err != nil {
		return nil, err
	}

	errorDelete := t.scheduleRepo.DeleteByTransId(ctx, insertTransportation)
	if errorDelete != nil {
		return nil, errorDelete
	}

	for _, year := range newCommandTransportation.DepartureRoute.Schedule {
		for _, month := range year.Month {
			for _, day := range month.DayPrice {
				for _, times := range newCommandTransportation.DepartureRoute.Time {
					var currency int
					if day.Currency == "USD" {
						currency = 1
					} else {
						currency = 0
					}
					priceObj := models.PriceObj{
						AdultPrice:    day.AdultPrice,
						ChildrenPrice: day.ChildrenPrice,
						Currency:      currency,
					}
					departureTimeOption, err := t.timeOptionsRepo.GetByTime(ctx, times.DepartureTime)
					if err != nil {
						return nil, err
					}
					arrivalTimeOption, err := t.timeOptionsRepo.GetByTime(ctx, times.ArrivalTime)
					if err != nil {
						return nil, err
					}
					price, _ := json.Marshal(priceObj)
					schedule := models.Schedule{
						Id:                    "",
						CreatedBy:             currentUserMerchant.MerchantEmail,
						CreatedDate:           time.Time{},
						ModifiedBy:            nil,
						ModifiedDate:          nil,
						DeletedBy:             nil,
						DeletedDate:           nil,
						IsDeleted:             0,
						IsActive:              0,
						TransId:               *insertTransportation,
						DepartureTime:         times.DepartureTime,
						ArrivalTime:           times.ArrivalTime,
						Day:                   day.Day,
						Month:                 month.Month,
						Year:                  year.Year,
						DepartureDate: day.DepartureDate,
						Price:                 string(price),
						DepartureTimeoptionId: &departureTimeOption.Id,
						ArrivalTimeoptionId:   &arrivalTimeOption.Id,
					}
					_, err = t.scheduleRepo.Insert(ctx, schedule)
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}
	var status string
	if newCommandTransportation.Status == 0 {
		status = "Draft"
	} else if newCommandTransportation.Status == 3 {
		status = "Publish"
	}
	response := models.ResponseCreateExperience{
		Id:      *insertTransportation,
		Message: "Success " + status,
	}

	return &response, nil

}

func (t transportationUsecase) PublishTransportation(c context.Context, newCommandTransportation models.NewCommandTransportation, token string) (*models.ResponseCreateExperience, error) {
	ctx, cancel := context.WithTimeout(c, t.contextTimeout)
	defer cancel()
	var response *models.ResponseCreateExperience
	if newCommandTransportation.Id == "" {
		create, err := t.CreateTransportation(ctx, newCommandTransportation, token)
		if err != nil {
			return nil, err
		}
		response = create
	} else {
		update, err := t.UpdateTransportation(ctx, newCommandTransportation, token)
		if err != nil {
			return nil, err
		}
		response = update
	}
	return response, nil
}
