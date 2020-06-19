package exp_Include

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByExpIdJoin(ctx context.Context,expId string)([]*models.ExperienceIncludeJoin ,error)
	Insert(ctx context.Context,a *models.ExperienceInclude)error
	Delete(ctx context.Context,expId string)error
}
