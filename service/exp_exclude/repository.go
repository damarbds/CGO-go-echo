package exp_exclude

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetByExpIdJoin(ctx context.Context,expId string)([]*models.ExperienceExcludeJoin ,error)
	Insert(ctx context.Context,a *models.ExperienceExclude)error
	Delete(ctx context.Context,expId string)error
}
