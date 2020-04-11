package repository

import (
	"context"
	"database/sql"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/transactions/payment"
)

type paymentRepository struct {
	Conn *sql.DB
}

// NewPaymentRepository will create an object that represent the article.Repository interface
func NewPaymentRepository(Conn *sql.DB) payment.Repository {
	return &paymentRepository{Conn}
}

func (p paymentRepository) Insert(ctx context.Context, pay *models.Transaction) (*models.Transaction, error) {
	id := guuid.New()
	pay.Id = id.String()
	q := `INSERT transactions SET id = ?, created_by = ?, created_date = ?, modified_by = ?, 
	modified_date = ?, deleted_by = ?, deleted_date = ?, is_deleted = ?, is_active = ?, booking_type = ?, 
	booking_exp_id = ?, promo_id = ?, payment_method_id = ?, experience_payment_id = ?, status = ?, total_price = ?, currency = ?`

	stmt, err := p.Conn.PrepareContext(ctx, q)
	if err != nil {
		return nil, err
	}

	_, err = stmt.ExecContext(
		ctx,
		pay.Id,
		pay.CreatedBy,
		pay.CreatedDate,
		pay.ModifiedBy,
		pay.ModifiedDate,
		pay.DeletedBy,
		pay.DeletedDate,
		pay.IsDeleted,
		pay.IsActive,
		pay.BookingType,
		pay.BookingExpId,
		pay.PromoId,
		pay.PaymentMethodId,
		pay.ExperiencePaymentId,
		pay.Status,
		pay.TotalPrice,
		pay.Currency,
	)
	if err != nil {
		return nil, err
	}

	return pay, nil
}

func (p paymentRepository) ConfirmPayment(ctx context.Context, confirmIn *models.ConfirmPaymentIn) error {
	query := `
	UPDATE
		transactions,
		booking_exps
	SET
		transactions.status = ?,
		booking_exps.status = ?
	WHERE
		booking_exp_id = booking_exps.id
		AND transactions.id = ?`

	stmt, err := p.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx,
		confirmIn.TransactionStatus,
		confirmIn.BookingStatus,
		confirmIn.TransactionID,
	)
	if err != nil {
		return err
	}

	return nil
}
