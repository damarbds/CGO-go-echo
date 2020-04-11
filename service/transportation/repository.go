package transportation

import (
	"context"

	"github.com/models"
)

type Repository interface {
	Insert(ctx context.Context, transportation models.Transportation) (*string, error)
	Update(ctx context.Context, transportation models.Transportation) (*string, error)
	List(ctx context.Context) ([]*models.TimesOption, error)
}
