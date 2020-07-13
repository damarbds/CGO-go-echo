package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	ExclusionServiceRepo "github.com/service/exclusion_service/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var (
	mockExclusionService = []models.ExclusionService{
		models.ExclusionService{
			Id:                   1,
			CreatedBy:            "test",
			CreatedDate:          time.Now(),
			ModifiedBy:           nil,
			ModifiedDate:         nil,
			DeletedBy:            nil,
			DeletedDate:          nil,
			IsDeleted:            0,
			IsActive:             1,
			ExclusionServiceName: "Test ExclusionService 1",
			ExclusionServiceType: 1,
		},
		models.ExclusionService{
			Id:                   1,
			CreatedBy:            "test",
			CreatedDate:          time.Now(),
			ModifiedBy:           nil,
			ModifiedDate:         nil,
			DeletedBy:            nil,
			DeletedDate:          nil,
			IsDeleted:            0,
			IsActive:             1,
			ExclusionServiceName: "Test ExclusionService 2",
			ExclusionServiceType: 1,
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
		AddRow(len(mockExclusionService))

	query := `SELECT count\(\*\) AS count FROM exclusion_services WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExclusionServiceRepo.NewExclusionServiceRepository(db)

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

	query := `SELECT count\(\*\) AS count FROM exclusion_services WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExclusionServiceRepo.NewExclusionServiceRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "exclusion_service_name", "exclusion_service_type"}).
		AddRow(mockExclusionService[0].Id, mockExclusionService[0].CreatedBy, mockExclusionService[0].CreatedDate, mockExclusionService[0].ModifiedBy,
			mockExclusionService[0].ModifiedDate, mockExclusionService[0].DeletedBy, mockExclusionService[0].DeletedDate, mockExclusionService[0].IsDeleted,
			mockExclusionService[0].IsActive, mockExclusionService[0].ExclusionServiceName, mockExclusionService[0].ExclusionServiceType).
		AddRow(mockExclusionService[1].Id, mockExclusionService[1].CreatedBy, mockExclusionService[1].CreatedDate, mockExclusionService[1].ModifiedBy,
			mockExclusionService[1].ModifiedDate, mockExclusionService[1].DeletedBy, mockExclusionService[1].DeletedDate, mockExclusionService[1].IsDeleted,
			mockExclusionService[1].IsActive, mockExclusionService[1].ExclusionServiceName, mockExclusionService[1].ExclusionServiceType)

	query := `SELECT \*\ FROM exclusion_services where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExclusionServiceRepo.NewExclusionServiceRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "exclusion_service_name", "exclusion_service_type"}).
		AddRow(mockExclusionService[0].Id, mockExclusionService[0].CreatedBy, mockExclusionService[0].CreatedDate, mockExclusionService[0].ModifiedBy,
			mockExclusionService[0].ModifiedDate, mockExclusionService[0].DeletedBy, mockExclusionService[0].DeletedDate, mockExclusionService[0].IsDeleted,
			mockExclusionService[0].IsActive, mockExclusionService[0].ExclusionServiceName, mockExclusionService[0].ExclusionServiceType).
		AddRow(mockExclusionService[1].Id, mockExclusionService[1].CreatedBy, mockExclusionService[1].CreatedDate, mockExclusionService[1].ModifiedBy,
			mockExclusionService[1].ModifiedDate, mockExclusionService[1].DeletedBy, mockExclusionService[1].DeletedDate, mockExclusionService[1].IsDeleted,
			mockExclusionService[1].IsActive, mockExclusionService[1].ExclusionServiceName, mockExclusionService[1].ExclusionServiceType)

	query := `SELECT \*\ FROM exclusion_services where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExclusionServiceRepo.NewExclusionServiceRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "exclusion_service_name", "exclusion_service_type"}).
		AddRow(mockExclusionService[0].Id, mockExclusionService[0].CreatedBy, mockExclusionService[0].CreatedDate, mockExclusionService[0].ModifiedBy,
			mockExclusionService[0].ModifiedDate, mockExclusionService[0].DeletedBy, mockExclusionService[0].DeletedDate, mockExclusionService[0].IsDeleted,
			mockExclusionService[0].IsActive, mockExclusionService[0].ExclusionServiceName, mockExclusionService[0].ExclusionServiceType).
		AddRow(mockExclusionService[1].Id, mockExclusionService[1].CreatedBy, mockExclusionService[1].CreatedDate, mockExclusionService[1].ModifiedBy,
			mockExclusionService[1].ModifiedDate, mockExclusionService[1].DeletedBy, mockExclusionService[1].DeletedDate, mockExclusionService[1].IsDeleted,
			mockExclusionService[1].ModifiedDate, mockExclusionService[1].ExclusionServiceName, mockExclusionService[1].ExclusionServiceType)

	query := `SELECT \*\ FROM exclusion_services where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExclusionServiceRepo.NewExclusionServiceRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "exclusion_service_name", "exclusion_service_type"}).
		AddRow(mockExclusionService[0].Id, mockExclusionService[0].CreatedBy, mockExclusionService[0].CreatedDate, mockExclusionService[0].ModifiedBy,
			mockExclusionService[0].ModifiedDate, mockExclusionService[0].DeletedBy, mockExclusionService[0].DeletedDate, mockExclusionService[0].IsDeleted,
			mockExclusionService[0].IsActive, mockExclusionService[0].ExclusionServiceName, mockExclusionService[0].ExclusionServiceType).
		AddRow(mockExclusionService[1].Id, mockExclusionService[1].CreatedBy, mockExclusionService[1].CreatedDate, mockExclusionService[1].ModifiedBy,
			mockExclusionService[1].ModifiedDate, mockExclusionService[1].DeletedBy, mockExclusionService[1].DeletedDate, mockExclusionService[1].IsDeleted,
			mockExclusionService[1].ModifiedDate, mockExclusionService[1].ExclusionServiceName, mockExclusionService[1].ExclusionServiceType)

	query := `SELECT \*\ FROM exclusion_services where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExclusionServiceRepo.NewExclusionServiceRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.Fetch(context.TODO(), 0, 0)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
