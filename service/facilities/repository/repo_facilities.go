package repository

import (
	"context"
	"database/sql"

	"github.com/models"
	"github.com/service/facilities"
	"github.com/sirupsen/logrus"
)

type facilityRepository struct {
	Conn *sql.DB
}

func NewFacilityRepository(Conn *sql.DB) facilities.Repository {
	return &facilityRepository{Conn: Conn}
}

func (f facilityRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Facilities, error) {
	rows, err := f.Conn.QueryContext(ctx, query, args...)
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

	result := make([]*models.Facilities, 0)
	for rows.Next() {
		t := new(models.Facilities)
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
			&t.FacilityName,
			&t.IsNumerable,
			&t.FacilityIcon,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (f facilityRepository) List(ctx context.Context) ([]*models.Facilities, error) {
	query := `SELECT * FROM facilities WHERE is_deleted = 0 and is_active = 1`

	res, err := f.fetch(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}
