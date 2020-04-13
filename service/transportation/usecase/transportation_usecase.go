package usecase

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
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

func NewTransportationUsecase(tr transportation.Repository, mr merchant.Usecase, s schedule.Repository, tmo time_options.Repository, timeout time.Duration) transportation.Usecase {
	return &transportationUsecase{
		transportationRepo: tr,
		merchantUsecase:    mr,
		scheduleRepo:       s,
		timeOptionsRepo:    tmo,
		contextTimeout:     timeout,
	}
}

func (t transportationUsecase) FilterSearchTrans(
	ctx context.Context,
	search,
	qStatus,
	sortBy,
	harborSourceId,
	harborDestId,
	depDate,
	class string,
	isReturn bool,
	depTimeOptions,
	arrTimeOptions,
	guest,
	page,
	limit,
	offset int,
) (*models.FilterSearchTransWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	query := `
	SELECT
		s.id as schedule_id,
		s.departure_date,
		s.departure_time,
		s.arrival_time,
		s.price,
		s.trans_id,
		t.trans_name,
		t.trans_images,
		t.trans_status,
		h.id as harbor_source_id,
		h.harbors_name as harbor_source_name,
		hdest.id as harbor_dest_id,
		hdest.harbors_name as harbor_dest_name
	FROM
		schedules s
		JOIN transportations t ON s.trans_id = t.id
		JOIN harbors h ON t.harbors_source_id = h.id
		JOIN harbors hdest ON t.harbors_dest_id = hdest.id
	WHERE
		s.is_deleted = 0
		AND s.is_active = 1`

	queryCount := `
	SELECT
		COUNT(*)
	FROM
		schedules s
		JOIN transportations t ON s.trans_id = t.id
		JOIN harbors h ON t.harbors_source_id = h.id
		JOIN harbors hdest ON t.harbors_dest_id = hdest.id
	WHERE
		s.is_deleted = 0
		AND s.is_active = 1`

	if qStatus != "" {
		var status int
		if qStatus == "preview" {
			status = 0
		} else if qStatus == "draft" {
			status = 1
		} else if qStatus == "published" {
			status = 2
		} else if qStatus == "unpublished" {
			status = 3
		} else if qStatus == "archived" {
			status = 4
		}

		query = query + ` AND t.trans_status =` + strconv.Itoa(status)
		queryCount = queryCount + ` AND t.trans_status =` + strconv.Itoa(status)
	}
	if search != "" {
		keyword := `'%` + search + `%'`
		query = query + ` AND (LOWER(t.trans_name) LIKE LOWER(` + keyword + `) OR LOWER(h.harbors_name) LIKE LOWER(` + keyword + `) OR LOWER(hdest.harbors_name) LIKE LOWER(` + keyword + `))`
		queryCount = queryCount + ` AND (LOWER(t.trans_name) LIKE LOWER(` + keyword + `) OR LOWER(h.harbors_name) LIKE LOWER(` + keyword + `) OR LOWER(hdest.harbors_name) LIKE LOWER(` + keyword + `))`
	}
	if harborSourceId != "" {
		query = query + ` AND t.harbors_source_id =` + harborSourceId
		queryCount = queryCount + ` AND t.harbors_source_id =` + harborSourceId
		if isReturn {
			query = query + ` AND t.harbors_source_id =` + harborDestId
			queryCount = queryCount + ` AND t.harbors_source_id =` + harborDestId
		}
	}
	if harborDestId != "" {
		query = query + ` AND t.harbors_dest_id =` + harborDestId
		queryCount = queryCount + ` AND t.harbors_dest_id =` + harborDestId
		if isReturn {
			query = query + ` AND t.harbors_dest_id =` + harborSourceId
			queryCount = queryCount + ` AND t.harbors_dest_id =` + harborSourceId
		}
	}
	if guest != 0 {
		query = query + ` AND t.harbors_dest_id =` + strconv.Itoa(guest)
		queryCount = queryCount + ` AND t.harbors_dest_id =` + strconv.Itoa(guest)
	}
	if depDate != "" {
		query = query + ` AND s.departure_date =` + depDate
		queryCount = queryCount + ` AND s.departure_date =` + depDate
	}
	if class != "" {
		query = query + ` AND t.class =` + class
		queryCount = queryCount + ` AND t.class =` + class
	}
	if depTimeOptions != 0 {
		query = query + ` AND s.departure_timeoption_id =` + strconv.Itoa(depTimeOptions)
		queryCount = queryCount + ` AND s.departure_timeoption_id =` + strconv.Itoa(depTimeOptions)
	}
	if arrTimeOptions != 0 {
		query = query + ` AND s.departure_timeoption_id =` + strconv.Itoa(arrTimeOptions)
		queryCount = queryCount + ` AND s.departure_timeoption_id =` + strconv.Itoa(arrTimeOptions)
	}

	transList, err := t.transportationRepo.FilterSearch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	trans := make([]*models.TransportationSearchObj, len(transList))
	for i, t := range transList {
		var transImages []models.CoverPhotosObj
		if errUnmarshal := json.Unmarshal([]byte(t.TransImages), &transImages); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		var transPrice models.TransPriceObj
		if errUnmarshal := json.Unmarshal([]byte(t.Price), &transPrice); errUnmarshal != nil {
			return nil, errUnmarshal
		}
		transPrice.PriceType = "pax"
		if transPrice.Currency == 1 {
			transPrice.CurrencyLabel = "USD"
		} else {
			transPrice.CurrencyLabel = "IDR"
		}

		departureTime, _ := time.Parse("15:04:05", t.DepartureTime)
		arrivalTime, _ := time.Parse("15:04:05", t.ArrivalTime)

		tripHour := arrivalTime.Hour() - departureTime.Hour()
		tripMinute := arrivalTime.Minute() - departureTime.Minute()
		tripDuration := strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`

		var transStatus string
		if t.TransStatus == 0 {
			transStatus = "Preview"
		} else if t.TransStatus == 1 {
			transStatus = "Draft"
		} else if t.TransStatus == 2 {
			transStatus = "Published"
		} else if t.TransStatus == 3 {
			transStatus = "Unpublished"
		} else if t.TransStatus == 4 {
			transStatus = "Archived"
		}

		trans[i] = &models.TransportationSearchObj{
			ScheduleId:            t.ScheduleId,
			DepartureDate:         t.DepartureDate,
			DepartureTime:         t.DepartureTime,
			ArrivalTime:           t.ArrivalTime,
			TripDuration:          tripDuration,
			TransportationId:      t.TransId,
			TransportationName:    t.TransName,
			TransportationImages:  transImages,
			TransportationStatus:  transStatus,
			HarborSourceId:        t.HarborSourceId,
			HarborSourceName:      t.HarborSourceName,
			HarborDestinationId:   t.HarborDestId,
			HarborDestinationName: t.HarborDestName,
			Price:                 transPrice,
		}
	}
	if sortBy != "" {
		if sortBy == "priceup" {
			sort.SliceStable(trans, func(i, j int) bool {
				return trans[i].Price.AdultPrice > trans[j].Price.AdultPrice
			})
			sort.SliceStable(trans, func(i, j int) bool {
				return trans[i].Price.ChildrenPrice > trans[j].Price.ChildrenPrice
			})
		} else if sortBy == "pricedown" {
			sort.SliceStable(trans, func(i, j int) bool {
				return trans[i].Price.AdultPrice < trans[j].Price.AdultPrice
			})
			sort.SliceStable(trans, func(i, j int) bool {
				return trans[i].Price.ChildrenPrice < trans[j].Price.ChildrenPrice
			})
		}
	}
	totalRecords, _ := t.transportationRepo.CountFilterSearch(ctx, queryCount)
	totalPage := int(math.Ceil(float64(totalRecords) / float64(limit)))
	prev := page
	next := page
	if page != 1 {
		prev = page - 1
	}

	if page != totalPage {
		next = page + 1
	}

	meta := &models.MetaPagination{
		Page:          page,
		Total:         totalPage,
		TotalRecords:  totalRecords,
		Prev:          prev,
		Next:          next,
		RecordPerPage: len(trans),
	}

	response := &models.FilterSearchTransWithPagination{
		Data: trans,
		Meta: meta,
	}

	return response, nil
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
							DepartureDate:         day.DepartureDate,
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
						DepartureDate:         day.DepartureDate,
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
							DepartureDate:         day.DepartureDate,
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
						DepartureDate:         day.DepartureDate,
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
