package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	HarborsRepo "github.com/service/Harbors/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var(
	harborsType = 1
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
func TestGetAllWithJoinCPC(t *testing.T) {
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

	query := `SELECT                                         \*\                                               FROM harbors WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

	anArticle, err := a.List(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetAllWithJoinCPCError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockHarbors := []models.Harbors{
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 1",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 2",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].HarborsName, mockHarbors[0].HarborsName, mockHarbors[0].HarborsIcon).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].HarborsName, mockHarbors[1].HarborsName, mockHarbors[1].HarborsIcon)

	query := `SELECT                                         \*\                                               FROM Harborss WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

	anArticle, err := a.List(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
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
	mockHarbors := []models.Harbors{
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 1",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 2",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].IsActive, mockHarbors[0].HarborsName, mockHarbors[0].HarborsIcon).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].IsActive, mockHarbors[1].HarborsName, mockHarbors[1].HarborsIcon)

	query := `SELECT \*\ FROM Harborss where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

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
	mockHarbors := []models.Harbors{
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 1",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 2",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].IsActive, mockHarbors[0].HarborsName, mockHarbors[0].HarborsIcon).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].IsActive, mockHarbors[1].HarborsName, mockHarbors[1].HarborsIcon)

	query := `SELECT \*\ FROM Harborss where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

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
	mockHarbors := []models.Harbors{
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 1",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 2",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].HarborsIcon, mockHarbors[0].HarborsName, mockHarbors[0].HarborsIcon).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].HarborsIcon, mockHarbors[1].HarborsName, mockHarbors[1].HarborsIcon)

	query := `SELECT \*\ FROM Harborss where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

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
	mockHarbors := []models.Harbors{
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 1",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 2",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].HarborsIcon, mockHarbors[0].HarborsName, mockHarbors[0].HarborsIcon).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].HarborsIcon, mockHarbors[1].HarborsName, mockHarbors[1].HarborsIcon)

	query := `SELECT \*\ FROM Harborss where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

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

	mockHarbors := []models.Harbors{
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 1",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
		models.Harbors{
			Id:           2,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 2",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].IsActive, mockHarbors[0].HarborsName, mockHarbors[0].HarborsIcon)

	query := `SELECT \*\ FROM Harborss WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

	num := 1
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
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"})

	query := `SELECT \*\ FROM Harborss WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

	num := 4
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

	mockHarbors := []models.Harbors{
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 1",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
		models.Harbors{
			Id:           2,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 2",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].HarborsName, mockHarbors[0].HarborsName, mockHarbors[0].HarborsIcon)

	query := `SELECT \*\ FROM Harborss WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

	num := 1
	anArticle, err := a.GetById(context.TODO(), num)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	mockHarbors := []models.Harbors{
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 1",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 2",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].IsActive, mockHarbors[0].HarborsName, mockHarbors[0].HarborsIcon).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].IsActive, mockHarbors[1].HarborsName, mockHarbors[1].HarborsIcon)

	query := `SELECT \*\ FROM Harborss WHERE Harbors_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

	HarborsName := "Test Harbors 2"
	anArticle, err := a.GetByName(context.TODO(), HarborsName)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
func TestGetByNameNotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"})

	query := `SELECT \*\ FROM Harborss WHERE Harbors_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

	HarborsName := "Test Harbors 2"
	anArticle, err := a.GetByName(context.TODO(), HarborsName)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetByNameErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	mockHarbors := []models.Harbors{
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 1",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
		models.Harbors{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			HarborsName:  "Test Harbors 2",
			HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "Harbors_name", "Harbors_icon"}).
		AddRow(mockHarbors[0].Id, mockHarbors[0].CreatedBy, mockHarbors[0].CreatedDate, mockHarbors[0].ModifiedBy,
			mockHarbors[0].ModifiedDate, mockHarbors[0].DeletedBy, mockHarbors[0].DeletedDate, mockHarbors[0].IsDeleted,
			mockHarbors[0].HarborsName, mockHarbors[0].HarborsName, mockHarbors[0].HarborsIcon).
		AddRow(mockHarbors[1].Id, mockHarbors[1].CreatedBy, mockHarbors[1].CreatedDate, mockHarbors[1].ModifiedBy,
			mockHarbors[1].ModifiedDate, mockHarbors[1].DeletedBy, mockHarbors[1].DeletedDate, mockHarbors[1].IsDeleted,
			mockHarbors[1].HarborsName, mockHarbors[1].HarborsName, mockHarbors[1].HarborsIcon)

	query := `SELECT \*\ FROM Harborss WHERE Harbors_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := HarborsRepo.NewHarborsRepository(db)

	HarborsName := "Test Harbors 2"
	anArticle, err := a.GetByName(context.TODO(), HarborsName)
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

	query := "UPDATE Harborss SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := HarborsRepo.NewHarborsRepository(db)

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

	query := "UPDATE Harborss SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id, id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := HarborsRepo.NewHarborsRepository(db)

	err = a.Delete(context.TODO(), id, deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {
	user := "test"
	now := time.Now()
	a := models.Harbors{
		Id:           1,
		CreatedBy:    user,
		CreatedDate:  now,
		ModifiedBy:   &user,
		ModifiedDate: &now,
		DeletedBy:    &user,
		DeletedDate:  &now,
		IsDeleted:    0,
		IsActive:     0,
		HarborsName:  "test Harbors 1",
		HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT Harborss SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , Harbors_name=\\?,  				Harbors_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.HarborsName,
		a.HarborsIcon).WillReturnResult(sqlmock.NewResult(1, 1))

	i := HarborsRepo.NewHarborsRepository(db)

	id, err := i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	user := "test"
	now := time.Now()
	a := models.Harbors{
		Id:           1,
		CreatedBy:    user,
		CreatedDate:  now,
		ModifiedBy:   &user,
		ModifiedDate: &now,
		DeletedBy:    &user,
		DeletedDate:  &now,
		IsDeleted:    0,
		IsActive:     0,
		HarborsName:  "test Harbors 1",
		HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT Harborss SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , Harbors_name=\\?,  				Harbors_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.HarborsName,
		a.HarborsIcon, a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := HarborsRepo.NewHarborsRepository(db)

	_, err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := models.Harbors{
		Id:           1,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &modifyBy,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		HarborsName:  "test Harbors 1",
		HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE Harborss set modified_by=\?, modified_date=\? ,Harbors_name=\?,Harbors_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.HarborsName, ar.HarborsIcon, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := HarborsRepo.NewHarborsRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.NoError(t, err)
	assert.Nil(t, err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := models.Harbors{
		Id:           1,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &modifyBy,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		HarborsName:  "test Harbors 1",
		HarborsIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Harbors/8941695193938718058.jpg",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE Harborss set modified_by=\?, modified_date=\? ,Harbors_name=\?,Harbors_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.HarborsName, ar.HarborsIcon, ar.Id, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := HarborsRepo.NewHarborsRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.Error(t, err)
}
