package user

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.User, nextCursor string, err error)
	GetByID(ctx context.Context, id string) (*models.User, error)
	GetByUserEmail(ctx context.Context, userEmail string) (*models.User, error)
	Update(ctx context.Context, ar *models.User) error
	Insert(ctx context.Context, a *models.User) error
	Delete(ctx context.Context, id string,deleted_by string) error
	GetCreditByID(ctx context.Context, id string) (int, error)
}