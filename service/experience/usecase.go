package experience

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetByID(ctx context.Context, id string) (*models.ExperienceDto, error)
	SearchExp(ctx context.Context, harborID, cityID string) ([]*models.ExpSearchObject, error)
	FilterSearchExp(ctx context.Context, cityID string,harborsId string,expTypeId string,startDate string,endDate string,guest string,trip string,bottomPrice string,upPrice string) ([]*models.ExpSearchObject, error)
	GetUserDiscoverPreference(ctx context.Context,page *int , size *int)([]*models.ExpUserDiscoverPreferenceDto,error)
}
