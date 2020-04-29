package repository

import (
	"context"
	"database/sql"
	"github.com/service/promo_merchant"
	"github.com/sirupsen/logrus"

	"github.com/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type promoMerchantRepository struct {
	Conn *sql.DB
}
func (m *promoMerchantRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.PromoMerchant, error) {
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

	result := make([]*models.PromoMerchant, 0)
	for rows.Next() {
		t := new(models.PromoMerchant)
		err = rows.Scan(
			&t.Id,
			&t.PromoId,
			&t.MerchantId ,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m promoMerchantRepository) GetByMerchantId(ctx context.Context, merchantId string, promoId string) (res []*models.PromoMerchant,err error) {
	query := `SELECT * FROM promo_merchants WHERE promo_id= ?`
	if merchantId != ""{
		query = query + ` and merchant_id = '` + merchantId + `' `
	}
	list, err := m.fetch(ctx, query, promoId)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list
	}

	return
}

// NewpromoRepository will create an object that represent the article.Repository interface
func NewpromoMerchantRepository(Conn *sql.DB) promo_merchant.Repository {
	return &promoMerchantRepository{Conn}
}
func (m promoMerchantRepository) Insert(ctx context.Context, a models.PromoMerchant) error {
	query := `INSERT promo_merchants SET promo_id=?,merchant_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.PromoId,a.MerchantId)
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

func (m promoMerchantRepository) DeleteByMerchantId(ctx context.Context, merchantId string,promoId string) error {
	query := "DELETE FROM promo_merchants WHERE merchant_id = ? AND promo_id=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, merchantId,promoId)
	if err != nil {

		return err
	}


	return nil
}
