package version_app

import (
	"github.com/models"
	"golang.org/x/net/context"
)

type Usecase interface {
	GetAllVersion(ctx context.Context,typeApp int)([]*models.VersionApp,error)
}
