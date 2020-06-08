package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/exclude"
	"github.com/sirupsen/logrus"
	"time"
)

type excludeRepository struct {
	Conn *sql.DB
}

func NewExcludeRepository(Conn *sql.DB) exclude.Repository {
	return &excludeRepository{Conn: Conn}
}

func (f excludeRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Exclude, error) {
	rows, err := f.Conn.QueryContext(ctx, query, args...)
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

	result := make([]*models.Exclude, 0)
	for rows.Next() {
		t := new(models.Exclude)
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
			&t.ExcludeName,
			&t.ExcludeIcon,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (f excludeRepository) List(ctx context.Context) ([]*models.Exclude, error) {
	query := `SELECT * FROM excludes WHERE is_deleted = 0 and is_active = 1`

	res, err := f.fetch(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

func (m *excludeRepository) Fetch(ctx context.Context, limit,offset int) ([]*models.Exclude, error) {
	if limit != 0 {
		query := `Select * FROM excludes where is_deleted = 0 AND is_active = 1 `

		//if search != ""{
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
		query := `Select * FROM includes where is_deleted = 0 AND is_active = 1 `

		//if search != ""{
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
func (m *excludeRepository) GetById(ctx context.Context, id int) (res *models.Exclude, err error) {
	query := `SELECT * FROM excludes WHERE id = ?`

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

func (m *excludeRepository) GetCount(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM excludes WHERE is_deleted = 0 and is_active = 1`

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

func (m *excludeRepository) Insert(ctx context.Context, a *models.Exclude) (*int, error) {
	query := `INSERT excludes SET created_by=? , created_date=? , modified_by=?, modified_date=? ,
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , exclude_name=?, 
				exclude_icon=? `
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.ExcludeName,
		a.ExcludeIcon)
	if err != nil {
		return nil,err
	}


	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	a.Id = int(lastID)
	return &a.Id,nil
}

func (m *excludeRepository) Update(ctx context.Context, a *models.Exclude) error {
	query := `UPDATE excludes set modified_by=?, modified_date=? ,exclude_name=?,exclude_icon=?  WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.ExcludeName, a.ExcludeIcon, a.Id)
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

func (m *excludeRepository) Delete(ctx context.Context, id int, deletedBy string) error {
	query := `UPDATE excludes SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0,id)
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


func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}