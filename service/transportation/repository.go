package transportation

import (
	"context"
	"github.com/models"
)

type Repository interface {
	List(ctx context.Context) ([]*models.TimesOption, error)
}
