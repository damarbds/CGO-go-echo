package repository_test

import (
	"context"
	ExChangeRateRepo "github.com/misc/currency/repository"
	"github.com/models"
	"testing"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetByDate(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockExChangeRate := []models.ExChangeRate{
		models.ExChangeRate{
			Id:    1,
			Date:  "2020-09-17",
			From:  "USD",
			To:    "IDR",
			Rates: 12789798,
		},
	}
	rows := sqlmock.NewRows([]string{"id", "date", "from", "to", "rates"}).
		AddRow(mockExChangeRate[0].Id, mockExChangeRate[0].Date, mockExChangeRate[0].From, mockExChangeRate[0].To,
			mockExChangeRate[0].Rates)

	query := `SELECT e.\*\ FROM ex_change_rates e WHERE e.from = \? AND e.to =\? order by e.date desc LIMIT 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExChangeRateRepo.NewExChangeRatesRepository(db)

	anArticle, err := a.GetByDate(context.TODO(), mockExChangeRate[0].From,mockExChangeRate[0].To)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Equal(t, anArticle.Rates, mockExChangeRate[0].Rates)
}
func TestGetByDateError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockExChangeRate := []models.ExChangeRate{
		models.ExChangeRate{
			Id:    1,
			Date:  "2020-09-17",
			From:  "USD",
			To:    "IDR",
			Rates: 12789798,
		},
	}
	rows := sqlmock.NewRows([]string{"id", "date", "from", "to", "rates"}).
		AddRow(mockExChangeRate[0].Id, mockExChangeRate[0].Date, mockExChangeRate[0].From, mockExChangeRate[0].To,
			mockExChangeRate[0].Date)

	query := `SELECT e.\*\ FROM ex_change_rates e WHERE e.from = \? AND e.to =\? order by e.date desc LIMIT 1`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ExChangeRateRepo.NewExChangeRatesRepository(db)

	anArticle, err := a.GetByDate(context.TODO(), mockExChangeRate[0].From,mockExChangeRate[0].To)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestInsert(t *testing.T) {
	//user := "test"
	//now := time.Now()
	a := models.ExChangeRate{
		Id:    1,
		Date:  "2020-09-17",
		From:  "USD",
		To:    "IDR",
		Rates: 121212,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT ex_change_rates SET ex_change_rates.date=\? , ex_change_rates.from=\? , ex_change_rates.to=\?, rates=\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Date, a.From,a.To,a.Rates).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExChangeRateRepo.NewExChangeRatesRepository(db)

	err = i.Insert(context.TODO(), &a)
	assert.NoError(t, err)
	//assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	a := models.ExChangeRate{
		Id:    1,
		Date:  "2020-09-17",
		From:  "USD",
		To:    "IDR",
		Rates: 121212,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT ex_change_rates SET ex_change_rates.date=\? , ex_change_rates.from=\? , ex_change_rates.to=\?, rates=\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Date, a.From,a.To,a.Rates,a.Rates).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExChangeRateRepo.NewExChangeRatesRepository(db)

	err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
func TestInsertErrorExecQuery(t *testing.T) {
	a := models.ExChangeRate{
		Id:    1,
		Date:  "2020-09-17",
		From:  "USD",
		To:    "IDR",
		Rates: 121212,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT ex_change_rates SET ex_change_rates.date=\?asdsadsa , ex_change_rates.from=\? , ex_change_rates.to=\?, rates=\?`

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Date, a.From,a.To,a.Rates).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ExChangeRateRepo.NewExChangeRatesRepository(db)

	err = i.Insert(context.TODO(), &a)
	assert.Error(t, err)
}
