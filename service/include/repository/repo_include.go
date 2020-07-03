package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/include"
	"github.com/sirupsen/logrus"
	"time"
)

type includeRepository struct {
	Conn *sql.DB
}


func NewIncludeRepository(Conn *sql.DB) include.Repository {
	return &includeRepository{Conn: Conn}
}

func (m includeRepository) GetByName(ctx context.Context, name string) (res *models.Include,err error) {
	query := `SELECT * FROM includes WHERE include_name = ?`

	list, err := m.fetch(ctx, query, name)
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
func (f includeRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Include, error) {
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

	result := make([]*models.Include, 0)
	for rows.Next() {
		t := new(models.Include)
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
			&t.IncludeName,
			&t.IncludeIcon,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (f includeRepository) List(ctx context.Context) ([]*models.Include, error) {
	query := `SELECT * FROM includes WHERE is_deleted = 0 and is_active = 1`

	res, err := f.fetch(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

func (m *includeRepository) Fetch(ctx context.Context, limit,offset int) ([]*models.Include, error) {

	query := `SELECT * FROM includes where is_deleted = 0 AND is_active = 1 `
	if limit != 0 {
		query = query + `ORDER BY created_date desc LIMIT ? OFFSET ?`
		res, err := m.fetch(ctx, query, limit, offset)
		if err != nil {
			return nil, err
		}
		return res, err

	} else {
		query = query + `ORDER BY created_date desc`
		res, err := m.fetch(ctx, query)
		if err != nil {
			return nil, err
		}
		return res, err
	}
}
func (m *includeRepository) GetById(ctx context.Context, id int) (res *models.Include, err error) {
	query := `SELECT * FROM includes WHERE id = ?`

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

func (m *includeRepository) GetCount(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM includes WHERE is_deleted = 0 and is_active = 1`

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

func (m *includeRepository) Insert(ctx context.Context, a *models.Include) (*int, error) {
	query := `INSERT includes SET created_by=? , created_date=? , modified_by=?, modified_date=? ,
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , include_name=?, 
				include_icon=? `
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,a.CreatedBy, a.CreatedDate, nil, nil, nil, nil, 0, 1, a.IncludeName,
		a.IncludeIcon)
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

func (m *includeRepository) Update(ctx context.Context, a *models.Include) error {
	query := `UPDATE includes set modified_by=?, modified_date=? ,include_name=?,include_icon=?  WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.IncludeName, a.IncludeIcon, a.Id)
	if err != nil {
		return err
	}

	return nil
}

func (m *includeRepository) Delete(ctx context.Context, id int, deletedBy string) error {
	query := `UPDATE includes SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0,id)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

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