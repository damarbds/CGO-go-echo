package time_options

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Repository interface {
	//Insert(ctx context.Context, a models.TimesOption) (*int, error)
	GetByTime(ctx context.Context, time string) (*models.TimesOption, error)
	TimeOptions(ctx context.Context) ([]*models.TimesOption, error)
}
