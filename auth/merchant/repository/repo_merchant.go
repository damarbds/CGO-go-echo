package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/auth/merchant"
	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type merchantRepository struct {
	Conn *sql.DB
}

// NewmerchantRepository will create an object that represent the article.Repository interface
func NewmerchantRepository(Conn *sql.DB) merchant.Repository {
	return &merchantRepository{Conn}
}

func (m *merchantRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Merchant, error) {
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

	result := make([]*models.Merchant, 0)
	for rows.Next() {
		t := new(models.Merchant)
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
			&t.MerchantName,
			&t.MerchantDesc,
			&t.MerchantEmail,
			&t.Balance,
			&t.PhoneNumber,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *merchantRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Merchant, string, error) {
	query := `SELECT * FROM merchants WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}

	res, err := m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedDate)
	}

	return res, nextCursor, err
}
func (m *merchantRepository) GetByID(ctx context.Context, id string) (res *models.Merchant, err error) {
	query := `SELECT * FROM merchants WHERE id = ?`

	list, err := m.fetch(ctx, query, id)
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
func (m *merchantRepository) GetByMerchantEmail(ctx context.Context, merchantEmail string) (res *models.Merchant, err error) {
	query := `SELECT * FROM merchants WHERE merchant_email = ?`

	list, err := m.fetch(ctx, query, merchantEmail)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}
	return
}
func (m *merchantRepository) Insert(ctx context.Context, a *models.Merchant) error {
	query := `INSERT merchants SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , merchant_name=? , merchant_desc=? , merchant_email=? ,balance=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx,a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.MerchantName, a.MerchantDesc,
		a.MerchantEmail, a.Balance)
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

func (m *merchantRepository) Delete(ctx context.Context, id string,deleted_by string) error {
	query := `UPDATE  merchants SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deleted_by, time.Now(),1,0)
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
func (m *merchantRepository) Update(ctx context.Context, ar *models.Merchant) error {
	query := `UPDATE merchants set modified_by=?, modified_date=? , merchant_name=? , 
				merchant_desc=? , merchant_email=? , balance=? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, ar.ModifiedBy, time.Now(), ar.MerchantName, ar.MerchantDesc, ar.MerchantEmail,
								ar.Balance,ar.Id)
	if err != nil {
		return err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)

		return err
	}

	return nil
}

// DecodeCursor will decode cursor from user for mysql
func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

// EncodeCursor will encode cursor from mysql to user
func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
