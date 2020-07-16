package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	ExpAvailabilityRepo "github.com/service/exp_availability/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var (
	merchantId          = "asldkajsdiqoueeqoi"
	mockExpAvailability = []models.ExpAvailability{
		models.ExpAvailability{
			Id:                   "asdjaskldjsaldj",
			CreatedBy:            "test",
			CreatedDate:          time.Now(),
			ModifiedBy:           nil,
			ModifiedDate:         nil,
			DeletedBy:            nil,
			DeletedDate:          nil,
			IsDeleted:            0,
			IsActive:             1,
			ExpAvailabilityMonth: "April",
			ExpAvailabilityDate:  "2020-01-02",
			ExpAvailabilityYear:  2020,
			ExpId:                "qrqeqweqweqw",
		},
		models.ExpAvailability{
			Id:                   "asdjaskldjsaldj",
			CreatedBy:            "test",
			CreatedDate:          time.Now(),
			ModifiedBy:           nil,
			ModifiedDate:         nil,
			DeletedBy:            nil,
			DeletedDate:          nil,
			IsDeleted:            0,
			IsActive:             1,
			ExpAvailabilityMonth: "May",
			ExpAvailabilityDate:  "2020-01-02",
			ExpAvailabilityYear:  2020,
			ExpId:                "qrqeqweqweqw",
		},
	}
)

func TestGetCountDate(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(len(mockExpAvailability))

	query := `SELECT
		count\(DISTINCT b.booking_date,b.exp_id\) AS count
	FROM
		transactions t
	JOIN booking_exps b on b.id = t.booking_exp_id
	JOIN experiences e on e.id = b.exp_id
	JOIN merchants m on m.id = e.merchant_id
	WHERE
		DATE \(b.booking_date\) = \? AND 
		t.status in \(0,1,2,5\) AND 
		e.merchant_id = \? AND 
		t.is_deleted = 0 AND 
		t.is_active = 1 `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	res, err := a.GetCountDate(context.TODO(), mockExpAvailability[0].ExpAvailabilityDate, merchantId)
	assert.NoError(t, err)
	assert.Equal(t, res, 2, "")
}
func TestGetCountDateErrorFetch(t *testing.T) {
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

	query := `SELECT
		count\(DISTINCT b.booking_date,b.exp_id\) AS count
	FROM
		transactions t
	JOIN booking_exps b on b.id = t.booking_exp_id
	JOIN experiences e on e.id = b.exp_id
	JOIN merchants m on m.id = e.merchant_id
	WHERE
		DATE \(b.booking_date\) = \? AND 
		t.status in \(0,1,2,5\) AND 
		e.merchant_id = \? AND 
		t.is_deleted = 0 AND 
		t.is_active = 1 `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	_, err = a.GetCountDate(context.TODO(), mockExpAvailability[0].ExpAvailabilityDate, merchantId)
	assert.Error(t, err)
}
func TestGetByExpIds(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "exp_availability_month", "exp_availability_date",
		"exp_availability_year", "exp_id"}).
		AddRow(mockExpAvailability[0].Id, mockExpAvailability[0].CreatedBy, mockExpAvailability[0].CreatedDate, mockExpAvailability[0].ModifiedBy,
			mockExpAvailability[0].ModifiedDate, mockExpAvailability[0].DeletedBy, mockExpAvailability[0].DeletedDate, mockExpAvailability[0].IsDeleted,
			mockExpAvailability[0].IsActive, mockExpAvailability[0].ExpAvailabilityMonth, mockExpAvailability[0].ExpAvailabilityDate,
			mockExpAvailability[0].ExpAvailabilityYear, mockExpAvailability[0].ExpId).
		AddRow(mockExpAvailability[1].Id, mockExpAvailability[1].CreatedBy, mockExpAvailability[1].CreatedDate, mockExpAvailability[1].ModifiedBy,
			mockExpAvailability[1].ModifiedDate, mockExpAvailability[1].DeletedBy, mockExpAvailability[1].DeletedDate, mockExpAvailability[1].IsDeleted,
			mockExpAvailability[1].IsActive, mockExpAvailability[1].ExpAvailabilityMonth, mockExpAvailability[1].ExpAvailabilityDate,
			mockExpAvailability[1].ExpAvailabilityYear, mockExpAvailability[1].ExpId)
	var expIds []*string
	var expId = []string{"qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw"}
	for _, id := range expId {
		expIds = append(expIds, &id)
	}
	query := `SELECT \*\ FROM exp_availabilities 
				WHERE is_deleted = 0 AND 
						is_active = 1 AND 
						exp_availability_month = \? AND 
						exp_availability_year = \? `
	for index, id := range expIds {
		if index == 0 && index != (len(expId)-1) {
			query = query + ` AND \(exp_id = '` + *id + `' `
		} else if index == 0 && index == (len(expId)-1) {
			query = query + ` AND \(exp_id = '` + *id + `' \) `
		} else if index == (len(expId) - 1) {
			query = query + ` OR  exp_id = '` + *id + `' \) `
		} else {
			query = query + ` OR  exp_id = '` + *id + `' `
		}
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	anArticle, err := a.GetByExpIds(context.TODO(), expIds, mockExpAvailability[0].ExpAvailabilityYear, mockExpAvailability[0].ExpAvailabilityMonth)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetByExpIdsErrorFetchQuery(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "exp_availability_month", "exp_availability_date",
		"exp_availability_year", "exp_id"}).
		AddRow(mockExpAvailability[0].Id, mockExpAvailability[0].CreatedBy, mockExpAvailability[0].CreatedDate, mockExpAvailability[0].ModifiedBy,
			mockExpAvailability[0].ModifiedDate, mockExpAvailability[0].DeletedBy, mockExpAvailability[0].DeletedDate, mockExpAvailability[0].IsDeleted,
			mockExpAvailability[0].IsActive, mockExpAvailability[0].ExpAvailabilityMonth, mockExpAvailability[0].ExpAvailabilityDate,
			mockExpAvailability[0].ExpAvailabilityYear, mockExpAvailability[0].ExpId).
		AddRow(mockExpAvailability[1].Id, mockExpAvailability[1].CreatedBy, mockExpAvailability[1].CreatedDate, mockExpAvailability[1].ModifiedBy,
			mockExpAvailability[1].ModifiedDate, mockExpAvailability[1].DeletedBy, mockExpAvailability[1].DeletedDate, mockExpAvailability[1].IsDeleted,
			mockExpAvailability[1].IsActive, mockExpAvailability[1].ExpAvailabilityMonth, mockExpAvailability[1].ExpAvailabilityDate,
			mockExpAvailability[1].ExpAvailabilityYear, mockExpAvailability[1].ExpId)
	var expIds []*string
	var expId = []string{"qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw"}
	for _, id := range expId {
		expIds = append(expIds, &id)
	}
	query := `SELECT * FROM exp_availabilities 
				WHERE is_deleted = 0 AND dasdsadsa
						is_active = 1 AND 
						exp_availability_month = ? AND 
						exp_availability_year = ? `
	for index, id := range expIds {
		if index == 0 && index != (len(expId)-1) {
			query = query + ` AND (exp_id = '` + *id + `' `
		} else if index == 0 && index == (len(expId)-1) {
			query = query + ` AND (exp_id = '` + *id + `' ) `
		} else if index == (len(expId) - 1) {
			query = query + ` OR  exp_id = '` + *id + `' ) `
		} else {
			query = query + ` OR  exp_id = '` + *id + `' `
		}
	}
	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	anArticle, err := a.GetByExpIds(context.TODO(), expIds, mockExpAvailability[0].ExpAvailabilityYear, mockExpAvailability[0].ExpAvailabilityMonth)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
	//assert.Len(t, anArticle, 2)
}
func TestGetByExpId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"id", "created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active", "exp_availability_month", "exp_availability_date",
		"exp_availability_year", "exp_id"}).
		AddRow(mockExpAvailability[0].Id, mockExpAvailability[0].CreatedBy, mockExpAvailability[0].CreatedDate, mockExpAvailability[0].ModifiedBy,
			mockExpAvailability[0].ModifiedDate, mockExpAvailability[0].DeletedBy, mockExpAvailability[0].DeletedDate, mockExpAvailability[0].IsDeleted,
			mockExpAvailability[0].IsActive, mockExpAvailability[0].ExpAvailabilityMonth, mockExpAvailability[0].ExpAvailabilityDate,
			mockExpAvailability[0].ExpAvailabilityYear, mockExpAvailability[0].ExpId).
		AddRow(mockExpAvailability[1].Id, mockExpAvailability[1].CreatedBy, mockExpAvailability[1].CreatedDate, mockExpAvailability[1].ModifiedBy,
			mockExpAvailability[1].ModifiedDate, mockExpAvailability[1].DeletedBy, mockExpAvailability[1].DeletedDate, mockExpAvailability[1].IsDeleted,
			mockExpAvailability[1].IsActive, mockExpAvailability[1].ExpAvailabilityMonth, mockExpAvailability[1].ExpAvailabilityDate,
			mockExpAvailability[1].ExpAvailabilityYear, mockExpAvailability[1].ExpId)

	query := `SELECT \*\ FROM exp_availabilities WHERE is_deleted = 0 AND is_active = 1 AND exp_id = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	//limit := 10
	//offset := 0
	anArticle, err := a.GetByExpId(context.TODO(), mockExpAvailability[0].ExpId)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestInsert(t *testing.T) {

	a := mockExpAvailability[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT exp_availabilities SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_availability_month=\?,
				exp_availability_date=\?,exp_availability_year=\?,exp_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExpAvailabilityMonth,
		a.ExpAvailabilityDate, a.ExpAvailabilityYear, a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	id, err := i.Insert(context.TODO(), a)
	assert.NoError(t, err)
	assert.Equal(t, id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	a := mockExpAvailability[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT exp_availabilities SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_availability_month=\?,
				exp_availability_date=\?,exp_availability_year=\?,exp_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.ExpAvailabilityMonth,
		a.ExpAvailabilityDate, a.ExpAvailabilityYear, a.ExpId,a.ExpId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	_, err = i.Insert(context.TODO(), a)

	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockExpAvailability[0]
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

	query := `UPDATE exp_availabilities SET modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_availability_month=\?,
				exp_availability_date=\?,exp_availability_year=\?,exp_id=\? 
				WHERE id =\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, a.ModifiedDate, nil, nil, 0, 1, a.ExpAvailabilityMonth,
		a.ExpAvailabilityDate, a.ExpAvailabilityYear, a.ExpId, a.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	err = u.Update(context.TODO(), a)
	assert.NoError(t, err)
	assert.Nil(t, err)
}
func TestUpdateErrorExec(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	a := mockExpAvailability[0]
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

	query := `UPDATE exp_availabilities SET modified_by=\?, modified_date=\? , 
				deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? , exp_availability_month=\?,
				exp_availability_date=\?,exp_availability_year=\?,exp_id=\? 
				WHERE id =\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.ModifiedBy, a.ModifiedDate, nil, nil, 0, 1, a.ExpAvailabilityMonth,
		a.ExpAvailabilityDate, a.ExpAvailabilityYear, a.ExpId, a.Id,a.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	u := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	err = u.Update(context.TODO(), a)
	assert.Error(t, err)
}
func TestDeletes(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	ids := []string{"asdjaskldjsaldj", "asdjaskldjsaldj", "asdjaskldjsaldj", "asdjaskldjsaldj", "asdjaskldjsaldj", "asdjaskldjsaldj"}

	query := `UPDATE  exp_availabilities SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`
	for index, id := range ids {
		if index == 0 && index != (len(ids)-1) {
			query = query + ` AND \(id !=` + id
		} else if index == 0 && index == (len(ids)-1) {
			query = query + ` AND \(id !=` + id + ` \) `
		} else if index == (len(ids) - 1) {
			query = query + ` OR id !=` + id + ` \) `
		} else {
			query = query + ` OR id !=` + id
		}
	}

	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, mockExpAvailability[0].ExpId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	err = a.Deletes(context.TODO(), ids, mockExpAvailability[0].ExpId,deletedBy)
	assert.NoError(t, err)
}
func TestDeletesErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	ids := []string{"asdjaskldjsaldj", "asdjaskldjsaldj", "asdjaskldjsaldj", "asdjaskldjsaldj", "asdjaskldjsaldj", "asdjaskldjsaldj"}

	query := `UPDATE  exp_availabilities SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`
	for index, id := range ids {
		if index == 0 && index != (len(ids)-1) {
			query = query + ` AND \(id !=` + id
		} else if index == 0 && index == (len(ids)-1) {
			query = query + ` AND \(id !=` + id + ` \) `
		} else if index == (len(ids) - 1) {
			query = query + ` OR id !=` + id + ` \) `
		} else {
			query = query + ` OR id !=` + id
		}
	}

	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, mockExpAvailability[0].ExpId,mockExpAvailability[0].ExpId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	err = a.Deletes(context.TODO(), ids, mockExpAvailability[0].ExpId,deletedBy)
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

	query := `UPDATE exp_availabilities SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`

	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, mockExpAvailability[0].ExpId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	err = a.DeleteByExpId(context.TODO(),mockExpAvailability[0].ExpId,deletedBy)
	assert.NoError(t, err)
}
func TestDeleteByExpIdErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE exp_availabilities SET deleted_by=\? , deleted_date=\? , is_deleted=\? , is_active=\? WHERE exp_id=\?`

	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, mockExpAvailability[0].ExpId,mockExpAvailability[0].ExpId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ExpAvailabilityRepo.NewExpavailabilityRepository(db)

	err = a.DeleteByExpId(context.TODO(),mockExpAvailability[0].ExpId,deletedBy)
	assert.Error(t, err)
}
