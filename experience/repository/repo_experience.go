package repository

import (
	"context"
	"database/sql"
	"encoding/base64"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/experience"
	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type experienceRepository struct {
	Conn *sql.DB
}

// NewexperienceRepository will create an object that represent the article.Repository interface
func NewexperienceRepository(Conn *sql.DB) experience.Repository {
	return &experienceRepository{Conn}
}

func (m *experienceRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Experience, error) {
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

	result := make([]*models.Experience, 0)
	for rows.Next() {
		t := new(models.Experience)
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
			&t.ExpTitle,
			&t.ExpType,
			&t.ExpTripType,
			&t.ExpBookingType,
			&t.ExpDesc,
			&t.ExpMaxGuest,
			&t.ExpPickupPlace,
			&t.ExpPickupTime,
			&t.ExpPickupPlaceLongitude,
			&t.ExpPickupPlaceLatitude,
			&t.ExpPickupPlaceMapsName,
			&t.ExpInternary,
			&t.ExpFacilities,
			&t.ExpInclusion,
			&t.ExpRules,
			&t.Status,
			&t.Rating,
			&t.ExpLocationLatitude,
			&t.ExpLocationLongitude,
			&t.ExpLocationName,
			&t.ExpCoverPhoto,
			&t.ExpDuration,
			&t.MinimumBookingId,
			&t.MerchantId,
			&t.HarborsId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *experienceRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Experience, string, error) {
	query := `SELECT * FROM experiences WHERE created_at > ? ORDER BY created_at LIMIT ? `

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
func (m *experienceRepository) GetByID(ctx context.Context, id string) (res *models.Experience, err error) {
	query := `SELECT * FROM experiences WHERE id = ? AND is_deleted = 0 AND is_active = 1`

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
func (m *experienceRepository) GetByExperienceEmail(ctx context.Context, experienceEmail string) (res *models.Experience, err error) {
	query := `SELECT * FROM experiences WHERE experience_email = ?`

	list, err := m.fetch(ctx, query, experienceEmail)
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
//func (m *experienceRepository) Insert(ctx context.Context, a *models.Experience) error {
//	query := `INSERT experiences SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , experience_name=? , experience_desc=? , experience_email=? ,balance=?`
//	stmt, err := m.Conn.PrepareContext(ctx, query)
//	if err != nil {
//		return err
//	}
//	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.experienceName, a.experienceDesc,
//		a.experienceEmail, a.Balance)
//	if err != nil {
//		return err
//	}
//
//	//lastID, err := res.RowsAffected()
//	if err != nil {
//		return err
//	}
//
//	//a.Id = lastID
//	return nil
//}

func (m *experienceRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE  experiences SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deleted_by, time.Now(), 1, 0)
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
//func (m *experienceRepository) Update(ctx context.Context, ar *models.Experience) error {
//	query := `UPDATE experiences set modified_by=?, modified_date=? , experience_name=? ,
//				experience_desc=? , experience_email=? , balance=? WHERE id = ?`
//
//	stmt, err := m.Conn.PrepareContext(ctx, query)
//	if err != nil {
//		return nil
//	}
//
//	res, err := stmt.ExecContext(ctx, ar.ModifiedBy, time.Now(), ar.experienceName, ar.experienceDesc, ar.experienceEmail,
//		ar.Balance, ar.Id)
//	if err != nil {
//		return err
//	}
//	affect, err := res.RowsAffected()
//	if err != nil {
//		return err
//	}
//	if affect != 1 {
//		err = fmt.Errorf("Weird  Behaviour. Total Affected: %d", affect)
//
//		return err
//	}
//
//	return nil
//}

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
