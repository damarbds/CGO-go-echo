package repository

import (
	"database/sql"
	"github.com/models"
	exp_availability "github.com/service/exp_availability"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type exp_availabilityRepository struct {
	Conn *sql.DB
}

// NewExpexp_availabilityRepository will create an object that represent the exp_exp_availability.Repository interface
func NewExpavailabilityRepository(Conn *sql.DB) exp_availability.Repository {
	return &exp_availabilityRepository{Conn}
}
func (m *exp_availabilityRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExpAvailability, error) {
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

	result := make([]*models.ExpAvailability, 0)
	for rows.Next() {
		t := new(models.ExpAvailability)
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
			&t.ExpAvailabilityMonth	,
			&t.ExpAvailabilityDate		,
			&t.ExpAvailabilityYear	 ,
			&t.ExpId				,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (e exp_availabilityRepository) GetByExpId(ctx context.Context, expId string) ([]*models.ExpAvailability, error) {
	query := `SELECT * FROM exp_availabilities WHERE exp_id = ?`

	list, err := e.fetch(ctx, query, expId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return list, nil
}

