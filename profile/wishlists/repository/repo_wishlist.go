package repository

import (
	"context"
	"database/sql"
	guuid "github.com/google/uuid"
	"github.com/models"
	"github.com/profile/wishlists"
)

type wishListRepository struct {
	Conn *sql.DB
}

func NewWishListRepository(Conn *sql.DB) wishlists.Repository {
	return &wishListRepository{Conn:Conn}
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


