package wishlists

import (
	"context"
	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context, wl *models.Wishlist) (*models.Wishlist, error)
	List(ctx context.Context, userID string,limit, offset int) ([]*models.WishlistObj, error)
	Count(ctx context.Context, userID string) (int, error)
	GetByUserAndExpId(ctx context.Context, userID string, expId string) ([]*models.WishlistObj, error)
}
