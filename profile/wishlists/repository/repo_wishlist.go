package repository

import (
	"context"
	"database/sql"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/profile/wishlists"
	"github.com/sirupsen/logrus"
)

type wishListRepository struct {
	Conn *sql.DB
}

func NewWishListRepository(Conn *sql.DB) wishlists.Repository {
	return &wishListRepository{Conn: Conn}
}

func (m *wishListRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.WishlistObj, error) {
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

	result := make([]*models.WishlistObj, 0)
	for rows.Next() {
		t := new(models.WishlistObj)
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
			&t.TransId,
			&t.ExpId,
			&t.UserId,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (w wishListRepository) List(ctx context.Context, userID string) ([]*models.WishlistObj, error) {
	query := `SELECT * FROM wishlists WHERE user_id = ? AND is_deleted = 0 AND is_active = 1`

	res, err := w.fetch(ctx, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return res, nil
}

func (m *wishListRepository) GetByUserAndExpId(ctx context.Context, userID string, expId string) ([]*models.WishlistObj, error) {
	query := `SELECT * FROM wishlists WHERE user_id = ? AND exp_id = ? AND is_deleted = 0 AND is_active = 1`

	res, err := m.fetch(ctx, query, userID, expId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrNotFound
		}
		return nil, err
	}

	return res, nil
}
func (w wishListRepository) Insert(ctx context.Context, wl *models.Wishlist) (*models.Wishlist, error) {
	id := guuid.New()
	wl.Id = id.String()

	q := `INSERT wishlists SET id = ?, created_by = ?, created_date = ?, modified_by = ?, 
	modified_date = ?, deleted_by = ?, deleted_date = ?, is_deleted = ?, is_active = ?, user_id = ?, `

	if wl.TransId != "" {
		q = q + `trans_id = ?`
		stmt, err := w.Conn.PrepareContext(ctx, q)
		if err != nil {
			return nil, err
		}
		_, err = stmt.ExecContext(
			ctx,
			wl.Id,
			wl.CreatedBy,
			wl.CreatedDate,
			wl.ModifiedBy,
			wl.ModifiedDate,
			wl.DeletedBy,
			wl.DeletedDate,
			wl.IsDeleted,
			wl.IsActive,
			wl.UserId,
			wl.TransId,
		)
		if err != nil {
			return nil, err
		}
	}

	if wl.ExpId != "" {
		q = q + `exp_id = ?`
		stmt, err := w.Conn.PrepareContext(ctx, q)
		if err != nil {
			return nil, err
		}
		_, err = stmt.ExecContext(
			ctx,
			wl.Id,
			wl.CreatedBy,
			wl.CreatedDate,
			wl.ModifiedBy,
			wl.ModifiedDate,
			wl.DeletedBy,
			wl.DeletedDate,
			wl.IsDeleted,
			wl.IsActive,
			wl.UserId,
			wl.ExpId,
		)
		if err != nil {
			return nil, err
		}
	}

	return wl, nil
}
