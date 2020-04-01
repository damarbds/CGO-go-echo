package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	inspiration "github.com/service/exp_inspiration"
	"github.com/sirupsen/logrus"
)

type expInspirationRepository struct {
	Conn *sql.DB
}

// NewExpInspirationRepository will create an object that represent the exp_payment.Repository interface
func NewExpInspirationRepository(Conn *sql.DB) inspiration.Repository {
	return &expInspirationRepository{Conn}
}

func (m *expInspirationRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExpInspirationObject, error) {
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

	result := make([]*models.ExpInspirationObject, 0)
	for rows.Next() {
		t := new(models.ExpInspirationObject)
		err = rows.Scan(
			&t.ExpInspirationID,
			&t.ExpId,
			&t.ExpTitle,
			&t.ExpDesc,
			&t.ExpCoverPhoto,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (e expInspirationRepository) GetExpInspirations(ctx context.Context) ([]*models.ExpInspirationObject, error) {
	query := `
	SELECT
		id as exp_inspiration_id,
		exp_id,
		exp_title,
		exp_desc,
		exp_cover_photo
	FROM
		exp_inspirations
	WHERE
		is_deleted = 0
		AND is_active = 1`

	list, err := e.fetch(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return list, nil
}

