package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	FilterActivityTypeRepo "github.com/service/filter_activity_type/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	expId = "asdasdasdasd"
	mockFilterActivityType = []models.FilterActivityType{
		models.FilterActivityType{
			Id:           1,
			CreatedBy:    "Test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExpTypeId:    1,
			ExpId:        &expId,
		},
		models.FilterActivityType{
			Id:           2,
			CreatedBy:    "Test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     1,
			ExpTypeId:    1,
			ExpId:        &expId,
		},
	}

)

func TestInsert(t *testing.T) {
	a := mockFilterActivityType[0]

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT filter_activity_types SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_type_id=\? , exp_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExpTypeId, a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := FilterActivityTypeRepo.NewFilterActivityTypeRepository(db)

	err = i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	//assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {

	a := mockFilterActivityType[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT filter_activity_types SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? ,
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_type_id=\? , exp_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExpTypeId, a.ExpId,a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := FilterActivityTypeRepo.NewFilterActivityTypeRepository(db)

	err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestDeleteByExpId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "DELETE FROM filter_activity_types WHERE exp_id = \\?"

	num := mockFilterActivityType[0].ExpId

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := FilterActivityTypeRepo.NewFilterActivityTypeRepository(db)

	err = a.DeleteByExpId(context.TODO(), *num)
	assert.NoError(t, err)
}
func TestDeleteByExpIdErrorExecQueryString(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	query := "DELETE FROM filter_activity_types WHERE exp_id = \\?asdasdsasa"

	num := mockFilterActivityType[0].ExpId

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(num).WillReturnResult(sqlmock.NewResult(12, 1))

	a := FilterActivityTypeRepo.NewFilterActivityTypeRepository(db)

	err = a.DeleteByExpId(context.TODO(), *num)
	assert.Error(t, err)
}
