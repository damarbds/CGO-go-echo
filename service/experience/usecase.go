package experience

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*models.ExperienceDto, error)
	SearchExp(ctx context.Context, harborID, cityID string) ([]*models.ExpSearchObject, error)
	GetUserDiscoverPreference(ctx context.Context,page *int , size *int)([]*models.ExpUserDiscoverPreferenceDto,error)
}
