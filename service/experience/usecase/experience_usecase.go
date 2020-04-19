package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/service/filter_activity_type"
	"math"
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
	filterATRepo     filter_activity_type.Repository
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
	fac filter_activity_type.Repository,
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
		filterATRepo:fac,
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

func (m experienceUsecase) GetPublishedExpCount(ctx context.Context, token string) (*models.Count, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	currentMerchant, err := m.mUsecase.ValidateTokenMerchant(ctx, token)
	if err != nil {
		return nil, err
	}

	count, err := m.experienceRepo.GetPublishedExpCount(ctx, currentMerchant.Id)
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

		countRating, err := m.reviewsRepo.CountRating(ctx, 0, exp.Id)
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
		countRating, err := m.reviewsRepo.CountRating(ctx, 0, element.Id)
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
				Id:          element.Id,
				ExpTitle:    element.ExpTitle,
				ExpType:     expType,
				Rating:      element.Rating,
				CountRating: countRating,
				Currency:    currency,
				//Price:        expPayment[0].Price,
				Payment_type: priceItemType,
				Cover_Photo:  coverPhotos,
			}
			if expPayment != nil {
				expDto.Price = expPayment[0].Price
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
					Id:          element.Id,
					ExpTitle:    element.ExpTitle,
					ExpType:     expType,
					Rating:      element.Rating,
					CountRating: countRating,
					Currency:    currency,
					//Price:        expPayment[0].Price,
					Payment_type: priceItemType,
					Cover_Photo:  coverPhotos,
				}
				if expPayment != nil {
					expDto.Price = expPayment[0].Price
				}
				cityDto.Item = append(cityDto.Item, expDto)
				expListDto = append(expListDto, &cityDto)
			} else {
				for _, dto := range expListDto {
					if dto.CityId == element.CityId {
						expDto := models.ExperienceUserDiscoverPreferenceDto{
							Id:          element.Id,
							ExpTitle:    element.ExpTitle,
							ExpType:     expType,
							Rating:      element.Rating,
							CountRating: countRating,
							Currency:    currency,
							//Price:        expPayment[0].Price,
							Payment_type: priceItemType,
							Cover_Photo:  coverPhotos,
						}
						if expPayment != nil {
							expDto.Price = expPayment[0].Price
						}
						dto.Item = append(dto.Item, expDto)
					}
				}
			}
		}
	}
	return expListDto, nil
}

func (m experienceUsecase) FilterSearchExp(
	ctx context.Context,
	isMerchant bool,
	search,
	token,
	qStatus,
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
	page int,
	limit int,
	offset int,
) (*models.FilterSearchWithPagination, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	var activityTypeArray []int
	if activityType != "" && activityType != "[]" {
		if errUnmarshal := json.Unmarshal([]byte(activityType), &activityTypeArray); errUnmarshal != nil {
			return nil, models.ErrInternalServerError
		}
	}
	query := `
	select
		e.id,
		e.exp_title,
		e.exp_type,
		e.status as exp_status,
		e.rating,
		e.exp_location_latitude as latitude,
		e.exp_location_longitude as longitude,
		e.exp_cover_photo as cover_photo,
		province_name AS province
	from 
		experiences e
	JOIN harbors ha ON e.harbors_id = ha.id
	JOIN cities ci ON ha.city_id = ci.id
	JOIN provinces p ON ci.province_id = p.id`

	qCount := `
	select COUNT(*) from experiences e
	JOIN harbors ha ON e.harbors_id = ha.id
	JOIN cities ci ON ha.city_id = ci.id
	JOIN provinces p ON ci.province_id = p.id`

	if bottomPrice != "" && upPrice != "" && qStatus != "draft"{
		query = query + ` join experience_payments ep on ep.exp_id = e.id`
		qCount = qCount + ` join experience_payments ep on ep.exp_id = e.id`
	}
	if startDate != "" && endDate != "" && qStatus != "draft"{
		query = query + ` join exp_availabilities ead on ead.exp_id = e.id`
	}
	//if len(activityTypeArray) != 0 {
	//	query = query + ` join filter_activity_types fat on fat.exp_id = e.id`
	//	qCount = qCount + ` join filter_activity_types fat on fat.exp_id = e.id`
	//}
	if cityID != "" && qStatus != "draft"{
		query = query + ` join harbors h on h.id = e.harbors_id`
		qCount = qCount + ` join harbors h on h.id = e.harbors_id`
	}

	query = query + ` WHERE e.is_deleted = 0 AND e.is_active = 1`
	qCount = qCount + ` WHERE e.is_deleted = 0 AND e.is_active = 1`

	if isMerchant {
		if token == "" {
			return nil, models.ErrUnAuthorize
		}

		currentMerchant, err := m.mUsecase.ValidateTokenMerchant(ctx, token)
		if err != nil {
			return nil, err
		}

		query = query + ` AND e.merchant_id = '` + currentMerchant.Id + `'`
		qCount = qCount + ` AND e.merchant_id = '` + currentMerchant.Id + `'`
	}

	if search != "" {
		keyword := `'%` + search + `%'`
		query = query + ` AND LOWER(e.exp_title) LIKE LOWER(` + keyword + `)`
		qCount = qCount + ` AND LOWER(e.exp_title) LIKE LOWER(` + keyword + `)`
	}
	if cityID != "" {
		city_id, _ := strconv.Atoi(cityID)
		query = query + ` AND h.city_id = ` + strconv.Itoa(city_id)
		qCount = qCount + ` AND h.city_id = ` + strconv.Itoa(city_id)
	} else if harborsId != "" {
		query = query + ` AND e.harbors_id = '` + harborsId + `'`
		qCount = qCount + ` AND e.harbors_id = '` + harborsId + `'`
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
			query = query + ` AND e.status IN (2,3,4)`
			qCount = qCount + ` AND e.status IN (2,3,4)`
		}else {
			query = query + ` AND e.status =` + strconv.Itoa(status)
			qCount = qCount + ` AND e.status =` + strconv.Itoa(status)
		}

	}
	if guest != "" {
		guests, _ := strconv.Atoi(guest)
		query = query + ` AND e.exp_max_guest <=` + strconv.Itoa(guests)
		qCount = qCount + ` AND e.exp_max_guest <=` + strconv.Itoa(guests)
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
		qCount = qCount + ` AND e.exp_trip_type = '` + tripType + `'`
	}

	if len(activityTypeArray) != 0 {
		for index, id := range activityTypeArray {
			if index == 0 && index != (len(activityTypeArray)-1) {
				query = query + ` AND (e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )`
				qCount = qCount + ` AND (e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )`
			} else if index == 0 && index == (len(activityTypeArray)-1) {
				query = query + ` AND (e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )` + ` ) `
				qCount = qCount + ` AND (e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )` + ` ) `
			} else if index == (len(activityTypeArray) - 1) {
				query = query + ` OR e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )` + ` )`
				qCount = qCount + ` OR e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )` + ` )`
			} else {
				query = query + ` OR e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )`
				qCount = qCount + ` OR e.id = (SELECT distinct exp_id FROM filter_activity_types where exp_id = e.id and exp_type_id = ` + strconv.Itoa(id) + ` )`
			}
		}

	}
	if bottomPrice != "" && upPrice != "" && qStatus != "draft"{
		bottomprices, _ := strconv.ParseFloat(bottomPrice, 64)
		upprices, _ := strconv.ParseFloat(upPrice, 64)

		query = query + ` AND (ep.price between ` + fmt.Sprint(bottomprices) + ` AND ` + fmt.Sprint(upprices) + `)`
		qCount = qCount + ` AND (ep.price between ` + fmt.Sprint(bottomprices) + ` AND ` + fmt.Sprint(upprices) + `)`
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

	if startDate != "" && endDate != "" {
		var startDates []string

		layoutFormat := "2006-01-02"
		start, errDateDob := time.Parse(layoutFormat, startDate)
		if errDateDob != nil {
			return nil, errDateDob
		}
		end, errDateDob := time.Parse(layoutFormat, endDate)
		if errDateDob != nil {
			return nil, errDateDob
		}
		startDates = append(startDates, start.Format("2006-01-02"))
	datess:

		start = start.AddDate(0, 0, 1)
		startDates = append(startDates, start.Format("2006-01-02"))
		if start == end {

		} else {
			startDates = append(startDates, start.String())
			goto datess
		}

		for index, id := range startDates {
			if index == 0 && index != (len(startDates)-1) {
				query = query + ` AND (ead.exp_availability_date like '%` + id + `%' `
				qCount = qCount + ` AND (ead.exp_availability_date like '%` + id + `%' `
			} else if index == 0 && index == (len(startDates)-1) {
				query = query + ` AND (ead.exp_availability_date like '%` + id + `%' ) `
				qCount = qCount + ` AND (ead.exp_availability_date like '%` + id + `%' ) `
			} else if index == (len(startDates) - 1) {
				query = query + ` OR ead.exp_availability_date like '%` + id + `%' ) `
				qCount = qCount + ` OR ead.exp_availability_date like '%` + id + `%' ) `
			} else {
				query = query + ` OR ead.exp_availability_date like '%` + id + `%' `
				qCount = qCount + ` OR ead.exp_availability_date like '%` + id + `%' `
			}
		}
	}
	fmt.Println(query)
	expList, err := m.experienceRepo.QueryFilterSearch(ctx, query, limit, offset)
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
		var price float64
		var priceItemType string

		if expPayment != nil {

			price = expPayment[0].Price
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
		}

		countRating, err := m.reviewsRepo.CountRating(ctx, 0, exp.Id)
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
					return nil, models.ErrInternalServerError
				}
				expPhoto.ExpPhotoImage = expPhotoImage
				listPhotos = append(listPhotos, expPhoto)
			}
		}

		var transStatus string
		if exp.ExpStatus == 0 {
			transStatus = "Preview"
		} else if exp.ExpStatus == 1 {
			transStatus = "Draft"
		} else if exp.ExpStatus == 2 {
			transStatus = "Published"
		} else if exp.ExpStatus == 3 {
			transStatus = "Unpublished"
		} else if exp.ExpStatus == 4 {
			transStatus = "Archived"
		}

		results[i] = &models.ExpSearchObject{
			Id:          exp.Id,
			ExpTitle:    exp.ExpTitle,
			ExpType:     expType,
			ExpStatus:   transStatus,
			Rating:      exp.Rating,
			CountRating: countRating,
			Currency:    currency,
			Price:       price,
			PaymentType: priceItemType,
			Longitude:   exp.Longitude,
			Latitude:    exp.Latitude,
			Province:    exp.Province,
			CoverPhoto:  coverPhoto,
			ListPhoto:   listPhotos,
		}
	}
	totalRecords, _ := m.experienceRepo.CountFilterSearch(ctx, qCount)
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
		RecordPerPage: len(results),
	}

	response := &models.FilterSearchWithPagination{
		Data: results,
		Meta: meta,
	}

	return response, nil

}
func (m experienceUsecase) GetExpInspirations(ctx context.Context) ([]*models.ExpInspirationDto, error) {
	ctx, cancel := context.WithTimeout(ctx, m.contextTimeout)
	defer cancel()

	query, err := m.inspirationRepo.GetExpInspirations(ctx)
	var results []*models.ExpInspirationDto
	for _, element := range query {
		getCountReview, err := m.reviewsRepo.CountRating(ctx, 0, element.ExpId)
		if err != nil {
			return nil, err
		}
		var expType []string
		if element.ExpType != "" {
			if errUnmarshal := json.Unmarshal([]byte(element.ExpType), &expType); errUnmarshal != nil {
				return nil, models.ErrInternalServerError
			}
		}
		dto := models.ExpInspirationDto{
			ExpInspirationID: element.ExpInspirationID,
			ExpId:            element.ExpId,
			ExpTitle:         element.ExpTitle,
			ExpDesc:          element.ExpDesc,
			ExpCoverPhoto:    element.ExpCoverPhoto,
			ExpType:          expType,
			Rating:           element.Rating,
			CountRating:      getCountReview,
		}
		results = append(results, &dto)
	}
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

		countRating, err := m.reviewsRepo.CountRating(ctx, 0, exp.Id)
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
		MinimumBookingId:        &commandExperience.MinimumBookingId,
		MerchantId:              currentUserMerchant.Id,
		HarborsId:               &commandExperience.HarborsId,
	}
	if *experiences.HarborsId == "" && experiences.Status == 1 {
		experiences.HarborsId = nil
	}
	if *experiences.MinimumBookingId == "" && experiences.Status == 1 {
		experiences.MinimumBookingId = nil
	}
	insertToExperience, err := m.experienceRepo.Insert(ctx, &experiences)

	for _,element := range commandExperience.ExpType{
		getExpType ,err := m.typesRepo.GetByName(ctx,element)
		if err != nil {
			return nil ,err
		}
		filterActivityT := models.FilterActivityType{
			Id:           0,
			CreatedBy:    currentUserMerchant.MerchantEmail,
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExpTypeId:    getExpType.ExpTypeID,
			ExpId:        insertToExperience,
		}
		insertToFilterAT := m.filterATRepo.Insert(ctx,&filterActivityT)
		if insertToFilterAT != nil {
			return nil,insertToFilterAT
		}
	}

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

		id, err := m.expPhotos.Insert(ctx, &expPhoto)
		if err != nil {
			return nil, err
		}
		element.Id = *id
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

		id, err := m.paymentRepo.Insert(ctx, payments)
		if err != nil {
			return nil, err
		}
		element.Id = id
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

		id, err := m.exp_availablitiy.Insert(ctx, expAvailability)
		if err != nil {
			return nil, err
		}
		element.Id = id
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
		id, err := m.adOnsRepo.Insert(ctx, addOns)
		if err != nil {
			return nil, err
		}
		element.Id = id
	}

	var status string
	if commandExperience.Status == 1 {
		status = "Draft"
	} else if commandExperience.Status == 2 {
		status = "Publish"
	}
	response := models.ResponseCreateExperience{
		Id:      *insertToExperience,
		Message: "Success " + status,
	}
	return &response, nil

}
func (m experienceUsecase) UpdateExperience(c context.Context, commandExperience models.NewCommandExperience, token string) (*models.ResponseCreateExperience, error) {
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
		Id:                      commandExperience.Id,
		CreatedBy:               currentUserMerchant.MerchantEmail,
		CreatedDate:             time.Time{},
		ModifiedBy:              &currentUserMerchant.MerchantEmail,
		ModifiedDate:            &time.Time{},
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
		MinimumBookingId:        &commandExperience.MinimumBookingId,
		MerchantId:              currentUserMerchant.Id,
		HarborsId:               &commandExperience.HarborsId,
	}
	if *experiences.HarborsId == "" && experiences.Status == 1 {
		experiences.HarborsId = nil
	}
	if *experiences.MinimumBookingId == "" && experiences.Status == 1 {
		experiences.MinimumBookingId = nil
	}
	insertToExperience, err := m.experienceRepo.Update(ctx, &experiences)

	var photoIds []string
	for _, element := range commandExperience.ExpPhotos {
		if element.Id == "" {

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

			id, err := m.expPhotos.Insert(ctx, &expPhoto)
			if err != nil {
				return nil, err
			}
			photoIds = append(photoIds, *id)
			element.Id = *id
		} else {
			images, _ := json.Marshal(element.ExpPhotoImage)
			expPhoto := models.ExpPhotos{
				Id:             element.Id,
				CreatedBy:      "",
				CreatedDate:    time.Time{},
				ModifiedBy:     &currentUserMerchant.MerchantEmail,
				ModifiedDate:   &time.Time{},
				DeletedBy:      nil,
				DeletedDate:    nil,
				IsDeleted:      0,
				IsActive:       0,
				ExpPhotoFolder: element.Folder,
				ExpPhotoImage:  string(images),
				ExpId:          *insertToExperience,
			}

			_, err = m.expPhotos.Update(ctx, &expPhoto)
			photoIds = append(photoIds, element.Id)
		}
	}

	_ = m.expPhotos.Deletes(ctx, photoIds, *insertToExperience, currentUserMerchant.MerchantEmail)

	var expPaymentIds []string
	for _, element := range commandExperience.ExpPayment {
		if element.Id == "" {
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

			id, err := m.paymentRepo.Insert(ctx, payments)
			if err != nil {
				return nil, err
			}
			expPaymentIds = append(expPaymentIds, id)
			element.Id = id

		} else {
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
				Id:               element.Id,
				CreatedBy:        "",
				CreatedDate:      time.Time{},
				ModifiedBy:       &currentUserMerchant.MerchantEmail,
				ModifiedDate:     &time.Time{},
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

			err = m.paymentRepo.Update(ctx, payments)

			expPaymentIds = append(expPaymentIds, element.Id)
		}

	}

	_ = m.paymentRepo.Deletes(ctx, expPaymentIds, *insertToExperience, currentUserMerchant.MerchantEmail)

	var expAvailabilityIds []string
	for _, element := range commandExperience.ExpAvailability {
		if element.Id == "" {
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

			id, err := m.exp_availablitiy.Insert(ctx, expAvailability)
			if err != nil {
				return nil, err
			}
			expAvailabilityIds = append(expAvailabilityIds, id)
			element.Id = id
		} else {
			date, _ := json.Marshal(element.Date)
			expAvailability := models.ExpAvailability{
				Id:                   element.Id,
				CreatedBy:            "",
				CreatedDate:          time.Time{},
				ModifiedBy:           &currentUserMerchant.MerchantEmail,
				ModifiedDate:         &time.Time{},
				DeletedBy:            nil,
				DeletedDate:          nil,
				IsDeleted:            0,
				IsActive:             0,
				ExpAvailabilityMonth: element.Month,
				ExpAvailabilityDate:  string(date),
				ExpAvailabilityYear:  element.Year,
				ExpId:                *insertToExperience,
			}

			err = m.exp_availablitiy.Update(ctx, expAvailability)
			expAvailabilityIds = append(expAvailabilityIds, element.Id)
		}

	}

	_ = m.exp_availablitiy.Deletes(ctx, expAvailabilityIds, *insertToExperience, currentUserMerchant.MerchantEmail)

	var addOnIds []string
	for _, element := range commandExperience.ExperienceAddOn {
		if element.Id == "" {
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
			id, err := m.adOnsRepo.Insert(ctx, addOns)
			if err != nil {
				return nil, err
			}
			addOnIds = append(addOnIds, id)
			element.Id = id
		} else {
			var currency int
			if element.Currency == "USD" {
				currency = 1
			} else {
				currency = 0
			}
			addOns := models.ExperienceAddOn{
				Id:           element.Id,
				CreatedBy:    "",
				CreatedDate:  time.Time{},
				ModifiedBy:   &currentUserMerchant.MerchantEmail,
				ModifiedDate: &time.Time{},
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
			err = m.adOnsRepo.Update(ctx, addOns)
			addOnIds = append(addOnIds, element.Id)
		}

	}
	_ = m.adOnsRepo.Deletes(ctx, addOnIds, *insertToExperience, currentUserMerchant.MerchantEmail)

	var status string
	if commandExperience.Status == 1 {
		status = "Draft"
	} else if commandExperience.Status == 2 {
		status = "Publish"
	}
	response := models.ResponseCreateExperience{
		Id:      *insertToExperience,
		Message: "Success " + status,
	}
	return &response, nil

}
func (m experienceUsecase) PublishExperience(c context.Context, commandExperience models.NewCommandExperience, token string) (*models.ResponseCreateExperience, error) {
	ctx, cancel := context.WithTimeout(c, m.contextTimeout)
	defer cancel()
	var response *models.ResponseCreateExperience
	if commandExperience.Id == "" {
		create, err := m.CreateExperience(ctx, commandExperience, token)
		if err != nil {
			return nil, err
		}
		response = create
	} else {
		update, err := m.UpdateExperience(ctx, commandExperience, token)
		if err != nil {
			return nil, err
		}
		response = update
	}
	return response, nil
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
				Id:            element.Id,
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
				Id:    element.Id,
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
	countRating, err := m.reviewsRepo.CountRating(ctx, 0, res.Id)
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
		CountRating:             countRating,
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
