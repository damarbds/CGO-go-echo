package repository

import (
	"context"
	"database/sql"
	"encoding/base64"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/models"
	"github.com/service/cpc"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type cpcRepository struct {
	Conn *sql.DB
}

// NewcpcRepository will create an object that represent the article.Repository interface
func NewcpcRepository(Conn *sql.DB) cpc.Repository {
	return &cpcRepository{Conn}
}

func (m *cpcRepository) fetchCity(ctx context.Context, query string, args ...interface{}) ([]*models.City, error) {
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

	result := make([]*models.City, 0)
	for rows.Next() {
		t := new(models.City)
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
			&t.CityName		,
			&t.CityDesc,
			&t.CityPhotos,
			&t.ProvinceId			,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *cpcRepository) FetchCity(ctx context.Context, cursor string, num int64) ([]*models.City, string, error) {
	query := `SELECT * FROM cities WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}

	res, err := m.fetchCity(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedDate)
	}

	return res, nextCursor, err
}
func (m *cpcRepository) GetCityByID(ctx context.Context, id int) (res *models.City, err error) {
	query := `SELECT * FROM cities WHERE id = ?`

	list, err := m.fetchCity(ctx, query, id)
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

func (m *cpcRepository) fetchProvince(ctx context.Context, query string, args ...interface{}) ([]*models.Province, error) {
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

	result := make([]*models.Province, 0)
	for rows.Next() {
		t := new(models.Province)
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
			&t.ProvinceName,
			&t.CountryId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *cpcRepository) FetchProvince(ctx context.Context, cursor string, num int64) ([]*models.Province, string, error) {
	query := `SELECT * FROM provinces WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}

	res, err := m.fetchProvince(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedDate)
	}

	return res, nextCursor, err
}
func (m *cpcRepository) GetProvinceByID(ctx context.Context, id int) (res *models.Province, err error) {
	query := `SELECT * FROM provinces WHERE id = ?`

	list, err := m.fetchProvince(ctx, query, id)
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
//func (m *cpcRepository) GetBycpcEmail(ctx context.Context, cpcEmail string) (res *models.cpc, err error) {
//	query := `SELECT * FROM cpcs WHERE cpc_email = ?`
//
//	list, err := m.fetch(ctx, query, cpcEmail)
//	if err != nil {
//		return
//	}
//
//	if len(list) > 0 {
//		res = list[0]
//	} else {
//		return nil, models.ErrNotFound
//	}
//	return
//}
//func (m *cpcRepository) Insert(ctx context.Context, a *models.cpc) error {
//	query := `INSERT cpcs SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , cpc_name=? , cpc_desc=? , cpc_email=? ,balance=?`
//	stmt, err := m.Conn.PrepareContext(ctx, query)
//	if err != nil {
//		return err
//	}
//	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.cpcName, a.cpcDesc,
//		a.cpcEmail, a.Balance)
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

func (m *cpcRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE  cpcs SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=?`
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

//func (m *cpcRepository) Update(ctx context.Context, ar *models.cpc) error {
//	query := `UPDATE cpcs set modified_by=?, modified_date=? , cpc_name=? ,
//				cpc_desc=? , cpc_email=? , balance=? WHERE id = ?`
//
//	stmt, err := m.Conn.PrepareContext(ctx, query)
//	if err != nil {
//		return nil
//	}
//
//	res, err := stmt.ExecContext(ctx, ar.ModifiedBy, time.Now(), ar.cpcName, ar.cpcDesc, ar.cpcEmail,
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
