package repository_test

import (
	"context"
	"testing"
	"time"

	TimesOptionsRepo "github.com/service/time_options/repository"

	"github.com/models"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	expId           = "asdasdasdasd"
	mockTimesOption = []models.TimesOption{
		models.TimesOption{
			Id:           1,
			CreatedBy:    "Test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			StartTime:    "00:00:00",
			EndTime:      "06:00:00",
		},
		models.TimesOption{
			Id:           2,
			CreatedBy:    "Test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			StartTime:    "00:00:00",
			EndTime:      "06:00:00",
		},
	}
)

//func TestInsert(t *testing.T) {
//	a := mockTimesOption[0]
//
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	query := `INSERT filter_activity_types SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
//				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_type_id=\? , exp_id=\?`
//	prep := mock.ExpectPrepare(query)
//	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExpTypeId, a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))
//
//	i := TimesOptionRepo.NewTimesOptionRepository(db)
//
//	err = i.Insert(context.TODO(), &a)
//	assert.NoError(t, err)
//	//assert.Equal(t, *id, a.Id)
//}
//func TestInsertErrorExec(t *testing.T) {
//
//	a := mockTimesOption[0]
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//	query := `INSERT filter_activity_types SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
//				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_type_id=\? , exp_id=\?`
//	prep := mock.ExpectPrepare(query)
//	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExpTypeId, a.ExpId, a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))
//
//	i := TimesOptionRepo.NewTimesOptionRepository(db)
//
//	err = i.Insert(context.TODO(), &a)
//	assert.Error(t, err)
//}
func TestGetByTime(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "start_time", "end_time"}).
		AddRow(mockTimesOption[0].Id, mockTimesOption[0].CreatedBy, mockTimesOption[0].CreatedDate, mockTimesOption[0].ModifiedBy,
			mockTimesOption[0].ModifiedDate, mockTimesOption[0].DeletedBy, mockTimesOption[0].DeletedDate, mockTimesOption[0].IsDeleted,
			mockTimesOption[0].IsActive, mockTimesOption[0].StartTime, mockTimesOption[0].EndTime).
		AddRow(mockTimesOption[1].Id, mockTimesOption[1].CreatedBy, mockTimesOption[1].CreatedDate, mockTimesOption[1].ModifiedBy,
			mockTimesOption[1].ModifiedDate, mockTimesOption[1].DeletedBy, mockTimesOption[1].DeletedDate, mockTimesOption[1].IsDeleted,
			mockTimesOption[1].IsActive, mockTimesOption[1].StartTime, mockTimesOption[1].EndTime)

	query := `SELECT \*\ FROM times_options where \? >= start_time  AND  \? <= end_time `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := TimesOptionsRepo.NewTimeOptionsRepository(db)

	times := "03:00:00"
	anArticle, err := a.GetByTime(context.TODO(), times)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Equal(t,anArticle.StartTime,mockTimesOption[0].StartTime)
	//assert.Len(t, anArticle, 2)
}
func TestGetByTimeErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "start_time", "end_time"}).
		AddRow(mockTimesOption[0].Id, mockTimesOption[0].CreatedBy, mockTimesOption[0].CreatedDate, mockTimesOption[0].ModifiedBy,
			mockTimesOption[0].ModifiedDate, mockTimesOption[0].DeletedBy, mockTimesOption[0].DeletedDate, mockTimesOption[0].IsDeleted,
			mockTimesOption[0].IsActive, mockTimesOption[0].StartTime, mockTimesOption[0].EndTime).
		AddRow(mockTimesOption[1].Id, mockTimesOption[1].CreatedBy, mockTimesOption[1].CreatedDate, mockTimesOption[1].ModifiedBy,
			mockTimesOption[1].ModifiedDate, mockTimesOption[1].DeletedBy, mockTimesOption[1].DeletedDate, mockTimesOption[1].IsDeleted,
			mockTimesOption[1].ModifiedDate, mockTimesOption[1].StartTime, mockTimesOption[1].EndTime)

	query := `SELECT \*\ FROM times_options where \? >= start_time  AND  \? <= end_time dadsadsa`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := TimesOptionsRepo.NewTimeOptionsRepository(db)

	times := "03:00:00"
	_, err = a.GetByTime(context.TODO(), times)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Len(t, anArticle, 2)
}
func TestTimeOptions(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "start_time", "end_time"}).
		AddRow(mockTimesOption[0].Id, mockTimesOption[0].CreatedBy, mockTimesOption[0].CreatedDate, mockTimesOption[0].ModifiedBy,
			mockTimesOption[0].ModifiedDate, mockTimesOption[0].DeletedBy, mockTimesOption[0].DeletedDate, mockTimesOption[0].IsDeleted,
			mockTimesOption[0].IsActive, mockTimesOption[0].StartTime, mockTimesOption[0].EndTime).
		AddRow(mockTimesOption[1].Id, mockTimesOption[1].CreatedBy, mockTimesOption[1].CreatedDate, mockTimesOption[1].ModifiedBy,
			mockTimesOption[1].ModifiedDate, mockTimesOption[1].DeletedBy, mockTimesOption[1].DeletedDate, mockTimesOption[1].IsDeleted,
			mockTimesOption[1].IsActive, mockTimesOption[1].StartTime, mockTimesOption[1].EndTime)


	query := `SELECT \*\ FROM times_options WHERE is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := TimesOptionsRepo.NewTimeOptionsRepository(db)

	anArticle, err := a.TimeOptions(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestTimeOptionsErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "start_time", "end_time"}).
		AddRow(mockTimesOption[0].Id, mockTimesOption[0].CreatedBy, mockTimesOption[0].CreatedDate, mockTimesOption[0].ModifiedBy,
			mockTimesOption[0].ModifiedDate, mockTimesOption[0].DeletedBy, mockTimesOption[0].DeletedDate, mockTimesOption[0].IsDeleted,
			mockTimesOption[0].IsActive, mockTimesOption[0].StartTime, mockTimesOption[0].EndTime).
		AddRow(mockTimesOption[1].Id, mockTimesOption[1].CreatedBy, mockTimesOption[1].CreatedDate, mockTimesOption[1].ModifiedBy,
			mockTimesOption[1].ModifiedDate, mockTimesOption[1].DeletedBy, mockTimesOption[1].DeletedDate, mockTimesOption[1].IsDeleted,
			mockTimesOption[1].ModifiedDate, mockTimesOption[1].StartTime, mockTimesOption[1].EndTime)


	query := `SELECT \*\ FROM times_options WHERE is_deleted = 0 AND is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := TimesOptionsRepo.NewTimeOptionsRepository(db)

	_, err = a.TimeOptions(context.TODO())
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	//assert.Len(t, anArticle, 2)
}
