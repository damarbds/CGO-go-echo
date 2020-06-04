package repository

import (
	"context"
	"database/sql"
	"github.com/service/currency"
	"time"

	"github.com/models"
	"github.com/sirupsen/logrus"
)

type currencyRepository struct {
	Conn *sql.DB
}

func NewCurrencyRepository(Conn *sql.DB) currency.Repository {
	return &currencyRepository{Conn: Conn}
}

func (f currencyRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Currency, error) {
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

	result := make([]*models.Currency, 0)
	for rows.Next() {
		t := new(models.Currency)
		err = rows.Scan(
			&t.Id,
			&t.Code 	,
			&t.Name 	,
			&t.Symbol,
			&t.CreatedBy,
			&t.CreatedDate ,
			&t.ModifiedBy  ,
			&t.ModifiedDate ,
			&t.DeletedBy   ,
			&t.DeletedDate  ,
			&t.IsDeleted  ,
			&t.IsActive    ,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}


func (m *currencyRepository) Fetch(ctx context.Context, limit,offset int) ([]*models.Currency, error) {
	if limit != 0 {
		query := `Select * FROM currencies where is_deleted = 0 AND is_active = 1 `

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
		query := `Select * FROM currencies where is_deleted = 0 AND is_active = 1 `

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
func (m *currencyRepository) GetById(ctx context.Context, id int) (res *models.Currency, err error) {
	query := `SELECT * FROM currencies WHERE id = ?`

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

func (m *currencyRepository) GetCount(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM currencies WHERE is_deleted = 0 and is_active = 1`

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

func (m *currencyRepository) Insert(ctx context.Context, a *models.Currency) (*int, error) {
	query := `INSERT currencies SET created_by=? , created_date=? , modified_by=?, modified_date=? ,
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , code=?,name=? , 
				symbol=? `
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	res, err := stmt.ExecContext(ctx,a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.Code,a.Name,
		a.Symbol)
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

func (m *currencyRepository) Update(ctx context.Context, a *models.Currency) error {
	query := `UPDATE currencies set modified_by=?, modified_date=? ,code=?,name=? , 
				symbol=?   WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.Code,a.Name, a.Symbol,a.Id)
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

func (m *currencyRepository) Delete(ctx context.Context, id int, deletedBy string) error {
	query := `UPDATE currencies SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
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