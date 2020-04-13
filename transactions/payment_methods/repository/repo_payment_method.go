package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/sirupsen/logrus"
	pm "github.com/transactions/payment_methods"
)

type paymentMethodRepository struct {
	Conn *sql.DB
}

// NewPaymentMethodRepository will create an object that represent the article.Repository interface
func NewPaymentMethodRepository(Conn *sql.DB) pm.Repository {
	return &paymentMethodRepository{Conn}
}

func (m *paymentMethodRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.PaymentMethod, error) {
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

	result := make([]*models.PaymentMethod, 0)
	for rows.Next() {
		t := new(models.PaymentMethod)
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
			&t.Name,
			&t.Type,
			&t.Desc,
			&t.Icon,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (p paymentMethodRepository) Fetch(ctx context.Context) ([]*models.PaymentMethod, error) {
	query := `SELECT * FROM payment_methods WHERE is_deleted = 0 AND is_active = 1 ORDER BY created_date desc`

	res, err := p.fetch(ctx, query)
	if err != nil {
		return nil, err
	}

	return res, nil
}
