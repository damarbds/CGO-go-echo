package repository

import (
	"context"
	"database/sql"
	"github.com/misc/currency"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type exChangeRatesRepository struct {
	Conn *sql.DB
}


// NewexperienceRepository will create an object that represent the article.repository interface
func NewExChangeRatesRepository(Conn *sql.DB) currency.Repository {
	return &exChangeRatesRepository{Conn}
}

func (m *exChangeRatesRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExChangeRate, error) {
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

	result := make([]*models.ExChangeRate, 0)
	for rows.Next() {
		t := new(models.ExChangeRate)
		err = rows.Scan(
			&t.Id,
			&t.Date,
			&t.From,
			&t.To,
			&t.Rates,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}


func (m exChangeRatesRepository) GetByDate(ctx context.Context, from, to string) (res *models.ExChangeRate,err error) {
	query := `SELECT e.* FROM ex_change_rates e WHERE e.from = ? AND e.to =? order by e.date desc LIMIT 1`

	list, err := m.fetch(ctx, query, from,to)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}