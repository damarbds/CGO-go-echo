package repository

import (
	"context"
	"database/sql"
	"time"

	guuid "github.com/google/uuid"

	"github.com/models"
	payment "github.com/service/exp_payment"
	"github.com/sirupsen/logrus"
)

type expPaymentRepository struct {
	Conn *sql.DB
}


// NewExpPaymentRepository will create an object that represent the exp_payment.repository interface
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

func (m *expPaymentRepository) GetById(ctx context.Context, id string) ([]*models.ExperiencePaymentJoinType, error) {
	query := `SELECT ep.*,ept.exp_payment_type_name as payment_type_name ,ept.exp_payment_type_desc as payment_type_desc 
			FROM experience_payments ep
			JOIN experience_payment_types ept on ept.id = ep.exp_payment_type_id
			WHERE ep.id = ? `

	list, err := m.fetch(ctx, query, id)
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
func (m *expPaymentRepository) Insert(ctx context.Context, a models.ExperiencePayment) (string, error) {
	id := guuid.New()
	a.Id = id.String()
	query := `INSERT experience_payments SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? , 
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , exp_payment_type_id=?,exp_id=?,
				price_item_type=?,currency=?,price=?,custom_price=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return "", err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.ExpPaymentTypeId, a.ExpId,
		a.PriceItemType, a.Currency, a.Price, a.CustomPrice)
	if err != nil {
		return "", err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return a.Id, nil
}
func (m *expPaymentRepository) Update(ctx context.Context, a models.ExperiencePayment) error {
	query := `UPDATE experience_payments SET modified_by=?, modified_date=? , 
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , exp_payment_type_id=?,exp_id=?,
				price_item_type=?,currency=?,price=?,custom_price=? 
				WHERE id =?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, a.ModifiedBy, a.ModifiedDate, nil, nil, 0, 1, a.ExpPaymentTypeId, a.ExpId,
		a.PriceItemType, a.Currency, a.Price, a.CustomPrice, a.Id)
	if err != nil {
		return err
	}

	//lastID, err := res.RowsAffected()
	//if err != nil {
	//	return err
	//}

	//a.Id = lastID
	return nil
}
func (m *expPaymentRepository) Deletes(ctx context.Context, ids []string, expId string, deletedBy string) error {
	query := `UPDATE  experience_payments SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE exp_id=?`
	for index, id := range ids {
		if index == 0 && index != (len(ids)-1) {
			query = query + ` AND (id !=` + id
		} else if index == 0 && index == (len(ids)-1) {
			query = query + ` AND (id !=` + id + ` ) `
		} else if index == (len(ids) - 1) {
			query = query + ` OR id !=` + id + ` ) `
		} else {
			query = query + ` OR id !=` + id
		}
	}
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0, expId)
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


func (m *expPaymentRepository) DeleteByExpId(ctx context.Context, expId string, deletedBy string) error {
	query := `UPDATE experience_payments SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE exp_id=?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0, expId)
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
