package repository_test

import (
	"context"
	"github.com/models"
	excludeRepo "github.com/service/exclude/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
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
	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(len(mockExclude))

	query := `SELECT count\(\*\) AS count FROM excludes WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	res, err := a.GetCount(context.TODO())
	assert.NoError(t, err)
	assert.Equal(t, 2, res, "")
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

	query := `SELECT count\(\*\) AS count FROM excludes WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	_, err = a.GetCount(context.TODO())
	assert.Error(t, err)
}
func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
		mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
		mockExclude[0].IsActive,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon).
		AddRow(mockExclude[1].Id, mockExclude[1].CreatedBy,mockExclude[1].CreatedDate,mockExclude[1].ModifiedBy,
		mockExclude[1].ModifiedDate,mockExclude[1].DeletedBy,mockExclude[1].DeletedDate,mockExclude[1].IsDeleted,
		mockExclude[1].IsActive,mockExclude[1].ExcludeName,mockExclude[1].ExcludeIcon)

	query := `SELECT                                         \*\                                               FROM excludes WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	anArticle, err := a.List(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestListError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
			mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
			mockExclude[0].IsActive,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon).
		AddRow(mockExclude[1].Id, mockExclude[1].CreatedBy,mockExclude[1].CreatedDate,mockExclude[1].ModifiedBy,
			mockExclude[1].ModifiedDate,mockExclude[1].DeletedBy,mockExclude[1].DeletedDate,mockExclude[1].IsDeleted,
			mockExclude[1].ExcludeName,mockExclude[1].ExcludeName,mockExclude[1].ExcludeIcon)
	query := `SELECT                                         \*\                                               FROM excludes WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

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
	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
			mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
			mockExclude[0].IsActive,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon).
		AddRow(mockExclude[1].Id, mockExclude[1].CreatedBy,mockExclude[1].CreatedDate,mockExclude[1].ModifiedBy,
			mockExclude[1].ModifiedDate,mockExclude[1].DeletedBy,mockExclude[1].DeletedDate,mockExclude[1].IsDeleted,
			mockExclude[1].IsActive,mockExclude[1].ExcludeName,mockExclude[1].ExcludeIcon)

	query := `SELECT \*\ FROM excludes where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.Fetch(context.TODO(), limit,offset)
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
	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
			mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
			mockExclude[0].IsActive,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon).
		AddRow(mockExclude[1].Id, mockExclude[1].CreatedBy,mockExclude[1].CreatedDate,mockExclude[1].ModifiedBy,
			mockExclude[1].ModifiedDate,mockExclude[1].DeletedBy,mockExclude[1].DeletedDate,mockExclude[1].IsDeleted,
			mockExclude[1].IsActive,mockExclude[1].ExcludeName,mockExclude[1].ExcludeIcon)


	query := `SELECT \*\ FROM excludes where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), 0,0)
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
	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
			mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
			mockExclude[0].ExcludeIcon,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon).
		AddRow(mockExclude[1].Id, mockExclude[1].CreatedBy,mockExclude[1].CreatedDate,mockExclude[1].ModifiedBy,
			mockExclude[1].ModifiedDate,mockExclude[1].DeletedBy,mockExclude[1].DeletedDate,mockExclude[1].IsDeleted,
			mockExclude[1].ExcludeIcon,mockExclude[1].ExcludeName,mockExclude[1].ExcludeIcon)

	query := `SELECT \*\ FROM excludes where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.Fetch(context.TODO(), limit,offset)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
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
	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
			mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
			mockExclude[0].ExcludeIcon,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon).
		AddRow(mockExclude[1].Id, mockExclude[1].CreatedBy,mockExclude[1].CreatedDate,mockExclude[1].ModifiedBy,
			mockExclude[1].ModifiedDate,mockExclude[1].DeletedBy,mockExclude[1].DeletedDate,mockExclude[1].IsDeleted,
			mockExclude[1].ExcludeIcon,mockExclude[1].ExcludeName,mockExclude[1].ExcludeIcon)


	query := `SELECT \*\ FROM excludes where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), 0,0)
	assert.Error(t, err)
	assert.Nil(t,anArticle)
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

	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
			mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
			mockExclude[0].IsActive,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon)

	query := `SELECT \*\ FROM excludes WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

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
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"})

	query := `SELECT \*\ FROM excludes WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

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

	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
			mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
			mockExclude[0].IsActive,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon).
		AddRow(mockExclude[1].Id, mockExclude[1].CreatedBy,mockExclude[1].CreatedDate,mockExclude[1].ModifiedBy,
			mockExclude[1].ModifiedDate,mockExclude[1].DeletedBy,mockExclude[1].DeletedDate,mockExclude[1].IsDeleted,
			mockExclude[1].ExcludeIcon,mockExclude[1].ExcludeName,mockExclude[1].ExcludeIcon)


	query := `SELECT \*\ FROM excludes WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

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

	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
			mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
			mockExclude[0].IsActive,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon).
		AddRow(mockExclude[1].Id, mockExclude[1].CreatedBy,mockExclude[1].CreatedDate,mockExclude[1].ModifiedBy,
			mockExclude[1].ModifiedDate,mockExclude[1].DeletedBy,mockExclude[1].DeletedDate,mockExclude[1].IsDeleted,
			mockExclude[1].IsActive,mockExclude[1].ExcludeName,mockExclude[1].ExcludeIcon)


	query := `SELECT \*\ FROM excludes WHERE exclude_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	includeName := "Test Include 2"
	anArticle, err := a.GetByName(context.TODO(), includeName)
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
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"})


	query := `SELECT \*\ FROM excludes WHERE exclude_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	includeName := "Test Include 2"
	anArticle, err := a.GetByName(context.TODO(), includeName)
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

	mockExclude := []models.Exclude{
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 1",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Exclude{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExcludeName:  "Test Include 2",
			ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","exclude_name","exclude_icon"}).
		AddRow(mockExclude[0].Id, mockExclude[0].CreatedBy,mockExclude[0].CreatedDate,mockExclude[0].ModifiedBy,
			mockExclude[0].ModifiedDate,mockExclude[0].DeletedBy,mockExclude[0].DeletedDate,mockExclude[0].IsDeleted,
			mockExclude[0].IsActive,mockExclude[0].ExcludeName,mockExclude[0].ExcludeIcon).
		AddRow(mockExclude[1].Id, mockExclude[1].CreatedBy,mockExclude[1].CreatedDate,mockExclude[1].ModifiedBy,
			mockExclude[1].ModifiedDate,mockExclude[1].DeletedBy,mockExclude[1].DeletedDate,mockExclude[1].IsDeleted,
			mockExclude[1].ExcludeIcon,mockExclude[1].ExcludeName,mockExclude[1].ExcludeIcon)


	query := `SELECT \*\ FROM excludes WHERE exclude_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := excludeRepo.NewExcludeRepository(db)

	includeName := "Test Include 2"
	anArticle, err := a.GetByName(context.TODO(), includeName)
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

	query := "UPDATE excludes SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := excludeRepo.NewExcludeRepository(db)

	err = a.Delete(context.TODO(), id,deletedBy)
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

	query := "UPDATE excludes SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := excludeRepo.NewExcludeRepository(db)

	err = a.Delete(context.TODO(), id,deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {
	user := "test"
	now := time.Now()
	a := models.Exclude{
		Id:           1,
		CreatedBy:    user,
		CreatedDate:  now,
		ModifiedBy:   &user,
		ModifiedDate: &now,
		DeletedBy:    &user,
		DeletedDate:  &now,
		IsDeleted:    0,
		IsActive:     0,
		ExcludeName:  "test include 1",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT excludes SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , exclude_name=\\?,  				exclude_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExcludeName,
		a.ExcludeIcon).WillReturnResult(sqlmock.NewResult(1, 1))

	i := excludeRepo.NewExcludeRepository(db)

	id, err := i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	user := "test"
	now := time.Now()
	a := models.Exclude{
		Id:           1,
		CreatedBy:    user,
		CreatedDate:  now,
		ModifiedBy:   &user,
		ModifiedDate: &now,
		DeletedBy:    &user,
		DeletedDate:  &now,
		IsDeleted:    0,
		IsActive:     0,
		ExcludeName:  "test include 1",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT excludes SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , exclude_name=\\?,  				exclude_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExcludeName,
		a.ExcludeIcon,a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := excludeRepo.NewExcludeRepository(db)

	_, err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := models.Exclude{
		Id:           1,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &modifyBy,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		ExcludeName:  "test include 1",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE excludes set modified_by=\?, modified_date=\? ,exclude_name=\?,exclude_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.ExcludeName, ar.ExcludeIcon, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := excludeRepo.NewExcludeRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.NoError(t, err)
	assert.Nil(t,err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := models.Exclude{
		Id:           1,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &modifyBy,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		ExcludeName:  "test include 1",
		ExcludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE excludes set modified_by=\?, modified_date=\? ,exclude_name=\?,exclude_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.ExcludeName, ar.ExcludeIcon, ar.Id,ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := excludeRepo.NewExcludeRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.Error(t, err)
}
