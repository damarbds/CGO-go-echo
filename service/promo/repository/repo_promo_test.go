package repository_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/models"
	PromoRepo "github.com/service/promo/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var(
	date = time.Now().String()
	isAnyTripPeriod = 1
	desc = "ini description"
	disc float32 = 12.2
	mockPromo = []models.Promo{
		models.Promo{
			Id:                 "asdqeqrqrasdsad",
			CreatedBy:          "test",
			CreatedDate:        time.Now(),
			ModifiedBy:         nil,
			ModifiedDate:       nil,
			DeletedBy:          nil,
			DeletedDate:        nil,
			IsDeleted:          0,
			IsActive:           1,
			PromoCode:          "Test1",
			PromoName:          "Test 1",
			PromoDesc:          "Test 1",
			PromoValue:         1,
			PromoType:          1,
			PromoImage:         "asdasdasdas",
			StartDate:          &date,
			EndDate:            &date,
			StartTripPeriod:    &date,
			EndTripPeriod:      &date,
			IsAnyTripPeriod:    &isAnyTripPeriod,
			HowToGet:           &desc,
			HowToUse:           &desc,
			TermCondition:      &desc,
			Disclaimer:         &desc,
			MaxDiscount:        &disc,
			MaxUsage:           &isAnyTripPeriod,
			ProductionCapacity: &isAnyTripPeriod,
			CurrencyId:         &isAnyTripPeriod,
			PromoProductType:   &isAnyTripPeriod,
		},
		models.Promo{
			Id:                 "tytetrtrete",
			CreatedBy:          "test",
			CreatedDate:        time.Now(),
			ModifiedBy:         nil,
			ModifiedDate:       nil,
			DeletedBy:          nil,
			DeletedDate:        nil,
			IsDeleted:          0,
			IsActive:           1,
			PromoCode:          "Test2",
			PromoName:          "Test 2",
			PromoDesc:          "Test 2",
			PromoValue:         1,
			PromoType:          1,
			PromoImage:         "asdasdasdas",
			StartDate:          &date,
			EndDate:            &date,
			StartTripPeriod:    &date,
			EndTripPeriod:      &date,
			IsAnyTripPeriod:    &isAnyTripPeriod,
			HowToGet:           &desc,
			HowToUse:           &desc,
			TermCondition:      &desc,
			Disclaimer:         &desc,
			MaxDiscount:        &disc,
			MaxUsage:           &isAnyTripPeriod,
			ProductionCapacity: &isAnyTripPeriod,
			CurrencyId:         &isAnyTripPeriod,
			PromoProductType:   &isAnyTripPeriod,
		},
	}

	merchantIds = []string{"adad","qeqe","adasd","jkjffsf","qpeiuqioeu","kjsdfdsj;f"}
)
func TestCount(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(len(mockPromo))

	query := `SELECT count\(\*\) AS count FROM promos WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	res, err := a.GetCount(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, res, 2, "")
}
func TestCountFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow("test")

	query := `SELECT count\(\*\) AS count FROM promos WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	_, err = a.GetCount(context.TODO())
	assert.Error(t, err)
}
func TestFetchWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"}).
		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
			mockPromo[0].IsActive, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
		 	mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod).
		AddRow(mockPromo[1].Id, mockPromo[1].CreatedBy, mockPromo[1].CreatedDate, mockPromo[1].ModifiedBy,
		mockPromo[1].ModifiedDate, mockPromo[1].DeletedBy, mockPromo[1].DeletedDate, mockPromo[1].IsDeleted,
		mockPromo[1].IsActive, mockPromo[1].PromoCode, mockPromo[1].PromoName, mockPromo[1].PromoDesc, mockPromo[1].PromoValue,
		mockPromo[1].PromoType, mockPromo[1].PromoImage, mockPromo[1].StartDate, mockPromo[1].EndDate, mockPromo[1].MaxUsage,
		mockPromo[1].ProductionCapacity, mockPromo[1].CurrencyId, mockPromo[1].PromoProductType, mockPromo[1].StartTripPeriod,
		mockPromo[1].EndTripPeriod, mockPromo[1].HowToGet, mockPromo[1].HowToUse, mockPromo[1].TermCondition, mockPromo[1].Disclaimer,
		mockPromo[1].MaxDiscount, mockPromo[1].IsAnyTripPeriod)

	limit := 10
	offset := 0
	page := &offset
	size := &limit
	trans := true
	exp := true
	promoId := mockPromo[0].Id
	search := "test"
	sortBy := "newest"
	query := `SELECT p.\*\ FROM promos p `
	if len(merchantIds) != 0 {
		query = query + ` JOIN promo_merchants pm ON p.id = pm.promo_id `
	}
	if trans == true || exp == true {
		query = query + ` JOIN promo_experience_transports pet ON p.id = pet.promo_id `
	}
	query = query + ` WHERE p.is_deleted = 0 AND p.is_active = 1 `
	if trans == true {
		query = query + ` AND pet.transportation_id != '' `
	}
	if exp == true {
		query = query + ` AND pet.experience_id != '' `
	}

	for index, id := range merchantIds {
		if index == 0 && index != (len(merchantIds)-1) {
			query = query + ` AND \(pm.merchant_id ='` + id + `' `
		} else if index == 0 && index == (len(merchantIds)-1) {
			query = query + ` AND \(pm.merchant_id ='` + id + `' \) `
		} else if index == (len(merchantIds) - 1) {
			query = query + ` OR pm.merchant_id ='` + id + `' \) `
		} else {
			query = query + ` OR pm.merchant_id ='` + id + `' `
		}
	}
	if promoId != "" {
		query = query + ` AND p.id ='` + promoId + `' `
	}
	if search != "" {
		query = query + ` AND \(promo_name LIKE '%` + search + `%'` +
			`OR promo_desc LIKE '%` + search + `%' ` +
			`OR start_date LIKE '%` + search + `%' ` +
			`OR end_date LIKE '%` + search + `%' ` +
			`OR promo_code LIKE '%` + search + `%' ` +
			`OR max_usage LIKE '%` + search + `%' ` + `\) `
	}
	if sortBy == "newest" {
		query = query + ` ORDER BY created_date desc `
	} else if sortBy == "latest" {
		query = query + ` ORDER BY created_date asc `
	} else {
		query = query + ` ORDER BY created_date desc `
	}

	if page != nil && size != nil {
		query = query + ` LIMIT ` + strconv.Itoa(*size) + ` OFFSET ` + strconv.Itoa(*page) + ` `
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	anArticle, err := a.Fetch(context.TODO(), page, size,search,trans,exp,merchantIds,sortBy,promoId)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestFetchWithoutPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"}).
		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
			mockPromo[0].IsActive, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
			mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod).
		AddRow(mockPromo[1].Id, mockPromo[1].CreatedBy, mockPromo[1].CreatedDate, mockPromo[1].ModifiedBy,
			mockPromo[1].ModifiedDate, mockPromo[1].DeletedBy, mockPromo[1].DeletedDate, mockPromo[1].IsDeleted,
			mockPromo[1].IsActive, mockPromo[1].PromoCode, mockPromo[1].PromoName, mockPromo[1].PromoDesc, mockPromo[1].PromoValue,
			mockPromo[1].PromoType, mockPromo[1].PromoImage, mockPromo[1].StartDate, mockPromo[1].EndDate, mockPromo[1].MaxUsage,
			mockPromo[1].ProductionCapacity, mockPromo[1].CurrencyId, mockPromo[1].PromoProductType, mockPromo[1].StartTripPeriod,
			mockPromo[1].EndTripPeriod, mockPromo[1].HowToGet, mockPromo[1].HowToUse, mockPromo[1].TermCondition, mockPromo[1].Disclaimer,
			mockPromo[1].MaxDiscount, mockPromo[1].IsAnyTripPeriod)
	sortBy:= ""
	query := `SELECT p.\*\ FROM promos p `
	query = query + ` WHERE p.is_deleted = 0 AND p.is_active = 1 `
	if sortBy == "newest" {
		query = query + ` ORDER BY created_date desc `
	} else if sortBy == "latest" {
		query = query + ` ORDER BY created_date asc `
	} else {
		query = query + ` ORDER BY created_date desc `
	}

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), nil, nil,"",false,false,make([]string,0),"","")
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
//func TestFetchWithPaginationSearch(t *testing.T) {
//	db, mock, err := sqlmock.New()
//
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
//		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
//		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
//		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
//		"is_any_trip_period"}).
//		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
//			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
//			mockPromo[0].IsActive, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
//			mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
//			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
//			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
//			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod).
//		AddRow(mockPromo[1].Id, mockPromo[1].CreatedBy, mockPromo[1].CreatedDate, mockPromo[1].ModifiedBy,
//			mockPromo[1].ModifiedDate, mockPromo[1].DeletedBy, mockPromo[1].DeletedDate, mockPromo[1].IsDeleted,
//			mockPromo[1].IsActive, mockPromo[1].PromoCode, mockPromo[1].PromoName, mockPromo[1].PromoDesc, mockPromo[1].PromoValue,
//			mockPromo[1].PromoType, mockPromo[1].PromoImage, mockPromo[1].StartDate, mockPromo[1].EndDate, mockPromo[1].MaxUsage,
//			mockPromo[1].ProductionCapacity, mockPromo[1].CurrencyId, mockPromo[1].PromoProductType, mockPromo[1].StartTripPeriod,
//			mockPromo[1].EndTripPeriod, mockPromo[1].HowToGet, mockPromo[1].HowToUse, mockPromo[1].TermCondition, mockPromo[1].Disclaimer,
//			mockPromo[1].MaxDiscount, mockPromo[1].IsAnyTripPeriod)
//	search := mockPromo[0].PromoName
//	query := `SELECT \*\ FROM promos where is_deleted = 0 AND is_active = 1`
//	query = query + ` AND (promo_name LIKE '%` + search + `%'` +
//		`OR promo_desc LIKE '%` + search + `%' ` +
//		`OR start_date LIKE '%` + search + `%' ` +
//		`OR end_date LIKE '%` + search + `%' ` +
//		`OR promo_code LIKE '%` + search + `%' ` +
//		`OR max_usage LIKE '%` + search + `%' ` + `) `
//
//	query = query + ` ORDER BY created_date desc LIMIT \? OFFSET \? `
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := PromoRepo.NewpromoRepository(db)
//
//	limit := 10
//	offset := 0
//	anArticle, err := a.Fetch(context.TODO(), &limit, &offset,search)
//	//assert.NotEmpty(t, nextCursor)
//	assert.NoError(t, err)
//	assert.Len(t, anArticle, 2)
//}
//func TestFetchWithoutPaginationSearch(t *testing.T) {
//	db, mock, err := sqlmock.New()
//
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
//		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
//		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
//		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
//		"is_any_trip_period"}).
//		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
//			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
//			mockPromo[0].IsActive, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
//			mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
//			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
//			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
//			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod).
//		AddRow(mockPromo[1].Id, mockPromo[1].CreatedBy, mockPromo[1].CreatedDate, mockPromo[1].ModifiedBy,
//			mockPromo[1].ModifiedDate, mockPromo[1].DeletedBy, mockPromo[1].DeletedDate, mockPromo[1].IsDeleted,
//			mockPromo[1].IsActive, mockPromo[1].PromoCode, mockPromo[1].PromoName, mockPromo[1].PromoDesc, mockPromo[1].PromoValue,
//			mockPromo[1].PromoType, mockPromo[1].PromoImage, mockPromo[1].StartDate, mockPromo[1].EndDate, mockPromo[1].MaxUsage,
//			mockPromo[1].ProductionCapacity, mockPromo[1].CurrencyId, mockPromo[1].PromoProductType, mockPromo[1].StartTripPeriod,
//			mockPromo[1].EndTripPeriod, mockPromo[1].HowToGet, mockPromo[1].HowToUse, mockPromo[1].TermCondition, mockPromo[1].Disclaimer,
//			mockPromo[1].MaxDiscount, mockPromo[1].IsAnyTripPeriod)
//
//	query := `SELECT \*\ FROM promos where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := PromoRepo.NewpromoRepository(db)
//
//	//limit := 10
//	//offset := 0
//	anArticle, err := a.Fetch(context.TODO(), nil, nil,"")
//	//assert.NotEmpty(t, nextCursor)
//	assert.NoError(t, err)
//	assert.Len(t, anArticle, 2)
//}
func TestFetchWithPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"}).
		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
			mockPromo[0].ModifiedDate, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
			mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod).
		AddRow(mockPromo[1].Id, mockPromo[1].CreatedBy, mockPromo[1].CreatedDate, mockPromo[1].ModifiedBy,
			mockPromo[1].ModifiedDate, mockPromo[1].DeletedBy, mockPromo[1].DeletedDate, mockPromo[1].IsDeleted,
			mockPromo[1].ModifiedDate, mockPromo[1].PromoCode, mockPromo[1].PromoName, mockPromo[1].PromoDesc, mockPromo[1].PromoValue,
			mockPromo[1].PromoType, mockPromo[1].PromoImage, mockPromo[1].StartDate, mockPromo[1].EndDate, mockPromo[1].MaxUsage,
			mockPromo[1].ProductionCapacity, mockPromo[1].CurrencyId, mockPromo[1].PromoProductType, mockPromo[1].StartTripPeriod,
			mockPromo[1].EndTripPeriod, mockPromo[1].HowToGet, mockPromo[1].HowToUse, mockPromo[1].TermCondition, mockPromo[1].Disclaimer,
			mockPromo[1].MaxDiscount, mockPromo[1].IsAnyTripPeriod)

	query := `SELECT \*\ FROM promos where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.Fetch(context.TODO(), &limit, &offset,"",false,false,make([]string,0),"","")
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
	//assert.Len(t, anArticle, 2)
}
func TestFetchWithoutPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"}).
		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
			mockPromo[0].ModifiedDate, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
			mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod).
		AddRow(mockPromo[1].Id, mockPromo[1].CreatedBy, mockPromo[1].CreatedDate, mockPromo[1].ModifiedBy,
			mockPromo[1].ModifiedDate, mockPromo[1].DeletedBy, mockPromo[1].DeletedDate, mockPromo[1].IsDeleted,
			mockPromo[1].IsActive, mockPromo[1].PromoCode, mockPromo[1].PromoName, mockPromo[1].PromoDesc, mockPromo[1].PromoValue,
			mockPromo[1].PromoType, mockPromo[1].PromoImage, mockPromo[1].StartDate, mockPromo[1].EndDate, mockPromo[1].MaxUsage,
			mockPromo[1].ProductionCapacity, mockPromo[1].CurrencyId, mockPromo[1].PromoProductType, mockPromo[1].StartTripPeriod,
			mockPromo[1].EndTripPeriod, mockPromo[1].HowToGet, mockPromo[1].HowToUse, mockPromo[1].TermCondition, mockPromo[1].Disclaimer,
			mockPromo[1].MaxDiscount, mockPromo[1].IsAnyTripPeriod)

	query := `SELECT \*\ FROM promos where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), nil, nil,"",false,false,make([]string,0),"","")
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()


	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"}).
		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
			mockPromo[0].IsActive, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
			mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod)

	query := `SELECT \*\ FROM promos WHERE is_deleted = 0 and is_active = 1 and id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	num := mockPromo[0].Id
	anArticle, err := a.GetById(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByIDNotfound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()


	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"})

	query := `SELECT \*\ FROM promos WHERE is_deleted = 0 and is_active = 1 and id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	num := mockPromo[0].Id
	anArticle, err := a.GetById(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByIDErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()



	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"}).
		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
			mockPromo[0].ModifiedDate, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
			mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod)

	query := `SELECT \*\ FROM promos WHERE is_deleted = 0 and is_active = 1 and id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	num := mockPromo[0].Id
	anArticle, err := a.GetById(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByCodeWithMerchantId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"}).
		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
			mockPromo[0].IsActive, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
			mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod)

	promoType := strconv.Itoa(*mockPromo[0].PromoProductType)
	promoUseDate := time.Now().Format("2006-01-02")
	query := `SELECT p.\*\ 
				FROM 
					promos p `

	query = query + ` WHERE 
						BINARY p.promo_code = \?  AND 
						p.is_deleted = 0 AND 
						p.is_active = 1 `
	if promoType != "" {
		query = query + ` AND 
						p.promo_product_type in \(0,` + promoType + `\) `
	}

	if promoUseDate != "" {
		query = query + ` AND \(DATE\('` + promoUseDate + `'\) >= p.start_date AND 
								DATE\('` + promoUseDate + `'\) <= p.end_date\)  `
	}


	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	code := mockPromo[0].PromoCode
	anArticle, err := a.GetByCode(context.TODO(), code,promoType,"","","","","",promoUseDate)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByCodeWithMerchantIdNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"})

	merchantId := "adqerqewqsDSADAD"
	query := `SELECT p.\*\ 
				FROM 
					promos p
				JOIN promo_merchants pm on pm.promo_id = p.id
				WHERE 
					BINARY p.promo_code = \? AND 
					p.promo_product_type = \? AND 
 					p.is_deleted = 0 AND 
					p.is_active = 1 AND
					pm.merchant_id = '` + merchantId  + `'`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	code := mockPromo[0].PromoCode
	promoType := mockPromo[0].PromoProductType
	anArticle, err := a.GetByCode(context.TODO(), code,strconv.Itoa(*promoType),merchantId,"","","","","")
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByCodeWithMerchantIdErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "promo_code", "promo_name", "promo_desc", "promo_value",
		"promo_type", "promo_image", "start_date", "end_date", "max_usage", "production_capacity", "currency_id", "promo_product_type",
		"start_trip_period", "end_trip_period", "how_to_get", "how_to_use", "term_condition", "disclaimer", "max_discount",
		"is_any_trip_period"}).
		AddRow(mockPromo[0].Id, mockPromo[0].CreatedBy, mockPromo[0].CreatedDate, mockPromo[0].ModifiedBy,
			mockPromo[0].ModifiedDate, mockPromo[0].DeletedBy, mockPromo[0].DeletedDate, mockPromo[0].IsDeleted,
			mockPromo[0].ModifiedDate, mockPromo[0].PromoCode, mockPromo[0].PromoName, mockPromo[0].PromoDesc, mockPromo[0].PromoValue,
			mockPromo[0].PromoType, mockPromo[0].PromoImage, mockPromo[0].StartDate, mockPromo[0].EndDate, mockPromo[0].MaxUsage,
			mockPromo[0].ProductionCapacity, mockPromo[0].CurrencyId, mockPromo[0].PromoProductType, mockPromo[0].StartTripPeriod,
			mockPromo[0].EndTripPeriod, mockPromo[0].HowToGet, mockPromo[0].HowToUse, mockPromo[0].TermCondition, mockPromo[0].Disclaimer,
			mockPromo[0].MaxDiscount, mockPromo[0].IsAnyTripPeriod)
	merchantId := "adqerqewqsDSADAD"
	query := `SELECT p.\*\ 
				FROM 
					promos p
				JOIN promo_merchants pm on pm.promo_id = p.id
				WHERE 
					BINARY p.promo_code = \? AND 
					p.promo_product_type = \? AND 
 					p.is_deleted = 0 AND 
					p.is_active = 1 AND
					pm.merchant_id = '` + merchantId  + `'`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := PromoRepo.NewpromoRepository(db)

	code := mockPromo[0].PromoCode
	promoType := mockPromo[0].PromoProductType
	anArticle, err := a.GetByCode(context.TODO(), code,strconv.Itoa(*promoType),merchantId,"","","","","")
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE promos SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := mockPromo[0].Id
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := PromoRepo.NewpromoRepository(db)

	err = a.Delete(context.TODO(), id, deletedBy)
	assert.NoError(t, err)
}
func TestDeleteErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE promos SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := mockPromo[0].Id
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id, id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := PromoRepo.NewpromoRepository(db)

	err = a.Delete(context.TODO(), id, deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {
	//user := "test"
	//now := time.Now()
	a := mockPromo[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT promos SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , promo_code=\?,promo_name=\? , 
				promo_desc=\? ,promo_value=\?,promo_type=\?,promo_image=\?,start_date=\?,end_date=\?,currency_id	=\?,
				max_usage=\?,production_capacity=\?,promo_product_type=\?,start_trip_period=\?,end_trip_period=\?,
				how_to_get=\?,how_to_use=\?,term_condition=\?,disclaimer=\?,max_discount=\?,is_any_trip_period=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.PromoCode, a.PromoName,
		a.PromoDesc, a.PromoValue, a.PromoType, a.PromoImage, a.StartDate, a.EndDate, a.CurrencyId, a.MaxUsage,
		a.ProductionCapacity,a.PromoProductType,a.StartTripPeriod,a.EndTripPeriod,a.HowToGet,a.HowToUse,a.TermCondition,
		a.Disclaimer,a.MaxDiscount,a.IsAnyTripPeriod).WillReturnResult(sqlmock.NewResult(1, 1))

	i := PromoRepo.NewpromoRepository(db)

	id, err := i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	//user := "test"
	//now := time.Now()
	a := mockPromo[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT promos SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , promo_code=\?,promo_name=\? , 
				promo_desc=\? ,promo_value=\?,promo_type=\?,promo_image=\?,start_date=\?,end_date=\?,currency_id	=\?,
				max_usage=\?,production_capacity=\?,promo_product_type=\?,start_trip_period=\?,end_trip_period=\?,
				how_to_get=\?,how_to_use=\?,term_condition=\?,disclaimer=\?,max_discount=\?,is_any_trip_period=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.PromoCode, a.PromoName,
		a.PromoDesc, a.PromoValue, a.PromoType, a.PromoImage, a.StartDate, a.EndDate, a.CurrencyId, a.MaxUsage,
		a.ProductionCapacity,a.PromoProductType,a.StartTripPeriod,a.EndTripPeriod,a.HowToGet,a.HowToUse,a.TermCondition,
		a.Disclaimer,a.MaxDiscount,a.IsAnyTripPeriod,a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := PromoRepo.NewpromoRepository(db)

	_, err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockPromo[0]
	a.ModifiedDate = &now
	a.ModifiedBy = &modifyBy

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE promos set modified_by=\?, modified_date=\? , promo_code=\?,promo_name=\? , promo_desc=\? ,promo_value=\?,
				promo_type=\?,promo_image=\?,start_date=\?,end_date=\?,currency_id=\?,max_usage=\?,production_capacity=\?, 
				promo_product_type=\?,start_trip_period=\?,end_trip_period=\?,
				how_to_get=\?,how_to_use=\?,term_condition=\?,disclaimer=\?,max_discount=\?,is_any_trip_period=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, a.ModifiedDate, a.PromoCode, a.PromoName, a.PromoDesc, a.PromoValue,
		a.PromoType, a.PromoImage, a.StartDate, a.EndDate, a.CurrencyId, a.MaxUsage, a.ProductionCapacity, a.PromoProductType,
		a.StartTripPeriod, a.EndTripPeriod, a.HowToGet, a.HowToUse, a.TermCondition, a.Disclaimer,a.MaxDiscount,
		a.IsAnyTripPeriod, a.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := PromoRepo.NewpromoRepository(db)

	err = u.Update(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Nil(t, err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockPromo[0]
	a.ModifiedDate = &now
	a.ModifiedBy = &modifyBy

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE promos set modified_by=\?, modified_date=\? , promo_code=\?,promo_name=\? , promo_desc=\? ,promo_value=\?,
				promo_type=\?,promo_image=\?,start_date=\?,end_date=\?,currency_id=\?,max_usage=\?,production_capacity=\?, 
				promo_product_type=\?,start_trip_period=\?,end_trip_period=\?,
				how_to_get=\?,how_to_use=\?,term_condition=\?,disclaimer=\?,max_discount=\?,is_any_trip_period=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, time.Now(), a.PromoCode, a.PromoName, a.PromoDesc, a.PromoValue,
		a.PromoType, a.PromoImage, a.StartDate, a.EndDate, a.CurrencyId, a.MaxUsage, a.ProductionCapacity, a.PromoProductType,
		a.StartTripPeriod, a.EndTripPeriod, a.HowToGet, a.HowToUse, a.TermCondition, a.Disclaimer,a.MaxDiscount,
		a.IsAnyTripPeriod, a.Id,a.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := PromoRepo.NewpromoRepository(db)

	err = u.Update(context.TODO(), &a)
	assert.Error(t, err)
}
