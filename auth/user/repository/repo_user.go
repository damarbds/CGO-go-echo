package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/auth/user"
	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type userRepository struct {
	Conn *sql.DB
}



// NewuserRepository will create an object that represent the article.repository interface
func NewuserRepository(Conn *sql.DB) user.Repository {
	return &userRepository{Conn}
}


func (m *userRepository) UpdatePointByID(ctx context.Context, point float64, id string) error {
	query := `UPDATE users SET points = points - ? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	_, err = stmt.ExecContext(ctx, point, id)
	if err != nil {
		return err
	}
	//affect, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}
	//if affect != 1 {
	//	err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)
	//
	//	return err
	//}

	return nil
}

func (m *userRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM users WHERE is_deleted = 0 and is_active = 1`

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

func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}

func (m *userRepository) List(ctx context.Context, limit, offset int, search string) ([]*models.User, error) {
	query := `SELECT * FROM users WHERE is_deleted = 0 and is_active = 1 `
	if search != "" {
		query = query + ` AND ( user_email LIKE '%` + search + `%' ` +
			`OR full_name LIKE '%` + search + `%' ` +
			`OR phone_number LIKE '%` + search + `%' ` +
			`OR address LIKE '%` + search + `%' ` +
			`OR dob LIKE '%` + search + `%' ` +
			`OR points LIKE '%` + search + `%' )`
	}
	query = query + ` LIMIT ? OFFSET ?`
	list, err := m.fetch(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (m *userRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.User, error) {
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

	result := make([]*models.User, 0)
	for rows.Next() {
		t := new(models.User)
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
			&t.UserEmail,
			&t.FullName,
			&t.PhoneNumber,
			&t.VerificationSendDate,
			&t.VerificationCode,
			&t.ProfilePictUrl,
			&t.Address,
			&t.Dob,
			&t.Gender,
			&t.IdType,
			&t.IdNumber,
			&t.ReferralCode,
			&t.Points,
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

func (m *userRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.User, string, error) {
	query := `SELECT * FROM users WHERE created_at > ? ORDER BY created_at LIMIT ? `

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
func (m *userRepository) GetByID(ctx context.Context, id string) (res *models.User, err error) {
	query := `SELECT * FROM users WHERE id = ?`

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

func (m *userRepository) GetCreditByID(ctx context.Context, id string) (int, error) {
	var points int
	query := `SELECT points FROM users WHERE id = ?`

	err := m.Conn.QueryRowContext(ctx, query, id).Scan(&points)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, models.ErrNotFound
		}
		return 0, err
	}

	return points, err
}

func (m *userRepository) GetByUserEmail(ctx context.Context, userEmail string) (res *models.User, err error) {
	query := `SELECT * FROM users WHERE is_deleted = 0 AND is_active = 1 AND  user_email = ?`

	list, err := m.fetch(ctx, query, userEmail)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, nil
	}
	return
}

func (m *userRepository) GetByUserNumberOTP(ctx context.Context, phoneNumber string, otp string) (res *models.User, err error) {
	if otp == "" {
		query := `SELECT * FROM users WHERE phone_number = ?`

		list, err := m.fetch(ctx, query, phoneNumber)
		if err != nil {
			return nil, err
		}

		if len(list) > 0 {
			res = list[0]
		} else {
			return nil, models.ErrNotFound
		}
	} else {
		query := `SELECT * FROM users WHERE phone_number = ? AND verification_code =?`

		list, err := m.fetch(ctx, query, phoneNumber, otp)
		if err != nil {
			return nil, err
		}

		if len(list) > 0 {
			res = list[0]
		} else {
			return nil, models.ErrNotFound
		}
	}
	return
}

func (m *userRepository) Insert(ctx context.Context, a *models.User) error {
	query := `INSERT users SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , user_email=? , full_name=? , phone_number=? ,verification_send_date=?,verification_code=?,profile_pict_url=?,address=?,dob=?,gender=?,id_type=?,id_number=?,referral_code=?,points=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.UserEmail, a.FullName,
		a.PhoneNumber, a.VerificationSendDate, a.VerificationCode, a.ProfilePictUrl, a.Address, a.Dob, a.Gender, a.IdType, a.IdNumber, a.ReferralCode, a.Points)
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

func (m *userRepository) SubscriptionUser(ctx context.Context, s *models.Subscribe) error {
	query := `INSERT subscribes SET  created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , subscriber_email=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, s.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, s.SubscriberEmail)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	if err != nil {
		return err
	}

	//sendingEmail := models.SendingEmail{
	//	Subject:           "Subscribe",
	//	Message:           "Push",
	//	AttachmentFileUrl: "",
	//	FileName:          "",
	//	From:              "cgo indonesia",
	//	To:                email,
	//}
	//_,err := s.isUsecase.SendingEmail(&sendingEmail)

	//a.Id = lastID
	return nil
}

func (m *userRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE users SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
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
func (m *userRepository) Update(ctx context.Context, a *models.User) error {
	query := `UPDATE users set modified_by=?, modified_date=? , user_email=? , full_name=? , phone_number=? ,verification_send_date=?,verification_code=?,profile_pict_url=?,address=?,dob=?,gender=?,id_type=?,id_number=?,referral_code=?,points=? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	res, err := stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.UserEmail, a.FullName,
		a.PhoneNumber, a.VerificationSendDate, a.VerificationCode, a.ProfilePictUrl, a.Address, a.Dob, a.Gender, a.IdType, a.IdNumber, a.ReferralCode, a.Points, a.Id)
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
