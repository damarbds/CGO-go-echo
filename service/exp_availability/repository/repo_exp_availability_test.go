package repository_test
//
//import (
//	"context"
//	"testing"
//	"time"
//
//	"github.com/models"
//	ExpAvailabilityRepo "github.com/service/minimum_booking/repository"
//
//	"github.com/stretchr/testify/assert"
//	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
//)
//
//var (
//	mockExpAvailability = []models.ExpAvailability{
//		models.ExpAvailability{
//			Id:                   "asdjaskldjsaldj",
//			CreatedBy:            "test",
//			CreatedDate:          time.Now(),
//			ModifiedBy:           nil,
//			ModifiedDate:         nil,
//			DeletedBy:            nil,
//			DeletedDate:          nil,
//			IsDeleted:            0,
//			IsActive:             1,
//			ExpAvailabilityMonth: "April",
//			ExpAvailabilityDate:  "05-05-2020",
//			ExpAvailabilityYear:  2020,
//			ExpId:                "qrqeqweqweqw",
//		},
//		models.ExpAvailability{
//			Id:                   "asdjaskldjsaldj",
//			CreatedBy:            "test",
//			CreatedDate:          time.Now(),
//			ModifiedBy:           nil,
//			ModifiedDate:         nil,
//			DeletedBy:            nil,
//			DeletedDate:          nil,
//			IsDeleted:            0,
//			IsActive:             1,
//			ExpAvailabilityMonth: "May",
//			ExpAvailabilityDate:  "05-05-2020",
//			ExpAvailabilityYear:  2020,
//			ExpId:                "qrqeqweqweqw",
//		},
//	}
//)
//
//func TestCount(t *testing.T) {
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
//	rows := sqlmock.NewRows([]string{"count"}).
//		AddRow(len(mockExpAvailability))
//
//	query := `SELECT count\(\*\) AS count FROM minimum_bookings WHERE is_deleted = 0 and is_active = 1`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := ExpAvailabilityRepo.NewExpAvailabilityRepository(db)
//
//	res, err := a.GetCount(context.TODO())
//	assert.NoError(t, err)
//	assert.Equal(t, res, 2, "")
//}
//func TestCountFetch(t *testing.T) {
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
//	rows := sqlmock.NewRows([]string{"count"}).
//		AddRow("test")
//
//	query := `SELECT count\(\*\) AS count FROM minimum_bookings WHERE is_deleted = 0 and is_active = 1`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := ExpAvailabilityRepo.NewExpAvailabilityRepository(db)
//
//	_, err = a.GetCount(context.TODO())
//	assert.Error(t, err)
//}
//func TestFetchWithPagination(t *testing.T) {
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
//		"deleted_date", "is_deleted", "is_active", "minimum_booking_desc", "minimum_booking_amount"}).
//		AddRow(mockExpAvailability[0].Id, mockExpAvailability[0].CreatedBy, mockExpAvailability[0].CreatedDate, mockExpAvailability[0].ModifiedBy,
//			mockExpAvailability[0].ModifiedDate, mockExpAvailability[0].DeletedBy, mockExpAvailability[0].DeletedDate, mockExpAvailability[0].IsDeleted,
//			mockExpAvailability[0].IsActive, mockExpAvailability[0].ExpAvailabilityDesc, mockExpAvailability[0].ExpAvailabilityAmount).
//		AddRow(mockExpAvailability[1].Id, mockExpAvailability[1].CreatedBy, mockExpAvailability[1].CreatedDate, mockExpAvailability[1].ModifiedBy,
//			mockExpAvailability[1].ModifiedDate, mockExpAvailability[1].DeletedBy, mockExpAvailability[1].DeletedDate, mockExpAvailability[1].IsDeleted,
//			mockExpAvailability[1].IsActive, mockExpAvailability[1].ExpAvailabilityDesc, mockExpAvailability[1].ExpAvailabilityAmount)
//
//	query := `SELECT \*\ FROM minimum_bookings where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := ExpAvailabilityRepo.NewExpAvailabilityRepository(db)
//
//	limit := 10
//	offset := 0
//	anArticle, err := a.GetAll(context.TODO(), limit, offset)
//	//assert.NotEmpty(t, nextCursor)
//	assert.NoError(t, err)
//	assert.Len(t, anArticle, 2)
//}
//func TestFetchWithoutPagination(t *testing.T) {
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
//		"deleted_date", "is_deleted", "is_active", "minimum_booking_desc", "minimum_booking_amount"}).
//		AddRow(mockExpAvailability[0].Id, mockExpAvailability[0].CreatedBy, mockExpAvailability[0].CreatedDate, mockExpAvailability[0].ModifiedBy,
//			mockExpAvailability[0].ModifiedDate, mockExpAvailability[0].DeletedBy, mockExpAvailability[0].DeletedDate, mockExpAvailability[0].IsDeleted,
//			mockExpAvailability[0].IsActive, mockExpAvailability[0].ExpAvailabilityDesc, mockExpAvailability[0].ExpAvailabilityAmount).
//		AddRow(mockExpAvailability[1].Id, mockExpAvailability[1].CreatedBy, mockExpAvailability[1].CreatedDate, mockExpAvailability[1].ModifiedBy,
//			mockExpAvailability[1].ModifiedDate, mockExpAvailability[1].DeletedBy, mockExpAvailability[1].DeletedDate, mockExpAvailability[1].IsDeleted,
//			mockExpAvailability[1].IsActive, mockExpAvailability[1].ExpAvailabilityDesc, mockExpAvailability[1].ExpAvailabilityAmount)
//
//	query := `SELECT \*\ FROM minimum_bookings where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := ExpAvailabilityRepo.NewExpAvailabilityRepository(db)
//
//	//limit := 10
//	//offset := 0
//	anArticle, err := a.GetAll(context.TODO(), 0, 0)
//	//assert.NotEmpty(t, nextCursor)
//	assert.NoError(t, err)
//	assert.Len(t, anArticle, 2)
//}
//func TestFetchWithPaginationErrorFetch(t *testing.T) {
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
//		"deleted_date", "is_deleted", "is_active", "minimum_booking_desc", "minimum_booking_amount"}).
//		AddRow(mockExpAvailability[0].Id, mockExpAvailability[0].CreatedBy, mockExpAvailability[0].CreatedDate, mockExpAvailability[0].ModifiedBy,
//			mockExpAvailability[0].ModifiedDate, mockExpAvailability[0].DeletedBy, mockExpAvailability[0].DeletedDate, mockExpAvailability[0].IsDeleted,
//			mockExpAvailability[0].ModifiedDate, mockExpAvailability[0].ExpAvailabilityDesc, mockExpAvailability[0].ExpAvailabilityAmount).
//		AddRow(mockExpAvailability[1].Id, mockExpAvailability[1].CreatedBy, mockExpAvailability[1].CreatedDate, mockExpAvailability[1].ModifiedBy,
//			mockExpAvailability[1].ModifiedDate, mockExpAvailability[1].DeletedBy, mockExpAvailability[1].DeletedDate, mockExpAvailability[1].IsDeleted,
//			mockExpAvailability[1].IsActive, mockExpAvailability[1].ExpAvailabilityDesc, mockExpAvailability[1].ExpAvailabilityAmount)
//
//	query := `SELECT \*\ FROM minimum_bookings where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := ExpAvailabilityRepo.NewExpAvailabilityRepository(db)
//
//	limit := 10
//	offset := 0
//	anArticle, err := a.GetAll(context.TODO(), limit, offset)
//	//assert.NotEmpty(t, nextCursor)
//	assert.Error(t, err)
//	assert.Nil(t, anArticle)
//	//assert.Len(t, anArticle, 2)
//}
//func TestFetchWithoutPaginationErrorFetch(t *testing.T) {
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
//		"deleted_date", "is_deleted", "is_active", "minimum_booking_desc", "minimum_booking_amount"}).
//		AddRow(mockExpAvailability[0].Id, mockExpAvailability[0].CreatedBy, mockExpAvailability[0].CreatedDate, mockExpAvailability[0].ModifiedBy,
//			mockExpAvailability[0].ModifiedDate, mockExpAvailability[0].DeletedBy, mockExpAvailability[0].DeletedDate, mockExpAvailability[0].IsDeleted,
//			mockExpAvailability[0].ModifiedDate, mockExpAvailability[0].ExpAvailabilityDesc, mockExpAvailability[0].ExpAvailabilityAmount).
//		AddRow(mockExpAvailability[1].Id, mockExpAvailability[1].CreatedBy, mockExpAvailability[1].CreatedDate, mockExpAvailability[1].ModifiedBy,
//			mockExpAvailability[1].ModifiedDate, mockExpAvailability[1].DeletedBy, mockExpAvailability[1].DeletedDate, mockExpAvailability[1].IsDeleted,
//			mockExpAvailability[1].IsActive, mockExpAvailability[1].ExpAvailabilityDesc, mockExpAvailability[1].ExpAvailabilityAmount)
//
//	query := `SELECT \*\ FROM minimum_bookings where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := ExpAvailabilityRepo.NewExpAvailabilityRepository(db)
//
//	//limit := 10
//	//offset := 0
//	anArticle, err := a.GetAll(context.TODO(), 0, 0)
//	assert.Error(t, err)
//	assert.Nil(t, anArticle)
//}
