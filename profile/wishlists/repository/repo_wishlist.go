package repository

import (
	"context"
	"database/sql"
	"time"

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

func (m *wishListRepository) Count(ctx context.Context, userID string) (int, error) {
	query := `SELECT COUNT(*) as count FROM wishlists WHERE user_id = ? AND is_deleted = 0 AND is_active = 1`
	rows, err := m.Conn.QueryContext(ctx, query, userID)
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

func (w wishListRepository) List(ctx context.Context, userID string, limit, offset int,expId string) ([]*models.WishlistObj, error) {
	res := make([]*models.WishlistObj, 0)
	query := `SELECT * FROM wishlists WHERE user_id = ? AND is_deleted = 0 AND is_active = 1`
	if expId != ""{
		query = query + ` AND exp_id = '` + expId + `' `
	}
	if limit != 0 {
		query = query + ` LIMIT ? OFFSET ?`
		result, err := w.fetch(ctx, query, userID, limit, offset)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
		res = result
	} else {
		result, err := w.fetch(ctx, query, userID)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
		res = result
	}

	return res, nil
}

func (m *wishListRepository) GetByUserAndExpId(ctx context.Context, userID string, expId string,transId string) ([]*models.WishlistObj, error) {
	var res []*models.WishlistObj
	if expId != ""{
		query := `SELECT * FROM wishlists WHERE user_id = ? AND exp_id = ? AND is_deleted = 0 AND is_active = 1`

		result, err := m.fetch(ctx, query, userID, expId)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
		res = result
	}else if transId != ""{
		query := `SELECT * FROM wishlists WHERE user_id = ? AND trans_id = ? AND is_deleted = 0 AND is_active = 1`

		result, err := m.fetch(ctx, query, userID, transId)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, models.ErrNotFound
			}
			return nil, err
		}
		res = result
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
func (m *wishListRepository) DeleteByUserIdAndExpIdORTransId(ctx context.Context, userId string, expId string, transId string,deletedBy string) error {
	if expId != ""{
		query := `UPDATE wishlists SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE user_id =? AND exp_id=?`
		stmt, err := m.Conn.PrepareContext(ctx, query)
		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0,userId,expId)
		if err != nil {
			return err
		}

		//lastID, err := res.RowsAffected()
		if err != nil {
			return err
		}

	}else if transId != ""{
		query := `UPDATE wishlists SET deleted_by=? , deleted_date=? , is_deleted=? , is_active=? WHERE user_id =? AND trans_id =?`
		stmt, err := m.Conn.PrepareContext(ctx, query)
		if err != nil {
			return err
		}

		_, err = stmt.ExecContext(ctx, deletedBy, time.Now(), 1, 0,userId,transId)
		if err != nil {
			return err
		}

		//lastID, err := res.RowsAffected()
		if err != nil {
			return err
		}
	}

	//a.Id = lastID
	return nil
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
