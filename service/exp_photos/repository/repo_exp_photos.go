package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"

	guuid "github.com/google/uuid"

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


// Newexp_photosRepository will create an object that represent the article.repository interface
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
			&t.ExpPhotoImage,
			&t.ExpId,
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
func (m *exp_photosRepository) GetByExperienceID(ctx context.Context, id string) (res []*models.ExpPhotos, err error) {
	query := `SELECT * FROM exp_photos WHERE exp_id = ? AND is_deleted = 0 AND is_active = 1`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list
	} else {
		return nil, nil
	}
	return
}

func (m *exp_photosRepository) Insert(ctx context.Context, a *models.ExpPhotos) (*string, error) {
	id := guuid.New()
	a.Id = id.String()
	query := `INSERT exp_photos SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , 
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , exp_photo_folder=?,exp_photo_image=?,
				exp_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.ExpPhotoFolder,
		a.ExpPhotoImage, a.ExpId)
	if err != nil {
		return nil, err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id, nil
}
func (m *exp_photosRepository) Update(ctx context.Context, a *models.ExpPhotos) (*string, error) {
	query := `UPDATE exp_photos SET modified_by=?, modified_date=? , 
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , exp_photo_folder=?,exp_photo_image=?,
				exp_id=? WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.ExpPhotoFolder,
		a.ExpPhotoImage, a.ExpId, a.Id)
	if err != nil {
		return nil, err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return &a.Id, nil
}
func (m *exp_photosRepository) Deletes(ctx context.Context, ids []string, expId string, deletedBy string) error {
	query := `UPDATE exp_photos SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE exp_id=?`
	for index, id := range ids {
		if index == 0 && index != (len(ids)-1) {
			query = query + ` AND (id !=` + id
		} else if index == 0 && index == (len(ids)-1) {
			query = query + ` AND (id !=` + id + ` ) `
		} else if index == (len(ids) - 1) {
			query = query + ` OR id !=` + id + ` ) `
		} else {
			query = query + ` OR id !=` + id
		}
	}
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0, expId)
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

func (m *exp_photosRepository) DeleteByExpId(ctx context.Context, expId string, deletedBy string) error {
	query := `UPDATE exp_photos SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE exp_id=?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0, expId)
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
