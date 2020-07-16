package repository

import (
	"context"
	"database/sql"

	"github.com/models"
	types "github.com/service/exp_types"
	"github.com/sirupsen/logrus"
)

type expTypeRepository struct {
	Conn *sql.DB
}

// NewExpTypeRepository will create an object that represent the exp_payment.repository interface
func NewExpTypeRepository(Conn *sql.DB) types.Repository {
	return &expTypeRepository{Conn}
}

func (m *expTypeRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExpTypeObject, error) {
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

	result := make([]*models.ExpTypeObject, 0)
	for rows.Next() {
		t := new(models.ExpTypeObject)
		err = rows.Scan(
			&t.ExpTypeID,
			&t.ExpTypeName,
			&t.ExpTypeIcon,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (e expTypeRepository) GetExpTypes(ctx context.Context) ([]*models.ExpTypeObject, error) {
	query := `
	SELECT
		id AS exp_type_id,
		exp_type_name,
		COALESCE(exp_type_icon,"") AS exp_type_icon
	FROM
		experience_types
	WHERE
		is_active = 1
		AND is_deleted = 0`

	list, err := e.fetch(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return list, nil
}
func (e *expTypeRepository) GetByName(ctx context.Context, name string) (*models.ExpTypeObject, error) {
	query := `
	SELECT
		id AS exp_type_id,
		exp_type_name,
		COALESCE(exp_type_icon,"") AS exp_type_icon
	FROM
		experience_types
	WHERE
		is_active = 1
		AND is_deleted = 0 AND exp_type_name = ?`

	list, err := e.fetch(ctx, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	if len(list) == 0 {
		return nil, models.ErrNotFound
	}
	return list[0], nil
}
