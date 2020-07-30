package usecase

import (
	"encoding/json"
	"github.com/misc/currency"
	"math"
	"sort"
	"strconv"
	"time"

	"github.com/service/exp_facilities"
	"github.com/service/facilities"
	"github.com/service/transportation/delivery/events"

	"github.com/auth/merchant"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/service/schedule"
	"github.com/service/time_options"
	"github.com/service/transportation"
	"github.com/transactions/transaction"
	"golang.org/x/net/context"
)

type transportationUsecase struct {
	transactionRepo    transaction.Repository
	transportationRepo transportation.Repository
	merchantUsecase    merchant.Usecase
	scheduleRepo       schedule.Repository
	timeOptionsRepo    time_options.Repository
	contextTimeout     time.Duration
	expFacilitiesRepo  exp_facilities.Repository
	facilitiesRepo     facilities.Repository
	currencyUsecase	currency.Usecase
}


func NewTransportationUsecase(currencyUsecase	currency.Usecase,expFacilitiesRepo exp_facilities.Repository, facilitiesRepo facilities.Repository, transr transaction.Repository, tr transportation.Repository, mr merchant.Usecase, s schedule.Repository, tmo time_options.Repository, timeout time.Duration) transportation.Usecase {
	return &transportationUsecase{
		currencyUsecase:currencyUsecase,
		expFacilitiesRepo:  expFacilitiesRepo,
		facilitiesRepo:     facilitiesRepo,
		transactionRepo:    transr,
		transportationRepo: tr,
		merchantUsecase:    mr,
		scheduleRepo:       s,
		timeOptionsRepo:    tmo,
		contextTimeout:     timeout,
	}
}

func (t transportationUsecase) GetPublishedTransCount(ctx context.Context, token string) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()
	currentMerchant, err := t.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, models.ErrUnAuthorize
	}

	count, err := t.transportationRepo.GetTransCount(ctx, currentMerchant.Id)
	if err != nil {
		return nil, err
	}
	return &models.Count{Count: count}, nil

}

func (t transportationUsecase) GetAllTransport(ctx context.Context) ([]*models.MasterDataTransport, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()


	count, err := t.transportationRepo.GetAllTransport(ctx)
	if err != nil {
		return nil, err
	}

	result := []*models.MasterDataTransport{}
	for _, element := range count {
		res := models.MasterDataTransport{
			Id: element.Id,
			TransportName: element.TransName,
		}
		result = append(result, &res)
	}
	return result, nil

}

func (t transportationUsecase) GetDetail(ctx context.Context, id string) (*models.TransportationDto, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	getDetailTrans, err := t.transportationRepo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	var transObj models.RouteObj
	var returnTransObj models.RouteObj
	if getDetailTrans.ReturnTransId != nil {
		getreturnTrans, err := t.transportationRepo.GetById(ctx, *getDetailTrans.ReturnTransId)
		schedules := make([]models.YearObj, 0)
		times := make([]models.TimeObj, 0)
		if err != nil {
			return nil, err
		}
		returnTrans := models.RouteObj{
			Id:            getreturnTrans.Id,
			HarborsIdFrom: getreturnTrans.HarborsSourceId,
			HarborsIdTo:   getreturnTrans.HarborsDestId,
			Time:          times,
			Schedule:      schedules,
		}
		getScheduleTimeReturn, err := t.scheduleRepo.GetTimeByTransId(ctx, getreturnTrans.Id)

		for _, element := range getScheduleTimeReturn {
			timeObj := models.TimeObj{
				DepartureTime: element.DepartureTime,
				ArrivalTime:   element.ArrivalTime,
			}
			returnTrans.Time = append(returnTrans.Time, timeObj)
		}

		getScheduleYearReturn, err := t.scheduleRepo.GetYearByTransId(ctx, getreturnTrans.Id,returnTrans.Time[0].ArrivalTime,returnTrans.Time[0].DepartureTime)
		if err != nil {
			return nil, err
		}
		for _, year := range getScheduleYearReturn {
			month := make([]models.MonthObj, 0)
			schedule := models.YearObj{
				Year:  year.Year,
				Month: month,
			}
			getScheduleMonthReturn, err := t.scheduleRepo.GetMonthByTransId(ctx, getreturnTrans.Id, schedule.Year,returnTrans.Time[0].ArrivalTime,returnTrans.Time[0].DepartureTime)
			if err != nil {
				return nil, err
			}
			for _, monthElement := range getScheduleMonthReturn {
				day := make([]models.DayPriceObj, 0)
				monthMap := models.MonthObj{
					Month:    monthElement.Month,
					DayPrice: day,
				}
				getScheduleDayReturn, err := t.scheduleRepo.GetDayByTransId(ctx, getreturnTrans.Id, schedule.Year, monthMap.Month,returnTrans.Time[0].ArrivalTime,returnTrans.Time[0].DepartureTime)
				if err != nil {
					return nil, err
				}
				for _, dayElement := range getScheduleDayReturn {
					var price models.PriceObj
					var currency string
					if dayElement.Price != "" {
						if errUnmarshal := json.Unmarshal([]byte(dayElement.Price), &price); errUnmarshal != nil {
							return nil, errUnmarshal
						}
					}
					if price.Currency == 1 {
						currency = "USD"
					} else {
						currency = "IDR"
					}
					dayMap := models.DayPriceObj{
						DepartureDate: dayElement.DepartureDate.Format("2006-01-02"),
						Day:           dayElement.Day,
						AdultPrice:    price.AdultPrice,
						ChildrenPrice: price.ChildrenPrice,
						Currency:      currency,
					}
					monthMap.DayPrice = append(monthMap.DayPrice, dayMap)
				}
				schedule.Month = append(schedule.Month, monthMap)
			}
			returnTrans.Schedule = append(returnTrans.Schedule, schedule)
		}
		returnTransObj = returnTrans
	}

	schedules := make([]models.YearObj, 0)
	times := make([]models.TimeObj, 0)
	if err != nil {
		return nil, err
	}
	trans := models.RouteObj{
		Id:            getDetailTrans.Id,
		HarborsIdFrom: getDetailTrans.HarborsSourceId,
		HarborsIdTo:   getDetailTrans.HarborsDestId,
		Time:          times,
		Schedule:      schedules,
	}
	getScheduleTime, err := t.scheduleRepo.GetTimeByTransId(ctx, getDetailTrans.Id)
	for _, element := range getScheduleTime {
		timeObj := models.TimeObj{
			DepartureTime: element.DepartureTime,
			ArrivalTime:   element.ArrivalTime,
		}
		trans.Time = append(trans.Time, timeObj)
	}

	getScheduleYear, err := t.scheduleRepo.GetYearByTransId(ctx, getDetailTrans.Id,trans.Time[0].ArrivalTime,trans.Time[0].DepartureTime)
	if err != nil {
		return nil, err
	}
	for _, year := range getScheduleYear {
		month := make([]models.MonthObj, 0)
		schedule := models.YearObj{
			Year:  year.Year,
			Month: month,
		}
		getScheduleMonth, err := t.scheduleRepo.GetMonthByTransId(ctx, getDetailTrans.Id, schedule.Year,trans.Time[0].ArrivalTime,trans.Time[0].DepartureTime)
		if err != nil {
			return nil, err
		}
		for _, monthElement := range getScheduleMonth {
			day := make([]models.DayPriceObj, 0)
			monthMap := models.MonthObj{
				Month:    monthElement.Month,
				DayPrice: day,
			}
			getScheduleDay, err := t.scheduleRepo.GetDayByTransId(ctx, getDetailTrans.Id, schedule.Year, monthMap.Month,trans.Time[0].ArrivalTime,trans.Time[0].DepartureTime)
			if err != nil {
				return nil, err
			}
			for _, dayElement := range getScheduleDay {
				var price models.PriceObj
				var currency string
				if dayElement.Price != "" {
					if errUnmarshal := json.Unmarshal([]byte(dayElement.Price), &price); errUnmarshal != nil {
						return nil, errUnmarshal
					}
				}
				if price.Currency == 1 {
					currency = "USD"
				} else {
					currency = "IDR"
				}

				dayMap := models.DayPriceObj{
					DepartureDate: dayElement.DepartureDate.Format("2006-01-02"),
					Day:           dayElement.Day,
					AdultPrice:    price.AdultPrice,
					ChildrenPrice: price.ChildrenPrice,
					Currency:      currency,
				}
				monthMap.DayPrice = append(monthMap.DayPrice, dayMap)
			}
			schedule.Month = append(schedule.Month, monthMap)
		}
		trans.Schedule = append(trans.Schedule, schedule)
	}
	transObj = trans
	var boatDetails models.BoatDetailsObj
	if getDetailTrans.BoatDetails != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailTrans.BoatDetails), &boatDetails); errUnmarshal != nil {
			return nil, errUnmarshal
		}
	}
	facilities := make([]models.ExpFacilitiesObject, 0)
	//if getDetailTrans.TransFacilities != nil {
	//	if errUnmarshal := json.Unmarshal([]byte(*getDetailTrans.TransFacilities), &facilities); errUnmarshal != nil {
	//		return nil, errUnmarshal
	//	}
	//}
	getFacilities, err := t.expFacilitiesRepo.GetJoin(ctx, "", getDetailTrans.Id)
	for _, element := range getFacilities {
		facility := models.ExpFacilitiesObject{
			Name:   element.FacilityName,
			Icon:   *element.FacilityIcon,
			Amount: element.Amount,
		}
		facilities = append(facilities, facility)
	}
	var transImage []models.CoverPhotosObj
	if getDetailTrans.TransImages != "" {
		if errUnmarshal := json.Unmarshal([]byte(getDetailTrans.TransImages), &transImage); errUnmarshal != nil {
			return nil, errUnmarshal
		}
	}
	result := models.TransportationDto{
		Id:              getDetailTrans.Id,
		TransName:       getDetailTrans.TransName,
		TransCapacity:   getDetailTrans.TransCapacity,
		TransTitle:      getDetailTrans.TransTitle,
		Status:          getDetailTrans.TransStatus,
		BoatDetails:     boatDetails,
		Transcoverphoto: getDetailTrans.Transcoverphoto,
		Class:           getDetailTrans.Class,
		Facilities:      facilities,
		TransImages:     transImage,
		DepartureRoute:  transObj,
		ReturnRoute:     &returnTransObj,
	}
	return &result, nil
}
func (t transportationUsecase) UpdateStatus(ctx context.Context, status int, id string, token string) (*models.NewCommandChangeStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	currentMerchant, err := t.merchantUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	errorUpdate := t.transportationRepo.UpdateStatus(ctx, status, id, currentMerchant.MerchantEmail)
	if errorUpdate != nil {
		return nil, errorUpdate
	}
	result := models.NewCommandChangeStatus{
		ExpId:   "",
		TransId: id,
		Status:  status,
	}
	return &result, nil
}
func (t transportationUsecase) FilterSearchTrans(
	ctx context.Context,
	isMerchant bool,
	token,
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
	returnTransId string,
	notReturn string,
	currencyPrice string,
) (*models.FilterSearchTransWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, t.contextTimeout)
	defer cancel()

	var query string
	var queryCount string

	if isMerchant == true && qStatus == "draft" {
		query = `
	SELECT 
		(SELECT id FROM schedules where trans_id = t.id LIMIT 0,1) as schedule_id,
		(SELECT departure_date FROM schedules where trans_id = t.id LIMIT 0,1) as departure_date,
		(SELECT departure_time FROM schedules where trans_id = t.id LIMIT 0,1) as departure_time,
		(SELECT arrival_time FROM schedules where trans_id = t.id LIMIT 0,1) as arrival_time,
        (SELECT price FROM schedules where trans_id = t.id LIMIT 0,1) as price,
		t.id as trans_id,
		t.trans_name,
		t.trans_images,
		t.trans_status,
		(SELECT id FROM harbors where id = t.harbors_source_id LIMIT 0,1) as harbor_source_id,
		(SELECT harbors_name FROM harbors where id = t.harbors_source_id LIMIT 0,1) as harbor_source_name,
		(SELECT id FROM harbors where id = t.harbors_dest_id LIMIT 0,1) as harbor_dest_id,
		(SELECT harbors_name FROM harbors where id = t.harbors_dest_id LIMIT 0,1) as harbor_dest_name,
		m.merchant_name,
		m.merchant_picture,
		t.class,
		t.trans_facilities,
		t.trans_capacity
	FROM
		transportations t
		JOIN merchants m on t.merchant_id = m.id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1
		AND t.is_return = 0`

		queryCount = `
	SELECT
		COUNT(t.id)
	FROM
		transportations t
		JOIN merchants m on t.merchant_id = m.id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1
        AND t.is_return = 0`
	} else if isMerchant == true {
		query = `
	SELECT 
		(SELECT id FROM schedules where trans_id = t.id LIMIT 0,1) as schedule_id,
		(SELECT departure_date FROM schedules where trans_id = t.id LIMIT 0,1) as departure_date,
		(SELECT departure_time FROM schedules where trans_id = t.id LIMIT 0,1) as departure_time,
		(SELECT arrival_time FROM schedules where trans_id = t.id LIMIT 0,1) as arrival_time,
        (SELECT price FROM schedules where trans_id = t.id LIMIT 0,1) as price,
		t.id as trans_id,
		t.trans_name,
		t.trans_images,
		t.trans_status,
		h.id as harbor_source_id,
		h.harbors_name as harbor_source_name,
		hdest.id as harbor_dest_id,
		hdest.harbors_name as harbor_dest_name,
		m.merchant_name,
		m.merchant_picture,
		t.class,
		t.trans_facilities,
		t.trans_capacity,
		csource.id as city_source_id,
		csource.city_name as city_source_name,
		cdest.id as city_dest_id,
		cdest.city_name as city_dest_name,
		psource.id as province_source_id,
		psource.province_name as province_source_name,
		pdest.id as province_dest_id,
		pdest.province_name as province_dest_name,
		t.boat_details,
		t.return_trans_id
	FROM
		transportations t 
		JOIN harbors h ON t.harbors_source_id = h.id
		JOIN harbors hdest ON t.harbors_dest_id = hdest.id
		JOIN cities csource ON h.city_id = csource.id
		JOIN cities cdest ON hdest.city_id = cdest.id
		JOIN provinces psource ON csource.province_id = psource.id
		JOIN provinces pdest ON cdest.province_id = pdest.id
		JOIN merchants m on t.merchant_id = m.id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1
		AND t.is_return = 0`

		queryCount = `
	SELECT
		COUNT(t.id)
	FROM
		transportations t 
		JOIN harbors h ON t.harbors_source_id = h.id
		JOIN harbors hdest ON t.harbors_dest_id = hdest.id
		JOIN cities csource ON h.city_id = csource.id
		JOIN cities cdest ON hdest.city_id = cdest.id
		JOIN provinces psource ON csource.province_id = psource.id
		JOIN provinces pdest ON cdest.province_id = pdest.id
		JOIN merchants m on t.merchant_id = m.id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1
		AND t.is_return = 0`
	} else {
		query = `
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
		hdest.harbors_name as harbor_dest_name,
		m.merchant_name,
		m.merchant_picture,
		t.class,
		t.trans_facilities,
		t.trans_capacity,
		csource.id as city_source_id,
		csource.city_name as city_source_name,
		cdest.id as city_dest_id,
		cdest.city_name as city_dest_name,
		psource.id as province_source_id,
		psource.province_name as province_source_name,
		pdest.id as province_dest_id,
		pdest.province_name as province_dest_name,
		t.boat_details,
		t.return_trans_id
	FROM
		schedules s
		JOIN transportations t ON s.trans_id = t.id
		JOIN harbors h ON t.harbors_source_id = h.id
		JOIN harbors hdest ON t.harbors_dest_id = hdest.id
		JOIN cities csource ON h.city_id = csource.id
		JOIN cities cdest ON hdest.city_id = cdest.id
		JOIN provinces psource ON csource.province_id = psource.id
		JOIN provinces pdest ON cdest.province_id = pdest.id
		JOIN merchants m on t.merchant_id = m.id
	WHERE
	t.is_deleted = 0
	AND t.is_active = 1`

		queryCount = `
	SELECT
		COUNT(*)
	FROM
		schedules s
		JOIN transportations t ON s.trans_id = t.id
		JOIN harbors h ON t.harbors_source_id = h.id
		JOIN harbors hdest ON t.harbors_dest_id = hdest.id
		JOIN cities csource ON h.city_id = csource.id
		JOIN cities cdest ON hdest.city_id = cdest.id
		JOIN provinces psource ON csource.province_id = psource.id
		JOIN provinces pdest ON cdest.province_id = pdest.id
		JOIN merchants m on t.merchant_id = m.id
	WHERE
	t.is_deleted = 0
	AND t.is_active = 1`
	}

	if isMerchant {
		if token == "" {
			return nil, models.ErrUnAuthorize
		}

		currentMerchant, err := t.merchantUsecase.ValidateTokenMerchant(ctx, token)
		if err != nil {
			return nil, err
		}

		query = query + ` AND t.merchant_id = '` + currentMerchant.Id + `'`
		queryCount = queryCount + ` AND t.merchant_id = '` + currentMerchant.Id + `'`
	}

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

		if qStatus == "inService" {
			query = query + ` AND t.trans_status IN (2,3)`
			queryCount = queryCount + ` AND t.trans_status IN (2,3)`
		} else {
			query = query + ` AND t.trans_status =` + strconv.Itoa(status)
			queryCount = queryCount + ` AND t.trans_status =` + strconv.Itoa(status)
		}

	}
	if search != "" {
		if isMerchant == true && qStatus == "draft" {
			keyword := `'%` + search + `%'`
			query = query + ` AND (LOWER(t.trans_name) LIKE LOWER(` + keyword + `))`
			queryCount = queryCount + ` AND (LOWER(t.trans_name) LIKE LOWER(` + keyword + `))`
		} else {
			keyword := `'%` + search + `%'`
			query = query + ` AND (LOWER(t.trans_name) LIKE LOWER(` + keyword + `) OR LOWER(h.harbors_name) LIKE LOWER(` + keyword + `) OR LOWER(hdest.harbors_name) LIKE LOWER(` + keyword + `))`
			queryCount = queryCount + ` AND (LOWER(t.trans_name) LIKE LOWER(` + keyword + `) OR LOWER(h.harbors_name) LIKE LOWER(` + keyword + `) OR LOWER(hdest.harbors_name) LIKE LOWER(` + keyword + `))`
		}

	}
	if harborSourceId != "" {
		query = query + ` AND t.harbors_source_id = '` + harborSourceId + `'`
		queryCount = queryCount + ` AND t.harbors_source_id = '` + harborSourceId + `'`
	}
	if harborDestId != "" {
		query = query + ` AND t.harbors_dest_id = '` + harborDestId + `'`
		queryCount = queryCount + ` AND t.harbors_dest_id = '` + harborDestId + `'`
	}
	if isReturn {
		query = query + ` AND t.is_return = 0 AND t.return_trans_id	!= '' `
		queryCount = queryCount + ` AND t.is_return = 0 AND t.return_trans_id != '' `
	} else if notReturn != "" {
		query = query + ` AND t.is_return = 0 `
		queryCount = queryCount + ` AND t.is_return = 0 `
	}
	if returnTransId != "" {
		query = query + ` AND t.id = '` + returnTransId + `'`
		queryCount = queryCount + ` AND t.id  = '` + returnTransId + `'`
	}
	//if guest != 0 {
	//	query = query + ` AND t.trans_capacity <=` + strconv.Itoa(guest)
	//	queryCount = queryCount + ` AND t.trans_capacity <=` + strconv.Itoa(guest)
	//}
	if depDate != "" {
		query = query + ` AND s.departure_date = '` + depDate + `'`
		queryCount = queryCount + ` AND s.departure_date = '` + depDate + `'`
	}
	if class != "" {
		query = query + ` AND t.class = '` + class + `'`
		queryCount = queryCount + ` AND t.class = '` + class + `'`
	}
	if depTimeOptions != 0 {
		query = query + ` AND s.departure_timeoption_id =` + strconv.Itoa(depTimeOptions)
		queryCount = queryCount + ` AND s.departure_timeoption_id =` + strconv.Itoa(depTimeOptions)
	}
	if arrTimeOptions != 0 {
		query = query + ` AND s.departure_timeoption_id =` + strconv.Itoa(arrTimeOptions)
		queryCount = queryCount + ` AND s.departure_timeoption_id =` + strconv.Itoa(arrTimeOptions)
	}

	if sortBy != "" {
		if sortBy == "newest" {
			query = query + ` ORDER BY t.created_date DESC`
		} else if sortBy == "latest" {
			query = query + ` ORDER BY t.created_date ASC`
		}
	}
	transList, err := t.transportationRepo.FilterSearch(ctx, query, limit, offset, isMerchant, qStatus)
	if err != nil {
		return nil, err
	}

	trans := make([]*models.TransportationSearchObj, 0)
	for _, element := range transList {

		var boatDetails models.BoatDetailsObj
		if element.BoatDetails != "" {
			if errUnmarshal := json.Unmarshal([]byte(element.BoatDetails), &boatDetails); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}

		var transImages []models.CoverPhotosObj
		if element.TransImages != "" {
			if errUnmarshal := json.Unmarshal([]byte(element.TransImages), &transImages); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}
		var transPrice models.TransPriceObj
		if element.Price != nil {
			if errUnmarshal := json.Unmarshal([]byte(*element.Price), &transPrice); errUnmarshal != nil {
				return nil, errUnmarshal
			}
			transPrice.PriceType = "pax"
			if transPrice.Currency == 1 {
				transPrice.CurrencyLabel = "USD"
			} else {
				transPrice.CurrencyLabel = "IDR"
			}
			if currencyPrice == "USD"{
				if transPrice.CurrencyLabel == "IDR"{
					convertCurrency ,_ := t.currencyUsecase.ExchangeRatesApi(ctx,"IDR","USD")
					calculatePriceAdult := convertCurrency.Rates.USD * transPrice.AdultPrice
					transPrice.AdultPrice = calculatePriceAdult
					calculatePriceChildren := convertCurrency.Rates.USD * transPrice.ChildrenPrice
					transPrice.ChildrenPrice = calculatePriceChildren
					transPrice.CurrencyLabel = "USD"
				}
			}else if currencyPrice =="IDR"{
				if transPrice.CurrencyLabel == "USD"{
					convertCurrency ,_ := t.currencyUsecase.ExchangeRatesApi(ctx,"USD","IDR")
					calculatePriceAdult := convertCurrency.Rates.USD * transPrice.AdultPrice
					transPrice.AdultPrice = calculatePriceAdult
					calculatePriceChildren := convertCurrency.Rates.USD * transPrice.ChildrenPrice
					transPrice.ChildrenPrice = calculatePriceChildren
					transPrice.CurrencyLabel = "IDR"
				}
			}
		}

		var tripDuration string
		if element.DepartureTime != nil && element.ArrivalTime != nil {
			departureTime, _ := time.Parse("15:04:00", *element.DepartureTime)
			arrivalTime, _ := time.Parse("15:04:00", *element.ArrivalTime)

			tripHour := arrivalTime.Hour() - departureTime.Hour()
			tripMinute := arrivalTime.Minute() - departureTime.Minute()
			tripDuration = strconv.Itoa(tripHour) + `h ` + strconv.Itoa(tripMinute) + `m`
		}

		var transStatus string
		if element.TransStatus == 0 {
			transStatus = "Preview"
		} else if element.TransStatus == 1 {
			transStatus = "Draft"
		} else if element.TransStatus == 2 {
			transStatus = "Published"
		} else if element.TransStatus == 3 {
			transStatus = "Unpublished"
		} else if element.TransStatus == 4 {
			transStatus = "Archived"
		}
		transFacilities := make([]models.ExpFacilitiesObject, 0)
		if element.TransFacilities != nil {
			if errUnmarshal := json.Unmarshal([]byte(*element.TransFacilities), &transFacilities); errUnmarshal != nil {
				return nil, errUnmarshal
			}
		}
		transDto := &models.TransportationSearchObj{
			ScheduleId:            element.ScheduleId,
			DepartureDate:         element.DepartureDate,
			DepartureTime:         element.DepartureTime,
			ArrivalTime:           element.ArrivalTime,
			TripDuration:          &tripDuration,
			TransportationId:      element.TransId,
			TransportationName:    element.TransName,
			TransportationImages:  transImages,
			TransportationStatus:  transStatus,
			HarborSourceId:        element.HarborSourceId,
			HarborSourceName:      element.HarborSourceName,
			HarborDestinationId:   element.HarborDestId,
			HarborDestinationName: element.HarborDestName,
			Price:                 transPrice,
			MerchantName:          element.MerchantName,
			MerchantPicture:       element.MerchantPicture,
			Class:                 element.Class,
			TransFacilities:       transFacilities,
			CitySourceId:          element.CitySourceId,
			CitySourceName:        element.CitySourceName,
			CityDestId:            element.CityDestId,
			CityDestName:          element.CityDestName,
			ProvinceSourceId:      element.ProvinceSourceId,
			ProvinceSourceName:    element.ProvinceSourceName,
			ProvinceDestId:        element.ProvinceDestId,
			ProvinceDestName:      element.ProvinceDestName,
			BoatSpecification:     boatDetails,
			ReturnTransId:         element.ReturnTransId,
		}
		if guest != 0 {
			getbooking, _ := t.transactionRepo.GetCountByTransId(ctx, element.TransId,false,"")
			getCapacity := element.TransCapacity
			var getbookingCount int
			if getbooking != nil {
				for _,booking := range getbooking{
					guestDesc := make([]models.GuestDescObj, 0)
					if errUnmarshal := json.Unmarshal([]byte(*booking), &guestDesc); errUnmarshal != nil {
						return nil, models.ErrInternalServerError
					}
					getbookingCount = getbookingCount + len(guestDesc)
				}
			}
			remainingSeat := getCapacity - getbookingCount

			if guest <= remainingSeat {
				trans = append(trans, transDto)
			}
		} else {
			trans = append(trans, transDto)
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
	var harborsSourceId *string
	var harborsDestId *string
	if len(newCommandTransportation.TransImages) != 0 {
		transImagesConvert, _ := json.Marshal(newCommandTransportation.TransImages)
		transImages = string(transImagesConvert)
	}
	boatDetails, _ := json.Marshal(newCommandTransportation.BoatDetails)
	facilites, _ := json.Marshal(newCommandTransportation.Facilities)
	harborsDestId = newCommandTransportation.DepartureRoute.HarborsIdTo
	harborsSourceId = newCommandTransportation.DepartureRoute.HarborsIdFrom
	isReturn := 0
	var transFacilities string
	transFacilities = string(facilites)
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
		TransFacilities: &transFacilities,
		IsReturn:        &isReturn,
	}
	if newCommandTransportation.ReturnRoute != nil {
		if newCommandTransportation.ReturnRoute.HarborsIdFrom != nil && newCommandTransportation.ReturnRoute.HarborsIdTo != nil {
			harborsDestReturnId := newCommandTransportation.ReturnRoute.HarborsIdTo
			harborsSourceReturnId := newCommandTransportation.ReturnRoute.HarborsIdFrom
			isReturnn := 1
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
				TransFacilities: &transFacilities,
				IsReturn:        &isReturnn,
			}
			insertTransportationReturn, err := t.transportationRepo.Insert(ctx, transportationReturn)
			if err != nil {
				return nil, err
			}

			for _, element := range newCommandTransportation.Facilities {
				getFacilitiesId, err := t.facilitiesRepo.GetByName(ctx, element.Name)
				if err != nil {
					return nil, err
				}
				facilities := models.ExperienceFacilities{
					Id:           0,
					ExpId:        nil,
					TransId:      insertTransportationReturn,
					FacilitiesId: getFacilitiesId.Id,
					Amount:       element.Amount,
				}
				err = t.expFacilitiesRepo.Insert(ctx, &facilities)
				if err != nil {
					return nil, err
				}
			}

			schedule := models.NewCommandSchedule{
				TimeObj:   newCommandTransportation.ReturnRoute.Time,
				CreatedBy: currentUserMerchant.MerchantEmail,
				TransId:   *insertTransportationReturn,
				Schedule:  newCommandTransportation.ReturnRoute.Schedule,
			}
			events.UserCreated.Trigger(events.UserCreatedPayload{
				Schedule: schedule,
			})

			transportation.ReturnTransId = insertTransportationReturn
		}
	}
	insertTransportation, err := t.transportationRepo.Insert(ctx, transportation)
	if err != nil {
		return nil, err
	}

	for _, element := range newCommandTransportation.Facilities {
		getFacilitiesId, err := t.facilitiesRepo.GetByName(ctx, element.Name)
		if err != nil {
			return nil, err
		}
		facilities := models.ExperienceFacilities{
			Id:           0,
			ExpId:        nil,
			TransId:      insertTransportation,
			FacilitiesId: getFacilitiesId.Id,
			Amount:       element.Amount,
		}
		err = t.expFacilitiesRepo.Insert(ctx, &facilities)
		if err != nil {
			return nil, err
		}
	}

	schedule := models.NewCommandSchedule{
		TimeObj:   newCommandTransportation.DepartureRoute.Time,
		CreatedBy: currentUserMerchant.MerchantEmail,
		TransId:   *insertTransportation,
		Schedule:  newCommandTransportation.DepartureRoute.Schedule,
	}
	events.UserCreated.Trigger(events.UserCreatedPayload{
		Schedule: schedule,
	})

	var status string
	if newCommandTransportation.Status == 1 {
		status = "Draft"
	} else if newCommandTransportation.Status == 2 {
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
	var harborsSourceId *string
	var harborsDestId *string
	if len(newCommandTransportation.TransImages) != 0 {
		transImagesConvert, _ := json.Marshal(newCommandTransportation.TransImages)
		transImages = string(transImagesConvert)
	}
	facilites, _ := json.Marshal(newCommandTransportation.Facilities)
	boatDetails, _ := json.Marshal(newCommandTransportation.BoatDetails)
	harborsDestId = newCommandTransportation.DepartureRoute.HarborsIdTo
	harborsSourceId = newCommandTransportation.DepartureRoute.HarborsIdFrom
	var transFacilities string
	transFacilities = string(facilites)
	isReturn := 0
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
		TransFacilities: &transFacilities,
		IsReturn:        &isReturn,
	}
	if newCommandTransportation.ReturnRoute != nil {
		if newCommandTransportation.ReturnRoute.HarborsIdFrom != nil && newCommandTransportation.ReturnRoute.HarborsIdTo != nil {
			harborsDestReturnId := newCommandTransportation.ReturnRoute.HarborsIdTo
			harborsSourceReturnId := newCommandTransportation.ReturnRoute.HarborsIdFrom
			isReturnn := 1
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
				TransFacilities: &transFacilities,
				IsReturn:        &isReturnn,
			}
			insertTransportationReturn, err := t.transportationRepo.Update(ctx, transportationReturn)
			if err != nil {
				return nil, err
			}
			err = t.expFacilitiesRepo.Delete(ctx, "", *insertTransportationReturn)
			if err != nil {
				return nil, err
			}
			for _, element := range newCommandTransportation.Facilities {
				getFacilitiesId, err := t.facilitiesRepo.GetByName(ctx, element.Name)
				if err != nil {
					return nil, err
				}
				facilities := models.ExperienceFacilities{
					Id:           0,
					ExpId:        nil,
					TransId:      insertTransportationReturn,
					FacilitiesId: getFacilitiesId.Id,
					Amount:       element.Amount,
				}
				err = t.expFacilitiesRepo.Insert(ctx, &facilities)
				if err != nil {
					return nil, err
				}
			}
			errorDelete := t.scheduleRepo.DeleteByTransId(ctx, insertTransportationReturn)
			if errorDelete != nil {
				return nil, errorDelete
			}
			schedule := models.NewCommandSchedule{
				TimeObj:   newCommandTransportation.ReturnRoute.Time,
				CreatedBy: currentUserMerchant.MerchantEmail,
				TransId:   *insertTransportationReturn,
				Schedule:  newCommandTransportation.ReturnRoute.Schedule,
			}
			events.UserCreated.Trigger(events.UserCreatedPayload{
				Schedule: schedule,
			})

			transportation.ReturnTransId = insertTransportationReturn
		}
	}
	insertTransportation, err := t.transportationRepo.Update(ctx, transportation)
	if err != nil {
		return nil, err
	}
	err = t.expFacilitiesRepo.Delete(ctx, "", *insertTransportation)
	if err != nil {
		return nil, err
	}
	for _, element := range newCommandTransportation.Facilities {
		getFacilitiesId, err := t.facilitiesRepo.GetByName(ctx, element.Name)
		if err != nil {
			return nil, err
		}
		facilities := models.ExperienceFacilities{
			Id:           0,
			ExpId:        nil,
			TransId:      insertTransportation,
			FacilitiesId: getFacilitiesId.Id,
			Amount:       element.Amount,
		}
		err = t.expFacilitiesRepo.Insert(ctx, &facilities)
		if err != nil {
			return nil, err
		}
	}
	errorDelete := t.scheduleRepo.DeleteByTransId(ctx, insertTransportation)
	if errorDelete != nil {
		return nil, errorDelete
	}

	schedule := models.NewCommandSchedule{
		TimeObj:   newCommandTransportation.DepartureRoute.Time,
		CreatedBy: currentUserMerchant.MerchantEmail,
		TransId:   *insertTransportation,
		Schedule:  newCommandTransportation.DepartureRoute.Schedule,
	}
	events.UserCreated.Trigger(events.UserCreatedPayload{
		Schedule: schedule,
	})

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

func (t transportationUsecase) scheduleMonthInsert(ctx context.Context, months []models.MonthObj, timeObj []models.TimeObj, year models.YearObj, insertTransportationReturn *string, createdBy string) error {
	for _, month := range months {
		go t.scheduleDayInsert(ctx, month.DayPrice, timeObj, year, insertTransportationReturn, createdBy, month)
	}
	return nil
}
func (t transportationUsecase) scheduleDayInsert(ctx context.Context, dayPrices []models.DayPriceObj, timeObj []models.TimeObj, year models.YearObj, insertTransportationReturn *string, createdBy string, month models.MonthObj) error {
	for _, day := range dayPrices {
		go t.scheduleTimeInsert(ctx, day, timeObj, year, insertTransportationReturn, createdBy, month)
	}
	return nil
}
func (t transportationUsecase) scheduleTimeInsert(ctx context.Context, day models.DayPriceObj, timeObj []models.TimeObj, year models.YearObj, insertTransportationReturn *string, createdBy string, month models.MonthObj) error {
	for _, times := range timeObj {
		go t.scheduleLastInsert(ctx, day, times, year, insertTransportationReturn, createdBy, month)
	}
	return nil
}

func (t transportationUsecase) scheduleLastInsert(ctx context.Context, day models.DayPriceObj, times models.TimeObj, year models.YearObj, insertTransportationReturn *string, createdBy string, month models.MonthObj) error {
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
		return err
	}
	arrivalTimeOption, err := t.timeOptionsRepo.GetByTime(ctx, times.ArrivalTime)
	if err != nil {
		return err
	}
	price, _ := json.Marshal(priceObj)
	schedule := models.Schedule{
		Id:                    "",
		CreatedBy:             createdBy,
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
	//_, err = t.scheduleRepo.Insert(ctx, schedule)
	//if err != nil {
	//	return err
	//}
	go t.scheduleRepo.Insert(ctx, schedule)
	return nil
}
