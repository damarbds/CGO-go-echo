package admin

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	Update(ctx context.Context, ar *models.NewCommandAdmin, user string) error
	Create(ctx context.Context, ar *models.NewCommandAdmin, user string) error
	Login(ctx context.Context, ar *models.Login) (*models.GetToken, error)
	ValidateTokenAdmin(ctx context.Context, token string) (*models.AdminDto, error)
	GetAdminInfo(ctx context.Context, token string) (*models.AdminDto, error)
}
