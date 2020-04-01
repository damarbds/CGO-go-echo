package repository

import (
	"context"
	"database/sql"
	"encoding/base64"

	"time"

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


// NewpromoRepository will create an object that represent the article.Repository interface
func NewpromoRepository(Conn *sql.DB) promo.Repository {
	return &promoRepository{Conn}
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
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}


func (m *promoRepository) Fetch(ctx context.Context, page *int, size *int) ([]*models.Promo, error) {
	if page != nil && size != nil{
		query := `Select * FROM promos where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc LIMIT ? OFFSET ? `

		res, err := m.fetch(ctx, query, size,page)
		if err != nil {
			return nil, err
		}
		return res, err

	}else {
		query := `Select * FROM promos where is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

		res, err := m.fetch(ctx, query)
		if err != nil {
			return nil, err
		}
		return res, err
	}


	return nil, nil
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
