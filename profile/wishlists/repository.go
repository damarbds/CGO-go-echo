package wishlists

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context, wl *models.Wishlist) (*models.Wishlist, error)
	List(ctx context.Context, userID string) ([]*models.WishlistObj, error)
}
