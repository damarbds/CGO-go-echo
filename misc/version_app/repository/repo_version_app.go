package repository

import (
	"context"
	"database/sql"
	"github.com/misc/version_app"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

type versionAPPRepository struct {
	Conn *sql.DB
}

// NewExpPaymentRepository will create an object that represent the exp_payment.repository interface
func NewVersionAPPRepositoryRepository(Conn *sql.DB) version_app.Repository {
	return &versionAPPRepository{Conn}
}

func (m *versionAPPRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.VersionApp, error) {
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

	result := make([]*models.VersionApp, 0)
	for rows.Next() {
		t := new(models.VersionApp)
		err = rows.Scan(
			&t.Id,
			&t.VersionCode,
			&t.VersionName,
			&t.Type,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *versionAPPRepository) GetAll(ctx context.Context, typeApp int) ([]*models.VersionApp, error) {
	query := `
	SELECT
		*
	FROM
		version_apps
	WHERE
		version_apps.type = ?`

	res, err := m.fetch(ctx, query, typeApp)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

