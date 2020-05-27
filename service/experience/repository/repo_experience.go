package repository

import (
	"context"
	"database/sql"
	"encoding/base64"

	guuid "github.com/google/uuid"

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

// NewexperienceRepository will create an object that represent the article.repository interface
func NewexperienceRepository(Conn *sql.DB) experience.Repository {
	return &experienceRepository{Conn}
}

func (m *experienceRepository) UpdateRating(ctx context.Context, exp models.Experience) error {
	query := `UPDATE experiences SET modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , rating=?,guide_review=?,activities_review=?,service_review =?,
				cleanliness_review=?,value_review=? WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, exp.ModifiedBy, time.Now(), nil, nil, 0, 1, exp.Rating, exp.GuideReview, exp.ActivitiesReview,
		exp.ServiceReview, exp.CleanlinessReview, exp.ValueReview, exp.Id)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return nil
}
func (m *experienceRepository) GetExpCount(ctx context.Context, merchantId string) (int, error) {
	query := `
	SELECT
		count(*) AS count
	FROM
		experiences
	WHERE
		merchant_id = ?`

	rows, err := m.Conn.QueryContext(ctx, query, merchantId)
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

func (m *experienceRepository) GetExpPendingTransactionCount(ctx context.Context, merchantId string) (int, error) {
	query := `
	SELECT
		count(*) AS count
	FROM
		experiences e
		JOIN booking_exps b ON e.id = b.exp_id
		JOIN transactions t ON b.id = t.booking_exp_id
	WHERE
		merchant_id = ?
		and(t.status = 0
			OR t.status IS NULL)
		AND b.status = 0`

	rows, err := m.Conn.QueryContext(ctx, query, merchantId)
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

func (m *experienceRepository) GetExpFailedTransactionCount(ctx context.Context, merchantId string) (int, error) {
	query := `
	SELECT
		count(*) AS count
	FROM
		experiences e
		JOIN booking_exps b ON e.id = b.exp_id
		JOIN transactions t ON b.id = t.booking_exp_id
	WHERE
		merchant_id = ?
		and t.status = 4`

	rows, err := m.Conn.QueryContext(ctx, query, merchantId)
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

func (m *experienceRepository) GetPublishedExpCount(ctx context.Context, merchantId string) (int, error) {
	query := `
	SELECT
		count(*) AS count
	FROM
		experiences
	WHERE
		merchant_id = ?
		AND status = 2`

	rows, err := m.Conn.QueryContext(ctx, query, merchantId)
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

func (m *experienceRepository) GetSuccessBookCount(ctx context.Context, merchantId string) (int, error) {
	query := `
	SELECT
		COUNT(*) AS count
	FROM
		experiences e
		LEFT JOIN booking_exps b ON e.id = b.exp_id
	WHERE
		e.merchant_id = ?
		AND b.status = 1`

	rows, err := m.Conn.QueryContext(ctx, query, merchantId)
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

func (m *experienceRepository) GetByCategoryID(ctx context.Context, categoryId int) ([]*models.ExpSearch, error) {
	query := `
	SELECT
		e.id,
		exp_title,
		exp_type,
		e.exp_location_latitude as latitude ,
		e.exp_location_longitude as longitude, 
		rating,
		exp_cover_photo as cover_photo
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

func (m *experienceRepository) CountFilterSearch(ctx context.Context, query string) (int, error) {
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

func (m *experienceRepository) QueryFilterSearch(ctx context.Context, query string, limit, offset int) ([]*models.ExpSearch, error) {
	query = query + ` LIMIT ? OFFSET ?`
	res, err := m.fetchSearchExp(ctx, query, limit, offset)
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
		exp.status as exp_status,
		exp_location_latitude as latitude ,
		exp_location_longitude as longitude, 
		rating,
		exp_cover_photo as cover_photo,
		p.province_name as province
	FROM
		experiences exp
		JOIN harbors ON harbors.id = harbors_id
		JOIN cities ON cities.id = harbors.city_id
		JOIN provinces p ON cities.province_id = p.id
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
			&t.ProvinceId,
			&t.ProvinceName,
			&t.CityId,
			&t.CityName,
			&t.CityDesc,
			&t.CityPhotos,
			&t.IdHarbors,
			&t.HarborsName,
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
			&t.GuideReview,
			&t.ActivitiesReview,
			&t.ServiceReview,
			&t.CleanlinessReview,
			&t.ValueReview,
			&t.ExpPaymentDeadlineAmount,
			&t.ExpPaymentDeadlineType,
			&t.IsCustomisedByUser,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *experienceRepository) fetchUserDiscoverPreferenceProvince(ctx context.Context, query string, args ...interface{}) ([]*models.ExpUserDiscoverPreference, error) {
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
			&t.ProvinceId,
			&t.ProvinceName,
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
			&t.GuideReview,
			&t.ActivitiesReview,
			&t.ServiceReview,
			&t.CleanlinessReview,
			&t.ValueReview,
			&t.ExpPaymentDeadlineAmount,
			&t.ExpPaymentDeadlineType,
			&t.IsCustomisedByUser,
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
			&t.ExpStatus,
			&t.Rating,
			&t.Latitude,
			&t.Longitude,
			&t.CoverPhoto,
			&t.Province,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m *experienceRepository) fetchJoinForegnKey(ctx context.Context, query string, args ...interface{}) ([]*models.ExperienceJoinForegnKey, error) {
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

	result := make([]*models.ExperienceJoinForegnKey, 0)
	for rows.Next() {
		t := new(models.ExperienceJoinForegnKey)
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
			&t.GuideReview,
			&t.ActivitiesReview,
			&t.ServiceReview,
			&t.CleanlinessReview,
			&t.ValueReview,
			&t.ExpPaymentDeadlineAmount,
			&t.ExpPaymentDeadlineType,
			&t.IsCustomisedByUser,
			&t.MinimumBookingAmount,
			&t.MinimumBookingDesc,
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
			&t.GuideReview,
			&t.ActivitiesReview,
			&t.ServiceReview,
			&t.CleanlinessReview,
			&t.ValueReview,
			&t.ExpPaymentDeadlineAmount,
			&t.ExpPaymentDeadlineType,
			&t.IsCustomisedByUser,
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
func (m *experienceRepository) GetUserDiscoverPreference(ctx context.Context, page *int, size *int) ([]*models.ExpUserDiscoverPreference, error) {

	if page != nil && size != nil {
		query := `select c.city_id, c.city_name, city.city_desc,city_photos,a.* from cgo_indonesia.experiences a
			join cgo_indonesia.harbors b on b.id = a.harbors_id
			join 
			(
				select c.id as city_id, c.city_name as city_name from cgo_indonesia.user_preference_exps upe
				join cgo_indonesia.harbors h on h.id = upe.harbors_id
				join cgo_indonesia.cities c on c.id = h.city_id
				where upe.is_active = 1 and upe.is_deleted = 0
				order by upe.amount desc LIMIT ?,? 
			) c on c.city_id = b.city_id
            join cgo_indonesia.cities city on city.id = c.city_id;`

		res, err := m.fetchUserDiscoverPreference(ctx, query, page, size)
		if err != nil {
			return nil, err
		}
		return res, err
	} else {
		query := `select c.city_id, c.city_name, city.city_desc,city.city_photos,a.* from cgo_indonesia.experiences a
			join cgo_indonesia.harbors b on b.id = a.harbors_id
			join 
			(
				select c.id as city_id, c.city_name as city_name from cgo_indonesia.user_preference_exps upe
				join cgo_indonesia.harbors h on h.id = upe.harbors_id
				join cgo_indonesia.cities c on c.id = h.city_id
				where upe.is_active = 1 and upe.is_deleted = 0
				order by upe.amount desc  LIMIT 0,3
			) c on c.city_id = b.city_id
            join cgo_indonesia.cities city on city.id = c.city_id;`

		res, err := m.fetchUserDiscoverPreference(ctx, query)
		if err != nil {
			return nil, err
		}
		return res, err
	}

}

func (m *experienceRepository) GetUserDiscoverPreferenceByHarborsIdOrProvince(ctx context.Context, harborsId *string, provinceId *int) ([]*models.ExpUserDiscoverPreference, error) {

	if harborsId != nil {
		query := `select 
					province.id as province_id ,
					province.province_name,
					c.city_id, 
					c.city_name, 
					city.city_desc,
					city_photos,	
					b.id as id_harbors,
					b.harbors_name ,
					a.*
			from 
				experiences a
			join harbors b on b.id = a.harbors_id
			join 
			(
				select c.id as city_id, c.city_name as city_name,p.id as province_id from temp_user_preferences upe
				join harbors h on h.id = upe.harbors_id
				join cities c on c.id = h.city_id
				join provinces p on p.id = c.province_id
				where h.id = ?
			) c on c.city_id = b.city_id
            join cities city on city.id = c.city_id
			join provinces province on province.id = c.province_id
			WHERE a.harbors_id = ? AND a.is_deleted = 0 AND a.is_active = 1 AND a.status = 2`

		res, err := m.fetchUserDiscoverPreference(ctx, query, harborsId, harborsId)
		if err != nil {
			return nil, err
		}
		return res, err
	} else if provinceId != nil {
		query := `select 
					province.id as province_id ,
					province.province_name,
					city.city_desc,
					a.*
			from 
				experiences a
			join harbors b on b.id = a.harbors_id
            join cities city on city.id = b.city_id
			join provinces province on province.id = city.province_id
			WHERE province.id = ? AND a.is_deleted = 0 AND a.is_active = 1 AND a.status = 2`

		res, err := m.fetchUserDiscoverPreferenceProvince(ctx, query, provinceId)
		if err != nil {
			return nil, err
		}
		return res, err
	}
	return nil, nil
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
func (m *experienceRepository) GetByID(ctx context.Context, id string) (res *models.ExperienceJoinForegnKey, err error) {
	query := `SELECT e.*,m.minimum_booking_amount,m.minimum_booking_desc FROM cgo_indonesia.experiences e 
				join cgo_indonesia.minimum_bookings m on m.id = e.minimum_booking_id
 				WHERE  e.is_deleted = 0 AND e.is_active = 1 AND e.id=?`

	list, err := m.fetchJoinForegnKey(ctx, query, id)
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
func (m *experienceRepository) SelectIdGetByMerchantId(ctx context.Context, merchantId string) (res []*string, err error) {
	query := `SELECT id FROM experiences WHERE merchant_id = ?`
	rows, err := m.Conn.QueryContext(ctx, query, merchantId)
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

func (m *experienceRepository) Insert(ctx context.Context, a *models.Experience) (*string, error) {
	a.Id = guuid.New().String()
	query := `INSERT experiences SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , exp_title=?,exp_type=?,exp_trip_type=?,exp_booking_type=?,
				exp_desc=?,exp_max_guest=?,exp_pickup_place=?,exp_pickup_time=?,exp_pickup_place_longitude=?,
				exp_pickup_place_latitude=?,exp_pickup_place_maps_name=?,exp_itinerary=?,exp_facilities=?,exp_inclusion=?,
				exp_rules=?,status=?,rating=?,exp_location_latitude=?,exp_location_longitude=?,exp_location_name=?,
				exp_cover_photo=?,exp_duration=?,minimum_booking_id=?,merchant_id=?,harbors_id=?,exp_payment_deadline_amount=?,
				exp_payment_deadline_type=?,is_customised_by_user=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.ExpTitle, a.ExpType, a.ExpTripType,
		a.ExpBookingType, a.ExpDesc, a.ExpMaxGuest, a.ExpPickupPlace, a.ExpPickupTime, a.ExpPickupPlaceLongitude,
		a.ExpPickupPlaceLatitude, a.ExpPickupPlaceMapsName, a.ExpInternary, a.ExpFacilities, a.ExpInclusion,
		a.ExpRules, a.Status, a.Rating, a.ExpLocationLatitude, a.ExpLocationLongitude, a.ExpLocationName,
		a.ExpCoverPhoto, a.ExpDuration, a.MinimumBookingId, a.MerchantId, a.HarborsId, a.ExpPaymentDeadlineAmount,
		a.ExpPaymentDeadlineType, a.IsCustomisedByUser)
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

func (m *experienceRepository) UpdateStatus(ctx context.Context, status int, id string, user string) error {
	query := `UPDATE experiences SET modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , status=?
				WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, user, time.Now(), nil, nil, 0, 1, status, id)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return nil
}
func (m *experienceRepository) Update(ctx context.Context, a *models.Experience) (*string, error) {
	query := `UPDATE experiences SET modified_by=?, modified_date=? , deleted_by=? , 
				deleted_date=? , is_deleted=? , is_active=? , exp_title=?,exp_type=?,exp_trip_type=?,exp_booking_type=?,
				exp_desc=?,exp_max_guest=?,exp_pickup_place=?,exp_pickup_time=?,exp_pickup_place_longitude=?,
				exp_pickup_place_latitude=?,exp_pickup_place_maps_name=?,exp_itinerary=?,exp_facilities=?,exp_inclusion=?,
				exp_rules=?,status=?,exp_location_latitude=?,exp_location_longitude=?,exp_location_name=?,
				exp_cover_photo=?,exp_duration=?,minimum_booking_id=?,merchant_id=?,harbors_id=? ,exp_payment_deadline_amount=?,
				exp_payment_deadline_type=?,is_customised_by_user=? WHERE id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.ExpTitle, a.ExpType, a.ExpTripType,
		a.ExpBookingType, a.ExpDesc, a.ExpMaxGuest, a.ExpPickupPlace, a.ExpPickupTime, a.ExpPickupPlaceLongitude,
		a.ExpPickupPlaceLatitude, a.ExpPickupPlaceMapsName, a.ExpInternary, a.ExpFacilities, a.ExpInclusion,
		a.ExpRules, a.Status, a.ExpLocationLatitude, a.ExpLocationLongitude, a.ExpLocationName,
		a.ExpCoverPhoto, a.ExpDuration, a.MinimumBookingId, a.MerchantId, a.HarborsId, a.ExpPaymentDeadlineAmount,
		a.ExpPaymentDeadlineType, a.IsCustomisedByUser, a.Id)
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
