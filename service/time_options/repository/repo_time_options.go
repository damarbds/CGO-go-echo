package repository

import (
	"database/sql"
	"github.com/models"
	"github.com/service/time_options"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"time"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type timeOptionsRepository struct {
	Conn *sql.DB
}


// NewpromoRepository will create an object that represent the article.Repository interface
func NewTimeOptionsRepository(Conn *sql.DB) time_options.Repository {
	return &timeOptionsRepository{Conn}
}


func (t timeOptionsRepository) Insert(ctx context.Context, a models.TimesOption) (*int, error) {
	query := `INSERT times_options SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , start_time=?,end_time=?`
	stmt, err := t.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil,err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.StartTime,a.EndTime)
	if err != nil {
		return nil,err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id,nil
}
func (m *timeOptionsRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.TimesOption, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
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

	result := make([]*models.TimesOption, 0)
	for rows.Next() {
		t := new(models.TimesOption)
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
			&t.StartTime ,
			&t.EndTime ,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (t timeOptionsRepository) GetByTime(ctx context.Context, time string) (*models.TimesOption, error) {
	query := `SELECT * FROM times_options where ? >= start_time  AND  ? <= end_time `

	res, err := t.fetch(ctx, query, time,time)
	if err != nil {
		return nil, err
	}

	return res[0], nil
}