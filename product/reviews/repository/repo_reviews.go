package repository

import (
	"context"
	"database/sql"
	"github.com/models"
	"github.com/product/reviews"
	"github.com/sirupsen/logrus"
	"strconv"
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
			&t.Name ,
			&t.UserId  ,
			&t.GuideReview ,
			&t.ActivitiesReview ,
			&t.ServiceReview ,
			&t.CleanlinessReview ,
			&t.ValueReview ,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (r reviewRepository) CountRating(ctx context.Context, rating int, expID string) (int, error) {
	query := `SELECT COUNT(*) as count FROM reviews r WHERE exp_id = ? AND is_deleted = 0 AND is_active = 1`
	if rating != 0 {
		query = query + ` AND r.values = ` + strconv.Itoa(rating)
	}
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

func (r reviewRepository) GetByExpId(ctx context.Context, expID, sortBy string, rating, limit, offset int) ([]*models.Review, error) {
	query := `select * from reviews r where exp_id = ? AND is_deleted = 0 AND is_active = 1`
	if rating != 0 {
		query = query + ` AND r.values = ` + strconv.Itoa(rating)
	}
	if sortBy != "" {
		if sortBy == "ratingup" {
			query = query + ` ORDER BY r.values DESC`
		} else if sortBy == "ratingdown" {
			query = query + ` ORDER BY r.values ASC`
		} else if sortBy == "latestdate" {
			query = query + ` ORDER BY created_date DESC`
		} else if sortBy == "oldestdate" {
			query = query + ` ORDER BY created_date ASC`
		}
	}
	query = query + ` LIMIT ? OFFSET ?`
	res, err := r.fetch(ctx, query, expID, limit, offset)
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
