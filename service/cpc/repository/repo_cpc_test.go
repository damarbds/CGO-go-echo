package repository_test

import (
	"context"
	"github.com/models"
	cpcRepo "github.com/service/cpc/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var (
	imageTestPath = "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg"
	code = 0
	mockCity = []models.City{
		models.City{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			CityName:  "Bogor",
			CityDesc:  "Bogor adalah kota hujan",
			CityPhotos:&imageTestPath,
			ProvinceId:1,
		},
		models.City{
			Id:           2,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			CityName:  "Jakarta",
			CityDesc:  "Jakarta adalah kota maced",
			CityPhotos:&imageTestPath,
			ProvinceId:1,
		},
	}
	mockProvince = []models.Province{
		models.Province{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ProvinceName:  "Jawa Barat",
			ProvinceNameTransportation:&imageTestPath,
			CountryId:1,
		},
		models.Province{
			Id:           2,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ProvinceName:  "DKI Jakarta",
			ProvinceNameTransportation:&imageTestPath,
			CountryId:1,
		},
	}
	mockCountry = []models.Country{
		models.Country{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			CountryName:  "Indonesia",
			Iso:          &imageTestPath,
			Name:         &imageTestPath,
			NiceName:     &imageTestPath,
			Iso3:         &imageTestPath,
			NumCode:      &code,
			PhoneCode:    &code,
		},
		models.Country{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			CountryName:  "Indonesia",
			Iso:          &imageTestPath,
			Name:         &imageTestPath,
			NiceName:     &imageTestPath,
			Iso3:         &imageTestPath,
			NumCode:      &code,
			PhoneCode:    &code,
		},
	}
)
func TestCountCity(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(len(mockCity))

	query := `SELECT count\(\*\) AS count FROM cities WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	res, err := a.GetCountCity(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, res, 2, "")
}
func TestCountCityErrorFetch(t *testing.T) {
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

	query := `SELECT count\(\*\) AS count FROM cities WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	_, err = a.GetCountCity(context.TODO())
	assert.Error(t, err)
}
func TestFetchCityWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","city_name","city_desc","city_photos","province_id"}).
		AddRow(mockCity[0].Id, mockCity[0].CreatedBy,mockCity[0].CreatedDate,mockCity[0].ModifiedBy,
		mockCity[0].ModifiedDate,mockCity[0].DeletedBy,mockCity[0].DeletedDate,mockCity[0].IsDeleted,
		mockCity[0].IsActive,mockCity[0].CityName,mockCity[0].CityDesc,mockCity[0].CityPhotos,mockCity[0].ProvinceId).
		AddRow(mockCity[1].Id, mockCity[1].CreatedBy,mockCity[1].CreatedDate,mockCity[1].ModifiedBy,
		mockCity[1].ModifiedDate,mockCity[1].DeletedBy,mockCity[1].DeletedDate,mockCity[1].IsDeleted,
		mockCity[1].IsActive,mockCity[1].CityName,mockCity[1].CityDesc,mockCity[1].CityPhotos,mockCity[1].ProvinceId)

	query := `SELECT \*\ FROM cities where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.FetchCity(context.TODO(), limit,offset)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestFetchCityWithoutPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","city_name","city_desc","city_photos","province_id"}).
		AddRow(mockCity[0].Id, mockCity[0].CreatedBy,mockCity[0].CreatedDate,mockCity[0].ModifiedBy,
			mockCity[0].ModifiedDate,mockCity[0].DeletedBy,mockCity[0].DeletedDate,mockCity[0].IsDeleted,
			mockCity[0].IsActive,mockCity[0].CityName,mockCity[0].CityDesc,mockCity[0].CityPhotos,mockCity[0].ProvinceId).
		AddRow(mockCity[1].Id, mockCity[1].CreatedBy,mockCity[1].CreatedDate,mockCity[1].ModifiedBy,
			mockCity[1].ModifiedDate,mockCity[1].DeletedBy,mockCity[1].DeletedDate,mockCity[1].IsDeleted,
			mockCity[1].IsActive,mockCity[1].CityName,mockCity[1].CityDesc,mockCity[1].CityPhotos,mockCity[1].ProvinceId)

	query := `SELECT \*\ FROM cities where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.FetchCity(context.TODO(), 0,0)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestFetchCityWithPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","city_name","city_desc","city_photos","province_id"}).
		AddRow(mockCity[0].Id, mockCity[0].CreatedBy,mockCity[0].CreatedDate,mockCity[0].ModifiedBy,
			mockCity[0].ModifiedDate,mockCity[0].DeletedBy,mockCity[0].DeletedDate,mockCity[0].IsDeleted,
			mockCity[0].IsActive,mockCity[0].CityName,mockCity[0].CityDesc,mockCity[0].CityPhotos,mockCity[0].ProvinceId).
		AddRow(mockCity[1].Id, mockCity[1].CreatedBy,mockCity[1].CreatedDate,mockCity[1].ModifiedBy,
			mockCity[1].ModifiedDate,mockCity[1].DeletedBy,mockCity[1].DeletedDate,mockCity[1].IsDeleted,
			mockCity[1].CityPhotos,mockCity[1].CityName,mockCity[1].CityDesc,mockCity[1].CityPhotos,mockCity[1].ProvinceId)

	query := `SELECT \*\ FROM cities where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.FetchCity(context.TODO(), limit,offset)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
	//assert.Len(t, anArticle, 2)
}
func TestFetchCityWithoutPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","city_name","city_desc","city_photos","province_id"}).
		AddRow(mockCity[0].Id, mockCity[0].CreatedBy,mockCity[0].CreatedDate,mockCity[0].ModifiedBy,
			mockCity[0].ModifiedDate,mockCity[0].DeletedBy,mockCity[0].DeletedDate,mockCity[0].IsDeleted,
			mockCity[0].IsActive,mockCity[0].CityName,mockCity[0].CityDesc,mockCity[0].CityPhotos,mockCity[0].ProvinceId).
		AddRow(mockCity[1].Id, mockCity[1].CreatedBy,mockCity[1].CreatedDate,mockCity[1].ModifiedBy,
			mockCity[1].ModifiedDate,mockCity[1].DeletedBy,mockCity[1].DeletedDate,mockCity[1].IsDeleted,
			mockCity[1].CityPhotos,mockCity[1].CityName,mockCity[1].CityDesc,mockCity[1].CityPhotos,mockCity[1].ProvinceId)

	query := `SELECT \*\ FROM cities where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.FetchCity(context.TODO(), 0,0)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
}
func TestGetCityByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","city_name","city_desc","city_photos","province_id"}).
		AddRow(mockCity[0].Id, mockCity[0].CreatedBy,mockCity[0].CreatedDate,mockCity[0].ModifiedBy,
			mockCity[0].ModifiedDate,mockCity[0].DeletedBy,mockCity[0].DeletedDate,mockCity[0].IsDeleted,
			mockCity[0].IsActive,mockCity[0].CityName,mockCity[0].CityDesc,mockCity[0].CityPhotos,mockCity[0].ProvinceId).
		AddRow(mockCity[1].Id, mockCity[1].CreatedBy,mockCity[1].CreatedDate,mockCity[1].ModifiedBy,
			mockCity[1].ModifiedDate,mockCity[1].DeletedBy,mockCity[1].DeletedDate,mockCity[1].IsDeleted,
			mockCity[1].IsActive,mockCity[1].CityName,mockCity[1].CityDesc,mockCity[1].CityPhotos,mockCity[1].ProvinceId)

	query := `SELECT \*\ FROM cities WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	num := 1
	anArticle, err := a.GetCityByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByCityIDNotfound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","city_name","city_desc","city_photos","province_id"})

	query := `SELECT \*\ FROM cities WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	num := 4
	anArticle, err := a.GetCityByID(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByCityIDErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","city_name","city_desc","city_photos","province_id"}).
		AddRow(mockCity[0].Id, mockCity[0].CreatedBy,mockCity[0].CreatedDate,mockCity[0].ModifiedBy,
			mockCity[0].ModifiedDate,mockCity[0].DeletedBy,mockCity[0].DeletedDate,mockCity[0].IsDeleted,
			mockCity[0].CityPhotos,mockCity[0].CityName,mockCity[0].CityDesc,mockCity[0].CityPhotos,mockCity[0].ProvinceId)

	query := `SELECT \*\ FROM cities WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	num := 1
	anArticle, err := a.GetCityByID(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestDeleteCity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE cities SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.DeleteCity(context.TODO(), id,deletedBy)
	assert.NoError(t, err)
}
func TestDeleteCityErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE cities SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.DeleteCity(context.TODO(), id,deletedBy)
	assert.Error(t, err)
}
func TestInsertCity(t *testing.T) {

	a := mockCity[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT cities SET created_by=\\?,created_date=\\?,modified_by=\\?,modified_date=\\?,deleted_by=\\?,deleted_date=\\?,is_deleted=\\?,is_active=\\?,city_name=\\?,city_desc=\\?,city_photos=\\?,province_id=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.CityName,
		a.CityDesc,a.CityPhotos,a.ProvinceId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := cpcRepo.NewcpcRepository(db)

	id, err := i.InsertCity(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertCityErrorExec(t *testing.T) {

	a := mockCity[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT cities SET created_by=\\?,created_date=\\?,modified_by=\\?,modified_date=\\?,deleted_by=\\?,deleted_date=\\?,is_deleted=\\?,is_active=\\?,city_name=\\?,city_desc=\\?,city_photos=\\?,province_id=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.CityName,
		a.CityDesc,a.CityPhotos,a.ProvinceId,a.ProvinceId).WillReturnResult(sqlmock.NewResult(1, 1))


	i := cpcRepo.NewcpcRepository(db)

	_, err = i.InsertCity(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdateCity(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockCity[0]
	ar.ModifiedBy = &modifyBy
	ar.ModifiedDate = &now
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE cities set modified_by=\\?, modified_date=\\? ,city_name=\\?,city_desc=\\?,city_photos=\\?,province_id WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.CityName, ar.CityDesc,ar.CityPhotos,ar.ProvinceId, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.UpdateCity(context.TODO(), &ar)
	assert.NoError(t, err)
	assert.Nil(t,err)
}
func TestUpdateCityErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockCity[0]
	ar.ModifiedBy = &modifyBy
	ar.ModifiedDate = &now
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE cities set modified_by=\\?, modified_date=\\? ,city_name=\\?,city_desc=\\?,city_photos=\\?,province_id WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.CityName, ar.CityDesc,ar.CityPhotos,ar.ProvinceId, ar.Id,ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.UpdateCity(context.TODO(), &ar)
	assert.Error(t, err)
}

func TestCountProvince(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(len(mockProvince))

	query := `SELECT count\(\*\) AS count FROM provinces WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	res, err := a.GetCountProvince(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, res, 2, "")
}
func TestCountProvinceErrorFetch(t *testing.T) {
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

	query := `SELECT count\(\*\) AS count FROM provinces WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	_, err = a.GetCountProvince(context.TODO())
	assert.Error(t, err)
}
func TestFetchProvinceWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","province_name","country_id","province_name_transportation"}).
		AddRow(mockProvince[0].Id, mockProvince[0].CreatedBy,mockProvince[0].CreatedDate,mockProvince[0].ModifiedBy,
		mockProvince[0].ModifiedDate,mockProvince[0].DeletedBy,mockProvince[0].DeletedDate,mockProvince[0].IsDeleted,
		mockProvince[0].IsActive,mockProvince[0].ProvinceName,mockProvince[0].CountryId,mockProvince[0].ProvinceNameTransportation).
		AddRow(mockProvince[1].Id, mockProvince[1].CreatedBy,mockProvince[1].CreatedDate,mockProvince[1].ModifiedBy,
		mockProvince[1].ModifiedDate,mockProvince[1].DeletedBy,mockProvince[1].DeletedDate,mockProvince[1].IsDeleted,
		mockProvince[1].IsActive,mockProvince[1].ProvinceName,mockProvince[1].CountryId,mockProvince[1].ProvinceNameTransportation)

	query := `SELECT \*\ FROM provinces where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.FetchProvince(context.TODO(), limit,offset)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestFetchProvinceWithoutPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","province_name","country_id","province_name_transportation"}).
		AddRow(mockProvince[0].Id, mockProvince[0].CreatedBy,mockProvince[0].CreatedDate,mockProvince[0].ModifiedBy,
			mockProvince[0].ModifiedDate,mockProvince[0].DeletedBy,mockProvince[0].DeletedDate,mockProvince[0].IsDeleted,
			mockProvince[0].IsActive,mockProvince[0].ProvinceName,mockProvince[0].CountryId,mockProvince[0].ProvinceNameTransportation).
		AddRow(mockProvince[1].Id, mockProvince[1].CreatedBy,mockProvince[1].CreatedDate,mockProvince[1].ModifiedBy,
			mockProvince[1].ModifiedDate,mockProvince[1].DeletedBy,mockProvince[1].DeletedDate,mockProvince[1].IsDeleted,
			mockProvince[1].IsActive,mockProvince[1].ProvinceName,mockProvince[1].CountryId,mockProvince[1].ProvinceNameTransportation)

	query := `SELECT \*\ FROM provinces where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.FetchProvince(context.TODO(), 0,0)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestFetchProvinceWithPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","province_name","country_id","province_name_transportation"}).
		AddRow(mockProvince[0].Id, mockProvince[0].CreatedBy,mockProvince[0].CreatedDate,mockProvince[0].ModifiedBy,
			mockProvince[0].ModifiedDate,mockProvince[0].DeletedBy,mockProvince[0].DeletedDate,mockProvince[0].IsDeleted,
			mockProvince[0].IsActive,mockProvince[0].ProvinceName,mockProvince[0].CountryId,mockProvince[0].ProvinceNameTransportation).
		AddRow(mockProvince[1].Id, mockProvince[1].CreatedBy,mockProvince[1].CreatedDate,mockProvince[1].ModifiedBy,
			mockProvince[1].ModifiedDate,mockProvince[1].DeletedBy,mockProvince[1].DeletedDate,mockProvince[1].IsDeleted,
			mockProvince[1].ProvinceName,mockProvince[1].ProvinceName,mockProvince[1].CountryId,mockProvince[1].ProvinceNameTransportation)

	query := `SELECT \*\ FROM provinces where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.FetchProvince(context.TODO(), limit,offset)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
	//assert.Len(t, anArticle, 2)
}
func TestFetchProvinceWithoutPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","province_name","country_id","province_name_transportation"}).
		AddRow(mockProvince[0].Id, mockProvince[0].CreatedBy,mockProvince[0].CreatedDate,mockProvince[0].ModifiedBy,
			mockProvince[0].ModifiedDate,mockProvince[0].DeletedBy,mockProvince[0].DeletedDate,mockProvince[0].IsDeleted,
			mockProvince[0].IsActive,mockProvince[0].ProvinceName,mockProvince[0].CountryId,mockProvince[0].ProvinceNameTransportation).
		AddRow(mockProvince[1].Id, mockProvince[1].CreatedBy,mockProvince[1].CreatedDate,mockProvince[1].ModifiedBy,
			mockProvince[1].ModifiedDate,mockProvince[1].DeletedBy,mockProvince[1].DeletedDate,mockProvince[1].IsDeleted,
			mockProvince[1].ProvinceName,mockProvince[1].ProvinceName,mockProvince[1].CountryId,mockProvince[1].ProvinceNameTransportation)

	query := `SELECT \*\ FROM provinces where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.FetchProvince(context.TODO(), 0,0)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
}
func TestGetProvinceByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","province_name","country_id","province_name_transportation"}).
		AddRow(mockProvince[0].Id, mockProvince[0].CreatedBy,mockProvince[0].CreatedDate,mockProvince[0].ModifiedBy,
			mockProvince[0].ModifiedDate,mockProvince[0].DeletedBy,mockProvince[0].DeletedDate,mockProvince[0].IsDeleted,
			mockProvince[0].IsActive,mockProvince[0].ProvinceName,mockProvince[0].CountryId,mockProvince[0].ProvinceNameTransportation).
		AddRow(mockProvince[1].Id, mockProvince[1].CreatedBy,mockProvince[1].CreatedDate,mockProvince[1].ModifiedBy,
			mockProvince[1].ModifiedDate,mockProvince[1].DeletedBy,mockProvince[1].DeletedDate,mockProvince[1].IsDeleted,
			mockProvince[1].IsActive,mockProvince[1].ProvinceName,mockProvince[1].CountryId,mockProvince[1].ProvinceNameTransportation)

	query := `SELECT \*\ FROM provinces WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	num := 1
	anArticle, err := a.GetProvinceByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByProvinceIDNotfound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","province_name","country_id","province_name_transportation"})

	query := `SELECT \*\ FROM provinces WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	num := 4
	anArticle, err := a.GetProvinceByID(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByProvinceIDErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","province_name","country_id","province_name_transportation"}).
		AddRow(mockProvince[0].Id, mockProvince[0].CreatedBy,mockProvince[0].CreatedDate,mockProvince[0].ModifiedBy,
			mockProvince[0].ModifiedDate,mockProvince[0].DeletedBy,mockProvince[0].DeletedDate,mockProvince[0].IsDeleted,
			mockProvince[0].ProvinceName,mockProvince[0].ProvinceName,mockProvince[0].CountryId,mockProvince[0].ProvinceNameTransportation)


	query := `SELECT \*\ FROM provinces WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	num := 1
	anArticle, err := a.GetProvinceByID(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestDeleteProvince(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE provinces SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.DeleteProvince(context.TODO(), id,deletedBy)
	assert.NoError(t, err)
}
func TestDeleteProvinceErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE provinces SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.DeleteCity(context.TODO(), id,deletedBy)
	assert.Error(t, err)
}
func TestInsertProvince(t *testing.T) {

	a := mockProvince[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT provinces SET created_by=\\?,created_date=\\?,modified_by=\\?,modified_date=\\?,deleted_by=\\?,deleted_date=\\?,is_deleted=\\?,is_active=\\?,province_name=\\?,country_id=\\?,province_name_transportation=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ProvinceName,
		a.CountryId,a.ProvinceNameTransportation).WillReturnResult(sqlmock.NewResult(1, 1))

	i := cpcRepo.NewcpcRepository(db)

	id, err := i.InsertProvince(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertProvinceErrorExec(t *testing.T) {

	a := mockProvince[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT provinces SET created_by=\\?,created_date=\\?,modified_by=\\?,modified_date=\\?,deleted_by=\\?,deleted_date=\\?,is_deleted=\\?,is_active=\\?,province_name=\\?,country_id=\\?,province_name_transportation=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ProvinceName,
		a.CountryId,a.ProvinceNameTransportation,a.ProvinceNameTransportation).WillReturnResult(sqlmock.NewResult(1, 1))


	i := cpcRepo.NewcpcRepository(db)

	_, err = i.InsertProvince(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdateProvince(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockProvince[0]
	ar.ModifiedBy = &modifyBy
	ar.ModifiedDate = &now
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE provinces set modified_by=\\?, modified_date=\\? ,province_name=\\?,country_id=\\?,province_name_transportation=\\? WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.ProvinceName, ar.CountryId,ar.ProvinceNameTransportation, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.UpdateProvince(context.TODO(), &ar)
	assert.NoError(t, err)
	assert.Nil(t,err)
}
func TestUpdateProvinceErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockProvince[0]
	ar.ModifiedBy = &modifyBy
	ar.ModifiedDate = &now
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE cities set modified_by=\\?, modified_date=\\? ,province_name=\\?,country_id=\\?,province_name_transportation=\\? WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.ProvinceName, ar.CountryId,ar.ProvinceNameTransportation, ar.Id,ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.UpdateProvince(context.TODO(), &ar)
	assert.Error(t, err)
}

func TestCountCountry(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(len(mockProvince))

	query := `SELECT count\(\*\) AS count FROM countries WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	res, err := a.GetCountCountry(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, res, 2, "")
}
func TestCountCountryErrorFetch(t *testing.T) {
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

	query := `SELECT count\(\*\) AS count FROM countries WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	_, err = a.GetCountCountry(context.TODO())
	assert.Error(t, err)
}
func TestFetchCountryWithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","country_name","iso","name","nice_name","iso3","num_code","phone_code"}).
		AddRow(mockCountry[0].Id, mockCountry[0].CreatedBy,mockCountry[0].CreatedDate,mockCountry[0].ModifiedBy,
		mockCountry[0].ModifiedDate,mockCountry[0].DeletedBy,mockCountry[0].DeletedDate,mockCountry[0].IsDeleted,
		mockCountry[0].IsActive,mockCountry[0].CountryName,mockCountry[0].Iso,mockCountry[0].Name,mockCountry[0].NiceName,
		mockCountry[0].Iso3,mockCountry[0].NumCode,mockCountry[0].PhoneCode).
		AddRow(mockCountry[1].Id, mockCountry[1].CreatedBy,mockCountry[1].CreatedDate,mockCountry[1].ModifiedBy,
			mockCountry[1].ModifiedDate,mockCountry[1].DeletedBy,mockCountry[1].DeletedDate,mockCountry[1].IsDeleted,
			mockCountry[1].IsActive,mockCountry[1].CountryName,mockCountry[1].Iso,mockCountry[1].Name,mockCountry[1].NiceName,
			mockCountry[1].Iso3,mockCountry[1].NumCode,mockCountry[1].PhoneCode)

	query := `SELECT \*\ FROM countries where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.FetchCountry(context.TODO(), limit,offset)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestFetchCountryWithoutPagination(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","country_name","iso","name","nice_name","iso3","num_code","phone_code"}).
		AddRow(mockCountry[0].Id, mockCountry[0].CreatedBy,mockCountry[0].CreatedDate,mockCountry[0].ModifiedBy,
			mockCountry[0].ModifiedDate,mockCountry[0].DeletedBy,mockCountry[0].DeletedDate,mockCountry[0].IsDeleted,
			mockCountry[0].IsActive,mockCountry[0].CountryName,mockCountry[0].Iso,mockCountry[0].Name,mockCountry[0].NiceName,
			mockCountry[0].Iso3,mockCountry[0].NumCode,mockCountry[0].PhoneCode).
		AddRow(mockCountry[1].Id, mockCountry[1].CreatedBy,mockCountry[1].CreatedDate,mockCountry[1].ModifiedBy,
			mockCountry[1].ModifiedDate,mockCountry[1].DeletedBy,mockCountry[1].DeletedDate,mockCountry[1].IsDeleted,
			mockCountry[1].IsActive,mockCountry[1].CountryName,mockCountry[1].Iso,mockCountry[1].Name,mockCountry[1].NiceName,
			mockCountry[1].Iso3,mockCountry[1].NumCode,mockCountry[1].PhoneCode)

	query := `SELECT \*\ FROM countries where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.FetchCountry(context.TODO(), 0,0)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestFetchCountryWithPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","country_name","iso","name","nice_name","iso3","num_code","phone_code"}).
		AddRow(mockCountry[0].Id, mockCountry[0].CreatedBy,mockCountry[0].CreatedDate,mockCountry[0].ModifiedBy,
			mockCountry[0].ModifiedDate,mockCountry[0].DeletedBy,mockCountry[0].DeletedDate,mockCountry[0].IsDeleted,
			mockCountry[0].IsActive,mockCountry[0].CountryName,mockCountry[0].Iso,mockCountry[0].Name,mockCountry[0].NiceName,
			mockCountry[0].Iso3,mockCountry[0].NumCode,mockCountry[0].PhoneCode).
		AddRow(mockCountry[1].Id, mockCountry[1].CreatedBy,mockCountry[1].CreatedDate,mockCountry[1].ModifiedBy,
			mockCountry[1].ModifiedDate,mockCountry[1].DeletedBy,mockCountry[1].DeletedDate,mockCountry[1].IsDeleted,
			mockCountry[1].Name,mockCountry[1].CountryName,mockCountry[1].Iso,mockCountry[1].Name,mockCountry[1].NiceName,
			mockCountry[1].Iso3,mockCountry[1].NumCode,mockCountry[1].PhoneCode)

	query := `SELECT \*\ FROM countries where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.FetchCountry(context.TODO(), limit,offset)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
	//assert.Len(t, anArticle, 2)
}
func TestFetchCountryWithoutPaginationErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","country_name","iso","name","nice_name","iso3","num_code","phone_code"}).
		AddRow(mockCountry[0].Id, mockCountry[0].CreatedBy,mockCountry[0].CreatedDate,mockCountry[0].ModifiedBy,
			mockCountry[0].ModifiedDate,mockCountry[0].DeletedBy,mockCountry[0].DeletedDate,mockCountry[0].IsDeleted,
			mockCountry[0].IsActive,mockCountry[0].CountryName,mockCountry[0].Iso,mockCountry[0].Name,mockCountry[0].NiceName,
			mockCountry[0].Iso3,mockCountry[0].NumCode,mockCountry[0].PhoneCode).
		AddRow(mockCountry[1].Id, mockCountry[1].CreatedBy,mockCountry[1].CreatedDate,mockCountry[1].ModifiedBy,
			mockCountry[1].ModifiedDate,mockCountry[1].DeletedBy,mockCountry[1].DeletedDate,mockCountry[1].IsDeleted,
			mockCountry[1].Name,mockCountry[1].CountryName,mockCountry[1].Iso,mockCountry[1].Name,mockCountry[1].NiceName,
			mockCountry[1].Iso3,mockCountry[1].NumCode,mockCountry[1].PhoneCode)

	query := `SELECT \*\ FROM countries where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.FetchCountry(context.TODO(), 0,0)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
}
func TestGetCountryByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","country_name","iso","name","nice_name","iso3","num_code","phone_code"}).
		AddRow(mockCountry[0].Id, mockCountry[0].CreatedBy,mockCountry[0].CreatedDate,mockCountry[0].ModifiedBy,
			mockCountry[0].ModifiedDate,mockCountry[0].DeletedBy,mockCountry[0].DeletedDate,mockCountry[0].IsDeleted,
			mockCountry[0].IsActive,mockCountry[0].CountryName,mockCountry[0].Iso,mockCountry[0].Name,mockCountry[0].NiceName,
			mockCountry[0].Iso3,mockCountry[0].NumCode,mockCountry[0].PhoneCode).
		AddRow(mockCountry[1].Id, mockCountry[1].CreatedBy,mockCountry[1].CreatedDate,mockCountry[1].ModifiedBy,
			mockCountry[1].ModifiedDate,mockCountry[1].DeletedBy,mockCountry[1].DeletedDate,mockCountry[1].IsDeleted,
			mockCountry[1].IsActive,mockCountry[1].CountryName,mockCountry[1].Iso,mockCountry[1].Name,mockCountry[1].NiceName,
			mockCountry[1].Iso3,mockCountry[1].NumCode,mockCountry[1].PhoneCode)

	query := `SELECT \*\ FROM countries WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	num := 1
	anArticle, err := a.GetCountryByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByCountryIDNotfound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","country_name","iso","name","nice_name","iso3","num_code","phone_code"})

	query := `SELECT \*\ FROM countries WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	num := 4
	anArticle, err := a.GetCountryByID(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByCountryIDErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","country_name","iso","name","nice_name","iso3","num_code","phone_code"}).
		AddRow(mockCountry[0].Id, mockCountry[0].CreatedBy,mockCountry[0].CreatedDate,mockCountry[0].ModifiedBy,
			mockCountry[0].ModifiedDate,mockCountry[0].DeletedBy,mockCountry[0].DeletedDate,mockCountry[0].IsDeleted,
			mockCountry[0].Name,mockCountry[0].CountryName,mockCountry[0].Iso,mockCountry[0].Name,mockCountry[0].NiceName,
			mockCountry[0].Iso3,mockCountry[0].NumCode,mockCountry[0].PhoneCode)

	query := `SELECT \*\ FROM countries WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := cpcRepo.NewcpcRepository(db)

	num := 1
	anArticle, err := a.GetCountryByID(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestDeleteCountry(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE countries SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.DeleteCountry(context.TODO(), id,deletedBy)
	assert.NoError(t, err)
}
func TestDeleteCountryErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE countries SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.DeleteCountry(context.TODO(), id,deletedBy)
	assert.Error(t, err)
}
func TestInsertCountry(t *testing.T) {

	a := mockCountry[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT countries SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? ,deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , country_name=\\?,iso=\\?,name=\\?,nice_name=\\?, iso3=\\?,num_code=\\?,phone_code=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.CountryName,a.Iso,a.Name,
		a.NiceName,a.Iso3,a.NumCode,a.PhoneCode).WillReturnResult(sqlmock.NewResult(1, 1))

	i := cpcRepo.NewcpcRepository(db)

	id, err := i.InsertCountry(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertCountryErrorExec(t *testing.T) {

	a := mockCountry[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT countries SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? ,deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , country_name=\\?,iso=\\?,name=\\?,nice_name=\\?, iso3=\\?,num_code=\\?,phone_code=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.CountryName,a.Iso,a.Name,
		a.NiceName,a.Iso3,a.NumCode,a.PhoneCode,a.PhoneCode).WillReturnResult(sqlmock.NewResult(1, 1))


	i := cpcRepo.NewcpcRepository(db)

	_, err = i.InsertCountry(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdateCountry(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockCountry[0]
	ar.ModifiedBy = &modifyBy
	ar.ModifiedDate = &now
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE countries set modified_by=\\?, modified_date=\\? ,country_name=\\?,iso=\\?,name=\\?,nice_name=\\?, iso3=\\?,num_code=\\?,phone_code=\\? WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.CountryName,ar.Iso,ar.Name,
		ar.NiceName,ar.Iso3,ar.NumCode,ar.PhoneCode, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.UpdateCountry(context.TODO(), &ar)
	assert.NoError(t, err)
	assert.Nil(t,err)
}
func TestUpdateCountryErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockCountry[0]
	ar.ModifiedBy = &modifyBy
	ar.ModifiedDate = &now
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "UPDATE countries set modified_by=\\?, modified_date=\\? ,country_name=\\?,iso=\\?,name=\\?,nice_name=\\?, iso3=\\?,num_code=\\?,phone_code=\\? WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.CountryName,ar.Iso,ar.Name,
		ar.NiceName,ar.Iso3,ar.NumCode,ar.PhoneCode, ar.Id,ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := cpcRepo.NewcpcRepository(db)

	err = a.UpdateCountry(context.TODO(), &ar)
	assert.Error(t, err)
}
