package experience

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	CreateExperience(ctx context.Context,commandExperience models.NewCommandExperience,token string)(*models.ResponseCreateExperience,error)
	GetByID(ctx context.Context, id string) (*models.ExperienceDto, error)
	SearchExp(ctx context.Context, harborID, cityID string) ([]*models.ExpSearchObject, error)
	FilterSearchExp(ctx context.Context, cityID string, harborsId string, expTypeId string, startDate string, endDate string, guest string, trip string, bottomPrice string, upPrice string) ([]*models.ExpSearchObject, error)
	GetUserDiscoverPreference(ctx context.Context, page *int, size *int) ([]*models.ExpUserDiscoverPreferenceDto, error)
	GetExpTypes(ctx context.Context) ([]*models.ExpTypeObject, error)
	GetExpInspirations(ctx context.Context) ([]*models.ExpInspirationObject, error)
	GetByCategoryID(ctx context.Context, categoryId int) ([]*models.ExpSearchObject, error)
	GetSuccessBookCount(ctx context.Context, token string) (*models.Count, error)
	GetExpCount(ctx context.Context, token string) (*models.Count, error)
	GetExpPendingTransactionCount(ctx context.Context, token string) (*models.Count, error)
	GetExpFailedTransactionCount(ctx context.Context, token string) (*models.Count, error)
}
