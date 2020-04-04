package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	payment "github.com/service/exp_payment"
	"github.com/sirupsen/logrus"
)

type expPaymentRepository struct {
	Conn *sql.DB
}

// NewExpPaymentRepository will create an object that represent the exp_payment.Repository interface
func NewExpPaymentRepository(Conn *sql.DB) payment.Repository {
	return &expPaymentRepository{Conn}
}

func (m *expPaymentRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExperiencePayment, error) {
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

	result := make([]*models.ExperiencePayment, 0)
	for rows.Next() {
		t := new(models.ExperiencePayment)
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
			&t.ExpPaymentTypeId,
			&t.ExpId,
			&t.PriceItemType,
			&t.Currency,
			&t.Price,
			&t.CustomPrice,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (e expPaymentRepository) GetByExpID(ctx context.Context, expID string) (*models.ExperiencePayment, error) {
	query := `SELECT * FROM experience_payments WHERE exp_id = ?`

	list, err := e.fetch(ctx, query, expID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}else if len(list) == 0 {
		return nil,nil
	}

	return list[0], nil
}
