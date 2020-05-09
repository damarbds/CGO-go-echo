package repository

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/product/reviews"
	"github.com/sirupsen/logrus"
)

type reviewRepository struct {
	Conn *sql.DB
}

// NewReviewRepository will create an object that represent the exp_payment.repository interface
func NewReviewRepository(Conn *sql.DB) reviews.Repository {
	return &reviewRepository{Conn}
}

func (m *reviewRepository) Insert(ctx context.Context, a models.Review) (string, error) {
	a.Id = guuid.New().String()
	query := `INSERT reviews SET id=? , created_by=? , created_date=? , modified_by=?, modified_date=? ,
				deleted_by=? , deleted_date=? , is_deleted=? , is_active=? , reviews.values=?,reviews.desc=?,exp_id=?,user_id=?,
				guide_review=?,activities_review=?,service_review=?,cleanliness_review=?,value_review=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return "", err
	}
	_, err = stmt.ExecContext(ctx, a.Id, a.CreatedBy, time.Now(), nil, nil, nil, nil, 0, 1, a.Values, a.Desc,
		a.ExpId, a.UserId, a.GuideReview, a.ActivitiesReview, a.ServiceReview, a.CleanlinessReview, a.ValueReview)
	if err != nil {
		return "", err
	}

	//lastID, err := res.RowsAffected()
	if err != nil {
		return "", err
	}

	//a.Id = lastID
	return a.Id, nil
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
			&t.UserId,
			&t.GuideReview,
			&t.ActivitiesReview,
			&t.ServiceReview,
			&t.CleanlinessReview,
			&t.ValueReview,
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

func (r reviewRepository) GetByExpId(ctx context.Context, expID, sortBy string, rating, limit, offset int, userId string) ([]*models.Review, error) {
	query := `select * from reviews r where exp_id = ? AND is_deleted = 0 AND is_active = 1`
	if rating != 0 {
		query = query + ` AND r.values = ` + strconv.Itoa(rating)
	}
	if userId != "" {
		query = query + ` AND r.user_id = '` + userId + `' `
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
