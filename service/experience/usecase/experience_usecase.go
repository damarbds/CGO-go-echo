package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/product/experience_add_ons"

	"github.com/auth/merchant"
	"github.com/product/reviews"
	"github.com/service/cpc"
	"github.com/service/exp_availability"
	"github.com/service/exp_photos"
	"github.com/service/harbors"

	"github.com/models"
	inspiration "github.com/service/exp_inspiration"
	payment "github.com/service/exp_payment"
	types "github.com/service/exp_types"
	"github.com/service/experience"
)

type experienceUsecase struct {
	adOnsRepo        experience_add_ons.Repository
	experienceRepo   experience.Repository
	harborsRepo      harbors.Repository
	cpcRepo          cpc.Repository
	paymentRepo      payment.Repository
	reviewsRepo      reviews.Repository
	typesRepo        types.Repository
	inspirationRepo  inspiration.Repository
	expPhotos        exp_photos.Repository
	mUsecase         merchant.Usecase
	contextTimeout   time.Duration
	exp_availablitiy exp_availability.Repository
}

// NewexperienceUsecase will create new an experienceUsecase object representation of experience.Usecase interface
func NewexperienceUsecase(
	adOns experience_add_ons.Repository,
	ea exp_availability.Repository,
	ps exp_photos.Repository,
	a experience.Repository,
	h harbors.Repository,
	c cpc.Repository,
	p payment.Repository,
	r reviews.Repository,
	t types.Repository,
	i inspiration.Repository,
	m merchant.Usecase,
	timeout time.Duration,
) experience.Usecase {
	return &experienceUsecase{
		adOnsRepo:        adOns,
		exp_availablitiy: ea,
		experienceRepo:   a,
		harborsRepo:      h,
		cpcRepo:          c,
		paymentRepo:      p,
		reviewsRepo:      r,
		typesRepo:        t,
		inspirationRepo:  i,
		mUsecase:         m,
		contextTimeout:   timeout,
		expPhotos:        ps,
	}
}

func (m experienceUsecase) GetExpPendingTransactionCount(ctx context.Context, token string) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	currentMerchant, err := m.mUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	count, err := m.experienceRepo.GetExpPendingTransactionCount(ctx, currentMerchant.Id)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}

func (m experienceUsecase) GetExpFailedTransactionCount(ctx context.Context, token string) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	currentMerchant, err := m.mUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	count, err := m.experienceRepo.GetExpFailedTransactionCount(ctx, currentMerchant.Id)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}

func (m experienceUsecase) GetExpCount(ctx context.Context, token string) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	currentMerchant, err := m.mUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	count, err := m.experienceRepo.GetExpCount(ctx, currentMerchant.Id)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}

func (m experienceUsecase) GetSuccessBookCount(ctx context.Context, token string) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	currentMerchant, err := m.mUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	count, err := m.experienceRepo.GetSuccessBookCount(ctx, currentMerchant.Id)
	if err != nil {
		return nil, err
	}

	return &models.Count{Count: count}, nil
}

func (m experienceUsecase) GetByCategoryID(ctx context.Context, categoryId int) ([]*models.ExpSearchObject, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	expList, err := m.experienceRepo.GetByCategoryID(ctx, categoryId)
	if err != nil {
		return nil, err
	}

	results := make([]*models.ExpSearchObject, len(expList))
	for i, exp := range expList {
		var expType []string
		if errUnmarshal := json.Unmarshal([]byte(exp.ExpType), &expType); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}

		expPayment, err := m.paymentRepo.GetByExpID(ctx, exp.Id)
		if err != nil {
			return nil, err
		}

		var currency string
		if expPayment[0].Currency == 1 {
			currency = "USD"
		} else {
			currency = "IDR"
		}

		var priceItemType string
		if expPayment[0].PriceItemType == 1 {
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
			Price:       expPayment[0].Price,
			PaymentType: priceItemType,
		}
	}

	return results, nil
}

func (m experienceUsecase) GetUserDiscoverPreference(ctx context.Context, page *int, size *int) ([]*models.ExpUserDiscoverPreferenceDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	expList, err := m.experienceRepo.GetUserDiscoverPreference(ctx, page, size)
	if err != nil {
		return nil, err
	}
	var expListDto []*models.ExpUserDiscoverPreferenceDto

	for _, element := range expList {

		var expType []string
		if errUnmarshal := json.Unmarshal([]byte(element.ExpType), &expType); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
		countRating, err := m.reviewsRepo.CountRating(ctx, element.Id)
		if err != nil {
			return nil, err
		}
		//expPhotos, err := m.expPhotos.GetByExperienceID(ctx,element.Id)
		//if err != nil {
		//	return nil, models.ErrInternalServerError
		//}

		var coverPhotos models.CoverPhotosObj
		var cityPhotos []models.CoverPhotosObj
		if element.ExpCoverPhoto != nil {
			covertPhoto := models.CoverPhotosObj{
				Original:  *element.ExpCoverPhoto,
				Thumbnail: "",
			}
			coverPhotos = covertPhoto
			//if errUnmarshal := json.Unmarshal([]byte(expPhotos[0].ExpPhotoImage), &coverPhotos); errUnmarshal != nil {
			//	return nil,models.ErrInternalServerError
			//}
		}
		if element.CityPhotos != nil {
			if errUnmarshal := json.Unmarshal([]byte(*element.CityPhotos), &cityPhotos); errUnmarshal != nil {
				return nil, models.ErrInternalServerError
			}
		}

		expPayment, err := m.paymentRepo.GetByExpID(ctx, element.Id)
		if err != nil {
			return nil, err
		}

		var priceItemType string
		var currency string
		if expPayment != nil {
			if expPayment[0].Currency == 1 {
				currency = "USD"
			} else {
				currency = "IDR"
			}

			if expPayment[0].PriceItemType == 1 {
				priceItemType = "Per Pax"
			} else {
				priceItemType = "Per Trip"
			}
		} else {
			priceItemType = ""
			currency = ""
		}

		//expDto := models.ExperienceUserDiscoverPreferenceDto{}

		if len(expListDto) == 0 {
			cityDto := models.ExpUserDiscoverPreferenceDto{
				CityId:     element.CityId,
				City:       element.CityName,
				CityDesc:   element.CityDesc,
				Item:       nil,
				CityPhotos: cityPhotos,
			}
			expDto := models.ExperienceUserDiscoverPreferenceDto{
				Id:           element.Id,
				ExpTitle:     element.ExpTitle,
				ExpType:      expType,
				Rating:       element.Rating,
				CountRating:  countRating,
				Currency:     currency,
				Price:        expPayment[0].Price,
				Payment_type: priceItemType,
				Cover_Photo:  coverPhotos,
			}
			cityDto.Item = append(cityDto.Item, expDto)
			expListDto = append(expListDto, &cityDto)
		} else if len(expListDto) != 0 {
			var searchDto *models.ExpUserDiscoverPreferenceDto
			for _, dto := range expListDto {
				if dto.CityId == element.CityId {
					searchDto = dto
				}
			}
			if searchDto == nil {
				cityDto := models.ExpUserDiscoverPreferenceDto{
					CityId:     element.CityId,
					City:       element.CityName,
					CityDesc:   element.CityDesc,
					CityPhotos: cityPhotos,
					Item:       nil,
				}
				expDto := models.ExperienceUserDiscoverPreferenceDto{
					Id:           element.Id,
					ExpTitle:     element.ExpTitle,
					ExpType:      expType,
					Rating:       element.Rating,
					CountRating:  countRating,
					Currency:     currency,
					Price:        expPayment[0].Price,
					Payment_type: priceItemType,
					Cover_Photo:  coverPhotos,
				}
				cityDto.Item = append(cityDto.Item, expDto)
				expListDto = append(expListDto, &cityDto)
			} else {
				for _, dto := range expListDto {
					if dto.CityId == element.CityId {
						expDto := models.ExperienceUserDiscoverPreferenceDto{
							Id:           element.Id,
							ExpTitle:     element.ExpTitle,
							ExpType:      expType,
							Rating:       element.Rating,
							CountRating:  countRating,
							Currency:     currency,
							Price:        expPayment[0].Price,
							Payment_type: priceItemType,
							Cover_Photo:  coverPhotos,
						}
						dto.Item = append(dto.Item, expDto)
					}
				}
			}
		}

		//expListDto = append(expListDto,&expDto)
	}
	return expListDto, nil
}

func (m experienceUsecase) FilterSearchExp(
	ctx context.Context,
	cityID string,
	harborsId string,
	activityType string,
	startDate string,
	endDate string,
	guest string,
	trip string,
	bottomPrice string,
	upPrice string,
	sortBy string,
) ([]*models.ExpSearchObject, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	var activityTypeArray []int
	if activityType != "" {
		if errUnmarshal := json.Unmarshal([]byte(activityType), &activityTypeArray); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	query := `
	select
		e.id,
		e.exp_title,
		e.exp_type,
		e.rating,
		e.exp_location_latitude as latitude,
		e.exp_location_longitude as longitude,
		e.exp_cover_photo as cover_photo 
	from 
		experiences e`

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
	if cityID != "" {
		city_id, _ := strconv.Atoi(cityID)
		query = query + ` where h.city_id = ` + strconv.Itoa(city_id)
	} else if harborsId != "" {
		query = query + ` where e.harbors_id = '` + harborsId + `'`
	} else {
		//return nil, models.ErrBadParamInput
	}
	if guest != "" {
		guests, _ := strconv.Atoi(guest)
		query = query + ` AND e.exp_max_guest =` + strconv.Itoa(guests)
	}
	if trip != "" {
		trips, _ := strconv.Atoi(trip)
		var tripType string
		if trips == 0 {
			tripType = "Private Trip"
		} else if trips == 1 {
			tripType = "Share Trip"
		} else {
			return nil, models.ErrInternalServerError
		}
		query = query + ` AND e.exp_trip_type = '` + tripType + `'`
	}

	if len(activityTypeArray) != 0 {
		//types, _ := strconv.Atoi(activityType)
		for index, id := range activityTypeArray {
			if index == 0 && index != (len(activityTypeArray)-1) {
				query = query + ` AND (fat.id =` + strconv.Itoa(id)
			} else if index == 0 && index == (len(activityTypeArray)-1) {
				query = query + ` AND (fat.id =` + strconv.Itoa(id) + ` ) `
			} else if index == (len(activityTypeArray) - 1) {
				query = query + ` OR fat.id =` + strconv.Itoa(id) + ` ) `
			} else {
				query = query + ` OR fat.id =` + strconv.Itoa(id)
			}
		}

	}
	if bottomPrice != "" && upPrice != "" {
		bottomprices, _ := strconv.ParseFloat(bottomPrice, 64)
		upprices, _ := strconv.ParseFloat(upPrice, 64)

		query = query + ` AND (ep.price between ` + fmt.Sprint(bottomprices) + ` AND ` + fmt.Sprint(upprices) + `)`
		if sortBy != "" {
			if sortBy == "priceup" {
				query = query + ` ORDER BY ep.price DESC`
			} else if sortBy == "pricedown" {
				query = query + ` ORDER BY ep.price ASC`
			}
		}
	}
	if sortBy != "" {
		if sortBy == "ratingup" {
			query = query + ` ORDER BY e.rating DESC`
		} else if sortBy == "ratingdown" {
			query = query + ` ORDER BY e.rating ASC`
		}
	}

	fmt.Println(query)

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
		var expType []string
		if errUnmarshal := json.Unmarshal([]byte(exp.ExpType), &expType); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}

		expPayment, err := m.paymentRepo.GetByExpID(ctx, exp.Id)
		if err != nil {
			return nil, err
		}

		var currency string
		if expPayment[0].Currency == 1 {
			currency = "USD"
		} else {
			currency = "IDR"
		}

		var priceItemType string
		if expPayment[0].PriceItemType == 1 {
			priceItemType = "Per Pax"
		} else {
			priceItemType = "Per Trip"
		}

		countRating, err := m.reviewsRepo.CountRating(ctx, exp.Id)
		if err != nil {
			return nil, err
		}
		coverPhoto := models.CoverPhotosObj{
			Original:  exp.CoverPhoto,
			Thumbnail: "",
		}
		var listPhotos []models.ExpPhotosObj
		expPhotoQuery, errorQuery := m.expPhotos.GetByExperienceID(ctx, exp.Id)
		if errorQuery != nil {
			return nil, errorQuery
		}
		if expPhotoQuery != nil {
			for _, element := range expPhotoQuery {
				expPhoto := models.ExpPhotosObj{
					Folder:        element.ExpPhotoFolder,
					ExpPhotoImage: nil,
				}
				var expPhotoImage []models.CoverPhotosObj
				errObject := json.Unmarshal([]byte(element.ExpPhotoImage), &expPhotoImage)
				if errObject != nil {
					//fmt.Println("Error : ",err.Error())
					return nil, models.ErrInternalServerError
				}
				expPhoto.ExpPhotoImage = expPhotoImage
				listPhotos = append(listPhotos, expPhoto)
			}
		}
		results[i] = &models.ExpSearchObject{
			Id:          exp.Id,
			ExpTitle:    exp.ExpTitle,
			ExpType:     expType,
			Rating:      exp.Rating,
			CountRating: countRating,
			Currency:    currency,
			Price:       expPayment[0].Price,
			PaymentType: priceItemType,
			Longitude:   exp.Longitude,
			Latitude:    exp.Latitude,
			CoverPhoto:  coverPhoto,
			ListPhoto:   listPhotos,
		}
	}

	return results, nil

}
func (m experienceUsecase) GetExpInspirations(ctx context.Context) ([]*models.ExpInspirationObject, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	results, err := m.inspirationRepo.GetExpInspirations(ctx)
	if err != nil {
		return nil, err
	}

	return results, nil
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
		var expType []string
		if errUnmarshal := json.Unmarshal([]byte(exp.ExpType), &expType); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
		expPayment, err := m.paymentRepo.GetByExpID(ctx, exp.Id)
		if err != nil {
			return nil, err
		}

		var currency string
		if expPayment[0].Currency == 1 {
			currency = "USD"
		} else {
			currency = "IDR"
		}

		var priceItemType string
		if expPayment[0].PriceItemType == 1 {
			priceItemType = "Per Pax"
		} else {
			priceItemType = "Per Trip"
		}

		countRating, err := m.reviewsRepo.CountRating(ctx, exp.Id)
		if err != nil {
			return nil, err
		}

		var listPhotos []models.ExpPhotosObj
		var coverPhotos models.CoverPhotosObj
		if exp.CoverPhoto != "" {
			coverPhotos.Original = exp.CoverPhoto
			coverPhotos.Thumbnail = ""
		}
		expPhotoQuery, errorQuery := m.expPhotos.GetByExperienceID(ctx, exp.Id)
		if errorQuery != nil {
			return nil, errorQuery
		}
		if expPhotoQuery != nil {
			for _, element := range expPhotoQuery {
				expPhoto := models.ExpPhotosObj{
					Folder:        element.ExpPhotoFolder,
					ExpPhotoImage: nil,
				}
				var expPhotoImage []models.CoverPhotosObj
				errObject := json.Unmarshal([]byte(element.ExpPhotoImage), &expPhotoImage)
				if errObject != nil {
					//fmt.Println("Error : ",err.Error())
					return nil, models.ErrInternalServerError
				}
				expPhoto.ExpPhotoImage = expPhotoImage
				listPhotos = append(listPhotos, expPhoto)
			}
		}
		results[i] = &models.ExpSearchObject{
			Id:          exp.Id,
			ExpTitle:    exp.ExpTitle,
			ExpType:     expType,
			Rating:      exp.Rating,
			CountRating: countRating,
			Currency:    currency,
			Price:       expPayment[0].Price,
			PaymentType: priceItemType,
			CoverPhoto:  coverPhotos,
			ListPhoto:   listPhotos,
		}
	}

	return results, nil
}

func (m experienceUsecase) CreateExperience(c context.Context, commandExperience models.NewCommandExperience, token string) (*models.ResponseCreateExperience, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	currentUserMerchant, err := m.mUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	//if commandExperience.ExpType != ""
	expItinerary, _ := json.Marshal(commandExperience.ExpInternary)
	expFacilities, _ := json.Marshal(commandExperience.ExpFacilities)
	expInclusion, _ := json.Marshal(commandExperience.ExpInclusion)
	expRules, _ := json.Marshal(commandExperience.ExpRules)
	expTypes, _ := json.Marshal(commandExperience.ExpType)
	experiences := models.Experience{
		Id:                      "",
		CreatedBy:               currentUserMerchant.MerchantEmail,
		CreatedDate:             time.Time{},
		ModifiedBy:              nil,
		ModifiedDate:            nil,
		DeletedBy:               nil,
		DeletedDate:             nil,
		IsDeleted:               0,
		IsActive:                0,
		ExpTitle:                commandExperience.ExpTitle,
		ExpType:                 string(expTypes),
		ExpTripType:             commandExperience.ExpTripType,
		ExpBookingType:          commandExperience.ExpBookingType,
		ExpDesc:                 commandExperience.ExpDesc,
		ExpMaxGuest:             commandExperience.ExpMaxGuest,
		ExpPickupPlace:          commandExperience.ExpPickupPlace,
		ExpPickupTime:           commandExperience.ExpPickupTime,
		ExpPickupPlaceLongitude: commandExperience.ExpPickupPlaceLongitude,
		ExpPickupPlaceLatitude:  commandExperience.ExpPickupPlaceLatitude,
		ExpPickupPlaceMapsName:  commandExperience.ExpPickupPlaceMapsName,
		ExpInternary:            string(expItinerary),
		ExpFacilities:           string(expFacilities),
		ExpInclusion:            string(expInclusion),
		ExpRules:                string(expRules),
		Status:                  commandExperience.Status,
		Rating:                  0,
		ExpLocationLatitude:     commandExperience.ExpLocationLatitude,
		ExpLocationLongitude:    commandExperience.ExpLocationLongitude,
		ExpLocationName:         commandExperience.ExpLocationName,
		ExpCoverPhoto:           commandExperience.ExpCoverPhoto,
		ExpDuration:             commandExperience.ExpDuration,
		MinimumBookingId:        commandExperience.MinimumBookingId,
		MerchantId:              currentUserMerchant.Id,
		HarborsId:               commandExperience.HarborsId,
	}
	insertToExperience, err := m.experienceRepo.Insert(ctx, &experiences)
	for _, element := range commandExperience.ExpPhotos {
		images, _ := json.Marshal(element.ExpPhotoImage)
		expPhoto := models.ExpPhotos{
			Id:             "",
			CreatedBy:      currentUserMerchant.MerchantEmail,
			CreatedDate:    time.Time{},
			ModifiedBy:     nil,
			ModifiedDate:   nil,
			DeletedBy:      nil,
			DeletedDate:    nil,
			IsDeleted:      0,
			IsActive:       0,
			ExpPhotoFolder: element.Folder,
			ExpPhotoImage:  string(images),
			ExpId:          *insertToExperience,
		}

		_, err = m.expPhotos.Insert(ctx, &expPhoto)
	}

	for _, element := range commandExperience.ExpPayment {
		var priceItemType int
		if element.PriceItemType == "Per Pax" {
			priceItemType = 1
		} else {
			priceItemType = 0
		}
		var currency int
		if element.Currency == "USD" {
			currency = 1
		} else {
			currency = 0
		}
		payments := models.ExperiencePayment{
			Id:               "",
			CreatedBy:        currentUserMerchant.MerchantEmail,
			CreatedDate:      time.Time{},
			ModifiedBy:       nil,
			ModifiedDate:     nil,
			DeletedBy:        nil,
			DeletedDate:      nil,
			IsDeleted:        0,
			IsActive:         0,
			ExpPaymentTypeId: element.PaymentTypeId,
			ExpId:            *insertToExperience,
			PriceItemType:    priceItemType,
			Currency:         currency,
			Price:            element.Price,
			CustomPrice:      nil,
		}

		err = m.paymentRepo.Insert(ctx, payments)
	}

	for _, element := range commandExperience.ExpAvailability {
		date, _ := json.Marshal(element.Date)
		expAvailability := models.ExpAvailability{
			Id:                   "",
			CreatedBy:            currentUserMerchant.MerchantEmail,
			CreatedDate:          time.Time{},
			ModifiedBy:           nil,
			ModifiedDate:         nil,
			DeletedBy:            nil,
			DeletedDate:          nil,
			IsDeleted:            0,
			IsActive:             0,
			ExpAvailabilityMonth: element.Month,
			ExpAvailabilityDate:  string(date),
			ExpAvailabilityYear:  element.Year,
			ExpId:                *insertToExperience,
		}

		err = m.exp_availablitiy.Insert(ctx, expAvailability)
	}

	for _, element := range commandExperience.ExperienceAddOn {
		var currency int
		if element.Currency == "USD" {
			currency = 1
		} else {
			currency = 0
		}
		addOns := models.ExperienceAddOn{
			Id:           "",
			CreatedBy:    currentUserMerchant.MerchantEmail,
			CreatedDate:  time.Time{},
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     0,
			Name:         element.Name,
			Desc:         element.Desc,
			Currency:     currency,
			Amount:       element.Amount,
			ExpId:        *insertToExperience,
		}
		err = m.adOnsRepo.Insert(ctx, addOns)
	}

	response := models.ResponseCreateExperience{
		Id:      *insertToExperience,
		Message: "Success Publish",
	}
	return &response, nil

}
func (m experienceUsecase) GetByID(c context.Context, id string) (*models.ExperienceDto, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()

	res, err := m.experienceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	var expPhotos []models.ExpPhotosObj
	expPhotoQuery, errorQuery := m.expPhotos.GetByExperienceID(ctx, res.Id)
	if expPhotoQuery != nil {
		for _, element := range expPhotoQuery {
			expPhoto := models.ExpPhotosObj{
				Folder:        element.ExpPhotoFolder,
				ExpPhotoImage: nil,
			}
			var expPhotoImage []models.CoverPhotosObj
			errObject := json.Unmarshal([]byte(element.ExpPhotoImage), &expPhotoImage)
			if errObject != nil {
				//fmt.Println("Error : ",err.Error())
				return nil, models.ErrInternalServerError
			}
			expPhoto.ExpPhotoImage = expPhotoImage
			expPhotos = append(expPhotos, expPhoto)
		}
	}
	var expPayment []models.ExpPaymentObj
	expPaymentQuery, errorQuery := m.paymentRepo.GetByExpID(ctx, res.Id)
	for _, elementPayment := range expPaymentQuery {
		var currency string
		if elementPayment.Currency == 1 {
			currency = "USD"
		} else {
			currency = "IDR"
		}

		var priceItemType string
		if elementPayment.PriceItemType == 1 {
			priceItemType = "Per Pax"
		} else {
			priceItemType = "Per Trip"
		}
		expPayobj := models.ExpPaymentObj{
			Id:              elementPayment.Id,
			Currency:        currency,
			Price:           elementPayment.Price,
			PriceItemType:   priceItemType,
			PaymentTypeId:   elementPayment.ExpPaymentTypeId,
			PaymentTypeName: elementPayment.ExpPaymentTypeName,
			PaymentTypeDesc: elementPayment.ExpPaymentTypeDesc,
		}
		expPayment = append(expPayment, expPayobj)
	}

	var expAvailability []models.ExpAvailablitityObj
	expAvailabilityQuery, errorQuery := m.exp_availablitiy.GetByExpId(ctx, res.Id)
	if errorQuery != nil {
		return nil, errorQuery
	}
	if expAvailabilityQuery != nil {
		for _, element := range expAvailabilityQuery {
			expA := models.ExpAvailablitityObj{
				Year:  element.ExpAvailabilityYear,
				Month: element.ExpAvailabilityMonth,
				Date:  nil,
			}
			var date []string
			errObject := json.Unmarshal([]byte(element.ExpAvailabilityDate), &date)
			if errObject != nil {
				//fmt.Println("Error : ",err.Error())
				return nil, models.ErrInternalServerError
			}
			expA.Date = date
			expAvailability = append(expAvailability, expA)
		}
	}

	var expType []string
	errObject := json.Unmarshal([]byte(res.ExpType), &expType)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil, models.ErrInternalServerError
	}
	expItinerary := models.ExpItineraryObject{}
	errObject = json.Unmarshal([]byte(res.ExpInternary), &expItinerary)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil, models.ErrInternalServerError
	}
	var expFacilities []models.ExpFacilitiesObject
	errObject = json.Unmarshal([]byte(res.ExpFacilities), &expFacilities)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil, models.ErrInternalServerError
	}
	var expInclusion []models.ExpInclusionObject
	errObject = json.Unmarshal([]byte(res.ExpInclusion), &expInclusion)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil, models.ErrInternalServerError
	}
	var expRules []models.ExpRulesObject
	errObject = json.Unmarshal([]byte(res.ExpRules), &expRules)
	if errObject != nil {
		//fmt.Println("Error : ",err.Error())
		return nil, models.ErrInternalServerError
	}
	harbors, err := m.harborsRepo.GetByID(ctx, res.HarborsId)
	city, err := m.cpcRepo.GetCityByID(ctx, harbors.CityId)
	province, err := m.cpcRepo.GetProvinceByID(ctx, city.ProvinceId)
	minimumBooking := models.MinimumBookingObj{
		MinimumBookingDesc:   res.MinimumBookingDesc,
		MinimumBookingAmount: res.MinimumBookingAmount,
	}
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
		ExpAvailability:         expAvailability,
		ExpPayment:              expPayment,
		ExpPhotos:               expPhotos,
		Status:                  res.Status,
		Rating:                  res.Rating,
		ExpLocationLatitude:     res.ExpLocationLatitude,
		ExpLocationLongitude:    res.ExpLocationLongitude,
		ExpLocationName:         res.ExpLocationName,
		ExpCoverPhoto:           res.ExpCoverPhoto,
		ExpDuration:             res.ExpDuration,
		MinimumBooking:          minimumBooking,
		MerchantId:              res.MerchantId,
		HarborsName:             harbors.HarborsName,
		City:                    city.CityName,
		Province:                province.ProvinceName,
	}
	return &experiences, nil
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */
