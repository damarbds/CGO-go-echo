package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/models"
	currencyRepo "github.com/service/currency/repository"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)
var (
	mockCurrency = []models.Currency{
		models.Currency{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     0,
			Code:         "testCode",
			Name:         "testName",
			Symbol:       "testSymbol",
		},
		models.Currency{
			Id:           1,
			CreatedBy:    "test",
			CreatedDate:  time.Now(),
			ModifiedBy:   nil,
			ModifiedDate: nil,
			DeletedBy:    nil,
			DeletedDate:  nil,
			IsDeleted:    0,
			IsActive:     0,
			Code:         "testCode",
			Name:         "testName",
			Symbol:       "testSymbol",
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
		AddRow(len(mockCurrency))

	query := `SELECT count\(\*\) AS count FROM currencies WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := currencyRepo.NewCurrencyRepository(db)

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

	query := `SELECT count\(\*\) AS count FROM currencies WHERE is_deleted = 0 and is_active = 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := currencyRepo.NewCurrencyRepository(db)

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


	rows := sqlmock.NewRows([]string{"id","code","name","symbol" ,"created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active"}).
		AddRow(mockCurrency[0].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol,mockCurrency[0].CreatedBy, mockCurrency[0].CreatedDate, mockCurrency[0].ModifiedBy,
			mockCurrency[0].ModifiedDate, mockCurrency[0].DeletedBy, mockCurrency[0].DeletedDate, mockCurrency[0].IsDeleted,
			mockCurrency[0].IsActive).
		AddRow(mockCurrency[1].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol, mockCurrency[1].CreatedBy, mockCurrency[1].CreatedDate, mockCurrency[1].ModifiedBy,
			mockCurrency[1].ModifiedDate, mockCurrency[1].DeletedBy, mockCurrency[1].DeletedDate, mockCurrency[1].IsDeleted,
			mockCurrency[1].IsActive)

	query := `SELECT \*\ FROM currencies where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := currencyRepo.NewCurrencyRepository(db)

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


	rows := sqlmock.NewRows([]string{"id","code","name","symbol" ,"created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active"}).
		AddRow(mockCurrency[0].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol,mockCurrency[0].CreatedBy, mockCurrency[0].CreatedDate, mockCurrency[0].ModifiedBy,
			mockCurrency[0].ModifiedDate, mockCurrency[0].DeletedBy, mockCurrency[0].DeletedDate, mockCurrency[0].IsDeleted,
			mockCurrency[0].IsActive).
		AddRow(mockCurrency[1].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol, mockCurrency[1].CreatedBy, mockCurrency[1].CreatedDate, mockCurrency[1].ModifiedBy,
			mockCurrency[1].ModifiedDate, mockCurrency[1].DeletedBy, mockCurrency[1].DeletedDate, mockCurrency[1].IsDeleted,
			mockCurrency[1].IsActive)

	query := `SELECT \*\ FROM currencies where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := currencyRepo.NewCurrencyRepository(db)

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


	rows := sqlmock.NewRows([]string{"id","code","name","symbol" ,"created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active"}).
		AddRow(mockCurrency[0].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol,mockCurrency[0].CreatedBy, mockCurrency[0].CreatedDate, mockCurrency[0].ModifiedBy,
			mockCurrency[0].ModifiedDate, mockCurrency[0].DeletedBy, mockCurrency[0].DeletedDate, mockCurrency[0].IsDeleted,
			mockCurrency[0].IsActive).
		AddRow(mockCurrency[1].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol, mockCurrency[1].CreatedBy, mockCurrency[1].CreatedDate, mockCurrency[1].ModifiedBy,
			mockCurrency[1].ModifiedDate, mockCurrency[1].DeletedBy, mockCurrency[1].DeletedDate, mockCurrency[1].IsDeleted,
			mockCurrency[1].Name)

	query := `SELECT \*\ FROM currencies where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT \? OFFSET \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := currencyRepo.NewCurrencyRepository(db)

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


	rows := sqlmock.NewRows([]string{"id","code","name","symbol" ,"created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active"}).
		AddRow(mockCurrency[0].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol,mockCurrency[0].CreatedBy, mockCurrency[0].CreatedDate, mockCurrency[0].ModifiedBy,
			mockCurrency[0].ModifiedDate, mockCurrency[0].DeletedBy, mockCurrency[0].DeletedDate, mockCurrency[0].IsDeleted,
			mockCurrency[0].IsActive).
		AddRow(mockCurrency[1].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol, mockCurrency[1].CreatedBy, mockCurrency[1].CreatedDate, mockCurrency[1].ModifiedBy,
			mockCurrency[1].ModifiedDate, mockCurrency[1].DeletedBy, mockCurrency[1].DeletedDate, mockCurrency[1].IsDeleted,
			mockCurrency[1].ModifiedDate)

	query := `SELECT \*\ FROM currencies where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := currencyRepo.NewCurrencyRepository(db)

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



	rows := sqlmock.NewRows([]string{"id","code","name","symbol" ,"created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active"}).
		AddRow(mockCurrency[0].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol,mockCurrency[0].CreatedBy, mockCurrency[0].CreatedDate, mockCurrency[0].ModifiedBy,
			mockCurrency[0].ModifiedDate, mockCurrency[0].DeletedBy, mockCurrency[0].DeletedDate, mockCurrency[0].IsDeleted,
			mockCurrency[0].IsActive)

	query := `SELECT \*\ FROM currencies WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := currencyRepo.NewCurrencyRepository(db)

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


	rows := sqlmock.NewRows([]string{"id","code","name","symbol" ,"created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active"})

	query := `SELECT \*\ FROM currencies WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := currencyRepo.NewCurrencyRepository(db)

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



	rows := sqlmock.NewRows([]string{"id","code","name","symbol" ,"created_by", "created_date", "modified_by", "modified_date", "deleted_by",
		"deleted_date", "is_deleted", "is_active"}).
		AddRow(mockCurrency[0].Id, mockCurrency[0].Code,mockCurrency[0].Name,mockCurrency[0].Symbol,mockCurrency[0].CreatedBy, mockCurrency[0].CreatedDate, mockCurrency[0].ModifiedBy,
			mockCurrency[0].ModifiedDate, mockCurrency[0].DeletedBy, mockCurrency[0].DeletedDate, mockCurrency[0].IsDeleted,
			mockCurrency[0].ModifiedDate)

	query := `SELECT \*\ FROM currencies WHERE id = \\?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := currencyRepo.NewCurrencyRepository(db)

	num := 1
	anArticle, err := a.GetById(context.TODO(), num)
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

	query := "UPDATE currencies SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := currencyRepo.NewCurrencyRepository(db)

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

	query := "UPDATE currencies SET deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? WHERE id =\\?"
	id := 2
	deletedBy := "test"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(deletedBy, time.Now(), 1, 0, id, id).WillReturnResult(sqlmock.NewResult(2, 1))

	a := currencyRepo.NewCurrencyRepository(db)

	err = a.Delete(context.TODO(), id, deletedBy)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {

	a := mockCurrency[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT currencies SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? ,code=\\?,name=\\? ,symbol=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.Code,
		a.Name,a.Symbol).WillReturnResult(sqlmock.NewResult(1, 1))

	i := currencyRepo.NewCurrencyRepository(db)

	id, err := i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	//user := "test"
	//now := time.Now()
	a := mockCurrency[0]
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "INSERT currencies SET created_by=\\? , created_date=\\? , modified_by=\\?, modified_date=\\? , 				deleted_by=\\? , deleted_date=\\? , is_deleted=\\? , is_active=\\? ,code=\\?,name=\\? ,symbol=\\? "
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.Code,
		a.Name,a.Symbol,a.Symbol).WillReturnResult(sqlmock.NewResult(1, 1))

	i := currencyRepo.NewCurrencyRepository(db)

	_, err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestUpdate(t *testing.T) {
	now := time.Now()
	modifyBy := "test"
	ar := mockCurrency[0]
	ar.ModifiedDate = &now
	ar.ModifiedBy = &modifyBy
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE currencies set modified_by=\?, modified_date=\? ,code=\?,name=\? ,symbol=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.Code, ar.Name,ar.Symbol, ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := currencyRepo.NewCurrencyRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.NoError(t, err)
	assert.Nil(t, err)
}
func TestUpdateErrorExec(t *testing.T) {
	//now := time.Now()
	//modifyBy := "test"
	ar := mockCurrency[0]

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `UPDATE currencies set modified_by=\?, modified_date=\? ,code=\?,name=\? ,symbol=\? WHERE id = \?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.ModifiedBy, ar.ModifiedDate, ar.Code, ar.Name,ar.Symbol, ar.Id,ar.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	a := currencyRepo.NewCurrencyRepository(db)

	err = a.Update(context.TODO(), &ar)
	assert.Error(t, err)
}
