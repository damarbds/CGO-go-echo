package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	MinimumBookingRepo "github.com/service/minimum_booking/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	mockMinimumBooking = []models.MinimumBooking{
		models.MinimumBooking{
			Id:                   "adadasdqrqr",
			CreatedBy:            "test",
			CreatedDate:          time.Now(),
			ModifiedBy:           nil,
			ModifiedDate:         nil,
			DeletedBy:            nil,
			DeletedDate:          nil,
			IsDeleted:            0,
			IsActive:             1,
			MinimumBookingDesc:   "Test MinimumBooking 1",
			MinimumBookingAmount: 1,
		},
		models.MinimumBooking{
			Id:                   "klklklklklk",
			CreatedBy:            "test",
			CreatedDate:          time.Now(),
			ModifiedBy:           nil,
			ModifiedDate:         nil,
			DeletedBy:            nil,
			DeletedDate:          nil,
			IsDeleted:            0,
			IsActive:             1,
			MinimumBookingDesc:   "Test MinimumBooking 2",
			MinimumBookingAmount: 1,
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
		AddRow(len(mockMinimumBooking))

	query := `SELECT count\(\*\) AS count FROM minimum_bookings WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := MinimumBookingRepo.NewMinimumBookingRepository(db)

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

	query := `SELECT count\(\*\) AS count FROM minimum_bookings WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := MinimumBookingRepo.NewMinimumBookingRepository(db)

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
		"deleted_date", "is_deleted", "is_active", "minimum_booking_desc", "minimum_booking_amount"}).
		AddRow(mockMinimumBooking[0].Id, mockMinimumBooking[0].CreatedBy, mockMinimumBooking[0].CreatedDate, mockMinimumBooking[0].ModifiedBy,
			mockMinimumBooking[0].ModifiedDate, mockMinimumBooking[0].DeletedBy, mockMinimumBooking[0].DeletedDate, mockMinimumBooking[0].IsDeleted,
			mockMinimumBooking[0].IsActive, mockMinimumBooking[0].MinimumBookingDesc, mockMinimumBooking[0].MinimumBookingAmount).
		AddRow(mockMinimumBooking[1].Id, mockMinimumBooking[1].CreatedBy, mockMinimumBooking[1].CreatedDate, mockMinimumBooking[1].ModifiedBy,
			mockMinimumBooking[1].ModifiedDate, mockMinimumBooking[1].DeletedBy, mockMinimumBooking[1].DeletedDate, mockMinimumBooking[1].IsDeleted,
			mockMinimumBooking[1].IsActive, mockMinimumBooking[1].MinimumBookingDesc, mockMinimumBooking[1].MinimumBookingAmount)

	query := `SELECT \*\ FROM minimum_bookings where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := MinimumBookingRepo.NewMinimumBookingRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.GetAll(context.TODO(), limit, offset)
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
		"deleted_date", "is_deleted", "is_active", "minimum_booking_desc", "minimum_booking_amount"}).
		AddRow(mockMinimumBooking[0].Id, mockMinimumBooking[0].CreatedBy, mockMinimumBooking[0].CreatedDate, mockMinimumBooking[0].ModifiedBy,
			mockMinimumBooking[0].ModifiedDate, mockMinimumBooking[0].DeletedBy, mockMinimumBooking[0].DeletedDate, mockMinimumBooking[0].IsDeleted,
			mockMinimumBooking[0].IsActive, mockMinimumBooking[0].MinimumBookingDesc, mockMinimumBooking[0].MinimumBookingAmount).
		AddRow(mockMinimumBooking[1].Id, mockMinimumBooking[1].CreatedBy, mockMinimumBooking[1].CreatedDate, mockMinimumBooking[1].ModifiedBy,
			mockMinimumBooking[1].ModifiedDate, mockMinimumBooking[1].DeletedBy, mockMinimumBooking[1].DeletedDate, mockMinimumBooking[1].IsDeleted,
			mockMinimumBooking[1].IsActive, mockMinimumBooking[1].MinimumBookingDesc, mockMinimumBooking[1].MinimumBookingAmount)

	query := `SELECT \*\ FROM minimum_bookings where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := MinimumBookingRepo.NewMinimumBookingRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.GetAll(context.TODO(), 0, 0)
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
		"deleted_date", "is_deleted", "is_active", "minimum_booking_desc", "minimum_booking_amount"}).
		AddRow(mockMinimumBooking[0].Id, mockMinimumBooking[0].CreatedBy, mockMinimumBooking[0].CreatedDate, mockMinimumBooking[0].ModifiedBy,
			mockMinimumBooking[0].ModifiedDate, mockMinimumBooking[0].DeletedBy, mockMinimumBooking[0].DeletedDate, mockMinimumBooking[0].IsDeleted,
			mockMinimumBooking[0].ModifiedDate, mockMinimumBooking[0].MinimumBookingDesc, mockMinimumBooking[0].MinimumBookingAmount).
		AddRow(mockMinimumBooking[1].Id, mockMinimumBooking[1].CreatedBy, mockMinimumBooking[1].CreatedDate, mockMinimumBooking[1].ModifiedBy,
			mockMinimumBooking[1].ModifiedDate, mockMinimumBooking[1].DeletedBy, mockMinimumBooking[1].DeletedDate, mockMinimumBooking[1].IsDeleted,
			mockMinimumBooking[1].IsActive, mockMinimumBooking[1].MinimumBookingDesc, mockMinimumBooking[1].MinimumBookingAmount)

	query := `SELECT \*\ FROM minimum_bookings where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := MinimumBookingRepo.NewMinimumBookingRepository(db)

	limit := 10
	offset := 0
	anArticle, err := a.GetAll(context.TODO(), limit, offset)
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
		"deleted_date", "is_deleted", "is_active", "minimum_booking_desc", "minimum_booking_amount"}).
		AddRow(mockMinimumBooking[0].Id, mockMinimumBooking[0].CreatedBy, mockMinimumBooking[0].CreatedDate, mockMinimumBooking[0].ModifiedBy,
			mockMinimumBooking[0].ModifiedDate, mockMinimumBooking[0].DeletedBy, mockMinimumBooking[0].DeletedDate, mockMinimumBooking[0].IsDeleted,
			mockMinimumBooking[0].ModifiedDate, mockMinimumBooking[0].MinimumBookingDesc, mockMinimumBooking[0].MinimumBookingAmount).
		AddRow(mockMinimumBooking[1].Id, mockMinimumBooking[1].CreatedBy, mockMinimumBooking[1].CreatedDate, mockMinimumBooking[1].ModifiedBy,
			mockMinimumBooking[1].ModifiedDate, mockMinimumBooking[1].DeletedBy, mockMinimumBooking[1].DeletedDate, mockMinimumBooking[1].IsDeleted,
			mockMinimumBooking[1].IsActive, mockMinimumBooking[1].MinimumBookingDesc, mockMinimumBooking[1].MinimumBookingAmount)

	query := `SELECT \*\ FROM minimum_bookings where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := MinimumBookingRepo.NewMinimumBookingRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.GetAll(context.TODO(), 0, 0)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
