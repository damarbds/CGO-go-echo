package repository

import (
	"context"
	"database/sql"
	"github.com/models"

	"github.com/sirupsen/logrus"
	"github.com/transactions/transaction"
)

type transactionRepository struct {
	Conn *sql.DB
}

func NewTransactionRepository(Conn *sql.DB) transaction.Repository {
	return &transactionRepository{Conn: Conn}
}

func (t transactionRepository) List(ctx context.Context, status string, limit, offset int) ([]*models.TransactionOut, error) {
	var transactionStatus int
	var bookingStatus int

	query := `
	SELECT
		t.id as transaction_id,
		exp_id,
		e.exp_type,
		booking_exp_id,
		order_id as booking_code,
		b.created_date as booking_date,
		b.booking_date as check_in_date,
		b.booked_by,
		guest_desc,
		b.booked_by_email as email
	FROM
		transactions t
		JOIN booking_exps b ON t.booking_exp_id = b.id
		JOIN experiences e ON b.exp_id = e.id`

	queryWithoutStatus := query + ` WHERE t.is_deleted = 0
		AND t.is_active = 1 ORDER BY t.created_date DESC LIMIT ? OFFSET ?`
	list, err := t.fetch(ctx, queryWithoutStatus, limit, offset)
	if status != "" {

		if status == "pending" {
			transactionStatus = 0
		} else if status == "waitingApproval" {
			transactionStatus = 1
		} else if status == "confirm" {
			transactionStatus = 2
		}
		queryWithStatus := query + ` WHERE t.is_deleted = 0
		AND t.is_active = 1 AND t.status = ? ORDER BY t.created_date DESC LIMIT ? OFFSET ?`
		list, err = t.fetch(ctx, queryWithStatus, transactionStatus, limit, offset)

		if status == "failed" {
			transactionStatus = 3
			cancelledStatus := 4
			queryWithStatus = query + ` WHERE t.is_deleted = 0
		AND t.is_active = 1 AND t.status IN (?,?) ORDER BY t.created_date DESC LIMIT ? OFFSET ?`
			list, err = t.fetch(ctx, queryWithStatus, transactionStatus, cancelledStatus, limit, offset)
		}

		if status == "boarded" {
			transactionStatus = 1
			bookingStatus = 3
			queryWithStatus = query + ` WHERE t.is_deleted = 0
		AND t.is_active = 1 AND t.status = ? AND b.status = ? ORDER BY t.created_date DESC LIMIT ? OFFSET ?≈`
			list, err = t.fetch(ctx, queryWithStatus, transactionStatus, bookingStatus, limit, offset)
		}
	}
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (t transactionRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.TransactionOut, error) {
	rows, err := t.Conn.QueryContext(ctx, query, args...)
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

	result := make([]*models.TransactionOut, 0)
	for rows.Next() {
		t := new(models.TransactionOut)
		err = rows.Scan(
			&t.TransactionId,
			&t.ExpId,
			&t.ExpType,
			&t.BookingExpId,
			&t.BookingCode,
			&t.BookingDate,
			&t.CheckInDate,
			&t.BookedBy,
			&t.GuestDesc,
			&t.Email,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (t transactionRepository) CountSuccess(ctx context.Context) (int, error) {
	query := `SELECT count(*) as count FROM transactions WHERE is_deleted = 0 AND is_active = 1 AND status = 2`

	rows, err := t.Conn.QueryContext(ctx, query)
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

func (t transactionRepository) Count(ctx context.Context, status string) (int, error) {
	query := `
	SELECT
		count(*) as count
	FROM transactions t
		JOIN booking_exps b ON t.booking_exp_id = b.id
	WHERE
		t.is_deleted = 0
		AND t.is_active = 1`

	rows, err := t.Conn.QueryContext(ctx, query)
	var transactionStatus int
	if status == "pending" {
		transactionStatus = 0
	} else if status == "waitingApproval" {
		transactionStatus = 1
	} else if status == "confirm" {
		transactionStatus = 2
	}
	queryWithStatus := query + ` AND t.status = ?`
	rows, err = t.Conn.QueryContext(ctx, queryWithStatus, transactionStatus)

	if status == "failed" {
		transactionStatus = 3
		cancelledStatus := 4
		queryWithStatus = query + ` AND t.status IN (?,?)`
		rows, err = t.Conn.QueryContext(ctx, query, transactionStatus, cancelledStatus)
	}

	if status == "boarded" {
		transactionStatus = 1
		bookingStatus := 3
		queryWithStatus = query + ` AND t.status = ? AND b.status = ?`
		rows, err = t.Conn.QueryContext(ctx, queryWithStatus, transactionStatus, bookingStatus)
	}
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


func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
