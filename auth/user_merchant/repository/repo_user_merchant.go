package repository

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/auth/user_merchant"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type userMerchantRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewuserMerchantRepository(Conn *sql.DB) user_merchant.Repository {
	return &userMerchantRepository{Conn}
}
func (m *userMerchantRepository) fetchWithMerchant(ctx context.Context, query string, args ...interface{}) ([]*models.UserMerchantWithMerchant, error) {
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

	result := make([]*models.UserMerchantWithMerchant, 0)
	for rows.Next() {
		t := new(models.UserMerchantWithMerchant)
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
			&t.FullName,
			&t.Email,
			&t.PhoneNumber,
			&t.MerchantId,
			&t.FCMToken,
			&t.MerchantName,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *userMerchantRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.UserMerchant, error) {
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

	result := make([]*models.UserMerchant, 0)
	for rows.Next() {
		t := new(models.UserMerchant)
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
			&t.FullName,
			&t.Email,
			&t.PhoneNumber,
			&t.MerchantId,
			&t.FCMToken,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *userMerchantRepository) Fetch(ctx context.Context, cursor string, num int64) (res []*models.UserMerchant, nextCursor string, err error) {
	query := `SELECT * FROM user_merchants WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor = ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedDate)
	}

	return res, nextCursor, err
}

func (m *userMerchantRepository) GetByID(ctx context.Context, id string) (res *models.UserMerchant, err error) {
	query := `SELECT * FROM user_merchants WHERE is_deleted = 0 and is_active = 1 and id = ?`

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
func (m *userMerchantRepository) GetUserByMerchantId(ctx context.Context, merchantId string) (res []*models.UserMerchant, err error) {
	query := `SELECT * FROM user_merchants WHERE is_deleted = 0 and is_active = 1 and merchant_id = ?`

	list, err := m.fetch(ctx, query, merchantId)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list
	} else {
		return nil, models.ErrNotFound
	}

	return
}
func (m *userMerchantRepository) GetByUserEmail(ctx context.Context, userEmail string) (res *models.UserMerchant, err error) {
	query := `SELECT * FROM user_merchants WHERE is_deleted = 0 AND is_active = 1 AND email = ?`

	list, err := m.fetch(ctx, query, userEmail)
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
func (m *userMerchantRepository) UpdateFCMToken(ctx context.Context, tokenFCM string, usermerchantId string) error {
	query := `UPDATE user_merchants set fcm_token=?
				WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, tokenFCM,usermerchantId)
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
func (m *userMerchantRepository) Update(ctx context.Context, a *models.UserMerchant) error {
	query := `UPDATE user_merchants set modified_by=?, modified_date=? , full_name=?,email=? , phone_number=? ,merchant_id=?
				WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.FullName, a.Email, a.PhoneNumber, a.MerchantId, a.Id)
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

func (m *userMerchantRepository) Insert(ctx context.Context, a *models.UserMerchant) error {
	query := `INSERT user_merchants SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? ,
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , full_name=?,email=? , phone_number=? ,merchant_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.FullName, a.Email, a.PhoneNumber, a.MerchantId)
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

func (m *userMerchantRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE user_merchants SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deleted_by, time.Now(), 1, 0, id)
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

func (m *userMerchantRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM user_merchants WHERE is_deleted = 0 and is_active = 1`

	rows, err := m.Conn.QueryContext(ctx, query)
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

func (m userMerchantRepository) List(ctx context.Context, limit, offset int, search string) ([]*models.UserMerchantWithMerchant, error) {
	query := `SELECT um.* ,m.merchant_name
					FROM 
						user_merchants um 
					JOIN merchants m on m.id = um.merchant_id
					WHERE um.is_deleted = 0 and um.is_active = 1 `
	if search != "" {
		query = query + ` AND ( um.email LIKE '%` + search + `%' ` +
			`OR um.full_name LIKE '%` + search + `%' ` +
			`OR um.phone_number LIKE '%` + search + `%' )`
	}
	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetchWithMerchant(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
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
