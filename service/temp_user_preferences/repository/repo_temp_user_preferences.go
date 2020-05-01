package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/temp_user_preferences"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type tempUserPreferencesRepository struct {
	Conn *sql.DB
}

// NewexperienceRepository will create an object that represent the article.Repository interface
func NewtempUserPreferencesRepository(Conn *sql.DB) temp_user_preferences.Repository {
	return &tempUserPreferencesRepository{Conn}
}
func (m *tempUserPreferencesRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.TempUserPreference, error) {
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

	result := make([]*models.TempUserPreference, 0)
	for rows.Next() {
		t := new(models.TempUserPreference)
		err = rows.Scan(
			&t.Id,
			&t.ProvinceId ,
			&t.CityId ,
			&t.HarborsId ,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *tempUserPreferencesRepository) GetAll(ctx context.Context, page *int, size *int) ([]*models.TempUserPreference, error) {
	if page != nil && size != nil {
		query := `Select * FROM temp_user_preferences `

		query = query + ` LIMIT ? OFFSET ? `
		res, err := m.fetch(ctx, query, size, page)
		if err != nil {
			return nil, err
		}
		return res, err

	} else {
		query := `Select * FROM temp_user_preferences `

		//query = query + ` ORDER BY created_date desc `
		res, err := m.fetch(ctx, query)
		if err != nil {
			return nil, err
		}
		return res, err
	}
}
