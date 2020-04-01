package usecase

import (
	"context"
	"encoding/json"
	"github.com/service/cpc"
	"github.com/service/harbors"
	"github.com/service/reviews"
	"time"

	"github.com/models"
	payment "github.com/service/exp_payment"
	types "github.com/service/exp_types"
	"github.com/service/experience"
)

type experienceUsecase struct {
	experienceRepo experience.Repository
	harborsRepo    harbors.Repository
	cpcRepo        cpc.Repository
	paymentRepo    payment.Repository
	reviewsRepo reviews.Repository
	typesRepo types.Repository
	contextTimeout time.Duration
}

// NewexperienceUsecase will create new an experienceUsecase object representation of experience.Usecase interface
func NewexperienceUsecase(
	a experience.Repository,
	h harbors.Repository,
	c cpc.Repository,
	p payment.Repository,
	r reviews.Repository,
	t types.Repository,
	timeout time.Duration,
) experience.Usecase {
	return &experienceUsecase{
		experienceRepo:   a,
		harborsRepo:	h,
		cpcRepo:	c,
		paymentRepo: p,
		reviewsRepo: r,
		typesRepo: t,
		contextTimeout:   timeout,
	}
}

func (m experienceUsecase) GetExpTypes(ctx context.Context) ([]*models.ExpTypeObject, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	results, err := m.typesRepo.GetExpTypes(ctx)
	if err != nil {
		return nil, err
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
