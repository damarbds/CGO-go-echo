package repository

import (
	"context"
	"database/sql"
	"encoding/base64"

	"time"

	"github.com/sirupsen/logrus"

	"github.com/models"
	"github.com/service/harbors"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type harborsRepository struct {
	Conn *sql.DB
}


// NewharborsRepository will create an object that represent the article.Repository interface
func NewharborsRepository(Conn *sql.DB) harbors.Repository {
	return &harborsRepository{Conn}
}

func (m *harborsRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Harbors, error) {
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

	result := make([]*models.Harbors, 0)
	for rows.Next() {
		t := new(models.Harbors)
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
			&t.HarborsName,
			&t.HarborsLongitude,
			&t.HarborsLatitude,
			&t.HarborsImage,
			&t.CityId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *harborsRepository) fetchWithJoinCPC(ctx context.Context, query string, args ...interface{}) ([]*models.HarborsWCPC, error) {
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

	result := make([]*models.HarborsWCPC, 0)
	for rows.Next() {
		t := new(models.HarborsWCPC)
		err = rows.Scan(
			&t.Id,
			&t.HarborsName,
			&t.HarborsLongitude,
			&t.HarborsLatitude,
			&t.HarborsImage,
			&t.CityId,
			&t.CityName,
			&t.ProvinceName,
			&t.CountryName,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *harborsRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Harbors, string, error) {
	query := `SELECT * FROM harbors WHERE created_at > ? ORDER BY created_at LIMIT ? `

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

func (m *harborsRepository) GetAllWithJoinCPC(ctx context.Context,page *int , size *int,search string) ([]*models.HarborsWCPC, error) {

	search = "%" + search + "%"
	if page != nil && size != nil && search != ""{
		query := `Select h.id, h.harbors_name,h.harbors_longitude,h.harbors_latitude,h.harbors_image,h.city_id , c.city_name,p.province_name,co.country_name from cgo_indonesia.harbors h
			join cities c on h.city_id = c.id
			join provinces p on c.province_id = p.id
			join countries co on p.country_id = co.id
			where h.is_active = 1 and h.is_deleted = 0 
			AND (h.harbors_name LIKE ? OR c.city_name LIKE ? OR p.province_name LIKE ?) 
            ORDER BY h.created_date desc LIMIT ? OFFSET ? `

		res, err := m.fetchWithJoinCPC(ctx, query, search,search,search,page, size)
		if err != nil {
			return nil, err
		}
		return res, err

	} else if page == nil && size == nil && search != ""{
		query := `Select h.id, h.harbors_name,h.harbors_longitude,h.harbors_latitude,h.harbors_image,h.city_id , c.city_name,p.province_name,co.country_name from cgo_indonesia.harbors h
			join cities c on h.city_id = c.id
			join provinces p on c.province_id = p.id
			join countries co on p.country_id = co.id
			where h.is_active = 1 and h.is_deleted = 0 
			AND (h.harbors_name LIKE ? OR c.city_name LIKE ? OR p.province_name LIKE ?)`

		res, err := m.fetchWithJoinCPC(ctx, query, search,search,search)
		if err != nil {
			return nil, err
		}
		return res, err

	}else {
		query := `Select h.id, h.harbors_name,h.harbors_longitude,h.harbors_latitude,h.harbors_image,h.city_id , c.city_name,p.province_name,co.country_name from cgo_indonesia.harbors h
			join cities c on h.city_id = c.id
			join provinces p on c.province_id = p.id
			join countries co on p.country_id = co.id
			where h.is_active = 1 and h.is_deleted = 0`

		res, err := m.fetchWithJoinCPC(ctx, query)
		if err != nil {
			return nil, err
		}
		return res, err
	}


	return nil, nil
}
func (m *harborsRepository) GetByID(ctx context.Context, id string) (res *models.Harbors, err error) {
	query := `SELECT * FROM harbors WHERE id = ? AND is_deleted = 0  AND is_active = 1`

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
//func (m *harborsRepository) GetByharborsEmail(ctx context.Context, harborsEmail string) (res *models.Harbors, err error) {
//	query := `SELECT * FROM harborss WHERE harbors_email = ?`
//
//	list, err := m.fetch(ctx, query, harborsEmail)
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
//func (m *harborsRepository) Insert(ctx context.Context, a *models.Harbors) error {
//	query := `INSERT harborss SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , harbors_name=? , harbors_desc=? , harbors_email=? ,balance=?`
//	stmt, err := m.Conn.PrepareContext(ctx, query)
//	if err != nil {
//		return err
//	}
//	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.harborsName, a.harborsDesc,
//		a.harborsEmail, a.Balance)
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

func (m *harborsRepository) Delete(ctx context.Context, id string, deleted_by string) error {
	query := `UPDATE  harborss SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=?`
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
//func (m *harborsRepository) Update(ctx context.Context, ar *models.Harbors) error {
//	query := `UPDATE harborss set modified_by=?, modified_date=? , harbors_name=? ,
//				harbors_desc=? , harbors_email=? , balance=? WHERE id = ?`
//
//	stmt, err := m.Conn.PrepareContext(ctx, query)
//	if err != nil {
//		return nil
//	}
//
//	res, err := stmt.ExecContext(ctx, ar.ModifiedBy, time.Now(), ar.harborsName, ar.harborsDesc, ar.harborsEmail,
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
