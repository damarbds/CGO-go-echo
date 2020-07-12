package repository_test

import (
	"context"
	"github.com/models"
	includeRepo "github.com/service/include/repository"
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
	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(len(mockInclude))

	query := `SELECT count\(\*\) AS count FROM includes WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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

	query := `SELECT count\(\*\) AS count FROM includes WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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
	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IsActive,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].IsActive,mockInclude[1].IncludeName,mockInclude[1].IncludeIcon)

	query := `SELECT                                         \*\                                               FROM includes WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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
	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IncludeName,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].IncludeName,mockInclude[1].IncludeName,mockInclude[1].IncludeIcon)

	query := `SELECT                                         \*\                                               FROM includes WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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
	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IsActive,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].IsActive,mockInclude[1].IncludeName,mockInclude[1].IncludeIcon)

	query := `SELECT \*\ FROM includes where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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
	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IsActive,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].IsActive,mockInclude[1].IncludeName,mockInclude[1].IncludeIcon)

	query := `SELECT \*\ FROM includes where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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
	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IncludeIcon,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].IncludeIcon,mockInclude[1].IncludeName,mockInclude[1].IncludeIcon)

	query := `SELECT \*\ FROM includes where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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
	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IncludeIcon,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].IncludeIcon,mockInclude[1].IncludeName,mockInclude[1].IncludeIcon)

	query := `SELECT \*\ FROM includes where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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

	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           2,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IsActive,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon)

	query := `SELECT \*\ FROM includes WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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
		"deleted_date","is_deleted","is_active","include_name","include_icon"})

	query := `SELECT \*\ FROM includes WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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

	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           2,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IncludeName,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon)

	query := `SELECT \*\ FROM includes WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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

	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IsActive,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].IsActive,mockInclude[1].IncludeName,mockInclude[1].IncludeIcon)

	query := `SELECT \*\ FROM includes WHERE include_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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
		"deleted_date","is_deleted","is_active","include_name","include_icon"})

	query := `SELECT \*\ FROM includes WHERE include_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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

	mockInclude := []models.Include{
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 1",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
		models.Include{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			IncludeName:  "Test Include 2",
			IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
		},
	}
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date","is_deleted","is_active","include_name","include_icon"}).
		AddRow(mockInclude[0].Id, mockInclude[0].CreatedBy,mockInclude[0].CreatedDate,mockInclude[0].ModifiedBy,
			mockInclude[0].ModifiedDate,mockInclude[0].DeletedBy,mockInclude[0].DeletedDate,mockInclude[0].IsDeleted,
			mockInclude[0].IncludeName,mockInclude[0].IncludeName,mockInclude[0].IncludeIcon).
		AddRow(mockInclude[1].Id, mockInclude[1].CreatedBy,mockInclude[1].CreatedDate,mockInclude[1].ModifiedBy,
			mockInclude[1].ModifiedDate,mockInclude[1].DeletedBy,mockInclude[1].DeletedDate,mockInclude[1].IsDeleted,
			mockInclude[1].IncludeName,mockInclude[1].IncludeName,mockInclude[1].IncludeIcon)

	query := `SELECT \*\ FROM includes WHERE include_name = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := includeRepo.NewIncludeRepository(db)

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

	query := "UPDATE includes SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := includeRepo.NewIncludeRepository(db)

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

	query := "UPDATE includes SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0,id,id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := includeRepo.NewIncludeRepository(db)

	err = a.Delete(context.TODO(), id,deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {
	user := "test"
	now := time.Now()
	a := models.Include{
		Id:           1,
		CreatedBy:    user,
		CreatedDate:  now,
		ModifiedBy:   &user,
		ModifiedDate: &now,
		DeletedBy:    &user,
		DeletedDate:  &now,
		IsDeleted:    0,
		IsActive:     0,
		IncludeName:  "test include 1",
		IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT includes SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , include_name=\\?,  				include_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.IncludeName,
		a.IncludeIcon).WillReturnResult(sqlmock.NewResult(1, 1))

	i := includeRepo.NewIncludeRepository(db)

	id, err := i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	user := "test"
	now := time.Now()
	a := models.Include{
		Id:           1,
		CreatedBy:    user,
		CreatedDate:  now,
		ModifiedBy:   &user,
		ModifiedDate: &now,
		DeletedBy:    &user,
		DeletedDate:  &now,
		IsDeleted:    0,
		IsActive:     0,
		IncludeName:  "test include 1",
		IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT includes SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? , include_name=\\?,  				include_icon=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.IncludeName,
		a.IncludeIcon,a.Id).WillReturnResult(sqlmock.NewResult(1, 1))

	i := includeRepo.NewIncludeRepository(db)

	_, err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := models.Include{
		Id:           1,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &modifyBy,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		IncludeName:  "test include 1",
		IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE includes set modified_by=\?, modified_date=\? ,include_name=\?,include_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.IncludeName, ar.IncludeIcon, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := includeRepo.NewIncludeRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.NoError(t, err)
	assert.Nil(t,err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := models.Include{
		Id:           1,
		CreatedBy:    "",
		CreatedDate:  time.Time{},
		ModifiedBy:   &modifyBy,
		ModifiedDate: &now,
		DeletedBy:    nil,
		DeletedDate:  nil,
		IsDeleted:    0,
		IsActive:     0,
		IncludeName:  "test include 1",
		IncludeIcon:  "https://cgostorage.blob.core.windows.net/cgo-storage/Master/Include/8941695193938718058.jpg",
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE includes set modified_by=\?, modified_date=\? ,include_name=\?,include_icon=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.IncludeName, ar.IncludeIcon, ar.Id,ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := includeRepo.NewIncludeRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.Error(t, err)
}
