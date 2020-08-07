package repository

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"

	"github.com/service/schedule"

	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type scheduleRepository struct {
	Conn *sql.DB
}

func (t scheduleRepository) GetBookingBySchedule(ctx context.Context, transId string, departureDate string, arrivalTime string, departureTime string) ([]*models.Schedule, error) {
	q := `
	SELECT s.* FROM booking_exps b
	JOIN schedules s ON b.schedule_id = s.id 
	WHERE b.trans_id = ? 
	AND s.departure_date = ? 
	AND s.arrival_time =? 
	AND s.departure_time =? `

	list, err := t.fetch(ctx, q, transId, departureDate, arrivalTime, departureTime)
	if err != nil {
		return nil, err
	}

	return list, err
}

// NewpromoRepository will create an object that represent the article.repository interface
func NewScheduleRepository(Conn *sql.DB) schedule.Repository {
	return &scheduleRepository{Conn}
}
func (t scheduleRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Schedule, error) {
	rows, err := t.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Schedule, 0)
	for rows.Next() {
		t := new(models.Schedule)
		err = rows.Scan(
			&t.Id,
			&t.CreatedBy,
			&t.CreatedDate,
			&t.ModifiedBy,
			&t.ModifiedDate,
			&t.DeletedBy,
			&t.DeletedDate,
			&t.IsDeleted,
			&t.IsActive,
			&t.TransId,
			&t.DepartureTime,
			&t.ArrivalTime,
			&t.Day,
			&t.Month,
			&t.Year,
			&t.DepartureDate,
			&t.Price,
			&t.DepartureTimeoptionId,
			&t.ArrivalTimeoptionId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (t scheduleRepository) fetchDtos(ctx context.Context, query string, args ...interface{}) ([]*models.ScheduleDtos, error) {
	rows, err := t.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.ScheduleDtos, 0)
	for rows.Next() {
		t := new(models.ScheduleDtos)
		err = rows.Scan(
			&t.DepartureDate,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (s scheduleRepository) Insert(ctx context.Context, a models.Schedule) (*string, error) {
	//a.Id = guuid.New().String()
	query := `INSERT schedules SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , trans_id=?,arrival_time=?,departure_time=?,day=?,
				month=?,year=?,departure_date=?,price=?,departure_timeoption_id=?,arrival_timeoption_id=?`
	stmt, err := s.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.TransId, a.ArrivalTime,
		a.DepartureTime, a.Day, a.Month, a.Year, a.DepartureDate, a.Price, a.DepartureTimeoptionId, a.ArrivalTimeoptionId)
	if err != nil {
		return nil, err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id, nil
}

func (t scheduleRepository) Update(ctx context.Context, a models.Schedule) error {
	query := `UPDATE schedules SET modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , trans_id=?,arrival_time=?,departure_time=?,day=?,
				month=?,year=?,departure_date=?,price=?,departure_timeoption_id=?,arrival_timeoption_id=?
				WHERE id =?`
	stmt, err := t.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.ModifiedBy, a.ModifiedDate, nil, nil, 0, 1, a.TransId, a.ArrivalTime,
		a.DepartureTime, a.Day, a.Month, a.Year, a.DepartureDate, a.Price, a.DepartureTimeoptionId, a.ArrivalTimeoptionId, a.Id)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return nil
}
func (s scheduleRepository) DeleteByTransId(ctx context.Context, transId *string) error {
	query := "DELETE FROM schedules WHERE trans_id = ? AND " +
		"id not in (SELECT schedule_id FROM booking_exps)"

	stmt, err := s.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, transId)
	if err != nil {

		return err
	}

	//rowsAfected, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}
	//
	//if rowsAfected != 1 {
	//	err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", rowsAfected)
	//	return err
	//}

	return nil
}

func (t scheduleRepository) GetScheduleByTransIds(ctx context.Context, transId []*string, year int, month string) ([]*models.ScheduleDtos, error) {
	res := make([]*models.ScheduleDtos, 0)
	query := `SELECT distinct departure_date FROM schedules WHERE month =? AND year = ?`
	if len(transId) != 0 {
		for index, id := range transId {
			if index == 0 && index != (len(transId)-1) {
				query = query + ` AND (trans_id LIKE '%` + *id + `%' `
			} else if index == 0 && index == (len(transId)-1) {
				query = query + ` AND (trans_id LIKE '%` + *id + `%' ) `
			} else if index == (len(transId) - 1) {
				query = query + ` OR  trans_id LIKE '%` + *id + `%' ) `
			} else {
				query = query + ` OR  trans_id LIKE '%` + *id + `%' `
			}
		}
		resp, err := t.fetchDtos(ctx, query, month, year)
		if err != nil {
			//if err == sql.ErrNoRows {
			//	return nil, models.ErrNotFound
			//}
			return nil, err
		}
		res = resp
	}

	return res, nil
}

func (t scheduleRepository) GetTimeByTransId(ctx context.Context, transId string) ([]*models.ScheduleTime, error) {
	query := `SELECT DISTINCT departure_time,arrival_time FROM schedules WHERE trans_id = ?`

	rows, err := t.Conn.QueryContext(ctx, query, transId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.ScheduleTime, 0)
	for rows.Next() {
		t := new(models.ScheduleTime)
		err = rows.Scan(
			&t.DepartureTime,
			&t.ArrivalTime,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (t scheduleRepository) GetYearByTransId(ctx context.Context, transId string, arrivalTime string, departureTime string) ([]*models.ScheduleYear, error) {
	query := `SELECT DISTINCT year FROM schedules WHERE trans_id = ? AND arrival_time = ? AND departure_time = ?`

	rows, err := t.Conn.QueryContext(ctx, query, transId, arrivalTime, departureTime)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.ScheduleYear, 0)
	for rows.Next() {
		t := new(models.ScheduleYear)
		err = rows.Scan(
			&t.Year,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (t scheduleRepository) GetMonthByTransId(ctx context.Context, transId string, year int, arrivalTime string, departureTime string) ([]*models.ScheduleMonth, error) {
	query := `SELECT DISTINCT year,month FROM schedules WHERE trans_id = ? AND year =? AND arrival_time = ? AND departure_time = ?`

	rows, err := t.Conn.QueryContext(ctx, query, transId, year, arrivalTime, departureTime)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.ScheduleMonth, 0)
	for rows.Next() {
		t := new(models.ScheduleMonth)
		err = rows.Scan(
			&t.Year,
			&t.Month,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (t scheduleRepository) GetDayByTransId(ctx context.Context, transId string, year int, month string, arrivalTime string, departureTime string) ([]*models.ScheduleDay, error) {
	query := `SELECT DISTINCT year,month,day,departure_date,price FROM schedules WHERE trans_id = ? AND year =? AND month=? AND arrival_time = ? AND departure_time = ?`

	rows, err := t.Conn.QueryContext(ctx, query, transId, year, month, arrivalTime, departureTime)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.ScheduleDay, 0)
	for rows.Next() {
		t := new(models.ScheduleDay)
		err = rows.Scan(
			&t.Year,
			&t.Month,
			&t.Day,
			&t.DepartureDate,
			&t.Price,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
func (m scheduleRepository) GetCountSchedule(ctx context.Context, merchantId string, date string) (int, error) {
	query := `
	SELECT
		count(DISTINCT b.booking_date,b.trans_id) AS count
	FROM
		transactions t
	JOIN booking_exps b on b.id = t.booking_exp_id OR b.order_id = t.order_id
	JOIN transportations trans on trans.id = b.trans_id
	JOIN merchants m on m.id = trans.merchant_id
	WHERE
		DATE (b.booking_date) = ? AND 
		t.status in (0,1,2,5) AND 
		trans.merchant_id = ? AND
		t.is_deleted = 0 AND
		t.is_active = 1 `
	//
	//for index, id := range transId {
	//	if index == 0 && index != (len(transId)-1) {
	//		query = query + ` AND (b.trans_id LIKE '%` + *id + `%' `
	//	} else if index == 0 && index == (len(transId)-1) {
	//		query = query + ` AND (b.trans_id LIKE '%` + *id + `%' ) `
	//	} else if index == (len(transId) - 1) {
	//		query = query + ` OR  b.trans_id LIKE '%` + *id + `%' ) `
	//	} else {
	//		query = query + ` OR  b.trans_id LIKE '%` + *id + `%' `
	//	}
	//}
	rows, err := m.Conn.QueryContext(ctx, query, date, merchantId)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	count, err := checkCount(rows)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return count, nil
}
