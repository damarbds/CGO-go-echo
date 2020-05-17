package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"time"

	guuid "github.com/google/uuid"

	"github.com/sirupsen/logrus"

	"github.com/models"
	"github.com/service/promo"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type promoRepository struct {
	Conn *sql.DB
}

// NewpromoRepository will create an object that represent the article.repository interface
func NewpromoRepository(Conn *sql.DB) promo.Repository {
	return &promoRepository{Conn}
}

func (m *promoRepository) GetCount(ctx context.Context) (int, error) {
	query := `SELECT count(*) AS count FROM promos WHERE is_deleted = 0 and is_active = 1`

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
func (m *promoRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Promo, error) {
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

	result := make([]*models.Promo, 0)
	for rows.Next() {
		t := new(models.Promo)
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
			&t.PromoCode,
			&t.PromoName,
			&t.PromoDesc,
			&t.PromoValue,
			&t.PromoType,
			&t.PromoImage,
			&t.StartDate,
			&t.EndDate,
			&t.MaxUsage,
			&t.ProductionCapacity,
			&t.CurrencyId,
			&t.PromoProductType,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *promoRepository) Delete(ctx context.Context, id string, deletedBy string) error {
	query := `UPDATE promos SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE id =?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0, id)
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

func (m *promoRepository) Insert(ctx context.Context, a *models.Promo) (string, error) {
	a.Id = guuid.New().String()
	query := `INSERT promos SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? ,
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , promo_code=?,promo_name=? , 
				promo_desc=? ,promo_value=?,promo_type=?,promo_image=?,start_date=?,end_date=?,currency_id	=?,
				max_usage=?,production_capacity=?,promo_product_type=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return "", err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.PromoCode, a.PromoName,
		a.PromoDesc, a.PromoValue, a.PromoType, a.PromoImage, a.StartDate, a.EndDate, a.CurrencyId, a.MaxUsage, a.ProductionCapacity,a.PromoProductType)
	if err != nil {
		return "", err
	}

	//lastID, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	//a.Id = lastID
	return a.Id, nil
}

func (m *promoRepository) Update(ctx context.Context, a *models.Promo) error {
	query := `UPDATE promos set modified_by=?, modified_date=? , promo_code=?,promo_name=? , promo_desc=? ,promo_value=?,
				promo_type=?,promo_image=?,start_date=?,end_date=?,currency_id=?,max_usage=?,production_capacity=?, 
				promo_product_type=? WHERE id = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return nil
	}

	_, err = stmt.ExecContext(ctx, a.ModifiedBy, time.Now(), a.PromoCode, a.PromoName, a.PromoDesc, a.PromoValue,
		a.PromoType, a.PromoImage, a.StartDate, a.EndDate, a.CurrencyId, a.MaxUsage, a.ProductionCapacity, a.PromoProductType,a.Id)
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

func (m *promoRepository) GetById(ctx context.Context, id string) (res *models.Promo, err error) {
	query := `SELECT * FROM promos WHERE is_deleted = 0 and is_active = 1 and id = ?`

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

func (m *promoRepository) GetByCode(ctx context.Context, code string,promoType *int,merchantId string) ([]*models.Promo, error) {
	var query string
	if merchantId != ""{
		query = `SELECT p.* 
				FROM 
					promos p
				JOIN promo_merchants pm on pm.promo_id = p.id
				WHERE 
					p.promo_code = ? AND 
					p.promo_product_type = ? AND 
 					p.is_deleted = 0 AND 
					p.is_active = 1 AND
					pm.merchant_id = '` + merchantId  + `'`
	}else {
		query = `SELECT * FROM promos WHERE promo_code = ? AND promo_product_type = ? AND is_deleted = 0 AND is_active = 1`
	}

	res, err := m.fetch(ctx, query, code,promoType)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (m *promoRepository) Fetch(ctx context.Context, page *int, size *int, search string) ([]*models.Promo, error) {
	if page != nil && size != nil {
		query := `Select * FROM promos where is_deleted = 0 AND is_active = 1 `

		if search != "" {
			query = query + `AND (promo_name LIKE '%` + search + `%'` +
				`OR promo_desc LIKE '%` + search + `%' ` +
				`OR start_date LIKE '%` + search + `%' ` +
				`OR end_date LIKE '%` + search + `%' ` +
				`OR promo_code LIKE '%` + search + `%' ` +
				`OR max_usage LIKE '%` + search + `%' ` + `) `
		}
		query = query + ` ORDER BY created_date desc LIMIT ? OFFSET ? `
		res, err := m.fetch(ctx, query, size, page)
		if err != nil {
			return nil, err
		}
		return res, err

	} else {
		query := `Select * FROM promos where is_deleted = 0 AND is_active = 1 `

		if search != "" {
			query = query + `AND (promo_name LIKE '%` + search + `%'` +
				`OR promo_desc LIKE '%` + search + `%' ` +
				`OR start_date LIKE '%` + search + `%' ` +
				`OR end_date LIKE '%` + search + `%' ` +
				`OR promo_code LIKE '%` + search + `%' ` +
				`OR max_usage LIKE '%` + search + `%' ` + `) `
		}
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
