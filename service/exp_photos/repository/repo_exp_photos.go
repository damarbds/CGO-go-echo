package repository

import (
	"context"
	"database/sql"
	"encoding/base64"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/models"
	"github.com/service/exp_photos"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type exp_photosRepository struct {
	Conn *sql.DB
}

// Newexp_photosRepository will create an object that represent the article.Repository interface
func Newexp_photosRepository(Conn *sql.DB) exp_photos.Repository {
	return &exp_photosRepository{Conn}
}

func (m *exp_photosRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExpPhotos, error) {
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

	result := make([]*models.ExpPhotos, 0)
	for rows.Next() {
		t := new(models.ExpPhotos)
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
			&t.ExpPhotoFolder,
			&t.ExpPhotoImage ,
			&t.ExpId	,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *exp_photosRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.ExpPhotos, string, error) {
	query := `SELECT * FROM exp_photos WHERE created_at > ? ORDER BY created_at LIMIT ? `

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
func (m *exp_photosRepository) GetByID(ctx context.Context, id string) (res *models.ExpPhotos, err error) {
	query := `SELECT * FROM exp_photos WHERE id = ? AND is_deleted = 0 AND is_active = 1`

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
func (m *exp_photosRepository) GetByExperienceID(ctx context.Context, id string) (res[] *models.ExpPhotos, err error) {
	query := `SELECT * FROM exp_photos WHERE exp_id = ? AND is_deleted = 0 AND is_active = 1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list
	} else {
		return nil, models.ErrNotFound
	}
	return
}

//func (m *exp_photosRepository) Insert(ctx context.Context, a *models.exp_photos) error {
//	query := `INSERT exp_photoss SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , exp_photos_name=? , exp_photos_desc=? , exp_photos_email=? ,balance=?`
//	stmt, err := m.Conn.PrepareContext(ctx, query)
//	if err != nil {
//		return err
//	}
//	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.exp_photosName, a.exp_photosDesc,
//		a.exp_photosEmail, a.Balance)
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

func (m *exp_photosRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE  exp_photoss SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=?`
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

//func (m *exp_photosRepository) Update(ctx context.Context, ar *models.exp_photos) error {
//	query := `UPDATE exp_photoss set modified_by=?, modified_date=? , exp_photos_name=? ,
//				exp_photos_desc=? , exp_photos_email=? , balance=? WHERE id = ?`
//
//	stmt, err := m.Conn.PrepareContext(ctx, query)
//	if err != nil {
//		return nil
//	}
//
//	res, err := stmt.ExecContext(ctx, ar.ModifiedBy, time.Now(), ar.exp_photosName, ar.exp_photosDesc, ar.exp_photosEmail,
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
