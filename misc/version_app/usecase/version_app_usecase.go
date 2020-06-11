package usecase

import (
	"github.com/misc/version_app"
	"github.com/models"
	"golang.org/x/net/context"
	"time"
)

type versionAPPUsecase struct {
	versionAPPRepo       version_app.Repository
	contextTimeout  time.Duration
}

func NewVersionAPPUsecase(versionAPPRepo version_app.Repository, timeout time.Duration) version_app.Usecase {
	return &versionAPPUsecase{
		versionAPPRepo:       versionAPPRepo,
		contextTimeout:  timeout,
	}
}

func (v versionAPPUsecase) GetAllVersion(ctx context.Context, typeApp int) ([]*models.VersionApp, error) {
	ctx, cancel := context.WithTimeout(ctx, v.contextTimeout)
	defer cancel()

	res, err := v.versionAPPRepo.GetAll(ctx,typeApp)
	if err != nil {
		return nil, err
	}
	return res,nil

}
