package repository

import (
	"context"
	"database/sql"
	"github.com/service/reviews"
	"github.com/sirupsen/logrus"
)

type reviewRepository struct {
	Conn *sql.DB
}

// NewReviewRepository will create an object that represent the exp_payment.Repository interface
func NewReviewRepository(Conn *sql.DB) reviews.Repository {
	return &reviewRepository{Conn}
}

func (r reviewRepository) CountRating(ctx context.Context, expID string) (int, error) {
	query := `SELECT COUNT(*) as count FROM reviews WHERE exp_id = ?`

	rows, err := r.Conn.QueryContext(ctx, query, expID)
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
