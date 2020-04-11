package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/transportation"
	"github.com/sirupsen/logrus"
)

type transportationRepository struct {
	Conn *sql.DB
}

func NewTransportationRepository(Conn *sql.DB) transportation.Repository {
	return &transportationRepository{Conn:Conn}
}

func (t transportationRepository) List(ctx context.Context) ([]*models.TimesOption, error) {
	query := `SELECT * FROM times_options WHERE is_deleted = 0 AND is_active = 1`

	list, err := t.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return list, err
}

func (t transportationRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.TimesOption, error) {
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
			&t.StartTime,
			&t.EndTime,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}


