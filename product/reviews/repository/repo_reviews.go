package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/product/reviews"
	"github.com/sirupsen/logrus"
)

type reviewRepository struct {
	Conn *sql.DB
}

// NewReviewRepository will create an object that represent the exp_payment.Repository interface
func NewReviewRepository(Conn *sql.DB) reviews.Repository {
	return &reviewRepository{Conn}
}

func (m *reviewRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Review, error) {
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

	result := make([]*models.Review, 0)
	for rows.Next() {
		t := new(models.Review)
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
			&t.Values,
			&t.Desc,
			&t.ExpId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
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

func (r reviewRepository) GetByExpId(ctx context.Context, expID string) ([]*models.Review, error) {
	query := `select * from reviews where exp_id = ? AND is_deleted = 0 AND is_active = 1`

	res, err := r.fetch(ctx, query, expID)
	if err != nil {
		return nil, err
	}
	return res, err
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
