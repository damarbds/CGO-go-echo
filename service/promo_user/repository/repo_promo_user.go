package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/service/promo_user"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type promoUserRepository struct {
	Conn *sql.DB
}

func (m *promoUserRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.PromoUser, error) {
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

	result := make([]*models.PromoUser, 0)
	for rows.Next() {
		t := new(models.PromoUser)
		err = rows.Scan(
			&t.Id,
			&t.PromoId,
			&t.UserId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
func (m promoUserRepository) GetByUserId(ctx context.Context, userId string, promoId string) (res []*models.PromoUser, err error) {
	query := `SELECT * FROM promo_users WHERE promo_id= ?`
	if userId != "" {
		query = query + ` and user_id = '` + userId + `' `
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
func NewpromoUserRepository(Conn *sql.DB) promo_user.Repository {
	return &promoUserRepository{Conn}
	}

func (m promoUserRepository) Insert(ctx context.Context, a models.PromoUser) error {
	query := `INSERT promo_users SET promo_id=?,user_id=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.PromoId, a.UserId)
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

func (m promoUserRepository) DeleteByUserId(ctx context.Context, userId string, promoId string) error {
	query := "DELETE FROM promo_users WHERE user_id = ? AND promo_id=?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, userId, promoId)
	if err != nil {

		return err
	}

	return nil
}
