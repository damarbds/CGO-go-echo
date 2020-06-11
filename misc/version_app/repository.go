package version_app

import (
	"context"
	"github.com/models"
)

type Repository interface {
	GetAll(ctx context.Context,typeApp int)([]*models.VersionApp,error)
}