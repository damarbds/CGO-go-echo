package repository

import (
	"database/sql"
	"github.com/models"
	"github.com/service/exclusion_service"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type exclusionServiceRepository struct {
	Conn *sql.DB
}


// NewExpexp_availabilityRepository will create an object that represent the exp_exp_availability.repository interface
func NewExclusionServiceRepository(Conn *sql.DB) exclusion_service.Repository {
	return &exclusionServiceRepository{Conn}
}

func (m *exclusionServiceRepository) GetCount(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM exclusion_services WHERE is_deleted = 0 and is_active = 1`

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
func (m *exclusionServiceRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExclusionService, error) {
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

	result := make([]*models.ExclusionService, 0)
	for rows.Next() {
		t := new(models.ExclusionService)
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
			&t.ExclusionServiceName,
			&t.ExclusionServiceType,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m exclusionServiceRepository) Fetch(ctx context.Context, limit, offset int) ([]*models.ExclusionService, error) {
	if limit != 0 {
		query := `Select * FROM exclusion_services where is_deleted = 0 AND is_active = 1 `

		//if search != "" {
		//	query = query + `AND (promo_name LIKE '%` + search + `%'` +
		//		`OR promo_desc LIKE '%` + search + `%' ` +
		//		`OR start_date LIKE '%` + search + `%' ` +
		//		`OR end_date LIKE '%` + search + `%' ` +
		//		`OR promo_code LIKE '%` + search + `%' ` +
		//		`OR max_usage LIKE '%` + search + `%' ` + `) `
		//}
		query = query + ` ORDER BY created_date desc LIMIT ? OFFSET ? `
		res, err := m.fetch(ctx, query, limit, offset)
		if err != nil {
			return nil, err
		}
		return res, err

	} else {
		query := `Select * FROM exclusion_services where is_deleted = 0 AND is_active = 1 `

		//if search != "" {
		//	query = query + `AND (promo_name LIKE '%` + search + `%'` +
		//		`OR promo_desc LIKE '%` + search + `%' ` +
		//		`OR start_date LIKE '%` + search + `%' ` +
		//		`OR end_date LIKE '%` + search + `%' ` +
		//		`OR promo_code LIKE '%` + search + `%' ` +
		//		`OR max_usage LIKE '%` + search + `%' ` + `) `
		//}
		query = query + ` ORDER BY created_date desc `
		res, err := m.fetch(ctx, query)
		if err != nil {
			return nil, err
		}
		return res, err
	}
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


