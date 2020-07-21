package repository_test

import (
	"context"
	"github.com/models"
	ScheduleRepo "github.com/service/schedule/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestGetCountSchedule(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	rows := sqlmock.NewRows([]string{"count"}).
		AddRow(2)

	query := `
	SELECT
		count\(DISTINCT b.booking_date,b.trans_id\) AS count
	FROM
		transactions t
	JOIN booking_exps b on b.id = t.booking_exp_id OR b.order_id = t.order_id
	JOIN transportations trans on trans.id = b.trans_id
	JOIN merchants m on m.id = trans.merchant_id
	WHERE
		DATE \(b.booking_date\) = \? AND 
		t.status in \(0,1,2,5\) AND 
		trans.merchant_id = \? AND
		t.is_deleted = 0 AND
		t.is_active = 1 `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)

	merchantId := "asldjlkdjalkdjasd"
	date := time.Now().Format("2006-01-02")
	res, err := a.GetCountSchedule(context.TODO(),merchantId,date)
	assert.NoError(t, err)
	assert.Equal(t, res, 2, "")
}
func TestGetCountScheduleErrorFetch(t *testing.T) {
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

	query := `
	SELECT
		count\(DISTINCT b.booking_date,b.trans_id\) AS count
	FROM
		transactions t
	JOIN booking_exps b on b.id = t.booking_exp_id OR b.order_id = t.order_id
	JOIN transportations trans on trans.id = b.trans_id
	JOIN merchants m on m.id = trans.merchant_id
	WHERE
		DATE \(b.booking_date\) = \? AND 
		t.status in \(0,1,2,5\) AND 
		trans.merchant_id = \? AND
		t.is_deleted = 0 AND
		t.is_active = 1 `

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)

	merchantId := "asldjlkdjalkdjasd"
	date := time.Now().Format("2006-01-02")
	_, err = a.GetCountSchedule(context.TODO(),merchantId,date)
	assert.Error(t, err)
}
func TestGetTimeByTransId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockSchedule := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	rows := sqlmock.NewRows([]string{"departure_time", "arrival_time"}).
		AddRow(mockSchedule[0].DepartureTime, mockSchedule[0].ArrivalTime).
		AddRow(mockSchedule[1].DepartureTime, mockSchedule[1].ArrivalTime)

	query := `SELECT DISTINCT departure_time,arrival_time FROM schedules WHERE trans_id = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)
	transId := "kljlkjdioqjei"
	anArticle, err := a.GetTimeByTransId(context.TODO(),transId)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetTimeByTransIdError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockSchedule := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	rows := sqlmock.NewRows([]string{"departure_time", "arrival_time"}).
		AddRow(mockSchedule[0].DepartureTime, mockSchedule[0].ArrivalTime).
		AddRow(nil, mockSchedule[1].ArrivalTime)

	query := `SELECT DISTINCT departure_time,arrival_time FROM schedules WHERE trans_id = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)
	transId := "kljlkjdioqjei"
	anArticle, err := a.GetTimeByTransId(context.TODO(),transId)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetYearByTransId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockScheduleTime := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	mockSchedule := []*models.ScheduleYear{
		&models.ScheduleYear{
			Year:2020,
		},
		&models.ScheduleYear{
			Year:2021,
		},
	}
	rows := sqlmock.NewRows([]string{"year"}).
		AddRow(mockSchedule[0].Year).
		AddRow(mockSchedule[1].Year)

	query := `SELECT DISTINCT year FROM schedules WHERE trans_id = \? AND arrival_time = \? AND departure_time = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)
	transId := "kljlkjdioqjei"
	anArticle, err := a.GetYearByTransId(context.TODO(),transId,mockScheduleTime[0].ArrivalTime,mockScheduleTime[0].DepartureTime)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetYearByTransIdError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockScheduleTime := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	mockSchedule := []*models.ScheduleYear{
		&models.ScheduleYear{
			Year:2020,
		},
		&models.ScheduleYear{
			Year:2021,
		},
	}
	rows := sqlmock.NewRows([]string{"year"}).
		AddRow(mockSchedule[0].Year).
		AddRow(time.Now())

	query := `SELECT DISTINCT year FROM schedules WHERE trans_id = \? AND arrival_time = \? AND departure_time = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)
	transId := "kljlkjdioqjei"
	anArticle, err := a.GetYearByTransId(context.TODO(),transId,mockScheduleTime[0].ArrivalTime,mockScheduleTime[0].DepartureTime)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetMonthByTransId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockScheduleTime := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	mockScheduleYear := []*models.ScheduleYear{
		&models.ScheduleYear{
			Year:2020,
		},
		&models.ScheduleYear{
			Year:2021,
		},
	}
	mockSchedule := []*models.ScheduleMonth{
		&models.ScheduleMonth{
			Year:  2020,
			Month: "May",
		},
		&models.ScheduleMonth{
			Year:  2020,
			Month: "June",
		},
	}
	rows := sqlmock.NewRows([]string{"year","month"}).
		AddRow(mockSchedule[0].Year,mockSchedule[0].Month).
		AddRow(mockSchedule[1].Year,mockSchedule[1].Month)

	query := `SELECT DISTINCT year,month FROM schedules WHERE trans_id = \? AND year =\? AND arrival_time = \? AND departure_time = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)
	transId := "kljlkjdioqjei"
	anArticle, err := a.GetMonthByTransId(context.TODO(),transId,mockScheduleYear[0].Year,mockScheduleTime[0].ArrivalTime,mockScheduleTime[0].DepartureTime)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetMonthByTransIdError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockScheduleTime := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	mockScheduleYear := []*models.ScheduleYear{
		&models.ScheduleYear{
			Year:2020,
		},
		&models.ScheduleYear{
			Year:2021,
		},
	}
	mockSchedule := []*models.ScheduleMonth{
		&models.ScheduleMonth{
			Year:  2020,
			Month: "May",
		},
		&models.ScheduleMonth{
			Year:  2020,
			Month: "June",
		},
	}
	rows := sqlmock.NewRows([]string{"year","month"}).
		AddRow(mockSchedule[0].Year,mockSchedule[0].Month).
		AddRow(mockSchedule[1].Year,nil)

	query := `SELECT DISTINCT year,month FROM schedules WHERE trans_id = \? AND year =\? AND arrival_time = \? AND departure_time = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)
	transId := "kljlkjdioqjei"
	anArticle, err := a.GetMonthByTransId(context.TODO(),transId,mockScheduleYear[0].Year,mockScheduleTime[0].ArrivalTime,mockScheduleTime[0].DepartureTime)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetDayByTransId(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockScheduleTime := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	mockScheduleYear := []*models.ScheduleYear{
		&models.ScheduleYear{
			Year:2020,
		},
		&models.ScheduleYear{
			Year:2021,
		},
	}
	mockScheduleMonth := []*models.ScheduleMonth{
		&models.ScheduleMonth{
			Year:  2020,
			Month: "May",
		},
		&models.ScheduleMonth{
			Year:  2020,
			Month: "June",
		},
	}
	mockSchedule := []*models.ScheduleDay{
		&models.ScheduleDay{
			Year:          2020,
			Month:         "May",
			Day:           "Thursday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
		&models.ScheduleDay{
			Year:          2020,
			Month:         "June",
			Day:           "Monday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
	}
	rows := sqlmock.NewRows([]string{"year","month","day","departure_date","price"}).
		AddRow(mockSchedule[0].Year,mockSchedule[0].Month,mockSchedule[0].Day,mockSchedule[0].DepartureDate,mockSchedule[0].Price).
		AddRow(mockSchedule[1].Year,mockSchedule[1].Month,mockSchedule[1].Day,mockSchedule[1].DepartureDate,mockSchedule[1].Price)

	query := `SELECT DISTINCT year,month,day,departure_date,price FROM schedules WHERE trans_id = \? AND year =\? AND month=\? AND arrival_time = \? AND departure_time = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)
	transId := "kljlkjdioqjei"
	anArticle, err := a.GetDayByTransId(context.TODO(),transId,mockScheduleYear[0].Year,mockScheduleMonth[0].Month,mockScheduleTime[0].ArrivalTime,mockScheduleTime[0].DepartureTime)
	//assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, anArticle, 2)
}
func TestGetDayByTransIdError(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockScheduleTime := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	mockScheduleYear := []*models.ScheduleYear{
		&models.ScheduleYear{
			Year:2020,
		},
		&models.ScheduleYear{
			Year:2021,
		},
	}
	mockScheduleMonth := []*models.ScheduleMonth{
		&models.ScheduleMonth{
			Year:  2020,
			Month: "May",
		},
		&models.ScheduleMonth{
			Year:  2020,
			Month: "June",
		},
	}
	mockSchedule := []*models.ScheduleDay{
		&models.ScheduleDay{
			Year:          2020,
			Month:         "May",
			Day:           "Thursday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
		&models.ScheduleDay{
			Year:          2020,
			Month:         "June",
			Day:           "Monday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
	}
	rows := sqlmock.NewRows([]string{"year","month","day","departure_date","price"}).
		AddRow(mockSchedule[0].Year,nil,mockSchedule[0].Day,mockSchedule[0].DepartureDate,mockSchedule[0].Price).
		AddRow(mockSchedule[1].Year,mockSchedule[1].Month,mockSchedule[1].Day,mockSchedule[1].DepartureDate,mockSchedule[1].Price)

	query := `SELECT DISTINCT year,month,day,departure_date,price FROM schedules WHERE trans_id = \? AND year =\? AND month=\? AND arrival_time = \? AND departure_time = \?`

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)
	transId := "kljlkjdioqjei"
	anArticle, err := a.GetDayByTransId(context.TODO(),transId,mockScheduleYear[0].Year,mockScheduleMonth[0].Month,mockScheduleTime[0].ArrivalTime,mockScheduleTime[0].DepartureTime)
	//assert.NotEmpty(t, nextCursor)
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestGetScheduleByTransIds(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()
	mockSchedule := []*models.ScheduleDtos{
		&models.ScheduleDtos{
			DepartureDate:time.Now(),
		},
		&models.ScheduleDtos{
			DepartureDate:time.Now(),
		},
	}
	rows := sqlmock.NewRows([]string{"departure_time"}).
		AddRow(mockSchedule[0].DepartureDate).
		AddRow(mockSchedule[1].DepartureDate)
	var transId []*string
	var ids = []string{"qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw"}
	for _,id := range ids{
		transId = append(transId,&id)
	}

	query := `SELECT distinct departure_date FROM schedules WHERE month =\? AND year = \?`
	if len(transId) != 0 {
		for index, id := range transId {
			if index == 0 && index != (len(transId)-1) {
				query = query + ` AND \(trans_id LIKE '%` + *id + `%' `
			} else if index == 0 && index == (len(transId)-1) {
				query = query + ` AND \(trans_id LIKE '%` + *id + `%' \) `
			} else if index == (len(transId) - 1) {
				query = query + ` OR  trans_id LIKE '%` + *id + `%' \) `
			} else {
				query = query + ` OR  trans_id LIKE '%` + *id + `%' `
			}
		}
	}

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)

	anArticle, err := a.GetScheduleByTransIds(context.TODO(), transId,2020,"June")
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}
//func TestGetScheduleByTransIdsNotFound(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//	}
//
//	//defer func() {
//	//	err = db.Close()
//	//	require.NoError(t, err)
//	//}()
//
//
//	rows := sqlmock.NewRows([]string{"departure_time"})
//	var transId []*string
//	var ids = []string{"qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw"}
//	for _,id := range ids{
//		transId = append(transId,&id)
//	}
//
//	query := `SELECT distinct departure_date FROM schedules WHERE month =\? AND year = \?`
//	if len(transId) != 0 {
//		for index, id := range transId {
//			if index == 0 && index != (len(transId)-1) {
//				query = query + ` AND \(trans_id LIKE '%` + *id + `%' `
//			} else if index == 0 && index == (len(transId)-1) {
//				query = query + ` AND \(trans_id LIKE '%` + *id + `%' \) `
//			} else if index == (len(transId) - 1) {
//				query = query + ` OR  trans_id LIKE '%` + *id + `%' \) `
//			} else {
//				query = query + ` OR  trans_id LIKE '%` + *id + `%' `
//			}
//		}
//	}
//
//	mock.ExpectQuery(query).WillReturnRows(rows)
//	a := ScheduleRepo.NewScheduleRepository(db)
//
//	anArticle, err := a.GetScheduleByTransIds(context.TODO(), transId,2020,"June")
//	assert.Error(t, err)
//	assert.Nil(t, anArticle)
//}
func TestGetScheduleByTransIdsErrorFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	mockSchedule := []*models.ScheduleDtos{
		&models.ScheduleDtos{
			DepartureDate:time.Now(),
		},
		&models.ScheduleDtos{
			DepartureDate:time.Now(),
		},
	}
	rows := sqlmock.NewRows([]string{"departure_time"}).
		AddRow(mockSchedule[0].DepartureDate).
		AddRow(1)
	var transId []*string
	var ids = []string{"qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw", "qrqeqweqweqw"}
	for _,id := range ids{
		transId = append(transId,&id)
	}

	query := `SELECT distinct departure_date FROM schedules WHERE month =\? AND year = \?`
	if len(transId) != 0 {
		for index, id := range transId {
			if index == 0 && index != (len(transId)-1) {
				query = query + ` AND \(trans_id LIKE '%` + *id + `%' `
			} else if index == 0 && index == (len(transId)-1) {
				query = query + ` AND \(trans_id LIKE '%` + *id + `%' \) `
			} else if index == (len(transId) - 1) {
				query = query + ` OR  trans_id LIKE '%` + *id + `%' \) `
			} else {
				query = query + ` OR  trans_id LIKE '%` + *id + `%' `
			}
		}
	}

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := ScheduleRepo.NewScheduleRepository(db)

	anArticle, err := a.GetScheduleByTransIds(context.TODO(), transId,2020,"June")
	assert.Error(t, err)
	assert.Nil(t, anArticle)
}
func TestDeleteByTransId(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "DELETE FROM schedules WHERE trans_id = \\?"

	transId := "qpoekqwpoekq"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(transId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ScheduleRepo.NewScheduleRepository(db)

	err = a.DeleteByTransId(context.TODO(), &transId)
	assert.NoError(t, err)
}
func TestDeleteByTransIdErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := "DELETE FROM schedules WHERE trans_id = \\?"

	transId := "qpoekqwpoekq"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(transId,transId).WillReturnResult(sqlmock.NewResult(2, 1))

	a := ScheduleRepo.NewScheduleRepository(db)

	err = a.DeleteByTransId(context.TODO(), &transId)
	assert.Error(t, err)
}
func TestInsert(t *testing.T) {
	//user := "test"
	//now := time.Now()
	mockScheduleTime := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	mockScheduleDay := []*models.ScheduleDay{
		&models.ScheduleDay{
			Year:          2020,
			Month:         "May",
			Day:           "Thursday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
		&models.ScheduleDay{
			Year:          2020,
			Month:         "June",
			Day:           "Monday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
	}
	a := models.Schedule{
		Id:                    "asdqweqwsad",
		CreatedBy:             "Test",
		CreatedDate:           time.Now(),
		ModifiedBy:            nil,
		ModifiedDate:          nil,
		DeletedBy:             nil,
		DeletedDate:           nil,
		IsDeleted:             0,
		IsActive:              1,
		TransId:               "asdqdsadqeqweqw",
		DepartureTime:         mockScheduleTime[0].DepartureTime,
		ArrivalTime:           mockScheduleTime[0].ArrivalTime,
		Day:                   mockScheduleDay[0].Day,
		Month:                 mockScheduleDay[0].Month,
		Year:                  mockScheduleDay[0].Year,
		DepartureDate:         "2020-09-17",
		Price:                 mockScheduleDay[0].Price,
		DepartureTimeoptionId: nil,
		ArrivalTimeoptionId:   nil,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT schedules SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , deleted_by=\? , 
				deleted_date=\? , is_deleted=\? , is_active=\? , trans_id=\?,arrival_time=\?,departure_time=\?,day=\?,
				month=\?,year=\?,departure_date=\?,price=\?,departure_timeoption_id=\?,arrival_timeoption_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.TransId, a.ArrivalTime,
		a.DepartureTime, a.Day, a.Month, a.Year, a.DepartureDate, a.Price, a.DepartureTimeoptionId, a.ArrivalTimeoptionId).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ScheduleRepo.NewScheduleRepository(db)

	id, err := i.Insert(context.TODO(), a)
	assert.NoError(t, err)
	assert.Equal(t, *id, a.Id)
}
func TestInsertErrorExec(t *testing.T) {
	mockScheduleTime := []*models.ScheduleTime{
		&models.ScheduleTime{
			DepartureTime: "07:00:00",
			ArrivalTime:   "08:00:00",
		},
		&models.ScheduleTime{
			DepartureTime: "09:00:00",
			ArrivalTime:   "10:00:00",
		},
	}
	mockScheduleDay := []*models.ScheduleDay{
		&models.ScheduleDay{
			Year:          2020,
			Month:         "May",
			Day:           "Thursday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
		&models.ScheduleDay{
			Year:          2020,
			Month:         "June",
			Day:           "Monday",
			DepartureDate: time.Now(),
			Price:         `{"adult_price":1000,"children_price":1000,"currency":0}`,
		},
	}
	a := models.Schedule{
		Id:                    "asdqweqwsad",
		CreatedBy:             "Test",
		CreatedDate:           time.Now(),
		ModifiedBy:            nil,
		ModifiedDate:          nil,
		DeletedBy:             nil,
		DeletedDate:           nil,
		IsDeleted:             0,
		IsActive:              1,
		TransId:               "asdqdsadqeqweqw",
		DepartureTime:         mockScheduleTime[0].DepartureTime,
		ArrivalTime:           mockScheduleTime[0].ArrivalTime,
		Day:                   mockScheduleDay[0].Day,
		Month:                 mockScheduleDay[0].Month,
		Year:                  mockScheduleDay[0].Year,
		DepartureDate:         "2020-09-17",
		Price:                 mockScheduleDay[0].Price,
		DepartureTimeoptionId: nil,
		ArrivalTimeoptionId:   nil,
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	//defer func() {
	//	err = db.Close()
	//	require.NoError(t, err)
	//}()

	query := `INSERT schedules SET id=\? , created_by=\? , created_date=\? , modified_by=\?, modified_date=\? , deleted_by=\? , 
				deleted_date=\? , is_deleted=\? , is_active=\? , trans_id=\?,arrival_time=\?,departure_time=\?,day=\?,
				month=\?,year=\?,departure_date=\?,price=\?,departure_timeoption_id=\?,arrival_timeoption_id=\?`
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.TransId, a.ArrivalTime,
		a.DepartureTime, a.Day, a.Month, a.Year, a.DepartureDate, a.Price, a.DepartureTimeoptionId, a.ArrivalTimeoptionId,a.Price).WillReturnResult(sqlmock.NewResult(1, 1))

	i := ScheduleRepo.NewScheduleRepository(db)

	_, err = i.Insert(context.TODO(), a)

	assert.Error(t, err)
}
