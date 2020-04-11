package repository

import (
	"context"
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/transactions/transaction"
)

type transactionRepository struct {
	Conn *sql.DB
}

func NewTransactionRepository(Conn *sql.DB) transaction.Repository {
	return &transactionRepository{Conn: Conn}
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

func checkCount(rows *sql.Rows) (count int, err error) {
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
