package admin

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	//Fetch(ctx context.Context, cursor string, num int64) (res []*models.Admin, nextCursor string, err error)
	GetByID(ctx context.Context, id string) (*models.Admin, error)
	GetByAdminEmail(ctx context.Context, merchantEmail string) (*models.Admin, error)
	Update(ctx context.Context, ar *models.Admin) error
	Insert(ctx context.Context, a *models.Admin) error
	Delete(ctx context.Context, id string, deleted_by string) error
}
