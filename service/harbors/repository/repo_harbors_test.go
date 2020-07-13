package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	HarborsRepo "github.com/service/harbors/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var(
	harborsType = 1
	provinceName = "Jawa Barat"
	mockHarbors = []models.Harbors{
		models.Harbors{
			Id:               "jlkjlkjlkjlkjlkj",
			CreatedBy:        "test",
			CreatedDate:      time.Now(),
			ModifiedBy:       nil,
			ModifiedDate:     nil,
			DeletedBy:        nil,
			DeletedDate:      nil,
			IsDeleted:        0,
			IsActive:         1,
			HarborsName:      "Harbors Test 1",
			HarborsLongitude: 1213,
			HarborsLatitude:  12313,
			HarborsImage:     "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
			CityId:           1,
			HarborsType:      &harborsType,
		},
		models.Harbors{
			Id:               "jlkjlkjlkjlkjlkj",
			CreatedBy:        "test",
			CreatedDate:      time.Now(),
			ModifiedBy:       nil,
			ModifiedDate:     nil,
			DeletedBy:        nil,
			DeletedDate:      nil,
			IsDeleted:        0,
			IsActive:         1,
			HarborsName:      "Harbors Test 1",
			HarborsLongitude: 1213,
			HarborsLatitude:  12313,
			HarborsImage:     "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
			CityId:           1,
			HarborsType:      &harborsType,
		},
	}
	mockHarborsJoin = []models.HarborsWCPC{
		models.HarborsWCPC{
			Id:               "dfgdgdgfdgfdgf",
			HarborsName:      "harbors Test 1",
			HarborsLongitude: 1231231231,
			HarborsLatitude:  43242342,
			HarborsImage:     "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
			CityId:           1,
			CityName:         "Bogor",
			ProvinceId:       1,
			ProvinceName:     &provinceName,
			CountryName:      "Indonesia",
		},
		models.HarborsWCPC{
			Id:               "dfgdgdgfdgfdgf",
			HarborsName:      "harbors Test 2",
			HarborsLongitude: 1231231231,
			HarborsLatitude:  43242342,
			HarborsImage:     "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
			CityId:           1,
			CityName:         "Bogor",
			ProvinceId:       1,
			ProvinceName:     &provinceName,
			CountryName:      "Indonesia",
		},
	}
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
		AddRow(len(mockHarbors))

	query := `SELECT count\(\*\) AS count FROM harbors WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

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

	query := `SELECT count\(\*\) AS count FROM harbors WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

	_, err = a.GetCount(context.TODO())
	assert.Error(t, err)
}
func TestGetAllWithJoinCPCWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "harbors_name", "harbors_longitude", "harbors_latitude",
		"harbors_image", "city_id", "city_name","province_id","province_name","country_name"}).
		AddRow(mockHarborsJoin[0].Id, mockHarborsJoin[0].HarborsName, mockHarborsJoin[0].HarborsLongitude,
		mockHarborsJoin[0].HarborsLatitude, mockHarborsJoin[0].HarborsImage, mockHarborsJoin[0].CityId,
		mockHarborsJoin[0].CityName,mockHarborsJoin[0].ProvinceId,mockHarborsJoin[0].ProvinceName,
		mockHarborsJoin[0].CountryName).
		AddRow(mockHarborsJoin[1].Id, mockHarborsJoin[1].HarborsName, mockHarborsJoin[1].HarborsLongitude,
		mockHarborsJoin[1].HarborsLatitude, mockHarborsJoin[1].HarborsImage, mockHarborsJoin[1].CityId,
		mockHarborsJoin[1].CityName,mockHarborsJoin[1].ProvinceId,mockHarborsJoin[1].ProvinceName,
		mockHarborsJoin[1].CountryName)

	query := `Select 
				h.id, 
				h.harbors_name,
				h.harbors_longitude,
				h.harbors_latitude,
				h.harbors_image,
				h.city_id ,
				c.city_name,
				p.id as province_id,
				p.province_name,
				co.country_name 
			from cgo_indonesia.harbors h
			join cities c on h.city_id = c.id
			join provinces p on c.province_id = p.id
			join countries co on p.country_id = co.id
			where h.is_active = 1 and h.is_deleted = 0`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)
	limit := 10
	offset := 0
	//search := ""
	anArticle, err := a.GetAllWithJoinCPC(context.TODO(),&limit,&offset,"","")
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetAllWithJoinCPCErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "harbors_name", "harbors_longitude", "harbors_latitude",
		"harbors_image", "city_id", "city_name","province_id","province_name","country_name"}).
		AddRow(mockHarborsJoin[0].Id, mockHarborsJoin[0].HarborsName, mockHarborsJoin[0].HarborsLongitude,
			mockHarborsJoin[0].CityName, mockHarborsJoin[0].HarborsImage, mockHarborsJoin[0].CityId,
			mockHarborsJoin[0].CityName,mockHarborsJoin[0].ProvinceId,mockHarborsJoin[0].ProvinceName,
			mockHarborsJoin[0].CountryName).
		AddRow(mockHarborsJoin[1].Id, mockHarborsJoin[1].HarborsName, mockHarborsJoin[1].HarborsLongitude,
			mockHarborsJoin[1].HarborsLatitude, mockHarborsJoin[1].HarborsImage, mockHarborsJoin[1].CityId,
			mockHarborsJoin[1].CityName,mockHarborsJoin[1].ProvinceId,mockHarborsJoin[1].ProvinceName,
			mockHarborsJoin[1].CountryName)

	query := `Select 
				h.id, 
				h.harbors_name,
				h.harbors_longitude,
				h.harbors_latitude,
				h.harbors_image,
				h.city_id ,
				c.city_name,
				p.id as province_id,
				p.province_name,
				co.country_name 
			from cgo_indonesia.harbors h
			join cities c on h.city_id = c.id
			join provinces p on c.province_id = p.id
			join countries co on p.country_id = co.id
			where h.is_active = 1 and h.is_deleted = 0`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

	limit := 10
	offset := 0
	//search := ""
	anArticle, err := a.GetAllWithJoinCPC(context.TODO(),&limit,&offset,"","")
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
//func TestGetAllWithJoinCPCSearchWithPagination(t *testing.T) {
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
//	rows := sqlmock.NewRows([]string{"id", "harbors_name", "harbors_longitude", "harbors_latitude",
//		"harbors_image", "city_id", "city_name","province_id","province_name","country_name"}).
//		AddRow(mockHarborsJoin[0].Id, mockHarborsJoin[0].HarborsName, mockHarborsJoin[0].HarborsLongitude,
//			mockHarborsJoin[0].HarborsLatitude, mockHarborsJoin[0].HarborsImage, mockHarborsJoin[0].CityId,
//			mockHarborsJoin[0].CityName,mockHarborsJoin[0].ProvinceId,mockHarborsJoin[0].ProvinceName,
//			mockHarborsJoin[0].CountryName).
//		AddRow(mockHarborsJoin[1].Id, mockHarborsJoin[1].HarborsName, mockHarborsJoin[1].HarborsLongitude,
//			mockHarborsJoin[1].HarborsLatitude, mockHarborsJoin[1].HarborsImage, mockHarborsJoin[1].CityId,
//			mockHarborsJoin[1].CityName,mockHarborsJoin[1].ProvinceId,mockHarborsJoin[1].ProvinceName,
//			mockHarborsJoin[1].CountryName)
//
//	query := `Select
//				h.id,
//				h.harbors_name,
//				h.harbors_longitude,
//				h.harbors_latitude,
//				h.harbors_image,
//				h.city_id ,
//				c.city_name,
//				p.id as province_id,
//				p.province_name,
//				co.country_name
//			from cgo_indonesia.harbors h
//			join cities c on h.city_id = c.id
//			join provinces p on c.province_id = p.id
//			join countries co on p.country_id = co.id
//			where h.is_active = 1 and h.is_deleted = 0
//			AND (h.harbors_name LIKE \? OR c.city_name LIKE \? OR p.province_name LIKE \?)`
//	query = query + ` ORDER BY h.created_date desc LIMIT ? OFFSET ? `
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := HarborsRepo.NewharborsRepository(db)
//	limit := 10
//	offset := 0
//	search := "Bogor"
//	anArticle, err := a.GetAllWithJoinCPC(context.TODO(),&limit,&offset,search,"")
//	//assert.NotEmpty(t, nextCursor)
//	assert.NoError(t, err)
//	assert.Len(t, anArticle, 2)
//}
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
		"deleted_date", "is_deleted", "is_active", "harbors_name", "harbors_longitude", "harbors_latitude",
		"harbors_image", "city_id", "harbors_type"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].IsActive, mockHarbors[0].HarborsName, mockHarbors[0].HarborsLongitude,
			mockHarbors[0].HarborsLatitude, mockHarbors[0].HarborsImage, mockHarbors[0].CityId,
			mockHarbors[0].HarborsType).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].IsActive, mockHarbors[1].HarborsName, mockHarbors[1].HarborsLongitude,
			mockHarbors[1].HarborsLatitude, mockHarbors[1].HarborsImage, mockHarbors[1].CityId,
			mockHarbors[1].HarborsType)

	query := `SELECT \*\ FROM harbors where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.Fetch(context.TODO(), limit, offset)
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
		"deleted_date", "is_deleted", "is_active", "harbors_name", "harbors_longitude", "harbors_latitude",
		"harbors_image", "city_id", "harbors_type"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].IsActive, mockHarbors[0].HarborsName, mockHarbors[0].HarborsLongitude,
			mockHarbors[0].HarborsLatitude, mockHarbors[0].HarborsImage, mockHarbors[0].CityId,
			mockHarbors[0].HarborsType).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].IsActive, mockHarbors[1].HarborsName, mockHarbors[1].HarborsLongitude,
			mockHarbors[1].HarborsLatitude, mockHarbors[1].HarborsImage, mockHarbors[1].CityId,
			mockHarbors[1].HarborsType)

	query := `SELECT \*\ FROM harbors where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), 0, 0)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
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
		"deleted_date", "is_deleted", "is_active", "harbors_name", "harbors_longitude", "harbors_latitude",
		"harbors_image", "city_id", "harbors_type"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].ModifiedDate, mockHarbors[0].HarborsName, mockHarbors[0].HarborsLongitude,
			mockHarbors[0].HarborsLatitude, mockHarbors[0].HarborsImage, mockHarbors[0].CityId,
			mockHarbors[0].HarborsType).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].IsActive, mockHarbors[1].HarborsName, mockHarbors[1].HarborsLongitude,
			mockHarbors[1].HarborsLatitude, mockHarbors[1].HarborsImage, mockHarbors[1].CityId,
			mockHarbors[1].HarborsType)

	query := `SELECT \*\ FROM harbors where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.Fetch(context.TODO(), limit, offset)
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
		"deleted_date", "is_deleted", "is_active", "harbors_name", "harbors_longitude", "harbors_latitude",
		"harbors_image", "city_id", "harbors_type"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].ModifiedDate, mockHarbors[0].HarborsName, mockHarbors[0].HarborsLongitude,
			mockHarbors[0].HarborsLatitude, mockHarbors[0].HarborsImage, mockHarbors[0].CityId,
			mockHarbors[0].HarborsType).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].IsActive, mockHarbors[1].HarborsName, mockHarbors[1].HarborsLongitude,
			mockHarbors[1].HarborsLatitude, mockHarbors[1].HarborsImage, mockHarbors[1].CityId,
			mockHarbors[1].HarborsType)

	query := `SELECT \*\ FROM harbors where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), 0, 0)
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
		"deleted_date", "is_deleted", "is_active", "harbors_name", "harbors_longitude", "harbors_latitude",
		"harbors_image", "city_id", "harbors_type"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].IsActive, mockHarbors[0].HarborsName, mockHarbors[0].HarborsLongitude,
			mockHarbors[0].HarborsLatitude, mockHarbors[0].HarborsImage, mockHarbors[0].CityId,
			mockHarbors[0].HarborsType).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].IsActive, mockHarbors[1].HarborsName, mockHarbors[1].HarborsLongitude,
			mockHarbors[1].HarborsLatitude, mockHarbors[1].HarborsImage, mockHarbors[1].CityId,
			mockHarbors[1].HarborsType)

	query := `SELECT \*\ FROM harbors WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

	num := "jlkjlkjlkjlkjlkj"
	anArticle, err := a.GetByID(context.TODO(), num)
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
		"deleted_date", "is_deleted", "is_active", "harbors_name", "harbors_longitude", "harbors_latitude",
		"harbors_image", "city_id", "harbors_type"})

	query := `SELECT \*\ FROM harbors WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

	num := "jlkjlkjlkjlkjlkj"
	anArticle, err := a.GetByID(context.TODO(), num)
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
		"deleted_date", "is_deleted", "is_active", "harbors_name", "harbors_longitude", "harbors_latitude",
		"harbors_image", "city_id", "harbors_type"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].ModifiedDate, mockHarbors[0].HarborsName, mockHarbors[0].HarborsLongitude,
			mockHarbors[0].HarborsLatitude, mockHarbors[0].HarborsImage, mockHarbors[0].CityId,
			mockHarbors[0].HarborsType)

	query := `SELECT \*\ FROM harbors WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewharborsRepository(db)

	num := "jlkjlkjlkjlkjlkj"
	anArticle, err := a.GetByID(context.TODO(), num)
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

	query := "UPDATE harbors SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := "jlkjlkjlkjlkjlkj"
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := HarborsRepo.NewharborsRepository(db)

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

	query := "UPDATE harbors SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := "jlkjlkjlkjlkjlkj"
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := HarborsRepo.NewharborsRepository(db)

	err = a.Delete(context.TODO(), id, deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {

	a := mockHarbors[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT harbors SET id=\?,created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , deleted_by=\? , 
				deleted_date=\? , is_deleted=\? , is_active=\? , harbors_name=\? , harbors_longitude=\? , harbors_latitude=\? ,
				harbors_image=\?,city_id=\?,harbors_type=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs( a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.HarborsName, a.HarborsLongitude,
		a.HarborsLatitude, a.HarborsImage, a.CityId,a.HarborsType).WillReturnResult(sqlmock.NewResult(2, 1))

	i := HarborsRepo.NewharborsRepository(db)

	id, err := i.Insert(context.TODO(), &a)
	//a.Id = *id
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {

	a := mockHarbors[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()


	query := `INSERT harbors SET id=\?,created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , deleted_by=\? , 
				deleted_date=\? , is_deleted=\? , is_active=\? , harbors_name=\? , harbors_longitude=\? , harbors_latitude=\? ,
				harbors_image=\?,city_id=\?,harbors_type=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs( a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.HarborsName, a.HarborsLongitude,
		a.HarborsLatitude, a.HarborsImage, a.CityId,a.HarborsType,a.Id,a.Id).WillReturnResult(sqlmock.NewResult(2, 1))

	i := HarborsRepo.NewharborsRepository(db)

	_, err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockHarbors[0]
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

	query := `UPDATE harbors set modified_by=\?, modified_date=\? ,  harbors_name=\? , harbors_longitude=\? , harbors_latitude=\? ,
				harbors_image=\?,city_id=\?,harbors_type=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, time.Now(), a.HarborsName, a.HarborsLongitude,
		a.HarborsLatitude, a.HarborsImage, a.CityId, a.HarborsType,a.Id).WillReturnResult(sqlmock.NewResult(2, 1))

	u := HarborsRepo.NewharborsRepository(db)

	err = u.Update(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Nil(t, err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockHarbors[0]
	a.ModifiedBy = &modifyBy
	a.ModifiedDate = &now

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE harbors set modified_by=\?, modified_date=\? ,  harbors_name=\? , harbors_longitude=\? , harbors_latitude=\? ,
				harbors_image=\?,city_id=\?,harbors_type=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, time.Now(), a.HarborsName, a.HarborsLongitude,
		a.HarborsLatitude, a.HarborsImage, a.CityId, a.HarborsType,a.Id,a.Id).WillReturnResult(sqlmock.NewResult(2, 1))


	u := HarborsRepo.NewharborsRepository(db)

	err = u.Update(context.TODO(), &a)

	assert.Error(t, err)
}
