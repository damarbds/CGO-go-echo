package repository

import (
	"context"
	"database/sql"

	"github.com/misc/notif"
	"github.com/models"
	"github.com/sirupsen/logrus"
)

type notifRepository struct {
	Conn *sql.DB
}

func NewNotifRepository(Conn *sql.DB) notif.Repository {
	return &notifRepository{Conn: Conn}
}

func (n notifRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Notification, error) {
	rows, err := n.Conn.QueryContext(ctx, query, args...)
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

	result := make([]*models.Notification, 0)
	for rows.Next() {
		t := new(models.Notification)
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
			&t.MerchantId,
			&t.Type,
			&t.Title,
			&t.Desc,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (n notifRepository) GetByMerchantID(ctx context.Context, merchantId string) ([]*models.Notification, error) {
	query := `
	SELECT
		*
	FROM
		notifications
	WHERE
		merchant_id = ?
		AND is_deleted = 0
		AND is_active = 1`

	res, err := n.fetch(ctx, query, merchantId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}
