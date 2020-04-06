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

func (m *expPaymentRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.ExperiencePaymentJoinType, error) {
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

	result := make([]*models.ExperiencePaymentJoinType, 0)
	for rows.Next() {
		t := new(models.ExperiencePaymentJoinType)
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
			&t.ExpPaymentTypeName,
			&t.ExpPaymentTypeDesc,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (e expPaymentRepository) GetByExpID(ctx context.Context, expID string) ([]*models.ExperiencePaymentJoinType, error) {
	query := `SELECT ep.*,ept.exp_payment_type_name as payment_type_name ,ept.exp_payment_type_desc as payment_type_desc FROM experience_payments ep
			JOIN experience_payment_types ept on ept.id = ep.exp_payment_type_id
			WHERE ep.exp_id = ? AND ep.is_deleted = 0 AND ep.is_active = 1`

	list, err := e.fetch(ctx, query, expID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	} else if len(list) == 0 {
		return nil, nil
	}

	return list, nil
}
