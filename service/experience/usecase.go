package experience

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	UpdateStatus(ctx context.Context,status int ,id string,token string)(*models.NewCommandChangeStatus,error)
	PublishExperience(ctx context.Context, commandExperience models.NewCommandExperience, token string) (*models.ResponseCreateExperience, error)
	CreateExperience(ctx context.Context, commandExperience models.NewCommandExperience, token string) (*models.ResponseCreateExperience, error)
	UpdateExperience(ctx context.Context, commandExperience models.NewCommandExperience, token string) (*models.ResponseCreateExperience, error)
	GetByID(ctx context.Context, id string,currency string) (*models.ExperienceDto, error)
	SearchExp(ctx context.Context, harborID, cityID string) ([]*models.ExpSearchObject, error)
	FilterSearchExp(ctx context.Context, isMerchant bool, search, token, status, cityID string, harborsId string, expTypeId string, startDate string, endDate string, guest string, trip string, bottomPrice string, upPrice string, sortBy string, page, limit, offset int,provinceId string,currency string) (*models.FilterSearchWithPagination, error)
	GetUserDiscoverPreference(ctx context.Context, page *int, size *int,currency string) ([]*models.ExpUserDiscoverPreferenceDto, error)
	GetExpTypes(ctx context.Context) ([]*models.ExpTypeObject, error)
	GetExpInspirations(ctx context.Context) ([]*models.ExpInspirationDto, error)
	GetByCategoryID(ctx context.Context, categoryId int) ([]*models.ExpSearchObject, error)
	GetSuccessBookCount(ctx context.Context, token string) (*models.Count, error)
	GetPublishedExpCount(ctx context.Context, token string) (*models.Count, error)
	GetExpPendingTransactionCount(ctx context.Context, token string) (*models.Count, error)
	GetExpFailedTransactionCount(ctx context.Context, token string) (*models.Count, error)
}
