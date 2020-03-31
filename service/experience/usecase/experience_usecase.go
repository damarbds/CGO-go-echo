package usecase

import (
	"context"
	"encoding/json"
	"github.com/service/cpc"
	"github.com/service/harbors"
	"time"

	"github.com/models"
	"github.com/service/experience"
)

type experienceUsecase struct {
	experienceRepo   experience.Repository
	harborsRepo 	harbors.Repository
	cpcRepo 		cpc.Repository
	contextTimeout   time.Duration
}

// NewexperienceUsecase will create new an experienceUsecase object representation of experience.Usecase interface
func NewexperienceUsecase(a experience.Repository, h harbors.Repository,c cpc.Repository,timeout time.Duration) experience.Usecase {
	return &experienceUsecase{
		experienceRepo:   a,
		harborsRepo:	h,
		cpcRepo:	c,
		contextTimeout:   timeout,
	}
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
