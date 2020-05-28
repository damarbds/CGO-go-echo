package repository

import (
	"database/sql"
	"time"

	"github.com/service/exp_availability"

	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type exp_availabilityRepository struct {
	Conn *sql.DB
}

// NewExpexp_availabilityRepository will create an object that represent the exp_exp_availability.repository interface
func NewExpavailabilityRepository(Conn *sql.DB) exp_availability.Repository {
	return &exp_availabilityRepository{Conn}
}
func (m *exp_availabilityRepository) GetByExpIds(ctx context.Context, expId []*string) ([]*models.ExpAvailability, error) {
	res := make([]*models.ExpAvailability,0)
	query := `SELECT * FROM exp_availabilities WHERE is_deleted = 0 AND is_active = 1 AND `
	if len(expId) != 0 {
		for index, id := range expId {
			if index == 0 && index != (len(expId)-1) {
				query = query + ` exp_id = '` + *id + `' `
			} else if index == 0 && index == (len(expId)-1) {
				query = query + ` exp_id = '` + *id + `' `
			} else if index == (len(expId) - 1) {
				query = query + ` OR  exp_id = '` + *id + `' `
			} else {
				query = query + ` OR  exp_id = '` + *id + `' `
			}
		}
		resp, err := m.fetch(ctx, query)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
	res = resp
	}
	return res, nil
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
func (m *exp_availabilityRepository) GetCountDate(ctx context.Context, date string, expId []*string) (int, error) {
	query := `
	SELECT
		count(DISTINCT b.booking_date,b.exp_id) AS count
	FROM
		transactions t
	JOIN booking_exps b on b.id = t.booking_exp_id
	WHERE
		DATE (b.booking_date) = ? AND 
		t.status = 2`

	for index, id := range expId {
		if index == 0 && index != (len(expId)-1) {
			query = query + ` AND (b.exp_id LIKE '%` + *id + `%' `
		} else if index == 0 && index == (len(expId)-1) {
			query = query + ` AND (b.exp_id LIKE '%` + *id + `%' ) `
		} else if index == (len(expId) - 1) {
			query = query + ` OR  b.exp_id LIKE '%` + *id + `%' ) `
		} else {
			query = query + ` OR  b.exp_id LIKE '%` + *id + `%' `
		}
	}
	rows, err := m.Conn.QueryContext(ctx, query, date)
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
func (m *exp_availabilityRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExpAvailability, error) {
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

	result := make([]*models.ExpAvailability, 0)
	for rows.Next() {
		t := new(models.ExpAvailability)
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
			&t.ExpAvailabilityMonth,
			&t.ExpAvailabilityDate,
			&t.ExpAvailabilityYear,
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
func (e exp_availabilityRepository) GetByExpId(ctx context.Context, expId string) ([]*models.ExpAvailability, error) {
	query := `SELECT * FROM exp_availabilities WHERE is_deleted = 0 AND is_active = 1 AND exp_id = ?`

	list, err := e.fetch(ctx, query, expId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return list, nil
}
func (m *exp_availabilityRepository) Insert(ctx context.Context, a models.ExpAvailability) (string, error) {
	id := guuid.New()
	a.Id = id.String()
	query := `INSERT exp_availabilities SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , 
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , exp_availability_month=?,
				exp_availability_date=?,exp_availability_year=?,exp_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return "", err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.ExpAvailabilityMonth,
		a.ExpAvailabilityDate, a.ExpAvailabilityYear, a.ExpId)
	if err != nil {
		return "", err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return a.Id, nil
}
func (m *exp_availabilityRepository) Update(ctx context.Context, a models.ExpAvailability) error {
	query := `UPDATE exp_availabilities SET modified_by=?, modified_date=? , 
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , exp_availability_month=?,
				exp_availability_date=?,exp_availability_year=?,exp_id=? 
				WHERE id =?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), nil, nil, 0, 1, a.ExpAvailabilityMonth,
		a.ExpAvailabilityDate, a.ExpAvailabilityYear, a.ExpId, a.Id)
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
func (m *exp_availabilityRepository) Deletes(ctx context.Context, ids []string, expId string, deletedBy string) error {
	query := `UPDATE  exp_availabilities SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE exp_id=?`
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

func (m *exp_availabilityRepository) DeleteByExpId(ctx context.Context, expId string, deletedBy string) error {
	query := `UPDATE exp_availabilities SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE exp_id=?`

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
