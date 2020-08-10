package repository

import (
	"context"
	"database/sql"
	"github.com/misc/notif"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"strconv"
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
			&t.ExpId ,
			&t.ScheduleId  ,
			&t.BookingExpId ,
			&t.IsRead ,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (n notifRepository) GetByMerchantID(ctx context.Context, merchantId string,limit,offset int) ([]*models.Notification, error) {
	query := `
	SELECT
		*
	FROM
		notifications
	WHERE
		merchant_id = ?
		AND is_deleted = 0
		AND is_active = 1 `

	if limit != 0 {
		query = query + ` ORDER BY created_date DESC LIMIT ` + strconv.Itoa(limit) +
			` OFFSET ` + strconv.Itoa(offset) + ` `
	}

	res, err := n.fetch(ctx, query, merchantId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

func (t notifRepository) GetCountByMerchantID(ctx context.Context, merchantId string) (int, error) {
	query := `SELECT count(*) as count FROM notifications WHERE is_deleted = 0 AND is_active = 1 `
	if merchantId != ""{
		query = query + ` AND merchant_id = '` + merchantId + `' `
	}
	rows, err := t.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	count, err := checkCount(rows)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}

	return count, nil
}
func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
func (m notifRepository) Insert(ctx context.Context, a models.Notification) error {
	//a.Id = guuid.New().String()
	query := `INSERT notifications SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? ,
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , merchant_id=?,type=? , title=? ,
				notifications.desc=?,exp_id=?,schedule_id=?,booking_exp_id=?,is_read=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.MerchantId,
			a.Type,a.Title,a.Desc,a.ExpId,a.ScheduleId,a.BookingExpId,a.IsRead)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	if err != nil {
		return err
	}

	//a.Id = lastID
	return nil
}
