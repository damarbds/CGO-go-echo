package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/product/reviews"
	"github.com/service/cpc"
	"github.com/service/harbors"
	"strconv"
	"time"

	"github.com/models"
	payment "github.com/service/exp_payment"
	"github.com/service/experience"
)

type experienceUsecase struct {
	experienceRepo experience.Repository
	harborsRepo    harbors.Repository
	cpcRepo        cpc.Repository
	paymentRepo    payment.Repository
	reviewsRepo    reviews.Repository
	contextTimeout time.Duration
}


// NewexperienceUsecase will create new an experienceUsecase object representation of experience.Usecase interface
func NewexperienceUsecase(
	a experience.Repository,
	h harbors.Repository,
	c cpc.Repository,
	p payment.Repository,
	r reviews.Repository,
	timeout time.Duration,
) experience.Usecase {
	return &experienceUsecase{
		experienceRepo:   a,
		harborsRepo:	h,
		cpcRepo:	c,
		paymentRepo: p,
		reviewsRepo: r,
		contextTimeout:   timeout,
	}
}

func (m experienceUsecase) GetUserDiscoverPreference(ctx context.Context,page *int,size *int) ([]*models.ExpUserDiscoverPreferenceDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	expList, err := m.experienceRepo.GetUserDiscoverPreference(ctx,page,size)
	if err != nil {
		return nil, err
	}
	var expListDto []*models.ExpUserDiscoverPreferenceDto

	for _, element := range expList {

		var	expType []string
		if errUnmarshal := json.Unmarshal([]byte(element.ExpType), &expType); errUnmarshal != nil {
			return nil,models.ErrInternalServerError
		}
		countRating, err := m.reviewsRepo.CountRating(ctx, element.Id)
		if err != nil {
			return nil, err
		}
		expPayment, err := m.paymentRepo.GetByExpID(ctx, element.Id)
		if err != nil {
			return nil, err
		}

		var currency string
		if expPayment.Currency == 1 {
			currency = "USD"
		} else {
			currency = "IDR"
		}
		var priceItemType string
		if expPayment.PriceItemType == 1 {
			priceItemType = "Per Pax"
		} else {
			priceItemType = "Per Trip"
		}
		//expDto := models.ExperienceUserDiscoverPreferenceDto{}

		if len(expListDto) == 0 {
			cityDto := models.ExpUserDiscoverPreferenceDto{
				CityId: element.CityId,
				City:     element.CityName,
				CityDesc: element.CityDesc,
				Item:     nil,
			}
			expDto := models.ExperienceUserDiscoverPreferenceDto{
				Id:           element.Id,
				ExpTitle:     element.ExpTitle,
				ExpType:      expType,
				Rating:       element.Rating,
				CountRating:  countRating,
				Currency:     currency,
				Price:        expPayment.Price,
				Payment_type: priceItemType,
			}
			cityDto.Item = append(cityDto.Item,expDto)
			expListDto = append(expListDto,&cityDto)
		}else if len(expListDto) != 0 {
			for _, dto := range expListDto {
				if dto.CityId == element.CityId{
					expDto := models.ExperienceUserDiscoverPreferenceDto{
						Id:           element.Id,
						ExpTitle:     element.ExpTitle,
						ExpType:      expType,
						Rating:       element.Rating,
						CountRating:  countRating,
						Currency:     currency,
						Price:        expPayment.Price,
						Payment_type: priceItemType,
					}
					dto.Item = append(dto.Item,expDto)
				}else if dto.CityId != element.CityId{
					cityDto := models.ExpUserDiscoverPreferenceDto{
						CityId: element.CityId,
						City:     element.CityName,
						CityDesc: element.CityDesc,
						Item:     nil,
					}
					expDto := models.ExperienceUserDiscoverPreferenceDto{
						Id:           element.Id,
						ExpTitle:     element.ExpTitle,
						ExpType:      expType,
						Rating:       element.Rating,
						CountRating:  countRating,
						Currency:     currency,
						Price:        expPayment.Price,
						Payment_type: priceItemType,
					}
					cityDto.Item = append(cityDto.Item,expDto)
					expListDto = append(expListDto,&cityDto)
				}
			}
		}

		//expListDto = append(expListDto,&expDto)
	}
	return expListDto, nil
}

func (m experienceUsecase) FilterSearchExp(ctx context.Context, cityID string, harborsId string, activityType string, startDate string, endDate string, guest string, trip string, bottomPrice string, upPrice string) ([]*models.ExpSearchObject, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	query := `select e.id,e.exp_title,e.exp_type,e.rating from experiences e`

	if bottomPrice != "" && upPrice != "" {
		query = query + ` join experience_payments ep on ep.exp_id = e.id`
	}
	//if startDate != "" && endDate != "" {
	//	query = query + ` join exp_availability_date ead on ead.exp_id = e.id`
	//}
	if activityType != "" {
		query = query + ` join filter_activity_types fat on fat.exp_id = e.id`
	}
	if cityID != "" {
		query = query + ` join harbors h on h.id = e.harbors_id`
	}
	if cityID != ""{
		city_id, _ := strconv.Atoi(cityID)
		query = query + ` where h.city_id = ` + strconv.Itoa(city_id)
	}else if harborsId != ""{
		query = query + ` where e.harbors_id = '` + harborsId + `'`
	}else {
		return nil, models.ErrBadParamInput
	}
	if guest != ""{
		guests, _ := strconv.Atoi(guest)
		query = query + ` AND e.exp_max_guest =` + strconv.Itoa(guests)
	}
	if trip != "" {
		trips, _ := strconv.Atoi(trip)
		var tripType string
		if trips == 0 {
			tripType = "Private Trip"
		}else if trips == 1 {
			tripType = "Share Trip"
		}else {
			return nil,models.ErrInternalServerError
		}
		query = query + ` AND e.exp_trip_type = '` + tripType + `'`
	}
	if activityType != "" {
		types, _ := strconv.Atoi(activityType)
		query = query + ` AND fat.id =` + strconv.Itoa(types)
	}
	if bottomPrice != "" && upPrice != "" {
		bottomprices, _ := strconv.ParseFloat(bottomPrice, 64)
		upprices, _ := strconv.ParseFloat(upPrice,64)

		query = query + ` AND (ep.price between ` +  fmt.Sprint(bottomprices) + ` AND ` +  fmt.Sprint(upprices) + `)`
	}
	//if startDate != "" && endDate != "" {
	//
	//	query = query + ` AND exp_availability_date like %` + s
	//}
	expList, err := m.experienceRepo.QueryFilterSearch(ctx, query)
	if err != nil {
		return nil, err
	}
	results := make([]*models.ExpSearchObject, len(expList))
	for i, exp := range expList {
		var	expType []string
		if errUnmarshal := json.Unmarshal([]byte(exp.ExpType), &expType); errUnmarshal != nil {
			return nil,models.ErrInternalServerError
		}

		expPayment, err := m.paymentRepo.GetByExpID(ctx, exp.Id)
		if err != nil {
			return nil, err
		}

		var currency string
		if expPayment.Currency == 1 {
			currency = "USD"
		} else {
			currency = "IDR"
		}

		var priceItemType string
		if expPayment.PriceItemType == 1 {
			priceItemType = "Per Pax"
		} else {
			priceItemType = "Per Trip"
		}

		countRating, err := m.reviewsRepo.CountRating(ctx, exp.Id)
		if err != nil {
			return nil, err
		}

		results[i] = &models.ExpSearchObject{
			Id:          exp.Id,
			ExpTitle:    exp.ExpTitle,
			ExpType:     expType,
			Rating:      exp.Rating,
			CountRating: countRating,
			Currency:    currency,
			Price:       expPayment.Price,
			PaymentType: priceItemType,
		}
	}

	return results, nil

}
func (m experienceUsecase) SearchExp(ctx context.Context, harborID, cityID string) ([]*models.ExpSearchObject, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	expList, err := m.experienceRepo.SearchExp(ctx, harborID, cityID)
	if err != nil {
		return nil, err
	}

	results := make([]*models.ExpSearchObject, len(expList))
	for i, exp := range expList {
		var	expType []string
		if errUnmarshal := json.Unmarshal([]byte(exp.ExpType), &expType); errUnmarshal != nil {
			return nil,models.ErrInternalServerError
		}

		expPayment, err := m.paymentRepo.GetByExpID(ctx, exp.Id)
		if err != nil {
			return nil, err
		}

		var currency string
		if expPayment.Currency == 1 {
			currency = "USD"
		} else {
			currency = "IDR"
		}

		var priceItemType string
		if expPayment.PriceItemType == 1 {
			priceItemType = "Per Pax"
		} else {
			priceItemType = "Per Trip"
		}

		countRating, err := m.reviewsRepo.CountRating(ctx, exp.Id)
		if err != nil {
			return nil, err
		}

		results[i] = &models.ExpSearchObject{
			Id:          exp.Id,
			ExpTitle:    exp.ExpTitle,
			ExpType:     expType,
			Rating:      exp.Rating,
			CountRating: countRating,
			Currency:    currency,
			Price:       expPayment.Price,
			PaymentType: priceItemType,
		}
	}

	return results, nil
}

func (m experienceUsecase)GetByID(c context.Context, id string) (*models.ExperienceDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err := m.experienceRepo.GetByID(ctx,id)
	if err != nil {
		return nil, err
	}
	var	expType []string
	errObject := json.Unmarshal([]byte(res.ExpType), &expType)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil,models.ErrInternalServerError
	}
	expItinerary := models.ExpItineraryObject{}
	errObject = json.Unmarshal([]byte(res.ExpInternary), &expItinerary)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil,models.ErrInternalServerError
	}
	var expFacilities  []models.ExpFacilitiesObject
	errObject = json.Unmarshal([]byte(res.ExpFacilities), &expFacilities)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil,models.ErrInternalServerError
	}
	var expInclusion  []models.ExpInclusionObject
	errObject = json.Unmarshal([]byte(res.ExpInclusion), &expInclusion)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil,models.ErrInternalServerError
	}
	var expRules  []models.ExpRulesObject
	errObject = json.Unmarshal([]byte(res.ExpRules), &expRules)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil,models.ErrInternalServerError
	}
	harbors , err := m.harborsRepo.GetByID(ctx,res.HarborsId)
	city, err := m.cpcRepo.GetCityByID(ctx,harbors.CityId)
	province, err := m.cpcRepo.GetProvinceByID(ctx,city.ProvinceId)
	experiences := models.ExperienceDto{
		Id:                      res.Id,
		ExpTitle:                res.ExpTitle,
		ExpType:                 expType,
		ExpTripType:             res.ExpTripType,
		ExpBookingType:          res.ExpBookingType,
		ExpDesc:                 res.ExpDesc,
		ExpMaxGuest:             res.ExpMaxGuest,
		ExpPickupPlace:          res.ExpPickupPlace,
		ExpPickupTime:           res.ExpPickupTime,
		ExpPickupPlaceLongitude: res.ExpPickupPlaceLongitude,
		ExpPickupPlaceLatitude:  res.ExpPickupPlaceLatitude,
		ExpPickupPlaceMapsName:  res.ExpPickupPlaceMapsName,
		ExpInternary:            expItinerary,
		ExpFacilities:           expFacilities,
		ExpInclusion:            expInclusion,
		ExpRules:                expRules,
		Status:                  res.Status,
		Rating:                  res.Rating,
		ExpLocationLatitude:     res.ExpLocationLatitude,
		ExpLocationLongitude:    res.ExpLocationLongitude,
		ExpLocationName:         res.ExpLocationName,
		ExpCoverPhoto:           res.ExpCoverPhoto,
		ExpDuration:             res.ExpDuration,
		MinimumBookingId:        res.MinimumBookingId,
		MerchantId:              res.MerchantId,
		HarborsName:              harbors.HarborsName,
		City:					city.CityName,
		Province:				province.ProvinceName,
	}
	return &experiences, nil
}


/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
