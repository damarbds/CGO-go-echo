package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/promo_experience_transport"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type promoExperienceTransportRepository struct {
	Conn *sql.DB
}

func NewpromoMerchantRepository(Conn *sql.DB) promo_experience_transport.Repository{
	return &promoExperienceTransportRepository{Conn}
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
func (m *promoExperienceTransportRepository) CountByPromoId(ctx context.Context, promoId string) (int, error) {
	query := `SELECT COUNT(*) as count FROM promo_experience_transports WHERE promo_id= ?`

	rows, err := m.Conn.QueryContext(ctx, query, promoId)
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

func (m *promoExperienceTransportRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.PromoExperienceTransport, error) {
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

	result := make([]*models.PromoExperienceTransport, 0)
	for rows.Next() {
		t := new(models.PromoExperienceTransport)
		err = rows.Scan(
			&t.Id,
			&t.PromoId,
			&t.ExperienceId,
			&t.TransportationId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m promoExperienceTransportRepository) GetByExperienceTransportId(ctx context.Context, experienceId string, transportId string,promoId string) (res []*models.PromoExperienceTransport, err error) {
	query := `SELECT * FROM promo_experience_transports WHERE promo_id= ?`
	if experienceId != "" {
		query = query + ` and experience_id = '` + experienceId + `' `
	} else if transportId != "" {
		query = query + ` and transportation_id = '` + transportId + `' `
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

// NewpromoRepository will create an object that represent the article.repository interface
//func NewpromoMerchantRepository(Conn *sql.DB) promo_merchant.Repository {
//	return &promoMerchantRepository{Conn}
//}

func (m promoExperienceTransportRepository) Insert(ctx context.Context, a models.PromoExperienceTransport) error {
	query := `INSERT promo_experience_transports SET promo_id=?,experience_id=?,transportation_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.PromoId, a.ExperienceId, a.TransportationId)
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

func (m promoExperienceTransportRepository) DeleteById(ctx context.Context, merchantId string, promoId string) error {
	query := "DELETE FROM promo_merchants WHERE merchant_id = ? AND promo_id=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, merchantId, promoId)
	if err != nil {

		return err
	}

	return nil
}
