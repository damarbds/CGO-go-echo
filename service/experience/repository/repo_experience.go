package repository

import (
"context"
"database/sql"
"encoding/base64"

"time"

"github.com/sirupsen/logrus"

"github.com/models"
"github.com/service/experience"
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

func (m *experienceRepository) GetByCategoryID(ctx context.Context, categoryId int) ([]*models.ExpSearch, error) {
	query := `
	SELECT
		e.id,
		exp_title,
		exp_type,
		rating
	FROM
		filter_activity_types f
		JOIN experiences e ON f.exp_id = e.id
	WHERE
		f.is_deleted = 0
		AND f.is_active = 1
		AND exp_type_id = ?`

	res, err := m.fetchSearchExp(ctx, query, categoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return res, nil
}

func (m *experienceRepository) QueryFilterSearch(ctx context.Context, query string) ([]*models.ExpSearch, error) {
	res, err := m.fetchSearchExp(ctx, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}
	return res, nil
}

func (m *experienceRepository) SearchExp(ctx context.Context, harborID, cityID string) ([]*models.ExpSearch, error) {
	query := `
	SELECT
		exp.id,
		exp_title,
		exp_type,
		rating
	FROM
		experiences exp
		JOIN harbors ON harbors.id = harbors_id
		JOIN cities ON cities.id = harbors.city_id
	WHERE
		(harbors_id = ? OR harbors.city_id = ?)
		AND exp.is_deleted = 0
		AND exp.is_active = 1`

	res, err := m.fetchSearchExp(ctx, query, harborID, cityID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return res, nil
}



func (m *experienceRepository) fetchUserDiscoverPreference(ctx context.Context, query string, args ...interface{}) ([]*models.ExpUserDiscoverPreference, error) {
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

	result := make([]*models.ExpUserDiscoverPreference, 0)
	for rows.Next() {
		t := new(models.ExpUserDiscoverPreference)
		err = rows.Scan(
			&t.CityId,
			&t.CityName,
			&t.CityDesc,
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

func (m *experienceRepository) fetchSearchExp(ctx context.Context, query string, args ...interface{}) ([]*models.ExpSearch, error) {
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

	result := make([]*models.ExpSearch, 0)
	for rows.Next() {
		t := new(models.ExpSearch)
		err = rows.Scan(
			&t.Id,
			&t.ExpTitle,
			&t.ExpType,
			&t.Rating,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
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

func (m *experienceRepository) GetIdByHarborsId(ctx context.Context, harborsId string) ([]*string, error) {
	query := `select id FROM experiences where is_deleted = 0 AND is_active = 1 AND status = 2 AND harbors_id = ?`

	rows, err := m.Conn.QueryContext(ctx, query, harborsId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	result := make([]*string, 0)
	for rows.Next() {
		t := new(string)
		err = rows.Scan(
			&t,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, err
}

func (m *experienceRepository) GetIdByCityId(ctx context.Context, cityId string) ([]*string, error) {
	query := `select e.id from experiences e 
			  join harbors h on h.id = e.harbors_id where e.is_active = 1 and e.is_deleted = 0 and
              e.status = 2 and h.city_id = ?`

	rows, err := m.Conn.QueryContext(ctx, query, cityId)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	result := make([]*string, 0)
	for rows.Next() {
		t := new(string)
		err = rows.Scan(
			&t,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, err
}
func (m *experienceRepository) GetUserDiscoverPreference(ctx context.Context,page *int,size *int) ([]*models.ExpUserDiscoverPreference, error) {

	if page != nil && size != nil{
		query := `select c.city_id, c.city_name, city.city_desc,a.* from cgo_indonesia.experiences a
			join cgo_indonesia.harbors b on b.id = a.harbors_id
			join 
			(
				select c.id as city_id, c.city_name as city_name from cgo_indonesia.user_preference_exps upe
				join cgo_indonesia.harbors h on h.id = upe.harbors_id
				join cgo_indonesia.cities c on c.id = h.city_id
				where upe.is_active = 1 and upe.is_deleted = 0
				order by upe.amount desc LIMIT ? OFFSET ? 
			) c on c.city_id = b.city_id
            join cgo_indonesia.cities city on city.id = c.city_id;`

		res, err := m.fetchUserDiscoverPreference(ctx, query,page,size)
		if err != nil {
			return nil, err
		}
		return res, err
	}else {
		query := `select c.city_id, c.city_name, city.city_desc,a.* from cgo_indonesia.experiences a
			join cgo_indonesia.harbors b on b.id = a.harbors_id
			join 
			(
				select c.id as city_id, c.city_name as city_name from cgo_indonesia.user_preference_exps upe
				join cgo_indonesia.harbors h on h.id = upe.harbors_id
				join cgo_indonesia.cities c on c.id = h.city_id
				where upe.is_active = 1 and upe.is_deleted = 0
				order by upe.amount desc  
			) c on c.city_id = b.city_id
            join cgo_indonesia.cities city on city.id = c.city_id;`

		res, err := m.fetchUserDiscoverPreference(ctx, query)
		if err != nil {
			return nil, err
		}
		return res, err
	}

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
